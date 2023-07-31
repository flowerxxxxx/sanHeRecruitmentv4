package mysqlService

import (
	"github.com/jinzhu/gorm"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"time"
)

type UpgradeService struct {
}

// AddUpgradeInfo 增加升级身份的记录
func (uc *UpgradeService) AddUpgradeInfo(username string, targetLevel, companyId, CompanyExist int, applyTime time.Time, timeId int64) error {
	var UpgradeInfo = mysqlModel.Upgrade{
		Qualification: 0,
		TargetLevel:   targetLevel,
		FromUsername:  username,
		CompanyId:     companyId,
		ApplyTime:     applyTime,
		CompanyExist:  CompanyExist,
		Show:          0,
		TimeId:        timeId,
	}
	err := dao.DB.Save(&UpgradeInfo).Error
	if err != nil {
		return err
	}
	return err
}

// UpgradeInfoChangerUser user 应用gorm事务，自行处理身份升级需要修改的数据
func (uc *UpgradeService) UpgradeInfoChangerUser(username, president string, targetLevel, companyId, CompanyExist int, applyTime time.Time, timeId int64) (err error) {
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		//
		if err := tx.Table("users").Where("username = ?", username).
			UpdateColumns(map[string]interface{}{
				"company_id": companyId,
				"user_level": targetLevel,
				"president":  president,
			}).
			Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}

		var UpgradeInfo = mysqlModel.Upgrade{
			Qualification: 1,
			TargetLevel:   targetLevel,
			FromUsername:  username,
			CompanyId:     companyId,
			ApplyTime:     applyTime,
			CompanyExist:  CompanyExist,
			Show:          1,
			TimeId:        timeId,
		}
		if err := tx.Table("upgrades").Save(&UpgradeInfo).Error; err != nil {
			return err
		}
		//if err := tx.Table("upgrades").Where("id = ?", upgradeIdInt).
		//	UpdateColumns(map[string]interface{}{
		//		"qualification": 1,
		//	}).Error; err != nil {
		//	return err
		//}

		if err := tx.Table("companies").Where("com_id = ?", companyId).
			UpdateColumns(map[string]interface{}{
				"com_status": 1,
			}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return
}

func (uc *UpgradeService) QueryUpgradeInfoByUsername(username string) (upgrade mysqlModel.Upgrade, err error) {
	var upgradeInfo mysqlModel.Upgrade
	err = dao.DB.Where("from_username=?", username).Where("qualification=?", 0).Last(&upgradeInfo).Error
	if err != nil {
		return upgradeInfo, NoRecord
	}
	return upgradeInfo, nil
}

func (uc *UpgradeService) QueryUpgradeInfoById(upId int) (upgrade mysqlModel.Upgrade, err error) {
	var upgradeInfo mysqlModel.Upgrade
	err = dao.DB.Where("id=?", upId).Last(&upgradeInfo).Error
	if err != nil {
		return upgradeInfo, NoRecord
	}
	return upgradeInfo, nil
}

func (uc *UpgradeService) QueryUpgradeInfoByTimeId(TimeId64 int64, companyIdInt int) (upgrade mysqlModel.Upgrade, err error) {
	var upgradeInfo mysqlModel.Upgrade
	err = dao.DB.Where("time_id = ?", TimeId64).Where("company_id=?", companyIdInt).Last(&upgradeInfo).Error
	if err != nil {
		return upgradeInfo, NoRecord
	}
	return upgradeInfo, nil
}

func (uc *UpgradeService) ModifyUpgradeShow(username string, companyId int) error {
	var upgradeInfo mysqlModel.Upgrade
	err := dao.DB.Where("from_username=?", username).Where("company_id=?", companyId).Last(&upgradeInfo).Error
	if err != nil {
		return NoRecord
	}
	upgradeInfo.Show = 1
	dao.DB.Save(&upgradeInfo)
	return err
}

