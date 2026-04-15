package mysql

import (
	"web/models"
	"web/pkg/tools"
)

func CheckUserExist(username string) (isExist bool, err error) {
	var user models.User

	err = db.Table("users").
		Where("username = ?", username).
		Limit(1).
		Find(&user).Error

	if err != nil {
		return false, err
	}
	if user.ID > 0 {
		return true, err
	}
	return false, err
}

// 插入新的用户
func InsertUser(user *models.User) error {
	// 对密码进行加密
	user.Password = tools.EncryptPassword(user.Password)
	// 执行SQL语句，入库
	err := db.Table("users").Create(user).Error
	return err
}

// 通过用户名查询用户信息
func GetInfoByUserName(userName string) (models.User, error) {
	var user models.User
	err := db.
		Where("username = ?", userName).
		First(&user).Error
	return user, err
}
