<template>
  <div class="percentage-widget">
    <div class="circle-container">
      <svg class="circle-svg" viewBox="0 0 100 100">
        <circle
          class="circle-bg"
          cx="50"
          cy="50"
          r="45"
          fill="none"
          stroke-width="8"
        />
        <circle
          class="circle-progress"
          cx="50"
          cy="50"
          r="45"
          fill="none"
          stroke-width="8"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="offset"
          :style="{ stroke: progressColor }"
        />
      </svg>
      <div class="circle-value">
        <span class="value">{{ modelValue }}</span>
        <span class="percent">%</span>
      </div>
    </div>

    <input
      type="range"
      class="slider"
      :value="modelValue"
      @input="$emit('update:modelValue', parseInt($event.target.value))"
      min="0"
      max="100"
      step="5"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Number,
    default: 50
  }
})

defineEmits(['update:modelValue'])

const circumference = 2 * Math.PI * 45

const offset = computed(() => {
  return circumference - (props.modelValue / 100) * circumference
})

const progressColor = computed(() => {
  if (props.modelValue >= 70) return 'var(--color-success)'
  if (props.modelValue >= 40) return 'var(--color-warning)'
  return 'var(--color-danger)'
})
</script>

<style scoped>
.percentage-widget {
  text-align: center;
}

.circle-container {
  position: relative;
  width: 140px;
  height: 140px;
  margin: 0 auto var(--spacing-md);
}

.circle-svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.circle-bg {
  stroke: var(--widget-bg);
}

.circle-progress {
  transition: stroke-dashoffset 0.3s ease;
  stroke-linecap: round;
}

.circle-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
}

.value {
  font-size: var(--font-size-xxl);
  font-weight: 700;
  color: var(--color-text);
}

.percent {
  font-size: var(--font-size-md);
  color: var(--color-text-muted);
}

.slider {
  width: 100%;
  height: 8px;
  border-radius: 4px;
  background: var(--widget-bg);
  appearance: none;
  outline: none;
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
</style>
