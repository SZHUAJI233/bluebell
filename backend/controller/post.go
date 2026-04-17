package controller

import (
	"web/controller/request"
	"web/controller/response"
	"web/models"
	"web/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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
