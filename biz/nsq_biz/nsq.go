package nsq_biz

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service"
)

type ReceiveMessage struct {
	InsertContent chan *websocketModel.InsertMysql
}

var RM = ReceiveMessage{
	InsertContent: make(chan *websocketModel.InsertMysql),
}

var producer *nsq.Producer            //会话消息 -> 消息队列 -> 数据库
var wsBroadcastProducer *nsq.Producer //单独抽象 -> ws服务器广播

//
////var nsqInsertMux sync.Mutex

func InitProducer() (err error) {
	producer, err = nsq.NewProducer(config.NsqConfig.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}

	wsBroadcastProducer, err = nsq.NewProducer(config.NsqConfigWsBroadcast.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}

	return err
}

// InitConsumer nsq消费者函数
func InitConsumer() {
	//producer 消费者
	c_producer, err := nsq.NewConsumer(config.NsqConfig.ConsumerTopic, config.NsqConfig.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	c_producer.AddHandler(&ConsumerT{})                                             // 添加消息处理
	if err := c_producer.ConnectToNSQD(config.NsqConfig.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

	//NsqConfigMsgPusher消费者
	cmsg, err := nsq.NewConsumer(config.NsqConfigMsgPusher.ConsumerTopic, config.NsqConfigMsgPusher.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	cmsg.AddHandler(&ConsumerT{})                                                      // 添加消息处理
	if err := cmsg.ConnectToNSQD(config.NsqConfigMsgPusher.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

}

// ReceiveToInsert 处理nsq消费者接收函数
func ReceiveToInsert() {
	ns := service.NsqService{}
	for {
		select {
		case MSG := <-RM.InsertContent:
			ns.CheckAndAddMsgForML(MSG)
			_ = ns.InsertMsg2(MSG)
		}
	}
}
