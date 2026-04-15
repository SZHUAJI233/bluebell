package tools

import (
	"crypto/md5"
	"encoding/hex"
)

// 加盐
const secret = "huaji.com"

// 密码加密
func EncryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
