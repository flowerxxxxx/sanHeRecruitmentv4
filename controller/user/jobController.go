package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/library/lruEngine"
	"sanHeRecruitment/models/BindModel/userBind"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/module/controllerModule"
	"sanHeRecruitment/module/recommendModule"
	"sanHeRecruitment/module/websocketModule"
	"sanHeRecruitment/service/esService"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/util"
	"sanHeRecruitment/util/copyUtil"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/messageUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"sanHeRecruitment/util/tokenUtil"
	"strconv"
	"time"
)

type JobController struct {
	*mysqlService.JobService
	*mysqlService.LabelService
	*mysqlService.CountService
	*mysqlService.CollectionService
	*mysqlService.ArticleService
	*mysqlService.DeliveryService
	*mysqlService.UserService
	*mysqlService.EducationService
	*mysqlService.InvitationService
	*mysqlService.CompanyService
	*mysqlService.DailySaverService
	*mysqlService.DockService
	controllerModule.JobConModule
	*esService.ArticleESservice
}

func JobControllerRouter(router *gin.RouterGroup) {
	j := JobController{}
	//获取工作信息小程序
	router.POST("/getJobInfo", j.GetJobInfos)
	//获取工作信息-web
	//router.POST("/getJobInfoWeb", j.GetJobInfosWeb)
	//获取标签
	router.GET("/getLabel/:labelType", j.GetLabel)
	//获取默认在父级带全部项的标签
	router.GET("/getLabelAll/:labelType", j.GetLabelAll)
	//获取职位详细页面信息
	router.GET("/getRecruitInfo/:art_id", j.GetRecruitInfo)
	//获取boss发布的招聘信息
	router.GET("/bossJobInfo/:labelType/:bossId/:pageNum", j.GetBossJobInfo)
	//获取公司有的招聘信息
	router.POST("/GetCompanyPub", j.GetCompanyPub)
	//首页模糊检索职位 招聘和需求可单独做
	router.GET("/fuzzyQueryAll/:fuzzyName", j.FuzzyQuery)
	//首页模糊检索公司
	router.POST("/fuzzyQueryCompanies", j.FuzzyQueryCompanies)
	//获取公司已经包含的工作/需求标签
	router.GET("/getCompanyPubLabel/:companyId/:type", j.QueryCompanyPubLabel)
	//获取公司信息
	router.GET("/getCompanyInfo/:companyName", j.QueryCompanyInfo)
	//模糊获取工作信息
	router.POST("/FuzzyQueryJobInfos", j.FuzzyQueryJobInfos)
}

func JobControllerRouterToken(router *gin.RouterGroup) {
	j := JobController{}
	router.POST("/checkResumeQualification", j.checkResumeQualification)
	router.GET("/getCollectStatus/:art_id", j.GetCollectStatus)
	router.POST("/deliverResume", j.DeliverResume) //投递简历
	//推荐招聘岗位
	router.POST("/getRecommendJobInfos", j.GetRecommendJobInfos)
	router.GET("/JobSearchProgress/:status/:pageNum", j.QueryJobSearchProgress)
	//router.GET("/GetInviterInfo/:pageNum", j.GetInviterInfo)
	//user 查看邀请投递 distinct去重以最新为准\n\t
	router.GET("/getInvitedJobs/:pageNum", j.QueryInvitedJobs)
	//删除邀请投递记录
	router.POST("/deleteInvitation", j.DeleteInvitation)
	//企业用户与普通用户对接记录，在求职页面与需求页面发起沟通前使用
	router.POST("/recordDockInfo", j.RecordDockInfo)
}

// FuzzyQueryCompanies 模糊
func (jc *JobController) FuzzyQueryCompanies(c *gin.Context) {
	var fBinder userBind.FuzzyQueryComs
	err := c.ShouldBind(&fBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	fuzzyCompanyList := jc.CompanyService.FuzzyQueryCompaniesPage(fBinder.FuzzyName, "1", 1, fBinder.PageNum)
	totalPage := jc.CountService.QueryAllFuzzyCompaniesTP(fBinder.FuzzyName, "1", 1)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "公司模糊查询成功",
		"data":   fuzzyCompanyList,
		"total":  totalPage,
	})
}

