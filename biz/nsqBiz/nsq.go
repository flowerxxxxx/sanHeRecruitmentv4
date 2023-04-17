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

var RM = ReceiveMessage{
	InsertContent: make(chan *websocketModel.InsertMysql),
}

var chatProducer *nsq.Producer
var mainToServiceProducer *nsq.Producer //单独抽象 -> ws服务器广播

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

func InitProducer() (err error) {
	chatProducer, err = nsq.NewProducer(config.NsqConfig.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}

	mainToServiceProducer, err = nsq.NewProducer(config.NsqMainToService.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}

	return err
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
