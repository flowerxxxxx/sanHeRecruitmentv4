package timeTask

import (
	"log"
	timeTaskModule2 "sanHeRecruitment/biz/timeTaskBiz"
	"sanHeRecruitment/service/mysqlService"
	"strconv"
	"time"
)

var chatService *mysqlService.ChatService
var timeTaskModule *timeTaskModule2.TimeTaskModule

// 删除过期消息，消息有效时长，三个月
// 删除数据库对应的本地图片
func deleteExpiredInfo() {
	nowUnix := time.Now().Unix()
	nowUnixStr := strconv.Itoa(int(nowUnix))
	PicData := chatService.FindExpiredChatPic(nowUnixStr)
	if len(PicData) != 0 {
		timeTaskModule.DeleteExpiredMsgPic(PicData)
	}
	err := chatService.DeleteExpiredMsg(nowUnixStr)
	if err != nil {
		log.Println("timeTask delete MSG error:", err)
	}
}
