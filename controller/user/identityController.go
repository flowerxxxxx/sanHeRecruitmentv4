package user

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"reflect"
	"sanHeRecruitment/biz/websocketBiz"
	"sanHeRecruitment/config"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/messageUtil"
	"sanHeRecruitment/util/saveUtil"
	"sanHeRecruitment/util/timeUtil"
	"sanHeRecruitment/util/tokenUtil"
	"sanHeRecruitment/util/uploadUtil"
	"strconv"
	"strings"
	"time"
)

// IdentityController 主结构体引用器
type IdentityController struct {
	*mysqlService.UserService
	*mysqlService.CompanyService
	*mysqlService.UpgradeService
	*mysqlService.JobService
	*mysqlService.ArticleService
	*mysqlService.VoucherService
}

// IdentityControllerRouter 身份控制台
func IdentityControllerRouter(router *gin.RouterGroup) {
	bc := IdentityController{}
	//上传凭证照片到服务器
	router.POST("/uploadPicVoucher", bc.SavePicVouchers)
	//上传新公司照片到服务器
	router.POST("/uploadComHeadVoucher", bc.SavePicHeadVouchers)
	// 根据url删除本地已经上传的照片凭证
	router.POST("/deletePicVoucher", bc.DeletePicVoucher)
	//实时根据现有参数模糊查询公司名称
	router.GET("/fuzzyQueryCompanies/:companyLevel/:companyName", bc.FuzzyQueryCompanies)
	//确定绑定状态
	router.POST("/confirmCompanyInfo", bc.ConfirmCompanyInfo)
}

// IdentityControllerRouterToken 身份控制台
func IdentityControllerRouterToken(router *gin.RouterGroup) {
	bc := IdentityController{}
	//核验升级资质
	router.GET("/checkUpgradeQualification", bc.CheckUpgradeQualification)
	//上传公司信息
	router.POST("/saveCompanyInfo", bc.SaveCompanyInfo)
	//上传身份升级请求
	router.POST("/uploadUpgradeRequest", bc.SaveUpgradeRequest)
	//存储升级凭证
	router.POST("/uploadUpgradeVouchers", bc.SaveUpgradeVouchers)
	//查询升级完成度
	router.GET("/checkUpgradeCompletion", bc.CheckUpgradeCompletion)
	// 切换身份 一端口通用
	router.POST("/changUserIdentityLevel", bc.ChangUserIdentityLevel)
	// 注销个人企业/服务机构身份
	router.POST("/cancelOrgIdentity", bc.CancelOrgIdentity)

}

// CancelOrgIdentity 注销个人企业/服务机构身份
func (bc *IdentityController) CancelOrgIdentity(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	err := bc.UserService.ResetUserLevel(username)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		log.Println("cancelOrgIdentity", err)
		return
	}
	userInfo := tokenUtil.GetUserClaims(c)
	err = bc.ArticleService.BatchDeletePub(userInfo.User.Id)
	if err != nil {
		log.Println("[errorLog]CancelOrgIdentity-BatchDeletePub error,", err)
		controller.ErrorResp(c, 216, "操作失败，服务器错误")
		log.Println("cancelOrgIdentity", err)
		return
	}
	controller.SuccessResp(c, "注销成功")
	return
}

// ChangUserIdentityLevel 切换身份 一端口通用
func (bc *IdentityController) ChangUserIdentityLevel(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	changeLevel := recJson["changeLevel"].(string)
	LevelInfo := bc.UserService.QueryUserLevel(username)
	if LevelInfo.UserLevel == 0 {
		controller.ErrorResp(c, 201, "普通用户不能操作")
		return
	}
	changeLevelInt, _ := strconv.Atoi(changeLevel)
	if changeLevelInt == 0 {
		err := bc.UserService.ModifyUserIdentityPin(username, changeLevelInt)
		if err != nil {
			controller.ErrorResp(c, 214, "服务器错误")
			log.Println("ChangUserIdentityLevel", err, "recJson:", recJson)
			return
		}
	} else {
		if changeLevelInt == LevelInfo.UserLevel {
			err := bc.UserService.ModifyUserIdentityPin(username, changeLevelInt)
			if err != nil {
				controller.ErrorResp(c, 215, "服务器错误")
				log.Println("ChangUserIdentityLevel", err, "recJson:", recJson)
				return
			}
		} else {
			controller.ErrorResp(c, 202, "目标身份与本人不符")
			return
		}
	}
	controller.SuccessResp(c, "身份切换成功")
}

