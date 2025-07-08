package db

import (
	"time"
)

type User struct {
	ID        int64  `gorm:"primaryKey"`
	ClientID  string `gorm:"index:unique"`
	RateLimit int
	CreatedAt time.Time
	UpdatedAt time.Time
}
