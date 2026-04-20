package mysql

import "web/models"

func CreatePost(p *models.Post) (err error) {
	return db.Create(p).Error
}

func GetPostDetailById(id uint64) (data *models.PostDetail, err error) {
	err = db.Table("posts").Where("post_id = ?", id).Select("post_id,author_id,community_id,status,title,content,created_at,updated_at").First(&data).Error
	return
}

func GetPostList() (data []*models.PostDetail, err error) {
	err = db.Table("posts").Select("post_id,author_id,community_id,status,title,content,created_at,updated_at").Find(&data).Error
	return
}
