package server

import (
	"web/dao/mysql"
	"web/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查找数据，查找到所有的community 并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetailByID(id uint) (*models.CommunityDetail, error) {
	// 查找数据，查找到所有的community 并返回
	return mysql.GetCommunityDetailByID(id)
}
