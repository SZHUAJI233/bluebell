package mysql

import (
	"database/sql"
	"web/models"

	"go.uber.org/zap"
)

func GetCommunityList() ([]*models.Community, error) {
	var communityList []*models.Community
	err := db.Find(&communityList).Error
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("GetCommunityList 查询为空")
			err = nil
		}
		return nil, err
	}
	return communityList, err
}

func GetCommunityDetailByID(id uint) (*models.Community, error) {
	var community *models.Community
	err := db.Where("id = ?", id).First(&community).Error
	if err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("社区不存在")
			err = nil
		}
		return nil, err
	}
	return community, err
}
