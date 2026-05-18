package controller

import (
	"strconv"
	"web/controller/request"
	"web/controller/response"
	"web/models"
	"web/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1. 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("c.ShouldBindJSON(p) failed：", zap.Error(err))
		response.Error(c, response.CodeInvalidParm)
		return
	}
	// 从c中取出当前请求的用户id
	userID, err := request.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("request.GetCurrentUser(c)：", zap.Error(err))
		response.Error(c, response.CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2. 创建帖子
	if err := server.CreatePost(p); err != nil {
		zap.L().Error("server.CreatePost(p) failed：", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.Success(c, nil)
}

// 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取参数（从url中获得帖子的id）
	pidStr := c.Param("id")
	zap.L().Error("c.Param(id) ", zap.Any("pidStr", pidStr))
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	zap.L().Error("strconv.ParseInt(pidStr, 10, 64) ", zap.Any("pid", pid))

	if err != nil {
		zap.L().Error("strconv.ParseInt failed: ", zap.Error(err))
		response.Error(c, response.CodeInvalidParm)
		return
	}
	// 2. 根据id取出帖子数据
	p, err := server.GetPostDetailById(uint64(pid))
	if err != nil {
		zap.L().Error("server.GetPostDetailById failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.Success(c, p)
}

// 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size, err := request.GetPageInfo(c)
	if err != nil {
		zap.L().Error("server.GetPostInfo(c) failed: ", zap.Error(err))
	}
	// 获取数据
	data, err := server.GetPostList(page, size)
	if err != nil {
		zap.L().Error("server.GetPostList() failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	// 返回响应
	response.Success(c, data)
}

// 根据时间或分数获取帖子列表
// 根据前端传来的参数（分数、创建时间、）排序
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数： /api/v1/post2?page=1&size=10&order=time

	// 获取排序参数
	// 指定初始默认参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBindQuery(p); err != nil {

	}

	// 获取数据
	data, err := server.GetPostList2(p)
	if err != nil {
		zap.L().Error("server.GetPostList() failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}

	// 返回响应
	response.Success(c, data)
}
