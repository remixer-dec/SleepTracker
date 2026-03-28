<template>
  <div class="mood-widget">
    <button
      v-for="mood in moods"
      :key="mood.value"
      class="mood-btn"
      :class="{ active: modelValue === mood.value }"
      @click="$emit('update:modelValue', mood.value)"
    >
      <svg class="mood-face" viewBox="0 0 50 50">
        <circle cx="25" cy="25" r="23" fill="none" :stroke="mood.color" stroke-width="2"/>
        <circle cx="17" cy="20" r="3" :fill="mood.color"/>
        <circle cx="33" cy="20" r="3" :fill="mood.color"/>
        <path :d="mood.mouth" :stroke="mood.color" stroke-width="2.5" fill="none" stroke-linecap="round"/>
      </svg>
      <span class="mood-label">{{ mood.label }}</span>
    </button>
  </div>
</template>

<script setup>
defineProps({
  modelValue: {
    type: Number,
    default: 1
  }
})

defineEmits(['update:modelValue'])

const moods = [
  {
    value: 0,
    label: 'Bad',
    color: '#b85c5c',
    mouth: 'M15 35 Q25 28 35 35'
  },
  {
    value: 1,
    label: 'Okay',
    color: '#d4a054',
    mouth: 'M15 32 L35 32'
  },
  {
    value: 2,
    label: 'Good',
    color: '#7a9a6d',
    mouth: 'M15 30 Q25 38 35 30'
  }
]
</script>

<style scoped>
.mood-widget {
  display: flex;
  justify-content: center;
  gap: var(--spacing-lg);
}

.mood-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-md);
  border-radius: var(--border-radius);
  background: rgba(255, 255, 255, 0.03);
  transition: all 0.2s;
}

.mood-btn:hover {
  background: rgba(255, 255, 255, 0.06);
}

.mood-btn.active {
  background: rgba(255, 255, 255, 0.08);
}

.mood-face {
  width: 60px;
  height: 60px;
}

.mood-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.mood-btn.active .mood-label {
  color: var(--color-text);
}
</style>
