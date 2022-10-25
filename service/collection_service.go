package service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
)

type CollectionService struct {
}

// CollectArticle 添加收藏记录
func (c *CollectionService) CollectArticle(username string, art_id int, labelType string) (err error) {
	var col mysqlModel.Collection
	err = dao.DB.Where("username=?", username).Where("art_id=?", art_id).Find(&col).Error
	if err == nil {
		return Existed
	}
	col.Username = username
	col.Art_id = art_id
	col.Classification = labelType
	err = dao.DB.Save(&col).Error
	if err != nil {
		return MysqlErr
	}
	return
}

// QueryColRec 查询收藏记录
func (c *CollectionService) QueryColRec(artId, username string) (colInfo mysqlModel.Collection, err error) {
	var col mysqlModel.Collection
	err = dao.DB.Where("username=?", username).Where("art_id=?", artId).Find(&col).Error
	if err != nil {
		return col, NoRecord
	}
	return col, nil
}

// DeleteRecord 取消收藏记录
func (c *CollectionService) DeleteRecord(username string, art_id int) {
	dao.DB.Where("username=?", username).Where("art_id=?", art_id).Delete(&mysqlModel.Collection{})
}
