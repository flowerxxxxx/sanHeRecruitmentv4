package timeTaskModule

import (
	"log"
	"os"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"strings"
)

type TimeTaskModule struct {
}

// DeleteExpiredMsgPic 删除过期消息对应的图片
func (t *TimeTaskModule) DeleteExpiredMsgPic(picData []*websocketModel.Trainer) {
	for _, pic := range picData {
		picUrl := pic.Content
		pos := strings.Index(picUrl, "/uploadPic")
		finalPicUrl := config.PicSaverPath + picUrl[pos+10:]
		err := os.Remove(finalPicUrl)
		if err != nil {
			log.Println("expire msg file remove Error!,err:", err)
		}
	}
}
