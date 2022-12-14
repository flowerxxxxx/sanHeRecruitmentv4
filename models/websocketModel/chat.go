package websocketModel

//添加message_type 字段，1为文本消息，2为图片，考虑增加可发送的文档（3）

type Trainer struct {
	Id           int    `json:"id" gorm:"primary_key"`
	Userid       string `json:"-"`            //用户名
	Content      string `json:"content"`      // 内容
	Start_time   int64  `json:"start_time"`   // 创建时间
	End_time     int64  `json:"-"`            // 过期时间
	Read         int    `json:"read"`         // 已读
	Message_type int    `json:"message_type"` //消息类型
	FromUsername string `json:"-"`            //发送者
	ToUsername   string `json:"-"`            //接受者
}

//func InsertMsg(userid string, content string, read int, expire int64) (err error) {
//	comment := Trainer{
//		Userid:     userid,
//		Content:    content,
//		Start_time: time.Now().Unix(),
//		End_time:   time.Now().Unix() + expire,
//		Read:       read,
//	}
//	err = dao.DB.Save(&comment).Error
//	if err != nil {
//		return err
//	}
//	return
//}
//
//func InsertMsg2(msg *nsq_biz.InsertMysql) (err error) {
//	comment := Trainer{
//		Userid:       msg.Id,
//		Content:      msg.Content,
//		Start_time:   time.Now().Unix(),
//		End_time:     time.Now().Unix() + msg.Expire,
//		Read:         msg.Read,
//		Message_type: msg.MessageType,
//		FromUsername: msg.FromUsername,
//		ToUsername:   msg.ToUsername,
//	}
//	err = dao.DB.Save(&comment).Error
//	if err != nil {
//		return err
//	}
//	return
//}
