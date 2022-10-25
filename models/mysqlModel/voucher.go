package mysqlModel

// Voucher 简历凭证存储表
type Voucher struct {
	ID         int    `json:"id" gorm:"primary_key"`
	Username   string `json:"username"`
	VoucherPic string `json:"voucher_pic"`
	CompanyId  int    `json:"company_id"`
	TimeId     int64  `json:"time_id"`
}
