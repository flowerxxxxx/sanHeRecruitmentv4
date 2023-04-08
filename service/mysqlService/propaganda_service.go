package mysqlService

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/timeUtil"
)

type PropagandaService struct {
}

func (ps *PropagandaService) QueryMaxRecommendCount() (maxNum int, err error) {
	var maxRecCount mysqlModel.MaxRecommendCount
	err = dao.DB.Table("propagandas").
		Select("MAX(recommend) as max_num").Find(&maxRecCount).Error
	if err != nil {
		return -1, err
	}
	return maxRecCount.MaxNum, nil
}

func (ps *PropagandaService) ChangeTopProStatus(id, desStatus int) (err error) {
	var ProInfo mysqlModel.Propaganda
	err = dao.DB.Table("propagandas").
		Where("id = ?", id).Find(&ProInfo).Error
	if err != nil {
		return NoRecord
	}
	ProInfo.Recommend = desStatus
	err = dao.DB.Table("propagandas").Save(&ProInfo).Error
	return err
}

// AddProInfo 添加宣传的内容
func (ps *PropagandaService) AddProInfo(
	uploadTime *timeUtil.MyTime, url, uploader, content, title string, uploadType int) (err error) {
	var ProInfo mysqlModel.Propaganda
	ProInfo.UploadTime = uploadTime
	ProInfo.UpdateTime = uploadTime
	ProInfo.Content = content
	ProInfo.Type = uploadType
	ProInfo.Url = url
	ProInfo.Uploader = uploader
	ProInfo.Title = title
	err = dao.DB.Save(&ProInfo).Error
	return err
}

func (ps *PropagandaService) EditProInfo(id int,
	upDateTime *timeUtil.MyTime, url, uploader, content, title string, uploadType int) (err error) {
	var ProInfo mysqlModel.Propaganda
	err = dao.DB.Where("id = ?", id).Find(&ProInfo).Error
	if err != nil {
		return NoRecord
	}
	ProInfo.UpdateTime = upDateTime
	ProInfo.Content = content
	ProInfo.Type = uploadType
	ProInfo.Title = title
	ProInfo.Url = url
	ProInfo.Uploader = uploader
	err = dao.DB.Save(&ProInfo).Error
	return err
}

// QueryProInfos 查询宣传栏的信息
func (ps *PropagandaService) QueryProInfos(host string) ([]mysqlModel.PropagandaOutHead, error) {
	var proInfos []mysqlModel.PropagandaOutHead
	err := dao.DB.Table("propagandas").
		Order("recommend desc,id desc").Find(&proInfos).Error
	for i, m := 0, len(proInfos); i < m; i++ {
		proInfos[i].Url = formatUtil.GetPicHeaderBody(host, proInfos[i].Url)
	}
	return proInfos, err
}

// QueryOneProInfo  查询一条宣传栏的信息
func (ps *PropagandaService) QueryOneProInfo(ProId int, host string) (mysqlModel.PropagandaOutContent, error) {
	var proInfo mysqlModel.PropagandaOutContent
	err := dao.DB.Table("propagandas").
		Where("id = ?", ProId).Find(&proInfo).Error
	proInfo.Url = formatUtil.GetPicHeaderBody(host, proInfo.Url)
	return proInfo, err
}

// DeleteProInfo 删除宣传栏内容
func (ps *PropagandaService) DeleteProInfo(ProId int) (err error) {
	err = dao.DB.Table("propagandas").
		Where("id = ?", ProId).Delete(&mysqlModel.Propaganda{}).Error
	return
}
