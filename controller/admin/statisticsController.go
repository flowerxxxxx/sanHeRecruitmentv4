package admin

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"sanHeRecruitment/controller"
	"sanHeRecruitment/models/BindModel/adminBind"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/module/backupModule"
	"sanHeRecruitment/module/controllerModule"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/util/excelUtil"
)

//本层的逻辑全部用于统计以及系统静态存储操作

type StatisticsController struct {
	*mysqlService.UserService
	*mysqlService.CompanyService
	*mysqlService.UpgradeService
	*mysqlService.JobService
	controllerModule.StatisticsModule
}

// StatisticsControllerRouterToken 统计者控制层 本层的逻辑全部用于统计以及系统静态存储操作
func StatisticsControllerRouterToken(router *gin.RouterGroup) {
	sCon := StatisticsController{}
	//统计用户相关总数
	router.GET("/TotalUsers", sCon.TotalUsers)
	//统计公司相关总数
	router.GET("/TotalCompanies", sCon.TotalCompanies)
	//统计在线用户总数
	router.GET("/TotalOnline", sCon.TotalOnline)
	//手动备份数据
	router.POST("/backupData", sCon.BackupData)
	//公司审核升级数据xls文件导出
	router.POST("/waitingUpgradeExport", sCon.WaitingUpgradeExport)
	//公司发布的全部需求信息xls文件导出
	router.POST("/CompanyPubExport", sCon.CompanyPubExport)

}

// CompanyPubExport 公司发布的全部需求信息xls文件导出
func (sc StatisticsController) CompanyPubExport(c *gin.Context) {
	var cpBinder adminBind.ComPubBinder
	err := c.ShouldBind(&cpBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	companyInfo, err2 := sc.CompanyService.QueryCompanyInfoById(cpBinder.ComId, c.Request.Host)
	if err2 != nil {
		controller.ErrorResp(c, 202, "公司不存在")
		return
	}
	comArtInfo := sc.JobService.GetAllCompanyJobInfoXls(cpBinder.ComId)
	if len(comArtInfo) == 0 {
		controller.ErrorResp(c, 202, "导出失败，暂无内容")
		return
	}
	res := sc.StatisticsModule.ComPubSliChanger(comArtInfo)
	content := excelUtil.ToExcel([]string{`公司名称`, `公司规模`, `需求标题`, `观看量`, `发布者昵称`, `发布者姓名`, `标签`, `最低报酬`, `最高报酬`, `需求内容`, `需求类型`, `状态`}, res)
	controller.ResponseXls(c, content, url.QueryEscape("公司发布一览-"+companyInfo.CompanyName))
}

// WaitingUpgradeExport xls文件导出,审核升级项
func (sc StatisticsController) WaitingUpgradeExport(c *gin.Context) {
	var WUBinder adminBind.WaitingUpBinder
	err := c.ShouldBind(&WUBinder)
	if err != nil {
		controller.ErrorResp(c, 201, "参数绑定失败")
		return
	}
	xlsInputer, err := sc.UpgradeService.QueryAllWaitingUpgradeXls(WUBinder.TargetLevel, -1)
	if err != nil {
		if err == mysqlService.NoRecord {
			controller.ErrorResp(c, 202, "导出失败，暂无内容")
			return
		} else {
			controller.ErrorResp(c, 211, "导出失败，服务器错误")
			log.Println("WaitingUpgradeExport failed,err:", err)
			return
		}
	}
	res := sc.StatisticsModule.UpgradeSliChanger(xlsInputer)
	if WUBinder.TargetLevel == 1 {
		content := excelUtil.ToExcel([]string{`公司名称`, `公司联系方式`, `公司地址`, `公司描述`, `公司状态`, `通过情况`, `申请人`, `申请时间`}, res)
		role := "Company"
		controller.ResponseXls(c, content, "Upgrade-"+role)
		return
	} else {
		content := excelUtil.ToExcel([]string{`机构名称`, `机构联系方式`, `机构地址`, `机构描述`, `机构状态`, `通过情况`, `申请人`, `申请时间`}, res)
		role := "Serve"
		controller.ResponseXls(c, content, "Upgrade-"+role)
		return
	}
}

// BackupData 手动备份数据
func (sc StatisticsController) BackupData(c *gin.Context) {
	ziperName, err := backupModule.Backer()
	if err != nil {
		controller.ErrorResp(c, 211, "服务器错误")
		log.Println("admin BackupData failed,err:", err)
		return
	}
	Url := "https://" + c.Request.Host + "/backup/" + ziperName
	controller.SuccessResp(c, "recommend gets success", Url)
	return
}

// TotalUsers 统计用户相关总数
func (sc StatisticsController) TotalUsers(c *gin.Context) {
	totalCounter, _ := sc.UserService.TotalCount()
	controller.SuccessResp(c, "total users ok", totalCounter)
}

// TotalCompanies 统计公司相关总数
func (sc StatisticsController) TotalCompanies(c *gin.Context) {
	totalCounter, _ := sc.CompanyService.TotalCount()
	controller.SuccessResp(c, "total companies ok", totalCounter)
}

// TotalOnline 统计公司相关总数
func (sc StatisticsController) TotalOnline(c *gin.Context) {
	totalOnliner := websocketModel.ReadTotalRecManClients()
	controller.SuccessResp(c, "TotalOnline ok", totalOnliner)
}
