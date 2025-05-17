package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"spotify-mock-api/internal/handlers"
	"spotify-mock-api/internal/models"
	"spotify-mock-api/internal/utils"
	"time"
)

var db *gorm.DB

func main() {
	// Initialize SQLite database
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// 2) Auto‐migrate all your models
	if err := db.AutoMigrate(
		&models.Artist{},
		&models.Album{},
		&models.Song{},
		&models.Playlist{},
		&models.Podcast{},
		&models.LibraryEntry{},
		&models.User{},
	); err != nil {
		log.Fatal("migration failed:", err)
	}

	// 3) Seed default data on first run
	if err := seedDefaults(db); err != nil {
		log.Fatal("seeding defaults failed:", err)
	}

	r := gin.Default()

	// CORS wrapper
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:19006", "http://localhost:8081"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	trackH := handlers.NewTrackHandler(db)

	// Serve the same MP3 file from /media/song.mp3
	r.Static("/media", "./media")

	// Track endpoints
	r.GET("/tracks/:id", trackH.GetTrackByID)
	r.GET("/tracks/:id/audio", handlers.GetTrackAudio)

	//Playlist
	r.GET("/users/:userId/recent-playlists", handlers.GetRecentPlaylistsByUser(db))
	r.GET("/library", handlers.GetLibraryData(db))

	r.GET("/me", handlers.GetCurrentUser(db))

	// Search endpoint
	r.GET("/search", handlers.GetSearch(db))

	// Start server
	localIP := utils.GetLocalIP()
	port := ":8080"

	fmt.Printf("🚀 Server running at http://%s%s\n", localIP, port)

	http.ListenAndServe("0.0.0.0"+port, handler)
}

// define a struct matching defaults.json
type Defaults struct {
	Songs          []models.Song         `json:"songs"`
	Albums         []models.Album        `json:"albums"`
	Artists        []models.Artist       `json:"artists"`
	Playlists      []models.Playlist     `json:"playlists"`
	Podcasts       []models.Podcast      `json:"podcasts"`
	LibraryEntries []models.LibraryEntry `json:"libraryEntries"`
	Users          []models.User         `json:"users"`
}

func seedDefaults(db *gorm.DB) error {
	b, err := ioutil.ReadFile("data/defaults.json")
	if err != nil {
		return fmt.Errorf("read defaults.json: %w", err)
	}

	var defs Defaults
	if err := json.Unmarshal(b, &defs); err != nil {
		return fmt.Errorf("unmarshal defaults.json: %w", err)
	}

	// Seed Artists
	var artistCount int64
	db.Model(&models.Artist{}).Count(&artistCount)
	if artistCount == 0 {
		if err := db.Create(&defs.Artists).Error; err != nil {
			return fmt.Errorf("insert artists: %w", err)
		}
		log.Printf("seeded %d artists", len(defs.Artists))
	}

	// Seed Albums
	var albumCount int64
	db.Model(&models.Album{}).Count(&albumCount)
	if albumCount == 0 {
		if err := db.Create(&defs.Albums).Error; err != nil {
			return fmt.Errorf("insert albums: %w", err)
		}
		log.Printf("seeded %d albums", len(defs.Albums))
	}

	// Seed Songs
	var songCount int64
	db.Model(&models.Song{}).Count(&songCount)
	if songCount == 0 {
		if err := db.Create(&defs.Songs).Error; err != nil {
			return fmt.Errorf("insert songs: %w", err)
		}
		log.Printf("seeded %d songs", len(defs.Songs))
	}

	// Seed Playlists
	var plCount int64
	db.Model(&models.Playlist{}).Count(&plCount)
	if plCount == 0 {
		for i := range defs.Playlists {
			if defs.Playlists[i].LastUpdated.IsZero() {
				defs.Playlists[i].LastUpdated = time.Now()
			}
		}
		if err := db.Create(&defs.Playlists).Error; err != nil {
			return fmt.Errorf("insert playlists: %w", err)
		}
		log.Printf("seeded %d playlists", len(defs.Playlists))
	}

	// Seed Podcasts
	var pdCount int64
	db.Model(&models.Podcast{}).Count(&pdCount)
	if pdCount == 0 {
		if err := db.Create(&defs.Podcasts).Error; err != nil {
			return fmt.Errorf("insert podcasts: %w", err)
		}
		log.Printf("seeded %d podcasts", len(defs.Podcasts))
	}

	// Seed LibraryEntries
	var libCount int64
	db.Model(&models.LibraryEntry{}).Count(&libCount)
	if libCount == 0 {
		if err := db.Create(&defs.LibraryEntries).Error; err != nil {
			return fmt.Errorf("insert library entries: %w", err)
		}
		log.Printf("seeded %d library entries", len(defs.LibraryEntries))
	}

	// Seed User
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		if _, err := os.Stat("media/avatar.jpg"); err != nil {
			log.Println("warning: avatar.jpg not found; using placeholder")
		}
		if err := db.Create(&defs.Users).Error; err != nil {
			return fmt.Errorf("insert users: %w", err)
		}
		log.Println("seeded default user profile")
	}

	return nil
}
