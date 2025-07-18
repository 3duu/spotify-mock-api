package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"spotify-mock-api/internal/models"
	"strconv"
	"time"
)

type PlaylistResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Cover       string    `json:"cover"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
}

// PlaylistDetailResponse is the full payload for GET /playlists/:id
type PlaylistDetailResponse struct {
	ID         int                    `json:"id"`
	Title      string                 `json:"title"`
	Cover      string                 `json:"cover"`
	OwnerName  string                 `json:"ownerName"`
	OwnerImage string                 `json:"ownerImage"`
	Duration   string                 `json:"duration"` // e.g. "5h 59m"
	Tracks     []models.TrackResponse `json:"tracks"`
}

// GetRecentPlaylistsByUser returns up to 10 most‐recently updated playlists
// for the given user ID, mapping them into PlaylistResponse.
func GetRecentPlayedByUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDParam := c.Param("userId")
		if userIDParam == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}

		// Parse userId to int
		userID, err := strconv.Atoi(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userId must be an integer"})
			return
		}

		// Fetch the 10 most recently updated playlists for this user
		var pls []models.Playlist
		if err := db.
			Where("user_id = ?", userID).
			Order("last_updated DESC").
			Limit(10).
			Find(&pls).
			Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load playlists"})
			return
		}

		// Map to API response type, counting tracks in the join table
		resp := make([]PlaylistResponse, len(pls))
		for i, p := range pls {
			var trackCount int64
			if err := db.
				Model(&models.PlaylistSong{}).
				Where("playlist_id = ?", p.ID).
				Count(&trackCount).
				Error; err != nil {
				// if counting fails, just fall back to zero
				trackCount = 0
			}

			subtitle := fmt.Sprintf("Playlist • %d tracks", trackCount)
			resp[i] = PlaylistResponse{
				ID:       p.ID,
				Title:    p.Title,
				Subtitle: subtitle,
				Cover:    p.Cover,
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}

// GetPlaylistDetail loads one playlist and returns its full detail
func GetPlaylistDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistID := c.Param("id")

		// 1) Load playlist, its owner, and its tracks (+ artists + cover fields)
		var pl models.Playlist
		if err := db.
			Preload("Owner").
			Preload("Songs", func(db *gorm.DB) *gorm.DB {
				return db.
					Preload("Artist").
					Preload("Album")
			}).
			First(&pl, playlistID).
			Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "playlist not found"})
			return
		}

		// 2) Build the slice of TrackResponse
		tracks := make([]models.TrackResponse, len(pl.Songs))
		var totalSec int
		for i, t := range pl.Songs {
			// sum up durations
			totalSec += t.Duration

			tracks[i] = models.TrackResponse{
				ID:         t.ID,
				Title:      t.Title,
				Artist:     t.Artist.Name,
				AlbumArt:   t.Album.Cover,
				Downloaded: false,
				Duration:   t.Duration,
				Album:      t.Album.Title,
			}
		}

		h := totalSec / 3600
		m := (totalSec % 3600) / 60
		durationStr := fmt.Sprintf("%dh %02dm", h, m)

		// 4) Assemble the response
		resp := PlaylistDetailResponse{
			ID:         pl.ID,
			Title:      pl.Title,
			Cover:      pl.Cover,
			OwnerName:  pl.Owner.Name,
			OwnerImage: pl.Owner.Image,
			Duration:   durationStr,
			Tracks:     tracks,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func AddTrackToPlaylist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		plID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
			return
		}

		var body struct {
			TrackID int `json:"track_id"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing track_id"})
			return
		}

		// create join record
		entry := models.PlaylistSong{
			PlaylistID: plID,
			SongID:     body.TrackID,
		}
		// ON CONFLICT DO NOTHING → no more UNIQUE violations
		if err := db.
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&entry).
			Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not add track"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// DELETE /playlists/:id/tracks/:trackId
func RemoveTrackFromPlaylist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		plID, err1 := strconv.Atoi(c.Param("id"))
		trID, err2 := strconv.Atoi(c.Param("trackId"))
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid IDs"})
			return
		}

		if err := db.
			Where("playlist_id = ? AND song_id = ?", plID, trID).
			Delete(&models.PlaylistSong{}).
			Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not remove track"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

type createPlaylistRequest struct {
	Title string `json:"title" binding:"required"`
	Cover string `json:"cover"`
}

// POST /playlists
// Body: { "title": "My New Playlist", "cover": "/media/my-cover.jpg" }
func CreatePlaylist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1) bind JSON
		var body createPlaylistRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
			return
		}

		userID := 1 // TODO: replace with actual authenticated user ID

		// 2) check for existing playlist with same title
		var existing models.Playlist
		err := db.
			Where("user_id = ? AND title = ?", userID, body.Title).
			First(&existing).
			Error

		if err == nil {
			// found one → reject
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("you already have a playlist named %q", body.Title),
			})
			return
		} else if err != gorm.ErrRecordNotFound {
			// some other DB error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify playlist name"})
			return
		}

		// 3) create new playlist
		pl := models.Playlist{
			Title:  body.Title,
			Cover:  body.Cover,
			UserID: userID,
		}
		if err := db.Create(&pl).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create playlist"})
			return
		}

		// 4) return the new playlist
		resp := PlaylistResponse{
			ID:          pl.ID,
			Title:       pl.Title,
			Subtitle:    fmt.Sprintf("Playlist • 0 tracks"),
			Cover:       pl.Cover,
			LastUpdated: pl.LastUpdated,
		}
		c.JSON(http.StatusCreated, resp)
	}
}

// PUT    /playlists/:id
// Body: { "title": "New Name", "cover": "/media/new.jpg" }
func UpdatePlaylistMeta(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		plID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
			return
		}

		var body struct {
			Title string `json:"title"`
			Cover string `json:"cover"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		if err := db.Model(&models.Playlist{}).
			Where("id = ?", plID).
			Updates(models.Playlist{Title: body.Title, Cover: body.Cover}).
			Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update playlist"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// PUT    /playlists/:id/reorder
// Body: { "track_ids": [3,5,2,1] }
func ReorderPlaylist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		plID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid playlist ID"})
			return
		}
		var body struct {
			TrackIDs []int `json:"track_ids"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		// update join table ordering: assumes PlaylistSong has Position field
		for pos, tid := range body.TrackIDs {
			if err := db.Model(&models.PlaylistSong{}).
				Where("playlist_id = ? AND song_id = ?", plID, tid).
				Update("position", pos).
				Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not reorder"})
				return
			}
		}
		c.Status(http.StatusNoContent)
	}
}
