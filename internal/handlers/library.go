package handlers

import (
	"net/http"
	"spotify-mock-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLibraryData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []models.Playlist
		var albums []models.Album
		var podcasts []models.Podcast

		db.Find(&playlists)
		db.Find(&albums)
		db.Find(&podcasts)

		// Default cover image for any album missing one
		/*const defaultCover = "/media/album-art.jpg"
		for i := range albums {
			if albums[i].Cover == "" {
				albums[i].Cover = defaultCover
			}
		}*/

		c.JSON(http.StatusOK, models.LibraryData{
			Playlists: playlists,
			Albums:    albums,
			Podcasts:  podcasts,
		})
	}
}
