<template>
  <div class="streak-display glass-card">
    <div
      class="fire-container"
      v-if="!isAntiStreak && streak >= 2"
      :style="fireContainerStyle"
    >
      <FireAnimation :intensity="fireIntensity" />
    </div>

    <div class="streak-main">
      <div class="streak-label">
        {{ isAntiStreak ? t("streak.antiStreak") : t("streak.current") }}
      </div>
      <div class="streak-value" :class="{ 'anti-streak': isAntiStreak }">
        {{ Math.abs(streak) }}
      </div>
      <div class="streak-unit">{{ t("streak.days") }}</div>
      <div class="level-badge">{{ t("streak.level") }} {{ level.level }}</div>
    </div>

    <div class="streak-stats">
      <div class="stat">
        <div class="stat-value">{{ record }}</div>
        <div class="stat-label">{{ t("streak.record") }}</div>
      </div>
      <div class="stat">
        <div class="stat-value">{{ weeklyStreaks }}</div>
        <div class="stat-label">{{ t("streak.weekly") }}</div>
      </div>
    </div>

    <div class="level-progress">
      <div class="level-bar">
        <div
          class="level-fill"
          :style="{ width: level.progress * 100 + '%' }"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import { getFireIntensity } from "../utils/calculations.js";
import FireAnimation from "./FireAnimation.vue";

const { t } = useI18n();

const props = defineProps({
  streak: {
    type: Number,
    default: 0,
  },
  record: {
    type: Number,
    default: 0,
  },
  weeklyStreaks: {
    type: Number,
    default: 0,
  },
  level: {
    type: Object,
    default: () => ({ level: 1, progress: 0 }),
  },
});

const isAntiStreak = computed(() => props.streak < 0);
const fireIntensity = computed(() => getFireIntensity(props.streak));
const fireContainerStyle = computed(() => {
  const scale = Math.min(1 + props.streak / 10, 10);
  const filter =
    props.streak > 10
      ? `saturate(${Math.min(100 + (props.streak - 10) * 9, 1000)}%)`
      : "saturate(100%)";

  return {
    transform: `translateX(-50%) scale(${scale})`,
    filter,
  };
});
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
  pointer-events: none;
  z-index: 0;
  transform-origin: bottom center;
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
  color: var(--color-text);
  line-height: 1;
  margin: var(--spacing-xs) 0;
}

.streak-value.anti-streak {
  color: var(--color-danger);
}

.rage-icon {
  font-size: 0.6em;
  margin-right: var(--spacing-xs);
  animation: ragePulse 0.5s ease-in-out infinite alternate;
}

@keyframes ragePulse {
  from {
    transform: scale(1);
  }
  to {
    transform: scale(1.1);
  }
}

.streak-unit {
  font-size: var(--font-size-sm);
  color: var(--color-text);
}

.level-badge {
  display: inline-block;
  margin-top: var(--spacing-xs);
  padding: 2px var(--spacing-sm);
  background: var(--color-bg-overlay);
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-xs);
  color: var(--color-accent);
  position: relative;
  z-index: 2;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

.streak-stats {
  display: flex;
  justify-content: space-around;
  gap: var(--spacing-xl);
  margin: var(--spacing-md) 0;
  position: relative;
  z-index: 1;
}

@media (max-width: 767px) {
  .streak-stats {
    gap: var(--spacing-sm);
  }
}

.stat {
  text-align: center;
  min-width: 80px;
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
  background: var(--widget-bg);
  border-radius: 2px;
  overflow: hidden;
}

.level-fill {
  height: 100%;
  background: linear-gradient(
    90deg,
    var(--color-accent),
    var(--color-accent-soft)
  );
  border-radius: 2px;
  transition: width 0.3s ease;
}
</style>