// FuzzyQueryJobInfos 模糊获取工作信息
func (jc *JobController) FuzzyQueryJobInfos(c *gin.Context) {
	var fBinder userBind.FuzzyQueryJobs
	err := c.ShouldBind(&fBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	fuzzyJobInfo, errEs := jc.ArticleESservice.FuzzyArticlesQuery(fBinder.PageNum, fBinder.FuzzyName, fBinder.QueryType)
	if errEs != nil {
		log.Println("FuzzyQueryJobInfos ArticleESservice.FuzzyArticlesQuery failed,err", errEs)
	}
	if len(fuzzyJobInfo) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data":   []mysqlModel.UserComArticle{},
			//"totalPage": TotalPageNum,
			"msg": "模糊招聘信息获取成功",
		})
		return
	}
	UCAs := []mysqlModel.UserComArticle{}
	uids := []int{}
	cids := []int{}
	for i, ul := 0, len(fuzzyJobInfo); i < ul; i++ {
		uids = append(uids, fuzzyJobInfo[i].BossId)
		cids = append(cids, fuzzyJobInfo[i].CompanyId)
	}
	//var uInfos []mysqlModel.UserNH
	//var cInfos []mysqlModel.CompanyLite
	//if len(uids) == 1 {
	//	uInfos = jc.UserService.QueryUserLite(uids[0])
	//}else {
	uInfos := jc.JobService.FuzzyQueryBaseUser(uids)
	//}
	//if len(cids) == 1 {
	//	cInfos = jc.CompanyService.QueryBaseCom(cids[0])
	//}else {
	cInfos := jc.JobService.FuzzyQueryBaseCom(cids)
	//
	//}
	for i, ul := 0, len(fuzzyJobInfo); i < ul; i++ {
		uca := mysqlModel.UserComArticle{}
		uca.ArtID = fuzzyJobInfo[i].ArtId
		errCF := copyUtil.CopyFields(&uca, fuzzyJobInfo[i])

		if errCF != nil {
			log.Println("FuzzyQueryJobInfos CopyFields err,", errCF)
		}

		for _, item := range uInfos {
			if item.User_id == fuzzyJobInfo[i].BossId {
				uca.Nickname = item.Nickname
				uca.HeadPic = item.Head_pic
			}
		}
		for _, item := range cInfos {
			if item.ComId == fuzzyJobInfo[i].CompanyId {
				uca.CompanyName = item.CompanyName
				uca.PersonScale = item.PersonScale
				uca.ComLevel = item.ComLevel
			}
		}
		if fBinder.QueryType == "request" {
			uca.CompanyName = fmt.Sprintf("%v*****", string([]rune(uca.CompanyName)[:1]))
		}
		uca.TagsOut = sqlUtil.SqlStringToSli(uca.Tags)
		uca.HeadPic = formatUtil.GetPicHeaderBody(c.Request.Host, uca.HeadPic)
		UCAs = append(UCAs, uca)
	}
	//fuzzyJobInfo := jc.JobService.FuzzyQueryJobs(fBinder.FuzzyName, fBinder.QueryType, c.Request.Host, fBinder.PageNum, 0)
	//TotalPageNum := jc.CountService.GetFuzzyQueryJobsTP(fBinder.FuzzyName, fBinder.QueryType, 0)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   UCAs,
		//"totalPage": TotalPageNum,
		"msg": "模糊招聘信息获取成功",
	})
}

