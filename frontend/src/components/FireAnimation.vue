<template>
  <div class="fire" :style="fireStyle">
    <svg
      :width="size"
      :height="size * 1.5"
      viewBox="0 0 50 75"
      class="fire-svg"
    >
      <defs>
        <linearGradient id="fireGradient" x1="0%" y1="100%" x2="0%" y2="0%">
          <stop offset="0%" stop-color="#ff6b35" />
          <stop offset="40%" stop-color="#f7931e" />
          <stop offset="70%" stop-color="#ffc107" />
          <stop offset="100%" stop-color="#ffeb3b" />
        </linearGradient>
        <filter id="fireGlow">
          <feGaussianBlur stdDeviation="2" result="coloredBlur"/>
          <feMerge>
            <feMergeNode in="coloredBlur"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
      </defs>

      <g filter="url(#fireGlow)">
        <path
          class="flame flame-1"
          fill="url(#fireGradient)"
          d="M25 75 C10 60 5 45 10 30 C15 15 25 5 25 0 C25 5 35 15 40 30 C45 45 40 60 25 75"
          :style="{ animationDuration: flameDuration + 's' }"
        />
        <path
          class="flame flame-2"
          fill="url(#fireGradient)"
          opacity="0.8"
          d="M25 75 C15 65 12 50 15 38 C18 26 25 18 25 12 C25 18 32 26 35 38 C38 50 35 65 25 75"
          :style="{ animationDuration: (flameDuration * 0.8) + 's' }"
        />
        <path
          class="flame flame-3"
          fill="#ffeb3b"
          opacity="0.6"
          d="M25 75 C18 68 16 55 18 45 C20 35 25 28 25 22 C25 28 30 35 32 45 C34 55 32 68 25 75"
          :style="{ animationDuration: (flameDuration * 0.6) + 's' }"
        />
      </g>
    </svg>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  intensity: {
    type: Number,
    default: 0.5
  }
})

const size = computed(() => 40 + (props.intensity * 60))
const flameDuration = computed(() => 1.5 - (props.intensity * 0.8))

const fireStyle = computed(() => ({
  opacity: 0.3 + (props.intensity * 0.5)
}))
</script>

<style scoped>
.fire {
  display: flex;
  justify-content: center;
}

.fire-svg {
  overflow: visible;
}

.flame {
  transform-origin: center bottom;
}

.flame-1 {
  animation: flicker1 1.5s ease-in-out infinite alternate;
}

.flame-2 {
  animation: flicker2 1.2s ease-in-out infinite alternate;
}

.flame-3 {
  animation: flicker3 0.9s ease-in-out infinite alternate;
}

@keyframes flicker1 {
  0% {
    transform: scaleY(1) scaleX(1);
  }
  50% {
    transform: scaleY(1.05) scaleX(0.95);
  }
  100% {
    transform: scaleY(0.95) scaleX(1.02);
  }
}

@keyframes flicker2 {
  0% {
    transform: scaleY(0.98) scaleX(1.02) translateX(-1px);
  }
  100% {
    transform: scaleY(1.02) scaleX(0.98) translateX(1px);
  }
}

@keyframes flicker3 {
  0% {
    transform: scaleY(1) translateY(0);
    opacity: 0.6;
  }
  100% {
    transform: scaleY(0.9) translateY(2px);
    opacity: 0.8;
  }
}
</style>
