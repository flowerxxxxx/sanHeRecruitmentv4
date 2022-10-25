package service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
)

type DockService struct {
}

func (ds *DockService) AddDockRecord(comId, bossId, UserId, artId int,
	dockTime *timeUtil.MyTime, bossName, userName, pubTitle string) (err error) {
	var dockInfo mysqlModel.Dock
	dockInfo.DockTime = dockTime
	dockInfo.UserId = UserId
	dockInfo.ComId = comId
	dockInfo.BossId = bossId
	dockInfo.ArtId = artId
	dockInfo.BossName = bossName
	dockInfo.UserName = userName
	dockInfo.PubTitle = pubTitle
	err = dao.DB.Table("docks").Save(&dockInfo).Error
	return
}

func (ds *DockService) QueryDockInfosManyMaps(comId, pageNum int) []*mysqlModel.Dock {
	var dockInfo []*mysqlModel.Dock
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	dao.DB.Table("docks").Select("id,dock_time,docks.art_id,docks.user_id,docks.boss_id,docks.com_id,"+
		"users.`name` AS user_name,boss_users.`name` AS boss_name,title").
		Joins("INNER JOIN users AS boss_users ON boss_users.user_id = docks.boss_id "+
			" INNER JOIN users ON users.user_id = docks.user_id "+
			" INNER JOIN articles ON articles.art_id = docks.art_id").
		Where("com_id = ?", comId).Order("id desc").Limit(webPageSize).Offset(sqlPage).
		Find(&dockInfo)
	return dockInfo
}

func (ds *DockService) QueryDockInfos(comId, pageNum int) []*mysqlModel.Dock {
	var dockInfo []*mysqlModel.Dock
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	dao.DB.Table("docks").Select("id,dock_time,docks.art_id,docks.user_id,docks.boss_id,docks.com_id,"+
		"docks.boss_name,docks.user_name,docks.pub_title").
		Where("com_id = ?", comId).Order("id desc").Limit(webPageSize).Offset(sqlPage).
		Find(&dockInfo)
	return dockInfo
}
