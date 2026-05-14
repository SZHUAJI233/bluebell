package request

import (
	"errors"
	"strconv"
	"web/middlerware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurrentUserID(c *gin.Context) (userID uint64, err error) {
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

func GetPageInfo(c *gin.Context) (int64, int64, error) {
	// 获取分页参数
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(pageStr, 10, 64) failed: ", zap.Error(err))
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(sizeStr, 10, 64) failed: ", zap.Error(err))
		size = 10
	}
	return page, size, err
}
