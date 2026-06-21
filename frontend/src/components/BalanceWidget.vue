<template>
  <div class="balance-widget">
    <span class="balance-label">{{ t("stats.balance") }}</span>
    <span class="balance-value" :style="{ color: balanceColor }">
      {{ formattedBalance }}
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  habit: { type: Object, required: true },
  entries: { type: Array, default: () => [] },
  monthDate: { type: Date, required: true }
})

const balance = computed(() => {
  const year = props.monthDate.getFullYear()
  const month = props.monthDate.getMonth()
  const goal = props.habit.goal || 0

  let total = 0
  for (const entry of props.entries) {
    const [y, m] = entry.date.split('-').map(Number)
    if (y === year && m - 1 === month) {
      total += (entry.value ?? 0) - goal
    }
  }
  return total
})

const formattedBalance = computed(() => {
  const val = Math.round(balance.value * 10) / 10
  return val > 0 ? `+${val}` : `${val}`
})

const balanceColor = computed(() => {
  if (balance.value > 0) return 'var(--color-success)'
  if (balance.value < 0) return 'var(--color-danger)'
  return 'var(--color-text-muted)'
})
</script>

<style scoped>
.balance-widget {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-xs) var(--spacing-md);
}

@media (max-width: 767px) {
  .balance-widget {
    padding-bottom: var(--spacing-sm);
  }
}

.balance-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-dim);
}

.balance-value {
  font-size: var(--font-size-md);
  font-weight: 600;
}
</style>
