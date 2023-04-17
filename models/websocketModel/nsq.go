package websocketModel

type InsertMysql struct {
	Id           string `json:"id"`
	Content      string `json:"content"`
	Read         int    `json:"read"`
	Expire       int64  `json:"expire"`
	MessageType  int    `json:"message_type"`
	FromUsername string `json:"from_username"`
	ToUsername   string `json:"to_username"`
}

type ToServiceMiddle struct {
	ToUsername string
	MsgContent string
}
