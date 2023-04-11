package mysqlService

import (
	"sanHeRecruitment/config"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
)

type CountService struct {
}

const pageSize = config.PageSize
const webPageSize = config.WebPageSize

// GetFuzzyQueryJobsTP 模糊获取工作信息总页数
func (c *CountService) GetFuzzyQueryJobsTP(fuzzyName, queryType string, ifAdmin int) (totalPage int) {
	var Total mysqlModel.Count
	jobTotalSql := dao.DB.Table("articles").Select("COUNT(*) as total_num").
		Where("`status` = ?", 1).Where("`show` = ?", 1).
		Where("articles.art_type = ?", queryType).
		Where("LOCATE(?,articles.title) > 0", fuzzyName)
	jobTotalSql.Find(&Total)
	if ifAdmin == 0 {
		return getTotalPage(Total)
	} else {
		return getTotalPageWeb(Total)
	}
}

// GetJobsTotalPage 获取对应标签的工作的总页数
func (c *CountService) GetJobsTotalPage(jobId int, region, labelType, pageType string) (totalPage int) {
	var Total mysqlModel.Count
	jobTotalSql := dao.DB.Table("articles").Select("COUNT(*) as total_num").
		Where("`status` = ?", 1).Where("`show` = ?", 1)
	if region != "" {
		jobTotalSql = jobTotalSql.Where("region = ?", region)
	}
	if jobId != 0 {
		jobTotalSql = jobTotalSql.Where("career_job_id = ?", jobId)
	} else {
		jobTotalSql = jobTotalSql.
			Joins("INNER JOIN labels on articles.career_job_id = labels.id").
			Where("labels.type = ?", labelType)
	}
	jobTotalSql.Find(&Total)
	if pageType == "micro" {
		return getTotalPage(Total)
	} else {
		return getTotalPageWeb(Total)
	}
}

func (c *CountService) GetEmployeesTotalPage(minAge, maxAge, gender, minDegreeLevel, jobLabel string) (totalPage int) {
	var Total mysqlModel.Count
	genderSql := ""
	if gender != "" {
		genderSql = " AND gender =? "
	}
	eduSql := ""
	if minDegreeLevel != "" {
		eduSql = " AND educations.degree_level >= ? "
	}
	if genderSql == "" && eduSql == "" {
		sql := "SELECT COUNT(*) AS total_num FROM (SELECT COUNT(*) as total_num  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ? GROUP BY educations.username) COUNT "
		dao.DB.Raw(sql, jobLabel, minAge, maxAge).Scan(&Total)
	} else if genderSql != "" && eduSql == "" {
		sql := "SELECT COUNT(*) AS total_num FROM (SELECT COUNT(*) as total_num  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ?" + genderSql +
			" GROUP BY educations.username) COUNT "
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, gender).Scan(&Total)
	} else if genderSql == "" && eduSql != "" {
		sql := "SELECT COUNT(*) AS total_num FROM (SELECT COUNT(*) as total_num  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ?" + eduSql +
			" GROUP BY educations.username) COUNT "
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, minDegreeLevel).Scan(&Total)
	} else {
		sql := "SELECT COUNT(*) AS total_num FROM (SELECT COUNT(*) as total_num  FROM `educations` INNER JOIN users ON users.username = educations.username  " +
			"WHERE user_level = 0 AND intended_position = ? AND age BETWEEN ? AND ?" + genderSql + eduSql +
			" GROUP BY educations.username) COUNT "
		dao.DB.Raw(sql, jobLabel, minAge, maxAge, gender, minDegreeLevel).Scan(&Total)
	}
	totalPage = getTotalPage(Total)
	return
}

func (c *CountService) GetDeliveryTotalPage(username string, qualification, read int) (totalPage int) {
	var Total mysqlModel.Count
	sql := "SELECT COUNT(*) AS total_num FROM `deliveries`" +
		"WHERE from_username = ? AND qualification = ? AND `read` = ?"
	dao.DB.Raw(sql, username, qualification, read).Scan(&Total)
	totalPage = getTotalPage(Total)
	return
}

func (c *CountService) GetAllDeliveryTotalPage(username string) (totalPage int) {
	var Total mysqlModel.Count
	sql := "SELECT COUNT(*) AS total_num FROM `deliveries`" +
		"WHERE from_username = ?"
	dao.DB.Raw(sql, username).Scan(&Total)
	totalPage = getTotalPage(Total)
	return
}

