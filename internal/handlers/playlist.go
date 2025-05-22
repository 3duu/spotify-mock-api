package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"spotify-mock-api/internal/models"
	"strconv"
	"time"
)

type PlaylistResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Cover       string    `json:"cover"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
}

// PlaylistDetailResponse is the full payload for GET /playlists/:id
type PlaylistDetailResponse struct {
	ID         int             `json:"id"`
	Title      string          `json:"title"`
	Cover      string          `json:"cover"`
	OwnerName  string          `json:"ownerName"`
	OwnerImage string          `json:"ownerImage"`
	Duration   string          `json:"duration"` // e.g. "5h 59m"
	Tracks     []TrackResponse `json:"tracks"`
}

// GetRecentPlaylistsByUser returns up to 10 most‐recently updated playlists
// for the given user ID, mapping them into PlaylistResponse.
func GetRecentPlaylistsByUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam := c.Param("userId")
		if userIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}

		// Parse userId to int
		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId must be an integer"})
			return
		}

		// Fetch the 10 most recently updated playlists for this user
		var pls []models.Playlist
		if err := db.
			Where("user_id = ?", userID).
			Order("last_updated DESC").
			Limit(10).
			Find(&pls).
			Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load playlists"})
			return
		}

		// Map to API response type
		resp := make([]PlaylistResponse, len(pls))
		for i, p := range pls {
			subtitle := fmt.Sprintf("Playlist • %d tracks", 0) // Replace 0 if you count tracks
			resp[i] = PlaylistResponse{
				ID:       p.ID,
				Title:    p.Title,
				Subtitle: subtitle,
				Cover:    p.Cover,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}

// GetPlaylistDetail loads one playlist and returns its full detail
func GetPlaylistDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID := c.Param("id")

		// 1) Load playlist, its owner, and its tracks (+ artists + cover fields)
		var pl models.Playlist
		err := db.
			Preload("Owner"). // assuming Playlist has an Owner   field -> User
			/*Preload("Tracks.Artist"). // assuming Playlist.Tracks []*Track
			Preload("Tracks").        // to get Track.Cover, Track.VideoFlag, Track.DownloadedFlag*/
			First(&pl, "id = ?", playlistID).Error
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}

		// 2) Build the slice of TrackResponse
		tracks := make([]TrackResponse, len(pl.Songs))
		var totalSec int
		for i, t := range pl.Songs {
			// sum up durations
			totalSec += t.Duration

			tracks[i] = TrackResponse{
				ID:         t.ID,
				Title:      t.Title,
				Artist:     t.Artist.Name,
				AlbumArt:   t.Album.Cover,
				Downloaded: false,
			}
		}

		// 3) Convert totalSec into "5h 59m"
		h := totalSec / 3600
		m := (totalSec % 3600) / 60
		durationStr := fmt.Sprintf("%dh %02dm", h, m)

		// 4) Assemble the response
		resp := PlaylistDetailResponse{
			ID:         pl.ID,
			Title:      pl.Title,
			Cover:      pl.Cover,
			OwnerName:  pl.Owner.Name,
			OwnerImage: pl.Owner.Image,
			Duration:   durationStr,
			Tracks:     tracks,
		}

		c.JSON(http.StatusOK, resp)
	}
}
