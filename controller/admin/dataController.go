package admin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sanHeRecruitment/config"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/BindModel/adminBind"
	"sanHeRecruitment/models/BindModel/userBind"
	"sanHeRecruitment/module/controllerModule"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/osUtil"
	"sanHeRecruitment/util/pageUtil"
	"sanHeRecruitment/util/saveUtil"
	"sanHeRecruitment/util/timeUtil"
	"sanHeRecruitment/util/tokenUtil"
	"sanHeRecruitment/util/uploadUtil"
	"strconv"
	"strings"
	"time"
)

//本文件的router包含小程序首页数据的管理、管理端部分数据的查看

// DataController 数据controller
type DataController struct {
	controllerModule.DataControlModule

	*service.DailySaverService
	*service.CountService
	*service.CompanyService
	*service.DockService
	*service.UserService
	*service.JobService
	*service.LabelService
	*service.ArticleService
	*service.PropagandaService
	*service.NoticeService
	*service.ConnectService
	*service.DescribeService
	*service.VipShowService
}

func DataControllerRouter(router *gin.RouterGroup) {
	//dataC := DataController{}
}

func DataControllerRouterToken(router *gin.RouterGroup) {
	dataC := DataController{}
	//上传流媒体（video&picture）
	router.POST("/SaveStreamingOrPic", dataC.SaveStreamingOrPic)
	//删除流媒体信息（video&picture）
	router.POST("/DeleteStreamOrPic", dataC.DeleteStreamOrPic)

	//获取每日数据（按照day_view或者day_delivery进行排序）
	router.POST("/GetDailyHotData", dataC.QueryDailyData)
	//获取每日热度标签（按照day_view或者day_delivery进行排序）
	router.POST("/GetDailyHotLabel", dataC.GetDailyHotLabel)
	//管理员查看公司信息
	router.POST("/getJobInfo", dataC.GetJobInfos)
	//获取所有公司 comLevel 1为企业机构 2为服务机构
	router.GET("/GetAllCompanies/:comLevel/:pageNum", dataC.GetAllCompanies)
	//获取企业机构/服务机构的对接记录
	router.GET("/getDockInfos/:comId/:pageNum", dataC.GetDockInfos)
	//查看公司/机构的人员
	router.GET("/getComEmployees/:comId/:pageNum", dataC.ShowComEmployees)
	//获取公司发布的待审核的
	router.POST("/GetComPubWaiting", dataC.GetComPubWaiting)
	//获取公司待审核已经包含的工作/需求标签
	router.GET("/GetComPubWaitingLabel/:companyName/:type", dataC.GetComPubWaitingLabel)
	//获取需要审核发布的 all全部，其余为招聘或者需求的类别(如java)
	router.GET("/getWaitingApply/:com_id/:queryLabel/:labelLevel/:pageNum", dataC.getWaitingApplyCom)
	//上传宣传栏目内容
	router.POST("/uploadPropagandaContent", dataC.SavePropagandaContent)
	//编辑宣传栏内容
	router.POST("/EditPropagandaContent", dataC.EditPropagandaContent)
	//删除宣传栏内容
	router.POST("/DeletePropagandaContent", dataC.DeletePropagandaContent)
	//上传公告栏内容
	router.POST("/uploadNotice", dataC.SaveNotice)
	//编辑公告栏内容
	router.POST("/editNotice", dataC.EditNotice)
	//删除公告栏内容
	router.POST("/deleteNotice", dataC.DeleteNotice)
	//上传平台联系方
	router.POST("/addPlatformConnection", dataC.AddPlatformConnection)
	//编辑平台联系方式
	router.POST("/editPlatformConnection", dataC.EditPlatformConnection)
	//删除平台联系方式
	router.POST("/deletePlatformConnection", dataC.DeletePlatformConnection)
	//上传平台简介
	router.POST("/addPlaDescription", dataC.AddPlaDescription)
	//编辑平台简介
	router.POST("/editPlaDescription", dataC.EditPlaDescription)
	//删除平台简介
	router.POST("/deletePlaDescription", dataC.DeletePlaDescription)
	//修改标签的导航栏推荐状态
	router.POST("/ChangeLabelRecommend", dataC.ChangeLabelRecommend)
	//模糊获取工作信息
	router.POST("/FuzzyQJobInfos", dataC.FuzzyQueryJobInfos)
	//添加会员风采
	router.POST("/AddVipStyle", dataC.AddVipStyle)
	//修改会员风采
	router.POST("/EditVipStyle", dataC.EditVipStyle)
	//删除会员风采
	router.POST("/DeleteVipStyle", dataC.DeleteVipStyle)
}

