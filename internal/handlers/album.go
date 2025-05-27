package handlers

import (
	"net/http"
	"spotify-mock-api/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AlbumDetailResponse matches the “PlaylistDetail” shape on the frontend:
type AlbumDetailResponse struct {
	ID         int             `json:"id" :"id"`
	Title      string          `json:"title" :"title"`
	Cover      string          `json:"cover" :"cover"`
	OwnerName  string          `json:"ownerName" :"owner_name"`   // here: artist name
	OwnerImage string          `json:"ownerImage" :"owner_image"` // could be blank or artist image
	Duration   string          `json:"duration" :"duration"`      // total playtime, e.g. "42m 15s"
	Tracks     []TrackResponse `json:"tracks" :"tracks"`
}

// GetAlbumDetail loads an album and its tracks + artist, then returns a unified response.
func GetAlbumDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// parse album ID
		rawID := c.Param("id")
		albumID, err := strconv.Atoi(rawID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid album ID"})
			return
		}

		// load album + artist + songs
		var album models.Album
		if err := db.
			Preload("Artist").
			Preload("Songs.Artist").
			First(&album, "album_id = ?", albumID).
			Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "album not found"})
			return
		}

		// build total duration
		totalSecs := 0
		for _, s := range album.Songs {
			totalSecs += s.Duration
		}
		dur := time.Duration(totalSecs) * time.Second
		durationStr := dur.Truncate(time.Second).String() // "1h2m3s"

		// map tracks
		tracks := make([]TrackResponse, len(album.Songs))
		for i, s := range album.Songs {
			tracks[i] = TrackResponse{
				ID:         s.ID,
				Title:      s.Title,
				Artist:     s.Artist.Name,
				AlbumArt:   album.Cover,
				Duration:   s.Duration,
				Downloaded: false,
			}
		}

		resp := AlbumDetailResponse{
			ID:         album.AlbumId,
			Title:      album.Title,
			Cover:      album.Cover,
			OwnerName:  album.Artist.Name,
			OwnerImage: "", // if you have an artist image, fill here
			Duration:   durationStr,
			Tracks:     tracks,
		}
		c.JSON(http.StatusOK, resp)
	}
}
