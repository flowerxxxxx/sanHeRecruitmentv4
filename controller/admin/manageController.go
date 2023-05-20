package admin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanHeRecruitment/biz/websocketBiz"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/BindModel/adminBind"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/util/messageUtil"
	"sanHeRecruitment/util/timeUtil"
	"sanHeRecruitment/util/tokenUtil"
	"strconv"
	"strings"
	"time"
)

// ManageController 管理controller
type ManageController struct {
	*mysqlService.UserService
	*mysqlService.LabelService
	*mysqlService.UpgradeService
	*mysqlService.ArticleService
	*mysqlService.CountService
	*mysqlService.DailySaverService
	*mysqlService.VoucherService
	*mysqlService.MassSendService
	*mysqlService.PropagandaService
	*mysqlService.NoticeService
	*mysqlService.VipShowService
}

func ManageControllerRouter(router *gin.RouterGroup) {
	//mc := ManageController{}

}

// TODO 2022.7.19 部分功能模块需要推送到公众号的协程尚未完成
// 2022.7.19 数据记录 每个【招聘/需求】每日的访问量和投递量 字段 id，art_id，day_view,day_delivery,date[日]
// 2022.7.19 日投递量/观看量统计【完善】
// 参考b站权重 指定自己的权重方案 推荐人才和招聘1

func ManageControllerRouterToken(router *gin.RouterGroup) {
	mc := ManageController{}
	//TODO 隐藏和取消隐藏文章（是否要和公司同级别 status0待审核 status1审核通过，show 1
	//TODO 对admin和boss/service隐藏没做权级区分
	router.POST("/manageArtStatus", mc.ManageArtStatus)
	// 添加label标签
	router.POST("/addLabel", mc.AddLabel)
	// 删除label标签
	router.POST("/deleteLabel", mc.DeleteLabel)
	// 修改label标签
	router.POST("/editLabel", mc.EditLabel)
	// TODO 管理公司vip，boss申请的并非为本人的vip而是公司的vip
	//router.POST("/mangeVip?")
	//同意升级
	router.POST("/upgradeAdmit", mc.UpgradeAdmit)
	//不同意升级
	router.POST("/upgradeDisAdmit", mc.UpgradeDisAdmit)
	//群发消息
	router.POST("/sendColonyMsg", mc.SendColonyMsg)
	//获取群发消息历史
	router.GET("/getMassPubHistory/:pageNum", mc.GetMassPubHistory)
	//招聘/需求 审核通过 status = 1 TODO  （2）  boss根据vip来发布无需审核的需求
	router.POST("/applyPassSuccess", mc.ApplyPassSuccess)
	// 招聘/需求 审核不通过 status = 2
	router.POST("/applyPassFail", mc.ApplyPassFail)
	//获取需要审核发布的 all全部，其余为招聘或者需求的类别(如java)
	//router.GET("/getWaitingApply/:queryLabel/:labelLevel/:pageNum", mc.GetApplyPubArt)
	//获取需要进行升级审核的
	router.POST("/GetWaitingUpgrade", mc.GetWaitingUpgrade)
	//查看用户基本信息
	router.POST("/queryUserBasicInfo", mc.QueryUserBasicInfo)
	//查看用户的升级凭证内容
	router.POST("/getUserUpVouchers", mc.GetUserUpVouchers)
	//管理员对发布的 招聘/需求 进行高权限删除
	router.POST("/AdminDeletePubInfo", mc.AdminDeletePubInfo)
	//发布待审核批量许可
	router.POST("/batchPubAdmit", mc.BatchPubAdmit)
	//删除 公司/机构 人员的 企业/服务机构 身份
	router.POST("/DeleteComUser", mc.DeleteComUser)
	//置顶需求发布
	router.POST("/ChangeTopPubStatus", mc.TopPub)
	//置顶焦点
	router.POST("/ChangePropagandaStatus", mc.TopPropaganda)
	//置顶公告
	router.POST("/ChangeNoticeStatus", mc.TopNotice)
	//置顶会员风采
	router.POST("/ChangeVipShowStatus", mc.TopVipShow)
	//优先需求标签
	router.POST("/TopLabel", mc.TopLabel)
}

