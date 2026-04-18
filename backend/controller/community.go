package controller

import (
	"strconv"
	"web/controller/response"
	"web/server"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区 (community_id, community_name) 以列表形式返回
	data, err := server.GetCommunityList()
	if err != nil {
		zap.L().Error("server.GetCommunityList() failed", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	response.Success(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	// 获取社区ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 0)
	uid := uint64(id)
	data, err := server.GetCommunityDetailByID(uid)
	if err != nil {
		zap.L().Error("server.GetCommunityList() failed", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
	}
	response.Success(c, data)
}
