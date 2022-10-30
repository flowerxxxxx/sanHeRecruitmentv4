package service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
)

type VipShowService struct {
}

func (vss *VipShowService) AddVipShowInfo(cover, content, publisher, title string) (err error) {
	var vssInfo = mysqlModel.VipShow{
		Cover:      cover,
		CreateTime: timeUtil.GetMyTimeNowPtr(),
		UpdateTime: timeUtil.GetMyTimeNowPtr(),
		Clicks:     0,
		Content:    content,
		Publisher:  publisher,
		Title:      title,
	}
	err = dao.DB.Table("vip_shows").Save(&vssInfo).Error
	return
}

func (vss *VipShowService) EditVipShowInfo(id int, cover, content, publisher, title string) (err error) {
	var vssInfo mysqlModel.VipShow
	err = dao.DB.Table("vip_shows").Where("id = ?", id).Find(&vssInfo).Error
	if err != nil {
		return NoRecord
	}
	vssInfo.Cover = cover
	vssInfo.Content = content
	vssInfo.Publisher = publisher
	vssInfo.UpdateTime = timeUtil.GetMyTimeNowPtr()
	vssInfo.Title = title
	err = dao.DB.Table("vip_shows").Save(&vssInfo).Error
	return
}

func (vss *VipShowService) DeleteVipShowInfo(id int) (err error) {
	err = dao.DB.Table("vip_shows").
		Where("id = ?", id).Delete(&mysqlModel.VipShow{}).Error
	return
}

func (vss *VipShowService) QueryVipShowInfos(pageNum int, host string) []mysqlModel.VipShowOut {
	var VipShowInfos []mysqlModel.VipShowOut
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, pageSize)
	dao.DB.Table("vip_shows").Select("*").
		Order("id desc").Offset(sqlPage).Limit(pageSize).Find(&VipShowInfos)
	for i, j := 0, len(VipShowInfos); i < j; i++ {
		VipShowInfos[i].Cover = formatUtil.GetPicHeaderBody(host, VipShowInfos[i].Cover)
	}
	return VipShowInfos
}

func (vss *VipShowService) QueryOneVipShowInfo(id int, host string) (mysqlModel.VipShowOut, error) {
	var VipShowInfo mysqlModel.VipShowOut
	err := dao.DB.Table("vip_shows").
		Where("id = ?", id).Find(&VipShowInfo).Error
	VipShowInfo.Cover = formatUtil.GetPicHeaderBody(host, VipShowInfo.Cover)
	return VipShowInfo, err
}