func (c *CountService) GetQueryDeliveriesTotalPage(bossId, qualification int) (totalPage int) {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("deliveries").Select("COUNT(*) AS total_num").
		Where("deliveries.boss_id = ?", bossId)
	if qualification != -1 {
		totalQ = totalQ.Where("qualification = ?", qualification)
	}
	totalQ.Scan(&Total)
	totalPage = getTotalPage(Total)
	return
}

// GetInvitedJobsTotalPage 查看邀请投递总页数
func (c *CountService) GetInvitedJobsTotalPage(desUsername string) (totalPage int) {
	var Total mysqlModel.Count
	dao.DB.Table("invitations").Select("COUNT(*) AS total_num").
		Where("des_username = ?", desUsername).Find(&Total)
	totalPage = getTotalPage(Total)
	return
}

// GetWaitingTotalPage 查看待审核总页数
func (c *CountService) GetWaitingTotalPage(waitType string, comId int) (totalPage int) {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("articles").Select("COUNT(*) AS total_num").
		//Joins("INNER JOIN labels on articles.career_job_id = labels.id").
		Where("status = ?", 0)
	if waitType != "all" {
		totalQ = totalQ.Where("articles.job_label = ?", waitType)
	}
	if comId != 0 {
		totalQ = totalQ.Where("company_id = ?", comId)
	}
	totalQ.Find(&Total)
	totalPage = getTotalPageWeb(Total)
	return
}

func (c *CountService) GetSonWaitingTotalPage(sonIdSli []int, comId int) (totalPage int) {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("articles").Select("COUNT(*) AS total_num").
		Where("status = ?", 0)
	if comId != 0 {
		totalQ = totalQ.Where("company_id = ?", comId)
	}
	totalQ.Where("career_job_id in (?)", sonIdSli).Find(&Total)
	totalPage = getTotalPageWeb(Total)
	return
}

func (c *CountService) NoticesInfosTP() int {
	var Total mysqlModel.Count
	dao.DB.Table("notices").Select("COUNT(*) AS total_num").
		Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) VipShowInfosTP() int {
	var Total mysqlModel.Count
	dao.DB.Table("vip_shows").Select("COUNT(*) AS total_num").
		Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) GetDailyInfoTotalPage(queryLabel, artType, queryType, queryDate string) (totalPage int) {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("dailysavers").Select("COUNT(*) AS total_num").
		Joins("INNER JOIN articles on articles.art_id = dailysavers.art_id "+
			"INNER JOIN labels on labels.id = articles.career_job_id").Where("date = ?", queryDate)
	if queryLabel == "allLabel" {
		totalQ = totalQ.Where("labels.type = ?", artType)
	} else {
		totalQ = totalQ.Where("labels.label = ?", queryLabel)
	}
	if queryType == "day_view" {
		totalQ = totalQ.Where("day_view > ?", 0)
	} else {
		totalQ = totalQ.Where("day_delivery > ?", 0)
	}
	totalQ.Order("dailysavers.id desc").Find(&Total)
	totalPage = getTotalPageWeb(Total)
	return
}

// GetWaitingUpgradeTP 获取待审核的升级总页数
func (c *CountService) GetWaitingUpgradeTP(qualification, targetLevel int) int {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("upgrades").Select("COUNT(*) AS total_num").
		Where("upgrades.show = ?", 1).
		Where("target_level = ?", targetLevel)
	if qualification != -1 {
		totalQ = totalQ.Where("qualification = ?", qualification)
	}
	totalQ.Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) MassPubHistoryTP() int {
	var Total mysqlModel.Count
	dao.DB.Table("mass_sends").Select("COUNT(*) AS total_num").
		Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) AdminInfosTP() int {
	var Total mysqlModel.Count
	dao.DB.Table("users").Select("COUNT(*) AS total_num").
		Where("role = ?", "admin").Where("user_level >= ?", 100).Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) UserCollectTP(username, classification string) int {
	var Total mysqlModel.Count
	dao.DB.Table("collections").Select("COUNT(*) AS total_num").
		Where("username = ?", username).Where("classification = ?", classification).Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) UserCollectTotal(username string) int {
	var Total mysqlModel.Count
	dao.DB.Table("collections").Select("COUNT(*) AS total_num").
		Where("username = ?", username).Find(&Total)
	return Total.TotalNum
}

