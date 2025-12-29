<template>
  <div class="habit-title" @click="handleClick">
    <h1 class="title">{{ habit.name }}</h1>
    <button v-if="authStore.isOwner" class="btn-icon edit-btn">
      <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
        <path d="M11.7 1.3a1 1 0 011.4 0l1.6 1.6a1 1 0 010 1.4l-9 9-3.4.8a.5.5 0 01-.6-.6l.8-3.4 9.2-8.8z"/>
      </svg>
    </button>
  </div>
</template>

<script setup>
import { useAuthStore } from '../stores/auth.js'

const props = defineProps({
  habit: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['edit'])
const authStore = useAuthStore()

function handleClick() {
  if (authStore.isOwner) {
    emit('edit')
  }
}
</script>

<style scoped>
.habit-title {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
  cursor: pointer;
}

.title {
  font-size: var(--font-size-xl);
  font-weight: 600;
  color: var(--color-text);
  text-align: center;
}

.edit-btn {
  opacity: 0.5;
  transition: opacity 0.2s;
}

.habit-title:hover .edit-btn {
  opacity: 1;
}
</style>
