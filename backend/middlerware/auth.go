package middlerware

import (
	"strings"
	"web/controller/response"
	"web/pkg/jwt"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

// 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取请求头
		authHeader := c.Request.Header.Get("Authorization")
		// 请求头为空
		if authHeader == "" {
			response.Error(c, response.CodeNeedLogin)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, response.CodeInvalidAuth)
			c.Abort()
			return
		}

		// 解析token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Error(c, response.CodeInvalidAuth)
			c.Abort()
			return
		}

		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set(ContextUserIDKey, mc.UserID) // 后续可以通过c.Get获取当前请求用户的信息
		c.Next()
	}
}
