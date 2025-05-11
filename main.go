package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"net/http"
	"os"
	"spotify-mock-api/internal/handlers"
	"spotify-mock-api/internal/models"
	"spotify-mock-api/internal/utils"
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
	r.GET("/library/recent-playlists", handlers.GetRecentPlaylists)

	r.GET("/me", handlers.GetCurrentUser)

	// Start server
	localIP := utils.GetLocalIP()
	port := ":8080"

	fmt.Printf("🚀 Server running at http://%s%s\n", localIP, port)

	http.ListenAndServe("0.0.0.0"+port, handler)
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
