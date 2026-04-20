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
	userID, err := request.GetCurrentUser(c)
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
	// 获取数据
	data, err := server.GetPostList()
	if err != nil {
		zap.L().Error("server.GetPostList() failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	// 返回响应
	response.Success(c, data)
}
