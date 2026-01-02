<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content slide-up">
      <h2 class="modal-title">{{ isEditing ? t('habits.edit') : t('habits.create') }}</h2>

      <div class="form-group">
        <label class="form-label">{{ t('habits.name') }}</label>
        <input
          v-model="form.name"
          type="text"
          class="form-input"
          :placeholder="t('habits.name')"
        />
      </div>

      <div class="form-group" v-if="!isEditing">
        <label class="form-label">{{ t('habits.type') }}</label>
        <div class="select-wrapper">
          <select v-model="form.type" class="custom-select form-input">
            <option v-for="type in habitTypes" :key="type" :value="type">
              {{ t(`habits.types.${type}`) }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group" v-if="showGoalField">
        <label class="form-label">{{ t('habits.goal') }}</label>
        <input
          v-model.number="form.goal"
          type="number"
          class="form-input"
          min="0"
          step="0.5"
        />
      </div>

      <div class="form-group" v-if="form.type === 'additive'">
        <div class="switch-row">
          <label class="form-label">{{ t('habits.saveProgress') }}</label>
          <label class="switch">
            <input type="checkbox" v-model="form.saveProgress">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>

      <div class="form-group" v-if="form.type === 'additive' && form.saveProgress">
        <div class="switch-row">
          <label class="form-label">{{ t('habits.resetOnStreakBreak') }}</label>
          <label class="switch">
            <input type="checkbox" v-model="form.resetOnStreakBreak">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>

      <div class="form-group" v-if="form.type === 'boolean'">
        <label class="form-label">{{ t('habits.positiveValue') }}</label>
        <div class="select-wrapper">
          <select v-model="form.positiveValue" class="custom-select form-input">
            <option value="yes">{{ t('entry.yes') }}</option>
            <option value="no">{{ t('entry.no') }}</option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <div class="switch-row">
          <label class="form-label">{{ t('habits.trackStreak') }}</label>
          <label class="switch">
            <input type="checkbox" v-model="form.trackStreak">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>

      <div class="form-group">
        <div class="switch-row">
          <label class="form-label">{{ t('habits.priority') }}</label>
          <label class="switch">
            <input type="checkbox" v-model="form.priority">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>

      <div class="form-group">
        <div class="switch-row">
          <label class="form-label">{{ t('habits.notifications') }}</label>
          <label class="switch">
            <input type="checkbox" v-model="form.notificationsOn">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>

      <div class="form-group" v-if="form.notificationsOn">
        <label class="form-label">{{ t('habits.notificationTime') }}</label>
        <input
          v-model="form.notificationTime"
          type="time"
          class="form-input"
        />
      </div>

      <div class="form-actions">
        <button v-if="isEditing" class="btn btn-danger" @click="confirmDelete">
          {{ t('habits.delete') }}
        </button>
        <div class="spacer"></div>
        <button class="btn btn-secondary" @click="$emit('close')">
          {{ t('entry.cancel') }}
        </button>
        <button class="btn btn-primary" @click="save" :disabled="!form.name">
          {{ t('entry.save') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  habit: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'save', 'delete'])

const habitTypes = ['boolean', 'reachable', 'number', 'percentage', 'additive', 'mood', 'rating']

const form = ref({
  name: '',
  type: 'reachable',
  goal: 8,
  positiveValue: 'yes',
  saveProgress: false,
  resetOnStreakBreak: false,
  trackStreak: true,
  priority: false,
  notificationsOn: false,
  notificationTime: '21:00'
})

const isEditing = computed(() => !!props.habit)

const showGoalField = computed(() => {
  return ['reachable', 'additive'].includes(form.value.type)
})

watch(() => props.habit, (newHabit) => {
  if (newHabit) {
    form.value = {
      name: newHabit.name || '',
      type: newHabit.type || 'reachable',
      goal: newHabit.goal || 8,
      positiveValue: newHabit.positiveValue || 'yes',
      saveProgress: newHabit.saveProgress || false,
      resetOnStreakBreak: newHabit.resetOnStreakBreak || false,
      trackStreak: newHabit.trackStreak !== undefined ? newHabit.trackStreak : true,
      priority: newHabit.priority || false,
      notificationsOn: newHabit.notificationsOn || false,
      notificationTime: newHabit.notificationTime || '21:00'
    }
  } else {
    form.value = {
      name: '',
      type: 'reachable',
      goal: 8,
      positiveValue: 'yes',
      saveProgress: false,
      resetOnStreakBreak: false,
      trackStreak: true,
      priority: false,
      notificationsOn: false,
      notificationTime: '21:00'
    }
  }
}, { immediate: true })

function save() {
  if (!form.value.name) return
  emit('save', { ...form.value })
}

function confirmDelete() {
  if (confirm(t('habits.confirmDelete'))) {
    emit('delete')
  }
}
</script>

<style scoped>
.switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.switch-row .form-label {
  margin-bottom: 0;
}

.spacer {
  flex: 1;
}

.form-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-lg);
}
</style>
