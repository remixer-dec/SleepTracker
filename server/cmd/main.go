package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sleeptracker/internal/auth"
	"sleeptracker/internal/db"
	"sleeptracker/internal/handlers"
	"sleeptracker/internal/middleware"
	"strings"
)

func main() {
	port := flag.String("port", "8080", "Server port")
	dbPath := flag.String("db", "data/sleeptracker.db", "Database file path")
	staticDir := flag.String("static", "static", "Static files directory")
	createJoinLink := flag.Bool("create-join-link", false, "Create a one-time join link")
	setExternalKey := flag.String("set-external-key", "", "Set external API key for integrations")
	flag.Parse()

	if err := os.MkdirAll(filepath.Dir(*dbPath), 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	secretKey, err := database.GetConfig("jwt_secret")
	if err == db.ErrNotFound {
		secretKey, err = auth.GenerateSecretKey()
		if err != nil {
			log.Fatalf("Failed to generate secret key: %v", err)
		}
		if err := database.SetConfig("jwt_secret", secretKey); err != nil {
			log.Fatalf("Failed to save secret key: %v", err)
		}
		log.Println("Generated new JWT secret key")
	} else if err != nil {
		log.Fatalf("Failed to get secret key: %v", err)
	}

	authenticator := auth.New(secretKey)

	if *setExternalKey != "" {
		if err := database.SetConfig("external_api_key", *setExternalKey); err != nil {
			log.Fatalf("Failed to set external API key: %v", err)
		}
		fmt.Println("External API key set successfully")
		return
	}

	joinTokenPath := filepath.Join(filepath.Dir(*dbPath), "join_token")

	if *createJoinLink {
		token, err := auth.GenerateJoinToken()
		if err != nil {
			log.Fatalf("Failed to generate join token: %v", err)
		}

		if err := os.WriteFile(joinTokenPath, []byte(token), 0600); err != nil {
			log.Fatalf("Failed to write join token: %v", err)
		}

		fmt.Printf("Join token created: %s\n", token)
		fmt.Printf("Use: /api/join?token=%s\n", token)
		return
	}

	h := handlers.New(database, authenticator)
	authMiddleware := middleware.NewAuthMiddleware(authenticator)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/habits", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetHabits(w, r)
		case http.MethodPost:
			h.CreateHabit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/habits/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetHabit(w, r)
		case http.MethodPut:
			h.UpdateHabit(w, r)
		case http.MethodDelete:
			h.DeleteHabit(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetEntries(w, r)
		case http.MethodPost:
			h.CreateEntry(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/entries/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			h.UpdateEntry(w, r)
		case http.MethodDelete:
			h.DeleteEntry(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/join", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Join(w, r)
	})

	mux.HandleFunc("/api/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.GetMe(w, r)
	})

	mux.HandleFunc("/api/achievements", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetAchievements(w, r)
		case http.MethodPost:
			h.UnlockAchievement(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/external/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ExternalUpdate(w, r)
	})

	if _, err := os.Stat(*staticDir); err == nil {
		fileServer := http.FileServer(http.Dir(*staticDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(*staticDir, r.URL.Path)

			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.NotFound(w, r)
				return
			}

			if _, err := os.Stat(path); os.IsNotExist(err) || isDir(path) {
				if !strings.Contains(r.URL.Path, ".") || isDir(path) {
					http.ServeFile(w, r, filepath.Join(*staticDir, "index.html"))
					return
				}
			}

			fileServer.ServeHTTP(w, r)
		})
	}

	handler := authMiddleware.Authenticate(mux)

	log.Printf("Server starting on port %s", *port)
	if err := http.ListenAndServe(":"+*port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&fs.ModeDir != 0
}
