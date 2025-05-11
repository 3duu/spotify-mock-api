package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Playlist struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	IconURL  string `json:"icon"`
}

// Handler for recent playlists
func GetRecentPlaylists(c *gin.Context) {
	log.Print("GetRecentPlaylists")
	// TODO: replace with real data from DB or cache
	playlists := []Playlist{
		{ID: "1", Title: "Liked songs", Subtitle: "Playlist • 27 músicas", IconURL: "/media/like.png"},
		{ID: "2", Title: "Novos episódios", Subtitle: "Atualizado em 23 de abr. de 2025", IconURL: "/media/podcast.png"},
		// … more entries …
	}
	c.JSON(http.StatusOK, playlists)
}
