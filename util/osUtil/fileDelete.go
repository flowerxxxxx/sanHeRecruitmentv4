package osUtil

import (
	"log"
	"os"
	"sanHeRecruitment/config"
	"strings"
)

func DeleteFile(TotalUrl string) {
	pos := strings.Index(TotalUrl, "/uploadPic")
	if pos == -1 {
		return
	}
	finalPicUrl := config.PicSaverPath + TotalUrl[pos+10:]
	err := os.Remove(finalPicUrl)
	if err != nil {
		log.Println("DeleteStreamOrPic file remove Error!")
		log.Printf("%s", err)
	}
}
