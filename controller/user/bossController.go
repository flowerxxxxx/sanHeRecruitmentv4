package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanHeRecruitment/biz/nsqBiz"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/BindModel/userBind"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/messageUtil"
	"sanHeRecruitment/util/saveUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"sanHeRecruitment/util/tokenUtil"
	"sanHeRecruitment/util/uploadUtil"
	"strconv"
	"strings"
	"time"
)

type BossController struct {
	*mysqlService.ArticleService
	*mysqlService.JobService
	*mysqlService.UserService
	*mysqlService.DeliveryService
	*mysqlService.EducationService
	*mysqlService.CountService
	*mysqlService.CompanyService
	*mysqlService.InvitationService
	*mysqlService.LabelService
}

func BossControllerRouter(router *gin.RouterGroup) {
	//boc := BossController{}

}

// 2022.7.16完善消息推送线程后，每当检测收到消息且处于消息页面才去请求getMsgList；为保证线程正常，每30-60s请求一下getMsgList

//待处理
//TODO boss/service 在注销企业身份后 自动删除其发布的全部 招聘/需求

func BossControllerRouterToken(router *gin.RouterGroup) {
	boc := BossController{}
	//隐藏/取消隐藏招聘 （hide/appear）
	router.POST("/changeArtStatus", boc.BossChangeArtStatus)
	//boss查看自己发布的招聘需求 目前为完全审核成功状态
	//all:全部  pass:审核通过  hide：隐藏  fail：未通过 wait:待审核
	router.GET("/bossGetPublish/:queryType/:type/:pageNum", boc.BossGetPublish)
	//boss删除自己发布的招聘/需求
	router.POST("/deletePubArt", boc.BossDeleteArt)
	//boss发布招聘 修改status和show
	router.POST("/pubEmployReq", boc.PubEmployReq)
	//boss或本人查看简历-教育经历
	router.GET("/queryResumeEdu/:desUsername", boc.QueryResumeEdu)
	//boss或本人查看简历-基础信息
	router.GET("/queryResumeBasic/:desUsername", boc.QueryResumeBasic)
	//邀请投递   邀请投递d-选择邀请到自己的招聘，发起沟通
	router.POST("/inviteDelivery", boc.InviteDelivery)
	//boss查看人力资源
	router.POST("/queryEmployees", boc.BossQueryEmployees)
	//boss端审核对简历操作 通过:pass 不通过:fail
	router.POST("/resumeManage", boc.ResumeGetManage)
	//boss/service 查看所属公司基本信息
	router.GET("/getCompanyInfo", boc.BossGetCompanyInfo)
	//boss/service 修改公司信息
	router.POST("/updateCompanyInfo", boc.UpdateCompanyInfo)
	// boss/service 修改公司头像
	router.POST("/updateCompanyHeadPic", boc.UpdateCompanyHeadPic)
	// boss 查看简历投递信息
	router.GET("/resumeDeliveries/:type/:pageNum", boc.QueryResumeDeliveries)
	// boss 查看简历设置已读
	router.POST("/bossReadResume", boc.BossReadResume)
	//上传照片到服务器
	router.POST("/uploadComHead", boc.SavePicVouchersCom)
}

// SavePicVouchersCom 上传照片凭证（公司照片，凭证）
func (boc *BossController) SavePicVouchersCom(c *gin.Context) {
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
	fileUrl, fileAddr := uploadUtil.SaveFormat(fileFormat, c.Request.Host)
	//if err := c.SaveUploadedFile(file, fileAddr); err != nil {
	//	c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
	//	return
	//}
	if err := saveUtil.SaveCompressCutImg(file, fileAddr); err != nil {
		c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		return
	}
	controller.SuccessResp(c, "照片凭证上传成功", formatUtil.GetPicHeaderBody(c.Request.Host, fileUrl))
}

