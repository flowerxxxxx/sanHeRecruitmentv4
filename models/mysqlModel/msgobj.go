package mysqlModel

type Msgobj struct {
	Id           int    `json:"id" gorm:"primary_key"`
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
	MsgTime      int64  `json:"start_time"`
	Read         int    `json:"read"`
	LastMsg      string `json:"content"`
	MessageType  int    `json:"message_type"`
}

type MsgObjUserOut struct {
	Msgobj
	HeadPic  string `json:"to_user_head_pic"`
	UserId   int    `json:"to_user_id"`
	Nickname string `json:"to_user_nickname"`
	Online   int    `json:"online"`
}
