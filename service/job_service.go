package service

import (
	"fmt"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"strconv"
)

type JobService struct {
}

// GetRecommendJobs 根据权重，用户目标职位推荐信息
func (j *JobService) GetRecommendJobs(jobId, host string, pageNum int) []mysqlModel.UserComArticle {
	var userArt []mysqlModel.UserComArticle
	//pageNumStr := strconv.Itoa((pageNum - 1) * 10)
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	//sql := "SELECT art_id,title,job_label,tags,`view`,boss_id,nickname,head_pic,weight,company_name,salary_min,salary_max,region,person_scale FROM `articles` INNER JOIN users on articles.boss_id = users.user_id  INNER JOIN companies on companies.com_id = articles.company_id"
	//careerAndRegion := " where career_job_id = ? and `status` = 1 and `show` = 1 order by weight desc,art_id desc limit ?,10"
	//////高效率
	//////pageSql := userid >= (SELECT userid FROM `times` LIMIT" + pageNumStr + ", 1) LIMIT 10;"
	//finalSql := sql + careerAndRegion
	//fmt.Println(finalSql)
	//dao.DB.Debug().Raw(finalSql, jobId, pageNumStr).Scan(&userArt)
	dao.DB.Table("articles").Select("art_id,title,job_label,tags,`view`,boss_id,nickname,head_pic,weight,company_name,salary_min,salary_max,region,person_scale,articles.recommend").
		Joins("INNER JOIN users on articles.boss_id = users.user_id").
		Joins("INNER JOIN companies on companies.com_id = articles.company_id").
		Where("career_job_id = ?", jobId).Where("`status` = ?", 1).Where("`show` = ?", 1).
		Order("recommend desc,weight desc").Offset(sqlPage).Limit(10).Find(&userArt)
	for i, ul := 0, len(userArt); i < ul; i++ {
		userArt[i].TagsOut = sqlUtil.SqlStringToSli(userArt[i].Tags)
		userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
	}
	return userArt
}

// FuzzyQueryJobs 根据工作id，地区，页码获取信息
func (j *JobService) FuzzyQueryJobs(fuzzyName, queryType, host string, pageNum, ifAdmin int) []mysqlModel.UserComArticle {
	var userArt []mysqlModel.UserComArticle
	pageSize := pageSize
	if ifAdmin == 1 {
		pageSize = webPageSize
	}
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, pageSize)
	dao.DB.Table("articles").Select("tags,art_id,career_job_id,job_label,title,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,companies.com_level,articles.recommend,articles.art_type").
		Joins("INNER JOIN users on articles .boss_id = users.user_id "+
			"INNER JOIN companies on companies.com_id = articles.company_id ").
		Where("`status` = ?", 1).Where("`show` = ?", 1).
		Where("articles.art_type = ?", queryType).
		Where("LOCATE(?,articles.title) > 0", fuzzyName).
		Order("recommend desc,art_id desc").
		Limit(pageSize).Offset(sqlPage).Find(&userArt)
	if ifAdmin == 1 {
		for i, m := 0, len(userArt); i < m; i++ {
			userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
		}
		return userArt
	}
	for i, ul := 0, len(userArt); i < ul; i++ {
		if queryType == "request" {
			userArt[i].CompanyName = fmt.Sprintf("%v*****", string([]rune(userArt[i].CompanyName)[:1]))
		}
		userArt[i].TagsOut = sqlUtil.SqlStringToSli(userArt[i].Tags)
		userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
	}
	return userArt
}

