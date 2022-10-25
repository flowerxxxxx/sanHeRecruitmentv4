package osUtil

import (
	"log"
	"os"
	"sanHeRecruitment/config"
	"strings"
)

func DeleteFile(TotalUrl string) {
	pos := strings.Index(TotalUrl, "/uploadPic")
	finalPicUrl := config.PicSaverPath + TotalUrl[pos+10:]
	go func() {
		err := os.Remove(finalPicUrl)
		if err != nil {
			log.Println("DeleteStreamOrPic file remove Error!")
			log.Printf("%s", err)
		}
	}()
}
