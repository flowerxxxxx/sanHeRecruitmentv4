package nsq_biz

import (
	"encoding/json"
	"log"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
)

// Producer nsq发布消息
func Producer(msg websocketModel.InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := producer.Publish(config.NsqConfig.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq_biz publish err :", err)
	}
}

// WsBroadcastProducer ws服务器广播消息
func WsBroadcastProducer(msg websocketModel.InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := wsBroadcastProducer.Publish(config.NsqConfig.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq_biz publish err :", err)
	}
}
