package mysqlModel

import "sanHeRecruitment/util/timeUtil"

// VipShow 会员风采
type VipShow struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Cover      string           `json:"cover"`
	CreateTime *timeUtil.MyTime `json:"create_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Clicks     int              `json:"clicks"`
	Content    string           `json:"content"`
	Publisher  string           `json:"publisher"`
}

// VipShowOut 会员风采输出
type VipShowOut struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Cover      string           `json:"cover"`
	CreateTime *timeUtil.MyTime `json:"create_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Clicks     int              `json:"clicks"`
	Content    string           `json:"content"`
}
