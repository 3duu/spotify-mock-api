package handlers

import (
	"net/http"
	"spotify-mock-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArtistDetailResponse matches the same shape
type ArtistDetailResponse struct {
	ID         int             `json:"id"`
	Title      string          `json:"title"`      // artist name
	Cover      string          `json:"cover"`      // could be a default artist image
	OwnerName  string          `json:"ownerName"`  // blank or same as title
	OwnerImage string          `json:"ownerImage"` // artist avatar
	Duration   string          `json:"duration"`   // sum of track durations
	Tracks     []TrackResponse `json:"tracks"`
}

// GetArtistDetail loads an artist, their top songs (or all songs), then returns unified response.
func GetArtistDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// parse artist ID
		rawID := c.Param("id")
		artistID, err := strconv.Atoi(rawID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid artist ID"})
			return
		}

		// fetch artist + songs + album (for art)
		var artist models.Artist
		if err := db.
			Preload("Songs.Album").
			Preload("Songs").
			First(&artist, "artist_id = ?", artistID).
			Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "artist not found"})
			return
		}

		// total duration
		totalSecs := 0
		for _, s := range artist.Songs {
			totalSecs += s.Duration
		}
		dur := time.Duration(totalSecs) * time.Second
		durationStr := dur.Truncate(time.Second).String()

		// map tracks
		tracks := make([]TrackResponse, len(artist.Songs))
		for i, s := range artist.Songs {
			tracks[i] = TrackResponse{
				ID:         s.ID,
				Title:      s.Title,
				Artist:     artist.Name,
				AlbumArt:   s.Album.Cover,
				Duration:   s.Duration,
				Downloaded: false,
			}
		}

		resp := ArtistDetailResponse{
			ID:         artist.ArtistId,
			Title:      artist.Name,
			Cover:      "", // optional artist cover
			OwnerName:  "",
			OwnerImage: artist.Image, // if you store one
			Duration:   durationStr,
			Tracks:     tracks,
		}
		c.JSON(http.StatusOK, resp)
	}
}
