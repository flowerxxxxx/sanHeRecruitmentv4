package user

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sanHeRecruitment/config"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/loginModel"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/models/wechatModel"
	"sanHeRecruitment/module/controllerModule"
	"sanHeRecruitment/module/recommendModule"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util"
	"sanHeRecruitment/util/saveUtil"
	"sanHeRecruitment/util/tokenUtil"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	*service.UserService
	*service.CollectionService
	*service.MsgObjService
	*service.ChatService
	*controllerModule.UserControllerModule
	*service.LabelService
	*service.JobService
	*service.EducationService
	*service.ArticleService
	*service.CountService
}

//var UserController = LoginController{
//	UserService: service.NewUserService(),
//}

func UserControllerRouter(router *gin.RouterGroup) {
	u := UserController{}
	//router.POST("/login", u.Login) //1
	router.POST("/wechatLogin", u.WeChatLogin)
}

func UserControllerRouterToken(router *gin.RouterGroup) {
	u := UserController{}
	router.POST("/testToken", u.TestToken)
	router.POST("/modifyPersonalInfo", u.ModifyPersonalInfo)
	router.POST("/uploadHeadPic", u.UploadHeadPic)
	router.POST("/uploadResume", u.UploadResume) //上传个人简历，文件版
	//收藏文章
	router.POST("/collectArticle", u.CollectArticle)
	router.POST("/cancelCollectStatus", u.CancelCollectStatus)
	router.GET("/getMsgList", u.GetMsgList)
	router.POST("/deleteMsgUser", u.DeleteOneMsgUser)
	router.POST("/showCollecArts", u.GetUserCollectArt)
	router.GET("/getPersonInfo", u.GetPersonInfo)
	//修改基础个人信息
	router.POST("/modifyBasicPersonalInfo", u.ModifyBasicPersonalInfo)
	//获取个人简历-教育经历
	router.GET("/personalResumeEdu", u.PersonalResumeEdu)
	//添加个人简历-教育经历
	router.POST("/addPersonalResumeEdu", u.AddPersonalResumeEdu)
	//修改个人简历-教育经历
	router.POST("/modifyPersonalResumeEdu", u.ModifyPersonalResumeEdu)
	//删除个人简历-教育经历
	router.POST("/deletePersonalResumeEdu", u.DeletePersonalResumeEdu)
	//修改简历-经历
	router.POST("/modifyPersonalResumePE", u.ModifyPersonalResumePE)
	//修改简历-技能
	router.POST("/modifyPersonalResumePR", u.ModifyPersonalResumePR)
	//修改简历-简介
	router.POST("/modifyPersonalResumePS", u.ModifyPersonalResumePS)
	//查询用户等级
	router.GET("/checkUserIdentityLevel", u.CheckUserIdentityLevel)
}

// CheckUserIdentityLevel 查询用户等级
func (u *UserController) CheckUserIdentityLevel(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	userLevelData := u.UserService.QueryUserLevel(username)
	controller.SuccessResp(c, "等级查询成功", userLevelData)
}

// Login 登录
func (u *UserController) Login(c *gin.Context) {
	login := struct {
		Username string `form:"username" json:"username"`
		Password string `form:"password" json:"password"`
	}{}
	errX := c.ShouldBind(&login)
	if errX != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errX.Error()})
		return
	}
	username := login.Username
	password := login.Password
	//fmt.Println(username, password)
	err := u.UserService.Login(username, password)
	if err != nil {
		if err == service.ErrorPasswordWrong {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "登陆失败,密码错误",
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
	userInfo, err := u.UserService.GetUserInfo(username, c.Request.Host)
	userInfo2 := loginModel.UserInfo{
		Id:       userInfo.User_id,
		Role:     userInfo.Role,
		UserName: username,
	}
	token, _ := tokenUtil.GenerateToken(userInfo2)
	//userData := &models.LoginInfo{
	//	AvatarUrl: userInfo.Head_pic,
	//	Nickname:  userInfo.Nickname,
	//	UserLevel: userInfo.UserLevel,
	//}
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "登录成功",
		"token":    token,
		"userData": userInfo,
	})
}

