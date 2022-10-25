package service

import (
	"errors"
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"strconv"
)

var (
	ErrorNotExisted    = errors.New("用户名不存在")
	ErrorPasswordWrong = errors.New("密码错误")
)

type UserService struct {
}

//func NewUserService() *UserService {
//	return &UserService{}
//}

func (us *UserService) TotalCount() (mysqlModel.UsersTotal, error) {
	var ut mysqlModel.UsersTotal
	err := dao.DB.Table("users").Select("COUNT( CASE WHEN user_level < 5 THEN 1 ELSE NULL END ) AS total_nums," +
		"COUNT( CASE WHEN user_level = 0 THEN 1 ELSE NULL END ) AS ordinary_nums," +
		"COUNT( CASE WHEN user_level = 1 THEN 1 ELSE NULL END ) AS com_nums," +
		"COUNT( CASE WHEN user_level = 2 THEN 1 ELSE NULL END ) AS ser_nums").
		Find(&ut).Error
	return ut, err
}

// Login 用户登录
func (us *UserService) Login(username, password string) error {
	var user mysqlModel.User
	err := dao.DB.Where("username=?", username).Find(&user).Error
	if err != nil {
		err2 := dao.DB.Where("email=?", username).Find(&user).Error
		if err2 != nil {
			return ErrorNotExisted
		}
	}
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return nil
}

func (us *UserService) WechatLogin(username string) error {
	var user mysqlModel.User
	err := dao.DB.Where("username=?", username).Find(&user).Error
	if err != nil {
		return ErrorNotExisted
	}
	return nil
}

