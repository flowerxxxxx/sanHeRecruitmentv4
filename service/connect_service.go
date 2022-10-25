package service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/timeUtil"
)

type ConnectService struct {
}

// SaveNewConnection 添加新的平台联系方式
func (cs *ConnectService) SaveNewConnection(
	desPerson, connect, conType string, uploadTime *timeUtil.MyTime) (err error) {
	var newConSaver mysqlModel.Connection
	newConSaver.Connect = connect
	newConSaver.DesPerson = desPerson
	newConSaver.Type = conType
	newConSaver.UploadTime = uploadTime
	newConSaver.UpdateTime = uploadTime
	err = dao.DB.Table("connections").Save(&newConSaver).Error
	return
}

// EditConnection 编辑平台联系方式
func (cs *ConnectService) EditConnection(id int, DesPerson, Connect, conType string, updateTime *timeUtil.MyTime) (err error) {
	var conSaver mysqlModel.Connection
	err = dao.DB.Table("connections").Where("id = ?", id).Find(&conSaver).Error
	if err != nil {
		return NoRecord
	}
	conSaver.DesPerson = DesPerson
	conSaver.Connect = Connect
	conSaver.Type = conType
	conSaver.UpdateTime = updateTime
	err = dao.DB.Table("connections").Save(&conSaver).Error
	return
}

// QueryConnectionInfos 获取平台联系方式
func (cs *ConnectService) QueryConnectionInfos() []mysqlModel.ConnectionOut {
	var conInfos []mysqlModel.ConnectionOut
	dao.DB.Table("connections").Select("id,des_person,connect,`type`").
		Find(&conInfos)
	return conInfos
}

// DeleteConnection 删除平台联系方式
func (cs *ConnectService) DeleteConnection(ConId int) (err error) {
	err = dao.DB.Table("connections").
		Where("id = ?", ConId).Delete(&mysqlModel.Connection{}).Error
	return
}
