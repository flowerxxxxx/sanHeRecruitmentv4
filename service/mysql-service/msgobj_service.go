package mysql_service

import (
	"log"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
)

type MsgObjService struct {
}

// GetMsgList 信息列表
func (m *MsgObjService) GetMsgList(fromUsername, host string) []mysqlModel.MsgObjUserOut {
	var msgObjs []mysqlModel.MsgObjUserOut
	dao.DB.Table("msgobjs").Select("msgobjs.id,from_username,to_username,msg_time,`read`,last_msg,message_type,head_pic,user_id,nickname").
		Joins("INNER JOIN users on users.username = msgobjs.to_username").
		Where("from_username=?", fromUsername).Order("id").
		Find(&msgObjs)
	for i, m := 0, len(msgObjs); i < m; i++ {
		msgObjs[i].HeadPic = formatUtil.GetPicHeaderBody(host, msgObjs[i].HeadPic)
	}
	return msgObjs
}

// DeleteMsg 删除指定id的记录
func (m *MsgObjService) DeleteMsg(msgId int, username string) error {
	var msgObj mysqlModel.Msgobj
	err := dao.DB.Where("id=?", msgId).Where("from_username=?", username).Find(&msgObj).Error
	if err != nil {
		return err
	}
	dao.DB.Delete(&msgObj)
	return nil
}

// BatchRead msgobjs 内最新消息设置已读
func (m *MsgObjService) BatchRead(FromUsername, ToUsername string) {
	//var msgObj mysqlModel.Msgobj
	//err := dao.DB.Where("from_username=?", FromUsername).Where("to_username=?", ToUsername).Find(&msgObj).Error
	//if err != nil {
	//	return
	//}
	//msgObj.Read = 1
	//dao.DB.Save(&msgObj)
	err := dao.DB.Table("msgobjs").Model(&mysqlModel.Msgobj{}).
		Where("from_username=?", FromUsername).Where("to_username=?", ToUsername).
		UpdateColumn(map[string]interface{}{"read": 1}).Error
	if err != nil {
		log.Println("ws BatchRead failed,err:", err)
		return
	}
}
