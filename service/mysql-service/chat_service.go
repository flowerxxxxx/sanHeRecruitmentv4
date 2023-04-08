package mysql_service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/models/websocketModel"
	"sync"
)

type ChatService struct {
}

// QELM 添加头像url
type QELM struct {
	websocketModel.Trainer
	ToUserHeadPic string `json:"to_user_head_pic"`
	//ToUsername    string `json:"to_username"`
	ToUserId       int    `json:"to_user_id"`
	ToUserNickname string `json:"to_user_nickname"`
}

// QueryEveryLastMsg 查询每个会话的最新消息
func (c *ChatService) QueryEveryLastMsg(personList []mysqlModel.Msgobj) []QELM {
	var MsgList []QELM
	wg := new(sync.WaitGroup)
	wgInt := len(personList)
	wg.Add(wgInt)
	for _, person := range personList {
		go func(person mysqlModel.Msgobj) {
			personChat := QELM{}
			//var personChat QELM
			//dao.DB.Raw("? UNION ? ORDER BY id DESC LIMIT 1",
			//	dao.DB.Table("trainers").Where("from_username = ?", person.ToUsername).Where("to_username =?", person.FromUsername),
			//	dao.DB.Table("trainers").Where("from_username = ?", person.FromUsername).Where("to_username =?", person.ToUsername),
			//).Scan(&personChat)
			sql := "SELECT id,content,start_time,`read`,message_type,from_username FROM trainers WHERE from_username = ? AND to_username = ? " +
				"UNION" +
				" SELECT id,content,start_time,`read`,message_type,from_username FROM trainers WHERE from_username = ? AND to_username =? " +
				" ORDER BY id DESC LIMIT 1"
			err := dao.DB.Raw(sql, person.ToUsername, person.FromUsername, person.FromUsername, person.ToUsername).Scan(&personChat).Error
			if err != nil {
				wg.Done()
				return
			}
			personInfo := mysqlModel.User{}
			personInfoSql := "select user_id,head_pic,nickname from users where username = ? "
			dao.DB.Raw(personInfoSql, person.ToUsername).Scan(&personInfo)
			personChat.Id = person.Id
			personChat.ToUserHeadPic = personInfo.Head_pic
			personChat.ToUserId = personInfo.User_id
			personChat.ToUserNickname = personInfo.Nickname
			MsgList = append(MsgList, personChat)
			wg.Done()
		}(person)
	}
	wg.Wait()
	return MsgList
}

// BatchRead 批量设置已读
func (c *ChatService) BatchRead(fromUsername, toUsername string) {
	dao.DB.Table("trainers").Where("from_username=?", fromUsername).Where("to_username=?", toUsername).Where("`read`=?", 0).Update(map[string]interface{}{"read": 1})
}

// FindExpiredChatPic 查找过期聊天图片
func (c *ChatService) FindExpiredChatPic(nowUnixStr string) []*websocketModel.Trainer {
	var Pics []*websocketModel.Trainer
	sql := "select content from trainers where message_type = 1 and end_time < ?"
	dao.DB.Raw(sql, nowUnixStr).Scan(&Pics)
	return Pics
}

// DeleteExpiredMsg 删除过期聊天记录
func (c *ChatService) DeleteExpiredMsg(nowUnixStr string) (err error) {
	sql := "DELETE FROM trainers WHERE end_time < ?"
	err = dao.DB.Exec(sql, nowUnixStr).Error
	return err
}
