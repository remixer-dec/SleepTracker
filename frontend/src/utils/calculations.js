export function formatDate(date) {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

export function parseDate(dateString) {
  const [year, month, day] = dateString.split('-').map(Number)
  return new Date(year, month - 1, day)
}

export function addDays(date, days) {
  const result = new Date(date)
  result.setDate(result.getDate() + days)
  return result
}

export function getDaysBetween(start, end) {
  const startDate = new Date(start)
  const endDate = new Date(end)
  const diffTime = endDate.getTime() - startDate.getTime()
  return Math.floor(diffTime / (1000 * 60 * 60 * 24))
}

export function isGreenValue(value, habit) {
  if (!habit) return false

  switch (habit.type) {
    case 'boolean':
      return (habit.positiveValue || 'yes') === 'no' ? value === 0 : value === 1
    case 'reachable':
      return value >= (habit.goal || 0)
    case 'number':
      return value > 0
    case 'percentage':
      return value >= 70
    case 'additive':
      if (habit.saveProgress) {
        return value > 0
      }
      return value >= (habit.goal || 0)
    case 'mood':
      return value >= 2
    case 'rating':
      return value >= 4
    default:
      return value > 0
  }
}

/**
 * Generates a CSS background string using color-mix.
 * Logic:
 * 0% - 75%: Red Zone (Danger)
 * 75% - 100%: Transition Zone (Danger -> Mid) [This makes 7.25 look Greenish]
 * 100%+: Green Zone (Mid -> High)
 */
export function getValueStyle(value, habit) {
  if (value === null || value === undefined) {
    return { backgroundColor: 'var(--color-heatmap-empty)' }
  }

  const goal = habit.goal || 8
  const minGood = habit.minGood || goal // Usually same as goal for 'reachable'
  
  // Helper for color-mix syntax
  // mix(colorA, colorB, percentage) -> percentage is how much of colorB
  const mix = (a, b, pct) => `color-mix(in srgb, ${a}, ${b} ${Math.round(pct * 100)}%)`

  // --- REACHABLE LOGIC (The 3-Zone Gradient) ---
  if (habit.type === 'reachable') {
    const ratio = value / goal
    // Zero
    if (ratio < 0.01) {
      return { backgroundColor: mix('red', 'var(--color-heatmap-danger)', 0.6) }
    }
    
    // 1. RED ZONE: 0 to 65% of goal
    // Visual: Faint Red -> Solid Red
    if (ratio < 0.65) {
      const intensity = 1 - (ratio / 0.65)
      return { backgroundColor: mix('transparent', 'var(--color-heatmap-danger)', 0.2 + intensity * 0.8) }
    }
    
    // 2. TRANSITION ZONE: 65% to 100%
    // Visual: Solid Red -> Solid Green (Mid)
    if (ratio < 1) {
      const progress = (ratio - 0.65) / 0.35
      return { backgroundColor: mix('var(--color-heatmap-danger)', 'var(--color-heatmap-mid)', progress) }
    }

    // 3. GREEN ZONE: 100%+
    // Visual: Mid Green -> High Green
    // 100% = Mid. 125%+ = High.
    const goodness = Math.min((value - goal) / (goal * 0.25), 1) // Fades to High over 2 hours (if goal is 8)
    return { backgroundColor: mix('var(--color-heatmap-mid)', 'var(--color-heatmap-high)', goodness) }
  }

  // --- DEFAULT LOGIC FOR OTHER HABITS ---
  const isGood = isGreenValue(value, habit)

  if (isGood) {
    return { backgroundColor: 'var(--color-heatmap-mid)' }
  } else {
    // Simple Red for non-reachable bad values
    return { backgroundColor: 'var(--color-heatmap-danger)' }
  }
}

export function calculateStreak(entries, habit) {
  if (!entries || entries.length === 0 || habit.trackStreak === false) {
    return { current: 0, longest: 0, isAntiStreak: false }
  }

  const sortedEntries = [...entries].sort((a, b) =>
    parseDate(b.date).getTime() - parseDate(a.date).getTime()
  )

  const entriesByDate = new Map()
  sortedEntries.forEach(e => entriesByDate.set(e.date, e))

  const today = formatDate(new Date())
  let currentStreak = 0
  let isAntiStreak = false
  let checkDate = today

  const todayEntry = entriesByDate.get(today)
  if (todayEntry) {
    isAntiStreak = !isGreenValue(todayEntry.value, habit)
  } else {
    const yesterday = formatDate(addDays(new Date(), -1))
    const yesterdayEntry = entriesByDate.get(yesterday)
    if (yesterdayEntry) {
      isAntiStreak = !isGreenValue(yesterdayEntry.value, habit)
      checkDate = yesterday
    }
  }

  while (true) {
    const entry = entriesByDate.get(checkDate)
    if (!entry) break

    const valueIsGreen = isGreenValue(entry.value, habit)
    if (isAntiStreak) {
      if (!valueIsGreen) {
        currentStreak++
      } else {
        break
      }
    } else {
      if (valueIsGreen) {
        currentStreak++
      } else {
        break
      }
    }
    checkDate = formatDate(addDays(parseDate(checkDate), -1))
  }

  let longest = 0
  let tempStreak = 0
  let prevDate = null

  for (const entry of sortedEntries) {
    if (isGreenValue(entry.value, habit)) {
      if (!prevDate || getDaysBetween(entry.date, prevDate) === 1) {
        tempStreak++
      } else {
        tempStreak = 1
      }
      longest = Math.max(longest, tempStreak)
      prevDate = entry.date
    } else {
      tempStreak = 0
      prevDate = null
    }
  }

  return {
    current: isAntiStreak ? -currentStreak : currentStreak,
    longest,
    isAntiStreak
  }
}

export function calculateWeeklyStreaks(entries, habit) {
  if (!entries || entries.length === 0) return 0;

  const now = new Date();
  let currentStreak = 0;
  
  const day = now.getDay();
  const diff = now.getDate() - (day === 0 ? 6 : day - 1);
  let currentWeekStart = new Date(now.setDate(diff));
  currentWeekStart.setHours(0, 0, 0, 0);

  while (true) {
    let allGreen = true;
    let hasAnyEntry = false;

    for (let d = 0; d < 7; d++) {
      const dateToCheck = new Date(currentWeekStart);
      dateToCheck.setDate(currentWeekStart.getDate() + d);
      
      if (dateToCheck > new Date()) continue;

      const dateStr = formatDate(dateToCheck);
      const entry = entries.find(e => e.date === dateStr);

      if (entry && isGreenValue(entry.value, habit)) {
        hasAnyEntry = true;
      } else {
        allGreen = false;
        break; 
      }
    }

    if (allGreen && hasAnyEntry) {
      currentStreak++;
      currentWeekStart.setDate(currentWeekStart.getDate() - 7);
    } else {
      break;
    }
  }

  return currentStreak;
}

export function calculateLevel(stats) {
  const xp = calculateXP(stats)
  let level = 1
  let xpNeeded = 100

  let totalXPForLevel = 0
  while (totalXPForLevel + xpNeeded <= xp) {
    totalXPForLevel += xpNeeded
    level++
    xpNeeded = Math.floor(xpNeeded * 1.2)
  }

  const xpInCurrentLevel = xp - totalXPForLevel
  const progress = xpInCurrentLevel / xpNeeded

  return { level, xp, xpNeeded, progress }
}

export function calculateXP(stats) {
  let xp = 0
  xp += stats.totalEntries * 10
  xp += stats.currentStreak * 5
  xp += stats.longestStreak * 20
  xp += stats.weeklyStreaks * 50
  xp += stats.perfectWeeks * 100
  xp += stats.achievements * 75
  return xp
}

export function calculateAccumulatedProgress(entries, habit, currentDate) {
  if (!habit.saveProgress || habit.type !== 'additive') return 0

  const goal = habit.goal || 10
  const sortedEntries = [...entries].sort((a, b) =>
    parseDate(b.date).getTime() - parseDate(a.date).getTime()
  )

  let accumulated = 0
  let checkDate = formatDate(addDays(parseDate(currentDate), -1))
  let foundBreak = false

  while (!foundBreak) {
    const entry = sortedEntries.find(e => e.date === checkDate)

    if (!entry) {
      if (habit.resetOnStreakBreak) {
        break
      }
      const prevDate = formatDate(addDays(parseDate(checkDate), -1))
      const prevEntry = sortedEntries.find(e => e.date === prevDate)
      if (!prevEntry) break
      checkDate = prevDate
      continue
    }

    const entryValue = entry.value || 0
    accumulated += entryValue

    if (accumulated >= goal) {
      accumulated = 0
      break
    }

    if (habit.resetOnStreakBreak && entryValue === 0) {
      accumulated = 0
      break
    }

    checkDate = formatDate(addDays(parseDate(checkDate), -1))
  }

  return accumulated
}

export function getFireIntensity(streak) {
  if (streak <= 0) return 0
  if (streak < 3) return 0.3
  if (streak < 7) return 0.5
  if (streak < 14) return 0.7
  if (streak < 30) return 0.85
  return 1
}

export function calculateStats(entries, habits, achievements) {
  let totalEntries = 0
  let longestStreak = 0
  let currentStreak = 0
  let weeklyStreaks = 0
  let perfectWeeks = 0
  let perfectMonths = 0

  habits.forEach((habit, index) => {
    const habitEntries = entries[index] || []
    totalEntries += habitEntries.filter(e => isGreenValue(e.value, habit)).length
    const streak = calculateStreak(habitEntries, habit)
    longestStreak = Math.max(longestStreak, streak.longest)
    if (!streak.isAntiStreak) {
      currentStreak = Math.max(currentStreak, streak.current)
    }
    weeklyStreaks += calculateWeeklyStreaks(habitEntries, habit)
  })

  return {
    totalEntries,
    longestStreak,
    currentStreak,
    weeklyStreaks,
    perfectWeeks,
    perfectMonths,
    habitsCreated: habits.length,
    habitsTracked: habits.length,
    achievements: achievements.length,
    earlyEntries: 0,
    nightEntries: 0,
    comebacks: 0,
    consistentWeeks: 0,
    improvements: 0,
    statsViews: 0,
    level: calculateLevel({
      totalEntries,
      currentStreak,
      longestStreak,
      weeklyStreaks,
      perfectWeeks,
      achievements: achievements.length
    }).level
  }
}