// GetJobs 根据工作id，地区，页码获取信息
func (j *JobService) GetJobs(region, labelType, host string, jobId, pageNum, pageSize, ifAdmin int) []mysqlModel.UserComArticle {
	var userArt []mysqlModel.UserComArticle
	pageNumStr := strconv.Itoa((pageNum - 1) * pageSize)
	jobsQ := dao.DB.Table("articles").Select("tags,art_id,career_job_id,job_label,title,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,companies.com_level,articles.recommend,articles.art_type").
		Joins("INNER JOIN users on articles .boss_id = users.user_id "+
			"INNER JOIN companies on companies.com_id = articles.company_id ").
		Where("`status` = ?", 1).Where("`show` = ?", 1)
	if region != "" {
		jobsQ = jobsQ.Where("region = ?", region)
	}
	if jobId != 0 {
		jobsQ = jobsQ.Where("career_job_id = ?", jobId)
	} else {
		jobsQ = jobsQ.Where("articles.art_type = ?", labelType)
	}
	jobsQ.Order("recommend desc,art_id desc").
		Limit(pageSize).Offset(pageNumStr).Find(&userArt)
	if ifAdmin == 1 {
		for i, m := 0, len(userArt); i < m; i++ {
			userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
		}
		return userArt
	}
	for i, ul := 0, len(userArt); i < ul; i++ {
		if labelType == "request" {
			userArt[i].CompanyName = fmt.Sprintf("%v*****", string([]rune(userArt[i].CompanyName)[:1]))
		}
		userArt[i].TagsOut = sqlUtil.SqlStringToSli(userArt[i].Tags)
		userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
	}
	return userArt
}

// GetOneArtInfo 获取详细招聘信息
func (j *JobService) GetOneArtInfo(artId, host string) (artInfo mysqlModel.OneArticleOut, err error) {
	var oneArt mysqlModel.OneArticleOut
	err = dao.DB.Table("articles").Select("articles.art_type,companies.com_level,president,job_label,tags,career_job_id,art_id,title,com_id,pic_url,content,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,create_time").
		Joins("INNER JOIN users on articles.boss_id = users.user_id").
		Joins("INNER JOIN companies on companies.com_id = articles.company_id").
		Where("art_id = ?", artId).Scan(&oneArt).Error
	if err != nil {
		return oneArt, err
	}
	oneArt.TagsOut = sqlUtil.SqlStringToSli(oneArt.Tags)
	oneArt.CreateTimeOut = oneArt.CreateTime.Format("2006-01-02")
	oneArt.HeadPic = formatUtil.GetPicHeaderBody(host, oneArt.HeadPic)
	oneArt.PicUrl = formatUtil.GetPicHeaderBody(host, oneArt.PicUrl)
	return oneArt, nil
}

// GetJobs2 more regions
func (j *JobService) GetJobs2(jobId string, region []string, pageNum int) []mysqlModel.UserComArticle {
	var userArt []mysqlModel.UserComArticle
	pageNumStr := strconv.Itoa((pageNum - 1) * 10)

	regionSql := dealSliToSql(region)
	fmt.Println(regionSql)
	sql := "SELECT id,boss_id,content,create_time,`view`,username,head_pic FROM `articles` INNER JOIN users on articles.boss_id = users.user_id"
	careerAndRegion := " where career_job_id = ? and region in ? order by art_id desc limit ?,10"
	//高效率
	//pageSql := userid >= (SELECT userid FROM `times` LIMIT" + pageNumStr + ", 1) LIMIT 10;"
	finalSql := sql + careerAndRegion
	//fmt.Println(finalSql)
	dao.DB.Raw(finalSql, jobId, regionSql, pageNumStr).Scan(&userArt)
	return userArt
}

func (j *JobService) GetCollectArts(username, classification, host string, pageNum int) []mysqlModel.UserComArticle {
	var userArts []mysqlModel.UserComArticle
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	//mainSql := "SELECT art_id,companies.com_level,title,job_label,tags,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale FROM `articles` " +
	//	" INNER JOIN users on articles.boss_id = users.user_id  " +
	//	" INNER JOIN companies on companies.com_id = articles.company_id " +
	//	"where art_id in (select art_id from collections where username = ? and classification = ? )"
	//limitSql := " and articles.`show` = 1 order by art_id desc limit ?,10"
	//finalSql := mainSql + limitSql
	////fmt.Println(finalSql)
	//dao.DB.Debug().Raw(finalSql, username, classification, sqlPage).Find(&userArts)

	dao.DB.Table("collections").Select("collections.art_id,companies.com_level,title,job_label,tags,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale").
		Joins("INNER JOIN articles ON articles.art_id = collections.art_id INNER JOIN companies ON articles.company_id = companies.com_id INNER JOIN users ON users.user_id = articles.boss_id").
		Where("collections.username = ?", username).Where("classification = ?", classification).
		Where("articles.`show` = ?", 1).Where("articles.`status` = ?", 1).Order("collections.id desc").
		Offset(sqlPage).Limit(pageSize).Find(&userArts)

	for i, ul := 0, len(userArts); i < ul; i++ {
		if classification != "job" {
			userArts[i].CompanyName = fmt.Sprintf("%v*****", string([]rune(userArts[i].CompanyName)[:1]))
		}
		userArts[i].TagsOut = sqlUtil.SqlStringToSli(userArts[i].Tags)
		userArts[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArts[i].HeadPic)
	}
	return userArts
}

