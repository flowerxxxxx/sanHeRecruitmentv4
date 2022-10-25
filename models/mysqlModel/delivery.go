package mysqlModel

import "time"

type Delivery struct {
	ID           int    `json:"id" gorm:"primary_key"`
	BossId       int    `json:"boss_id"`       //boss
	ArtId        int    `json:"art_id"`        //招聘id
	FromUsername string `json:"from_username"` //投递人
	// 成功资质 0未审核 1通过 2未通过
	Qualification int `json:"qualification"`
	// boss已读 0未读 1已读
	Read         int       `json:"read"`
	DeliveryTime time.Time `json:"-"`
}

// BossDeliveries boss查询简历投递结构体
type BossDeliveries struct {
	ID              int       `json:"id" gorm:"primary_key"`
	DeliveryTime    time.Time `json:"-"`
	DeliveryTimeOut string    `json:"delivery_time"`
	Name            string    `json:"name"`
	Title           string    `json:"title"`
	Qualification   int       `json:"qualification"`
	FromUsername    string    `json:"from_username"` //投递人
}
