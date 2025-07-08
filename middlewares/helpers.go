package middlewares

import (
	"throttling-api/models/responses"

	"github.com/gin-gonic/gin"
)

func abortWithError(c *gin.Context, status int, msg string) {
	c.JSON(status, responses.NewErrorResponse(msg))
	c.Abort()
}
