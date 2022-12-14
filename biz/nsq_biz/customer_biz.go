package nsq_biz

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/models/websocketModel"
)

// ConsumerT producer
type ConsumerT struct{}

// HandleMessage nsq消费者处理函数
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	//receive msg printer
	//fmt.Println("---获取消息---")
	//fmt.Println(string(msg.Body))
	var insertMsg *websocketModel.InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	//fmt.Println(insertMsg)
	//nsqInsertMux.Lock()
	RM.InsertContent <- insertMsg
	//nsqInsertMux.Unlock()
	return nil
}

// ConsumerMsgPusher NsqConfigMsgPusher
type ConsumerMsgPusher struct{}

// HandleMessage nsq消费者处理函数
func (*ConsumerMsgPusher) HandleMessage(msg *nsq.Message) error {
	var insertMsg *websocketModel.InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	RM.InsertContent <- insertMsg
	return nil
}
