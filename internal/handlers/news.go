package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"

	"spotify-mock-api/internal/models"
)

// NewsletterResponse mirrors the TS Newsletter interface:
// id, title, subtitle, image (URL), type.
type NewsletterResponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Image    string `json:"image"`
	Type     string `json:"type"` // e.g. "SONG"|"ALBUM"|"PODCAST"|"ARTIST"
}

// GetNewsletters returns all newsletters from the DB.
// Frontend calls GET /newsletters :contentReference[oaicite:3]{index=3}.
func GetNewsletters(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var nl []models.Newsletter
		if err := db.Find(&nl).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load newsletters"})
			return
		}

		resp := make([]NewsletterResponse, len(nl))
		for i, n := range nl {
			resp[i] = NewsletterResponse{
				ID:       n.ID,
				Title:    n.Title,
				Subtitle: n.Subtitle,
				Image:    n.ImageURL,
				Type:     n.Type,
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}

// GetAllPlaylists returns every playlist in your DB
func GetAllPlaylists(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pls []models.Playlist
		if err := db.Find(&pls).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load playlists"})
			return
		}

		resp := make([]PlaylistResponse, len(pls))
		for i, p := range pls {
			resp[i] = PlaylistResponse{
				ID:    p.ID,
				Title: p.Title,
				//Subtitle:    p.Subtitle,
				Cover:       p.Cover,
				LastUpdated: p.LastUpdated,
			}
		}
		c.JSON(http.StatusOK, resp)
	}
}
