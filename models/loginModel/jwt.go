package loginModel

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserInfo struct {
	Id       int
	Role     string
	UserName string
}

type MyClaims struct {
	User               UserInfo
	jwt.StandardClaims // 标准Claims结构体，可设置8个标准字段
}
