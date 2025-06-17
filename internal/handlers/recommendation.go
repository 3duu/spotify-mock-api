package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"spotify-mock-api/internal/models"
	"strings"
)

// GetRecommendations returns up to 20 random tracks matching the user's top genres.
func GetRecommendations(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, we hard‐code user_id = 1
		const userID = 1

		// 1) Find the top 3 genres the user has in their library_entries → songs.genres
		// Because genres is stored as JSON text, we unnest it via raw SQL.
		rows, err := db.Raw(`
			SELECT genre, COUNT(*) as cnt
			FROM library_entries le
			JOIN songs s ON s.id = le.song_id
			CROSS JOIN json_each(s.genres)
			WHERE le.user_id = ?
			GROUP BY genre
			ORDER BY cnt DESC
			LIMIT 3
		`, userID).Rows()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot compute favorite genres"})
			return
		}
		defer rows.Close()

		var topGenres []string
		for rows.Next() {
			var genre sql.NullString
			var cnt int
			if err := rows.Scan(&genre, &cnt); err == nil && genre.Valid {
				topGenres = append(topGenres, genre.String)
			}
		}
		if len(topGenres) == 0 {
			// fallback: if no library entries, just return some random tracks
			var fallback []models.TrackResponse
			db.Raw(`
				SELECT s.id, s.title, a.name as artist, s.genres, s.audio_url, al.cover
				FROM songs s
				JOIN artists a ON a.artist_id = s.artist_id
				JOIN albums al ON al.album_id   = s.album_id
				ORDER BY RANDOM()
				LIMIT 20
			`).Scan(&fallback)
			c.JSON(http.StatusOK, fallback)
			return
		}

		// 2) Query up to 20 random tracks whose genres JSON array contains any of topGenres
		// We build a WHERE clause like: genres LIKE '%"Rock"%' OR genres LIKE '%"Pop"%'
		conds := make([]string, len(topGenres))
		args := make([]interface{}, len(topGenres))
		for i, g := range topGenres {
			conds[i] = "s.genres LIKE ?"
			// match the JSON-encoded genre element: e.g. '%"Rock"%'
			args[i] = "%\"" + g + "\"%"
		}
		rawWhere := strings.Join(conds, " OR ")

		var recs []models.TrackResponse
		db.Raw(`
			SELECT s.id,
			       s.title,
			       a.name   AS artist,
			       s.genres AS genres,
			       s.audio_url,
			       al.cover AS cover
			FROM songs s
			JOIN artists a ON a.artist_id = s.artist_id
			JOIN albums  al ON al.album_id   = s.album_id
			WHERE `+rawWhere+`
			ORDER BY RANDOM()
			LIMIT 20
		`, args...).Scan(&recs)

		c.JSON(http.StatusOK, recs)
	}
}

func GetRecommendations2(db *gorm.DB) gin.HandlerFunc {
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

			/*for _, r := range recs {
				defaults.SetDefaults(r)
			}*/

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
