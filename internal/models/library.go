package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Playlist struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Cover       string    `json:"cover"`
	LastUpdated time.Time `gorm:"autoUpdateTime" json:"last_updated"`
	UserID      int       `json:"user_id"`
	Songs       []Song    `gorm:"many2many:playlist_songs;" json:"songs"`
}
type Artist struct {
	ArtistId int    `json:"artist_id" gorm:"primaryKey"`
	Name     string `json:"name"`
}
type Album struct {
	AlbumId  int    `json:"album_id" gorm:"primaryKey"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`           // foreign key column
	Artist   Artist `gorm:"foreignKey:ArtistID"` // association
	Cover    string `json:"cover"`
}
type Podcast struct {
	ID    int            `json:"id" gorm:"primaryKey"`
	Title string         `json:"title"`
	Hosts datatypes.JSON `json:"hosts" gorm:"type:json"`
	Cover string         `json:"cover"`
}

type Song struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   Artist `json:"artist"`
	Album    Album  `json:"album"`
	Genre    string `json:"genre"`
	Duration int    `json:"duration"`
}

type PlaylistSong struct {
	PlaylistID string
	SongID     string
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
