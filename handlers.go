package main

import (
	"net/http"
	"throttling-api/models/responses"

	"github.com/gin-gonic/gin"
)

func FooHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.NewSuccessResponse(true))
	}
}

func BarHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, responses.NewSuccessResponse(true))
	}
}
