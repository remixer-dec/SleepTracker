<template>
  <div class="heatmap-container">
    <div class="heatmap-header">
      <button class="btn-icon" @click="prevMonth">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path d="M12 4l-6 6 6 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
        </svg>
      </button>
      <span class="month-label">{{ monthLabel }}</span>
      <button class="btn-icon" @click="nextMonth" :disabled="isCurrentMonth">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path d="M8 4l6 6-6 6" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
        </svg>
      </button>
    </div>

    <div
      class="heatmap-grid-container"
      ref="gridContainer"
      @touchstart="handleTouchStart"
      @touchmove="handleTouchMove"
      @touchend="handleTouchEnd"
    >
      <div class="weekday-labels">
        <span v-for="day in weekdayLabels" :key="day">{{ day }}</span>
      </div>

      <div class="heatmap-grid" :style="gridStyle">
        <div
          v-for="day in displayDays"
          :key="day.date"
          class="heatmap-cell"
          :class="{
            'empty': !day.inMonth,
            'today': day.isToday,
            'has-entry': day.hasEntry
          }"
          :style="{ backgroundColor: day.color }"
          :title="day.hasEntry ? formatTooltip(day) : ''"
          @click="handleCellClick(day)"
        >
          <span class="cell-day" v-if="day.inMonth">{{ day.dayNumber }}</span>
        </div>
      </div>
    </div>

    <p class="swipe-hint text-dim text-xs text-center mt-sm">{{ t('actions.swipeHint') }}</p>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatDate, getValueColor } from '../utils/calculations.js'

const { t, locale } = useI18n()

const props = defineProps({
  habit: {
    type: Object,
    required: true
  },
  entries: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['select-date'])

const currentDate = ref(new Date())
const touchStartX = ref(0)
const touchEndX = ref(0)
const gridContainer = ref(null)

const monthLabel = computed(() => {
  const options = { year: 'numeric', month: 'long' }
  return currentDate.value.toLocaleDateString(locale.value, options)
})

const isCurrentMonth = computed(() => {
  const now = new Date()
  return currentDate.value.getMonth() === now.getMonth() &&
         currentDate.value.getFullYear() === now.getFullYear()
})

const weekdayLabels = computed(() => {
  const days = []
  const baseDate = new Date(2024, 0, 1) // Monday, January 1, 2024
  for (let i = 0; i < 7; i++) {
    const date = new Date(baseDate)
    date.setDate(date.getDate() + i)
    days.push(date.toLocaleDateString(locale.value, { weekday: 'narrow' }))
  }
  return days
})

const entriesByDate = computed(() => {
  const map = new Map()
  props.entries.forEach(entry => {
    map.set(entry.date, entry)
  })
  return map
})

const displayDays = computed(() => {
  const days = []
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  const today = formatDate(new Date())

  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startOffset = (firstDay.getDay() + 6) % 7

  for (let i = 0; i < startOffset; i++) {
    days.push({ date: '', inMonth: false, color: 'transparent' })
  }

  for (let d = 1; d <= lastDay.getDate(); d++) {
    const date = formatDate(new Date(year, month, d))
    const entry = entriesByDate.value.get(date)
    const hasEntry = !!entry

    days.push({
      date,
      dayNumber: d,
      inMonth: true,
      isToday: date === today,
      hasEntry,
      value: entry?.value,
      color: hasEntry ? getValueColor(entry.value, props.habit) : 'var(--color-heatmap-empty)'
    })
  }

  const remaining = 7 - (days.length % 7)
  if (remaining < 7) {
    for (let i = 0; i < remaining; i++) {
      days.push({ date: '', inMonth: false, color: 'transparent' })
    }
  }

  return days
})

const gridStyle = computed(() => ({
  gridTemplateRows: `repeat(${Math.ceil(displayDays.value.length / 7)}, 1fr)`
}))

function prevMonth() {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() - 1,
    1
  )
}

function nextMonth() {
  if (isCurrentMonth.value) return
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() + 1,
    1
  )
}

function handleCellClick(day) {
  if (day.inMonth) {
    emit('select-date', day.date)
  }
}

function formatTooltip(day) {
  if (!day.hasEntry) return ''
  const value = day.value
  if (props.habit.type === 'boolean') {
    return value ? t('entry.yes') : t('entry.no')
  }
  if (props.habit.type === 'mood') {
    const moods = ['😢', '😐', '😄']
    return moods[Math.min(Math.max(value - 1, 0), 2)]
  }
  if (props.habit.type === 'rating') {
    return '★'.repeat(value)
  }
  return String(value)
}

function handleTouchStart(e) {
  touchStartX.value = e.touches[0].clientX
}

function handleTouchMove(e) {
  touchEndX.value = e.touches[0].clientX
}

function handleTouchEnd() {
  const diff = touchStartX.value - touchEndX.value
  const threshold = 50

  if (Math.abs(diff) > threshold) {
    if (diff > 0) {
      nextMonth()
    } else {
      prevMonth()
    }
  }

  touchStartX.value = 0
  touchEndX.value = 0
}

watch(() => props.habit?.id, () => {
  currentDate.value = new Date()
})
</script>

<style scoped>
.heatmap-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.heatmap-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--spacing-sm) 0;
}

.month-label {
  font-size: var(--font-size-md);
  font-weight: 500;
  color: var(--color-text);
  text-transform: capitalize;
}

.heatmap-grid-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  touch-action: pan-y;
  user-select: none;
}

.weekday-labels {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: var(--spacing-xs);
}

.weekday-labels span {
  text-align: center;
  font-size: var(--font-size-xs);
  color: var(--color-text-dim);
  text-transform: uppercase;
}

.heatmap-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  flex: 1;
}

.heatmap-cell {
  aspect-ratio: 1;
  border-radius: var(--border-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: transform 0.1s, box-shadow 0.1s;
  position: relative;
}

.heatmap-cell:not(.empty):hover {
  transform: scale(1.05);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.heatmap-cell.empty {
  cursor: default;
  background: transparent !important;
}

.heatmap-cell.today {
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.6);
}

.cell-day {
  font-size: var(--font-size-xs);
  color: var(--color-text);
  opacity: 0.8;
}

.heatmap-cell.has-entry .cell-day {
  font-weight: 600;
  opacity: 1;
}

.swipe-hint {
  opacity: 0.5;
}

@media (min-width: 768px) {
  .swipe-hint {
    display: none;
  }
}
</style>
