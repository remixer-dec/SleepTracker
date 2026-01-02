package db

import (
	"encoding/json"
	"errors"
	"sleeptracker/internal/models"
	"time"

	"github.com/tidwall/buntdb"
)

var ErrNotFound = errors.New("not found")

type Database struct {
	db *buntdb.DB
}

func New(path string) (*Database, error) {
	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}

	err = db.CreateIndex("habits", "habit:*", buntdb.IndexJSON("order"))
	if err != nil && err != buntdb.ErrIndexExists {
		return nil, err
	}

	err = db.CreateIndex("entries_by_habit", "entry:*", buntdb.IndexJSON("habitId"))
	if err != nil && err != buntdb.ErrIndexExists {
		return nil, err
	}

	err = db.CreateIndex("entries_by_date", "entry:*", buntdb.IndexJSON("date"))
	if err != nil && err != buntdb.ErrIndexExists {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateHabit(habit *models.Habit) error {
	data, err := json.Marshal(habit)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("habit:"+habit.ID, string(data), nil)
		return err
	})
}

func (d *Database) GetHabit(id string) (*models.Habit, error) {
	var habit models.Habit

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("habit:" + id)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &habit)
	})

	if err != nil {
		return nil, err
	}
	return &habit, nil
}

func (d *Database) GetAllHabits() ([]models.Habit, error) {
	var habits []models.Habit

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("habits", func(key, value string) bool {
			var habit models.Habit
			if err := json.Unmarshal([]byte(value), &habit); err == nil {
				habits = append(habits, habit)
			}
			return true
		})
	})

	return habits, err
}

func (d *Database) UpdateHabit(habit *models.Habit) error {
	data, err := json.Marshal(habit)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Get("habit:" + habit.ID)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		_, _, err = tx.Set("habit:"+habit.ID, string(data), nil)
		return err
	})
}

func (d *Database) DeleteHabit(id string) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete("habit:" + id)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		return err
	})
}

func (d *Database) CreateEntry(entry *models.Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("entry:"+entry.ID, string(data), nil)
		return err
	})
}

func (d *Database) GetEntry(id string) (*models.Entry, error) {
	var entry models.Entry

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("entry:" + id)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &entry)
	})

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (d *Database) GetEntriesByHabit(habitID string) ([]models.Entry, error) {
	var entries []models.Entry

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if len(key) > 6 && key[:6] == "entry:" {
				var entry models.Entry
				if err := json.Unmarshal([]byte(value), &entry); err == nil {
					if entry.HabitID == habitID {
						entries = append(entries, entry)
					}
				}
			}
			return true
		})
	})

	return entries, err
}

func (d *Database) GetEntriesByHabitAndDateRange(habitID, startDate, endDate string) ([]models.Entry, error) {
	var entries []models.Entry

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if len(key) > 6 && key[:6] == "entry:" {
				var entry models.Entry
				if err := json.Unmarshal([]byte(value), &entry); err == nil {
					if entry.HabitID == habitID && entry.Date >= startDate && entry.Date <= endDate {
						entries = append(entries, entry)
					}
				}
			}
			return true
		})
	})

	return entries, err
}

func (d *Database) GetEntryByHabitAndDate(habitID, date string) (*models.Entry, error) {
	var found *models.Entry

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if len(key) > 6 && key[:6] == "entry:" {
				var entry models.Entry
				if err := json.Unmarshal([]byte(value), &entry); err == nil {
					if entry.HabitID == habitID && entry.Date == date {
						found = &entry
						return false
					}
				}
			}
			return true
		})
	})

	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrNotFound
	}
	return found, nil
}

func (d *Database) UpdateEntry(entry *models.Entry) error {
	entry.UpdatedAt = time.Now()
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Get("entry:" + entry.ID)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		_, _, err = tx.Set("entry:"+entry.ID, string(data), nil)
		return err
	})
}

func (d *Database) DeleteEntry(id string) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete("entry:" + id)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		return err
	})
}

func (d *Database) DeleteEntriesByHabit(habitID string) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		var keysToDelete []string
		err := tx.Ascend("", func(key, value string) bool {
			if len(key) > 6 && key[:6] == "entry:" {
				var entry models.Entry
				if err := json.Unmarshal([]byte(value), &entry); err == nil {
					if entry.HabitID == habitID {
						keysToDelete = append(keysToDelete, key)
					}
				}
			}
			return true
		})
		if err != nil {
			return err
		}
		for _, key := range keysToDelete {
			if _, err := tx.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *Database) SaveJoinLink(link *models.JoinLink) error {
	data, err := json.Marshal(link)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("joinlink:"+link.Token, string(data), nil)
		return err
	})
}

func (d *Database) GetJoinLink(token string) (*models.JoinLink, error) {
	var link models.JoinLink

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("joinlink:" + token)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &link)
	})

	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (d *Database) MarkJoinLinkUsed(token string) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		val, err := tx.Get("joinlink:" + token)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}

		var link models.JoinLink
		if err := json.Unmarshal([]byte(val), &link); err != nil {
			return err
		}

		link.Used = true
		data, err := json.Marshal(link)
		if err != nil {
			return err
		}

		_, _, err = tx.Set("joinlink:"+token, string(data), nil)
		return err
	})
}

func (d *Database) SaveUser(user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("user:"+user.ID, string(data), nil)
		return err
	})
}

func (d *Database) GetUser(id string) (*models.User, error) {
	var user models.User

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("user:" + id)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(val), &user)
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) SaveAchievement(achievement *models.Achievement) error {
	data, err := json.Marshal(achievement)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("achievement:"+achievement.ID, string(data), nil)
		return err
	})
}

func (d *Database) GetAchievementsByUser(userID string) ([]models.Achievement, error) {
	var achievements []models.Achievement

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if len(key) > 12 && key[:12] == "achievement:" {
				var achievement models.Achievement
				if err := json.Unmarshal([]byte(value), &achievement); err == nil {
					if achievement.UserID == userID {
						achievements = append(achievements, achievement)
					}
				}
			}
			return true
		})
	})

	return achievements, err
}

func (d *Database) GetOwnerExists() (bool, error) {
	exists := false

	err := d.db.View(func(tx *buntdb.Tx) error {
		return tx.Ascend("", func(key, value string) bool {
			if len(key) > 5 && key[:5] == "user:" {
				var user models.User
				if err := json.Unmarshal([]byte(value), &user); err == nil {
					if user.IsOwner {
						exists = true
						return false
					}
				}
			}
			return true
		})
	})

	return exists, err
}

func (d *Database) GetConfig(key string) (string, error) {
	var value string

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get("config:" + key)
		if err == buntdb.ErrNotFound {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		value = val
		return nil
	})

	return value, err
}

func (d *Database) SetConfig(key, value string) error {
	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set("config:"+key, value, nil)
		return err
	})
}
