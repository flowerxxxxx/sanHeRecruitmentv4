package service

import (
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
)

type NoticeService struct {
}

func (ns *NoticeService) ChangeTopNoticeStatus(id, desStatus int) (err error) {
	var ProInfo mysqlModel.Propaganda
	err = dao.DB.Table("notices").
		Where("id = ?", id).Find(&ProInfo).Error
	if err != nil {
		return NoRecord
	}
	ProInfo.Recommend = desStatus
	err = dao.DB.Table("notices").Save(&ProInfo).Error
	return err
}

func (ns *NoticeService) QueryMaxRecommendCount() (maxNum int, err error) {
	var maxRecCount mysqlModel.MaxRecommendCount
	err = dao.DB.Table("notices").
		Select("MAX(recommend) as max_num").Find(&maxRecCount).Error
	if err != nil {
		return -1, err
	}
	return maxRecCount.MaxNum, nil
}

// SaveNotice 保存公告信息
func (ns *NoticeService) SaveNotice(content, uploader, title string, uploadTime *timeUtil.MyTime) (err error) {
	var noticeSaver mysqlModel.Notice
	noticeSaver.Content = content
	noticeSaver.UpdateTime = uploadTime
	noticeSaver.UploadTime = uploadTime
	noticeSaver.Uploader = uploader
	noticeSaver.Title = title
	err = dao.DB.Table("notices").Save(&noticeSaver).Error
	return
}

// EditNotice 编辑公告栏内容
func (ns *NoticeService) EditNotice(id int, content, title string, updateTime *timeUtil.MyTime) (err error) {
	var noticeSaver mysqlModel.Notice
	err = dao.DB.Table("notices").Where("id = ?", id).Find(&noticeSaver).Error
	if err != nil {
		return
	}
	noticeSaver.Content = content
	noticeSaver.UpdateTime = updateTime
	noticeSaver.Title = title
	err = dao.DB.Table("notices").Save(&noticeSaver).Error
	return
}

// QueryNoticesInfos 获取公告栏信息
func (ns *NoticeService) QueryNoticesInfos(pageNum int) []mysqlModel.NoticeOutHead {
	var noticeInfos []mysqlModel.NoticeOutHead
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, pageSize)
	dao.DB.Table("notices").Select("id,title,upload_time,update_time,recommend").
		Order("recommend desc,id desc").Offset(sqlPage).Limit(pageSize).Find(&noticeInfos)
	return noticeInfos
}

// DeleteNotice 删除宣传栏内容
func (ns *NoticeService) DeleteNotice(NoticeId int) (err error) {
	err = dao.DB.Table("notices").
		Where("id = ?", NoticeId).Delete(&mysqlModel.Notice{}).Error
	return
}

// QueryOneNoticeInfo  查询一条公告栏信息
func (ns *NoticeService) QueryOneNoticeInfo(NoticeId int) (mysqlModel.NoticeOutContent, error) {
	var NoticeInfo mysqlModel.NoticeOutContent
	err := dao.DB.Table("notices").
		Where("id = ?", NoticeId).Find(&NoticeInfo).Error
	return NoticeInfo, err
}
