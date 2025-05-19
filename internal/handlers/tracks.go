package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"path/filepath"
	"spotify-mock-api/internal/models"
)

type TrackHandler struct {
	DB *gorm.DB
}

func NewTrackHandler(db *gorm.DB) *TrackHandler {
	return &TrackHandler{DB: db}
}

// getTrackByID retrieves song metadata and returns it along with an audio URL
func (h *TrackHandler) GetTrackByID(c *gin.Context) {
	var song models.Song // make sure Song is imported or duplicated here
	if err := h.DB.First(&song, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		return
	}

	// Determine request scheme for correct URL
	/*scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}*/
	audioURL := "/media/song.mp3" //fmt.Sprintf("%s://%s/media/song.mp3", scheme, c.Request.Host)

	c.JSON(http.StatusOK, gin.H{
		"id":        song.ID,
		"title":     song.Title,
		"artist":    song.Artist.Name,
		"artist_id": song.Artist.ArtistId,
		"album_art": "/media/album-art.jpg",
		"album_id":  song.AlbumID,
		"duration":  song.Duration,
		"audio_url": audioURL,
	})
}

// getTrackAudio streams the MP3 file (always the same file)
func GetTrackAudio(c *gin.Context) {
	log.Print("Get Audio")
	c.File(filepath.Join("media", "song.mp3"))
}

func GetArtistByID(c *gin.Context) {
	log.Print("Get Artist")
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "Mock Artist"})
}

func GetAlbumByID(c *gin.Context) {
	log.Print("Get Album")
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "title": "Mock Album"})
}

func GetPlaylistByID(c *gin.Context) {
	log.Print("Get Playlist")
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id, "name": "Mock Playlist"})
}

func Search(c *gin.Context) {
	log.Print("Search")
	query := c.Query("q")
	c.JSON(http.StatusOK, gin.H{"query": query, "results": []string{}})
}

func GetProfile(c *gin.Context) {
	log.Println("Get Profile")
	c.JSON(http.StatusOK, gin.H{"id": "user123", "display_name": "Mock User"})
}

func Login(c *gin.Context) {
	log.Print("Login User")
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
