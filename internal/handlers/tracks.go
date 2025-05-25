package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"path/filepath"
	"spotify-mock-api/internal/models"
	"strconv"
)

type TrackHandler struct {
	DB *gorm.DB
}

func NewTrackHandler(db *gorm.DB) *TrackHandler {
	return &TrackHandler{DB: db}
}

// getTrackByID retrieves song metadata and returns it along with an audio URL
func (h *TrackHandler) GetTrackByID(c *gin.Context) {
	var song models.Song
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.DB.
		Preload("Artist").
		Preload("Album").
		First(&song, "id = ?", id).
		Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	// build your audio URL
	audioURL := "/media/song.mp3"

	// read where this play came from:
	// e.g. GET /tracks/123?origin=playlist&originId=42
	originType := c.Query("origin")
	originID, _ := strconv.Atoi(c.Query("originId"))

	// record the play, including its origin
	h.DB.Create(&models.RecentPlay{
		UserID:      1,
		Type:        originType,
		ReferenceID: id,
		OriginID:    originID,
	})

	// return the metadata
	c.JSON(http.StatusOK, gin.H{
		"id":        song.ID,
		"title":     song.Title,
		"artist":    song.Artist.Name,
		"artist_id": song.Artist.ArtistId,
		"album_art": song.Album.Cover,
		"album_id":  song.AlbumID,
		"duration":  song.Duration,
		"audio_url": audioURL,
		"color":     song.Color,
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
