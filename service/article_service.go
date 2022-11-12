package service

import (
	"github.com/jinzhu/gorm"
	"log"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"strconv"
	"time"
)

type ArticleService struct {
}

// AddArtView 增加文章阅读量
func (as *ArticleService) AddArtView(artId int) {
	//var ArtInfo mysqlModel.Article
	//dao.DB.Where("art_id=?", artId).Find(&ArtInfo)
	//ArtInfo.View += 1
	//dao.DB.Save(&ArtInfo)
	err := dao.DB.Table("articles").Where("art_id=?", artId).
		Update("view", gorm.Expr("`view`+ ?", 1)).Error
	if err != nil {
		log.Println("AddArtView failed,err:", err)
	}
}

// AddDeliveryNum 增加简历投递量
func (as *ArticleService) AddDeliveryNum(artId string) {
	var ArtInfo mysqlModel.Article
	artIdInt, _ := strconv.Atoi(artId)
	dao.DB.Where("art_id=?", artIdInt).Find(&ArtInfo)
	ArtInfo.DeliveryNum += 1
	dao.DB.Save(&ArtInfo)
}

// AddArtCollectNum 增加招聘收藏次数
func (as *ArticleService) AddArtCollectNum(artId string) {
	var ArtInfo mysqlModel.Article
	artIdInt, _ := strconv.Atoi(artId)
	dao.DB.Where("art_id=?", artIdInt).Find(&ArtInfo)
	ArtInfo.Collect += 1
	dao.DB.Save(&ArtInfo)
}

// BossChangeArtShowStatus boss修改文章展示状态
func (as *ArticleService) BossChangeArtShowStatus(bossId, artId int, action string) error {
	var ArtInfo mysqlModel.Article
	err := dao.DB.Where("art_id=?", artId).Where("boss_id=?", bossId).Find(&ArtInfo).Error
	if err != nil {
		return NoRecord
	}
	if action == "hide" {
		ArtInfo.Show = 0
		err = dao.DB.Save(&ArtInfo).Error
		if err != nil {
			return ServiceErr
		}
	} else {
		ArtInfo.Show = 1
		err = dao.DB.Save(&ArtInfo).Error
		if err != nil {
			return ServiceErr
		}
	}
	return nil
}

// BossDeletePubArt boss删除发布的需求
func (as *ArticleService) BossDeletePubArt(bossId, artId int) error {
	err := dao.DB.Where("art_id=?", artId).Where("boss_id=?", bossId).Delete(&mysqlModel.Article{}).Error
	if err != nil {
		return NoRecord
	}
	return err
}

// AdminDeletePubInfo admin删除发布的需求
func (as *ArticleService) AdminDeletePubInfo(artId int) error {
	err := dao.DB.Where("art_id=?", artId).Delete(&mysqlModel.Article{}).Error
	if err != nil {
		return NoRecord
	}
	return err
}

// UpdateArtWeight 更新招聘/需求权值
func (as *ArticleService) UpdateArtWeight(artId int, artWeight float64) {
	//var artInfo mysqlModel.Article
	//dao.DB.Where("art_id=?", artId).Find(&artInfo)
	//artInfo.Weight = artWeight
	//dao.DB.Save(&artInfo)
	err := dao.DB.Table("articles").Where("art_id=?", artId).Model(&mysqlModel.Article{}).
		UpdateColumn(map[string]interface{}{"weight": artWeight}).Error
	if err != nil {
		log.Println("UpdateArtWeight failed,err:", err)
	}
}

func (as *ArticleService) QueryArtByID(artID int) (artInfo mysqlModel.Article, err error) {
	err = dao.DB.Where("art_id=?", artID).Find(&artInfo).Error
	return
}