// TopVipShow 置顶会员风采
func (mc *ManageController) TopVipShow(c *gin.Context) {
	var TopBinder adminBind.TopPub
	err := c.ShouldBind(&TopBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	if TopBinder.DesStatus == 1 {
		maxCount, err := mc.VipShowService.QueryMaxRecommendCount()
		if err != nil {
			log.Println("TopVipShow QueryMaxRecommendCount failed,err:", TopBinder)
			controller.ErrorResp(c, 211, "无发布内容")
			return
		}
		errX := mc.VipShowService.ChangeVipShowStatus(TopBinder.Id, maxCount+1)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopVipShow ChangeTopNoticeStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 212, "服务器错误")
				return
			}
		}
	} else {
		errX := mc.VipShowService.ChangeVipShowStatus(TopBinder.Id, 0)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopVipShow ChangeTopNoticeStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 213, "服务器错误")
				return
			}
		}
	}
	controller.SuccessResp(c, "置顶状态修改成功")
	return
}

// TopNotice 置顶公告
func (mc *ManageController) TopNotice(c *gin.Context) {
	var TopBinder adminBind.TopPub
	err := c.ShouldBind(&TopBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	if TopBinder.DesStatus == 1 {
		maxCount, err := mc.NoticeService.QueryMaxRecommendCount()
		if err != nil {
			log.Println("TopNotice QueryMaxRecommendCount failed,err:", err, "\n reqInfo:", TopBinder)
			controller.ErrorResp(c, 211, "无发布内容")
			return
		}
		errX := mc.NoticeService.ChangeTopNoticeStatus(TopBinder.Id, maxCount+1)
		if errX != nil {
			if errX == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopNotice ChangeTopNoticeStatus failed,err:", errX, "\n reqInfo:", TopBinder)
				controller.ErrorResp(c, 212, "服务器错误")
				return
			}
		}
	} else {
		errX := mc.NoticeService.ChangeTopNoticeStatus(TopBinder.Id, 0)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopNotice ChangeTopNoticeStatus failed,err:", err, "\n reqInfo:", TopBinder)
				controller.ErrorResp(c, 213, "服务器错误")
				return
			}
		}
	}
	controller.SuccessResp(c, "置顶状态修改成功")
	return
}

// TopPropaganda 置顶焦点
func (mc *ManageController) TopPropaganda(c *gin.Context) {
	var TopBinder adminBind.TopPub
	err := c.ShouldBind(&TopBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	if TopBinder.DesStatus == 1 {
		maxCount, err := mc.PropagandaService.QueryMaxRecommendCount()
		if err != nil {
			log.Println("TopPropaganda QueryMaxRecommendCount failed,err:", TopBinder)
			controller.ErrorResp(c, 211, "无发布内容")
			return
		}
		errX := mc.PropagandaService.ChangeTopProStatus(TopBinder.Id, maxCount+1)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopPropaganda ChangeTopProStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 212, "服务器错误")
				return
			}
		}
	} else {
		errX := mc.PropagandaService.ChangeTopProStatus(TopBinder.Id, 0)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该内容相关信息")
				return
			} else {
				log.Println("TopPropaganda ChangeTopProStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 213, "服务器错误")
				return
			}
		}
	}
	dao.Redis.Del("PropagandaInfo")
	controller.SuccessResp(c, "置顶状态修改成功")
	return
}

// TopLabel 优先需求标签
func (mc *ManageController) TopLabel(c *gin.Context) {
	var TopBinder adminBind.TopLabel
	err := c.ShouldBind(&TopBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	maxCount, err := mc.LabelService.QueryMaxLabelCount()
	if err != nil {
		log.Println("TopPub QueryMaxLabelCount failed,err:", TopBinder)
		controller.ErrorResp(c, 211, "无发布内容")
		return
	}
	errX := mc.LabelService.ChangeTopPubStatus(TopBinder.Id, maxCount+1)
	if errX != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "无该标签")
			return
		} else {
			log.Println("TopPub ChangeTopPubStatus failed,err:", TopBinder, "err:", errX)
			controller.ErrorResp(c, 212, "服务器错误")
			return
		}
	}
	controller.SuccessResp(c, "优先设置成功")
	return
}