// BossReadResume boss 查看简历设置已读
func (boc *BossController) BossReadResume(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	deliverId := recJson["deliverId"].(string)
	deliverIdInt, err := strconv.Atoi(deliverId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	bossInfo := tokenUtil.GetUserClaims(c)
	err = boc.DeliveryService.SetDeliveryR1ead(deliverIdInt, bossInfo.User.Id)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "信息不存在")
			return
		} else {
			controller.ErrorResp(c, 215, "服务器错误")
			log.Println("BossReadResume", err, "recJson:", recJson)
			return
		}
	}
	controller.SuccessResp(c, "已设置为已读")
	return
}

// QueryResumeDeliveries 查看简历投递管理 all-1  pass 1  fail 2   wait 0
func (boc *BossController) QueryResumeDeliveries(c *gin.Context) {
	bossInfo := tokenUtil.GetUserClaims(c)
	bossId := bossInfo.User.Id
	qualification := -1
	queryType := c.Param("type")
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	switch queryType {
	case "pass":
		qualification = 1
	case "fail":
		qualification = 2
	case "wait":
		qualification = 0
	case "all":
		qualification = -1
	default:
		controller.ErrorResp(c, 202, "类型参数错误")
		return
	}
	bossQuery, _ := boc.DeliveryService.BossQueryDeliveries(bossId, qualification, pageNumInt)
	bossQueryTP := boc.CountService.GetQueryDeliveriesTotalPage(bossId, qualification)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "查询成功",
		"data":      bossQuery,
		"totalPage": bossQueryTP,
	})
	return
}

// UpdateCompanyHeadPic boss/service 修改公司头像
func (boc *BossController) UpdateCompanyHeadPic(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	username := tokenUtil.GetUsernameByToken(c)
	pic_url := recJson["pic_url"].(string)
	userInfo, err := boc.UserService.GetUserInfo(username, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 201, "未查询到用户信息")
		log.Println("UpdateCompanyHeadPic GetUserInfo", err, "recJson:", recJson)
		return
	}
	companyInfo, err := boc.CompanyService.QueryCompanyInfoById(userInfo.CompanyID, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 214, "服务器错误")
		log.Println("UpdateCompanyHeadPic", err, "recJson:", recJson)
		return
	}
	//异步删除上版图片
	go func() {
		_ = saveUtil.DeletePicSaver(companyInfo.PicUrl)
	}()
	saveFlag := strings.Index(pic_url, "uploadPic/")
	if saveFlag == -1 {
		controller.ErrorResp(c, 201, "上传路径错误")
		return
	}
	pic_url = pic_url[saveFlag:]
	err = boc.CompanyService.UpdateCompanyHeadPic(companyInfo.ComId, pic_url, username)
	if err != nil {
		controller.ErrorResp(c, 215, "服务器错误")
		log.Println("UpdateCompanyHeadPic", err, "recJson:", recJson)
		return
	}
	controller.SuccessResp(c, "修改成功")
}

// UpdateCompanyInfo boss/service 修改公司信息
func (boc *BossController) UpdateCompanyInfo(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	scaleTag := recJson["scale_tag"].(string)
	personScale := recJson["person_scale"].(string)
	address := recJson["address"].(string)
	description := recJson["description"].(string)
	phone := recJson["phone"].(string)
	username := tokenUtil.GetUsernameByToken(c)
	userInfo, err := boc.UserService.GetUserInfo(username, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 213, "服务器错误")
		log.Println("UpdateCompanyInfo", err, "recJson:", recJson)
		return
	}
	companyInfo, err := boc.CompanyService.QueryCompanyInfoById(userInfo.CompanyID, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 214, "服务器错误")
		log.Println("UpdateCompanyInfo", err, "recJson:", recJson)
		return
	}
	err = boc.CompanyService.UpdateCompanyInfo(companyInfo.ComId,
		scaleTag, personScale, address, username, description, phone)
	if err != nil {
		controller.ErrorResp(c, 215, "服务器错误")
		log.Println("UpdateCompanyInfo", err, "recJson:", recJson)
		return
	}
	controller.SuccessResp(c, "修改成功")
}