// AddNewEmployeeReq 创建新的招聘
func (as *ArticleService) AddNewEmployeeReq(
	careerJobIdInt, bossId, CompanyId, salaryMin, salaryMax, bossVip int,
	title, content, region, jobLabel, tagStr, artType string, time time.Time,
) (err error) {
	status := 0
	if bossVip == 1 {
		status = 1
	}
	var artInfo = mysqlModel.Article{
		CareerJobId: careerJobIdInt,
		BossId:      bossId,
		CompanyId:   CompanyId,
		JobLabel:    jobLabel,
		Title:       title,
		Content:     content,
		CreateTime:  time,
		UpdateTime:  time,
		Collect:     0,
		View:        0,
		DeliveryNum: 0,
		Weight:      0,
		Status:      status,
		Region:      region,
		SalaryMin:   salaryMin,
		SalaryMax:   salaryMax,
		Show:        1,
		Tags:        tagStr,
		ArtType:     artType,
	}
	err = dao.DB.Save(&artInfo).Error
	if err != nil {
		return
	}
	return
}

// ModifyArtStatus admin修改招聘/请求发布状态
func (as *ArticleService) ModifyArtStatus(artId, desStatusNum int) (err error) {
	var artInfo mysqlModel.Article
	err = dao.DB.Where("art_id=?", artId).Find(&artInfo).Error
	if err != nil {
		return
	}
	artInfo.Status = desStatusNum
	err = dao.DB.Save(&artInfo).Error
	if err != nil {
		return
	}
	return
}

func (as *ArticleService) BatchDeletePub(bossID int) (err error) {
	sqlStr := "DELETE FROM articles WHERE boss_id = ?"
	err = dao.DB.Exec(sqlStr, bossID).Error
	return
}

func (as *ArticleService) BatchDeleteSonLabelPub(careerJobId int) (err error) {
	sqlStr := "DELETE FROM articles WHERE career_job_id = ?"
	err = dao.DB.Exec(sqlStr, careerJobId).Error
	return
}

func (as *ArticleService) BatchDeleteFatherLabelPub(sonLabelSli []int) (err error) {
	sqlStr := "DELETE FROM articles WHERE career_job_id in (?)"
	err = dao.DB.Exec(sqlStr, sonLabelSli).Error
	return
}

func (as *ArticleService) BatchChangePubStatus(companyId, desStatus int) (err error) {
	batchSql := "UPDATE articles SET `status` = ? "
	if companyId != 0 {
		batchSql = "UPDATE articles SET `status` = ? WHERE company_id = ? "
		err = dao.DB.Exec(batchSql, desStatus, companyId).Error
		return
	}
	err = dao.DB.Exec(batchSql, desStatus).Error
	return
}

func (as *ArticleService) QueryComPubInfos(comId, status int) (comPubs []*mysqlModel.Article, err error) {
	if comId != 0 {
		err = dao.DB.Table("articles").
			Where("company_id = ?", comId).
			Where("`status` = ?", status).
			Find(&comPubs).Error
	}
	err = dao.DB.Table("articles").
		Where("`status` = ?", status).
		Find(&comPubs).Error
	return
}

func (as *ArticleService) ChangeTopPubStatus(id, desStatus int) (err error) {
	var ArtInfo mysqlModel.Article
	err = dao.DB.Table("articles").
		Where("art_id = ?", id).Find(&ArtInfo).Error
	if err != nil {
		return NoRecord
	}
	ArtInfo.Recommend = desStatus
	err = dao.DB.Table("articles").Save(&ArtInfo).Error
	return err
}

func (as *ArticleService) QueryMaxRecommendCount() (maxNum int, err error) {
	var maxRecCount mysqlModel.MaxRecommendCount
	err = dao.DB.Table("articles").
		Select("MAX(recommend) as max_num").Find(&maxRecCount).Error
	if err != nil {
		return -1, err
	}
	return maxRecCount.MaxNum, nil
}

// BatchKillUserAndPub 应用gorm事务，删除用户身份并且批量删除用户发布的文章
func (as *ArticleService) BatchKillUserAndPub(bossId int) (err error) {
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		//
		sqlStr := "DELETE FROM articles WHERE boss_id = ?"
		if err := tx.Exec(sqlStr, bossId).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		err := tx.Table("users").Where("user_id = ?", bossId).
			UpdateColumns(map[string]interface{}{
				"identy_pin": 0,
				"user_level": 0,
				"company_id": 0,
			}).Error
		if err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})
	return
}
