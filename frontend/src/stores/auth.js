import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../utils/api.js'

export const useAuthStore = defineStore('auth', () => {
  const userId = ref('')
  const isOwner = ref(false)
  const isAuthenticated = ref(false)
  const isLoading = ref(true)

  const isGuest = computed(() => !isOwner.value)

  async function checkAuth() {
    isLoading.value = true
    try {
      const data = await api.get('me')
      userId.value = data.userId || ''
      isOwner.value = data.isOwner || false
      isAuthenticated.value = data.authenticated || false
    } catch (error) {
      console.error('Auth check failed:', error)
      userId.value = ''
      isOwner.value = false
      isAuthenticated.value = false
    } finally {
      isLoading.value = false
    }
  }

  async function join(token) {
    try {
      const data = await api.post(`join?token=${token}`)
      userId.value = data.userId
      isOwner.value = data.isOwner
      isAuthenticated.value = true
      if (data.token) {
        localStorage.setItem('token', data.token)
      }
      return true
    } catch (error) {
      console.error('Join failed:', error)
      return false
    }
  }

  function logout() {
    localStorage.removeItem('token')
    userId.value = ''
    isOwner.value = false
    isAuthenticated.value = false
  }

  return {
    userId,
    isOwner,
    isAuthenticated,
    isLoading,
    isGuest,
    checkAuth,
    join,
    logout
  }
})
