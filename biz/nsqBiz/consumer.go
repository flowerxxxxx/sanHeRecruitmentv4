package nsqBiz

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/config"
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
	var insertMsg *websocketModel.InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	RM.InsertContent <- insertMsg
	return nil
}

// Consumer nsq消费者函数
func Consumer() {
	//会话自插入
	ct, err := nsq.NewConsumer(config.NsqConfig.ConsumerTopic, config.NsqConfig.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	ct.AddHandler(&consumerT{})                                             // 添加消息处理
	if err := ct.ConnectToNSQD(config.NsqConfig.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

	//web -> service
	tc, err := nsq.NewConsumer(config.NsqConfig.ConsumerTopic, config.NsqConfig.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	tc.AddHandler(&consumerT{})                                             // 添加消息处理
	if err := tc.ConnectToNSQD(config.NsqConfig.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

}
