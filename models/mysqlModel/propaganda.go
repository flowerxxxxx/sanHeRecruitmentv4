package mysqlModel

import "sanHeRecruitment/util/timeUtil"

type Propaganda struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Uploader   string           `json:"uploader"`
	Url        string           `json:"url"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Title      string           `json:"title"`
	//流媒体类型 0图片 1视频
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// PropagandaOutHead 首页 焦点
type PropagandaOutHead struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Uploader   string           `json:"-"`
	Url        string           `json:"url"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	//流媒体类型 0图片 1视频
	Type  int    `json:"type"`
	Load  bool   `json:"load"`
	Title string `json:"title"`
}

// PropagandaOutContent 详细 焦点
type PropagandaOutContent struct {
	PropagandaOutHead
	Content string `json:"content"`
}