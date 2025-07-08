package main

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func failOnError(err error, msg string) {
	if err != nil {
		logrus.Panicf("%s: %s", msg, err)
	}
}

func main() {
	logrus.Info("Application is starting...")

	err := LoadConfig()
	failOnError(err, "Failed to create config")

	db, err := setupDB(AppConfig)
	failOnError(err, "Failed to create db connection")

	redis, err := setupRedis(AppConfig)
	failOnError(err, "Failed to create redis connection")

	router := gin.Default()
	RegisterRoutes(router, db, redis)
	router.Run(":8888")
}

func setupDB(cfg Config) (*gorm.DB, error) {
	return NewDbConnection(cfg)
}

func setupRedis(cfg Config) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	}), nil
}
