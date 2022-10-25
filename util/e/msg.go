package e

var codeMsg = map[Code]string{
	WebsocketSuccessMessage: "解析content内容",
	WebsocketSuccess:        "发送消息，请求历史记录操作成功",
	WebsocketEnd:            "请求历史记录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对霍夫信息里先回答成功",
	WebsocketLimit:          "请求收到限制",
}
