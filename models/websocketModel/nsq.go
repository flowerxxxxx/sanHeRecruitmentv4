package websocketModel

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"time"
)

type ReceiveMessage struct {
	InsertContent chan *InsertMysql
}

var RM = ReceiveMessage{
	InsertContent: make(chan *InsertMysql),
}

type InsertMysql struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	Read         int    `json:"read"`
	Expire       int64  `json:"expire"`
	MessageType  int    `json:"message_type"`
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
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

// Producer nsq发布消息
func Producer(msg InsertMysql) {
	data, _ := json.Marshal(msg)
	if err := producer.Publish(config.NsqConfig.ProducerTopic, data); err != nil { // 发布消息
		log.Println("[fatal Info]nsq publish err :", err)
	}
}

// ConsumerT nsq订阅消息
type ConsumerT struct{}

// HandleMessage nsq消费者处理函数
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	//TODO receive msg printer
	//fmt.Println("---获取消息---")
	//fmt.Println(string(msg.Body))
	var insertMsg *InsertMysql
	_ = json.Unmarshal(msg.Body, &insertMsg)
	//fmt.Println(insertMsg)
	//nsqInsertMux.Lock()
	RM.InsertContent <- insertMsg
	//nsqInsertMux.Unlock()
	return nil
}

// Consumer nsq消费者函数
func Consumer() {
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
}

// ReceiveToInsert 处理nsq消费者接收函数
func ReceiveToInsert() {
	for {
		// TODO ---接收消息并处理---
		//log.Println("---接收消息并处理---")
		select {
		case MSG := <-RM.InsertContent:
			CheckAndAddMsgForML(MSG)
			_ = InsertMsg2(MSG)
			//TODO "处理完成"
			//log.Println("处理完成")
		}
	}
}

// CheckAndAddMsgForML 消息列表在发消息的检测和添加机制
func CheckAndAddMsgForML(msg *InsertMysql) {
	msgObj1 := mysqlModel.Msgobj{}
	err := dao.DB.Where("from_username=?", msg.FromUsername).Where("to_username=?", msg.ToUsername).Find(&msgObj1).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData mysqlModel.Msgobj
		msgData.FromUsername = msg.FromUsername
		msgData.ToUsername = msg.ToUsername
		msgData.MsgTime = time.Now().Unix()
		msgData.LastMsg = msg.Content
		msgData.Read = 1
		msgData.MessageType = msg.MessageType
		dao.DB.Save(&msgData)
	} else {
		msgObj1.MsgTime = time.Now().Unix()
		msgObj1.LastMsg = msg.Content
		msgObj1.Read = 1
		msgObj1.MessageType = msg.MessageType
		dao.DB.Save(&msgObj1)
	}
	msgObj2 := mysqlModel.Msgobj{}
	err = dao.DB.Where("from_username=?", msg.ToUsername).Where("to_username=?", msg.FromUsername).Find(&msgObj2).Error
	if err != nil {
		//检索不到相关信息的情况则创建
		var msgData mysqlModel.Msgobj
		msgData.FromUsername = msg.ToUsername
		msgData.ToUsername = msg.FromUsername
		msgData.MsgTime = time.Now().Unix()
		msgData.LastMsg = msg.Content
		msgData.Read = msg.Read
		msgData.MessageType = msg.MessageType
		dao.DB.Save(&msgData)
	} else {
		msgObj2.MsgTime = time.Now().Unix()
		msgObj2.LastMsg = msg.Content
		msgObj2.Read = msg.Read
		msgObj2.MessageType = msg.MessageType
		dao.DB.Save(&msgObj2)
	}
}
