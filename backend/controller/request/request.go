package request

import (
	"errors"
	"web/middlerware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUser(c *gin.Context) (userID uint64, err error) {
	uid, ok := c.Get(middlerware.ContextUserIDKey)
	if !ok {
		zap.L().Error("c.Get(middlerware.ContextUserIDKey) failed")
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(uint64)
	if !ok {
		zap.L().Error("uid.(int64) failed")
		err = ErrorUserNotLogin
		return
	}
	return
}
