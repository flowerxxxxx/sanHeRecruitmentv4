package mysqlModel

import "sanHeRecruitment/util/timeUtil"

// Dock 存储对接数据的数据库结构体
type Dock struct {
	ID       int              `json:"id" gorm:"primary_key"`
	ArtId    int              `json:"art_id"`
	ComId    int              `json:"com_id"`
	BossId   int              `json:"boss_id"`
	UserId   int              `json:"user_id"`
	DockTime *timeUtil.MyTime `json:"dock_time" time_format:"2006-01-02 15:04:05"`
	BossName string           `json:"boss_name"`
	UserName string           `json:"user_name"`
	PubTitle string           `json:"pub_title"`
}
