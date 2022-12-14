package websocket_biz

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sanHeRecruitment/biz/nsq_biz"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service"
	"sanHeRecruitment/util/e"
	"sanHeRecruitment/wechatPubAcc"
)

var userSer *service.UserService

type WsModule struct {
}

// CreateID 创建trainer的userid
func (ws *WsModule) CreateID(uid, toUid string) string {
	return uid + "->" + toUid // 1 -> 2
}

// AddMsgPusher 根据count计数器添加msger
func (ws *WsModule) AddMsgPusher(stuNum string, cliRec *websocketModel.ClientRecMsg) {
	msgCount := websocketModel.ReceiveMsgManager.ClientCount[stuNum]
	if msgCount == 0 {
		websocketModel.EditRecManCliCount(stuNum, 1)
	} else {
		msgCount++
		websocketModel.EditRecManCliCount(stuNum, msgCount)
	}
	websocketModel.ReceiveMsgManager.ClientsRWM.Lock()
	websocketModel.ReceiveMsgManager.Clients[stuNum] = cliRec
	websocketModel.ReceiveMsgManager.ClientsRWM.Unlock()
}

// WsStart websocket管道通信监听
func (ws *WsModule) WsStart() {
	for {
		//TODO "---监听管道通信---"
		//log.Println("---监听管道通信---")
		select {
		case conn := <-websocketModel.Manager.Register:
			// TODO new conn printer
			//fmt.Printf("有新连接：%v\n", conn.ID)
			//fmt.Println(&Manager.Register)

			//注册前检测上一次是否存在
			//if connLastLive, ok := websocketModel.ReadManClient(conn.ID); ok {
			//	fmt.Println("check get last live")
			//	close(connLastLive.Send)
			//} else {
			//	fmt.Println("check get last failed")
			//}
			//开始注册
			websocketModel.Manager.ClientsRWM.Lock()
			websocketModel.Manager.Clients[conn.ID] = conn //将该连接放到用户管理上
			websocketModel.ManagerCliCountIncr(conn.ID)
			websocketModel.Manager.ClientsRWM.Unlock()
			replyMsg := &websocketModel.ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "服务器连接成功",
			}
			msg, _ := json.Marshal(replyMsg)
			conn.SocketMutex.Lock()
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			conn.SocketMutex.Unlock()
		case conn := <-websocketModel.Manager.Unregister:
			//fmt.Printf("连接中断%s\n", conn.ID)
			//_, ok := websocketModel.Manager.Clients[conn.ID]
			cliCount, ok := websocketModel.ReadCliCount(conn.ID)
			if !ok {
				log.Println("[May fatal error]Manager Unregister logic maybe err ")
				return
			}
			if cliCount > 1 {
				websocketModel.ManagerCliCountCutOne(conn.ID)
			} else if cliCount == 1 {
				if _, ok := websocketModel.ReadManClient(conn.ID); ok {
					replyMsg := &websocketModel.ReplyMsg{
						Code:    e.WebsocketEnd,
						Content: "连接中断",
					}
					msg, _ := json.Marshal(replyMsg)
					conn.SocketMutex.Lock()
					_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
					conn.SocketMutex.Unlock()
					//if _, okSend := <-conn.Send; okSend {
					//	fmt.Println("chec")
					//	close(conn.Send)
					//}
					websocketModel.Manager.ClientsRWM.Lock()
					delete(websocketModel.Manager.Clients, conn.ID)
					websocketModel.Manager.ClientsRWM.Unlock()
					//fmt.Println("Manager del succ")
				}
				websocketModel.DelManagerCliCount(conn.ID)
			}
		case broadcast := <-websocketModel.Manager.Broadcast: //1->2
			//start := time.Now()
			broadcastMessage := string(broadcast.Message)
			message := &websocketModel.BroadcastMsg{
				Message:     broadcastMessage,
				MessageType: broadcast.Type,
			}
			message2, _ := json.Marshal(message)
			SendId := broadcast.Client.SendID //2->1
			flag := false                     //默认对方是不在线的
			//去用户管理里寻找sendid，如果有则证明是该被发送者是在线的，如果没有则不在线
			conn, ok := websocketModel.ReadManClient(SendId)
			if ok {
				select {
				case conn.Send <- message2:
					flag = true
				default:
					websocketModel.Manager.ClientsRWM.Lock()
					delete(websocketModel.Manager.Clients, conn.ID)
					websocketModel.Manager.ClientsRWM.Unlock()
					close(conn.Send)
				}
			}

			id := broadcast.Client.ID //1->2
			if flag {
				// TODO WS online Printer
				//fmt.Println("对方在线")
				replyMsg := &websocketModel.ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				//RWMux.Lock()
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//RWMux.Unlock()
				newInsert := websocketModel.InsertMysql{
					Id:           id,
					Content:      message.Message,
					Read:         1,
					Expire:       int64(config.MsgExpiredTime),
					MessageType:  broadcast.Type,
					FromUsername: broadcast.Client.FromUsername,
					ToUsername:   broadcast.Client.ToUsername,
				}
				go nsq_biz.Producer(newInsert)
			} else {
				//fmt.Println("对方不在线")
				replyMsg := &websocketModel.ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线回答",
				}
				msg, _ := json.Marshal(replyMsg)
				//RWMux.Lock()
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				//RWMux.Unlock()
				newInsert := websocketModel.InsertMysql{
					Id:           id,
					Content:      message.Message,
					Read:         0,
					Expire:       int64(config.MsgExpiredTime),
					MessageType:  broadcast.Type,
					FromUsername: broadcast.Client.FromUsername,
					ToUsername:   broadcast.Client.ToUsername,
				}
				//建立goroutine向不在线但登录的用户推送消息提醒
				go func(fromUser, content string, messageType int) {
					// TODO 异步推送打印
					//fmt.Println("异步消息推送")
					//用户推送在线查找flag，检索不到在线即通过公众号推送
					findFlag := 0
					cliMap, ok := websocketModel.ReadRecManClient(broadcast.Client.ToUsername)
					if ok {
						fromUserNickname := userSer.QueryUserNickByUsername(fromUser)
						publishMsg := websocketModel.PublishMsg{
							FromUser:       fromUserNickname,
							MessageContent: content,
							MessageType:    messageType,
						}
						pubMsg, _ := json.Marshal(publishMsg)
						PushMsg := struct {
							HeartBeat int
						}{1}
						checkOnlineMsg, _ := json.Marshal(PushMsg)
						//检查clients是否存在
						err := websocketModel.ReceiveMsgManager.Clients[cliMap.ID].
							Socket.WriteMessage(websocket.TextMessage, checkOnlineMsg)
						if err != nil {
							//log.Println("Socket.WriteMessage failed,errInfo:", cliMap.ID)
							return
						}
						select {
						case websocketModel.ReceiveMsgManager.Clients[cliMap.ID].Send <- pubMsg:
							findFlag = 1
						default:
							log.Println("Clients[", cliMap.ID, "].Send has err closed")
							return
						}
					}
					if findFlag == 0 {
						//微信公众号推送
						if messageType == 1 {
							content = "[图片]"
						}
						fromUserNickname := userSer.QueryUserNickByUsername(fromUser)
						wechatPubAcc.ConversationMessagePush(broadcast.Client.ToUsername, fromUserNickname, content)
					}
				}(broadcast.Client.FromUsername, message.Message, message.MessageType)
				go nsq_biz.Producer(newInsert)
			}
		}
	}
}

// RecMsgStart msg publisher init start
func (ws *WsModule) RecMsgStart() {
	for {
		select {
		case conn := <-websocketModel.ReceiveMsgManager.Unregister:
			close(conn.Send)
			websocketModel.DelRecManCli(conn.ID)
			websocketModel.DelRecCliCountCli(conn.ID)
		}
	}
}
