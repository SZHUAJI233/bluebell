package server

import (
	"strconv"
	"web/dao/mysql"
	"web/dao/redis"
	"web/models"
	"web/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1. 生成postID
	p.PostID = snowflake.GenID()
	// 2. 保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("mysql.CreatePost(p) failed: ", zap.Error(err))
		return err
	}
	// redis 记录帖子创建时间
	if err = redis.CreatePost(strconv.Itoa(int(p.PostID))); err != nil {
		zap.L().Error("redis.CreatePost(strconv.Itoa(int(p.PostID))) failed: ", zap.Error(err))
		return err
	}
	// 3. 返回
	return
}

func GetPostDetailById(id uint64) (data *models.ApiPostDetail, err error) {
	// 查询并组合接口需要的数据
	post, err := mysql.GetPostDetailById(id)
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

// 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed: ", zap.Error(err))
	}
	data = make([]*models.ApiPostDetail, 0, len(postList))

	for _, post := range postList {
		// 根据作者id查询作者信息
		var user *models.User
		user, err = mysql.GetInfoByUserId(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetInfoByUserId(post.AuthorID) failed: ", zap.Error(err))
			return
		}
		// 根据社区id查询社区信息
		var community *models.CommunityDetail
		community, err = mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed: ", zap.Error(err))
			return
		}
		// 将信息写入api结构体
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			PostDetail:      post,
			CommunityDetail: community,
		}
		// 添加到list中
		data = append(data, postDetail)
	}
	return
}