// RecordDockInfo 企业用户与普通用户对接记录
func (jc *JobController) RecordDockInfo(c *gin.Context) {
	userInfo := tokenUtil.GetUserClaims(c)
	userId := userInfo.User.Id
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	artId := recJson["art_id"].(string)
	bossId := recJson["boss_id"].(string)
	comId := recJson["com_id"].(string)
	aiInt, err := strconv.Atoi(artId)
	biInt, err2 := strconv.Atoi(bossId)
	ciInt, err3 := strconv.Atoi(comId)
	if err != nil || err2 != nil || err3 != nil {
		controller.ErrorResp(c, 201, "操作失败，参数错误")
		return
	}
	userInf, err := jc.UserService.QueryBasicUserInfoByUserId(userId)
	bossInf, err2 := jc.UserService.QueryBasicUserInfoByUserId(biInt)
	artInfo, err3 := jc.ArticleService.QueryArtByID(aiInt)
	if err != nil || err2 != nil || err3 != nil {
		controller.ErrorResp(c, 211, "RecordDockFailed，服务器错误")
		log.Println("RecordDockInfo", err, err2, err3, "recJson:", recJson)
		return
	}
	err = jc.DockService.AddDockRecord(ciInt, biInt, userId, aiInt,
		&timeUtil.MyTime{Time: time.Now()}, bossInf.Name, userInf.Name, artInfo.Title)
	if err != nil {
		controller.ErrorResp(c, 212, "RecordDockFailed，服务器错误")
		log.Println("RecordDockInfo", err, "recJson:", recJson)
		return
	}
	controller.SuccessResp(c, "对接记录存储成功")
}

// QueryInvitedJobs 查看邀请投递
func (jc *JobController) QueryInvitedJobs(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页数参数错误")
		return
	}
	invInfos, _ := jc.InvitationService.QueryInvitationInfos(username, c.Request.Host, pageNumInt)
	totalPage := jc.CountService.GetInvitedJobsTotalPage(username)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "查询成功",
		"data":      invInfos,
		"totalPage": totalPage,
	})
	return

}

// DeleteInvitation 删除邀请投递记录
func (jc *JobController) DeleteInvitation(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	invitationId := recJson["invitationId"].(string)
	fmt.Println(invitationId)
	username := tokenUtil.GetUsernameByToken(c)
	err := jc.InvitationService.UserDeleteInfo(username, invitationId)
	if err != nil {
		controller.ErrorResp(c, 201, "id验证错误")
		return
	}
	controller.SuccessResp(c, "删除成功")
}

// FuzzyQuery 模糊查询
func (jc *JobController) FuzzyQuery(c *gin.Context) {
	queryFlag := c.Param("fuzzyName")
	fuzzyCompanyList := jc.CompanyService.FuzzyQueryCompanies(queryFlag, "1", 1)
	fuzzyJobList, _ := jc.JobService.FuzzyQueryJobOrReq(queryFlag, "job")
	fuzzyReqList, _ := jc.JobService.FuzzyQueryJobOrReq(queryFlag, "request")
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"msg":     "模糊查询成功",
		"company": fuzzyCompanyList,
		"job":     fuzzyJobList,
		"request": fuzzyReqList,
	})
	return
}

func (jc *JobController) QueryCompanyInfo(c *gin.Context) {
	companyName := c.Param("companyName")
	companyBasicInfo, err := jc.CompanyService.QueryCompanyBasicInfoByName(companyName, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 201, "该公司不存在")
		return
	}
	if companyBasicInfo.ComLevel == 2 {
		controller.ErrorResp(c, 202, "禁止访问该信息")
		return
	}
	controller.SuccessResp(c, "公司信息查找成功", companyBasicInfo)
}

func (jc *JobController) QueryCompanyPubLabel(c *gin.Context) {
	companyId := c.Param("companyId")
	jobLabel := c.Param("type")
	companyIdInt, err := strconv.Atoi(companyId)
	if err != nil {
		controller.ErrorResp(c, 201, "id 格式错误 ")
		return
	}
	companyLabel := jc.LabelService.QueryCompanyLabel(jobLabel, companyIdInt, 1)
	controller.SuccessResp(c, "["+jobLabel+"]标签查询成功", companyLabel)
}

