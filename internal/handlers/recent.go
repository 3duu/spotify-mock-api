package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"spotify-mock-api/internal/models"
	"time"
)

type RecentItemResponse struct {
	Type     string `json:"type"`
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"`
	Cover    string `json:"cover,omitempty"`
	IconURL  string `json:"icon_url,omitempty"`
	PlayedAt string `json:"played_at"` // ISO timestamp
}

func GetRecentPlays(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// for now, static user=1
		var recs []models.RecentPlay
		if err := db.
			Where("user_id = ?", 1).
			Order("played_at DESC").
			Limit(20).
			Find(&recs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load recents"})
			return
		}

		out := make([]RecentItemResponse, 0, len(recs))
		for _, r := range recs {
			item := RecentItemResponse{
				Type:     r.Type,
				ID:       r.ReferenceID,
				PlayedAt: r.PlayedAt.Format(time.RFC3339),
			}

			// fetch metadata based on type:
			switch r.Type {
			case "track":
				var s models.Song
				if err := db.Preload("Artist").First(&s, "id = ?", r.ReferenceID).Error; err == nil {
					item.Title = s.Title
					item.Subtitle = s.Artist.Name
					item.Cover = s.AudioURL // or album art
				}
			case "artist":
				var a models.Artist
				if err := db.First(&a, "artist_id = ?", r.ReferenceID).Error; err == nil {
					item.Title = a.Name
				}
			case "album":
				var a models.Album
				if err := db.First(&a, "album_id = ?", r.ReferenceID).Error; err == nil {
					item.Title = a.Title
					item.Cover = a.Cover
				}
			case "playlist":
				var p models.Playlist
				if err := db.First(&p, "id = ?", r.ReferenceID).Error; err == nil {
					item.Title = p.Title
					item.Cover = p.Cover
					item.Subtitle = fmt.Sprintf("Updated %s", p.LastUpdated.Format("2006-01-02"))
				}
			case "podcast":
				var p models.Podcast
				if err := db.First(&p, "id = ?", r.ReferenceID).Error; err == nil {
					item.Title = p.Title
					item.Cover = p.Cover
				}
			}
			out = append(out, item)
		}

		c.JSON(http.StatusOK, out)
	}
}
