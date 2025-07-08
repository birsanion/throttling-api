package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBHost     string `envconfig:"db_host" required:"true"`
	DBUser     string `envconfig:"db_user" required:"true"`
	DBPassword string `envconfig:"db_password" required:"true"`
	DBName     string `envconfig:"db_name" required:"true"`

	RedisAddr string `split_words:"true" validate:"required"`
	RedisDB   int    `split_words:"true" validate:"gte=0,lte=16"`

	ThrottleWindowDuration time.Duration `envconfig:"throttle_window_duration" default:"1m"`
}

var AppConfig Config

func LoadConfig() error {
	if err := envconfig.Process("", &AppConfig); err != nil {
		return err
	}

	return nil
}
