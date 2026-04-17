package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"type:bigint(20);primaryKey"`
	UserID   uint64 `gorm:"type:bigint(20);not null;uniqueIndex:idx_user_id"`
	Username string `gorm:"type:varchar(64);not null;uniqueIndex:idx_username"`
	Password string `gorm:"type:varchar(64);not null"`
	Email    string `gorm:"type:varchar(64)"`
	Gender   uint8  `gorm:"type:tinyint(4);not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
