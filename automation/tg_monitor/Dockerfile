FROM python:3.11-slim-bookworm AS builder
WORKDIR /build
COPY requirements.txt .
RUN pip install --no-cache-dir --target=/build/deps -r requirements.txt

FROM gcr.io/distroless/python3-debian12
WORKDIR /app
ENV PYTHONPATH=/app/deps

COPY --from=builder /build/deps /app/deps
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY sleep_bot.py /app/

USER 65532:65532
ENTRYPOINT ["python3", "sleep_bot.py", "-u"]