// GetJobInfoByBossId 根据bossid获取工作信息
func (j *JobService) GetJobInfoByBossId(bossId, queryType, host string, desStatus, desShow, pageNum int) []mysqlModel.UserComArtBoss {
	var userArt []mysqlModel.UserComArtBoss
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	// 链式查询，用于条件筛选
	sqlQuery := dao.DB.Table("articles").
		Select("`show`,status,com_level,art_id,job_label,tags,title,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,articles.art_type").
		Joins(" INNER JOIN users on articles.boss_id = users.user_id  " +
			"INNER JOIN companies on companies.com_id = articles.company_id ")
	if desStatus != -1 {
		sqlQuery = sqlQuery.Where("`status`=?", desStatus)
	}
	if desShow != -1 {
		sqlQuery = sqlQuery.Where("`show` = ?", desShow)
	}
	if queryType != "all" {
		sqlQuery = sqlQuery.Where("articles.art_type = ?", queryType)
	}
	sqlQuery = sqlQuery.Where("boss_id=?", bossId).
		Order("art_id desc").Limit(10).Offset(sqlPage)
	sqlQuery.Find(&userArt)
	for i, ul := 0, len(userArt); i < ul; i++ {
		if queryType != "job" {
			userArt[i].CompanyName = fmt.Sprintf("%v*****", string([]rune(userArt[i].CompanyName)[:1]))
		}
		userArt[i].TagsOut = sqlUtil.SqlStringToSli(userArt[i].Tags)
		userArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, userArt[i].HeadPic)
	}
	return userArt
}

func (j *JobService) QueryArtWeight(artId string) (aw mysqlModel.ArtWeightCount) {
	sql := "SELECT art_id,`view`,collect,articles.update_time,articles.delivery_num,person_scale FROM `articles` INNER JOIN companies ON companies.com_id = articles.company_id where art_id = ?"
	dao.DB.Raw(sql, artId).Scan(&aw)
	return
}

func (j *JobService) GetCompanyJobInfo(companyName, JobReqType, getType, host string, pageNumInt, status int) []mysqlModel.UserComArticle {
	var comArt []mysqlModel.UserComArticle
	sqlPage := sqlUtil.PageNumToSqlPage(pageNumInt, 10)
	sqlQuery := dao.DB.Table("articles").
		Select("art_id,title,com_level,`view`,boss_id,nickname,head_pic," +
			"company_name,salary_min,job_label,tags,salary_max,region,person_scale,articles.art_type").
		Joins(" INNER JOIN users on articles.boss_id = users.user_id  " +
			"INNER JOIN companies on companies.com_id = articles.company_id ")
	sqlQuery = sqlQuery.Where("company_name = ?", companyName).
		Where("`status`=?", status).
		Where("`show` = ?", 1).
		Order("art_id desc")
	if JobReqType != "" {
		sqlQuery = sqlQuery.Where("`job_label` = ?", JobReqType)
	} else {
		if getType != "" {
			//sqlQuery = sqlQuery.Joins("INNER JOIN labels on articles.career_job_id = labels.id").
			//	Where("labels.`type` = ?", getType)
			sqlQuery = sqlQuery.Where("articles.art_type = ?", getType)
		}
	}
	sqlQuery.Limit(10).Offset(sqlPage).Find(&comArt)
	for i, ul := 0, len(comArt); i < ul; i++ {
		if comArt[i].ComLevel == 2 || getType == "request" {
			comArt[i].CompanyName = fmt.Sprintf("%v*****", string([]rune(comArt[i].CompanyName)[:1]))
		}
		comArt[i].TagsOut = sqlUtil.SqlStringToSli(comArt[i].Tags)
		comArt[i].HeadPic = formatUtil.GetPicHeaderBody(host, comArt[i].HeadPic)
	}
	return comArt
}

