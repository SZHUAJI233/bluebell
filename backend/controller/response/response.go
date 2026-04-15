package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code":1001, // 程序中的错误码
	"msg":xx,	// 提示信息
	"data":{}	// 数据
}
*/

type Data struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} ` json:"data"`
}

func Error(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &Data{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ErrorWhitMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &Data{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Data{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

func SuccessWhitMsg(c *gin.Context, msg interface{}, data interface{}) {
	c.JSON(http.StatusOK, &Data{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}
