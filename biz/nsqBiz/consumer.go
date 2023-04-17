package nsqBiz

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/config"
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

// Consumer nsq消费者函数
func Consumer() {
	//会话自插入
	c, err := nsq.NewConsumer(config.NsqConfig.ConsumerTopic, config.NsqConfig.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	c.AddHandler(&ConsumerT{})                                             // 添加消息处理
	if err := c.ConnectToNSQD(config.NsqConfig.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

}
