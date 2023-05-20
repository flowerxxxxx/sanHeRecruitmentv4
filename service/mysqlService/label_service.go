package mysqlService

import (
	"github.com/jinzhu/gorm"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
)

type LabelService struct {
}

type Label struct {
	mysqlModel.Label
	Value    int           `json:"value"`
	Children []interface{} `json:"children"`
}

func (ls *LabelService) QueryMaxLabelCount() (maxNum int, err error) {
	var maxRecCount mysqlModel.MaxLabelCount
	err = dao.DB.Table("labels").
		Select("MAX(sort_num) as max_num").Find(&maxRecCount).Error
	if err != nil {
		return -1, err
	}
	return maxRecCount.MaxNum, nil
}

func (ls *LabelService) ChangeTopPubStatus(id, desSort int) (err error) {
	var Labels mysqlModel.Label
	err = dao.DB.Table("labels").
		Where("id = ?", id).Find(&Labels).Error
	if err != nil {
		return NoRecord
	}
	Labels.SortNum = desSort
	err = dao.DB.Table("labels").Save(&Labels).Error
	return err
}

// GetLabelTree 获取标签的json树
func (ls *LabelService) GetLabelTree(labelType string, value int) []interface{} {
	var Labels []Label
	dao.DB.Where("type = ?", labelType).Order("sort_num desc").Find(&Labels)
	data := makeLabelTree(Labels, 0, 0, value)
	return data
}

// 将含等级和父子关系的数据递归成树
func makeLabelTree(labels []Label, level int, parentId int, value int) []interface{} {
	level++
	dataList := []interface{}{}
	dataLength := len(labels)
	for i := 0; i < dataLength; i++ {
		if labels[i].Level == level && labels[i].ParentId == parentId {
			labels[i].Children = makeLabelTree(labels, labels[i].Level, labels[i].ID, 0)
			labels[i].Value = value
			value++
			dataList = append(dataList, labels[i])
		}
	}
	return dataList
}

// QueryLabelInfo 查找标签信息
func (ls *LabelService) QueryLabelInfo(labelId string) (labelInfo mysqlModel.Label, err error) {
	sql := "select * from labels where id = ?"
	err = dao.DB.Raw(sql, labelId).Scan(&labelInfo).Error
	if err != nil {
		return labelInfo, err
	}
	return labelInfo, nil
}

// QueryLabelById 查找标签内容
func (ls *LabelService) QueryLabelById(labelId int) (label mysqlModel.Label, err error) {
	var labelInfo mysqlModel.Label
	sql := "select label,`type` from labels where id = ?"
	err = dao.DB.Raw(sql, labelId).Scan(&labelInfo).Error
	if err != nil {
		return labelInfo, err
	}
	return labelInfo, nil
}

func (ls *LabelService) AddLabel(labelType, label string, parentId, parentLevel int) error {
	var newLabel mysqlModel.Label
	if parentLevel == 0 {
		newLabel.Level = 1
		newLabel.ParentId = 0
	} else {
		newLabel.Level = 2
		newLabel.ParentId = parentId
	}
	newLabel.Label = label
	newLabel.Type = labelType
	newLabel.SortNum = 0
	newLabel.Recommend = 0
	err := dao.DB.Save(&newLabel).Error
	if err != nil {
		return err
	}
	return err
}

// CheckDuplicateLabel 标签查重
func (ls *LabelService) CheckDuplicateLabel(label, labelType string, parentIdInt int) bool {
	var RepeatLabel mysqlModel.Label
	err := dao.DB.Table("labels").Where("`type` = ?", labelType).Where("parent_id = ?", parentIdInt).Where("label = ?", label).Find(&RepeatLabel).Error
	if err != nil {
		return true
	}
	return false
}

func (ls *LabelService) DeleteLabel(labelId, labelLevel int) error {
	err := dao.DB.Table("labels").Where("id=?", labelId).Delete(&mysqlModel.Label{}).Error
	if err != nil {
		return err
	}
	if labelLevel == 1 {
		dao.DB.Table("labels").Where("parent_id=?", labelId).Delete(&mysqlModel.Label{})
	}
	return err
}

func (ls *LabelService) QueryLabelByContent(label string) (newLabel mysqlModel.Label) {
	dao.DB.Where("label=?", label).Find(&newLabel)
	return
}

