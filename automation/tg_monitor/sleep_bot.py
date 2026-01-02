import asyncio
import logging
import json
import os
import requests
from datetime import datetime, timezone, timedelta

# Import ZoneInfo for timezone support
from zoneinfo import ZoneInfo
from telethon import TelegramClient
from telethon.tl.types import UserStatusOffline, UserStatusOnline

# ================= CONFIGURATION =================
# Telegram API Credentials
API_ID = int(os.environ["TG_DEV_APP_ID"])
API_HASH = os.environ["TG_DEV_API_HASH"]

# Paths
DATA_DIR = "/app/data"
SESSION_NAME = os.path.join(DATA_DIR, "sleep_tracker_session")
STATE_FILE = os.path.join(DATA_DIR, "sleep_state.json")

# Target
TARGET_USER = "@" + os.environ["TG_USERNAME"]

# API Configuration
API_URL = os.environ["REPORT_API_URL"]
API_KEY = os.environ["REPORT_API_KEY"]
HABIT_ID = os.environ["REPORT_API_HABITID"]

# Logic Settings
# Get Timezone from Env, default to UTC if not set
TZ_NAME = os.getenv("TZ", "UTC")
try:
    LOCAL_TZ = ZoneInfo(TZ_NAME)
except Exception:
    print(f"Timezone {TZ_NAME} not found, falling back to UTC")
    LOCAL_TZ = timezone.utc

CHECK_INTERVAL = int(os.getenv("CHECK_INTERVAL", "1800"))
SLEEP_THRESHOLD_HOURS = float(os.getenv("SLEEP_THRESHOLD_HOURS", "3.0"))
SAFE_WAKE_HOUR = int(os.getenv("SAFE_WAKE_HOUR", "14"))  # 2 PM Local Time
# =================================================

logging.basicConfig(
    format="%(asctime)s - %(levelname)s - %(message)s",
    level=logging.INFO,
    handlers=[logging.StreamHandler()],  # Distroless logs to stdout/docker logs
)
logger = logging.getLogger(__name__)

client = TelegramClient(SESSION_NAME, API_ID, API_HASH)


def get_now_local():
    """Returns current time in the configured timezone."""
    return datetime.now(LOCAL_TZ)


def load_state():
    if os.path.exists(STATE_FILE):
        try:
            with open(STATE_FILE, "r") as f:
                return json.load(f)
        except Exception:
            pass
    # Default state
    return {
        "last_log_date": (get_now_local() - timedelta(days=1)).strftime("%Y-%m-%d"),
        "last_seen_ts": datetime.now(timezone.utc).timestamp(),
    }


def save_state(state):
    with open(STATE_FILE, "w") as f:
        json.dump(state, f)


def send_api_request(date_str, value):
    headers = {"X-API-Key": API_KEY, "Content-Type": "application/json"}
    payload = {"habitId": HABIT_ID, "date": date_str, "value": value}

    try:
        logger.info(f"🚀 Sending API Request: Date={date_str}, Value={value}")
        resp = requests.post(API_URL, json=payload, headers=headers, timeout=10)
        if resp.status_code in [200, 201]:
            logger.info("✅ API Success")
            return True
        else:
            logger.error(f"❌ API Fail: {resp.status_code} - {resp.text}")
            return False
    except Exception as e:
        logger.error(f"❌ API Error: {e}")
        return False


async def main_loop():
    logger.info(f"🤖 Sleep Tracker Bot Started in Timezone: {TZ_NAME}")
    state = load_state()
    potential_sleep_start = None

    while True:
        try:
            entity = await client.get_entity(TARGET_USER)
            status = entity.status

            # Get current local time
            now_local = get_now_local()

            is_online = isinstance(status, UserStatusOnline)
            is_offline = isinstance(status, UserStatusOffline)

            # Determine Last Seen (Convert to Local Time)
            last_seen_local = None

            if is_offline:
                # Telethon gives naive UTC. Make it aware, then convert to local.
                last_seen_utc = status.was_online.replace(tzinfo=timezone.utc)
                last_seen_local = last_seen_utc.astimezone(LOCAL_TZ)
            elif is_online:
                last_seen_local = now_local

            if not (is_online or is_offline):
                logger.warning("⚠️ User status hidden/unknown. Waiting...")
                await asyncio.sleep(CHECK_INTERVAL)
                continue

            state["last_seen_ts"] = last_seen_local.timestamp()
            save_state(state)

            # --- SLEEP DETECTION ---
            if is_online:
                if potential_sleep_start:
                    # diff calculation works regardless of timezone as long as both are aware
                    duration = (now_local - potential_sleep_start).total_seconds() / 3600

                    if duration >= SLEEP_THRESHOLD_HOURS:
                        logger.info(f"🛌 Sleep Detected! Duration: {duration:.2f}h")
                        sleep_date_str = potential_sleep_start.strftime("%Y-%m-%d")

                        if send_api_request(sleep_date_str, round(duration, 1)):
                            state["last_log_date"] = sleep_date_str
                            save_state(state)
                    else:
                        logger.info(f"👀 User back online. Inactive {duration:.2f}h")

                    potential_sleep_start = None
                else:
                    logger.info("🟢 User is Online.")

            elif is_offline:
                time_diff = now_local - last_seen_local
                hours_offline = time_diff.total_seconds() / 3600

                logger.info(f"🔴 Offline. Last seen: {last_seen_local.strftime('%H:%M')} ({hours_offline:.2f}h ago)")

                if hours_offline >= SLEEP_THRESHOLD_HOURS:
                    if not potential_sleep_start:
                        logger.info("💤 Potential sleep detected...")
                        potential_sleep_start = last_seen_local
                    else:
                        potential_sleep_start = last_seen_local

            # --- "0" VALUE (INSOMNIA) LOGIC ---
            # Parse last log date (assume it was local date)
            last_log_dt = datetime.strptime(state["last_log_date"], "%Y-%m-%d").replace(tzinfo=LOCAL_TZ)
            expected_next_log_date = last_log_dt + timedelta(days=1)

            # Compare Local Date vs Expected Date
            if now_local.date() > expected_next_log_date.date():
                # Compare Local Hour vs Safe Wake Hour
                if now_local.hour >= SAFE_WAKE_HOUR:
                    missing_date_str = expected_next_log_date.strftime("%Y-%m-%d")
                    logger.info(f"🚫 No sleep detected for {missing_date_str}. It is past {SAFE_WAKE_HOUR}:00. Logging 0.")

                    if send_api_request(missing_date_str, 0):
                        state["last_log_date"] = missing_date_str
                        save_state(state)

        except Exception as e:
            logger.error(f"⚠️ Error in loop: {e}")

        await asyncio.sleep(CHECK_INTERVAL)


def run_api_sanity_check():
    from urllib.parse import urlparse

    try:
        parsed = urlparse(API_URL)
        sanity_url = f"{parsed.scheme}://{parsed.netloc}/api/me"

        headers = {"X-API-Key": API_KEY, "Content-Type": "application/json"}

        logger.info(f"🔎 Running API Sanity Check against: {sanity_url}")

        response = requests.get(sanity_url, headers=headers, timeout=10)

        if response.status_code == 200:
            data = response.json()
            logger.info(f"✅ API Sanity Check Passed! Connected as: {data.get('userId', 'Unknown')}")
        else:
            logger.error(f"❌ API Sanity Check Failed! Status: {response.status_code}. Response: {response.text}")

    except Exception as e:
        logger.error(f"❌ API Sanity Check Error: {e}")


if __name__ == "__main__":
    with client:
        run_api_sanity_check()
        client.loop.run_until_complete(main_loop())