// GetCompanyPub 获取改公司所属的职位或者需求 全部职业或者需求为 ""
func (jc *JobController) GetCompanyPub(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	companyId := recJson["companyId"].(string)
	JobReqType := recJson["jobReqType"].(string)
	pageNum := recJson["pageNum"].(string)
	getType := recJson["type"].(string)
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	companyIdInt, err := strconv.Atoi(companyId)
	if err != nil {
		controller.ErrorResp(c, 201, "公司id参数错误")
		return
	}
	comArtInfo := jc.JobService.GetCompanyJobInfo(JobReqType, getType, c.Request.Host, companyIdInt, pageNumInt, 1)
	Total := jc.CountService.CompanyPubTotal(JobReqType, getType, companyIdInt, 1)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "查询成功",
		"data":   comArtInfo,
		"total":  Total,
	})
}

// QueryJobSearchProgress 查询求职进展 all，read，pass，fail
func (jc *JobController) QueryJobSearchProgress(c *gin.Context) {
	username := tokenUtil.GetUsernameByToken(c)
	queryType := c.Param("status")
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	if queryType != "all" && queryType != "read" && queryType != "pass" &&
		queryType != "fail" {
		controller.ErrorResp(c, 201, "请求所属的参数错误")
		return
	}
	read := -1
	qualification := -1
	switch queryType {
	case "read":
		read = 1
	case "pass":
		qualification = 1
	case "fail":
		qualification = 2
	}
	if read == 1 {
		deliveryInfo, err := jc.DeliveryService.QueryDelivery(username, c.Request.Host, 0, read, pageNumInt)
		if err != nil {
			controller.ErrorResp(c, 213, "服务器错误")
			log.Println("RecordDockInfo", err)
			return
		}
		totalPageNum := jc.CountService.GetDeliveryTotalPage(username, 0, read)
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "查询成功",
			"data":   deliveryInfo,
			"total":  totalPageNum,
		})
	} else if qualification == -1 {
		deliveryInfo, err := jc.DeliveryService.QueryAllDelivery(username, c.Request.Host, pageNumInt)
		if err != nil {
			controller.ErrorResp(c, 214, "服务器错误")
			log.Println("RecordDockInfo", err)
			return
		}
		totalPageNum := jc.CountService.GetAllDeliveryTotalPage(username)
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "查询成功",
			"data":   deliveryInfo,
			"total":  totalPageNum,
		})
	} else {
		deliveryInfo, err := jc.DeliveryService.QueryDelivery(username, c.Request.Host, qualification, 1, pageNumInt)
		if err != nil {
			controller.ErrorResp(c, 215, "服务器错误")
			log.Println("RecordDockInfo", err)
			return
		}
		totalPageNum := jc.CountService.GetDeliveryTotalPage(username, qualification, 1)
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    "查询成功",
			"data":   deliveryInfo,
			"total":  totalPageNum,
		})
	}
	return
}

// FuzzyQueryJobs 模糊检索工作
func (jc *JobController) FuzzyQueryJobs(c *gin.Context) {
	//fuzzyJobName := c.Param("fuzzyName")

}

// GetLabel 获取工作标签 (多处通用) job city
func (jc *JobController) GetLabel(c *gin.Context) {
	labelType := c.Param("labelType")
	data := jc.LabelService.GetLabelTree(labelType, 1)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   data,
		"msg":    "标签获取成功",
	})
}

// GetLabelAll 获取工作标签 (多处通用) job city
func (jc *JobController) GetLabelAll(c *gin.Context) {
	labelType := c.Param("labelType")
	data := jc.LabelService.GetLabelTree(labelType, 1)
	allFlagLabel := mysqlModel.LabelOut{
		Label: mysqlModel.Label{
			ID:       0,
			Level:    1,
			ParentId: 0,
			Label:    "全部",
			Type:     labelType,
		},
		Value:    0,
		Children: []interface{}{},
	}
	newData := util.SliAddFromHead(data, allFlagLabel)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   newData,
		"msg":    "标签获取成功",
	})
}

