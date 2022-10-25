package mysqlModel

import "time"

// Upgrade 升级身份数据库结构体
type Upgrade struct {
	ID          int `json:"id" gorm:"primary_key"`
	CompanyId   int `json:"company_id"`
	TargetLevel int `json:"target_level"`
	//资质审核 0未审核 1通过 2未通过
	Qualification int       `json:"qualification"`
	FromUsername  string    `json:"from_username"`
	ApplyTime     time.Time `json:"apply_time"`
	//公司是否先前存在 0不存在 1存在
	CompanyExist int `json:"company_exist"`
	//管理员端是否显示（根据用户身份升级完成度） 0不显示 1显示
	Show   int   `json:"show"`
	TimeId int64 `json:"time_id"` //时间凭证
}

type WaitingUpgrade struct {
	ID          int `json:"id" gorm:"primary_key"`
	CompanyId   int `json:"company_id"`
	TargetLevel int `json:"target_level"`
	//资质审核 0未审核 1通过 2未通过
	Qualification int       `json:"qualification"`
	FromUsername  string    `json:"from_username"`
	ApplyTime     time.Time `json:"-"`
	//公司是否先前存在 0不存在 1存在
	CompanyExist int    `json:"company_exist"`
	CompanyName  string `json:"company_name"`
	Name         string `json:"name"`
	TimeId       int64  `json:"time_id"` //时间凭证
}

type WaitingUpgradeOut struct {
	WaitingUpgrade
	ApplyTimeOut string `json:"apply_time"`
}

type WaitingUpgradeXls struct {
	WaitingUpgradeOut
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Description string `json:"description"`
}