// ModifyUpgradeQualification 修改企业升级审核凭证 1通过 2未通过
func (uc *UpgradeService) ModifyUpgradeQualification(upId, qualification int) (err error) {
	var upgradeInfo mysqlModel.Upgrade
	dao.DB.Where("id=?", upId).Last(&upgradeInfo)
	upgradeInfo.Qualification = qualification
	err = dao.DB.Save(&upgradeInfo).Error
	return err
}

// QueryWaitingUpgrade 查询待升级的
func (uc *UpgradeService) QueryWaitingUpgrade(targetLevel, qualification, pageNum int) []mysqlModel.WaitingUpgradeOut {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, webPageSize)
	var waitingUpgradeInfo []mysqlModel.WaitingUpgradeOut
	qwuSql := dao.DB.Table("upgrades").
		Select("id,upgrades.company_id,target_level,upgrades.time_id,upgrades.qualification,upgrades.from_username,"+
			"apply_time,company_exist,company_name,users.name").
		Joins("INNER JOIN users on users.username = upgrades.from_username "+
			" INNER JOIN companies on companies.com_id = upgrades.company_id ").
		Where("upgrades.show = ?", 1).
		Where("target_level = ?", targetLevel)
	if qualification != -1 {
		qwuSql = qwuSql.Where("qualification = ?", qualification)
	}
	qwuSql.Order("upgrades.id desc").Limit(webPageSize).Offset(sqlPage).Find(&waitingUpgradeInfo)
	for i, n := 0, len(waitingUpgradeInfo); i < n; i++ {
		waitingUpgradeInfo[i].ApplyTimeOut = timeUtil.TimeFormatToStr(waitingUpgradeInfo[i].ApplyTime)
	}
	return waitingUpgradeInfo
}

// QueryAllWaitingUpgradeXls 查询全部的，xls导出者
func (uc *UpgradeService) QueryAllWaitingUpgradeXls(targetLevel, qualification int) ([]mysqlModel.WaitingUpgradeXls, error) {
	var waitingUpgradeInfo []mysqlModel.WaitingUpgradeXls
	qwuSql := dao.DB.Table("upgrades").
		Select("id,upgrades.company_id,target_level,upgrades.time_id,upgrades.qualification,upgrades.from_username,"+
			"apply_time,company_exist,company_name,users.name,companies.address,companies.phone,companies.description").
		Joins("INNER JOIN users on users.username = upgrades.from_username "+
			" INNER JOIN companies on companies.com_id = upgrades.company_id ").
		Where("upgrades.show = ?", 1).
		Where("target_level = ?", targetLevel)
	if qualification != -1 {
		qwuSql = qwuSql.Where("qualification = ?", qualification)
	}
	err := qwuSql.Order("upgrades.id desc").Find(&waitingUpgradeInfo).Error
	if err != nil {
		return waitingUpgradeInfo, NoRecord
	}
	for i, n := 0, len(waitingUpgradeInfo); i < n; i++ {
		waitingUpgradeInfo[i].ApplyTimeOut = timeUtil.TimeFormatToStr(waitingUpgradeInfo[i].ApplyTime)
	}
	return waitingUpgradeInfo, nil
}

// UpgradeInfoChanger 应用gorm事务，处理身份升级需要修改的数据
func (uc *UpgradeService) UpgradeInfoChanger(username string, companyId, targetLevel, upgradeIdInt, exist int) (err error) {
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		//
		if err := tx.Table("users").Where("username = ?", username).
			UpdateColumns(map[string]interface{}{
				"company_id": companyId,
				"user_level": targetLevel,
			}).
			Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err := tx.Table("upgrades").Where("id = ?", upgradeIdInt).
			UpdateColumns(map[string]interface{}{
				"qualification": 1,
			}).Error; err != nil {
			return err
		}
		if err := tx.Table("companies").Where("com_id = ?", companyId).
			UpdateColumns(map[string]interface{}{
				"com_status": 1,
			}).Error; err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return
}