// GetMsgList 用户获取消息列表
func (u *UserController) GetMsgList(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	fromUsername := claims.User.UserName
	msgList := u.MsgObjService.GetMsgList(fromUsername, c.Request.Host)
	//finalMsgList := u.ChatService.QueryEveryLastMsg(msgList)
	msgList = u.UserControllerModule.MsgListSortByStartTime(msgList)
	finalData := []interface{}{}
	//fmt.Println(websocketModel.ReceiveMsgManager.Clients)
	for i, j := 0, len(msgList); i < j; i++ {
		//_, ok := websocketModel.ReceiveMsgManager.Clients[msgList[i].ToUsername]
		if _, ok := websocketModel.ReadRecManClient(msgList[i].ToUsername); !ok {
			msgList[i].Online = 0
		} else {
			msgList[i].Online = 1
		}
	}
	////消息量化排序
	for _, msg := range msgList {
		if msg.Read == 0 {
			//fmt.Println(msg)
			finalData = util.SliAddFromHead(finalData, msg)
		} else {
			finalData = append(finalData, msg)
		}
		//fmt.Println(finalData)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "get msg done",
		"data":   finalData,
	})
}

// GetUserCollectArt 查看用户收藏的内容
func (u *UserController) GetUserCollectArt(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	classification := recJson["classification"].(string)
	pageNum := recJson["pageNum"].(string)
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	collecArts := u.JobService.GetCollectArts(username, classification, c.Request.Host, pageNumInt)
	totalPage := u.CountService.UserCollectTP(username, classification)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "收藏信息获取成功",
		"data":   collecArts,
		"total":  totalPage,
	})
}

// DeleteOneMsgUser 用户删除消息列表的单条聊天
func (u *UserController) DeleteOneMsgUser(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	msgId := recJson["msgId"].(string)
	fmt.Println(msgId)
	fmt.Println(username)
	msgIdInt, _ := strconv.Atoi(msgId)
	err := u.MsgObjService.DeleteMsg(msgIdInt, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "服务器逻辑层错误，该记录不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "删除成功",
	})
}

// TestToken 测试token
func (u *UserController) TestToken(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	fmt.Println(claims)
	x := mysqlModel.AddLabel()
	//xx := mysqlModel.FindArticles()
	//models.GetAllTypes()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"x":      x,
		//"xx":     xx,
		"msg": "success",
	})
}

// GetPersonInfo 个人中心获取个人信息
func (u *UserController) GetPersonInfo(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	userData, _ := u.UserService.GetUserInfo(username, c.Request.Host)
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "个人信息获取成功",
		"userData": userData,
	})
}

// CollectArticle 收藏文章
func (u *UserController) CollectArticle(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	ArtId := recJson["art_id"].(string)
	CareerJobId := recJson["career_job_id"].(string)
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	if username == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 210,
			"msg":    "未登录",
		})
		return
	}
	artIdInt, _ := strconv.Atoi(ArtId)
	labelInfo, err := u.LabelService.QueryLabelInfo(CareerJobId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "结果不存在",
		})
		return
	}
	userCollecNum := u.CountService.UserCollectTotal(username)
	if userCollecNum >= 300 {
		c.JSON(http.StatusOK, gin.H{
			"status": 203,
			"msg":    "收藏达到上限",
		})
		return
	}
	err = u.CollectionService.CollectArticle(username, artIdInt, labelInfo.Type)
	go func(artId string) {
		u.ArticleService.AddArtCollectNum(artId)
		recommendModule.DealArtRecommendWeight(artId)
	}(ArtId)
	if err != nil {
		if err == service.MysqlErr {
			c.JSON(http.StatusOK, gin.H{
				"status": 216,
				"msg":    "服务器错误",
			})
			log.Println("CollectArticle", err, "recJson:", recJson)
			return
		}
		if err == service.Existed {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "收藏失败，该文章已被收藏",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 217,
			"msg":    "服务器错误",
		})
		log.Println("CollectArticle", err, "recJson:", recJson)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "收藏成功",
	})
}

