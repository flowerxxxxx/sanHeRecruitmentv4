package logUtil

import (
	"log"
	"os"
	"sanHeRecruitment/config"
)

func LogOutInit() {
	//file := "./" +"logs/systemLogOut"+ ".txt"
	file := config.ErrorLogAddr
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetPrefix("[logTool]")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Ltime)
	return
}
