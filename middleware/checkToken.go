package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util/tokenUtil"
	"strings"
)

var us *service.UserService

func CheckToken(c *gin.Context) {
	token := tokenUtil.GetTokenFromHeader(c.Request.Header)
	if token == "none" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "未登录",
		})
		c.Abort()
		return
	}
	stuNum := dao.Redis.Get(token).Val()
	if stuNum == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "登录信息过期，请重新登录",
		})
		c.Abort()
		return
	}
	tokenUtil.ResetTokenTime(token)
}

func JWTAuth2() gin.HandlerFunc {
	return func(context *gin.Context) {
		auth := context.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			context.Abort()
			context.String(http.StatusOK, "未登录无权限")
			return
		}
		// 校验token，只要出错直接拒绝请求
		_, err := tokenUtil.ParseToken(auth)
		if err != nil {
			context.Abort()
			message := err.Error()
			context.JSON(http.StatusOK, message)
			return
		} else {
			println("tokenUtil 正确")
		}
		context.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			// 无token直接拒绝
			c.Abort()
			c.String(http.StatusForbidden, "未登录无权限")
			return
		}
		// 校验token
		claims, err := tokenUtil.ParseToken(auth)
		//fmt.Println(claims)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				// 若过期，调用续签函数
				newToken, _ := tokenUtil.RenewToken(claims)
				if newToken != "" {
					// 续签成功給返回头设置一个newtoken字段
					c.Header("newtoken", newToken)
					c.Request.Header.Set("Authorization", newToken)
					c.Next()
					return
				}
			}
			// Token验证失败或续签失败直接拒绝请求
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status": 999,
				"msg":    "登录过期",
			})
			return
		}
		// token未过期继续执行1其他中间件
		c.Next()
	}
}

func WeChatAppCheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		claims, _ := tokenUtil.ParseToken(auth)
		username := claims.User.UserName
		if username == "" {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status": 205,
				"msg":    "未登录",
			})
			return
		}
		//if err != nil {
		//	// Token验证失败或续签失败直接拒绝请求
		//	c.Abort()
		//	c.JSON(http.StatusOK, gin.H{
		//		"status": 999,
		//		"msg":    "登录过期",
		//	})
		//	return
		//}
	}
}

func BackerTokenRoleChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		//auth := c.Request.Header.Get("Authorization")
		token2 := c.Query("token")
		//claims, _ := tokenUtil.ParseToken(auth)
		//userRole := claims.User.Role
		claims2, _ := tokenUtil.ParseToken(token2)
		userRole2 := claims2.User.Role
		if userRole2 != "admin" {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status": 210,
				"msg":    "账号权限不足，非管理员身份",
			})
			return
		}
	}
}

func AdminCheckTokenRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		claims, _ := tokenUtil.ParseToken(auth)
		userRole := claims.User.Role
		if userRole != "admin" {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status": 210,
				"msg":    "账号权限不足，非管理员身份",
			})
			return
		}
	}
}

// StaticFileChecker 静态文件夹链接checker
func StaticFileChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		authPath := c.Request.URL.Path
		if len(authPath) == len("/uploadPic/") {
			c.Abort()
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "so sad to see it :(",
			})
			return
		}
	}
}

func BossCheckTokenRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		claims, _ := tokenUtil.ParseToken(auth)
		username := claims.User.UserName
		userLevelInfo := us.QueryUserLevel(username)
		if userLevelInfo.UserLevel == 0 {
			c.Abort()
			c.JSON(http.StatusOK, gin.H{
				"status": 210,
				"msg":    "账号权限不足，非企业或者机构身份",
			})
			return
		}
	}
}
