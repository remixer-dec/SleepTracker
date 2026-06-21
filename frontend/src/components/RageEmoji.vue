<template>
  <svg viewBox="0 0 100 100" class="rage-emoji">
    <circle cx="50" cy="50" r="45" fill="#8f4444"/>
    <g class="face">
      <path
        d="M25 35 L40 42"
        stroke="#1a1a1a"
        stroke-width="4"
        stroke-linecap="round"
        class="eyebrow left"
      />
      <path
        d="M75 35 L60 42"
        stroke="#1a1a1a"
        stroke-width="4"
        stroke-linecap="round"
        class="eyebrow right"
      />
      <circle cx="35" cy="45" r="5" fill="#1a1a1a"/>
      <circle cx="65" cy="45" r="5" fill="#1a1a1a"/>
      <path
        d="M30 70 Q50 55 70 70"
        stroke="#1a1a1a"
        stroke-width="4"
        stroke-linecap="round"
        fill="none"
        class="mouth"
      />
    </g>
  </svg>
</template>

<style scoped>
/* CSS overrides SVG presentation attributes */
.rage-emoji circle:first-child {
  fill: var(--color-rage);
}

.rage-emoji {
  width: 100%;
  height: 100%;
  max-width: 80vw;
  max-height: 80vh;
}

.face {
  animation: shake 0.3s ease-in-out infinite;
}

.eyebrow {
  animation: eyebrowShake 0.2s ease-in-out infinite alternate;
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-1px) rotate(-0.5deg);
  }
  75% {
    transform: translateX(1px) rotate(0.5deg);
  }
}

@keyframes eyebrowShake {
  from {
    transform: translateY(0);
  }
  to {
    transform: translateY(-1px);
  }
}
</style>
