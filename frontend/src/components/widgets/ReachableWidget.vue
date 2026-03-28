<template>
  <div class="reachable-widget">
    <div class="value-display">
      <span class="value" :style="{ color: valueColor }">{{ modelValue }}</span>
      <span class="unit">{{ unit }}</span>
    </div>

    <div class="goal-indicator">
      <span class="goal-text">{{ goalText }}</span>
    </div>

    <input
      type="range"
      class="slider"
      :value="modelValue"
      @input="$emit('update:modelValue', parseFloat($event.target.value))"
      :min="0"
      :max="maxValue"
      :step="0.5"
    />

    <div class="quick-values">
      <button
        v-for="val in quickValues"
        :key="val"
        class="quick-btn"
        :class="{ active: modelValue === val }"
        @click="$emit('update:modelValue', val)"
      >
        {{ val }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: Number,
    default: 8
  },
  goal: {
    type: Number,
    default: 8
  },
  unit: {
    type: String,
    default: 'hours'
  }
})

defineEmits(['update:modelValue'])

const maxValue = computed(() => Math.max(props.goal + 4, 12))

const quickValues = computed(() => {
  const values = []
  for (let i = 0; i <= maxValue.value; i += 2) {
    values.push(i)
  }
  return values
})

const valueColor = computed(() => {
  if (props.modelValue >= props.goal) {
    return 'var(--color-success)'
  }
  if (props.modelValue >= props.goal - 1) {
    return 'var(--color-warning)'
  }
  return 'var(--color-danger)'
})

const goalText = computed(() => {
  const diff = props.modelValue - props.goal
  if (diff >= 0) {
    return diff > 0 ? `+${diff} ${t('goal.above')}` : t('goal.reached')
  }
  return `${Math.abs(diff)} ${t('goal.below')}`
})
</script>

<style scoped>
.reachable-widget {
  text-align: center;
}

.value-display {
  margin-bottom: var(--spacing-sm);
}

.value {
  font-size: var(--font-size-huge);
  font-weight: 700;
}

.unit {
  font-size: var(--font-size-md);
  color: var(--color-text-muted);
  margin-left: var(--spacing-xs);
}

.goal-indicator {
  margin-bottom: var(--spacing-md);
}

.goal-text {
  font-size: var(--font-size-sm);
  color: var(--color-text-dim);
}

.slider {
  width: 100%;
  height: 8px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.06);
  appearance: none;
  outline: none;
  margin-bottom: var(--spacing-md);
}

.slider::-webkit-slider-thumb {
  appearance: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.slider::-moz-range-thumb {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  border: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.quick-values {
  display: flex;
  gap: var(--spacing-xs);
  justify-content: center;
  flex-wrap: wrap;
}

.quick-btn {
  width: 40px;
  height: 40px;
  border-radius: var(--border-radius-sm);
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-text);
  font-size: var(--font-size-sm);
}

.quick-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.quick-btn.active {
  background: var(--color-accent);
  color: var(--color-primary-dark);
}
</style>
