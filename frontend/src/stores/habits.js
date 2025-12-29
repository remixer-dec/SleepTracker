import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../utils/api.js'

export const useHabitsStore = defineStore('habits', () => {
  const habits = ref([])
  const selectedHabitId = ref(null)
  const entries = ref({})
  const isLoading = ref(false)

  const selectedHabit = computed(() => {
    return habits.value.find(h => h.id === selectedHabitId.value) || null
  })

  const currentEntries = computed(() => {
    return entries.value[selectedHabitId.value] || []
  })

  async function fetchHabits() {
    isLoading.value = true
    try {
      const data = await api.get('habits')
      habits.value = data || []
      if (habits.value.length > 0 && !selectedHabitId.value) {
        selectedHabitId.value = habits.value[0].id
      }
    } catch (error) {
      console.error('Failed to fetch habits:', error)
    } finally {
      isLoading.value = false
    }
  }

  async function createHabit(habit) {
    try {
      const data = await api.post('habits', habit)
      habits.value.push(data)
      if (!selectedHabitId.value) {
        selectedHabitId.value = data.id
      }
      return data
    } catch (error) {
      console.error('Failed to create habit:', error)
      throw error
    }
  }

  async function updateHabit(id, updates) {
    try {
      const data = await api.put(`habits/${id}`, updates)
      const index = habits.value.findIndex(h => h.id === id)
      if (index !== -1) {
        habits.value[index] = data
      }
      return data
    } catch (error) {
      console.error('Failed to update habit:', error)
      throw error
    }
  }

  async function deleteHabit(id) {
    try {
      await api.delete(`habits/${id}`)
      habits.value = habits.value.filter(h => h.id !== id)
      if (selectedHabitId.value === id) {
        selectedHabitId.value = habits.value[0]?.id || null
      }
      delete entries.value[id]
    } catch (error) {
      console.error('Failed to delete habit:', error)
      throw error
    }
  }

  async function fetchEntries(habitId, startDate, endDate) {
    try {
      let url = `entries?habitId=${habitId}`
      if (startDate && endDate) {
        url += `&start=${startDate}&end=${endDate}`
      }
      const data = await api.get(url)
      entries.value[habitId] = data || []
    } catch (error) {
      console.error('Failed to fetch entries:', error)
    }
  }

  async function createOrUpdateEntry(entry) {
    try {
      const data = await api.post('entries', entry)
      const habitEntries = entries.value[entry.habitId] || []
      const existingIndex = habitEntries.findIndex(e => e.date === entry.date)
      if (existingIndex !== -1) {
        habitEntries[existingIndex] = data
      } else {
        habitEntries.push(data)
      }
      entries.value[entry.habitId] = habitEntries
      return data
    } catch (error) {
      console.error('Failed to save entry:', error)
      throw error
    }
  }

  async function deleteEntry(id, habitId) {
    try {
      await api.delete(`entries/${id}`)
      entries.value[habitId] = (entries.value[habitId] || []).filter(e => e.id !== id)
    } catch (error) {
      console.error('Failed to delete entry:', error)
      throw error
    }
  }

  function selectHabit(id) {
    selectedHabitId.value = id
  }

  function getEntryByDate(habitId, date) {
    const habitEntries = entries.value[habitId] || []
    return habitEntries.find(e => e.date === date)
  }

  return {
    habits,
    selectedHabitId,
    entries,
    isLoading,
    selectedHabit,
    currentEntries,
    fetchHabits,
    createHabit,
    updateHabit,
    deleteHabit,
    fetchEntries,
    createOrUpdateEntry,
    deleteEntry,
    selectHabit,
    getEntryByDate
  }
})
