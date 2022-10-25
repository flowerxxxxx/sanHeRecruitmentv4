package tokenUtil

import (
	"math/rand"
	"net/http"
	"sanHeRecruitment/dao"
	"time"
)

const tokenOverdue = 2 //h

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func GetTokenFromHeader(header http.Header) (token string) {
	token = "none"
	for k, v := range header {
		if k == "Token" {
			token = v[0]
		}
	}
	return token
}

// 重置token时间
func ResetTokenTime(token string) {
	dao.Redis.Expire(token, time.Hour*tokenOverdue)
}