// TopPub 置顶发布
func (mc *ManageController) TopPub(c *gin.Context) {
	var TopBinder adminBind.TopPub
	err := c.ShouldBind(&TopBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	if TopBinder.DesStatus == 1 {
		maxCount, err := mc.ArticleService.QueryMaxRecommendCount()
		if err != nil {
			log.Println("TopPub QueryMaxRecommendCount failed,err:", TopBinder)
			controller.ErrorResp(c, 211, "无发布内容")
			return
		}
		errX := mc.ArticleService.ChangeTopPubStatus(TopBinder.Id, maxCount+1)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该文章相关信息")
				return
			} else {
				log.Println("TopPub ChangeTopPubStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 212, "服务器错误")
				return
			}
		}
	} else {
		errX := mc.ArticleService.ChangeTopPubStatus(TopBinder.Id, 0)
		if errX != nil {
			if err == mysqlService.NoRecord {
				controller.ErrorResp(c, 202, "无该文章相关信息")
				return
			} else {
				log.Println("TopPub ChangeTopPubStatus failed,err:", TopBinder)
				controller.ErrorResp(c, 213, "服务器错误")
				return
			}
		}
	}
	controller.SuccessResp(c, "置顶状态修改成功")
	return
}

// DeleteComUser 删除 公司/机构 人员的 企业/服务机构 身份
func (mc *ManageController) DeleteComUser(c *gin.Context) {
	recJson := make(map[string]interface{})
	_ = c.BindJSON(&recJson)
	userId := recJson["user_id"].(string)
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		controller.ErrorResp(c, 201, "用户id参数错误")
		return
	}
	err = mc.ArticleService.BatchKillUserAndPub(userIdInt)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "修改失败，无该用户信息")
			return
		} else {
			controller.ErrorResp(c, 211, "服务器错误")
			log.Println("DeleteComUser failed,err :", err, "\n data:", recJson)
			return
		}
	}

	bossIdStr := strconv.Itoa(userIdInt)
	userInfo, _ := mc.UserService.QueryUserInfoByUserId(bossIdStr)
	PubTem := messageUtil.BossBeDel(userInfo.Name, userInfo.Gender)
	websocketBiz.SysMsgPusher(userInfo.Username, PubTem)

	controller.SuccessResp(c, "人员身份删除成功")
}

// BatchPubAdmit 发布待审核批量许可 all:companyId 0
func (mc *ManageController) BatchPubAdmit(c *gin.Context) {
	recJson := make(map[string]interface{})
	_ = c.BindJSON(&recJson)
	companyId := recJson["comId"].(string)
	companyIdInt, err := strconv.Atoi(companyId)
	if err != nil {
		controller.ErrorResp(c, 201, "公司id参数错误")
		return
	}
	ComWaits, err := mc.ArticleService.QueryComPubInfos(companyIdInt, 0)
	if err != nil || len(ComWaits) == 0 {
		controller.SuccessResp(c, "无可通过数据")
		return
	}
	err = mc.ArticleService.BatchChangePubStatus(companyIdInt, 1)
	if err != nil {
		controller.ErrorResp(c, 211, "服务器错误")
		log.Println("BatchPubAdmit failed ,err :", err)
		return
	}
	go func() {
		for _, item := range ComWaits {
			bossIdStr := strconv.Itoa(item.BossId)
			userInfo, _ := mc.UserService.QueryUserInfoByUserId(bossIdStr)
			PubTem := messageUtil.BossPubPassOrNotTem(userInfo.Name, userInfo.Gender,
				item.Title, "", item.CreateTime, 1)
			websocketBiz.SysMsgPusher(userInfo.Username, PubTem)
		}
	}()
	controller.SuccessResp(c, "批量许可成功，批量同意数量："+strconv.Itoa(len(ComWaits)))
}

// GetMassPubHistory 获取群发消息历史
func (mc *ManageController) GetMassPubHistory(c *gin.Context) {
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	massPubHistory := mc.MassSendService.QueryHistory(pageNumInt)
	totalPage := mc.CountService.MassPubHistoryTP()
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "群发消息历史查询成功",
		"data":      massPubHistory,
		"totalPage": totalPage,
	})
}

