package server

import (
	"errors"
	"fmt"
	"web/dao/mysql"
	"web/models"
	"web/pkg/jwt"
	"web/pkg/snowflake"
	"web/pkg/tools"
)

// 存放业务逻辑的地方

// 注册
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户存不存在
	isExist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		// 数据库查询出错
		return err
	} else if isExist {
		return errors.New("用户已存在")
	}

	// 生成UID
	userID := snowflake.GenID()
	// 构造user实例
	u := models.User{
		UserID:   uint(userID),
		Username: p.Username,
		Password: p.Password,
	}

	// 保存进数据库
	mysql.InsertUser(&u)
	return

}

// 登录
func Login(p *models.ParamLogin) (token string, err error) {
	// 查询用户是否存在
	isExist, err := mysql.CheckUserExist(p.Username)
	if err != nil {
		// 查询数据库失败
		return "", fmt.Errorf("查询数据库失败：%v", err.Error())
	}
	if !isExist {
		return "", errors.New("用户不存在")
	}

	// 验证密码
	newPassword := tools.EncryptPassword(p.Password)
	user, err := mysql.GetInfoByUserName(p.Username)
	if err != nil {
		return "", fmt.Errorf("获取用户信息失败：%v", err.Error())
	}
	if newPassword != user.Password {
		return "", errors.New("密码错误")
	}

	token, err = jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return "", fmt.Errorf("生成token失败：%v", err.Error())
	}
	// 生成jwt
	return token, nil
}
