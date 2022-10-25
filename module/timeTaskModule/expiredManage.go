package timeTaskModule

import (
	"log"
	"os"
	"sanHeRecruitment/config"
	"sanHeRecruitment/models/websocketModel"
	"strings"
	"sync"
)

var wg sync.WaitGroup

type TimeTaskModule struct {
}

// DeleteExpiredMsgPic 删除过期消息对应的图片
func (t *TimeTaskModule) DeleteExpiredMsgPic(picData []*websocketModel.Trainer) {
	wg.Add(len(picData))
	for _, pic := range picData {
		picUrl := pic.Content
		pos := strings.Index(picUrl, "/uploadPic")
		finalPicUrl := config.PicSaverPath + picUrl[pos+10:]
		go func(finalPicUrl string) {
			err := os.Remove(finalPicUrl)
			if err != nil {
				log.Println("expire msg file remove Error!")
				log.Printf("%s", err)
			} else {
				log.Println("10000 expire msg file remove OK!")
			}
			wg.Done()
		}(finalPicUrl)
	}
	wg.Wait()
}
