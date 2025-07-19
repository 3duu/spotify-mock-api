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
		var recs []models.RecentPlay
		if err := db.
			Where("user_id = ?", 1).
			Order("played_at DESC").
			Limit(50). // Increase limit if you want 20 unique recents
			Find(&recs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load recents"})
			return
		}

		out := make([]RecentItemResponse, 0, len(recs))
		seen := make(map[string]bool)

		for _, r := range recs {
			// Key by type and ReferenceID to ensure uniqueness
			key := fmt.Sprintf("%s-%d", r.Type, r.ReferenceID)
			if seen[key] {
				continue // skip duplicates
			}
			seen[key] = true

			item := RecentItemResponse{
				Type:     r.Type,
				ID:       r.ReferenceID,
				PlayedAt: r.PlayedAt.Format(time.RFC3339),
			}

			// fetch metadata as before
			switch r.Type {
			case "track":
				var s models.Song
				if err := db.Preload("Artist").First(&s, "id = ?", r.ReferenceID).Error; err == nil {
					item.Title = s.Title
					item.Subtitle = s.Artist.Name
					item.Cover = "/media/album-art.jpg"
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
					item.Cover = "/media/album-art.jpg"
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

			// If you want exactly 8 items, break after collecting 8 uniques
			if len(out) >= 8 {
				break
			}
		}

		c.JSON(http.StatusOK, out)
	}
}

func GetRecentTracks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var recs []models.RecentPlay
		if err := db.
			Where("user_id = ?", 1).
			Order("played_at DESC").
			Limit(8). // Increase limit if you want 20 unique recents
			Find(&recs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load recents"})
			return
		}

		out := make([]models.TrackResponse, 0, len(recs))
		seen := make(map[string]bool)

		for _, r := range recs {
			// Key by type and ReferenceID to ensure uniqueness
			key := fmt.Sprintf("%s-%d", r.Type, r.ReferenceID)
			if seen[key] {
				continue // skip duplicates
			}
			seen[key] = true

			item := RecentItemResponse{
				Type:     r.Type,
				ID:       r.ReferenceID,
				PlayedAt: r.PlayedAt.Format(time.RFC3339),
			}

			// fetch metadata as before
			switch r.Type {
			case "track":
				var s models.Song
				if err := db.Preload("Artist").First(&s, "id = ?", r.ReferenceID).Error; err == nil {
					item.Title = s.Title
					item.Subtitle = s.Artist.Name
					item.Cover = "/media/album-art.jpg"
				}
			}

			out = append(out, models.TrackResponse{
				ID:         r.ReferenceID,
				Title:      item.Title,
				Artist:     item.Subtitle,
				AlbumArt:   item.Cover,
				Duration:   169,
				Downloaded: false,
			})

			// If you want exactly 8 items, break after collecting 8 uniques
			if len(out) >= 4 {
				break
			}
		}

		c.JSON(http.StatusOK, out)
	}
}
