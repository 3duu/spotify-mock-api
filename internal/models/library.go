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
	Owner       User      `gorm:"foreignKey:UserID"`
	Songs       []Song    `gorm:"many2many:playlist_songs;"`
	SongIDs     []int     `gorm:"-" json:"songs"`
}

type PlaylistSong struct {
	PlaylistID int `json:"playlist_id"`
	SongID     int `json:"song_id"`
}

type Song struct {
	ID    int    `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`

	// foreign key to Artist
	ArtistID int    `json:"artist_id"`
	Artist   Artist `gorm:"foreignKey:ArtistID"` // the actual relation

	// foreign key to Album
	AlbumID int   `json:"album_id"`
	Album   Album `gorm:"foreignKey:AlbumID"` // the actual relation

	Genres   datatypes.JSON `json:"genres"` // ← store raw JSON in SQLite
	Duration int            `json:"duration"`
	AudioURL string         `json:"audio_url" default:"/media/song.mp3"`
}

type Artist struct {
	ArtistId int    `json:"artist_id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Songs    []Song `gorm:"foreignKey:ArtistID"` // association to songs
	SongIDs  []int  `gorm:"-" json:"songs"`      // for frontend convenience
	Image    string `json:"image"`               // optional image for the artist
}
type Album struct {
	AlbumId  int    `json:"album_id" gorm:"primaryKey"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`           // foreign key column
	Artist   Artist `gorm:"foreignKey:ArtistID"` // association
	Cover    string `json:"cover"`
	Songs    []Song `gorm:"foreignKey:AlbumID"` // association to songs
	SongIDs  []int  `gorm:"-" json:"songs"`     // for
	Image    string `json:"image"`              // optional image for the album
}

type PodcastEpisode struct {
	ID          int
	Title       string `json:"title"`
	Duration    int    `json:"duration"`
	AudioURL    string `json:"audio_url"`
	Description string `json:"description"`
}

type Podcast struct {
	ID       int            `json:"id" gorm:"primaryKey"`
	Title    string         `json:"title"`
	Hosts    datatypes.JSON `json:"hosts" gorm:"type:json"`
	Cover    string         `json:"cover"`
	Episodes datatypes.JSON `json:"episodes" gorm:"type:json"`
}

// Bundle into a single response object
type LibraryData struct {
	Playlists []Playlist `json:"playlists"`
	Albums    []Album    `json:"albums"`
	Podcasts  []Podcast  `json:"podcasts"`
}

type LibraryEntry struct {
	gorm.Model
	ID          int    `gorm:"primaryKey"`
	Type        string // "playlist" | "album" | "podcast"
	ReferenceID string // points to the real Playlist.ID, Album.ID, etc.
	Title       string
	Subtitle    string
	IconURL     string
}

type RecentPlay struct {
	ID          uint      `gorm:"primaryKey" json:"-"`
	UserID      int       `json:"user_id"`
	Type        string    `json:"type"`         // "track", "artist", "album", "playlist", "podcast"
	ReferenceID int       `json:"reference_id"` // the ID of the item played
	OriginID    int       `json:"origin_id"`
	PlayedAt    time.Time `gorm:"autoCreateTime" json:"played_at"`
}

type Newsletter struct {
	gorm.Model
	ID      int    `gorm:"primaryKey" json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"date"`
	Type    string `json:"type"`
	ItemID  int    `json:"item_id"`
}
