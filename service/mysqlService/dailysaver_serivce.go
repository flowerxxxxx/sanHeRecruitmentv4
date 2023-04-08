package mysqlService

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"time"
)

type DailySaverService struct {
}

func (ds *DailySaverService) AddDailyView(artId int) (err error) {
	var dailyInfo mysqlModel.DailySaver
	NowDate := time.Now().Format("2006-01-02")
	errIn := dao.DB.Table("dailysavers").Where("art_id = ?", artId).Where("date =?", NowDate).Find(&dailyInfo).Error
	if errIn != nil {
		dailyInfo.ArtId = artId
		dailyInfo.DayDelivery = 0
		dailyInfo.DayView = 1
		dailyInfo.Date = NowDate
		err = dao.DB.Table("dailysavers").Save(&dailyInfo).Error
		return
	}
	dailyInfo.DayView = dailyInfo.DayView + 1
	dailyInfo.Date = NowDate
	err = dao.DB.Table("dailysavers").Save(&dailyInfo).Error
	return
}

func (ds *DailySaverService) AddDailyDelivery(artId int) (err error) {
	var dailyInfo mysqlModel.DailySaver
	NowDate := time.Now().Format("2006-01-02")
	errIn := dao.DB.Table("dailysavers").Where("art_id = ?", artId).Where("date =?", NowDate).Find(&dailyInfo).Error
	if errIn != nil {
		dailyInfo.ArtId = artId
		dailyInfo.DayDelivery = 1
		dailyInfo.DayView = 0
		dailyInfo.Date = NowDate
		err = dao.DB.Table("dailysavers").Save(&dailyInfo).Error
		return
	}
	dailyInfo.DayDelivery = dailyInfo.DayDelivery + 1
	err = dao.DB.Table("dailysavers").Save(&dailyInfo).Error
	return
}

// QueryDailyInfo 查询每日数据
func (ds *DailySaverService) QueryDailyInfo(pageNum int, queryType, queryLabel, queryDate, artType string) (dI []mysqlModel.DailyInfo, err error) {
	var DailyInfo []mysqlModel.DailyInfo
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	dailyQ := dao.DB.Table("dailysavers").Select("dailysavers.id,dailysavers.art_id,day_view,day_delivery,date,title").
		Joins("INNER JOIN articles on articles.art_id = dailysavers.art_id "+
			"INNER JOIN labels on labels.id = articles.career_job_id").
		Where("date = ?", queryDate)
	if queryLabel == "allLabel" {
		dailyQ = dailyQ.Where("labels.type = ?", artType)
	} else {
		dailyQ = dailyQ.Where("labels.label = ?", queryLabel)
	}
	if queryType == "day_view" {
		dailyQ = dailyQ.Where("day_view > ?", 0)
	} else {
		dailyQ = dailyQ.Where("day_delivery > ?", 0)
	}
	err = dailyQ.Order(queryType + " desc").Limit(webPageSize).Offset(sqlPage).Find(&DailyInfo).Error
	//for i, n := 0, len(DailyInfo); i < n; i++ {
	//	DailyInfo[i].DateOut = DailyInfo[i].Date.Format("2006-01-02")
	//}
	return DailyInfo, err
}

func (ds *DailySaverService) QueryDaily(queryType, queryLabel, queryDate, artType string) (dI []mysqlModel.DailyInfo, err error) {
	var DailyInfo []mysqlModel.DailyInfo
	dailyQ := dao.DB.Table("dailysavers").Select("dailysavers.id,dailysavers.art_id,day_view,day_delivery,date,title").
		Joins("INNER JOIN articles on articles.art_id = dailysavers.art_id "+
			"INNER JOIN labels on labels.id = articles.career_job_id").
		Where("date = ?", queryDate)
	if queryLabel == "allLabel" {
		dailyQ = dailyQ.Where("labels.type = ?", artType)
	} else {
		dailyQ = dailyQ.Where("labels.label = ?", queryLabel)
	}
	err = dailyQ.Find(&DailyInfo).Error
	return DailyInfo, err
}

func (ds *DailySaverService) QueryDailyLabel(queryDate, artType, queryType string) (dI []mysqlModel.DailyLabelInfo, err error) {
	var DailyInfo []mysqlModel.DailyLabelInfo
	dailyQ := dao.DB.Table("dailysavers").Select("day_view,day_delivery,label").
		Joins("INNER JOIN articles on articles.art_id = dailysavers.art_id "+
			"INNER JOIN labels on labels.id = articles.career_job_id").
		Where("date = ?", queryDate).Where("labels.type = ?", artType)
	if queryType == "day_view" {
		dailyQ = dailyQ.Where("day_view > ?", 0)
	} else {
		dailyQ = dailyQ.Where("day_delivery > ?", 0)
	}
	err = dailyQ.Find(&DailyInfo).Error
	return DailyInfo, err
}
