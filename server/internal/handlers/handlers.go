package handlers

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"os"
	"sleeptracker/internal/auth"
	"sleeptracker/internal/db"
	"sleeptracker/internal/middleware"
	"sleeptracker/internal/models"
	"strings"
	"time"
)

type Handlers struct {
	db   *db.Database
	auth *auth.Auth
}

func New(database *db.Database, a *auth.Auth) *Handlers {
	return &Handlers{
		db:   database,
		auth: a,
	}
}

func (h *Handlers) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (h *Handlers) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}

func (h *Handlers) clearOtherPriorities(exceptID string) error {
	habits, err := h.db.GetAllHabits()
	if err != nil {
		return err
	}

	for i := range habits {
		if habits[i].ID != exceptID && habits[i].Priority {
			habits[i].Priority = false
			if err := h.db.UpdateHabit(&habits[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (h *Handlers) GetHabits(w http.ResponseWriter, r *http.Request) {
	habits, err := h.db.GetAllHabits()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get habits")
		return
	}
	if habits == nil {
		habits = []models.Habit{}
	}
	h.respondJSON(w, http.StatusOK, habits)
}

func (h *Handlers) CreateHabit(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	var habit models.Habit
	if err := json.NewDecoder(r.Body).Decode(&habit); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := auth.GenerateUserID()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate ID")
		return
	}

	habit.ID = id
	habit.CreatedAt = time.Now()

	if habit.Priority {
		if err := h.clearOtherPriorities(id); err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to clear other priorities")
			return
		}
	}

	if err := h.db.CreateHabit(&habit); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create habit")
		return
	}

	h.respondJSON(w, http.StatusCreated, habit)
}

func (h *Handlers) GetHabit(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Habit ID required")
		return
	}

	habit, err := h.db.GetHabit(id)
	if err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Habit not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get habit")
		return
	}

	h.respondJSON(w, http.StatusOK, habit)
}

func (h *Handlers) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Habit ID required")
		return
	}

	existing, err := h.db.GetHabit(id)
	if err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Habit not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get habit")
		return
	}

	var updates models.Habit
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updates.ID = existing.ID
	updates.Type = existing.Type
	updates.CreatedAt = existing.CreatedAt

	if updates.Priority {
		if err := h.clearOtherPriorities(id); err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to clear other priorities")
			return
		}
	}

	if err := h.db.UpdateHabit(&updates); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update habit")
		return
	}

	h.respondJSON(w, http.StatusOK, updates)
}

func (h *Handlers) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/habits/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Habit ID required")
		return
	}

	if err := h.db.DeleteEntriesByHabit(id); err != nil && err != db.ErrNotFound {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete entries")
		return
	}

	if err := h.db.DeleteHabit(id); err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Habit not found")
		return
	} else if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete habit")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handlers) GetEntries(w http.ResponseWriter, r *http.Request) {
	habitID := r.URL.Query().Get("habitId")
	startDate := r.URL.Query().Get("start")
	endDate := r.URL.Query().Get("end")

	if habitID == "" {
		h.respondError(w, http.StatusBadRequest, "habitId query parameter required")
		return
	}

	var entries []models.Entry
	var err error

	if startDate != "" && endDate != "" {
		entries, err = h.db.GetEntriesByHabitAndDateRange(habitID, startDate, endDate)
	} else {
		entries, err = h.db.GetEntriesByHabit(habitID)
	}

	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get entries")
		return
	}

	if entries == nil {
		entries = []models.Entry{}
	}

	h.respondJSON(w, http.StatusOK, entries)
}

func (h *Handlers) CreateEntry(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	var entry models.Entry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existing, err := h.db.GetEntryByHabitAndDate(entry.HabitID, entry.Date)
	if err == nil {
		existing.Value = entry.Value
		existing.Note = entry.Note
		existing.UpdatedAt = time.Now()
		if err := h.db.UpdateEntry(existing); err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to update entry")
			return
		}
		h.respondJSON(w, http.StatusOK, existing)
		return
	}

	id, err := auth.GenerateUserID()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate ID")
		return
	}

	entry.ID = id
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	if err := h.db.CreateEntry(&entry); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create entry")
		return
	}

	h.respondJSON(w, http.StatusCreated, entry)
}

func (h *Handlers) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/entries/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Entry ID required")
		return
	}

	existing, err := h.db.GetEntry(id)
	if err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Entry not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get entry")
		return
	}

	var updates models.Entry
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existing.Value = updates.Value
	existing.Note = updates.Note

	if err := h.db.UpdateEntry(existing); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to update entry")
		return
	}

	h.respondJSON(w, http.StatusOK, existing)
}

