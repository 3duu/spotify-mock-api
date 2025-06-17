package models

type TrackResponse struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	ArtistID   int    `json:"artist_id"`
	AudioURL   string `json:"audio_url"`
	AlbumArt   string `json:"album_art,omitempty"`
	AlbumID    int    `json:"album_id"`
	Album      string `json:"album"`
	Downloaded bool   `json:"downloaded"`
	Duration   int    `json:"duration,omitempty"` // in seconds
	Color      string `json:"color,omitempty"`    // hex color code for UI
	Genres     string `json:"genres,omitempty"`
}