// GetRecommendJobInfos 获取推荐招聘信息（计算权重后排序）
func (jc *JobController) GetRecommendJobInfos(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	username := tokenUtil.GetUsernameByToken(c)
	userInfo, _ := jc.UserService.GetUserInfo(username, c.Request.Host)
	if userInfo.Intended_position == "" {
		controller.ErrorResp(c, 201, "请完善个人信息以激活智能推荐")
		return
	}
	labelInfo := jc.LabelService.QueryLabelByContent(userInfo.Intended_position)
	careerJobId := labelInfo.ID
	if careerJobId == 0 {
		controller.ErrorResp(c, 201, "目标职位类别更新，请重新设置")
		return
	}
	careerJobIdStr := strconv.Itoa(careerJobId)
	pageNum := recJson["pageNum"].(string)
	pageNumInt, _ := strconv.Atoi(pageNum)
	jobInfo := jc.JobService.GetRecommendJobs(careerJobIdStr, c.Request.Host, pageNumInt)
	TotalPageNum := jc.CountService.GetJobsTotalPage(careerJobId, "", "job", "micro")
	if TotalPageNum == 0 {
		TotalPageNum = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      jobInfo,
		"totalPage": TotalPageNum,
		"msg":       "招聘推荐信息获取成功",
	})
}

// GetJobInfos 获取招聘信息（查找招聘+boss信息）sql为最新发布
func (jc *JobController) GetJobInfos(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	region := recJson["region"].(string)
	careerJobId := recJson["careerJobId"].(string)
	pageNum := recJson["pageNum"].(string)
	labelType := recJson["labelType"].(string)
	cjiInt, err := strconv.Atoi(careerJobId)
	if err != nil {
		controller.ErrorResp(c, 201, "工作参数错误")
		return
	}
	pageNumInt, _ := strconv.Atoi(pageNum)
	jobInfo := jc.JobService.GetJobs(region, labelType, c.Request.Host, cjiInt, pageNumInt, 10, 0)

	TotalPageNum := jc.CountService.GetJobsTotalPage(cjiInt, region, labelType, "micro")
	if TotalPageNum == 0 {
		TotalPageNum = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      jobInfo,
		"totalPage": TotalPageNum,
		"msg":       "招聘信息获取成功",
	})
}

// GetJobInfosWeb 获取招聘信息（查找招聘+boss信息）sql为最新发布
func (jc *JobController) GetJobInfosWeb(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	region := recJson["region"].(string)
	careerJobId := recJson["careerJobId"].(string)
	pageNum := recJson["pageNum"].(string)
	labelType := recJson["labelType"].(string)
	cjiInt, err := strconv.Atoi(careerJobId)
	if err != nil {
		controller.ErrorResp(c, 201, "工作参数错误")
		return
	}
	pageNumInt, _ := strconv.Atoi(pageNum)
	jobInfo := jc.JobService.GetJobs(region, labelType, c.Request.Host, cjiInt, pageNumInt, 15, 0)
	TotalPageNum := jc.CountService.GetJobsTotalPage(cjiInt, region, labelType, "web")
	if TotalPageNum == 0 {
		TotalPageNum = 1
	}
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      jobInfo,
		"totalPage": TotalPageNum,
		"msg":       "招聘信息获取成功",
	})
}

