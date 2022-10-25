package tokenUtil

import (
	"errors"
	"sanHeRecruitment/models/loginModel"
	"time"
)

const TokenRenewTime = 60 * 60 * 3

func RenewToken(claims *loginModel.MyClaims) (string, error) {
	// 若token过期不超过TokenRenewTime则给它续签
	if withinLimit(claims.ExpiresAt, TokenRenewTime) {
		return GenerateToken(claims.User)
	}
	return "", errors.New("登录已过期")
}

// 计算过期时间是否超过l
func withinLimit(s int64, l int64) bool {
	e := time.Now().Unix()
	return e-s < l
}
