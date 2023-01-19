package timeTask

import (
	"github.com/robfig/cron"
	"log"
)

// InitTimer 定义gin-定时任务
func InitTimer() {
	c := cron.New()
	//删除即时通讯中的过期信息 执行周期：每个小时
	errD := c.AddFunc("0 0 0/1 * * ? ", deleteExpiredInfo)
	if errD != nil {
		log.Println("cron deleteExpiredInfo work failed,err:", errD)
	}

	//系统对数据库和即时通讯消息进行备份 执行周期：每日凌晨2点
	errB := c.AddFunc("0 0 2 1/5 * ? ", backer)
	if errB != nil {
		log.Println("cron backer work failed,err:", errB)
	}

	errExpireRemove := c.AddFunc("0 0 3 * * ? ", expireBackerRemove)
	if errExpireRemove != nil {
		log.Println("cron errExpireRemove work failed,err:", errB)
	}

	c.Start()
}
