package models

import "gorm.io/gorm"

type Playlist struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Cover    string
	Subtitle string `json:"subtitle"`
	IconURL  string `json:"icon"`
}
type Album struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Cover  string `json:"cover"`
}
type Podcast struct {
	ID    string   `json:"id"`
	Title string   `json:"title"`
	Hosts []string `json:"hosts"`
	Art   string   `json:"art"`
}

// Bundle into a single response object
type LibraryData struct {
	Playlists []Playlist `json:"playlists"`
	Albums    []Album    `json:"albums"`
	Podcasts  []Podcast  `json:"podcasts"`
}

type LibraryEntry struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Type        string // "playlist" | "album" | "podcast"
	ReferenceID string // points to the real Playlist.ID, Album.ID, etc.
	Title       string
	Subtitle    string
	IconURL     string
}
