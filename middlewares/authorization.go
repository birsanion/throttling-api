package middlewares

import (
	"errors"
	"net/http"
	"strings"

	models "throttling-api/models/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthorizationMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		prefix := "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			abortWithError(c, http.StatusUnauthorized, "unauthorized")
			return
		}

		clientID := auth[len(prefix):]

		var user models.User
		if err := db.First(&user, models.User{
			ClientID: clientID,
		}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				abortWithError(c, http.StatusUnauthorized, "unauthorized")
				return
			}

			abortWithError(c, http.StatusInternalServerError, "internal error")
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