// CancelCollectStatus 用户取消收藏
func (u *UserController) CancelCollectStatus(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	if username == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": 210,
			"msg":    "取消收藏失败，未登录",
		})
		return
	}
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	artId := recJson["art_id"].(string)
	artIdint, _ := strconv.Atoi(artId)
	u.CollectionService.DeleteRecord(username, artIdint)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "取消收藏成功",
	})
}

// ModifyPersonalInfo 修改用户个人信息
func (u *UserController) ModifyPersonalInfo(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	email := recJson["email"].(string)
	gender := recJson["gender"].(string)
	intendedPosition := recJson["intended_position"].(string)
	age := recJson["age"].(string)
	name := recJson["name"].(string)
	ageInt, _ := strconv.Atoi(age)
	nickname := recJson["nickname"].(string)
	u.UserService.ModifyPersonalInfo(username, email, gender, intendedPosition, name, nickname, ageInt)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "个人信息修改成功",
	})
}

func (u *UserController) ModifyBasicPersonalInfo(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	name := recJson["name"].(string)
	nickname := recJson["nickname"].(string)
	gender := recJson["gender"].(string)
	phoneNumber := recJson["phone"].(string)
	intendedPosition := recJson["intended_position"].(string)
	email := recJson["email"].(string)
	ageStr := recJson["age"].(string)
	age, _ := strconv.Atoi(ageStr)
	err := u.UserService.ModifyBasicPersonalInfo(username, email, name, gender, nickname, intendedPosition, phoneNumber, age)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "基本信息修改失败，服务器错误",
		})
		log.Println("CollectArticle", err, "recJson:", recJson)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "基本信息修改成功",
	})
}

// UploadHeadPic 上传个人头像
func (u *UserController) UploadHeadPic(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	file, err := c.FormFile("headPic")
	if err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fileFormat := file.Filename[strings.Index(file.Filename, "."):]
	if fileFormat != ".jpg" && fileFormat != ".png" && fileFormat != ".jpeg" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "支持格式:jpg,png,jpeg",
		})
		return
	}
	userInfo, err := u.UserService.GetUserInfo(username, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "用户不存在",
		})
		return
	}
	uuid := util.GetUUID()
	newFileName := uuid + "-" + strconv.Itoa(int(time.Now().Unix())) + fileFormat
	HeadPic := "uploadPic/" + newFileName
	fileAddr := config.PicSaverPath + "/" + newFileName
	//es := saveUtil.SaveCompressFile(file, fileAddr)
	//if es != nil {
	//	c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
	//	log.Println("Upload pic failed，err:", err)
	//	return
	//}
	es := saveUtil.SaveCompressCutImg(file, fileAddr)
	if es != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		log.Println("Upload pic failed，err:", err)
		return
	}
	err = u.UserService.ModifyHeadPic(username, HeadPic, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 216,
			"msg":    "服务器错误",
		})
		log.Println("UploadHeadPic", err)
		return
	}
	oldHeadPicAddr := userInfo.Head_pic
	if oldHeadPicAddr != "" {
		oldHeadPicAddr = oldHeadPicAddr[strings.LastIndex(oldHeadPicAddr, "/")+1:]
		_ = os.Remove(config.PicSaverPath + "/" + oldHeadPicAddr)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "头像上传成功",
	})
}

func FileToImage(header *multipart.FileHeader) (image image.Image, err error) {
	file, err := header.Open()
	ext := strings.ToLower(path.Ext(header.Filename))
	switch ext {
	case "jpeg", "jpg":
		image, err = jpeg.Decode(bufio.NewReader(file))
	case "png":
		image, err = png.Decode(bufio.NewReader(file))
	}
	return image, err
}

