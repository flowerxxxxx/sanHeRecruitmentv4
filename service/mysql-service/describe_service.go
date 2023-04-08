package mysql_service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/timeUtil"
)

type DescribeService struct {
}

// SaveNewDescription 添加新的平台简介
func (ds *DescribeService) SaveNewDescription(
	content, module, uploader string, uploadTime *timeUtil.MyTime) (err error) {
	var newDescribeSaver mysqlModel.Description
	newDescribeSaver.Content = content
	newDescribeSaver.Module = module
	newDescribeSaver.Uploader = uploader
	newDescribeSaver.UploadTime = uploadTime
	newDescribeSaver.UpdateTime = uploadTime
	err = dao.DB.Table("descriptions").Save(&newDescribeSaver).Error
	return
}

// EditDescription 编辑平台简介
func (ds *DescribeService) EditDescription(
	id int, content, module string, updateTime *timeUtil.MyTime) (err error) {
	var desSaver mysqlModel.Description
	err = dao.DB.Table("descriptions").Where("id = ?", id).Find(&desSaver).Error
	if err != nil {
		return NoRecord
	}
	desSaver.Content = content
	desSaver.Module = module
	desSaver.UpdateTime = updateTime
	err = dao.DB.Table("descriptions").Save(&desSaver).Error
	return
}

// QueryDescriptionInfos 获取平台简介
func (ds *DescribeService) QueryDescriptionInfos() []mysqlModel.DescriptionOut {
	var desInfos []mysqlModel.DescriptionOut
	dao.DB.Table("descriptions").Select("id,content,module,update_time").
		Find(&desInfos)
	return desInfos
}

// QueryModuleDescriptionInfo 分模块获取平台简介
func (ds *DescribeService) QueryModuleDescriptionInfo(module string) (mysqlModel.DescriptionOut, error) {
	var desInfos mysqlModel.DescriptionOut
	err := dao.DB.Table("descriptions").Select("id,content,module,update_time").
		Where("module = ?", module).Find(&desInfos).Error
	return desInfos, err
}

// DeleteDescription 删除平台简介
func (ds *DescribeService) DeleteDescription(DesId int) (err error) {
	err = dao.DB.Table("descriptions").
		Where("id = ?", DesId).Delete(&mysqlModel.Description{}).Error
	return
}
