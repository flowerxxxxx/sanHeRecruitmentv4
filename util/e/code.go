package e

import "errors"

type Code int

const (
	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004 //在线应答
	WebsocketOfflineReply   = 50005 //不在线应答
	WebsocketLimit          = 50006
	WebsocketHistoryMsg     = 50007 //历史消息
	WebsocketUpdate         = 50008 //ws更新
)

func (c Code) Msg() string {
	return codeMsg[c]
}

var (
	RedisNoVal   = errors.New("redis get no val")
	UnMarshalErr = errors.New("json unmarshall err")
)
