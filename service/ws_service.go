package service

import "sanHeRecruitment/dao"

type WsService struct {
}

// AddOnlineUserToRedis 将在线用户注册进redis缓存 1为会话在线 2为普通在线
func (ws *WsService) AddOnlineUserToRedis(username, coordinate string, AddType int) (err error) {
	adder := ""
	if AddType == 1 {
		adder = username + "_TalkOnline"
	} else if AddType == 2 {
		adder = username + "_PushOnline"
	}
	return dao.RedisDF.Do("SET", adder, coordinate).Err()

}
