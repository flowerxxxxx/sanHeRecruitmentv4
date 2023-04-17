package nsqBiz

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/models/websocketModel"
)

// consumerT 内部会话consumer
type consumerT struct{}

// toServiceConsumer web -> service
type toServiceConsumer struct{}

// HandleMessage nsq消费者处理函数
func (*toServiceConsumer) HandleMessage(msg *nsq.Message) error {
	var insertMsg *websocketModel.InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	RM.InsertContent <- insertMsg
	return nil
}

// HandleMessage nsq消费者处理函数
func (*consumerT) HandleMessage(msg *nsq.Message) error {
	var mainMsg *websocketModel.ToServiceMiddle
	_ = json.Unmarshal(msg.Body, &mainMsg)
	FM.ToServiceMiddleContent <- mainMsg
	return nil
}
