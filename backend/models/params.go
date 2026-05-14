package models

// 定义请求的参数结构体

// 注册参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// 登录参数
type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 投票数据
type ParamVoteData struct {
	// UserID 请求中获取
	PostID string `json:"post_id" binding:"required"`         // 帖子id
	Vote   int8   `json:"vote,string" binding:"oneof=1 0 -1"` // 赞成票（1）反对票（-1）未投票（0）
}