func (j *JobService) GetAllCompanyJobInfoXls(comId int) []mysqlModel.UserComArticleXls {
	var comArt []mysqlModel.UserComArticleXls
	sqlQuery := dao.DB.Table("articles").
		Select("title,`view`,nickname,users.name,articles.`status`,company_name,salary_min,job_label,tags,salary_max,person_scale,articles.art_type,articles.content").
		Joins(" INNER JOIN users on articles.boss_id = users.user_id  " +
			"INNER JOIN companies on companies.com_id = articles.company_id ")
	sqlQuery = sqlQuery.Where("articles.company_id = ?", comId).
		Order("articles.art_id desc").Find(&comArt)
	return comArt
}

func (j *JobService) FuzzyQueryJobOrReq(fuzzyJobName, fuzzyType string,
) ([]mysqlModel.FuzzyArtInfo, error) {
	var fuzzyList []mysqlModel.FuzzyArtInfo

	err := dao.DB.Table("articles").Select("art_id,title,articles.art_type,`show`").
		Where("articles.art_type = ?", fuzzyType).Where("LOCATE(?,articles.title) > 0", fuzzyJobName).
		Where("`show` = ?", 1).Order("art_id desc").Limit(50).
		Find(&fuzzyList).Error
	//sql := "SELECT art_id,title,`type` FROM articles " +
	//	"INNER JOIN labels ON labels.id = articles.career_job_id " +
	//	"WHERE `type` = ? AND  LOCATE(?,articles.title) > 0 " +
	//	"order by art_id desc limit 50"
	//err := dao.DB.Raw(sql, fuzzyType, fuzzyJobName).Scan(&fuzzyList).Error
	return fuzzyList, err
}

// QueryWaitingApply 查询待审核的(web)
func (as *ArticleService) QueryWaitingApply(queryType, host string, pageNum, comId int) (waitInfo []mysqlModel.UserComArticle) {
	var waitingInfo []mysqlModel.UserComArticle
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	sqlQuery := dao.DB.Table("articles").
		Select("`show`,career_job_id,status,art_id,title,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,articles.art_type").
		Joins(" INNER JOIN users on articles.boss_id = users.user_id  "+
			"INNER JOIN companies on companies.com_id = articles.company_id ").
		Where("status = ?", 0)
	if queryType != "all" {
		sqlQuery = sqlQuery.Where("articles.job_label = ?", queryType)
	}
	if comId != 0 {
		sqlQuery = sqlQuery.Where("articles.company_id = ?", comId)
	}
	sqlQuery.Order("art_id desc").Limit(webPageSize).Offset(sqlPage).Find(&waitingInfo)
	for i, m := 0, len(waitingInfo); i < m; i++ {
		waitingInfo[i].HeadPic = formatUtil.GetPicHeaderBody(host, waitingInfo[i].HeadPic)
	}
	return waitingInfo
}

// QuerySonWaitingApply 查询父类对应子类的待审核信息
func (as *ArticleService) QuerySonWaitingApply(sonIdSli []int, pageNum, comId int, host string) (waitInfo []mysqlModel.UserComArticle) {
	var waitingInfo []mysqlModel.UserComArticle
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	sqlQuery := dao.DB.Table("articles").
		Select("`show`,career_job_id,status,art_id,title,`view`,boss_id,nickname,head_pic,company_name,salary_min,salary_max,region,person_scale,articles.art_type").
		Joins(" INNER JOIN users on articles.boss_id = users.user_id  "+
			"INNER JOIN companies on companies.com_id = articles.company_id ").
		Where("status = ?", 0)
	if comId != 0 {
		sqlQuery = sqlQuery.Where("articles.company_id = ?", comId)
	}
	sqlQuery.Where("career_job_id in (?)", sonIdSli).Order("art_id desc").Limit(webPageSize).Offset(sqlPage).Find(&waitingInfo)
	for i, m := 0, len(waitingInfo); i < m; i++ {
		waitingInfo[i].HeadPic = formatUtil.GetPicHeaderBody(host, waitingInfo[i].HeadPic)
	}
	return waitingInfo
}

// 将切片转换为sql ('工商类','管理类')
func dealSliToSql(sli []string) string {
	sql := "("
	for _, v := range sli {
		sql = sql + "'" + v + "',"
	}
	fmt.Println(len(sql))
	sql = sql[:len(sql)-1] + ")"
	return sql
}
