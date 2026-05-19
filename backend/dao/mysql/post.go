package mysql

import (
	"strconv"
	"web/models"

	"go.uber.org/zap"
)

// 创建帖子
func CreatePost(p *models.Post) (err error) {
	return db.Create(p).Error
}

// 根据id查询帖子详细信息
func GetPostDetailById(id uint64) (data *models.PostDetail, err error) {
	err = db.Table("posts").
		Where("post_id = ?", id).
		Select("post_id,author_id,community_id,status,title,content,created_at,updated_at").
		First(&data).Error
	return
}

// 根据id列表查询帖子详细信息列表
func GetPostDetailListByIds(ids []string) (datas []*models.PostDetail, err error) {
	for _, id := range ids {

		idInt, err := strconv.Atoi(id)
		if err != nil {
			zap.L().Error("strconv.Atoi(id) failed: ", zap.Error(err))
			return nil, err
		}

		data, err := GetPostDetailById(uint64(idInt))
		if err != nil {
			zap.L().Error("GetPostDetailById(uint64(idStr)) failed: ", zap.Error(err), zap.Any("id: ", id))
			return nil, err
		}

		datas = append(datas, data)

	}
	return
}

// 获取帖子列表
func GetPostList(page, size int64) (data []*models.PostDetail, err error) {
	// 计算偏移量
	offset := (page - 1) * size
	err = db.Table("posts").
		Select("post_id,author_id,community_id,status,title,content,created_at,updated_at").
		Limit(int(size)).
		Offset(int(offset)).
		Find(&data).Error
	return
}
