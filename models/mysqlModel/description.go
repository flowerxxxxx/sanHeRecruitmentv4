package mysqlModel

import "sanHeRecruitment/util/timeUtil"

type Description struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Content    string           `json:"content"`
	Module     string           `json:"biz"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Uploader   string           `json:"uploader"`
}

type DescriptionOut struct {
	Id         int              `json:"id" gorm:"primary_key"`
	Content    string           `json:"content"`
	Module     string           `json:"module"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
}
