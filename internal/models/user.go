package models

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"` // URL to the avatar
}
