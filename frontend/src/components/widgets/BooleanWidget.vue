<template>
  <div class="boolean-widget">
    <button
      class="bool-btn"
      :class="{
        active: modelValue === 1,
        positive: positiveValue !== 'no',
        negative: positiveValue === 'no'
      }"
      @click="$emit('update:modelValue', 1)"
    >
      {{ positiveLabel }}
    </button>
    <button
      class="bool-btn"
      :class="{
        active: modelValue === 0,
        positive: positiveValue === 'no',
        negative: positiveValue !== 'no'
      }"
      @click="$emit('update:modelValue', 0)"
    >
      {{ negativeLabel }}
    </button>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  modelValue: {
    type: Number,
    default: 0
  },
  positiveLabel: {
    type: String,
    default: ''
  },
  positiveValue: {
    type: String,
    default: 'yes'
  }
})

defineEmits(['update:modelValue'])

const negativeLabel = computed(() => t('entry.no'))
</script>

<style scoped>
.boolean-widget {
  display: flex;
  gap: var(--spacing-md);
}

.bool-btn {
  flex: 1;
  padding: var(--spacing-md);
  border-radius: var(--border-radius);
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text);
  font-size: var(--font-size-lg);
  font-weight: 500;
  transition: all 0.2s;
}

.bool-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

.bool-btn.positive.active {
  background: var(--color-success);
  color: white;
}

.bool-btn.negative.active {
  background: var(--color-danger);
}
</style>
