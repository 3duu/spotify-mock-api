package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"spotify-mock-api/internal/models"
)

type PlaylistResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	IconURL  string `json:"icon"`
}

// GetRecentPlaylists queries the playlists table and returns them
func GetRecentPlaylists(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("GetRecentPlaylists")

		// 1) Fetch all playlists from the DB
		var pls []models.Playlist
		if err := db.Find(&pls).Error; err != nil {
			log.Printf("DB error fetching playlists: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load playlists"})
			return
		}

		// 2) Map to the API response type
		res := make([]PlaylistResponse, len(pls))
		for i, p := range pls {
			// You can customize Subtitle as you like (e.g. count of songs)
			subtitle := fmt.Sprintf("Playlist • %d tracks", 0)
			res[i] = PlaylistResponse{
				ID:       p.ID,
				Title:    p.Title,
				Subtitle: subtitle,
				IconURL:  p.Cover,
			}
		}

		// 3) Return JSON
		c.JSON(http.StatusOK, res)
	}
}