// BossGetCompanyInfo boss获取公司信息
func (boc *BossController) BossGetCompanyInfo(c *gin.Context) {
	bossUsername := tokenUtil.GetUsernameByToken(c)
	BossInfo, err := boc.UserService.GetUserInfo(bossUsername, c.Request.Host)
	companyInfo, err := boc.CompanyService.QueryCompanyInfoById(BossInfo.CompanyID, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 215, "服务器错误")
		log.Println("BossGetCompanyInfo", err)
		return
	}
	controller.SuccessResp(c, "查询成功", companyInfo)
}

// QueryResumeBasic boss或本人查看简历-基础信息
func (boc *BossController) QueryResumeBasic(c *gin.Context) {
	desUsername := c.Param("desUsername")
	perResumeInfo, err := boc.UserService.QueryUserResumeInfo(desUsername, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 201, "无匹配结果")
		return
	}
	controller.SuccessResp(c, "简历-除教育查询成功", perResumeInfo)
}

// QueryResumeEdu boss或本人查看简历-教育经历
func (boc *BossController) QueryResumeEdu(c *gin.Context) {
	desUsername := c.Param("desUsername")
	perEduInfo := boc.EducationService.QueryPersonalEdu(desUsername)
	controller.SuccessResp(c, "简历-教育查询成功", perEduInfo)
}

// InviteDelivery 邀请投递
func (boc *BossController) InviteDelivery(c *gin.Context) {
	bossUsername := tokenUtil.GetUsernameByToken(c)
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	artId := recJson["artId"].(string)
	desUsername := recJson["desUsername"].(string)
	inviteTime := timeUtil.GetNowTimeFormat()
	artIdInt, err := strconv.Atoi(artId)
	if err != nil {
		controller.ErrorResp(c, 201, "操作失败，[招聘/需求]id错误")
		return
	}
	invInfo, err := boc.InvitationService.QueryOneInvitation(artIdInt, desUsername)
	if err == nil &&
		invInfo.InviteTime.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		controller.ErrorResp(c, 202, "您今天已邀请过")
		return
	}
	err = boc.InvitationService.AddInvitationInfo(artIdInt, bossUsername, desUsername, inviteTime)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		log.Println("InviteDelivery", err, "recJson:", recJson)
		return
	}
	//消息通知协程
	go func() {
		artInfo, errPush := boc.ArticleService.QueryArtByID(artIdInt)
		deliveryUserBasicInfo, errPush := boc.UserService.QueryUserBasicInfo(desUsername, c.Request.Host)
		bossInfo, errPush := boc.UserService.GetUserInfo(bossUsername, c.Request.Host)
		bossCompanyInfo, errPush := boc.CompanyService.QueryCompanyInfoById(bossInfo.CompanyID, c.Request.Host)
		if errPush != nil {
			log.Println("InviteDelivery msg push failed,err:", errPush, "\nrequest body:", c.Request.Body)
			return
		}
		DeliveryManageTem := messageUtil.InviteDeliveryTem(
			deliveryUserBasicInfo.Name,
			deliveryUserBasicInfo.Gender,
			bossCompanyInfo.CompanyName,
			bossInfo.Name,
			artInfo.Title,
		)
		//websocketBiz.SysMsgPusher(desUsername, DeliveryManageTem)
		nsqBiz.ToServiceProducer(websocketModel.ToServiceMiddle{
			ToUsername: desUsername, MsgContent: DeliveryManageTem,
		})
	}()
	controller.SuccessResp(c, "邀请投递操作成功")
}

