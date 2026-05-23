package router

import (
	"net/http"
	"web/controller"
	"web/logger"
	"web/middlerware"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 将gin设置成发布模式。控制台将不会打印路由信息
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	// 注册路由
	v1.POST("/signup", controller.SingUpHandler)
	v1.POST("/login", controller.LoginHandler)

	v1.Use(middlerware.JWTAuthMiddleware()) // 应用jwt认证中间件

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		v1.GET("/postList", controller.GetPostListHandler)
		v1.GET("/postList2", controller.GetPostListHandler2)
		v1.GET("/communityPostList", controller.GetCommunityPostListHandler)

		// 投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	v1.GET("/ping", middlerware.JWTAuthMiddleware(), func(ctx *gin.Context) {
		// 如果是登录用户，判断请求头中是否有有效的 jwt token
		ctx.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