// CheckUpgradeCompletion 查询升级完成度
func (bc *IdentityController) CheckUpgradeCompletion(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	upgradeInfo, err := bc.UpgradeService.QueryUpgradeInfoByUsername(username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"CompletionLevel": 0,
			"msg":             "完善信息阶段",
			"upgradeLevel":    0,
			"status":          200,
		})
		return
	}
	upgradeVouchers := bc.VoucherService.QueryUpgradeVouchers(upgradeInfo.CompanyId, username, c.Request.Host, upgradeInfo.TimeId)
	if len(upgradeVouchers) == 0 {
		if upgradeInfo.CompanyExist == 0 {
			c.JSON(http.StatusOK, gin.H{
				"CompletionLevel": 2,
				"msg":             "上传凭证阶段",
				"action":          "register",
				"targetCompany":   upgradeInfo.CompanyId,
				//传送 -企业用户 -服务机构
				"upgradeLevel": upgradeInfo.TargetLevel,
				"time_id":      upgradeInfo.TimeId,
				"status":       200,
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"CompletionLevel": 2,
				"msg":             "上传凭证阶段",
				"action":          "bind",
				"targetCompany":   upgradeInfo.CompanyId,
				"upgradeLevel":    upgradeInfo.TargetLevel,
				"time_id":         upgradeInfo.TimeId,
				"status":          200,
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"CompletionLevel": 3,
			"msg":             "审核阶段",
			"upgradeLevel":    upgradeInfo.TargetLevel,
			"status":          200,
		})
		return
	}
}

// SaveUpgradeVouchers 存储升级凭证
func (bc *IdentityController) SaveUpgradeVouchers(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	upVouchers := recJson["upVouchers"]
	companyId := recJson["companyId"].(string)
	TimeId := recJson["time_id"].(string)
	TimeId64, err := strconv.ParseInt(TimeId, 10, 64)
	if err != nil {
		controller.ErrorResp(c, 210, "TimeId format err")
		return
	}
	companyIdInt, _ := strconv.Atoi(companyId)
	username := tokenUtil.GetUsernameByToken(c)
	Vouchers := reflect.ValueOf(upVouchers)
	if Vouchers.Len() > 3 {
		controller.ErrorResp(c, 201, "超出上传上限")
		return
	}
	for i := 0; i < Vouchers.Len(); i++ {
		upVoucher := Vouchers.Index(i).Interface().(string)
		saveFlag := strings.Index(upVoucher, "uploadPic/")
		if saveFlag == -1 {
			controller.ErrorResp(c, 203, "图片路径错误")
			return
		}
		upVoucher = upVoucher[saveFlag:]
		err := bc.VoucherService.SaveUpgradeVoucher(username, upVoucher, companyIdInt, TimeId64)
		err = bc.UpgradeService.ModifyUpgradeShow(username, companyIdInt)
		if err != nil {
			controller.ErrorResp(c, 215, "存储失败，服务器异常")
			log.Println("SaveUpgradeVouchers", err, "recJson:", recJson)
			return
		}
	}
	//upgradeInfo, err := bc.UpgradeService.QueryUpgradeInfoByTimeId(TimeId64, companyIdInt)
	//if upgradeInfo.TargetLevel == 2 {
	//	bc.UpgradeService.ModifyUpgradeQualification(upgradeInfo.ID, 1)
	//	_ = bc.UserService.ModifyPersonalInfoByUpgrade(upgradeInfo.FromUsername,
	//		upgradeInfo.CompanyId, upgradeInfo.TargetLevel)
	//	go func() {
	//		applyUserBasicInfo, _ := bc.UserService.QueryUserBasicInfo(upgradeInfo.FromUsername)
	//		succTem := messageUtil.UpgradeApplySuccessTem(
	//			applyUserBasicInfo.Name,
	//			applyUserBasicInfo.Gender,
	//			upgradeInfo.ApplyTime,
	//		)
	//		websocketBiz.SysMsgPusher(upgradeInfo.FromUsername, succTem)
	//	}()
	//}
	controller.SuccessResp(c, "凭证存储成功")
}

