package models

import "time"

type HabitType string

const (
	HabitTypeBoolean    HabitType = "boolean"
	HabitTypeReachable  HabitType = "reachable"
	HabitTypeNumber     HabitType = "number"
	HabitTypePercentage HabitType = "percentage"
	HabitTypeAdditive   HabitType = "additive"
	HabitTypeMood       HabitType = "mood"
	HabitTypeRating     HabitType = "rating"
)

type Habit struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Type               HabitType `json:"type"`
	Goal               float64   `json:"goal,omitempty"`
	PositiveValue      string    `json:"positiveValue,omitempty"`
	SaveProgress       bool      `json:"saveProgress"`
	ResetOnStreakBreak bool      `json:"resetOnStreakBreak"`
	TrackStreak        bool      `json:"trackStreak"`
	Priority           bool      `json:"priority"`
	Private            bool      `json:"private"`
	NotificationsOn    bool      `json:"notificationsOn"`
	NotificationTime   string    `json:"notificationTime,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
	Order              int       `json:"order"`
	ExternalSourceID   string    `json:"externalSourceId,omitempty"`
}

type Entry struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habitId"`
	Date      string    `json:"date"`
	Value     float64   `json:"value"`
	Note      string    `json:"note,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	ID        string    `json:"id"`
	IsOwner   bool      `json:"isOwner"`
	CreatedAt time.Time `json:"createdAt"`
}

type JoinLink struct {
	Token     string    `json:"token"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"createdAt"`
}

type Achievement struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Type        string    `json:"type"`
	UnlockedAt  time.Time `json:"unlockedAt"`
}

type Stats struct {
	CurrentStreak   int     `json:"currentStreak"`
	LongestStreak   int     `json:"longestStreak"`
	WeeklyStreaks   int     `json:"weeklyStreaks"`
	TotalEntries    int     `json:"totalEntries"`
	Level           int     `json:"level"`
	Experience      float64 `json:"experience"`
}
