<template>
  <div class="streak-display glass-card">
    <div class="fire-container" v-if="!isAntiStreak && streak > 0">
      <FireAnimation :intensity="fireIntensity" />
    </div>

    <div class="streak-main">
      <div class="streak-label">{{ isAntiStreak ? t('streak.antiStreak') : t('streak.current') }}</div>
      <div class="streak-value" :class="{ 'anti-streak': isAntiStreak }">
        {{ Math.abs(streak) }}
      </div>
      <div class="streak-unit">{{ t('streak.days') }}</div>
      <div class="level-badge">{{ t('streak.level') }} {{ level.level }}</div>
    </div>

    <div class="streak-stats">
      <div class="stat">
        <div class="stat-value">{{ record }}</div>
        <div class="stat-label">{{ t('streak.record') }}</div>
      </div>
      <div class="stat">
        <div class="stat-value">{{ weeklyStreaks }}</div>
        <div class="stat-label">{{ t('streak.weekly') }}</div>
      </div>
    </div>

    <div class="level-progress">
      <div class="level-bar">
        <div class="level-fill" :style="{ width: (level.progress * 100) + '%' }"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getFireIntensity } from '../utils/calculations.js'
import FireAnimation from './FireAnimation.vue'

const { t } = useI18n()

const props = defineProps({
  streak: {
    type: Number,
    default: 0
  },
  record: {
    type: Number,
    default: 0
  },
  weeklyStreaks: {
    type: Number,
    default: 0
  },
  level: {
    type: Object,
    default: () => ({ level: 1, progress: 0 })
  }
})

const isAntiStreak = computed(() => props.streak < 0)
const fireIntensity = computed(() => getFireIntensity(props.streak))
</script>

<style scoped>
.streak-display {
  padding: var(--spacing-md);
  text-align: center;
  position: relative;
  overflow: hidden;
}

.fire-container {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  pointer-events: none;
  z-index: 0;
}

.streak-main {
  position: relative;
  z-index: 1;
  margin-bottom: var(--spacing-sm);
}

.streak-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.streak-value {
  font-size: var(--font-size-huge);
  font-weight: 700;
  color: var(--color-success);
  line-height: 1;
  margin: var(--spacing-xs) 0;
}

.streak-value.anti-streak {
  color: var(--color-danger);
}

.streak-unit {
  font-size: var(--font-size-sm);
  color: var(--color-text-muted);
}

.level-badge {
  display: inline-block;
  margin-top: var(--spacing-xs);
  padding: 2px var(--spacing-sm);
  background: rgba(0, 180, 216, 0.2);
  border-radius: 12px;
  font-size: var(--font-size-xs);
  color: var(--color-accent);
}

.streak-stats {
  display: flex;
  justify-content: center;
  gap: var(--spacing-xl);
  margin: var(--spacing-md) 0;
  position: relative;
  z-index: 1;
}

.stat {
  text-align: center;
}

.stat-value {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--color-text);
}

.stat-label {
  font-size: var(--font-size-xs);
  color: var(--color-text-dim);
}

.level-progress {
  position: relative;
  z-index: 1;
}

.level-bar {
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  overflow: hidden;
}

.level-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-soft));
  border-radius: 2px;
  transition: width 0.3s ease;
}
</style>
