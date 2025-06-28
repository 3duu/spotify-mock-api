package handlers

import (
	"encoding/json"
	"net/http"
	"spotify-mock-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func podcastEpisodesToTrackResponses(episodes []models.PodcastEpisode, podcast models.Podcast) []models.TrackResponse {
	responses := make([]models.TrackResponse, len(episodes))
	for i, ep := range episodes {
		responses[i] = models.TrackResponse{
			ID:       ep.ID,
			Title:    ep.Title,
			Duration: ep.Duration,
			AudioURL: ep.AudioURL,
			Artist:   "",
			Album:    "",
			AlbumArt: podcast.Cover,
		}
	}
	return responses
}

func GetPodcastDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		rawID := c.Param("id")
		podcastID, err := strconv.Atoi(rawID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid podcast ID"})
			return
		}

		var podcast models.Podcast
		// Here, fixed to podcast with ID=1. To generalize, get id from params
		if err := db.First(&podcast, podcastID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Podcast not found"})
			return
		}

		// Episodes stored as JSON in podcast.Episodes
		var episodes []models.PodcastEpisode
		if len(podcast.Episodes) > 0 {
			if err := json.Unmarshal(podcast.Episodes, &episodes); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode episodes"})
				return
			}
		}

		resp := PlaylistDetailResponse{
			ID:    podcast.ID,
			Title: podcast.Title,
			Cover: podcast.Cover,
			//OwnerName:  "",
			//OwnerImage: artist.Image,
			//Duration:   podcast.,
			Tracks: podcastEpisodesToTrackResponses(episodes, podcast),
		}

		c.JSON(http.StatusOK, resp)
	}
}