// AdminDeletePubInfo 管理员对发布的信息进行删除
func (mc *ManageController) AdminDeletePubInfo(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	artId := recJson["art_id"].(string)
	deleteReason := recJson["deleteReason"].(string)
	ArtIdInt, err := strconv.Atoi(artId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	artInfo, err := mc.ArticleService.QueryArtByID(ArtIdInt)
	if err != nil {
		controller.ErrorResp(c, 202, "无对应发布信息")
		return
	}
	err = mc.ArticleService.AdminDeletePubInfo(ArtIdInt)
	if err != nil {
		controller.ErrorResp(c, 215, "删除失败，服务器错误")
		return
	}
	//通知协程
	go func() {
		bossIdStr := strconv.Itoa(artInfo.BossId)
		userInfo, _ := mc.UserService.QueryUserInfoByUserId(bossIdStr)
		PubTem := messageUtil.AdminDeletePubInfo(userInfo.Name, userInfo.Gender,
			artInfo.Title, deleteReason, artInfo.CreateTime, 2)
		websocketBiz.SysMsgPusher(userInfo.Username, PubTem)
	}()
	controller.SuccessResp(c, "删除成功")
}

func (mc *ManageController) GetUserUpVouchers(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	fromUsername := recJson["from_username"].(string)
	TimeId := recJson["time_id"].(string)
	TimeIdInt, err := strconv.Atoi(TimeId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	UpgradeVouchers := mc.VoucherService.AdminQueryUpgradeVouchers(TimeIdInt, fromUsername, c.Request.Host)
	controller.SuccessResp(c, "凭证查询成功", UpgradeVouchers)
}

// QueryUserBasicInfo 查看用户基本信息
func (mc *ManageController) QueryUserBasicInfo(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	desUsername := recJson["username"].(string)
	userBasicInfo, err := mc.UserService.QueryUserBasicInfo(desUsername, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 201, "目标用户不存在")
		return
	}
	controller.SuccessResp(c, "用户基础信息查询成功", userBasicInfo)
}

// GetWaitingUpgrade 获取需要进行升级审核的
// Qualification = 0 为待审核 -1 为全部 targetLevel = 1为企业机构 2为服务机构
func (mc *ManageController) GetWaitingUpgrade(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	Qualification := recJson["Qualification"].(string)
	targetLevel := recJson["targetLevel"].(string)
	pageNum := recJson["pageNum"].(string)
	quaInt, err1 := strconv.Atoi(Qualification)
	tlInt, err2 := strconv.Atoi(targetLevel)
	pageNumInt, err3 := strconv.Atoi(pageNum)
	if err1 != nil || err2 != nil || err3 != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	waitingUpInfo := mc.UpgradeService.QueryWaitingUpgrade(tlInt, quaInt, pageNumInt)
	totalPage := mc.CountService.GetWaitingUpgradeTP(quaInt, tlInt)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "待升级信息查询成功",
		"data":      waitingUpInfo,
		"totalPage": totalPage,
	})
}

// EditLabel 修改label标签
func (mc *ManageController) EditLabel(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	labelId := recJson["labelId"].(string)
	labelContent := recJson["labelContent"].(string)
	labelIdInt, err := strconv.Atoi(labelId)
	if err != nil {
		controller.ErrorResp(c, 201, "标签id参数格式错误")
		return
	}
	//TODO 查重，需要parent_id
	exist := strings.Index(labelContent, "&")
	if exist != -1 {
		controller.ErrorResp(c, 201, "不符合命名规则，不能存在'&'")
		return
	}
	err = mc.LabelService.EditLabel(labelIdInt, labelContent)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "数据库无匹配内容")
			return
		} else {
			controller.ErrorResp(c, 215, "修改失败，服务器错误")
			return
		}
	}
	controller.SuccessResp(c, "修改成功")
}

