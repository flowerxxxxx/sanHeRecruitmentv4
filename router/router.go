package router

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"sanHeRecruitment/config"
	"sanHeRecruitment/controller/admin"
	"sanHeRecruitment/controller/user"
	"sanHeRecruitment/middleware"
	"time"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//r.StaticFS("/uploadPic", http.Dir("D:\\uploadPicSaver"))
	r.Use(middleware.TlsHandler(), middleware.Cors(), middleware.LogerMiddleware())                         //
	r.Use(middleware.RateLimitMiddleware(time.Second, config.ConcurrentPeak, config.CurrentLimiterQuantum)) //限流器最高五百并发每秒添加100发
	////告诉gin框架模板引用的静态文件去哪里找
	//r.Static("/static", "static")
	////告诉gin框架去哪里找模板文件
	//r.LoadHTMLGlob("templates/*")
	//r.GET("/", controller.IndexHandler)
	//r.Use(static.Serve("/", static.LocalFile("dist", true))) //"github.com/gin-contrib/static"
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	UpLoader := r.Group("/")
	UpLoader.Use(middleware.StaticFileChecker())
	{
		UpLoader.StaticFS("/uploadPic", http.Dir(config.PicSaverPath))
	}

	//备份连接下载的静态文件访问
	BackerStatic := r.Group("/")
	BackerStatic.Use(middleware.BackerTokenRoleChecker())
	{
		BackerStatic.StaticFS("/backup", http.Dir(config.BackUpConfig.SavePath))
	}

	//用户端通用controller
	TokenNoUse := r.Group("/")
	TokenNoUse.Use()
	{
		user.UserControllerRouter(TokenNoUse)
		user.WsControllerRouter(TokenNoUse)
		user.JobControllerRouter(TokenNoUse)
		user.IdentityControllerRouter(TokenNoUse)
		user.BossControllerRouter(TokenNoUse)
		user.DataControllerRouter(TokenNoUse)
	}

	UserTokenUse := r.Group("/")
	UserTokenUse.Use(middleware.WeChatAppCheckToken())
	{
		user.UserControllerRouterToken(UserTokenUse)
		user.WsControllerRouterToken(UserTokenUse)
		user.JobControllerRouterToken(UserTokenUse)
		user.IdentityControllerRouterToken(UserTokenUse)
		user.DataControllerRouterToken(UserTokenUse)
	}

	BossTokenUse := r.Group("/boss")
	BossTokenUse.Use(middleware.WeChatAppCheckToken(), middleware.BossCheckTokenRole()) //,
	{
		user.BossControllerRouterToken(BossTokenUse)
	}

	AdminTokenNoUse := r.Group("/adminMange")
	AdminTokenNoUse.Use()
	{
		admin.AdminControllerRouter(AdminTokenNoUse)
		admin.ManageControllerRouter(AdminTokenNoUse)
		admin.DataControllerRouter(AdminTokenNoUse)
	}

	AdminTokenUse := r.Group("/adminMange")
	//AdminTokenUse.Use(middleware.JWTAuth())
	AdminTokenUse.Use(middleware.JWTAuth(), middleware.AdminCheckTokenRole())
	{
		admin.AdminControllerRouterToken(AdminTokenUse)
		admin.ManageControllerRouterToken(AdminTokenUse)
		admin.DataControllerRouterToken(AdminTokenUse)
		admin.StatisticsControllerRouterToken(AdminTokenUse)
	}

	//性能调优监视 TODO Gin自主隐藏，待优化
	authStr := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(config.ProducerUsername+":"+config.ProducerPassword)))
	pprofGroup := r.Group("/producer", func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth != authStr {
			c.Header("www-Authenticate", "Basic")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	})
	pprof.RouteRegister(pprofGroup, "sanHeRec_pprof")
	return r
}