// PubEmployReq boss发布招聘
func (boc *BossController) PubEmployReq(c *gin.Context) {
	//recJson := map[string]interface{}{}
	//_ = c.BindJSON(&recJson)
	//fmt.Println(recJson)
	//careerJobId := recJson["careerJobId"].(string)
	//title := recJson["title"].(string)
	//content := recJson["content"].(string)
	//region := recJson["region"].(string)
	//salaryMin := recJson["salaryMin"].(string)
	//salaryMax := recJson["salaryMax"].(string)
	var bossPubBinder userBind.PubEmployReqBind
	err := c.ShouldBind(&bossPubBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败111")
		return
	}
	if !sqlUtil.JudgeFormatLegal(bossPubBinder.TagList) {
		controller.ErrorResp(c, 202, "不支持英文逗号")
		return
	}
	username := tokenUtil.GetUsernameByToken(c)
	userInfo, _ := boc.UserService.GetUserInfo(username, c.Request.Host)
	companyInfo, _ := boc.CompanyService.QueryCompanyInfoById(userInfo.CompanyID, c.Request.Host)
	tagStr := sqlUtil.SliToSqlString(bossPubBinder.TagList)
	salaryMinInt, err1 := strconv.Atoi(bossPubBinder.SalaryMin)
	salaryMaxInt, err2 := strconv.Atoi(bossPubBinder.SalaryMax)
	careerJobIdInt, err3 := strconv.Atoi(bossPubBinder.CareerJobId)
	if err1 != nil || err2 != nil || err3 != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	jobLabel, err := boc.LabelService.QueryLabelById(careerJobIdInt)
	if err != nil {
		controller.ErrorResp(c, 202, "无该标签id对应的信息")
		return
	}
	nowTime := timeUtil.GetNowTimeFormat()
	err = boc.ArticleService.AddNewEmployeeReq(careerJobIdInt, userInfo.User_id, userInfo.CompanyID,
		salaryMinInt, salaryMaxInt, companyInfo.Vip, bossPubBinder.Title,
		bossPubBinder.Content, bossPubBinder.Region, jobLabel.Label, tagStr, jobLabel.Type, nowTime)
	if err != nil {
		controller.ErrorResp(c, 215, "发布失败，服务器错误")
		log.Println("PubEmployReq failed,err:", err, "\n request body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "发布成功")
	return
}

// ResumeGetManage boss端审核通过/不通过简历   通过:pass 不通过:fail
func (boc *BossController) ResumeGetManage(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	deliverId := recJson["deliverId"].(string)
	action := recJson["action"].(string)
	deliverIdInt, _ := strconv.Atoi(deliverId)
	deliveryInfo, err := boc.DeliveryService.QueryDeliveryById(deliverIdInt)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		log.Println("ResumeGetManage", err, "recJson:", recJson)
		return
	}
	//if deliveryInfo.Qualification != 0 {
	//	controller.ErrorResp(c, 202, "该简历已做出操作")
	//	return
	//}
	qualification := 0
	if action == "pass" {
		qualification = 1
	} else if action == "fail" {
		qualification = 2
	}
	err = boc.DeliveryService.ModifyDeliveryQualification(deliverIdInt, qualification)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		log.Println("ResumeGetManage", err, "recJson:", recJson)
		return
	}
	go func() {
		artInfo, _ := boc.ArticleService.QueryArtByID(deliveryInfo.ArtId)
		deliveryUserBasicInfo, _ := boc.UserService.QueryUserBasicInfo(deliveryInfo.FromUsername, c.Request.Host)
		DeliveryManageTem := messageUtil.ResumeManageTem(
			deliveryUserBasicInfo.Name,
			deliveryUserBasicInfo.Gender,
			artInfo.Title,
			deliveryInfo.DeliveryTime,
			qualification,
		)
		//websocketBiz.SysMsgPusher(deliveryInfo.FromUsername, DeliveryManageTem)
		nsqBiz.ToServiceProducer(websocketModel.ToServiceMiddle{
			ToUsername: deliveryInfo.FromUsername, MsgContent: DeliveryManageTem,
		})
	}()
	controller.SuccessResp(c, action+"操作成功")
	return
}

