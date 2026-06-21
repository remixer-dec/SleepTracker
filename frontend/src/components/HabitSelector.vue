<template>
  <div class="habit-selector">
    <div class="select-wrapper">
      <select
        v-model="selectedId"
        class="custom-select"
        :disabled="habits.length === 0"
      >
        <option v-if="habits.length === 0" value="">
          {{ t("habits.noHabits") }}
        </option>
        <option v-for="habit in habits" :key="habit.id" :value="habit.id">
          {{ habit.name }}
        </option>
      </select>
      <span class="select-arrow">
        <svg width="12" height="12" viewBox="0 0 12 12" fill="currentColor">
          <path d="M2 4l4 4 4-4" />
        </svg>
      </span>
    </div>
    <div class="selector-actions">
      <button
        v-if="authStore.isOwner"
        class="btn-icon action-btn"
        @click="$emit('add-habit')"
        :title="t('habits.create')"
      >
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
          <path
            d="M10 4v12M4 10h12"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
          />
        </svg>
      </button>
      <button
        class="btn-icon action-btn"
        @click="cycleTheme"
        :title="themeLabel"
      >
        <svg width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round">
          <circle cx="9" cy="9" r="7"/>
          <path d="M9 2 A7 7 0 0 1 9 16 Z" fill="currentColor" stroke="none"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "../stores/auth.js";
import { useHabitsStore } from "../stores/habits.js";

const { t } = useI18n();
const authStore = useAuthStore();
const habitsStore = useHabitsStore();

defineEmits(["add-habit"]);

const habits = computed(() => habitsStore.habits);

const selectedId = computed({
  get: () => habitsStore.selectedHabitId,
  set: (value) => habitsStore.selectHabit(value),
});

const themes = ['', 'warm', 'light']
const themeLabels = { '': 'Original', warm: 'Warm Stone', light: 'Light Warm' }

const currentTheme = ref(localStorage.getItem('theme') || '')

const themeLabel = computed(() => themeLabels[currentTheme.value] || 'Original')

function cycleTheme() {
  const nextIndex = (themes.indexOf(currentTheme.value) + 1) % themes.length
  currentTheme.value = themes[nextIndex]
  if (currentTheme.value) {
    document.documentElement.dataset.theme = currentTheme.value
    localStorage.setItem('theme', currentTheme.value)
  } else {
    delete document.documentElement.dataset.theme
    localStorage.removeItem('theme')
  }
}
</script>

<style scoped>
.habit-selector {
  display: flex;
  gap: var(--spacing-sm);
  align-items: center;
  width: 100%;
}

.select-wrapper {
  flex: 1;
}

.selector-actions {
  position: absolute;
  top: 23px;
  display: flex;
  gap: var(--spacing-xs);
}

.action-btn {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--input-bg);
  border-radius: var(--border-radius-sm);
}

.action-btn:hover {
  background: var(--btn-icon-hover-bg);
}

@media (max-width: 767px) {
  .habit-selector {
    margin-top: 10px;
  }
}
</style>
