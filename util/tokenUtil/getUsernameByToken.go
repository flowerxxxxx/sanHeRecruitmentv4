package tokenUtil

import (
	"github.com/gin-gonic/gin"
	"sanHeRecruitment/models/loginModel"
)

func GetUsernameByToken(ctx *gin.Context) (username string) {
	auth := ctx.Request.Header.Get("Authorization")
	claims, _ := ParseToken(auth)
	return claims.User.UserName
}

func GetUserClaims(ctx *gin.Context) *loginModel.MyClaims {
	auth := ctx.Request.Header.Get("Authorization")
	claims, _ := ParseToken(auth)
	return claims
}