// ConfirmCompanyInfo 完善信息模块确认公司并进行绑定验证
func (bc *IdentityController) ConfirmCompanyInfo(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	companyId := recJson["companyId"].(string)
	if companyId == "0" {
		c.JSON(http.StatusOK, gin.H{
			"ifExist": false,
			"msg":     "无匹配公司，无法绑定",
			"status":  200,
		})
		return
	}
	companyIdInt, _ := strconv.Atoi(companyId)
	companyInfo, err := bc.CompanyService.QueryCompanyInfoById(companyIdInt, c.Request.Host)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"ifExist": false,
			"msg":     "无公司相关信息，无法绑定",
			"status":  200,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"ifExist": true,
		"msg":     "公司相关信息查询成功，允许绑定",
		"data":    companyInfo,
		"status":  200,
	})
}

// FuzzyQueryCompanies 实时根据现有参数模糊查询公司名称
func (bc *IdentityController) FuzzyQueryCompanies(c *gin.Context) {
	fuzzyComName := c.Param("companyName")
	companyLevel := c.Param("companyLevel")
	comInfos := bc.CompanyService.FuzzyQueryCompanies(fuzzyComName, companyLevel, -1)
	controller.SuccessResp(c, "模糊检索成功", comInfos)
}

// DeletePicVoucher 根据url删除本地已经上传的照片凭证
func (bc *IdentityController) DeletePicVoucher(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	picUrl := recJson["picUrl"].(string)
	posIndex := strings.Index(picUrl, "/uploadPic")
	posVoucher := strings.Index(picUrl, "voucher")
	if posIndex == -1 || posVoucher == -1 {
		controller.ErrorResp(c, 201, "图片路径错误")
	}
	finalPicUrl := config.PicSaverPath + picUrl[posIndex+10:]
	go func() {
		err := os.Remove(finalPicUrl)
		if err != nil {
			log.Println("file remove Error!")
			log.Printf("%s", err)
		}
	}()
	controller.SuccessResp(c, "图片凭证删除成功")
}

// SaveUpgradeRequest 上传身份升级请求
func (bc *IdentityController) SaveUpgradeRequest(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	userPresident := recJson["userPresident"].(string)
	CompanyId := recJson["companyId"].(string)
	TargetLevel := recJson["targetLevel"].(string)
	CompanyIdInt, _ := strconv.Atoi(CompanyId)
	TargetLevelInt, _ := strconv.Atoi(TargetLevel)
	applyTime := timeUtil.GetNowTimeFormat()
	NowTime := time.Now()
	TimeId := NowTime.Unix()

	//err := bc.UpgradeService.AddUpgradeInfo(username, TargetLevelInt, CompanyIdInt, 1, applyTime, TimeId)
	err := bc.UpgradeService.UpgradeInfoChangerUser(username, userPresident, TargetLevelInt, CompanyIdInt, 1, applyTime, TimeId)
	//_ = bc.UserService.ModifyPersonalPresident(username, userPresident)

	go func() {
		applyUserBasicInfo, _ := bc.UserService.QueryUserBasicInfo(username, c.Request.Host)
		succTem := messageUtil.UpgradeApplySuccessTem(
			applyUserBasicInfo.Name,
			applyUserBasicInfo.Gender,
			NowTime,
		)
		websocketBiz.SysMsgPusher(username, succTem)
	}()

	if err != nil {
		controller.ErrorResp(c, 215, "升级凭证上传失败，服务器错误")
		log.Println("cancelOrgIdentity", err, "recJson:", recJson)
		return
	}
	controller.SuccessResp(c, "凭证上传成功")
}