func (c *CountService) CompanyPubTotal(JobReqType, getType string, companyId, status int) int {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("articles").Select("COUNT(*) AS total_num").
		Joins("INNER JOIN companies on companies.com_id = articles.company_id "+
			"INNER JOIN labels ON articles.career_job_id = labels.id").
		Where("articles.company_id = ?", companyId).
		Where("`status`=?", status).
		Where("`show` = ?", 1)
	if JobReqType != "" {
		totalQ = totalQ.Where("`job_label` = ?", JobReqType)
	} else {
		if getType != "" {
			totalQ = totalQ.Where("labels.`type` = ?", getType)
		}
	}
	totalQ.Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) CompanyPubTotalAdmin(companyName, JobReqType, getType string, status int) int {
	var Total mysqlModel.Count
	totalQ := dao.DB.Table("articles").Select("COUNT(*) AS total_num").
		Joins("INNER JOIN companies on companies.com_id = articles.company_id "+
			"INNER JOIN labels ON articles.career_job_id = labels.id").
		Where("company_name = ?", companyName).
		Where("`status`=?", status).
		Where("`show` = ?", 1)
	if JobReqType != "" {
		totalQ = totalQ.Where("`job_label` = ?", JobReqType)
	} else {
		if getType != "" {
			totalQ = totalQ.Where("labels.`type` = ?", getType)
		}
	}
	totalQ.Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) BossJobInfoTP(bossId, queryType string, desStatus, desShow int) int {
	var Total mysqlModel.Count
	// 链式查询，用于条件筛选
	sqlQuery := dao.DB.Table("articles").Select("COUNT(*) AS total_num").
		Joins("INNER JOIN labels on labels.id = articles.career_job_id")
	if desStatus != -1 {
		sqlQuery = sqlQuery.Where("`status`=?", desStatus)
	}
	if desShow != -1 {
		sqlQuery = sqlQuery.Where("`show` = ?", desShow)
	}
	if queryType != "all" {
		sqlQuery = sqlQuery.Where("labels.type = ?", queryType)
	}
	sqlQuery = sqlQuery.Where("boss_id=?", bossId).Find(&Total)
	return getTotalPage(Total)
}

func (c *CountService) QueryAllCompaniesTP(comLevel int) int {
	var Total mysqlModel.Count
	dao.DB.Table("companies").Select("COUNT(*) AS total_num").
		Where("com_level = ?", comLevel).Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) QueryAllFuzzyCompaniesTP(fuzzyComName, companyLevel string, desStatus int) int {
	var Total mysqlModel.Count
	queryQ := dao.DB.Table("companies").Select("COUNT(*) AS total_num").
		Where("LOCATE(?,companies.company_name) > 0", fuzzyComName)
	if companyLevel != "0" {
		queryQ = queryQ.Where("com_level = ?", companyLevel)
	}
	if desStatus != -1 {
		queryQ = queryQ.Where("com_status = ?", desStatus)
	}
	queryQ.Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) GetDockInfosTP(comId int) int {
	var Total mysqlModel.Count
	dao.DB.Table("docks").Select("COUNT(*) AS total_num").
		Where("com_id = ?", comId).Find(&Total)
	return getTotalPageWeb(Total)
}

func (c *CountService) GetComUsersTP(comId int) int {
	var Total mysqlModel.Count
	dao.DB.Table("users").Select("COUNT(*) AS total_num").
		Where("company_id = ?", comId).Find(&Total)
	return getTotalPageWeb(Total)
}

func getTotalPage(Total mysqlModel.Count) (totalPage int) {
	if Total.TotalNum%pageSize != 0 {
		totalPage = Total.TotalNum/pageSize + 1
	} else {
		totalPage = Total.TotalNum / pageSize
	}
	if totalPage == 0 {
		totalPage = 1
	}
	return
}

func getTotalPageWeb(Total mysqlModel.Count) (totalPage int) {
	if Total.TotalNum%webPageSize != 0 {
		totalPage = Total.TotalNum/webPageSize + 1
	} else {
		totalPage = Total.TotalNum / webPageSize
	}
	if totalPage == 0 {
		totalPage = 1
	}
	return
}

// 自由设置pageSize和获取页数
func getTotalPageFree(TotalNum, pageSize int) (totalPage int) {
	if TotalNum%pageSize != 0 {
		totalPage = TotalNum/pageSize + 1
	} else {
		totalPage = TotalNum / pageSize
	}
	if totalPage == 0 {
		totalPage = 1
	}
	return
}
