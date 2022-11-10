package mysqlModel

import "sanHeRecruitment/util/timeUtil"

type Notice struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Content    string           `json:"content"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Uploader   string           `json:"uploader"`
	Title      string           `json:"title"`
	Recommend  int              `json:"recommend"`
}

type NoticeOutHead struct {
	Id         int              `json:"id" gorm:"primary_key"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Title      string           `json:"title"`
	Recommend  int              `json:"recommend"`
}

type NoticeOutContent struct {
	NoticeOutHead
	Content string `json:"content"`
}