func (h *Handlers) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/entries/")
	if id == "" {
		h.respondError(w, http.StatusBadRequest, "Entry ID required")
		return
	}

	if err := h.db.DeleteEntry(id); err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Entry not found")
		return
	} else if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete entry")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handlers) Join(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		h.respondError(w, http.StatusBadRequest, "Token required")
		return
	}

	tokenPath := "data/join_token"
	info, err := os.Stat(tokenPath)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "Invalid token")
		return
	}
	if time.Since(info.ModTime()) > 24*time.Hour {
		os.Remove(tokenPath)
		h.respondError(w, http.StatusNotFound, "Invalid token")
		return
	}
	tokenBytes, err := os.ReadFile(tokenPath)
	os.Remove(tokenPath)
	tokenFromFile := strings.TrimSpace(string(tokenBytes))

	if err != nil || subtle.ConstantTimeCompare([]byte(tokenFromFile), []byte(token)) != 1 {
		h.respondError(w, http.StatusNotFound, "Invalid token")
		return
	}

	userID, err := auth.GenerateUserID()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate user ID")
		return
	}

	user := &models.User{
		ID:        userID,
		IsOwner:   true,
		CreatedAt: time.Now(),
	}

	if err := h.db.SaveUser(user); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	jwtToken, err := h.auth.GenerateToken(userID, true)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate JWT")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   365 * 24 * 60 * 60,
	})

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"token":   jwtToken,
		"userId":  userID,
		"isOwner": true,
	})
}

func (h *Handlers) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	isOwner := middleware.IsOwner(r.Context())

	h.respondJSON(w, http.StatusOK, map[string]interface{}{
		"userId":        userID,
		"isOwner":       isOwner,
		"authenticated": userID != "",
	})
}

func (h *Handlers) GetAchievements(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		h.respondJSON(w, http.StatusOK, []models.Achievement{})
		return
	}

	achievements, err := h.db.GetAchievementsByUser(userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get achievements")
		return
	}

	if achievements == nil {
		achievements = []models.Achievement{}
	}

	h.respondJSON(w, http.StatusOK, achievements)
}

func (h *Handlers) UnlockAchievement(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsOwner(r.Context()) {
		h.respondError(w, http.StatusForbidden, "Owner access required")
		return
	}

	userID := middleware.GetUserID(r.Context())

	var req struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := auth.GenerateUserID()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate ID")
		return
	}

	achievement := &models.Achievement{
		ID:         id,
		UserID:     userID,
		Type:       req.Type,
		UnlockedAt: time.Now(),
	}

	if err := h.db.SaveAchievement(achievement); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to save achievement")
		return
	}

	h.respondJSON(w, http.StatusCreated, achievement)
}

func (h *Handlers) ExternalUpdate(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		h.respondError(w, http.StatusUnauthorized, "API key required")
		return
	}

	storedKey, err := h.db.GetConfig("external_api_key")
	if err != nil || subtle.ConstantTimeCompare([]byte(storedKey), []byte(apiKey)) != 1 {
		h.respondError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	var req struct {
		HabitID  string  `json:"habitId"`
		Date     string  `json:"date"`
		Value    float64 `json:"value"`
		SourceID string  `json:"sourceId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	habit, err := h.db.GetHabit(req.HabitID)
	if err == db.ErrNotFound {
		h.respondError(w, http.StatusNotFound, "Habit not found")
		return
	}
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to get habit")
		return
	}

	if habit.ExternalSourceID != "" && habit.ExternalSourceID != req.SourceID {
		h.respondError(w, http.StatusForbidden, "Habit is linked to different source")
		return
	}

	existing, err := h.db.GetEntryByHabitAndDate(req.HabitID, req.Date)
	if err == nil {
		existing.Value = req.Value
		existing.UpdatedAt = time.Now()
		if err := h.db.UpdateEntry(existing); err != nil {
			h.respondError(w, http.StatusInternalServerError, "Failed to update entry")
			return
		}
		h.respondJSON(w, http.StatusOK, existing)
		return
	}

	id, err := auth.GenerateUserID()
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to generate ID")
		return
	}

	entry := &models.Entry{
		ID:        id,
		HabitID:   req.HabitID,
		Date:      req.Date,
		Value:     req.Value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.CreateEntry(entry); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to create entry")
		return
	}

	h.respondJSON(w, http.StatusCreated, entry)
}
