package mysqlService

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"strconv"
)

type EducationService struct {
}

// QueryEmployeeEduBasic 根据查找人力资源
func (es *EducationService) QueryEmployeeEduBasic(users []mysqlModel.UserName) (userEduBasic []mysqlModel.EmployEduInfo, err error) {
	var usersSLi []string
	for _, item := range users {
		usersSLi = append(usersSLi, item.Username)
	}
	userSql := sqlUtil.DealSliToSql(usersSLi)
	selectSql := "select users.username,school,major,degree,age,gender,`name` from educations as a INNER JOIN users ON users.username = a.username"
	conditionSql := " where users.username in " + userSql + " AND degree_level = (select max(degree_level) from educations where a.username=username)"
	finalSql := selectSql + conditionSql
	err = dao.DB.Raw(finalSql).Scan(&userEduBasic).Error
	return
}

// ScreeningResumes 筛选简历
func (es *EducationService) ScreeningResumes(minAge, maxAge, gender, minDegreeLevel, jobLabel string, pageNum int) (userInfos []mysqlModel.UserName) {
	pageNumStr := strconv.Itoa((pageNum - 1) * 10)
	genderSql := ""
	if gender != "" {
		genderSql = " AND gender =? "
	}
	eduSql := ""
	if minDegreeLevel != "" {
		eduSql = " AND educations.degree_level >= ? "
	}
	if genderSql == "" && eduSql == "" {
		sql := "SELECT educations.username  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ?  GROUP BY username LIMIT ?,10"
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, pageNumStr).Scan(&userInfos)
	} else if genderSql != "" && eduSql == "" {
		sql := "SELECT educations.username  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ? " + genderSql +
			" GROUP BY username LIMIT ?,10"
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, gender, pageNumStr).Scan(&userInfos)
	} else if genderSql == "" && eduSql != "" {
		sql := "SELECT educations.username  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ? " + eduSql +
			" GROUP BY username LIMIT ?,10"
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, minDegreeLevel, pageNumStr).Scan(&userInfos)
	} else {
		sql := "SELECT educations.username  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ? " + genderSql + eduSql +
			" GROUP BY username LIMIT ?,10"
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, gender, minDegreeLevel, pageNumStr).Scan(&userInfos)
	}
	return
}

// QueryPersonalEdu 查询个人教育经历
func (es *EducationService) QueryPersonalEdu(username string) []mysqlModel.Education {
	var PersonalEduInfo []mysqlModel.Education
	sql := "select * from educations where username = ?"
	dao.DB.Raw(sql, username).Scan(&PersonalEduInfo)
	return PersonalEduInfo
}

// AddPersonalEdu 添加个人教育经历
func (es *EducationService) AddPersonalEdu(username, school, major, start_time, end_time, degree string) error {
	var PersonalEduInfo mysqlModel.Education
	PersonalEduInfo.Username = username
	PersonalEduInfo.School = school
	PersonalEduInfo.Major = major
	PersonalEduInfo.StartTime = start_time
	PersonalEduInfo.EndTime = end_time
	PersonalEduInfo.Degree = degree
	PersonalEduInfo.DegreeLevel = mysqlModel.DegreeWeight[degree]
	err := dao.DB.Save(&PersonalEduInfo).Error
	if err != nil {
		return err
	}
	return err
}

// ModifyPersonalResumeEdu ModifyPersonalResume 修改个人简历
func (es *EducationService) ModifyPersonalResumeEdu(id int, username, school, major, start_time, end_time, degree string) error {
	var PersonalEduInfo mysqlModel.Education
	dao.DB.Where("id=?", id).Where("username=?", username).Find(&PersonalEduInfo)
	PersonalEduInfo.School = school
	PersonalEduInfo.Major = major
	PersonalEduInfo.StartTime = start_time
	PersonalEduInfo.EndTime = end_time
	PersonalEduInfo.Degree = degree
	PersonalEduInfo.DegreeLevel = mysqlModel.DegreeWeight[degree]
	err := dao.DB.Save(&PersonalEduInfo).Error
	if err != nil {
		return err
	}
	return err
}

// DeletePersonalResumeEdu DeletePersonalResume 删除个人简历
func (es *EducationService) DeletePersonalResumeEdu(id int, username string) (err error) {
	err = dao.DB.Where("id=?", id).Where("username=?", username).Delete(&mysqlModel.Education{}).Error
	if err != nil {
		return err
	}
	return
}
