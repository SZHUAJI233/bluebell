package mysql

import "web/models"

func CreatePost(p *models.Post) (err error) {
	return db.Create(p).Error
}
