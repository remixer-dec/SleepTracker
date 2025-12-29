<template>
  <div class="habit-selector">
    <div class="select-wrapper">
      <select
        v-model="selectedId"
        class="custom-select"
        :disabled="habits.length === 0"
      >
        <option v-if="habits.length === 0" value="">{{ t('habits.noHabits') }}</option>
        <option v-for="habit in habits" :key="habit.id" :value="habit.id">
          {{ habit.name }}
        </option>
      </select>
      <span class="select-arrow">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
          <path d="M2 4l4 4 4-4"/>
        </svg>
      </span>
    </div>
    <button
      v-if="authStore.isOwner"
      class="btn-icon add-btn"
      @click="$emit('add-habit')"
      :title="t('habits.create')"
    >
      <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
        <path d="M10 4v12M4 10h12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth.js'
import { useHabitsStore } from '../stores/habits.js'

const { t } = useI18n()
const authStore = useAuthStore()
const habitsStore = useHabitsStore()

defineEmits(['add-habit'])

const habits = computed(() => habitsStore.habits)

const selectedId = computed({
  get: () => habitsStore.selectedHabitId,
  set: (value) => habitsStore.selectHabit(value)
})
</script>

<style scoped>
.habit-selector {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
}

.select-wrapper {
  flex: 1;
}

.add-btn {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 65, 78, 0.4);
  border-radius: var(--border-radius-sm);
}

.add-btn:hover {
  background: rgba(0, 65, 78, 0.6);
}
</style>
