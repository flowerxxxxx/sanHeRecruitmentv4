package websocketBiz

import (
	"encoding/json"
	"sanHeRecruitment/biz/nsqBiz"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service/mysqlService"
	"sanHeRecruitment/wechatPubAcc"
)

//互斥锁

var wsModule *WsModule
var us *mysqlService.UserService

// SysMsgPusher 定制化系统消息推送
func SysMsgPusher(toUsername, sendMsg string) {
	client := &websocketModel.Client{
		ID:           wsModule.CreateID("employeeSystem", toUsername), //1->2
		SendID:       wsModule.CreateID(toUsername, "employeeSystem"), //2->1
		FromUsername: "employeeSystem",
		ToUsername:   toUsername,
	}
	msgCast := &websocketModel.Broadcast{
		Client:  client,
		Message: []byte(sendMsg), //发送过来的消息
		Type:    0,
	}
	broadcastMessage := string(msgCast.Message)
	message := &websocketModel.BroadcastMsg{
		Message:     broadcastMessage,
		MessageType: msgCast.Type,
	}
	message2, _ := json.Marshal(message)
	SendId := msgCast.Client.SendID //2->1
	flag := false                   //默认对方是不在线的

	//去用户管理里寻找sendid，如果有则证明是该被发送者是在线的，如果没有则不在线
	conn, ok := websocketModel.ReadManClient(SendId)
	if ok {
		if conn.SendOpen {
			conn.Send <- message2
			flag = true
		}
	}

	id := msgCast.Client.ID //1->2
	if flag {
		newInsert := websocketModel.InsertMysql{
			Id:           id,
			Content:      message.Message,
			Read:         1,
			Expire:       int64(config.MsgExpiredTime),
			MessageType:  msgCast.Type,
			FromUsername: msgCast.Client.FromUsername,
			ToUsername:   msgCast.Client.ToUsername,
		}
		go nsqBiz.ChatProducer(newInsert)
	} else {
		newInsert := websocketModel.InsertMysql{
			Id:           id,
			Content:      message.Message,
			Read:         0,
			Expire:       int64(config.MsgExpiredTime),
			MessageType:  msgCast.Type,
			FromUsername: msgCast.Client.FromUsername,
			ToUsername:   msgCast.Client.ToUsername,
		}
		//建立goroutine向不在线但登录的用户推送消息提醒
		go func(fromUser, content string, messageType int) {
			//异步pusher
			//fmt.Println("异步消息推送")
			//用户推送在线查找flag，检索不到在线即通过公众号推送
			findFlag := 0
			cliMap, ok := websocketModel.ReadRecManClient(msgCast.Client.ToUsername)
			if ok {
				fromUserNickname := userSer.QueryUserNickByUsername(fromUser)
				publishMsg := websocketModel.PublishMsg{
					FromUser:       fromUserNickname,
					MessageContent: content,
					MessageType:    messageType,
				}
				pubMsg, _ := json.Marshal(publishMsg)
				if websocketModel.ReceiveMsgManager.Clients[cliMap.ID].SendOpen {
					websocketModel.ReceiveMsgManager.Clients[cliMap.ID].Send <- pubMsg
					findFlag = 1
				}
			}
			if findFlag == 0 {
				//微信公众号推送
				//TODO 暂时关闭公众号推送
				if messageType == 1 {
					content = "[图片]"
				}
				fromUserNickname := userSer.QueryUserNickByUsername(fromUser)
				wechatPubAcc.ConversationMessagePush(msgCast.Client.ToUsername, fromUserNickname, content)
			}
		}(msgCast.Client.FromUsername, message.Message, message.MessageType)
		go nsqBiz.ChatProducer(newInsert)

	}
}

func FromMainToPush() {
	for {
		select {
		case MSG := <-nsqBiz.FM.ToServiceMiddleContent:
			SysMsgPusher(MSG.ToUsername, MSG.MsgContent)
		}
	}
}

// MassSendMsg 系统定制化身份群发信息(ALL -1,usual 0,boss 1,service 2)
func MassSendMsg(sendRole int, msg string) (err error) {
	userColony, err := us.QueryUserColony(sendRole)
	if err != nil {
		return err
	}
	go func() {
		for _, item := range userColony {
			//go func(toUsername string) {
			//	SysMsgPusher(toUsername, msg)
			//}(item.Username)
			SysMsgPusher(item.Username, msg)
		}
	}()
	return
}

// InitSystemAdminer 初始化系统人员
func InitSystemAdminer() error {
	_, err := us.GetUserInfo("employeeSystem", "")
	if err != nil {
		err = us.CreateSysCaller()
		if err != nil {
			return err
		}
	}
	_, err = us.GetUserInfo(config.AdminUsername, "")
	if err != nil {
		err = us.CreateSysAdmin()
		if err != nil {
			return err
		}
	}
	//_, err = us.GetUserInfo(config.ProducerUsername)
	//if err != nil {
	//	err = us.CreateSysAdminDeveloper()
	//	if err != nil {
	//		return err
	//	}
	//}
	return nil
}
