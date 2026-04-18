package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID          uint64 `gorm:"type:bigint(20);primaryKey" json:"id"`
	PostID      uint64 `gorm:"type:bigint(20);not null;uniqueIndex:idx_post_id" json:"post_id" `
	AuthorID    uint64 `gorm:"type:bigint(20);not null;index:idx_author_id" json:"author_id"`
	CommunityID uint64 `gorm:"type:bigint(20);not null;index:idx_community_id" json:"community_id"`
	Status      uint8  `gorm:"type:tinyint(4);not null;default:1;"`
	Title       string `gorm:"type:varchar(128);collate:utf8mb4_general_ci;not null" json:"title" binding:"required"`
	Content     string `gorm:"type:varchar(8192);collate:utf8mb4_general_ci;not null" json:"content" binding:"required"`

	CreatedAt time.Time      `json:"create_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

type PostDetail struct {
	PostID      uint64 `gorm:"type:bigint(20);not null;uniqueIndex:idx_post_id" json:"post_id" `
	AuthorID    uint64 `gorm:"type:bigint(20);not null;index:idx_author_id" json:"author_id"`
	CommunityID uint64 `gorm:"type:bigint(20);not null;index:idx_community_id" json:"community_id"`
	Status      uint8  `gorm:"type:tinyint(4);not null;default:1;"`
	Title       string `gorm:"type:varchar(128);collate:utf8mb4_general_ci;not null" json:"title" binding:"required"`
	Content     string `gorm:"type:varchar(8192);collate:utf8mb4_general_ci;not null" json:"content" binding:"required"`

	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// 帖子详情接口
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	*PostDetail                         // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
}
