package models

// Song represents the metadata for a track
type Song struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration int    `json:"duration"` // duration in seconds
}