// ApplyPassFail 招聘/需求 审核不通过
func (mc *ManageController) ApplyPassFail(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	articleId := recJson["articleId"].(string)
	failReason := recJson["failReason"].(string)
	articleIdInt, _ := strconv.Atoi(articleId)
	artInfo, _ := mc.ArticleService.QueryArtByID(articleIdInt)
	if artInfo.Status != 0 {
		controller.ErrorResp(c, 201, "操作失败，已对该[招聘/请求]做出操作")
		return
	}
	err := mc.ArticleService.ModifyArtStatus(articleIdInt, 2)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		return
	}
	go func() {
		bossIdStr := strconv.Itoa(artInfo.BossId)
		userInfo, _ := mc.UserService.QueryUserInfoByUserId(bossIdStr)
		PubTem := messageUtil.BossPubPassOrNotTem(userInfo.Name, userInfo.Gender,
			artInfo.Title, failReason, artInfo.CreateTime, 2)
		websocketBiz.SysMsgPusher(userInfo.Username, PubTem)
	}()
	controller.SuccessResp(c, "操作成功")
}

// ApplyPassSuccess 招聘/需求 审核通过
func (mc *ManageController) ApplyPassSuccess(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	articleId := recJson["articleId"].(string)
	articleIdInt, _ := strconv.Atoi(articleId)
	artInfo, _ := mc.ArticleService.QueryArtByID(articleIdInt)
	if artInfo.Status != 0 {
		controller.ErrorResp(c, 201, "操作失败，已对该[招聘/请求]做出操作")
		return
	}
	err := mc.ArticleService.ModifyArtStatus(articleIdInt, 1)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		return
	}
	go func() {
		bossIdStr := strconv.Itoa(artInfo.BossId)
		userInfo, _ := mc.UserService.QueryUserInfoByUserId(bossIdStr)
		PubTem := messageUtil.BossPubPassOrNotTem(userInfo.Name, userInfo.Gender,
			artInfo.Title, "", artInfo.CreateTime, 1)
		websocketBiz.SysMsgPusher(userInfo.Username, PubTem)
	}()
	controller.SuccessResp(c, "操作成功")
}

// SendColonyMsg 群发消息(ALL -1,usual 0,boss 1,service 2)
func (mc *ManageController) SendColonyMsg(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	adminUsername := tokenUtil.GetUsernameByToken(c)
	desUserRole := recJson["desUserRole"].(string)
	sendMsg := recJson["sendMsg"].(string)
	desRoleInt, _ := strconv.Atoi(desUserRole)
	if desRoleInt != -1 && desRoleInt != 0 && desRoleInt != 1 && desRoleInt != 2 {
		controller.ErrorResp(c, 201, "目标用户群体参数错误")
		return
	}
	err := mc.MassSendService.AddMassPubHistory(
		adminUsername, sendMsg, &timeUtil.MyTime{Time: time.Now()}, desRoleInt)
	if err != nil {
		controller.ErrorResp(c, 211, "操作失败，服务器错误")
		return
	}
	err = websocketBiz.MassSendMsg(desRoleInt, sendMsg)
	if err != nil {
		controller.ErrorResp(c, 202, "发送失败，未查询到对应群体信息")
		return
	}
	controller.SuccessResp(c, "目标群体群发成功")
	return
}

//// GetApplyPubArt 获取需要审核发布的 pageSize = 15
//func (mc *ManageController) GetApplyPubArt(c *gin.Context) {
//	queryLabel := c.Param("queryLabel")
//	pageNum := c.Param("pageNum")
//	labelLevel := c.Param("labelLevel")
//	labelLevelInt, err := strconv.Atoi(labelLevel)
//	if err != nil {
//		controller.ErrorResp(c, 201, "等级参数错误")
//		return
//	}
//	pageNumInt, err := strconv.Atoi(pageNum)
//	if err != nil {
//		controller.ErrorResp(c, 201, "页码参数错误")
//		return
//	}
//	if labelLevelInt == 1 {
//		fatherLabelInfo, err := mc.LabelService.QueryLabelInfoByLabel(queryLabel)
//		if err != nil {
//			controller.ErrorResp(c, 202, "无该标签对应信息")
//			return
//		}
//		sonLabels, _ := mc.LabelService.QuerySonLabelsById(fatherLabelInfo.ID)
//		if len(sonLabels) == 0 {
//			controller.ErrorResp(c, 203, "该标签无子集内容")
//			return
//		}
//		sonSli := make([]int, 0)
//		for _, item := range sonLabels {
//			sonSli = append(sonSli, item.ID)
//		}
//		waitingInfo := mc.ArticleService.QuerySonWaitingApply(sonSli, pageNumInt)
//		totalPage := mc.CountService.GetSonWaitingTotalPage(sonSli)
//		c.JSON(http.StatusOK, gin.H{
//			"status":    200,
//			"msg":       "待审核大类信息查询成功",
//			"data":      waitingInfo,
//			"totalPage": totalPage,
//		})
//		return
//	} else {
//		waitingInfo := mc.ArticleService.QueryWaitingApply(queryLabel, pageNumInt)
//		totalPage := mc.CountService.GetWaitingTotalPage(queryLabel)
//		c.JSON(http.StatusOK, gin.H{
//			"status":    200,
//			"msg":       "待审核信息查询成功",
//			"data":      waitingInfo,
//			"totalPage": totalPage,
//		})
//		return
//	}
//}

