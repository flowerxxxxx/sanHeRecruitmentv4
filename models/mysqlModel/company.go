package mysqlModel

import "time"

// Company 公司的数据库结构体
type Company struct {
	ComId        int       `json:"com_id" gorm:"primary_key"`
	PicUrl       string    `json:"pic_url"`
	CompanyName  string    `json:"company_name"`
	Description  string    `json:"description"`
	ScaleTag     string    `json:"scale_tag"`    //融资
	PersonScale  string    `json:"person_scale"` //人员规模
	Address      string    `json:"address"`
	UpdateTime   time.Time `json:"-"`
	UpdatePerson string    `json:"update_person"`
	ComLevel     int       `json:"com_level"`
	Vip          int       `json:"vip"`
	Phone        string    `json:"phone"`
	ComStatus    int       `json:"com_status"`
	ScaleLevel   int       `json:"scale_level"`
}

type CompanyName struct {
	ComId       int    `json:"com_id" gorm:"primary_key"`
	CompanyName string `json:"company_name"`
	ComLevel    int    `json:"com_level"`
}

// CompanyBasicInfo Company 公司的数据库结构体
type CompanyBasicInfo struct {
	ComId       int    `json:"com_id" gorm:"primary_key"`
	PicUrl      string `json:"pic_url"`
	CompanyName string `json:"company_name"`
	Description string `json:"description"`  //简称
	ScaleTag    string `json:"scale_tag"`    //融资
	PersonScale string `json:"person_scale"` //人员规模
	Address     string `json:"address"`
	ComLevel    int    `json:"com_level"`
	Vip         int    `json:"vip"`
	Phone       string `json:"phone"`
}

type CompanyOut struct {
	Company
	UpdateTimeOut string `json:"update_time"`
	Name          string `json:"name"`
}

type CompaniesTotal struct {
	TotalNums int `json:"total_nums"`
	Companies int `json:"companies"`
	Servers   int `json:"servers"`
}
