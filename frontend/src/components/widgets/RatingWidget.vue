<template>
  <div class="rating-widget">
    <div class="stars">
      <button
        v-for="star in 5"
        :key="star"
        class="star-btn"
        :class="{ filled: star <= modelValue }"
        @click="$emit('update:modelValue', star)"
      >
        <svg class="star-icon" viewBox="0 0 24 24">
          <path
            d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"
            :fill="star <= modelValue ? 'currentColor' : 'none'"
            stroke="currentColor"
            stroke-width="1.5"
          />
        </svg>
      </button>
    </div>
    <div class="rating-label">{{ ratingLabel }}</div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Number,
    default: 3
  }
})

defineEmits(['update:modelValue'])

const labels = ['', 'Poor', 'Fair', 'Good', 'Great', 'Excellent']

const ratingLabel = computed(() => labels[props.modelValue] || '')
</script>

<style scoped>
.rating-widget {
  text-align: center;
}

.stars {
  display: flex;
  justify-content: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.star-btn {
  padding: var(--spacing-xs);
  background: none;
  color: var(--color-text-dim);
  transition: all 0.15s;
}

.star-btn:hover {
  transform: scale(1.1);
}

.star-btn.filled {
  color: #ffc107;
}

.star-icon {
  width: 40px;
  height: 40px;
}

.rating-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
}
</style>
