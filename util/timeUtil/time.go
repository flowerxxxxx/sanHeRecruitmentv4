package timeUtil

import (
	"fmt"
	"strconv"
	"time"
)

// GetNowTimeFormat 获取当前时间的time.Time
func GetNowTimeFormat() time.Time {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	StartTime := timeStrToTime(nowTime)
	return StartTime
}

// TimeFormatToStr time.Time转str
func TimeFormatToStr(timeFormat time.Time) string {
	return timeFormat.Format("2006-01-02 15:04:05")
}

// TimeStrToTime 字符串格式时间转换为时间格式
func timeStrToTime(timer string) (TrueTime time.Time) {
	local, _ := time.ParseInLocation("2006-01-02 15:04:05", timer, time.Local)
	return local
}

func TimeStrToFormatStr(timer time.Time) (formatTimer time.Time) {
	newTimer := timer.Format("2006-01-02 15:04:05")
	timeStrToTime(newTimer)
	return
}

// MinutesToH 时间转换为小时（保存两位小数）
func MinutesToH(minutues int) (value float64) {
	a2 := float64(minutues)
	c := a2 / float64(60)
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", c), 64)
	return value
}
