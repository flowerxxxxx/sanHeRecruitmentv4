package mysql_service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"time"
)

type DeliveryService struct {
}

func (ds *DeliveryService) AddDeliveryService(bossId, artId int, fromUsername string, deliveryTime time.Time) error {
	var deliveryInfo mysqlModel.Delivery
	err := dao.DB.Where("boss_id=?", bossId).Where("art_id=?", artId).Find(&deliveryInfo).Error
	if err == nil {
		return HasFound
	}
	deliveryInfo.BossId = bossId
	deliveryInfo.ArtId = artId
	deliveryInfo.FromUsername = fromUsername
	deliveryInfo.Qualification = 0
	deliveryInfo.Read = 0
	deliveryInfo.DeliveryTime = deliveryTime
	err = dao.DB.Save(&deliveryInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (ds *DeliveryService) QueryDeliveryById(deliverId int) (deliveryInfo mysqlModel.Delivery, err error) {
	err = dao.DB.Where("id=?", deliverId).Find(&deliveryInfo).Error
	if err != nil {
		return deliveryInfo, err
	}
	return deliveryInfo, err
}

func (ds *DeliveryService) ModifyDeliveryQualification(deliverId, qualification int) (err error) {
	var deliveryInfo mysqlModel.Delivery
	dao.DB.Where("id=?", deliverId).Find(&deliveryInfo)
	deliveryInfo.Qualification = qualification
	err = dao.DB.Save(&deliveryInfo).Error
	return
}

// QueryDelivery 通过qualification查询
func (ds *DeliveryService) QueryDelivery(username, host string, qualification, read, pageNum int) (deliveryInfo []mysqlModel.JobProgressOut, err error) {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	sql := "SELECT deliveries.art_id,articles.title,articles.salary_min,articles.salary_max,users.nickname,users.president,qualification,delivery_time,`read`,head_pic " +
		"FROM `deliveries` INNER JOIN articles on articles.art_id = deliveries.art_id INNER JOIN users ON users.user_id = deliveries.boss_id " +
		"WHERE from_username = ? AND qualification = ? AND `read` = ? order by delivery_time desc limit ?,?"
	err = dao.DB.Raw(sql, username, qualification, read, sqlPage, 10).Scan(&deliveryInfo).Error
	for i := 0; i < len(deliveryInfo); i++ {
		deliveryInfo[i].DeliveryTimeOut = timeUtil.TimeFormatToStr(deliveryInfo[i].DeliveryTime)
		deliveryInfo[i].HeadPic = formatUtil.GetPicHeaderBody(host, deliveryInfo[i].HeadPic)
	}
	return
}

// QueryAllDelivery 通过username查询
func (ds *DeliveryService) QueryAllDelivery(username, host string, pageNum int) (deliveryInfo []mysqlModel.JobProgressOut, err error) {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	//sql := "SELECT deliveries.art_id,articles.title,articles.salary_min,articles.salary_max,users.nickname,users.president,qualification,delivery_time,`read`,head_pic " +
	//	"FROM `deliveries` INNER JOIN articles on articles.art_id = deliveries.art_id INNER JOIN users ON users.user_id = deliveries.boss_id " +
	//	"WHERE from_username = ? limit ?,?"
	//err = dao.DB.Raw(sql, username, sqlPage, 10).Scan(&deliveryInfo).Error
	err = dao.DB.Table("deliveries").Select("deliveries.art_id,articles.title,articles.salary_min,articles.salary_max,users.nickname,users.president,qualification,delivery_time,`read`,head_pic").
		Joins("INNER JOIN articles on articles.art_id = deliveries.art_id INNER JOIN users ON users.user_id = deliveries.boss_id").
		Where("from_username = ?", username).Order("delivery_time desc").Limit(10).Offset(sqlPage).Find(&deliveryInfo).Error
	for i := 0; i < len(deliveryInfo); i++ {
		deliveryInfo[i].DeliveryTimeOut = timeUtil.TimeFormatToStr(deliveryInfo[i].DeliveryTime)
		deliveryInfo[i].HeadPic = formatUtil.GetPicHeaderBody(host, deliveryInfo[i].HeadPic)
	}
	return
}

func (ds *DeliveryService) BossQueryDeliveries(bossId, qualification, pageNum int,
) (info []mysqlModel.BossDeliveries, err error) {
	var bossQuery []mysqlModel.BossDeliveries
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	bossQ := dao.DB.Table("deliveries").Select("deliveries.id,`name`,delivery_time,title,deliveries.qualification,from_username").
		Joins("INNER JOIN users on users.username = deliveries.from_username").
		Joins("INNER JOIN articles on articles.art_id = deliveries.art_id").
		Where("deliveries.boss_id = ?", bossId)
	if qualification != -1 {
		bossQ = bossQ.Where("qualification = ?", qualification)
	}
	err = bossQ.Order("delivery_time desc").Limit(10).Offset(sqlPage).Find(&bossQuery).Error
	if err != nil {
		return
	}
	for i, n := 0, len(bossQuery); i < n; i++ {
		bossQuery[i].DeliveryTimeOut = timeUtil.TimeFormatToStr(bossQuery[i].DeliveryTime)
	}
	return bossQuery, nil
}

func (ds *DeliveryService) SetDeliveryR1ead(deliverId, bossID int) (err error) {
	var deliveryInfo mysqlModel.Delivery
	err = dao.DB.Where("id = ?", deliverId).Where("boss_id=?", bossID).Find(&deliveryInfo).Error
	if err != nil {
		return NoRecord
	}
	deliveryInfo.Read = 1
	err = dao.DB.Save(&deliveryInfo).Error
	if err != nil {
		return err
	}
	return
}
