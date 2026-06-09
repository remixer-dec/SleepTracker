<template>
  <div class="app-container">
    <div v-if="showAntiStreakBackground" class="rage-background">
      <RageEmoji />
    </div>

    <HabitSelector
      v-if="!authStore.isLoading"
      @add-habit="openHabitModal(null)"
    />

    <div class="main-content" v-if="!authStore.isLoading && habitsStore.selectedHabit">
      <HabitTitle :habit="habitsStore.selectedHabit" @edit="openHabitModal(habitsStore.selectedHabit)" />
      <StreakDisplay
        v-if="habitsStore.selectedHabit.trackStreak !== false"
        :streak="currentStreak"
        :record="longestStreak"
        :weeklyStreaks="weeklyStreaks"
        :level="levelInfo"
      />
      <HeatmapCalendar
        :habit="habitsStore.selectedHabit"
        :entries="habitsStore.currentEntries"
        @select-date="openEntryModal"
        @month-change="handleMonthChange"
      />
      <BalanceWidget
        v-if="habitsStore.selectedHabit.goal"
        :habit="habitsStore.selectedHabit"
        :entries="habitsStore.currentEntries"
        :month-date="selectedMonth"
      />
    </div>

    <div v-else-if="!authStore.isLoading && habitsStore.habits.length === 0" class="empty-state">
      <p class="text-muted" v-cloak>{{ t('habits.noHabits') }}</p>
      <button v-if="authStore.isOwner" class="btn btn-primary mt-md" @click="openHabitModal(null)">
        {{ t('habits.addFirst') }}
      </button>
    </div>

    <div v-if="authStore.isLoading" class="loading-state">
      <div class="loader"></div>
    </div>

    <HabitModal
      v-if="showHabitModal"
      :habit="editingHabit"
      @close="showHabitModal = false"
      @save="saveHabit"
      @delete="deleteHabit"
    />

    <EntryModal
      v-if="showEntryModal"
      :habit="habitsStore.selectedHabit"
      :entry="editingEntry"
      :date="selectedDate"
      @close="showEntryModal = false"
      @save="saveEntry"
      @delete="deleteEntry"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from './stores/auth.js'
import { useHabitsStore } from './stores/habits.js'
import { useAchievementsStore } from './stores/achievements.js'
import { calculateStreak, calculateWeeklyStreaks, calculateLevel, calculateStats, formatDate } from './utils/calculations.js'
import HabitSelector from './components/HabitSelector.vue'
import HabitTitle from './components/HabitTitle.vue'
import StreakDisplay from './components/StreakDisplay.vue'
import HeatmapCalendar from './components/HeatmapCalendar.vue'
import BalanceWidget from './components/BalanceWidget.vue'
import HabitModal from './components/HabitModal.vue'
import EntryModal from './components/EntryModal.vue'
import RageEmoji from './components/RageEmoji.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const habitsStore = useHabitsStore()
const achievementsStore = useAchievementsStore()

const showHabitModal = ref(false)
const showEntryModal = ref(false)
const editingHabit = ref(null)
const editingEntry = ref(null)
const selectedDate = ref('')
const selectedMonth = ref(new Date())

const streakInfo = computed(() => {
  if (!habitsStore.selectedHabit || !habitsStore.currentEntries) {
    return { current: 0, longest: 0, isAntiStreak: false }
  }
  return calculateStreak(habitsStore.currentEntries, habitsStore.selectedHabit)
})

const currentStreak = computed(() => streakInfo.value.current)
const longestStreak = computed(() => streakInfo.value.longest)
const showAntiStreakBackground = computed(() => streakInfo.value.isAntiStreak)

const weeklyStreaks = computed(() => {
  if (!habitsStore.selectedHabit || !habitsStore.currentEntries) return 0
  return calculateWeeklyStreaks(habitsStore.currentEntries, habitsStore.selectedHabit)
})

const levelInfo = computed(() => {
  const allEntries = Object.values(habitsStore.entries)
  const stats = calculateStats(allEntries, habitsStore.habits, achievementsStore.unlockedAchievements)
  return calculateLevel(stats)
})

function openHabitModal(habit) {
  if (!authStore.isOwner) return
  editingHabit.value = habit
  showHabitModal.value = true
}

function openEntryModal(date) {
  if (!authStore.isOwner) return
  selectedDate.value = date
  editingEntry.value = habitsStore.getEntryByDate(habitsStore.selectedHabitId, date)
  showEntryModal.value = true
}

