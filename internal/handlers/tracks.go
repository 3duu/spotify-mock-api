package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"spotify-mock-api/internal/models"
)

var db *gorm.DB

// getTrackByID retrieves song metadata and returns it along with an audio URL
func GetTrackByID(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	if err := db.First(&song, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	// Determine request scheme for correct URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	audioURL := fmt.Sprintf("%s://%s/media/song.mp3", scheme, c.Request.Host)

	c.JSON(http.StatusOK, gin.H{
		"id":        song.ID,
		"title":     song.Title,
		"artist":    song.Artist,
		"album":     song.Album,
		"duration":  song.Duration,
		"audio_url": audioURL,
	})
}

// getTrackAudio streams the MP3 file (always the same file)
func GetTrackAudio(c *gin.Context) {
	c.File(filepath.Join("media", "song.mp3"))
}

func GetArtistByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "Mock Artist"})
}

func GetAlbumByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "title": "Mock Album"})
}

func GetPlaylistByID(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "Mock Playlist"})
}

func Search(c *gin.Context) {
	query := c.Query("q")
	c.JSON(http.StatusOK, gin.H{"query": query, "results": []string{}})
}

func GetProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id": "user123", "display_name": "Mock User"})
}

func Login(c *gin.Context) {
	// TODO: authenticate and return mock token
	c.JSON(http.StatusOK, gin.H{"access_token": "mock_token"})
}

func Logout(c *gin.Context) {
	// TODO: invalidate token
	c.Status(http.StatusNoContent)
}

func Play(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Pause(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Next(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Previous(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
