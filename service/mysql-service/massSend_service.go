package mysql_service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
)

type MassSendService struct {
}

func (mss *MassSendService) QueryHistory(pageNum int) []mysqlModel.MassSend {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	var massHis []mysqlModel.MassSend
	dao.DB.Table("mass_sends").Order("id desc").
		Limit(webPageSize).Offset(sqlPage).Find(&massHis)
	return massHis
}

func (mss *MassSendService) AddMassPubHistory(pubUsername, content string, massPubTime *timeUtil.MyTime, desRoleInt int) error {
	var massHis mysqlModel.MassSend
	massHis.PublishUsername = pubUsername
	massHis.Content = content
	massHis.PublishTime = massPubTime
	massHis.DesRole = desRoleInt
	err := dao.DB.Table("mass_sends").Save(&massHis).Error
	return err
}
