package handlers

import (
	"log"
	"net/http"
	"spotify-mock-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLibraryData(db *gorm.DB) gin.HandlerFunc {

	log.Print("GetLibraryData")

	return func(c *gin.Context) {
		var playlists []models.Playlist
		var albums []models.Album
		var podcasts []models.Podcast

		db.Find(&playlists)
		db.Find(&albums)
		db.Find(&podcasts)

		c.JSON(http.StatusOK, models.LibraryData{
			Playlists: playlists,
			Albums:    albums,
			Podcasts:  podcasts,
		})
	}
}
