package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"spotify-mock-api/internal/models"
)

func GetCurrentUser(c *gin.Context) {

	// pick scheme based on TLS
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	// build a full URL – or just return "/media/user_image.jpg"
	avatarURL := fmt.Sprintf("%s://%s/media/user_image.jpg", scheme, c.Request.Host)

	user := models.User{
		ID:    "1",
		Name:  "Eduardo Porto de Araujo",
		Image: avatarURL,
	}
	c.JSON(http.StatusOK, user)
}
