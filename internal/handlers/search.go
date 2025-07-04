package handlers

import (
	"net/http"
	models "spotify-mock-api/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArtistResponse matches the front‐end Artist type
type ArtistResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

// AlbumResponse matches the front‐end Album type
type AlbumResponse struct {
	AlbumID int    `json:"album_id"`
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Cover   string `json:"cover"`
}

// GetSearch handles GET /search?q=foo
func GetSearch(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))
		if q == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "q query param required"})
			return
		}
		wildcard := "%" + q + "%"

		// 1) Search Songs
		var songs []models.Song
		db.Preload("Artist"). // so s.Artist is populated in your response
					Joins("LEFT JOIN artists ON artists.artist_id = songs.artist_id").
					Where("songs.title LIKE ? OR artists.name LIKE ?", wildcard, wildcard).
					Find(&songs)
		tracks := make([]models.TrackResponse, len(songs))
		for i, s := range songs {
			tracks[i] = models.TrackResponse{
				ID:       s.ID,
				Title:    s.Title,
				Artist:   s.Artist.Name,
				AudioURL: "/media/song.mp3",
				AlbumArt: "/media/album-art.jpg",
			}
		}

		// 2) Search Artists
		var artists []models.Artist
		db.Where("name LIKE ?", wildcard).Find(&artists)
		artistRes := make([]ArtistResponse, len(artists))
		for i, a := range artists {
			artistRes[i] = ArtistResponse{
				ID:    a.ArtistId,
				Name:  a.Name,
				Image: "/media/album-art.jpg",
			}
		}

		// 3) Search Albums
		var albums []models.Album
		db.Where("title LIKE ?", wildcard).Find(&albums)
		albumRes := make([]AlbumResponse, len(albums))
		for i, al := range albums {
			albumRes[i] = AlbumResponse{
				AlbumID: al.AlbumId,
				Title:   al.Title,
				Artist:  al.Artist.Name,
				Cover:   al.Cover,
			}
		}

		// 4) Search Playlists
		var pls []models.Playlist
		db.Where("title LIKE ?", wildcard).Find(&pls)
		playRes := make([]PlaylistResponse, len(pls))
		for i, p := range pls {
			playRes[i] = PlaylistResponse{
				ID:    p.ID,
				Title: p.Title,
				Cover: p.Cover,
			}
		}

		// 5) Return combined JSON
		c.JSON(http.StatusOK, gin.H{
			"tracks":    tracks,
			"artists":   artistRes,
			"albums":    albumRes,
			"playlists": playRes,
		})
	}
}
