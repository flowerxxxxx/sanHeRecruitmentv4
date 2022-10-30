package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/service"
	"strconv"
)

type DataController struct {
	*service.PropagandaService
	*service.NoticeService
	*service.CountService
	*service.ConnectService
	*service.DescribeService
	*service.LabelService
	*service.UpgradeService
	*service.VipShowService
}

func DataControllerRouter(router *gin.RouterGroup) {
	dc := DataController{}
	//获取宣传栏信息（0图片 1视频）
	router.GET("/GetPropagandaInfos", dc.GetPropagandaInfos)
	//获取一个详细的宣传栏信息
	router.GET("/GetOnePropagandaInfo/:pro_id", dc.GetOnePropagandaInfo)
	//获取公告栏信息
	router.GET("/getNoticeInfos/:pageNum", dc.GetNoticeInfos)
	//获取一个详细的公告栏信息
	router.GET("/GetOneNoticeInfo/:notice_id", dc.GetOneNoticeInfo)
	//获取平台联系方式
	router.GET("/getPlaConnectionInfos", dc.GetPlatformConnectionInfos)
	//获取平台简介
	router.GET("/getPlaDescription/:module", dc.GetPlaDescription)
	//获取推荐导航
	router.GET("/getRecommendLabels", dc.getRecommendLabels)
	//获取会员风采
	router.GET("/getVipShows/:pageNum", dc.GetVipShows)
	//获取一个详细的会员风采信息
	router.GET("/GetOneVipShow/:vip_id", dc.GetOnePropagandaInfo)
}

func DataControllerRouterToken(router *gin.RouterGroup) {
	//dc := DataController{}
}

// GetOneVipShow 获取一个详细的会员风采信息
func (dc *DataController) GetOneVipShow(c *gin.Context) {
	vipId := c.Param("vip_id")
	vipIdInt, err := strconv.Atoi(vipId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	vipShowInfo, err := dc.VipShowService.QueryOneVipShowInfo(vipIdInt, c.Request.Host)
	if err != nil {
		c.String(http.StatusNotFound, "404 not found")
		return
	}
	controller.SuccessResp(c, "ok", vipShowInfo)
	return
}

// GetVipShows 获取会员风采信息
func (dc *DataController) GetVipShows(c *gin.Context) {
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	noticeInfos := dc.VipShowService.QueryVipShowInfos(pageNumInt, c.Request.Host)
	noticeTotal := dc.CountService.VipShowInfosTP()
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      noticeInfos,
		"totalPage": noticeTotal,
		"msg":       "会员风采信息息获取成功",
	})
}

func (dc *DataController) getRecommendLabels(c *gin.Context) {
	labelInfos, _ := dc.LabelService.QueryRecommendLabels(1)
	controller.SuccessResp(c, "recommend gets success", labelInfos)
	return
}

// GetOneNoticeInfo 获取一个详细的公告栏信息
func (dc *DataController) GetOneNoticeInfo(c *gin.Context) {
	proId := c.Param("notice_id")
	proIdInt, err := strconv.Atoi(proId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	proInfo, err := dc.NoticeService.QueryOneNoticeInfo(proIdInt)
	if err != nil {
		c.String(http.StatusNotFound, "404 not found")
		return
	}
	controller.SuccessResp(c, "ok", proInfo)
	return
}

// GetPlaDescription 获取平台简介 平台简介：platform；普通用户：0；企业：1，服务机构：2
func (dc *DataController) GetPlaDescription(c *gin.Context) {
	module := c.Param("module")
	desInfos, err := dc.DescribeService.QueryModuleDescriptionInfo(module)
	if err != nil {
		nullDesInfo := mysqlModel.DescriptionOut{Module: module}
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"data":   nullDesInfo,
			"msg":    "no message",
		})
		return
	}
	controller.SuccessResp(c, "ok", desInfos)
	return
}

// GetPlatformConnectionInfos 平台联系方式
func (dc *DataController) GetPlatformConnectionInfos(c *gin.Context) {
	proInfos := dc.ConnectService.QueryConnectionInfos()
	controller.SuccessResp(c, "ok", proInfos)
}

// GetNoticeInfos 获取公告栏信息
func (dc *DataController) GetNoticeInfos(c *gin.Context) {
	pageNum := c.Param("pageNum")
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	noticeInfos := dc.NoticeService.QueryNoticesInfos(pageNumInt)
	noticeTotal := dc.CountService.NoticesInfosTP()
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      noticeInfos,
		"totalPage": noticeTotal,
		"msg":       "公告栏信息获取成功",
	})
}

// GetOnePropagandaInfo 获取一个详细的宣传栏信息
func (dc *DataController) GetOnePropagandaInfo(c *gin.Context) {
	proId := c.Param("pro_id")
	proIdInt, err := strconv.Atoi(proId)
	if err != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	proInfo, err := dc.PropagandaService.QueryOneProInfo(proIdInt, c.Request.Host)
	if err != nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	controller.SuccessResp(c, "ok", proInfo)
}

// GetPropagandaInfos 获取宣传栏信息（0图片 1视频）
func (dc *DataController) GetPropagandaInfos(c *gin.Context) {
	proInfos, err := dc.PropagandaService.QueryProInfos(c.Request.Host)
	if err != nil {
		c.String(http.StatusNotFound, "not found")
		return
	}
	controller.SuccessResp(c, "propaganda infos ok", proInfos)
}
