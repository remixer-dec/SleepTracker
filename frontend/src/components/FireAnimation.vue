<template>
  <div class="fire" :style="containerStyle">
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
      </defs>

      <!-- Replaced <g filter="url(#fireGlow)"> with a CSS class -->
      <g class="flame-group">
        <path
          class="flame flame-1"
          fill="url(#fireGradient)"
          d="M25 75 C10 60 5 45 10 30 C15 15 25 5 25 0 C25 5 35 15 40 30 C45 45 40 60 25 75"
        />
        <path
          class="flame flame-2"
          fill="url(#fireGradient)"
          opacity="0.8"
          d="M25 75 C15 65 12 50 15 38 C18 26 25 18 25 12 C25 18 32 26 35 38 C38 50 35 65 25 75"
        />
        <path
          class="flame flame-3"
          fill="#ffeb3b"
          opacity="0.6"
          d="M25 75 C18 68 16 55 18 45 C20 35 25 28 25 22 C25 28 30 35 32 45 C34 55 32 68 25 75"
        />
      </g>
    </svg>
  </div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  intensity: {
    type: Number,
    default: 0.5,
  },
});

const size = computed(() => 40 + props.intensity * 60);
const containerStyle = computed(() => ({
  "--intensity-opacity": 0.3 + props.intensity * 0.5,
  "--flame-duration": 1.5 - props.intensity * 0.8 + "s",
  opacity: "var(--intensity-opacity)",
}));
</script>

<style scoped>
.fire {
  display: flex;
  justify-content: center;
  contain: layout paint;
}

.fire-svg {
  overflow: visible;
  transform: translateZ(0);
}

.flame-group {
  filter: blur(2px);
  transform: translateY(18px)
}

.flame {
  transform-origin: center bottom;
  will-change: transform, opacity;
}

.flame-1 {
  animation: flicker1 var(--flame-duration) ease-in-out infinite alternate;
}

.flame-2 {
  animation: flicker2 calc(var(--flame-duration) * 0.8) ease-in-out infinite
    alternate;
}

.flame-3 {
  animation: flicker3 calc(var(--flame-duration) * 0.6) ease-in-out infinite
    alternate;
}

@keyframes flicker1 {
  0% {
    transform: scale(1, 1);
  }
  50% {
    transform: scale(0.95, 1.05);
  }
  100% {
    transform: scale(1.02, 0.95);
  }
}

@keyframes flicker2 {
  0% {
    transform: scale(1.02, 0.98) translateX(-1px);
  }
  100% {
    transform: scale(0.98, 1.02) translateX(1px);
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
