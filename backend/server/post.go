package server

import (
	"web/dao/mysql"
	"web/models"
	"web/pkg/snowflake"

	"go.uber.org/zap"
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

func GetPostDetailById(id uint64) (data *models.ApiPostDetail, err error) {
	zap.L().Error("server.GetPostDetailById(id) ", zap.Any("id", id))

	// 查询并组合接口需要的数据
	post, err := mysql.GetPostDetailById(id)
	zap.L().Error("mysql.GetPostDetailById(id) ", zap.Any("post", post))
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(id) failed: ", zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetInfoByUserId(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetInfoByUserId(post.AuthorID) failed: ", zap.Error(err))
		return
	}
	// 根据社区id查询社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed: ", zap.Error(err))
		return
	}
	// 将信息写入api结构体
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		PostDetail:      post,
		CommunityDetail: community,
	}
	return
}
