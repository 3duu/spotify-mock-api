package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"spotify-mock-api/internal/models"
	"strings"
)

func GetRecommendations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		const userID = 1 // TODO: replace with real auth

		// helper to actually load TrackRec by a WHERE clause
		loadTracks := func(whereSQL string, args ...interface{}) ([]models.TrackResponse, error) {
			//conds := strings.Repeat("?", 1) // not used, but placeholder
			// we'll build our raw query dynamically below...
			var recs []models.TrackResponse
			query := `
				SELECT s.id,
				       s.title,
				       a.name   AS artist,
				       s.genres AS genres,
				       '/media/song.mp3' as audio_url, -- hardcoded for now
				       al.cover AS album_art,
				       al.title AS album,
				       s.duration AS duration,
				       a.artist_id,
				       al.album_id
				FROM songs s
				JOIN artists a ON a.artist_id = s.artist_id
				JOIN albums al ON al.album_id   = s.album_id
				WHERE ` + whereSQL + `
				ORDER BY RANDOM()
				LIMIT 20
			`
			if err := db.Raw(query, args...).Scan(&recs).Error; err != nil {
				return nil, err
			}

			return recs, nil
		}

		// 1) Try recent plays
		var recPlays []string
		if err := db.
			Raw(`SELECT reference_id 
			      FROM recent_plays 
			      WHERE user_id = ? 
			      ORDER BY played_at DESC 
			      LIMIT 20`, userID).
			Scan(&recPlays).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load recents"})
			return
		}

		var recs []models.TrackResponse
		var err error

		if len(recPlays) > 0 {
			// build WHERE id IN (?, ?, …)
			placeholders := strings.Repeat("?,", len(recPlays))
			placeholders = strings.TrimRight(placeholders, ",")
			where := "s.id IN (" + placeholders + ")"
			args := make([]interface{}, len(recPlays))
			for i, id := range recPlays {
				args[i] = id
			}
			recs, err = loadTracks(where, args...)
		}

		// 2) Fallback → user playlists
		if err == nil && len(recs) == 0 {
			// gather song_ids from ANY of this user’s playlists
			var songIDs []int
			if err = db.
				Raw(`
				  SELECT ps.song_id 
				  FROM playlist_songs ps
				  JOIN playlists p ON p.id = ps.playlist_id
				  WHERE p.user_id = ?
				`, userID).
				Scan(&songIDs).Error; err == nil && len(songIDs) > 0 {
				// convert to []interface{} for query
				args := make([]interface{}, len(songIDs))
				for i, id := range songIDs {
					args[i] = id
				}
				placeholders := strings.Repeat("?,", len(args))
				placeholders = strings.TrimRight(placeholders, ",")
				where := "s.id IN (" + placeholders + ")"
				recs, err = loadTracks(where, args...)
			}
		}

		// 3) Random fallback
		if err == nil && len(recs) == 0 {
			recs, err = loadTracks("1=1") // no WHERE, pure random
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not load recommendations"})
			return
		}

		c.JSON(http.StatusOK, recs)
	}
}