// SaveCompanyInfo 上传公司信息
func (bc *IdentityController) SaveCompanyInfo(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	TargetLevel := recJson["targetLevel"].(string)
	comHeadPic := recJson["comHeadPic"].(string)
	companyName := recJson["companyName"].(string)
	description := recJson["description"].(string) //公司简称
	scaleTag := recJson["scaleTag"].(string)       //融资情况
	personScale := recJson["personScale"].(string) //人员规模
	address := recJson["address"].(string)
	userPresident := recJson["userPresident"].(string)
	phone := recJson["phone"].(string)
	TargetLevelInt, err := strconv.Atoi(TargetLevel)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	UpdateTime := timeUtil.GetNowTimeFormat()
	if exist := strings.Index(companyName, "/"); exist != -1 || len(companyName) == 0 {
		controller.ErrorResp(c, 202, "公司名不符合规则")
		return
	}
	companyInfo, err := bc.CompanyService.QueryCompanyInfoByName(companyName)
	if err == nil {
		controller.ErrorResp(c, 203, "公司已存在", companyInfo)
		return
	}

	var newComHeadPic string
	if comHeadPic != "" {
		saveFlag := strings.Index(comHeadPic, "uploadPic/")
		if saveFlag == -1 {
			controller.ErrorResp(c, 204, "图片路径错误")
			return
		}
		newComHeadPic = comHeadPic[saveFlag:]
	}

	ea := bc.CompanyService.AddCompanyInfoTX(username, newComHeadPic, companyName, description,
		scaleTag, personScale, address, phone, userPresident, UpdateTime, TargetLevelInt)

	if ea != nil {
		controller.ErrorResp(c, 215, "服务器错误，公司信息上传失败")
		log.Println("SaveCompanyInfo", err, "recJson:", recJson)
		return
	}
	controller.SuccessResp(c, "公司信息上传成功")
	return

	//err = bc.CompanyService.AddCompanyInfo(username, newComHeadPic, companyName, description,
	//	scaleTag, personScale, address, phone, UpdateTime, TargetLevelInt)
	//newCompanyInfo, err := bc.CompanyService.QueryCompanyInfoByName(companyName)
	//TimeId := time.Now().Unix()
	//err = bc.UpgradeService.AddUpgradeInfo(username, TargetLevelInt, newCompanyInfo.ComId, 0, UpdateTime, TimeId)
	//go func() {
	//	errx := bc.UserService.ModifyPersonalPresident(username, userPresident)
	//	if errx != nil {
	//		log.Println("SaveCompanyInfo ModifyPersonalPresident err.info:", recJson)
	//	}
	//}()
}

// SavePicHeadVouchers 上传照片凭证（公司照片，凭证）
func (bc *IdentityController) SavePicHeadVouchers(c *gin.Context) {
	file, err := c.FormFile("pic_voucher")
	if err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fileFormat := file.Filename[strings.Index(file.Filename, "."):]
	//if fileFormat != ".jpg" && fileFormat != ".png" {
	//	controller.ErrorResp(c, 202, "仅支持jpg、png格式")
	//	return
	//}
	if formatFlag := uploadUtil.FormatJudge(fileFormat, ".jpg", ".png", ".jpeg"); !formatFlag {
		controller.ErrorResp(c, 202, "图片格式不支持")
		return
	}
	fileUrl, fileAddr := uploadUtil.SaveFormat("voucher"+fileFormat, c.Request.Host)
	if err := saveUtil.SaveCompressCutImg(file, fileAddr); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	controller.SuccessResp(c, "照片凭证上传成功", formatUtil.GetPicHeaderBody(c.Request.Host, fileUrl))
}

// SavePicVouchers 上传照片凭证（公司照片，凭证）
func (bc *IdentityController) SavePicVouchers(c *gin.Context) {
	file, err := c.FormFile("pic_voucher")
	if err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fileFormat := file.Filename[strings.Index(file.Filename, "."):]
	//if fileFormat != ".jpg" && fileFormat != ".png" {
	//	controller.ErrorResp(c, 202, "仅支持jpg、png格式")
	//	return
	//}
	if formatFlag := uploadUtil.FormatJudge(fileFormat, ".jpg", ".png", ".jpeg", ".webp"); !formatFlag {
		controller.ErrorResp(c, 202, "图片格式不支持")
		return
	}
	fileUrl, fileAddr := uploadUtil.SaveFormat("voucher"+fileFormat, c.Request.Host)
	if err := c.SaveUploadedFile(file, fileAddr); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	controller.SuccessResp(c, "照片凭证上传成功", formatUtil.GetPicHeaderBody(c.Request.Host, fileUrl))
}

// CheckUpgradeQualification 检测用户升级资格（检测基础信息是否完善）
func (bc *IdentityController) CheckUpgradeQualification(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	userBasicData, _ := bc.UserService.QueryUserBasicInfo(username, c.Request.Host)
	userDataByte, _ := json.Marshal(&userBasicData)
	userDataMap := make(map[string]interface{})
	//fmt.Println(userBasicData)
	_ = json.Unmarshal(userDataByte, &userDataMap)
	//fmt.Println(userDataMap)
	for _, v := range userDataMap {
		if v == "" || v == 0 {
			c.JSON(http.StatusOK, gin.H{
				"qualification": false,
				"msg":           "基础信息未完善",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"qualification": true,
		"msg":           "身份升级资格核验成功",
	})
}
