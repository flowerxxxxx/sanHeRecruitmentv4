package mysqlService

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/models/websocketModel"
	"time"
)

type NsqService struct {
}

// CheckAndAddMsgForML 消息列表在发消息的检测和添加机制
func (ns *NsqService) CheckAndAddMsgForML(msg *websocketModel.InsertMysql) {
	msgObj1 := mysqlModel.Msgobj{}
	err := dao.DB.Where("from_username=?", msg.FromUsername).Where("to_username=?", msg.ToUsername).Find(&msgObj1).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData mysqlModel.Msgobj
		msgData.FromUsername = msg.FromUsername
		msgData.ToUsername = msg.ToUsername
		msgData.MsgTime = time.Now().Unix()
		msgData.LastMsg = msg.Content
		msgData.Read = 1
		msgData.MessageType = msg.MessageType
		dao.DB.Save(&msgData)
	} else {
		msgObj1.MsgTime = time.Now().Unix()
		msgObj1.LastMsg = msg.Content
		msgObj1.Read = 1
		msgObj1.MessageType = msg.MessageType
		dao.DB.Save(&msgObj1)
	}
	msgObj2 := mysqlModel.Msgobj{}
	err = dao.DB.Where("from_username=?", msg.ToUsername).Where("to_username=?", msg.FromUsername).Find(&msgObj2).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData mysqlModel.Msgobj
		msgData.FromUsername = msg.ToUsername
		msgData.ToUsername = msg.FromUsername
		msgData.MsgTime = time.Now().Unix()
		msgData.LastMsg = msg.Content
		msgData.Read = msg.Read
		msgData.MessageType = msg.MessageType
		dao.DB.Save(&msgData)
	} else {
		msgObj2.MsgTime = time.Now().Unix()
		msgObj2.LastMsg = msg.Content
		msgObj2.Read = msg.Read
		msgObj2.MessageType = msg.MessageType
		dao.DB.Save(&msgObj2)
	}
}

func (ns *NsqService) InsertMsg(msg *websocketModel.InsertMysql) (err error) {
	comment := websocketModel.Trainer{
		Userid:       msg.Id,
		Content:      msg.Content,
		Start_time:   time.Now().Unix(),
		End_time:     time.Now().Unix() + msg.Expire,
		Read:         msg.Read,
		Message_type: msg.MessageType,
		FromUsername: msg.FromUsername,
		ToUsername:   msg.ToUsername,
	}
	err = dao.DB.Save(&comment).Error
	if err != nil {
		return err
	}
	return
}
