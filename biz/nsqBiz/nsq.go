package nsqBiz

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"sanHeRecruitment/service/mysqlService"
)

type ReceiveMessage struct {
	InsertContent chan *websocketModel.InsertMysql
}

type FromMainMsg struct {
	ToServiceMiddleContent chan *websocketModel.ToServiceMiddle
}

var RM = ReceiveMessage{
	InsertContent: make(chan *websocketModel.InsertMysql),
}

var FM = FromMainMsg{
	ToServiceMiddleContent: make(chan *websocketModel.ToServiceMiddle),
}

var chatProducer *nsq.Producer

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
	tc, err := nsq.NewConsumer(config.NsqMainToService.ConsumerTopic, config.NsqMainToService.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	tc.AddHandler(&toServiceConsumer{})                                            // 添加消息处理
	if err := tc.ConnectToNSQD(config.NsqMainToService.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

}

func InitProducer() (err error) {
	chatProducer, err = nsq.NewProducer(config.NsqConfig.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}

	return err
}