func (ls *LabelService) QueryCompanyLabel(jobLabel string, companyId, status int) (label []mysqlModel.Label) {
	dao.DB.Table("articles").Select("DISTINCT labels.id,labels.label,labels.`level`,labels.parent_id,labels.type,"+
		"articles.career_job_id").
		Joins("INNER JOIN labels ON labels.id = articles.career_job_id ").
		Where("company_id=?", companyId).
		Where("articles.`status` = ?", status).
		Where("type = ?", jobLabel).Find(&label)
	return
}

func (ls *LabelService) QueryCompanyLabelAdmin(companyName, jobLabel string, status int) (label []mysqlModel.Label) {
	dao.DB.Table("articles").Select("DISTINCT labels.id,labels.label,labels.`level`,labels.parent_id,labels.type,"+
		"companies.com_id,companies.company_name,articles.career_job_id").
		Joins("INNER JOIN companies ON companies.com_id = articles.company_id").
		Joins("INNER JOIN labels ON labels.id = articles.career_job_id ").
		Where("company_name=?", companyName).
		Where("articles.`status` = ?", status).
		Where("type = ?", jobLabel).Find(&label)
	return
}

//func (ls *LabelService) EditLabel(labelId int, labelContent string) (err error) {
//	var labelInfo mysqlModel.Label
//	err = dao.DB.Table("labels").Where("id = ?", labelId).Find(&labelInfo).Error
//	if err != nil {
//		return NoRecord
//	}
//	labelInfo.Label = labelContent
//	err = dao.DB.Table("labels").Save(&labelInfo).Error
//	return
//}

func (ls *LabelService) EditLabel(labelId int, labelContent string) (err error) {
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作
		//
		var labelInfo mysqlModel.Label
		if err := tx.Table("labels").Where("id = ?", labelId).Find(&labelInfo).Error; err != nil {
			// 返回任何错误都会回滚事务
			return NoRecord
		}
		labelInfo.Label = labelContent
		if err := tx.Table("labels").Save(&labelInfo).Error; err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		err := tx.Table("articles").Where("career_job_id = ?", labelId).
			UpdateColumns(map[string]interface{}{
				"job_label": labelContent,
			}).Error
		if err != nil {
			return err
		}
		// 返回 nil 提交事务
		return nil
	})
	return
}

func (ls *LabelService) QuerySonLabelsById(fatherId int) (labelInfos []mysqlModel.Label, err error) {
	var labelInfo []mysqlModel.Label
	err = dao.DB.Table("labels").Where("parent_id = ?", fatherId).Find(&labelInfo).Error
	return labelInfo, err
}

func (ls *LabelService) QueryLabelInfoByLabel(fatherName string) (labelInfos mysqlModel.Label, err error) {
	var LabelInfo mysqlModel.Label
	err = dao.DB.Table("labels").Where("label = ?", fatherName).Find(&LabelInfo).Error
	if err != nil {
		return LabelInfo, NoRecord
	}
	return LabelInfo, err
}

func (ls *LabelService) QueryRecommendLabels() ([]mysqlModel.RecommendLabel, error) {
	var LabelInfos []mysqlModel.RecommendLabel
	err := dao.DB.Table("labels").Select("id,label,`type`,recommend").
		Where("recommend != ?", 0).Order("recommend desc").Find(&LabelInfos).Error
	return LabelInfos, err
}

func (ls *LabelService) ChangeLabelRecommend(id, desReco int) (err error) {
	var LabelInfo mysqlModel.Label
	err = dao.DB.Table("labels").
		Where("id = ?", id).Find(&LabelInfo).Error
	if err != nil {
		return NoRecord
	}
	LabelInfo.Recommend = desReco
	err = dao.DB.Table("labels").Save(&LabelInfo).Error
	return err
}

func (ls *LabelService) QueryMaxRecommendCount() (maxNum int, err error) {
	var maxRecCount mysqlModel.MaxRecommendCount
	err = dao.DB.Table("labels").
		Select("MAX(recommend) as max_num").Find(&maxRecCount).Error
	if err != nil {
		return -1, err
	}
	return maxRecCount.MaxNum, nil
}
