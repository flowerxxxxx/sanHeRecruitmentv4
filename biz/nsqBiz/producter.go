package nsqBiz

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
)

// ChatProducer ChatProducer内部发布nsq消息
func ChatProducer(msg websocketModel.InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := chatProducer.Publish(config.NsqConfig.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq publish err :", err)
	}
}

// ToServiceProducer Web -> 会话x
func ToServiceProducer(msg websocketModel.InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := chatProducer.Publish(config.NsqMainToService.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq publish err :", err)
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
