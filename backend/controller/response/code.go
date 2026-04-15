package response

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParm
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidAuth
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParm:     "参数请求错误",
	CodeUserExist:       "用户名已存在",
	CodeUserNotExist:    "用户名不存在",
	CodeInvalidPassword: "用户或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeNeedLogin:   "需要登录",
	CodeInvalidAuth: "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
