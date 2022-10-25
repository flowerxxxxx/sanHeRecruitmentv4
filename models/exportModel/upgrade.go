package exportModel

type UpgradeExporter struct {
	CompanyName string `json:"company_name"` //公司名
	Phone       string `json:"phone"`        //联系方式
	Address     string `json:"address"`      //地址
	Description string `json:"description"`  //公司描述
	NewRegister string `json:"new_register"` //公司是否新注册
	Done        string `json:"done"`         //通过情况
	Applicant   string `json:"applicant"`    //申请人
	ApplyTime   string `json:"apply_time"`   //申请时间
}
