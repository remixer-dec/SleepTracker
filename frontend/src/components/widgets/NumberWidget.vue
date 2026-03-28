<template>
  <div class="number-widget">
    <div class="value-display">
      <button class="adjust-btn" @click="decrement">-</button>
      <input
        type="number"
        class="value-input"
        :value="modelValue"
        @input="$emit('update:modelValue', parseFloat($event.target.value) || 0)"
      />
      <button class="adjust-btn" @click="increment">+</button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Number,
    default: 0
  },
  step: {
    type: Number,
    default: 1
  }
})

const emit = defineEmits(['update:modelValue'])

function increment() {
  emit('update:modelValue', props.modelValue + props.step)
}

function decrement() {
  emit('update:modelValue', props.modelValue - props.step)
}
</script>

<style scoped>
.number-widget {
  display: flex;
  justify-content: center;
}

.value-display {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
}

.adjust-btn {
  width: 48px;
  height: 48px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.06);
  color: var(--color-text);
  font-size: var(--font-size-xl);
  font-weight: 300;
}

.adjust-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.value-input {
  width: 100px;
  text-align: center;
  font-size: var(--font-size-xxl);
  font-weight: 600;
  background: transparent;
  border: none;
  color: var(--color-text);
}

.value-input::-webkit-outer-spin-button,
.value-input::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.value-input[type=number] {
  -moz-appearance: textfield;
}
</style>