// PersonalResumeEdu 获取个人教育经历-Token
func (u *UserController) PersonalResumeEdu(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	perEduInfo := u.EducationService.QueryPersonalEdu(username)
	controller.SuccessResp(c, "简历-教育查询成功", perEduInfo)
}

// AddPersonalResumeEdu 添加教育经历
func (u *UserController) AddPersonalResumeEdu(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	school := recJson["school"].(string)
	major := recJson["major"].(string)
	start_time := recJson["start_time"].(string)
	end_time := recJson["end_time"].(string)
	degree := recJson["degree"].(string)
	err := u.EducationService.AddPersonalEdu(username, school, major, start_time, end_time, degree)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 217,
			"msg":    "简历-教育经历添加失败，服务器错误",
		})
		log.Println("CollectArticle", err, "recJson:", recJson)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "简历-教育经历添加成功",
	})
}

// ModifyPersonalResumePS  修改简历-技能
func (u *UserController) ModifyPersonalResumePS(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	personalSkill := recJson["personal_skill"].(string)
	err := u.UserService.ModifyPersonalResumePS(username, personalSkill)
	if err != nil {
		controller.ErrorResp(c, 201, "个人简历-技能修改失败")
		return
	}
	controller.SuccessResp(c, "个人简历-技能修改成功")
}

// ModifyPersonalResumePE 修改简历-经历
func (u *UserController) ModifyPersonalResumePE(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	projectExper := recJson["personal_experience"].(string)
	err := u.UserService.ModifyPersonalResumePE(username, projectExper)
	if err != nil {
		controller.ErrorResp(c, 201, "个人简历-经历修改失败")
		return
	}
	controller.SuccessResp(c, "个人简历-经历修改成功")
}

// ModifyPersonalResumePR 修改个人简历-个人简介部分
func (u *UserController) ModifyPersonalResumePR(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	personalRes := recJson["personal_resume"].(string)
	err := u.UserService.ModifyPersonalResumePR(username, personalRes)
	if err != nil {
		controller.ErrorResp(c, 201, "个人简历-个人简介修改失败")
		return
	}
	controller.SuccessResp(c, "个人简历-个人简介修改成功")
}

// DeletePersonalResumeEdu 删除个人简历的教育经历
func (u *UserController) DeletePersonalResumeEdu(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	id := recJson["id"].(string)
	idInt, _ := strconv.Atoi(id)
	err := u.EducationService.DeletePersonalResumeEdu(idInt, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "简历-教育经历删除失败，服务器错误",
		})
		log.Println("CollectArticle", err, "recJson:", recJson)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "简历-教育经历删除成功",
	})
}

// ModifyPersonalResumeEdu 修改教育经历
func (u *UserController) ModifyPersonalResumeEdu(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	id := recJson["id"].(string)
	school := recJson["school"].(string)
	major := recJson["major"].(string)
	start_time := recJson["start_time"].(string)
	end_time := recJson["end_time"].(string)
	degree := recJson["degree"].(string)
	idInt, _ := strconv.Atoi(id)
	err := u.EducationService.ModifyPersonalResumeEdu(idInt, username, school, major, start_time, end_time, degree)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "简历-教育经历修改失败，服务器错误",
		})
		log.Println("CollectArticle", err, "recJson:", recJson)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "简历-教育经历修改成功",
	})
}