// DeleteVipStyle 删除会员风采
func (dc *DataController) DeleteVipStyle(c *gin.Context) {
	var DelBinder adminBind.VipStyleDelBinder
	err := c.ShouldBind(&DelBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	oldInfo, errFound := dc.VipShowService.QueryOneVipShowInfo(DelBinder.Id, c.Request.Host)
	if errFound != nil {
		controller.ErrorResp(c, 202, "该内容已丢失")
		return
	}
	err = dc.VipShowService.DeleteVipShowInfo(DelBinder.Id)
	if err != nil {
		controller.ErrorResp(c, 211, "删除失败，无相关信息或服务器错误")
		log.Println("DeleteVipStyle failed,err:", err, "\n request info:", DelBinder)
		return
	}
	go func() {
		ed := saveUtil.DeletePicSaver(oldInfo.Cover)
		if ed != nil {
			log.Println("DeleteVipStyle DeletePicSaver failed,err:", ed)
		}
	}()
	controller.SuccessResp(c, "会员风采删除成功")
	return
}

// EditVipStyle 修改会员风采
func (dc *DataController) EditVipStyle(c *gin.Context) {
	var vipSBinder adminBind.VipStyleEditBinder
	err := c.ShouldBind(&vipSBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	newCoverUrl, err := formatUtil.SavePicHeaderCutter(vipSBinder.CoverUrl)
	if err != nil {
		controller.ErrorResp(c, 203, "url路径错误")
		return
	}
	oldInfo, errF := dc.VipShowService.QueryOneVipShowInfo(vipSBinder.Id, c.Request.Host)
	if errF != nil {
		controller.ErrorResp(c, 204, "无对应内容")
		return
	}
	if newCoverUrl != oldInfo.Cover {
		ed := saveUtil.DeletePicSaver(oldInfo.Cover)
		if ed != nil {
			log.Println("EditVipStyle DeletePicSaver failed,err:", ed)
		}
	}
	errAdd := dc.VipShowService.EditVipShowInfo(vipSBinder.Id, newCoverUrl, vipSBinder.Content, tokenUtil.GetUsernameByToken(c), vipSBinder.Title)
	if errAdd != nil {
		if errAdd == service.NoRecord {
			controller.ErrorResp(c, 202, "无id对应内容")
			return
		}
		controller.ErrorResp(c, 211, "会员风采修改失败，服务器错误")
		log.Println("EditVipStyle failed,err:", errAdd, "\nrequest info:", vipSBinder)
		return
	}
	controller.SuccessResp(c, "会员风采添加成功")
	return
}

// AddVipStyle 添加会员风采
func (dc *DataController) AddVipStyle(c *gin.Context) {
	var vipSBinder adminBind.VipStyleAddBinder
	err := c.ShouldBind(&vipSBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	newCoverUrl, err := formatUtil.SavePicHeaderCutter(vipSBinder.CoverUrl)
	if err != nil {
		controller.ErrorResp(c, 203, "url路径错误")
		return
	}
	errAdd := dc.VipShowService.AddVipShowInfo(newCoverUrl, vipSBinder.Content, tokenUtil.GetUsernameByToken(c), vipSBinder.Title)
	if errAdd != nil {
		controller.ErrorResp(c, 211, "会员风采添加失败，服务器错误")
		log.Println("AddVipStyle failed,err:", errAdd, "\nrequest info:", vipSBinder)
		return
	}
	controller.SuccessResp(c, "会员风采添加成功")
	return
}

// FuzzyQueryJobInfos 模糊获取工作信息
func (dc *DataController) FuzzyQueryJobInfos(c *gin.Context) {
	var fBinder userBind.FuzzyQueryJobs
	err := c.ShouldBind(&fBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	fuzzyJobInfo := dc.JobService.FuzzyQueryJobs(fBinder.FuzzyName, fBinder.QueryType, c.Request.Host, fBinder.PageNum, 1)
	TotalPageNum := dc.CountService.GetFuzzyQueryJobsTP(fBinder.FuzzyName, fBinder.QueryType, 1)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"data":      fuzzyJobInfo,
		"totalPage": TotalPageNum,
		"msg":       "模糊招聘信息获取成功",
	})
}

// ChangeLabelRecommend 标记推荐标签
func (dc *DataController) ChangeLabelRecommend(c *gin.Context) {
	var RecoBinder adminBind.AddRecoLabelBinder
	err := c.ShouldBind(&RecoBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.LabelService.ChangeLabelRecommend(RecoBinder.Id, RecoBinder.DesReco)
	if err != nil {
		if err == service.NoRecord {
			controller.ErrorResp(c, 202, "无该标签相关信息")
			return
		} else {
			log.Println("AddRecommendLabel failed,err:", RecoBinder)
			controller.ErrorResp(c, 212, "服务器错误")
			return
		}
	}
	controller.SuccessResp(c, "标签导航信息修改成功")
	return
}

// DeletePlaDescription 删除平台简介
func (dc *DataController) DeletePlaDescription(c *gin.Context) {
	var DelBinder adminBind.DelPlatDesBinder
	err := c.ShouldBind(&DelBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.DescribeService.DeleteDescription(DelBinder.Id)
	if err != nil {
		controller.ErrorResp(c, 211, "删除失败，无相关信息或服务器错误")
		log.Println("DeletePlaDescription failed,err:", err, "request body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "平台简介删除成功")
	return
}

// EditPlaDescription 编辑平台简介
func (dc *DataController) EditPlaDescription(c *gin.Context) {
	var editSaver adminBind.EditPlatDesBinder
	err := c.ShouldBind(&editSaver)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.DescribeService.EditDescription(
		editSaver.Id, editSaver.Content, editSaver.Module,
		timeUtil.GetMyTimeNowPtr())
	if err != nil {
		if err == service.NoRecord {
			controller.ErrorResp(c, 202, "修改失败，为找到匹配信息")
			return
		} else {
			controller.ErrorResp(c, 211, "修改失败，服务器错误")
			log.Println("EditPlaDescription failed,err:", err, "request:", c.Request.Body)
			return
		}
	}
	controller.SuccessResp(c, "平台简介编辑成功")
	return
}

// AddPlaDescription 上传平台简介
func (dc *DataController) AddPlaDescription(c *gin.Context) {
	var plaDesBinder adminBind.PlatDescriptionBinder
	err := c.ShouldBind(&plaDesBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.DescribeService.SaveNewDescription(plaDesBinder.Content, plaDesBinder.Module,
		tokenUtil.GetUsernameByToken(c), timeUtil.GetMyTimeNowPtr())
	if err != nil {
		controller.ErrorResp(c, 211, "上传失败，服务器错误")
		log.Println("AddPlaDescription failed,err:", err, "\nrequest body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "简介上传成功")
	return
}

// DeletePlatformConnection 删除平台联系方式
func (dc *DataController) DeletePlatformConnection(c *gin.Context) {
	var DelConBinder adminBind.DeleteConBinder
	err := c.ShouldBind(&DelConBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.ConnectService.DeleteConnection(DelConBinder.Id)
	if err != nil {
		controller.ErrorResp(c, 211, "删除失败，无相关信息或服务器错误")
		log.Println("DeletePlatformConnection failed,err:", err, "request body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "公告信息删除成功")
	return
}

// EditPlatformConnection 编辑平台联系方式
func (dc *DataController) EditPlatformConnection(c *gin.Context) {
	var editSaver adminBind.EditConnectionBinder
	err := c.ShouldBind(&editSaver)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.ConnectService.EditConnection(
		editSaver.Id, editSaver.DesPerson, editSaver.Connect, editSaver.Type,
		timeUtil.GetMyTimeNowPtr())
	if err != nil {
		if err == service.NoRecord {
			controller.ErrorResp(c, 202, "修改失败，为找到匹配信息")
			return
		} else {
			controller.ErrorResp(c, 211, "修改失败，服务器错误")
			log.Println("EditPlatformConnection failed,err:", err, "request:", editSaver)
			return
		}
	}
	controller.SuccessResp(c, "联系方式编辑成功")
	return
}

// AddPlatformConnection 上传平台联系方式  type :qq,email,phone...
func (dc *DataController) AddPlatformConnection(c *gin.Context) {
	var AddPlaConBinder adminBind.AddPlatformCon
	err := c.ShouldBind(&AddPlaConBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.ConnectService.SaveNewConnection(
		AddPlaConBinder.DesPerson, AddPlaConBinder.Connect, AddPlaConBinder.Type,
		timeUtil.GetMyTimeNowPtr())
	if err != nil {
		controller.ErrorResp(c, 211, "上传失败，服务器错误")
		log.Println("AddPlatformConnection failed,err:", err, "\nrequest body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "联系方式上传成功")
	return
}

// DeleteNotice 删除公告栏内容
func (dc *DataController) DeleteNotice(c *gin.Context) {
	var deleteBinder adminBind.DeleteNoticeBinder
	err := c.ShouldBind(&deleteBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.NoticeService.DeleteNotice(deleteBinder.Id)
	if err != nil {
		controller.ErrorResp(c, 211, "删除失败，无相关信息或服务器错误")
		log.Println("DeleteNotice failed,err:", err, "request body:", c.Request.Body)
		return
	}
	controller.SuccessResp(c, "公告信息删除成功")
	return
}

// EditNotice 编辑公告栏内容
func (dc *DataController) EditNotice(c *gin.Context) {
	var editSaver adminBind.EditNoticeBinder
	err := c.ShouldBind(&editSaver)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	err = dc.NoticeService.EditNotice(editSaver.Id, editSaver.Content, editSaver.Title,
		timeUtil.GetMyTimeNowPtr())
	if err != nil {
		if err == service.NoRecord {
			controller.ErrorResp(c, 202, "修改失败，为找到匹配信息")
			return
		} else {
			controller.ErrorResp(c, 211, "修改失败，服务器错误")
			log.Println("EditNotice failed,err:", err, "request:", editSaver)
			return
		}
	}
	controller.SuccessResp(c, "公告编辑成功")
	return
}

// SaveNotice 上传公告栏内容
func (dc *DataController) SaveNotice(c *gin.Context) {
	var noticeSaver adminBind.NoticeSaveBinder
	err := c.ShouldBind(&noticeSaver)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	uploader := tokenUtil.GetUsernameByToken(c)
	err = dc.NoticeService.SaveNotice(noticeSaver.Content, uploader, noticeSaver.Title,
		timeUtil.GetMyTimeNowPtr())
	if err != nil {
		controller.ErrorResp(c, 211, "服务器错误，上传失败")
		log.Println("SaveNotice failed err:", err, "upload Info:", noticeSaver)
		return
	}
	controller.SuccessResp(c, "公告上传成功")
	return
}

// DeletePropagandaContent 删除宣传栏内容
func (dc *DataController) DeletePropagandaContent(c *gin.Context) {
	var DeleteCont adminBind.DeleteProContBinder
	err := c.ShouldBind(&DeleteCont)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	WaitDeleteInfo, err := dc.PropagandaService.QueryOneProInfo(DeleteCont.ProId, c.Request.Host)
	if err != nil {
		controller.ErrorResp(c, 202, "无相关宣传信息")
		return
	}
	err = dc.PropagandaService.DeleteProInfo(DeleteCont.ProId)
	if err != nil {
		controller.ErrorResp(c, 211, "删除失败")
		log.Println("DeletePropagandaContent failed,err:", err)
		return
	}
	go osUtil.DeleteFile(WaitDeleteInfo.Url)
	controller.SuccessResp(c, "删除成功")
}

// GetJobInfos Admin招聘信息（查找招聘+boss信息）sql为最新发布
func (dc *DataController) GetJobInfos(c *gin.Context) {
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
	jobInfo := dc.JobService.GetJobs(region, labelType, c.Request.Host, cjiInt, pageNumInt, 15, 1)
	TotalPageNum := dc.CountService.GetJobsTotalPage(cjiInt, region, labelType, "web")
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

// EditPropagandaContent 编辑宣传栏内容
func (dc *DataController) EditPropagandaContent(c *gin.Context) {
	var edBinder adminBind.EditPropagandaContent
	err := c.ShouldBind(&edBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	uploader := tokenUtil.GetUsernameByToken(c)
	proInfo, errFound := dc.PropagandaService.QueryOneProInfo(edBinder.ID, c.Request.Host)
	if errFound != nil {
		controller.ErrorResp(c, 203, "该内容已丢失")
		return
	}
	saveFlag := strings.Index(edBinder.Url, "uploadPic/")
	if saveFlag == -1 {
		controller.ErrorResp(c, 204, "url路径错误")
		return
	}
	url := edBinder.Url[saveFlag:]
	if url != proInfo.Url {
		ed := saveUtil.DeletePicSaver(proInfo.Url)
		if ed != nil {
			log.Println("EditPropagandaContent DeletePicSaver failed,err:", ed)
		}
	}
	err = dc.PropagandaService.EditProInfo(edBinder.ID,
		&timeUtil.MyTime{Time: time.Now()}, url,
		uploader, edBinder.Content, edBinder.Title, edBinder.Type)
	if err != nil {
		if err == service.NoRecord {
			controller.ErrorResp(c, 202, "无相关信息")
			return
		} else {
			controller.ErrorResp(c, 211, "服务器错误")
			log.Println("SavePropagandaContent err:", err, "\n requestBody:", edBinder)
			return
		}
	}
	controller.SuccessResp(c, "宣传栏目内容上传成功")

}

// DeleteStreamOrPic 根据url删除本地已经上传的照片凭证
func (dc *DataController) DeleteStreamOrPic(c *gin.Context) {
	var DeleteStream adminBind.DeleteStreamBinder
	err := c.ShouldBind(&DeleteStream)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	streamUrl := DeleteStream.Url
	pos := strings.Index(streamUrl, "/uploadPic")
	finalPicUrl := config.PicSaverPath + streamUrl[pos+10:]
	go func() {
		err := os.Remove(finalPicUrl)
		if err != nil {
			log.Println("DeleteStreamOrPic file remove Error!")
			log.Printf("%s", err)
		}
	}()
	controller.SuccessResp(c, "图片凭证删除成功")
}

// SavePropagandaContent 上传宣传栏目内容
func (dc *DataController) SavePropagandaContent(c *gin.Context) {
	var propagandaBinder adminBind.PropagandaContentBinder
	err := c.ShouldBind(&propagandaBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	uploader := tokenUtil.GetUsernameByToken(c)
	saveFlag := strings.Index(propagandaBinder.Url, "uploadPic/")
	if saveFlag == -1 {
		controller.ErrorResp(c, 203, "url路径错误")
		return
	}
	url := propagandaBinder.Url[saveFlag:]
	err = dc.PropagandaService.AddProInfo(
		&timeUtil.MyTime{Time: time.Now()}, url,
		uploader, propagandaBinder.Content, propagandaBinder.Title, propagandaBinder.Type)
	if err != nil {
		controller.ErrorResp(c, 211, "服务器错误")
		log.Println("SavePropagandaContent err:", err, "\n requestBody:", propagandaBinder)
		return
	}
	controller.SuccessResp(c, "宣传栏目内容上传成功")
}

// SaveStreamingOrPic 上传流媒体（video&picture）
func (dc *DataController) SaveStreamingOrPic(c *gin.Context) {
	file, err := c.FormFile("uploadElement")
	if err != nil {
		c.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fileSizeMB := file.Size / 1024 / 1024 / 1024
	if fileSizeMB > 1 {
		controller.ErrorResp(c, 201, "上传失败，视频大于上限（1GB）")
	}
	fileFormat := file.Filename[strings.LastIndex(file.Filename, "."):]
	//fmt.Println(fileFormat)
	judgeFlag := uploadUtil.FormatJudge(fileFormat, ".mp4", ".jpg", ".png", ".jpeg", "gif")
	if judgeFlag == false {
		//c.String(http.StatusBadRequest, "上传失败，支持格式：\".mp4\", \".jpg\", \".png\",\".jpeg\"")
		c.JSON(http.StatusOK, gin.H{
			"status":  201,
			"message": "上传失败，支持格式：\".mp4\", \".jpg\", \".png\",\".jpeg\"",
			"errno":   1,
		})
		return
	}
	fileUrl, fileAddr := uploadUtil.SaveFormat(fileFormat, c.Request.Host)
	if err := c.SaveUploadedFile(file, fileAddr); err != nil {
		//c.String(http.StatusBadRequest, "保存失败 Error:%s", err.Error())
		c.JSON(http.StatusOK, gin.H{
			"status":  202,
			"message": "流媒体&图片上传失败，服务器错误",
			"errno":   1,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "流媒体&图片上传成功",
		"data": map[string]string{
			"url": formatUtil.GetPicHeaderBody(c.Request.Host, fileUrl),
		},
		"errno": 0,
	})
	return
	//controller.SuccessResp(c, "流媒体&图片上传成功", formatUtil.GetPicHeaderBody(c.Request.Host, fileUrl))
}

// GetApplyPubArt 获取需要审核发布的 pageSize = 15 comId = 0 为全部公司
func (dc *DataController) getWaitingApplyCom(c *gin.Context) {
	queryLabel := c.Param("queryLabel")
	pageNum := c.Param("pageNum")
	labelLevel := c.Param("labelLevel")
	comId := c.Param("com_id")
	labelLevelInt, err := strconv.Atoi(labelLevel)
	comIdInt, err2 := strconv.Atoi(comId)
	pageNumInt, err3 := strconv.Atoi(pageNum)
	if err != nil || err2 != nil || err3 != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	if labelLevelInt == 1 {
		fatherLabelInfo, err := dc.LabelService.QueryLabelInfoByLabel(queryLabel)
		if err != nil {
			controller.ErrorResp(c, 202, "无该标签对应信息")
			return
		}
		sonLabels, _ := dc.LabelService.QuerySonLabelsById(fatherLabelInfo.ID)
		//if len(sonLabels) == 0 {
		//	controller.ErrorResp(c, 203, "该标签无子集内容")
		//	return
		//}
		sonSli := make([]int, 0)
		for _, item := range sonLabels {
			sonSli = append(sonSli, item.ID)
		}
		waitingInfo := dc.ArticleService.QuerySonWaitingApply(sonSli, pageNumInt, comIdInt, c.Request.Host)
		totalPage := dc.CountService.GetSonWaitingTotalPage(sonSli, comIdInt)
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"msg":       "待审核大类信息查询成功",
			"data":      waitingInfo,
			"totalPage": totalPage,
		})
		return
	} else {
		waitingInfo := dc.ArticleService.QueryWaitingApply(queryLabel, c.Request.Host, pageNumInt, comIdInt)
		totalPage := dc.CountService.GetWaitingTotalPage(queryLabel, comIdInt)
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"msg":       "待审核信息查询成功",
			"data":      waitingInfo,
			"totalPage": totalPage,
		})
		return
	}
}

// GetComPubWaitingLabel 获取公司待审核的标签
func (dc *DataController) GetComPubWaitingLabel(c *gin.Context) {
	companyName := c.Param("companyName")
	jobLabel := c.Param("type")
	companyLabel := dc.LabelService.QueryCompanyLabel(companyName, jobLabel, 0)
	controller.SuccessResp(c, companyName+"["+jobLabel+"]标签查询成功", companyLabel)
}

// GetComPubWaiting 获取公司发布的待审核的
func (dc *DataController) GetComPubWaiting(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	companyName := recJson["companyName"].(string)
	JobReqType := recJson["jobReqType"].(string)
	pageNum := recJson["pageNum"].(string)
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	comArtInfo := dc.JobService.GetCompanyJobInfo(companyName, JobReqType, "", c.Request.Host, pageNumInt, 0)
	Total := dc.CountService.CompanyPubTotal(companyName, JobReqType, "", 0)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "查询成功",
		"data":   comArtInfo,
		"total":  Total,
	})
}

// ShowComEmployees 查看公司/机构的人员
func (dc *DataController) ShowComEmployees(c *gin.Context) {
	ComId := c.Param("comId")
	pageNum := c.Param("pageNum")
	ComIdInt, err := strconv.Atoi(ComId)
	pageNumInt, err2 := strconv.Atoi(pageNum)
	if err != nil || err2 != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	comUserInfos := dc.UserService.GetComUsers(pageNumInt, ComIdInt, c.Request.Host)
	totalPage := dc.CountService.GetComUsersTP(ComIdInt)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "公司/机构的人员查询成功",
		"data":      comUserInfos,
		"totalPage": totalPage,
	})
}

func (dc *DataController) GetDockInfos(c *gin.Context) {
	comId := c.Param("comId")
	pageNum := c.Param("pageNum")
	comIdInt, err := strconv.Atoi(comId)
	pageNumInt, err2 := strconv.Atoi(pageNum)
	if err != nil || err2 != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	docInfos := dc.DockService.QueryDockInfos(comIdInt, pageNumInt)
	totalPage := dc.CountService.GetDockInfosTP(comIdInt)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "全部公司信息查询成功",
		"data":      docInfos,
		"totalPage": totalPage,
	})
}

// GetAllCompanies 获取所有公司 comLevel 1为企业机构 2为服务机构
func (dc *DataController) GetAllCompanies(c *gin.Context) {
	comLevel := c.Param("comLevel")
	pageNum := c.Param("pageNum")
	comLInt, err := strconv.Atoi(comLevel)
	pageNumInt, err2 := strconv.Atoi(pageNum)
	if err != nil || err2 != nil {
		controller.ErrorResp(c, 201, "参数错误")
		return
	}
	comInfos, _ := dc.CompanyService.QueryAllCompanies(comLInt, pageNumInt, c.Request.Host)
	totalPage := dc.CountService.QueryAllCompaniesTP(comLInt)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       "全部公司信息查询成功",
		"data":      comInfos,
		"totalPage": totalPage,
	})
	return
}

// QueryDailyData 查询每日数据 按照投递量或查看量排序
// queryType是day_view或day_delivery，queryLabel是标签或者allLabel(全部),
// 当queryLabel是allLabel时，artType 需要是 job 或者 request 其他时候为""（达到全部job或者全部request效果）
func (dc *DataController) QueryDailyData(c *gin.Context) {
	//TODO 按照投递量或查看量排序
	//TODO order by view/delivery
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	pageNum := recJson["pageNum"].(string)
	queryType := recJson["queryType"].(string)
	queryLabel := recJson["queryLabel"].(string)
	queryDate := recJson["queryDate"].(string)
	artType := recJson["artType"].(string)
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	dailyInfo, _ := dc.DailySaverService.QueryDailyInfo(pageNumInt, queryType, queryLabel, queryDate, artType)
	totalPage := dc.CountService.GetDailyInfoTotalPage(queryLabel, artType, queryType, queryDate)
	if len(dailyInfo) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"msg":       queryDate + "无数据",
			"data":      dailyInfo,
			"totalPage": totalPage,
			"showMsg":   true,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"msg":       "每日数据获取成功",
			"data":      dailyInfo,
			"totalPage": totalPage,
			"showMsg":   false,
		})
		return
	}
}

// GetDailyHotLabel 获取热门标签顺序
func (dc *DataController) GetDailyHotLabel(c *gin.Context) {
	recJson := map[string]interface{}{}
	_ = c.BindJSON(&recJson)
	queryDate := recJson["queryDate"].(string)
	queryType := recJson["queryType"].(string)
	artType := recJson["artType"].(string)
	pageNum := recJson["pageNum"].(string)
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		controller.ErrorResp(c, 201, "页码参数错误")
		return
	}
	hotLabelData := dc.DataControlModule.GetDailyHotLabel(queryDate, artType, queryType)
	hotData, err := dc.DataControlModule.HotLabelCutSliByPageNum(hotLabelData, pageNumInt, 15)
	if err != nil {
		//controller.ErrorResp(c, 202, queryDate+"无数据",[]interface{}{})
		//return
		c.JSON(http.StatusOK, gin.H{
			"status":    200,
			"msg":       queryDate + "无数据",
			"data":      []interface{}{},
			"totalPage": 1,
			"showMsg":   true,
		})
		return
	}
	totalPage := pageUtil.TotalPageByTotalNum(len(hotLabelData), 15)
	c.JSON(http.StatusOK, gin.H{
		"status":    200,
		"msg":       artType + "-" + queryType + "热门标签查询成功",
		"data":      hotData,
		"totalPage": totalPage,
		"showMsg":   false,
	})
	return
}
