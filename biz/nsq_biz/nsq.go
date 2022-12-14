package nsq_biz

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
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

var producer *nsq.Producer

//
////var nsqInsertMux sync.Mutex

func InitProducer() (err error) {
	producer, err = nsq.NewProducer(config.NsqConfig.ProducerAddr, nsq.NewConfig()) // 新建生产者
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		panic(any(err))
	}
	return err
}

// Producer nsq发布消息
func Producer(msg websocketModel.InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := producer.Publish(config.NsqConfig.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq_biz publish err :", err)
	}
}

// ConsumerT nsq订阅消息
type ConsumerT struct{}

// HandleMessage nsq消费者处理函数
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	//TODO receive msg printer
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

// InitConsumer nsq消费者函数
func InitConsumer() {
	// TODO "---消息队列---"
	//log.Println("---消息队列---")
	c, err := nsq.NewConsumer(config.NsqConfig.ConsumerTopic, config.NsqConfig.ConsumerChannel, nsq.NewConfig()) // 新建一个消费者
	if err != nil {
		panic(any(err))
	}
	c.AddHandler(&ConsumerT{})                                             // 添加消息处理
	if err := c.ConnectToNSQD(config.NsqConfig.ConsumerAddr); err != nil { // 建立连接
		panic(any(err))
	}

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
