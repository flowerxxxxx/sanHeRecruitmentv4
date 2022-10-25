package mysqlModel

import "sanHeRecruitment/util/timeUtil"

type Connection struct {
	Id         int              `json:"id" gorm:"primary_key"`
	DesPerson  string           `json:"des_person"`
	Connect    string           `json:"connect"`
	UploadTime *timeUtil.MyTime `json:"upload_time"`
	UpdateTime *timeUtil.MyTime `json:"update_time"`
	Type       string           `json:"type"`
}

type ConnectionOut struct {
	Id        int    `json:"id" gorm:"primary_key"`
	DesPerson string `json:"des_person"`
	Connect   string `json:"connect"`
	Type      string `json:"type"`
}
