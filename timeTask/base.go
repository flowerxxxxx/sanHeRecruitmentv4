package timeTask

import (
	"github.com/robfig/cron"
)

// InitTimer 定义gin-定时任务
func InitTimer() {
	c := cron.New()
	////删除即时通讯中的过期信息 执行周期：每个小时
	//errD := c.AddFunc("0 0 0/1 * * ? ", deleteExpiredInfo)
	//if errD != nil {
	//	log.Println("cron deleteExpiredInfo work failed,err:", errD)
	//}
	//

	c.Start()
}