function handleMonthChange(monthDate) {
  selectedMonth.value = monthDate
}

async function saveHabit(habitData) {
  try {
    if (habitData.notificationsOn && 'Notification' in window && Notification.permission === 'default') {
      await Notification.requestPermission()
    }

    if (editingHabit.value) {
      await habitsStore.updateHabit(editingHabit.value.id, habitData)
    } else {
      await habitsStore.createHabit(habitData)
    }
    showHabitModal.value = false
    checkAchievements()
  } catch (error) {
    console.error('Failed to save habit:', error)
  }
}

async function deleteHabit() {
  if (!editingHabit.value) return
  try {
    await habitsStore.deleteHabit(editingHabit.value.id)
    showHabitModal.value = false
  } catch (error) {
    console.error('Failed to delete habit:', error)
  }
}

async function saveEntry(entryData) {
  try {
    await habitsStore.createOrUpdateEntry({
      habitId: habitsStore.selectedHabitId,
      date: selectedDate.value,
      ...entryData
    })
    showEntryModal.value = false
    checkAchievements()
  } catch (error) {
    console.error('Failed to save entry:', error)
  }
}

async function deleteEntry() {
  if (!editingEntry.value) return
  try {
    await habitsStore.deleteEntry(editingEntry.value.id, habitsStore.selectedHabitId)
    showEntryModal.value = false
  } catch (error) {
    console.error('Failed to delete entry:', error)
  }
}

async function checkAchievements() {
  const allEntries = Object.values(habitsStore.entries)
  const stats = calculateStats(allEntries, habitsStore.habits, achievementsStore.unlockedAchievements)
  await achievementsStore.checkAndUnlockAchievements(stats)
}

const notifiedToday = ref(new Set())

function checkNotifications() {
  if (!('Notification' in window) || Notification.permission !== 'granted') return
  if (!habitsStore.habits || habitsStore.habits.length === 0) return

  const now = new Date()
  const currentTime = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`
  const today = formatDate(now)

  habitsStore.habits.forEach(habit => {
    if (!habit.notificationsOn || !habit.notificationTime) return
    if (notifiedToday.value.has(`${habit.id}-${today}`)) return

    if (habit.notificationTime === currentTime) {
      sendNotification(habit.name, t('habits.reminderMessage') || `Don't forget to log your ${habit.name}!`)
      notifiedToday.value.add(`${habit.id}-${today}`)
    }
  })
}

async function setupNotifications() {
  if (!('Notification' in window)) return

  // Request permission if any habit has notifications enabled
  const hasNotifications = habitsStore.habits.some(h => h.notificationsOn)
  if (hasNotifications && Notification.permission === 'default') {
    await Notification.requestPermission()
  }

  // Run check immediately then every minute
  checkNotifications()
  setInterval(checkNotifications, 60000)
}

function sendNotification(title, body) {
  if (Notification.permission === 'granted') {
    new Notification(title, { body, icon: '/favicon.svg' })
  }
}

watch(() => habitsStore.selectedHabitId, async (newId) => {
  if (newId) {
    const today = new Date()
    const threeMonthsAgo = new Date(today.getFullYear(), today.getMonth() - 3, 1)
    await habitsStore.fetchEntries(newId, formatDate(threeMonthsAgo), formatDate(today))
  }
})

onMounted(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const joinToken = urlParams.get('token')

  if (joinToken) {
    const success = await authStore.join(joinToken)
    if (success) {
      window.history.replaceState({}, document.title, window.location.pathname)
    }
  } else {
    await authStore.checkAuth()
  }

  await habitsStore.fetchHabits()
  await achievementsStore.fetchAchievements()

  if (habitsStore.selectedHabitId) {
    const today = new Date()
    const threeMonthsAgo = new Date(today.getFullYear(), today.getMonth() - 3, 1)
    await habitsStore.fetchEntries(
      habitsStore.selectedHabitId,
      formatDate(threeMonthsAgo),
      formatDate(today)
    )
  }

  setupNotifications()
})
</script>

<style scoped>
v-cloak {
  display: none;
}
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
  min-height: 0;
  overflow: hidden;
}

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.loading-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loader {
  width: 40px;
  height: 40px;
  border: 3px solid var(--loader-track-color);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.rage-background {
  position: fixed;
  inset: 0;
  background: radial-gradient(circle at center, var(--rage-glow) 0%, transparent 70%);
  pointer-events: none;
  z-index: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.4;
}
</style>
