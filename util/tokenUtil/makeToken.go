package tokenUtil

import (
	"github.com/golang-jwt/jwt/v4"
	"sanHeRecruitment/models/loginModel"
	"time"
)

const TokenExpireDuration = time.Hour * 3
const WeChatExpireDuration = time.Hour * 24 * 31

var MySecret = []byte("sanHeProducer_YAN") // 生成签名的密钥

// GenerateToken 登录成功后调用，传入UserInfo结构体
func GenerateToken(userInfo loginModel.UserInfo) (string, error) {
	expirationTime := time.Now().Add(TokenExpireDuration) // 两个小时有效期
	//expirationTime := time.Now().Add(10 * time.Second) // 两个小时有效期
	claims := &loginModel.MyClaims{
		User: userInfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "yanmingyu55@gmail.com",
		},
	}
	// 生成Token，指定签名算法和claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	if tokenString, err := token.SignedString(MySecret); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func WeChatLoginGenerateToken(userInfo loginModel.UserInfo) (string, error) {
	expirationTime := time.Now().Add(WeChatExpireDuration) // 两个小时有效期
	//expirationTime := time.Now().Add(10 * time.Second) // 两个小时有效期
	claims := &loginModel.MyClaims{
		User: userInfo,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "yanmingyu55@gmail.com",
		},
	}
	// 生成Token，指定签名算法和claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 签名
	if tokenString, err := token.SignedString(MySecret); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}
