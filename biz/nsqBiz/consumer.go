package nsqBiz

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/models/websocketModel"
)

// ConsumerT nsq订阅消息
type ConsumerT struct{}

// HandleMessage nsq消费者处理函数
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	var insertMsg *websocketModel.InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	RM.InsertContent <- insertMsg
	return nil
}
