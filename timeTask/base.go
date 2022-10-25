package timeTask

import (
	"github.com/robfig/cron"
	"log"
)

// InitTimer 定义gin-定时任务
func InitTimer() {
	log.Println("Timer Starting...")
	c := cron.New()
	//删除即时通讯中的过期信息 执行周期：每个小时
	errD := c.AddFunc("0 0 0/1 * * ? ", deleteExpiredInfo)
	if errD != nil {
		log.Println("cron deleteExpiredInfo work err,err:", errD)
	}

	//系统对数据库和即时通讯消息进行备份 执行周期：每日凌晨2点
	errB := c.AddFunc("0 0 2 * * ? ", backer)
	if errB != nil {
		log.Println("cron backer work err,err:", errB)
	}

	c.Start()
}
