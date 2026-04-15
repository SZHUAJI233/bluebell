package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"
	"web/dao/redis"

	"github.com/dgrijalva/jwt-go"
)

// token有效时间
const TokenExpireDuration = time.Hour * 2

// 加盐
var mySercet = []byte("夏天夏天悄悄过去")

type MyClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func GenToken(userID uint, username string) (string, error) {
	// 创建
	c := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err := token.SignedString(mySercet) // 加盐加密
	if err != nil {
		return "", fmt.Errorf("token.SignedString 加密错误：%v", err)
	}

	// 存储token
	err = redis.Set(context.Background(), "token", tokenStr, TokenExpireDuration)
	if err != nil {
		return "", fmt.Errorf("redis存储token失败: %v", err)
	}

	return tokenStr, nil
}

// 解析token
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySercet, nil // 解盐
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		// 验证token
		oldToken, err := redis.Get(context.Background(), "token")
		if err != nil {
			return nil, fmt.Errorf("redis获取token失败: %v", err)
		} else if tokenString == oldToken {
			return mc, nil
		}
	}
	return nil, errors.New("invalid token")

}
