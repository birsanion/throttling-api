package main

import (
	"throttling-api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, redis *redis.Client) {
	auth := middlewares.AuthorizationMiddleware(db)

	fixedWindowThrottle := middlewares.NewFixedWindowThrottle(AppConfig.ThrottleWindowDuration, redis)
	slidingLogThrottle := middlewares.NewSlidingLogThrottle(AppConfig.ThrottleWindowDuration)

	router.GET("/foo", auth, middlewares.ThrottleMiddleware(fixedWindowThrottle, "foo"), FooHandler())
	router.GET("/bar", auth, middlewares.ThrottleMiddleware(slidingLogThrottle, "bar"), BarHandler())
}
