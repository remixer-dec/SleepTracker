<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content slide-up">
      <h2 class="modal-title">{{ t('entry.edit') }}</h2>
      <p class="date-label">{{ formattedDate }}</p>

      <div class="widget-container">
        <BooleanWidget
          v-if="habit.type === 'boolean'"
          v-model="form.value"
          :positiveLabel="habit.positiveValue || t('entry.yes')"
        />

        <ReachableWidget
          v-else-if="habit.type === 'reachable'"
          v-model="form.value"
          :goal="habit.goal || 8"
          :unit="t('entry.hours')"
        />

        <NumberWidget
          v-else-if="habit.type === 'number'"
          v-model="form.value"
        />

        <PercentageWidget
          v-else-if="habit.type === 'percentage'"
          v-model="form.value"
        />

        <AdditiveWidget
          v-else-if="habit.type === 'additive'"
          v-model="form.value"
          :goal="habit.goal || 10"
          :accumulated="accumulatedProgress"
          :saveProgress="habit.saveProgress"
        />

        <MoodWidget
          v-else-if="habit.type === 'mood'"
          v-model="form.value"
        />

        <RatingWidget
          v-else-if="habit.type === 'rating'"
          v-model="form.value"
        />
      </div>

      <div class="form-group">
        <label class="form-label">{{ t('entry.note') }}</label>
        <textarea
          v-model="form.note"
          class="form-input note-input"
          rows="2"
        ></textarea>
      </div>

      <div class="form-actions">
        <button v-if="entry" class="btn btn-danger" @click="confirmDelete">
          {{ t('entry.delete') }}
        </button>
        <div class="spacer"></div>
        <button class="btn btn-secondary" @click="$emit('close')">
          {{ t('entry.cancel') }}
        </button>
        <button class="btn btn-primary" @click="save">
          {{ t('entry.save') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useHabitsStore } from '../stores/habits.js'
import { calculateAccumulatedProgress } from '../utils/calculations.js'
import BooleanWidget from './widgets/BooleanWidget.vue'
import ReachableWidget from './widgets/ReachableWidget.vue'
import NumberWidget from './widgets/NumberWidget.vue'
import PercentageWidget from './widgets/PercentageWidget.vue'
import AdditiveWidget from './widgets/AdditiveWidget.vue'
import MoodWidget from './widgets/MoodWidget.vue'
import RatingWidget from './widgets/RatingWidget.vue'

const { t, locale } = useI18n()
const habitsStore = useHabitsStore()

const props = defineProps({
  habit: {
    type: Object,
    required: true
  },
  entry: {
    type: Object,
    default: null
  },
  date: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'save', 'delete'])

const form = ref({
  value: 0,
  note: ''
})

const formattedDate = computed(() => {
  const [year, month, day] = props.date.split('-').map(Number)
  const date = new Date(year, month - 1, day)
  return date.toLocaleDateString(locale.value, {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const accumulatedProgress = computed(() => {
  if (props.habit.type !== 'additive' || !props.habit.saveProgress) return 0
  const entries = habitsStore.entries[props.habit.id] || []
  return calculateAccumulatedProgress(entries, props.habit, props.date)
})

function getDefaultValue() {
  switch (props.habit.type) {
    case 'boolean': return 0
    case 'reachable': return props.habit.goal || 8
    case 'number': return 0
    case 'percentage': return 50
    case 'additive': return 0
    case 'mood': return 1
    case 'rating': return 3
    default: return 0
  }
}

watch([() => props.entry, () => props.date], () => {
  if (props.entry) {
    form.value = {
      value: props.entry.value,
      note: props.entry.note || ''
    }
  } else {
    form.value = {
      value: getDefaultValue(),
      note: ''
    }
  }
}, { immediate: true })

function save() {
  emit('save', { ...form.value })
}

function confirmDelete() {
  emit('delete')
}
</script>

<style scoped>
.date-label {
  color: var(--color-text-muted);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-md);
  text-transform: capitalize;
}

.widget-container {
  margin-bottom: var(--spacing-md);
}

.note-input {
  resize: none;
}

.spacer {
  flex: 1;
}
</style>