// UploadResume 上传个人简历-文件版
func (u *UserController) UploadResume(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	file, err := c.FormFile("resume")
	if err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fileFormat := file.Filename[strings.Index(file.Filename, "."):]
	if fileFormat != ".doc" && fileFormat != ".pdf" {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "简历仅支持pdf或者doc格式",
		})
		return
	}
	userInfo, err := u.UserService.GetUserInfo(username, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "用户不存在",
		})
		return
	}
	oldResume := userInfo.Resume
	if oldResume != "" {
		oldResume = oldResume[strings.LastIndex(oldResume, "/")+1:]
		_ = os.Remove("./uploadPic/" + oldResume)
	}
	uuid := util.GetUUID()
	newFileName := uuid + "-" + username + fileFormat
	Resume := "https://" + c.Request.Host + "/uploadPic/" + newFileName
	fileAddr := "./uploadPic/" + newFileName
	if err := c.SaveUploadedFile(file, fileAddr); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	err = u.UserService.ModifyResume(username, Resume, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "服务器错误",
		})
		log.Println("UploadResume", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "文件上传成功",
	})
}

// WXLogin 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况,作为微信登录的调用函数
func WXLogin(code string) (*wechatModel.WXLoginResp, error) {
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	//url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=${code}&grant_type=authorization_code"
	appId := config.WechatAppid
	secret := config.WechatSecret
	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	url = fmt.Sprintf(url, appId, secret, code)
	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := wechatModel.WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		return nil, err
	}
	// 判断微信接口返回的是否是一个异常情况
	if wxResp.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("ErrCode:%v  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
	}

	//src, err2 := wechatUtil.Dncrypt(encryptedData, wxResp.SessionKey, iv)
	//if err2 != nil {
	//	log.Println("wechatUtil.Dncrypt failed,err:", err2)
	//}
	//var s = map[string]interface{}{}
	//json.Unmarshal([]byte(src), &s)
	//fmt.Println(src)
	//fmt.Println("unionId:", s["unionId"])
	//wxResp.UnionId = s["unionId"].(string)
	return &wxResp, nil
}

// WeChatLogin 微信小程序登录
func (u *UserController) WeChatLogin(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	code := recJson["code"].(string)
	rawData := recJson["rawData"].(string)
	//iv := recJson["iv"].(string)
	//encryptedData := recJson["encryptedData"].(string)
	fmt.Println(recJson)
	var recRowData wechatModel.RawData
	_ = json.Unmarshal([]byte(rawData), &recRowData)
	//fmt.Println(recRowData)
	// 根据code获取 openID 和 session_key
	wxLoginResp, err := WXLogin(code)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": "server fail",
		})
		return
	}
	err = u.UserService.WechatLogin(wxLoginResp.OpenId)
	if err != nil {
		if err == service.ErrorNotExisted {
			_ = u.UserService.WechatRegister(wxLoginResp.OpenId, recRowData.NickName, recRowData.AvatarUrl)
		} else {
			c.JSON(215, gin.H{
				"msg": "服务器错误",
			})
			log.Println("CollectArticle", err, "recJson:", recJson)
			return
		}
	}
	userInfo, err := u.UserService.GetUserInfo(wxLoginResp.OpenId, c.Request.Host)
	userInfo2 := loginModel.UserInfo{
		Id:       userInfo.User_id,
		Role:     userInfo.Role,
		UserName: wxLoginResp.OpenId,
	}
	token, _ := tokenUtil.WeChatLoginGenerateToken(userInfo2)
	c.JSON(http.StatusOK, gin.H{
		"status":   200,
		"msg":      "登录成功",
		"token":    token,
		"userData": userInfo,
	})
	////fmt.Println(tokenUtil)
	//// 这里用openid和sessionkey的串接 进行MD5之后作为该用户的自定义登录态
	//mySession := GetMD5Encode(wxLoginResp.OpenId + wxLoginResp.SessionKey)
	//fmt.Println(wxLoginResp.OpenId)
	//// 接下来可以将openid 和 sessionkey, mySession 存储到数据库中,
	//// 但这里要保证mySession 唯一, 以便于用mySession去索引openid 和sessionkey
	//c.String(200, mySession)
	//c.JSON(http.StatusOK, gin.H{
	//	"msg": "succ",
	//})
}

// GetMD5Encode 将一个字符串进行MD5加密后返回加密后的字符串
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
