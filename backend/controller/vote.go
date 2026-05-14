package controller

import (
	"web/controller/request"
	"web/controller/response"
	"web/models"
	"web/server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 投票

type VoteData struct {
	// UserID 从请求中获取当前用户
	PostID    int64 `json:"post_id,string"`   // 帖子id
	Direction int   `json:"direction,string"` // 赞成票（1），反对票（-1）
}

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)

	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			response.Error(c, response.CodeInvalidParm)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		response.ErrorWhitMsg(c, response.CodeInvalidParm, errData)
		return
	}

	// 获取userID
	userID, err := request.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("request.GetCurrentUserID(c) failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}

	if err := server.PostVote(userID, p); err != nil {
		zap.L().Error("server.PostVote(userID, p) failed: ", zap.Error(err))
		response.Error(c, response.CodeServerBusy)
		return
	}
	response.Success(c, nil)
}
