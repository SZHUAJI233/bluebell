package controller

import (
	"web/controller/response"
	"web/models"
	"web/server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func SingUpHandler(c *gin.Context) {
	// 获取参数和参数校验
	var p models.ParamSignUp
	// 获取参数并自动绑定到结构体
	if err := c.ShouldBindJSON(&p); err != nil {
		// 记录日志
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是否是validationErrors类型（校验器类型）错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.Error(c, response.CodeInvalidParm)
			return
		}
		// 请求参数有误，直接返回响应
		response.ErrorWhitMsg(c, response.CodeInvalidParm, removeTopStruct(errs.Translate(trans)))
		return
	}
	// // 手动对请求参数进行详细的业务规则校验		简化至结构体，使用tag进行参数校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 {
	// 	zap.L().Error("SignUp with invalid param")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }

	// 业务处理
	if err := server.SignUp(&p); err != nil {
		zap.L().Error("server.SignUp failed", zap.Error(err))
		response.ErrorWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	// 返回响应
	response.SuccessWhitMsg(c, "注册成功", nil)
}

func LoginHandler(c *gin.Context) {
	// 1. 获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 记录日志
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是否是validationErrors类型（校验器类型）错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			response.Error(c, response.CodeInvalidParm)
			return
		}
		// 请求参数有误，直接返回响应
		response.ErrorWhitMsg(c, response.CodeInvalidParm, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务逻辑处理
	token, err := server.Login(p)
	if err != nil {
		zap.L().Error("server.Login failed", zap.Error(err))
		response.Error(c, response.CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	response.SuccessWhitMsg(c, "登录成功", token)
}
