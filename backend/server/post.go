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
func GetPostList(page, size int64) (datas []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed: ", zap.Error(err))
	}

	datas = make([]*models.ApiPostDetail, 0, len(postList))

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
		datas = append(datas, postDetail)
	}
	return
}

// 获取帖子列表并按照参数排序
func GetPostList2(p *models.ParamPostList) (datas []*models.ApiPostDetail, err error) {
	// 在redis中查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetPostIDsInOrder(p) failed: ", zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}

	// 根据id在Mysql数据库中查询帖子详细信息
	postList, err := mysql.GetPostDetailListByIds(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailListByIds(ids) failed: ", zap.Error(err))
		return
	}

	datas = make([]*models.ApiPostDetail, 0, len(postList))

	// 获取赞成票和反对票票数
	var agreeVotes, disagreeVotes []int64
	agreeVotes, disagreeVotes, err = redis.GetVotes(ids)
	if err != nil {
		zap.L().Error("redis.GetVotes(ids) failed: ", zap.Error(err))
		return
	}

	for index, post := range postList {

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

		post.AgreeVotes = agreeVotes[index]
		post.DisagreeVotes = disagreeVotes[index]

		// 将信息写入api结构体
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			PostDetail:      post,
			CommunityDetail: community,
		}
		// 添加到list中
		datas = append(datas, postDetail)
	}

	return
}
