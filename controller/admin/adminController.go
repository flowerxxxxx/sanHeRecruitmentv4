package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/loginModel"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/tokenUtil"
	"strconv"
)

//本文件的router包含管理员的的相关身份管理1

type AdminController struct {
	*service.UserService
	*service.CountService
}

func AdminControllerRouter(router *gin.RouterGroup) {
	ac := AdminController{}
	//登录
	router.POST("/adminLogin", ac.AdminLogin)
}

func AdminControllerRouterToken(router *gin.RouterGroup) {
	ac := AdminController{}
	//修改管理员密码
	router.POST("/ModifyPassword", ac.ModifyPassWord)
	//添加管理员
	router.POST("/addAdminer", ac.AddAdminer)
	//删除管理员
	router.POST("/deleteAdminer", ac.DeleteAdminer)
	//获取所有管理员信息
	router.GET("/showAdminInfos/:pageNum", ac.ShowAdminInfos)
}

// ShowAdminInfos 查看管理员信息
func (ac *AdminController) ShowAdminInfos(c *gin.Context) {
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	adminInfos, _ := ac.UserService.GetAdminerInfos(pageNumInt)
	totalPage := ac.CountService.AdminInfosTP()
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "查询成功",
		"totalPage": totalPage,
		"data":      adminInfos,
	})
}

// ModifyPassWord 主管理员修改密码
func (ac *AdminController) ModifyPassWord(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	newPwd := recJson["newPassword"].(string)
	desUsername := recJson["desUsername"].(string)

	ModifierInfo, _ := ac.UserService.GetUserInfo(tokenUtil.GetUsernameByToken(c), c.Request.Host)
	if ModifierInfo.UserLevel > 100 {
		controller.ErrorResp(c, 201, "管理员权限等级不足")
		return
	}
	desUserInfo, _ := ac.UserService.GetUserInfo(desUsername, c.Request.Host)
	if desUserInfo.Role != "admin" || desUserInfo.UserLevel <= 100 {
		controller.ErrorResp(c, 202, "目标用户角色等级均不符")
		return
	}
	newMd5Pwd := sqlUtil.GenMD5Password(newPwd)
	err := ac.UserService.AdminModifyPwd(desUsername, newMd5Pwd)
	if err != nil {
		controller.ErrorResp(c, 211, "修改失败，服务器错误")
		return
	}
	controller.SuccessResp(c, "密码修改成功")
}

func (ac *AdminController) AddAdminer(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	username := recJson["username"].(string)
	password := recJson["password"].(string)
	nickname := recJson["nickname"].(string)
	ModifierInfo, _ := ac.UserService.GetUserInfo(tokenUtil.GetUsernameByToken(c), c.Request.Host)
	if ModifierInfo.UserLevel > 100 {
		controller.ErrorResp(c, 201, "管理员权限等级不足")
		return
	}
	_, err := ac.UserService.GetUserInfo(username, c.Request.Host)
	if err == nil {
		controller.ErrorResp(c, 202, "用户名已存在")
		return
	}
	md5Pwd := sqlUtil.GenMD5Password(password)
	err = ac.UserService.AddAdminer(username, md5Pwd, nickname)
	if err != nil {
		controller.ErrorResp(c, 213, "创建失败，服务器错误")
		return
	}
	controller.SuccessResp(c, "添加成功")
}

func (ac *AdminController) DeleteAdminer(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	username := recJson["username"].(string)
	ModifierInfo, _ := ac.UserService.GetUserInfo(tokenUtil.GetUsernameByToken(c), c.Request.Host)
	if ModifierInfo.UserLevel > 100 {
		controller.ErrorResp(c, 201, "管理员权限等级不足")
		return
	}
	desUserInfo, _ := ac.UserService.GetUserInfo(username, c.Request.Host)
	if desUserInfo.Role != "admin" || desUserInfo.UserLevel <= 100 {
		controller.ErrorResp(c, 202, "目标用户角色等级均不符")
		return
	}
	err := ac.UserService.DeleteAdminer(username)
	if err != nil {
		controller.ErrorResp(c, 213, "删除失败，服务器错误")
		return
	}
	controller.SuccessResp(c, "删除成功")
}

// AdminLogin 登录
func (ac *AdminController) AdminLogin(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	username := recJson["username"].(string)
	password := recJson["password"].(string)
	passwordMD5 := sqlUtil.GenMD5Password(password)
	err := ac.UserService.Login(username, passwordMD5)
	if err != nil {
		if err == service.ErrorPasswordWrong {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "密码错误",
			})
			return
		}
		if err == service.ErrorNotExisted {
			c.JSON(http.StatusOK, gin.H{
				"status": 202,
				"msg":    "登陆失败,不存在该用户名",
			})
			return
		}
	}
	userInfo, err := ac.UserService.GetUserInfo(username, c.Request.Host)
	userInfo2 := loginModel.UserInfo{
		Id:       userInfo.User_id,
		Role:     userInfo.Role,
		UserName: username,
	}
	if userInfo.Role != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "登陆失败,用户身份非管理员",
		})
		return
	}
	token, _ := tokenUtil.GenerateToken(userInfo2)
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "登录成功",
		"token":    token,
		"userData": userInfo,
	})
}