// UpgradeAdmit 身份升级承认
func (mc *ManageController) UpgradeAdmit(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	upgradeId := recJson["upgradeId"].(string)
	upgradeIdInt, _ := strconv.Atoi(upgradeId)
	upgradeInfo, err := mc.UpgradeService.QueryUpgradeInfoById(upgradeIdInt)
	if err != nil {
		controller.ErrorResp(c, 204, "无对应信息")
		return
	}
	if upgradeInfo.Qualification != 0 {
		controller.ErrorResp(c, 201, "该申请已做出审核")
		return
	}
	//err = mc.UserService.ModifyPersonalInfoByUpgrade(upgradeInfo.FromUsername,
	//	upgradeInfo.CompanyId, upgradeInfo.TargetLevel)
	//if err != nil {
	//	controller.ErrorResp(c, 215, "升级失败，服务器错误")
	//	log.Println("UpgradeAdmit failed,err:", err, "req info:", recJson)
	//	return
	//}
	//errMUQ := mc.UpgradeService.ModifyUpgradeQualification(upgradeIdInt, 1)
	//if errMUQ != nil {
	//	controller.ErrorResp(c, 216, "升级失败，服务器错误")
	//	log.Println("UpgradeAdmit failed,err:", errMUQ, "req info:", recJson)
	//	return
	//}
	errUp := mc.UpgradeService.UpgradeInfoChanger(upgradeInfo.FromUsername,
		upgradeInfo.CompanyId, upgradeInfo.TargetLevel, upgradeIdInt, upgradeInfo.CompanyExist)
	if errUp != nil {
		controller.ErrorResp(c, 215, "升级失败，服务器错误")
		log.Println("UpgradeAdmit failed,err:", err, "req info:", recJson)
		return
	}
	//系统消息推送
	go func() {
		applyUserBasicInfo, _ := mc.UserService.QueryUserBasicInfo(upgradeInfo.FromUsername, c.Request.Host)
		succTem := messageUtil.UpgradeApplySuccessTem(
			applyUserBasicInfo.Name,
			applyUserBasicInfo.Gender,
			upgradeInfo.ApplyTime,
		)
		websocketBiz.SysMsgPusher(upgradeInfo.FromUsername, succTem)
	}()
	controller.SuccessResp(c, "审核通过操作成功")
	return
}

// UpgradeDisAdmit 身份升级拒绝承认
func (mc *ManageController) UpgradeDisAdmit(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	upgradeId := recJson["upgradeId"].(string)
	rejectReason := recJson["rejectReason"].(string)
	upgradeIdInt, _ := strconv.Atoi(upgradeId)
	upgradeInfo, err := mc.UpgradeService.QueryUpgradeInfoById(upgradeIdInt)
	if err != nil {
		controller.ErrorResp(c, 204, "操作失败，无对应信息")
		return
	}
	if upgradeInfo.Qualification != 0 {
		controller.ErrorResp(c, 201, "该申请已做出审核")
		return
	}
	errx := mc.UpgradeService.ModifyUpgradeQualification(upgradeIdInt, 2)
	if errx != nil {
		controller.ErrorResp(c, 214, "操作失败，服务器错误")
		log.Println("UpgradeDisAdmit failed,err:", errx, "\n req info:", recJson)
		return
	}
	go func() {
		applyUserBasicInfo, _ := mc.UserService.QueryUserBasicInfo(upgradeInfo.FromUsername, c.Request.Host)
		ErrTem := messageUtil.UpgradeApplyErrTem(
			applyUserBasicInfo.Name,
			applyUserBasicInfo.Gender,
			rejectReason,
			upgradeInfo.ApplyTime,
		)
		websocketBiz.SysMsgPusher(upgradeInfo.FromUsername, ErrTem)
	}()
	controller.SuccessResp(c, "审核不通过操作成功")
	return
}