func (us *UserService) WechatRegister(username, nickname, headPic string) (err error) {
	var newUser mysqlModel.User
	newUser.Username = username
	newUser.Role = "user"
	newUser.Nickname = nickname
	newUser.Head_pic = headPic
	newUser.UserLevel = 0
	newUser.IdentyPin = 0
	newUser.CompanyID = 0
	newUser.MakeTime = timeUtil.GetMyTimeNowPtr()
	err = dao.DB.Save(&newUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserInfo(username, host string) (userInfo mysqlModel.User, err error) {
	var user mysqlModel.User
	err = dao.DB.Where("username=?", username).Find(&user).Error
	if err != nil {
		return mysqlModel.User{}, ErrorNotExisted
	}
	user.Head_pic = formatUtil.GetPicHeaderBody(host, user.Head_pic)
	return user, nil
}

// QueryUserResumeInfo 搜索用户简历信息（含个人信息）
func (us *UserService) QueryUserResumeInfo(username, host string) (mysqlModel.ResumeUserInfo, error) {
	var ResumeInfo mysqlModel.ResumeUserInfo
	sql := "SELECT user_id,head_pic,username,email,nickname,gender,age,`name`,phone,personal_skill,personal_resume,project_experience,intended_position FROM `users` where username = ?"
	err := dao.DB.Raw(sql, username).Scan(&ResumeInfo).Error
	if err != nil {
		return ResumeInfo, ErrorNotExisted
	}
	ResumeInfo.Head_pic = formatUtil.GetPicHeaderBody(host, ResumeInfo.Head_pic)
	return ResumeInfo, nil
}

// QueryUserBasicInfo 检索用户个人基础信息
func (us *UserService) QueryUserBasicInfo(username, host string) (mysqlModel.BasicUserInfo, error) {
	var BasicInfo mysqlModel.BasicUserInfo
	sql := "SELECT user_id,username,email,head_pic,nickname,gender,age,`name`,phone FROM `users` where username =?"
	err := dao.DB.Raw(sql, username).Scan(&BasicInfo).Error
	if err != nil {
		return BasicInfo, ErrorNotExisted
	}
	BasicInfo.Head_pic = formatUtil.GetPicHeaderBody(host, BasicInfo.Head_pic)
	return BasicInfo, nil
}

// QueryUserInfoByUserId 通过user_id查找用户信息
func (us *UserService) QueryUserInfoByUserId(userId string) (userInfo mysqlModel.User, err error) {
	userIdInt, _ := strconv.Atoi(userId)
	err = dao.DB.Where("user_id=?", userIdInt).Find(&userInfo).Error
	return
}

// QueryBasicUserInfoByUserId 通过user_id查找用户信息
func (us *UserService) QueryBasicUserInfoByUserId(userId int) (userInfo mysqlModel.BasicUserInfo, err error) {
	err = dao.DB.Table("users").
		Select("user_id,username,email,head_pic,nickname,gender,age,`name`,phone").
		Where("user_id=?", userId).Find(&userInfo).Error
	return
}

// ModifyHeadPic 修改头像
func (us *UserService) ModifyHeadPic(username, headPicAddr, host string) (err error) {
	userInfo, err := us.GetUserInfo(username, host)
	if err != nil {
		return err
	}
	userInfo.Head_pic = headPicAddr
	err = dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return
}

// ModifyResume 修改简历
func (us *UserService) ModifyResume(username, resume, host string) (err error) {
	userInfo, err := us.GetUserInfo(username, host)
	if err != nil {
		return err
	}
	userInfo.Resume = resume
	err = dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return
}

// ModifyPersonalInfo 修改个人信息
func (us *UserService) ModifyPersonalInfo(username, email, gender, intended_position, name, nickname string, age int) {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.Email = email
	userInfo.Gender = gender
	userInfo.Intended_position = intended_position
	userInfo.Name = name
	userInfo.Nickname = nickname
	userInfo.Age = age
	dao.DB.Save(&userInfo)
}

func (us *UserService) QueryUserNickByUsername(username string) (nickname string) {
	var userInfo mysqlModel.User
	sql := "select nickname from users where username = '" + username + "'"
	dao.DB.Raw(sql).Scan(&userInfo)
	return userInfo.Nickname
}

// ModifyBasicPersonalInfo 修改个人基础信息
func (us *UserService) ModifyBasicPersonalInfo(username, email, name, gender, nickname, intendedPosition, phoneNumber string, age int) error {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.Email = email
	userInfo.Gender = gender
	userInfo.Name = name
	userInfo.Nickname = nickname
	userInfo.Age = age
	userInfo.Intended_position = intendedPosition
	userInfo.Phone = phoneNumber
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

// ModifyPersonalResumePE 修改个人简历-个人经历部分
func (us *UserService) ModifyPersonalResumePE(username, personal_experience string) error {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.ProjectExperience = personal_experience
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

// ModifyPersonalResumePS 修改个人简历-个人技能部分
func (us *UserService) ModifyPersonalResumePS(username, personalSkill string) error {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.PersonalSkill = personalSkill
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

// ModifyPersonalResumePR  修改个人简历-个人简介部分
func (us *UserService) ModifyPersonalResumePR(username, personal_resume string) error {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.PersonalResume = personal_resume
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) ModifyPersonalPresident(username, president string) error {
	var userInfo mysqlModel.User
	dao.DB.Where("username=?", username).Find(&userInfo)
	userInfo.President = president
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) QueryUserLevel(username string) (userLevel mysqlModel.UserLevel) {
	sql := "SELECT identy_pin,user_level FROM `users` WHERE username = ?"
	dao.DB.Raw(sql, username).Scan(&userLevel)
	return
}

func (us *UserService) ModifyUserIdentityPin(username string, changeLevel int) (err error) {
	var userInfo mysqlModel.User
	err = dao.DB.Where("username=?", username).Find(&userInfo).Error
	if err != nil {
		return
	}
	userInfo.IdentyPin = changeLevel
	err = dao.DB.Save(&userInfo).Error
	if err != nil {
		return
	}
	return
}

// ResetUserLevel 重置用户等级
func (us *UserService) ResetUserLevel(username string) (err error) {
	var userInfo mysqlModel.User
	err = dao.DB.Where("username=?", username).Find(&userInfo).Error
	if err != nil {
		return
	}
	userInfo.IdentyPin = 0
	userInfo.UserLevel = 0
	userInfo.CompanyID = 0
	err = dao.DB.Save(&userInfo).Error
	if err != nil {
		return
	}
	return
}

func (us *UserService) ModifyPersonalInfoByUpgrade(username string, companyId, targetLevel int) (err error) {
	var userInfo mysqlModel.User
	err = dao.DB.Where("username=?", username).Find(&userInfo).Error
	if err != nil {
		return err
	}
	userInfo.CompanyID = companyId
	userInfo.IdentyPin = 0
	userInfo.UserLevel = targetLevel
	err = dao.DB.Save(&userInfo).Error
	if err != nil {
		return err
	}
	return
}

// CreateSysCaller 创建sysCaller
func (us *UserService) CreateSysCaller() error {
	var userInfo mysqlModel.User
	userInfo.Username = "employeeSystem"
	userInfo.Nickname = "系统通知"
	userInfo.Role = "system"
	userInfo.UserLevel = 99
	userInfo.IdentyPin = 99
	userInfo.Name = "系统通知"
	userInfo.MakeTime = timeUtil.GetMyTimeNowPtr()
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return errors.New("系统通知初始化失败")
	}
	return err
}

func (us *UserService) CreateSysAdmin() error {
	var userInfo mysqlModel.User
	userInfo.Username = config.AdminUsername
	userInfo.Nickname = "三河系统管理员"
	userInfo.Password = sqlUtil.GenMD5Password(config.AdminPassword)
	userInfo.Role = "admin"
	userInfo.UserLevel = 100
	userInfo.IdentyPin = 100
	userInfo.Name = "管理员"
	userInfo.MakeTime = timeUtil.GetMyTimeNowPtr()
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return errors.New("系统管理员初始化失败")
	}
	return err
}

func (us *UserService) CreateSysAdminDeveloper() error {
	var userInfo mysqlModel.User
	userInfo.Username = config.ProducerUsername
	userInfo.Nickname = "三河系统开发者"
	userInfo.Password = sqlUtil.GenMD5Password(config.ProducerPassword)
	userInfo.Role = "admin"
	userInfo.UserLevel = 99
	userInfo.IdentyPin = 99
	userInfo.Name = "三河系统开发者"
	err := dao.DB.Save(&userInfo).Error
	if err != nil {
		return errors.New("系统管理员（开发者）初始化失败")
	}
	return err
}

func (us *UserService) QueryEmployees(jobLabel string) (userInfos []mysqlModel.BasicUserInfo) {
	selectSql := "SELECT user_id,age,username,nickname,name,email,phone,gender FROM `users`"
	conditionSql := "where user_level = 0 and intended_position = ? "
	finalSql := selectSql + conditionSql
	dao.DB.Raw(finalSql, jobLabel).Scan(&userInfos)
	return
}

// QueryUserColony 根据等级查询用户群体（username）
func (us *UserService) QueryUserColony(desUserLevel int) (userColony []mysqlModel.UserName, err error) {
	if desUserLevel == -1 {
		err = dao.DB.Table("users").Where("user_level < ?", 99).Find(&userColony).Error
	} else {
		err = dao.DB.Table("users").Where("user_level =?", desUserLevel).Find(&userColony).Error
	}
	return
}

func (us *UserService) AdminModifyPwd(username, newPWD string) (err error) {
	var userInfo mysqlModel.User
	err = dao.DB.Table("users").Where("username = ?", username).Find(&userInfo).Error
	if err != nil {
		return NoRecord
	}
	userInfo.Password = newPWD
	err = dao.DB.Table("users").Save(&userInfo).Error
	return
}

func (us *UserService) AddAdminer(username, password, nickname string) (err error) {
	var adminInfo mysqlModel.User
	adminInfo.Username = username
	adminInfo.UserLevel = 101
	adminInfo.IdentyPin = 101
	adminInfo.Role = "admin"
	adminInfo.Password = password
	adminInfo.Nickname = nickname
	adminInfo.Name = nickname
	err = dao.DB.Table("users").Save(&adminInfo).Error
	return
}

func (us *UserService) DeleteAdminer(username string) (err error) {
	err = dao.DB.Table("users").Where("username = ?", username).Delete(&mysqlModel.User{}).Error
	return
}

func (us *UserService) GetAdminerInfos(pageNum int) (userInfo []mysqlModel.User, err error) {
	var users []mysqlModel.User
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	err = dao.DB.Table("users").
		Where("role =?", "admin").Where("user_level >= ?", 100).Limit(webPageSize).Offset(sqlPage).Find(&users).Error
	if err != nil {
		return users, ErrorNotExisted
	}
	return users, nil
}

func (us *UserService) GetComUsers(pageNum, comId int, host string) []*mysqlModel.BasicUserInfo {
	var users []*mysqlModel.BasicUserInfo
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	dao.DB.Table("users").Where("company_id = ?", comId).Limit(webPageSize).
		Offset(sqlPage).Find(&users)
	for i, m := 0, len(users); i < m; i++ {
		users[i].Head_pic = formatUtil.GetPicHeaderBody(host, users[i].Head_pic)
	}
	return users
}

// KillUserComIdentity 重置用户的身份
func (us *UserService) KillUserComIdentity(userId int) (err error) {
	err = dao.DB.Model(&mysqlModel.User{}).Where("user_id = ?", userId).
		UpdateColumns(mysqlModel.User{UserLevel: 0, CompanyID: 0, IdentyPin: 0}).Error
	return
}

//  根据用户收藏的articleid查询收藏的需求
// sql:select articles.art_id,articles.content,articles.creat_time,articles.`view`,users.username,users.vip,users.head_pic,labels.type FROM user_article_lables ual LEFT JOIN users on users.user_id = ual.user_id LEFT JOIN articles on articles.art_id = ual.art_id LEFT JOIN labels on labels.lab_id = ual.lab_id WHERE articles.art_id in('1','2','3') LIMIT 0,10

// 通过收藏查询收藏的内容
//SELECT articles.art_id,articles.content,articles.creat_time,articles.`view`,users.username,users.vip,users.head_pic,labels.type FROM user_article_lables ual LEFT JOIN users ON users.user_id = ual.user_id LEFT JOIN articles ON articles.art_id = ual.art_id LEFT JOIN labels ON labels.lab_id = ual.lab_id WHERE articles.art_id IN(select art_id from collections where username = '20062111') ORDER BY art_id DESC LIMIT 0,10
