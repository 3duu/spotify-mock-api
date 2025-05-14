package handlers

import (
	"net/http"
	models "spotify-mock-api/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCurrentUser loads the user record from the DB (ID=1) and returns it.
func GetCurrentUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, 1).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			}
			return
		}

		// 2) Build a full URL for the avatar
		/*scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		user.Image = fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, user.Image)*/

		// 3) Return the JSON
		c.JSON(http.StatusOK, user)
	}
}
