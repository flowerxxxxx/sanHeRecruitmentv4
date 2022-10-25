package mysqlModel

import "sanHeRecruitment/util/timeUtil"

// User users数据库结构体
type User struct {
	User_id  int    `json:"user_id" map:"user_id" gorm:"primary_key" `
	Username string `json:"username" map:"username"`
	Password string `json:"-"`
	Email    string `json:"email" map:"email"`
	Head_pic string `json:"head_pic" map:"head_pic"`
	Role     string `json:"role" map:"role"`
	//Company           int    `json:"company"` //判断是企业还是个人   0为个人 1为企业
	//Vip               int    `json:"vip"`     //企业判断是否是会员   个人判断是否是企业用户
	CompanyID int `json:"company_id" map:"company_id" `
	//身份码 普通用户0 企业用户1 服务机构2
	IdentyPin         int              `json:"identy_pin" map:"identy_pin"` //用户当前等级 0个人用户，1企业用户，2服务机构
	UserLevel         int              `json:"user_level" map:"user_level"` //用户等级 0个人用户，1企业用户，2服务机构
	Gender            string           `json:"gender" map:"gender"`
	Intended_position string           `json:"intended_position" map:"intended_position"` //意向岗位
	Resume            string           `json:"-" map:"-"`                                 //简历 文件版
	Age               int              `json:"age" map:"age"`
	Name              string           `json:"name" map:"name"`
	Nickname          string           `json:"nickname" map:"nickname"`
	Phone             string           `json:"phone" map:"phone"`
	PersonalSkill     string           `json:"personal_skill" map:"personal_skill"`
	ProjectExperience string           `json:"project_experience" map:"project_experience"`
	PersonalResume    string           `json:"personal_resume" map:"personal_resume"`
	President         string           `json:"president"` //职位
	MakeTime          *timeUtil.MyTime `json:"make_time"`
}

// BasicUserInfo 用户基础信息
type BasicUserInfo struct {
	User_id  int    `json:"user_id" map:"user_id" gorm:"primary_key" `
	Username string `json:"username" map:"username"`
	Email    string `json:"-" map:"email"`
	Gender   string `json:"gender" map:"gender"`
	Age      int    `json:"age" map:"age"`
	Name     string `json:"name" map:"name"`
	Nickname string `json:"nickname" map:"nickname"`
	Phone    string `json:"-" map:"phone"`
	Head_pic string `json:"head_pic" map:"head_pic"`
}

// ResumeUserInfo 简历信息含个人基础信息
type ResumeUserInfo struct {
	BasicUserInfo
	Intended_position string `json:"intended_position" map:"intended_position"` //意向岗位
	PersonalSkill     string `json:"personal_skill" map:"personal_skill"`
	ProjectExperience string `json:"project_experience" map:"project_experience"`
	PersonalResume    string `json:"personal_resume" map:"personal_resume"`
}

// UserLevel 用户等级
type UserLevel struct {
	IdentyPin int `json:"identy_pin" map:"identy_pin"` //用户当前等级 0个人用户，1企业用户，2服务机构
	UserLevel int `json:"user_level" map:"user_level"` //用户等级 0个人用户，1企业用户，2服务机构
}

// UserResumeInfo 人力资源显示信息
type UserResumeInfo struct {
	Username string `json:"username" map:"username"`
	School   string `json:"school"`
	Major    string `json:"major"`
	Degree   string `json:"degree"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
	Name     string `json:"name"`
}

// UserName username用户集结构体
type UserName struct {
	Username string `json:"username"`
}

// UsersTotal users表格总数统计
type UsersTotal struct {
	TotalNums    int `json:"total_nums"`
	OrdinaryNums int `json:"ordinary_nums"`
	ComNums      int `json:"com_nums"`
	SerNums      int `json:"ser_nums"`
}
