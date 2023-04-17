package nsqBiz

import (
	"github.com/nsqio/go-nsq"
	"log"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service/mysqlService"
)

type ReceiveMessage struct {
	InsertContent chan *websocketModel.InsertMysql
}

var RM = ReceiveMessage{
	InsertContent: make(chan *websocketModel.InsertMysql),
}

var chatProducer *nsq.Producer

//
////var nsqInsertMux sync.Mutex

// ReceiveToInsert 处理nsq消费者接收函数
func ReceiveToInsert() {
	var ns mysqlService.NsqService
	for {
		select {
		case MSG := <-RM.InsertContent:
			ns.CheckAndAddMsgForML(MSG)
			insertErr := ns.InsertMsg(MSG)
			if insertErr != nil {
				log.Println("ReceiveToInsert InsertMsg failed,err:", insertErr)
			}
		}
	}
}