// ManageArtStatus 隐藏和取消隐藏文章 action:hide
func (mc *ManageController) ManageArtStatus(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	ArtId := recJson["artId"].(string)
	action := recJson["action"].(string)
	ArtIdInt, _ := strconv.Atoi(ArtId)
	artInfo, err := mc.ArticleService.QueryArtByID(ArtIdInt)
	if err != nil {
		controller.ErrorResp(c, 201, "文章不存在")
		return
	}
	err = mc.ArticleService.BossChangeArtShowStatus(artInfo.BossId, ArtIdInt, action)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "文章状态修改失败，身份不符或文章不存在")
			return
		} else {
			controller.ErrorResp(c, 215, "文章状态修改失败，服务器错误")
			return
		}
	}
	controller.SuccessResp(c, "文章状态修改成功")
	return
}

// AddLabel 标签管理-添加  parentLevel = 0为主级
func (mc *ManageController) AddLabel(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	labelType := recJson["labelType"].(string)
	parentId := recJson["parentId"].(string)
	parentLevel := recJson["parentLevel"].(string)
	label := recJson["label"].(string)
	parentIdInt, err := strconv.Atoi(parentId)
	if err != nil {
		controller.ErrorResp(c, 201, "parentId参数错误")
		return
	}
	parentLevelInt, err := strconv.Atoi(parentLevel)
	if err != nil {
		controller.ErrorResp(c, 201, "parentId参数错误")
		return
	}
	exist := strings.Index(label, "&")
	if exist != -1 {
		controller.ErrorResp(c, 201, "不符合命名规则，不能存在'&'")
		return
	}
	RepeatFlag := mc.LabelService.CheckDuplicateLabel(label, labelType, parentIdInt)
	if RepeatFlag == false {
		controller.ErrorResp(c, 201, "添加失败，标签重复")
		return
	}
	err = mc.LabelService.AddLabel(labelType, label, parentIdInt, parentLevelInt)
	if err != nil {
		controller.ErrorResp(c, 215, "添加失败，服务器错误")
		return
	}
	controller.SuccessResp(c, "标签添加成功")
}

// DeleteLabel 标签管理-删除
func (mc *ManageController) DeleteLabel(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	labelId := recJson["labelId"].(string)
	labelLevel := recJson["labelLevel"].(string)
	labelIdInt, err := strconv.Atoi(labelId)
	if err != nil {
		controller.ErrorResp(c, 201, "labelId参数错误")
		return
	}
	labelLevelInt, err := strconv.Atoi(labelLevel)
	if err != nil {
		controller.ErrorResp(c, 201, "labelLevel参数错误")
		return
	}
	switch labelLevelInt {
	case 1:
		if labelIdInt == 0 {
			controller.ErrorResp(c, 212, "已阻止高风险业务逻辑")
			log.Println("DeleteLabel err,ReqData:", recJson)
			return
		}
		sonLabels, err := mc.LabelService.QuerySonLabelsById(labelIdInt)
		if err == nil {
			sonLabelSli := make([]int, 0)
			for _, item := range sonLabels {
				sonLabelSli = append(sonLabelSli, item.ID)
			}
			errIn := mc.ArticleService.BatchDeleteFatherLabelPub(sonLabelSli)
			if errIn != nil {
				controller.ErrorResp(c, 215, "操作失败，服务器错误")
				return
			}
		}
	case 2:
		err = mc.ArticleService.BatchDeleteSonLabelPub(labelIdInt)
	default:
		controller.ErrorResp(c, 202, "等级参数错误")
		return
	}
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		return
	}
	//BatchDeleteSonLabelPub
	err = mc.LabelService.DeleteLabel(labelIdInt, labelLevelInt)
	if err != nil {
		controller.ErrorResp(c, 215, "操作失败，服务器错误")
		return
	}
	controller.SuccessResp(c, "标签删除成功")
}
