package tokenUtil

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"sanHeRecruitment/models/loginModel"
)

func ParseToken(tokenString string) (*loginModel.MyClaims, error) {
	claims := &loginModel.MyClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	// 若token只是过期claims是有数据的，若token无法解析claims无数据
	return claims, err
}

// 第二种方法通过jwt.ParseWithClaims返回的Token结构体取出Claims结构体
func ParseToken2(tokenString string) (*loginModel.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &loginModel.MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*loginModel.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token无法解析")
}
