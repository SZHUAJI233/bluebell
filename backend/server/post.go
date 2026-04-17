package server

import (
	"web/dao/mysql"
	"web/models"
	"web/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成postID
	p.PostID = snowflake.GenID()
	// 2. 保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		return err
	}
	// 3. 返回
	return
}
