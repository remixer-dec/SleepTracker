import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../utils/api.js'

export const ACHIEVEMENT_DEFINITIONS = [
  { id: 'first_entry', icon: '1', condition: (stats) => stats.totalEntries >= 1 },
  { id: 'week_streak', icon: '7', condition: (stats) => stats.longestStreak >= 7 },
  { id: 'two_week_streak', icon: '14', condition: (stats) => stats.longestStreak >= 14 },
  { id: 'month_streak', icon: '30', condition: (stats) => stats.longestStreak >= 30 },
  { id: 'quarter_streak', icon: '90', condition: (stats) => stats.longestStreak >= 90 },
  { id: 'year_streak', icon: '365', condition: (stats) => stats.longestStreak >= 365 },
  { id: 'perfect_week', icon: 'W', condition: (stats) => stats.perfectWeeks >= 1 },
  { id: 'perfect_month', icon: 'M', condition: (stats) => stats.perfectMonths >= 1 },
  { id: 'early_bird', icon: 'E', condition: (stats) => stats.earlyEntries >= 1 },
  { id: 'night_owl', icon: 'N', condition: (stats) => stats.nightEntries >= 1 },
  { id: 'comeback', icon: 'C', condition: (stats) => stats.comebacks >= 1 },
  { id: 'consistent', icon: 'K', condition: (stats) => stats.consistentWeeks >= 1 },
  { id: 'improver', icon: 'I', condition: (stats) => stats.improvements >= 1 },
  { id: 'multi_tracker', icon: '3', condition: (stats) => stats.habitsTracked >= 3 },
  { id: 'data_lover', icon: 'D', condition: (stats) => stats.statsViews >= 10 },
  { id: 'first_habit', icon: 'H', condition: (stats) => stats.habitsCreated >= 1 },
  { id: 'five_habits', icon: '5', condition: (stats) => stats.habitsCreated >= 5 },
  { id: 'hundred_entries', icon: '100', condition: (stats) => stats.totalEntries >= 100 },
  { id: 'five_hundred_entries', icon: '500', condition: (stats) => stats.totalEntries >= 500 },
  { id: 'level_five', icon: 'L5', condition: (stats) => stats.level >= 5 },
  { id: 'level_ten', icon: 'L10', condition: (stats) => stats.level >= 10 },
  { id: 'level_twenty', icon: 'L20', condition: (stats) => stats.level >= 20 }
]

export const useAchievementsStore = defineStore('achievements', () => {
  const unlockedAchievements = ref([])
  const statsViews = ref(0)

  const unlockedIds = computed(() => {
    return new Set(unlockedAchievements.value.map(a => a.type))
  })

  function isUnlocked(achievementId) {
    return unlockedIds.value.has(achievementId)
  }

  async function fetchAchievements() {
    try {
      const data = await api.get('achievements')
      unlockedAchievements.value = data || []
    } catch (error) {
      console.error('Failed to fetch achievements:', error)
    }
  }

  async function unlockAchievement(type) {
    if (isUnlocked(type)) return
    try {
      const data = await api.post('achievements', { type })
      unlockedAchievements.value.push(data)
    } catch (error) {
      console.error('Failed to unlock achievement:', error)
    }
  }

  async function checkAndUnlockAchievements(stats) {
    for (const achievement of ACHIEVEMENT_DEFINITIONS) {
      if (!isUnlocked(achievement.id) && achievement.condition(stats)) {
        await unlockAchievement(achievement.id)
      }
    }
  }

  function incrementStatsViews() {
    statsViews.value++
  }

  return {
    unlockedAchievements,
    statsViews,
    unlockedIds,
    isUnlocked,
    fetchAchievements,
    unlockAchievement,
    checkAndUnlockAchievements,
    incrementStatsViews
  }
})