// GetJobInfos2 获取招聘信息（查找招聘+boss信息）sql为最新发布 地区可为
// {"careerJobId":"193","region":["燕郊镇","金桥镇"],"pageNum":"1"}
func (jc *JobController) GetJobInfos2(c *gin.Context) {
	json1 := map[string]interface{}{}
	_ = c.BindJSON(&json1)
	region := json1["region"]
	careerJobId := json1["careerJobId"].(string)
	pageNum := json1["pageNum"].(string)
	pageNumInt, _ := strconv.Atoi(pageNum)
	regionSLi := []string{}
	region2 := region.([]interface{})
	for _, v := range region2 {
		regionSLi = append(regionSLi, v.(string))
	}
	fmt.Println(regionSLi)
	jobInfo := jc.JobService.GetJobs2(careerJobId, regionSLi, pageNumInt)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   jobInfo,
		"msg":    "招聘信息获取成功",
	})
}

// CheckResumeIfDone 检测提交简历的资格（信息完善）
func (jc *JobController) checkResumeQualification(c *gin.Context) {
	userClaim := tokenUtil.GetUserClaims(c)
	username := userClaim.User.UserName
	userId := userClaim.User.Id
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	bossId := recJson["boss_id"].(string)
	bossIdInt, _ := strconv.Atoi(bossId)
	if bossIdInt == userId {
		c.JSON(http.StatusOK, gin.H{
			"qualification": false,
			"msg":           "不能向自己投递",
		})
		return
	}
	//userInfo := jc.UserService.QueryUserLevel(username)
	//if userInfo.UserLevel != 0 || userInfo.IdentyPin != 0 {
	//	c.JSON(http.StatusOK, gin.H{
	//		"qualification": false,
	//		"msg":           "仅普通用户投递",
	//	})
	//	return
	//}
	userResumeData, _ := jc.UserService.QueryUserResumeInfo(username, c.Request.Host)
	userDataByte, _ := json.Marshal(&userResumeData)
	userDataMap := make(map[string]interface{})
	_ = json.Unmarshal(userDataByte, &userDataMap)
	for _, v := range userDataMap {
		if v == "" || v == 0 {
			c.JSON(http.StatusOK, gin.H{
				"qualification": false,
				"msg":           "简历信息未完善",
			})
			return
		}
	}
	perEduInfo := jc.EducationService.QueryPersonalEdu(username)
	if len(perEduInfo) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"qualification": false,
			"msg":           "简历信息未完善",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"qualification": true,
		"msg":           "投递资格核验成功",
	})
}

// GetRecruitInfo 获取详细的招聘信息
func (jc *JobController) GetRecruitInfo(c *gin.Context) {
	artId := c.Param("art_id")
	artIdInt, err := strconv.Atoi(artId)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数格式错误")
		return
	}

	artInfoB, ok := lruEngine.LruEngine.Get("RecruitInfo_" + artId)
	if ok {
		recInfo := mysqlModel.OneArticleOut{}
		errUmMar := json.Unmarshal(artInfoB.ByteSlice(), &recInfo)
		//fmt.Println("get from lru success")
		if errUmMar == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 200,
				"data":   recInfo,
				"msg":    "文章信息获取成功",
			})
			go func(artId string) {
				//ctx := context.Background()
				//valueCtx := context.WithValue(ctx,artId,artId)
				//fmt.Println(valueCtx.Value(artId))
				jc.ArticleService.AddArtView(artIdInt)
				errIn := jc.DailySaverService.AddDailyView(artIdInt)
				if errIn != nil {
					log.Println("[GoroutineErrLog]", errIn)
				}
				//重新处理权重
				recommendModule.DealArtRecommendWeight(artId)
			}(artId)
			return
		}
		log.Println("lru err")
	}

	artInfo, errGet := jc.JobConModule.GetRecruitInfoFromRedis(artId)
	if errGet != nil {
		artInfo, err = jc.JobService.GetOneArtInfo(artId, c.Request.Host)
		if err != nil {
			controller.ErrorResp(c, 202, "该需求不存在")
			return
		}
		jc.JobConModule.SaveRecruitInfoToRedis(artId, artInfo)

		recInfoByte, errMar := json.Marshal(artInfo)
		if errMar != nil {
			log.Println("SaveRecruitInfoToRedis Marshal failed,err:", errMar)
			return
		}
		lruEngine.LruEngine.Add("RecruitInfo_"+artId, lruEngine.ByteView{B: recInfoByte})
		//fmt.Println("lru add success")
	}

	go func(artId string) {
		//ctx := context.Background()
		//valueCtx := context.WithValue(ctx,artId,artId)
		//fmt.Println(valueCtx.Value(artId))
		jc.ArticleService.AddArtView(artIdInt)
		errIn := jc.DailySaverService.AddDailyView(artIdInt)
		if errIn != nil {
			log.Println("[GoroutineErrLog]", errIn)
		}
		//重新处理权重
		recommendModule.DealArtRecommendWeight(artId)
	}(artId)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   artInfo,
		"msg":    "文章信息获取成功",
	})
}

