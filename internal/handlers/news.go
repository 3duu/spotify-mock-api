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
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"` // send back as “YYYY-MM-DD”
	Type    string `json:"type"`
	ItemID  int    `json:"item_id"`
	Image   string `json:"image"` // URL to the newsletter image
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
			// pick the correct cover based on the newsletter type
			var cover string
			switch n.Type {
			case "ALBUM":
				var a models.Album
				if err := db.First(&a, n.ItemID).Error; err == nil {
					cover = a.Cover
				}
			case "PODCAST":
				var p models.Podcast
				if err := db.First(&p, n.ItemID).Error; err == nil {
					cover = p.Cover
				}
			case "ARTIST":
				var ar models.Artist
				if err := db.First(&ar, n.ItemID).Error; err == nil {
					// assuming Artist has an Image or Cover field
					cover = ar.Image
				}
			default:
				// fallback placeholder
				cover = "/media/default-newsletter.jpg"
			}

			resp[i] = NewsletterResponse{
				ID:      n.ID,
				Title:   n.Title,
				Content: n.Content,
				Date:    n.Date,
				Type:    n.Type,
				ItemID:  n.ItemID,
				Image:   cover,
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
