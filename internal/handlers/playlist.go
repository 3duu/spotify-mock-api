package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"spotify-mock-api/internal/models"
	"time"
)

type PlaylistResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	IconURL     string    `json:"icon"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
}

// GetRecentPlaylistsByUser returns up to 10 most‐recently updated playlists
// for the given user ID, mapping them into PlaylistResponse.
func GetRecentPlaylistsByUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
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
			// If you have a track‐count relationship, replace 0 with the real count
			subtitle := fmt.Sprintf("Playlist • %d tracks", 0)
			resp[i] = PlaylistResponse{
				ID:       p.ID,
				Title:    p.Title, // or p.Title if your model uses Title
				Subtitle: subtitle,
				IconURL:  p.Cover, // or p.IconURL
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}
