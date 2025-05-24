package handlers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// TrackRec is what we return for each recommended track
type TrackRec struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	Genres   []string `json:"genres"`
	AudioURL string   `json:"audio_url"`
	Cover    string   `json:"cover"` // you might use album cover here
}

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
			var fallback []TrackRec
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

		var recs []TrackRec
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
