package nsqBiz

import (
	"encoding/json"
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
