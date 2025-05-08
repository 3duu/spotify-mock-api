package spotify_mock_api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"spotify-mock-api/internal/handlers"
	"spotify-mock-api/internal/models"
)

var db *gorm.DB

func main() {
	// Initialize SQLite database
	var err error
	db, err = gorm.Open(sqlite.Open("songs.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Migrate the Song schema
	if err := db.AutoMigrate(&models.Song{}); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	// Load songs from JSON into database
	loadSongs("data/songs.json")

	r := gin.Default()

	// Serve the same MP3 file from /media/song.mp3
	r.Static("/media", "./media")

	// Track endpoints
	r.GET("/tracks/:id", handlers.GetTrackByID)
	r.GET("/tracks/:id/audio", handlers.GetTrackAudio)

	// Start server
	r.Run(":8080")
}

// loadSongs reads a JSON array of songs and inserts them into the DB
func loadSongs(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic("unable to open songs.json: " + err.Error())
	}
	defer file.Close()

	var songs []models.Song
	if err := json.NewDecoder(file).Decode(&songs); err != nil {
		panic("failed to decode JSON: " + err.Error())
	}

	for _, s := range songs {
		db.FirstOrCreate(&models.Song{}, s)
	}
}
