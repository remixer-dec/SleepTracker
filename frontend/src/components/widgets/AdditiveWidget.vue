<template>
  <div class="additive-widget">
    <div class="progress-container">
      <div class="progress-bar">
        <div
          class="progress-fill"
          :style="{ width: progressWidth, background: progressColor }"
        ></div>
      </div>
      <div class="progress-labels">
        <span>0</span>
        <span>{{ goal }}</span>
      </div>
    </div>

    <div class="value-display">
      <button class="adjust-btn" @click="decrement">-</button>
      <div class="value-info">
        <span class="value">{{ modelValue }}</span>
        <span class="goal-text">/ {{ goal }}</span>
      </div>
      <button class="adjust-btn" @click="increment">+</button>
    </div>

    <div class="status" :style="{ color: statusColor }">
      {{ statusText }}
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Number,
    default: 0
  },
  goal: {
    type: Number,
    default: 10
  }
})

const emit = defineEmits(['update:modelValue'])

const progressWidth = computed(() => {
  const percent = Math.min((props.modelValue / props.goal) * 100, 100)
  return `${percent}%`
})

const progressColor = computed(() => {
  const ratio = props.modelValue / props.goal
  if (ratio >= 1) return 'var(--color-success)'
  if (ratio >= 0.7) return 'var(--color-warning)'
  return 'var(--color-danger)'
})

const statusColor = computed(() => progressColor.value)

const statusText = computed(() => {
  const remaining = props.goal - props.modelValue
  if (remaining <= 0) return 'Goal completed!'
  return `${remaining} more to goal`
})

function increment() {
  emit('update:modelValue', props.modelValue + 1)
}

function decrement() {
  emit('update:modelValue', Math.max(0, props.modelValue - 1))
}
</script>

<style scoped>
.additive-widget {
  text-align: center;
}

.progress-container {
  margin-bottom: var(--spacing-md);
}

.progress-bar {
  height: 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 6px;
  transition: width 0.3s ease;
}

.progress-labels {
  display: flex;
  justify-content: space-between;
  margin-top: var(--spacing-xs);
  font-size: var(--font-size-xs);
  color: var(--color-text-dim);
}

.value-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-lg);
  margin-bottom: var(--spacing-sm);
}

.adjust-btn {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text);
  font-size: var(--font-size-xl);
  font-weight: 300;
}

.adjust-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

.value-info {
  text-align: center;
}

.value {
  font-size: var(--font-size-huge);
  font-weight: 700;
  color: var(--color-text);
}

.goal-text {
  font-size: var(--font-size-lg);
  color: var(--color-text-dim);
  margin-left: var(--spacing-xs);
}

.status {
  font-size: var(--font-size-sm);
}
</style>