// GetCollectStatus 检测招聘需求的收藏状态
func (jc *JobController) GetCollectStatus(c *gin.Context) {
	artId := c.Param("art_id")
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	username := claims.User.UserName
	_, err := jc.CollectionService.QueryColRec(artId, username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"msg":    false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    true,
	})
}

// GetBossJobInfo 获取boss发布的招聘信息
func (jc *JobController) GetBossJobInfo(c *gin.Context) {
	bossId := c.Param("bossId")
	pageNum := c.Param("pageNum")
	labelType := c.Param("labelType")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	jobInfos := jc.JobService.GetJobInfoByBossId(bossId, labelType, c.Request.Host, 1, 1, pageNumInt)
	totalPage := jc.CountService.BossJobInfoTP(bossId, labelType, 1, 1)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"data":   jobInfos,
		"msg":    "boss信息获取成功",
		"total":  totalPage,
	})
}

// DeliverResume 投递简历
func (jc *JobController) DeliverResume(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	auth := c.Request.Header.Get("Authorization")
	claims, _ := tokenUtil.ParseToken(auth)
	fromUsername := claims.User.UserName
	artIdStr := recJson["artId"].(string)
	bossIdStr := recJson["bossId"].(string)
	bossId, _ := strconv.Atoi(bossIdStr)
	artId, _ := strconv.Atoi(artIdStr)
	deliveryTime := timeUtil.GetNowTimeFormat()
	artInfo, errArt := jc.JobService.GetOneArtInfo(artIdStr, c.Request.Host)
	if errArt != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 404,
			"msg":    "该需求已丢失",
		})
		return
	}
	err := jc.DeliveryService.AddDeliveryService(bossId, artId, fromUsername, deliveryTime)
	if err != nil {
		if err == mysqlService.HasFound {
			c.JSON(http.StatusOK, gin.H{
				"status": 201,
				"msg":    "您已投递过简历，无法再次投递",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": 215,
			"msg":    "简历投递失败，服务器错误",
		})
		log.Println("RecordDockInfo", err, "recJson:", recJson)
		return
	}
	go func(artIdStr string) {
		jc.ArticleService.AddDeliveryNum(artIdStr)
		_ = jc.DailySaverService.AddDailyDelivery(artId)
		recommendModule.DealArtRecommendWeight(artIdStr)
	}(artIdStr)
	//消息推送匿名函数
	go func() {
		bossInfo, _ := jc.UserService.QueryUserInfoByUserId(bossIdStr)
		bossUsername := bossInfo.Username
		fromUserNickname := jc.UserService.QueryUserNickByUsername(fromUsername)
		title := artInfo.Title
		//wechatPubAcc.DeliveryResumeMessagePush(bossUsername, fromUserNickname, title)

		DeliveryTem := messageUtil.DeliveryPushTem(
			bossInfo.Name,
			bossInfo.Gender,
			fromUserNickname,
			title,
		)
		websocketModule.SysMsgPusher(bossUsername, DeliveryTem)

	}()
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "简历投递成功",
	})
}