// BossQueryEmployees boss查看人力资源
func (boc *BossController) BossQueryEmployees(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	jobLabel := recJson["jobLabel"].(string)
	gender := recJson["gender"].(string)
	minAge := recJson["minAge"].(string)
	maxAge := recJson["maxAge"].(string)
	minDegree := recJson["minDegree"].(string)
	pageNum := recJson["pageNum"].(string)
	pageNumInt, _ := strconv.Atoi(pageNum)
	minDegreeLevelStr := ""
	if minDegree != "" {
		minDegreeLevel := mysqlModel.DegreeWeight[minDegree]
		minDegreeLevelStr = strconv.Itoa(minDegreeLevel)
	}
	usernameS := boc.EducationService.ScreeningResumes(minAge, maxAge, gender, minDegreeLevelStr, jobLabel, pageNumInt)
	if len(usernameS) == 0 {
		controller.ErrorResp(c, 201, "暂无人力资源")
		return
	}
	userEduBasic, err := boc.EducationService.QueryEmployeeEduBasic(usernameS)
	if err != nil {
		controller.ErrorResp(c, 202, "无匹配结果")

		return
	}
	totalPage := boc.CountService.GetEmployeesTotalPage(minAge, maxAge, gender, minDegreeLevelStr, jobLabel)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "人力资源查询成功",
		"data":      userEduBasic,
		"totalPage": totalPage,
	})
	return
}

// BossChangeArtStatus   隐藏/取消隐藏文章  hide 3 /appear 1
func (boc *BossController) BossChangeArtStatus(c *gin.Context) {
	claims := tokenUtil.GetUserClaims(c)
	bossId := claims.User.Id
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	ArtId := recJson["artId"].(string)
	action := recJson["action"].(string)
	ArtIdInt, _ := strconv.Atoi(ArtId)
	err := boc.ArticleService.BossChangeArtShowStatus(bossId, ArtIdInt, action)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 201, "文章状态修改失败，身份不符或文章不存在")
			return
		} else {
			controller.ErrorResp(c, 215, "文章状态修改失败，服务器错误")
			log.Println("BossChangeArtStatus", err, "recJson:", recJson)
			return
		}
	}
	controller.SuccessResp(c, "文章状态修改成功")
}

// BossGetPublish BossGetPubArts boss查看自己发布的招聘需求 all:全部  pass:审核通过  hide：隐藏  fail：未通过 wait:待审核 TODO 招聘/其他
func (boc *BossController) BossGetPublish(c *gin.Context) {
	Type := c.Param("type")
	queryType := c.Param("queryType")
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	desStatus := -10
	desShow := -1
	switch Type {
	case "pass":
		desStatus = 1
	case "all":
		desStatus = -1
	case "fail":
		desStatus = 2
	case "wait":
		desStatus = 0
	case "hide":
		desShow = 0
		desStatus = -1
	}
	if desStatus == -10 {
		controller.ErrorResp(c, 201, "获取失败，无法匹配目标类型")
		return
	}
	claims := tokenUtil.GetUserClaims(c)
	BossId := claims.User.Id
	BossIdStr := strconv.Itoa(BossId)
	pubInfo := boc.JobService.BossGetJobInfoByBossId(BossIdStr, queryType, c.Request.Host, desStatus, desShow, pageNumInt)
	totalPage := boc.CountService.BossJobInfoTP(BossIdStr, queryType, desStatus, desShow)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   pubInfo,
		"msg":    "boss信息获取成功",
		"total":  totalPage,
	})
}

// BossDeleteArt boss删除自己发布的招聘/需求
func (boc *BossController) BossDeleteArt(c *gin.Context) {
	claims := tokenUtil.GetUserClaims(c)
	bossId := claims.User.Id
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	ArtId := recJson["artId"].(string)
	ArtIdInt, _ := strconv.Atoi(ArtId)
	err := boc.ArticleService.BossDeletePubArt(bossId, ArtIdInt)
	if err != nil {
		controller.ErrorResp(c, 201, "删除失败，身份不符")
		return
	}
	controller.SuccessResp(c, "删除成功")
}
