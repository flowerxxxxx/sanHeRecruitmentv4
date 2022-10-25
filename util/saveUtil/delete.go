package saveUtil

import (
	"log"
	"os"
	"sanHeRecruitment/config"
	"strings"
)

// DeletePicSaver 删除数据库存储在本地的图片
func DeletePicSaver(picUrl string) (err error) {
	pos := strings.Index(picUrl, "/uploadPic")
	finalPicUrl := config.PicSaverPath + picUrl[pos+10:]
	err = os.Remove(finalPicUrl)
	if err != nil {
		log.Println("file remove Error!")
		log.Printf("%s", err)
		return
	} else {
		log.Println("file remove OK!")
		return
	}
}
