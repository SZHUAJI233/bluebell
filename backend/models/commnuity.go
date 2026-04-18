package models

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	ID            uint64 `gorm:"type:bigint(20);primaryKey" json:"id"`
	CommunityID   uint64 `gorm:"type:bigint(20);not null;uniqueIndex:idx_community_id" json:"community_id"`
	CommunityName string `gorm:"type:varchar(128);not null;uniqueIndex:idx_community_name" json:"community_name"`
	Introduction  string `gorm:"type:varchar(256);not null" json:"introduction"`

	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

type CommunityDetail struct {
	CommunityID   uint64 `gorm:"type:bigint(20);not null;uniqueIndex:idx_community_id" json:"community_id"`
	CommunityName string `gorm:"type:varchar(128);not null;uniqueIndex:idx_community_name" json:"community_name"`
	Introduction  string `gorm:"type:varchar(256);not null" json:"introduction"`
}
