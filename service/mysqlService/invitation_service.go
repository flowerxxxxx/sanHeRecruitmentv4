package mysqlService

import (
	"fmt"
	"sanHeRecruitment/dao"
	"sanHeRecruitment/models/mysqlModel"
	"sanHeRecruitment/util/formatUtil"
	"sanHeRecruitment/util/sqlUtil"
	"sanHeRecruitment/util/timeUtil"
	"time"
)

type InvitationService struct {
}

func (is *InvitationService) AddInvitationInfo(artId int, bossUsername, desUsername string, inviteTime time.Time) (err error) {
	var invInfo = mysqlModel.Invitation{
		ArtId:        artId,
		BossUsername: bossUsername,
		DesUsername:  desUsername,
		InviteTime:   inviteTime,
	}
	err = dao.DB.Save(&invInfo).Error
	return
}

func (is *InvitationService) QueryInvitation(username string, pageNum int) (invInfo []mysqlModel.Invitation, err error) {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	fmt.Println(sqlPage)
	err = dao.DB.Where("des_username=?", username).Order("id desc").Limit("?,?").Find(&invInfo).Error
	return
}

func (is *InvitationService) QueryOneInvitation(artId int, desUsername string) (invInfo mysqlModel.Invitation, err error) {
	err = dao.DB.Where("art_id=?", artId).Where("des_username=?", desUsername).Find(&invInfo).Error
	return
}

func (is *InvitationService) UserDeleteInfo(username, id string) (err error) {
	err = dao.DB.Debug().Table("invitations").Where("id = ?", id).Where("des_username = ?", username).Delete(&mysqlModel.Invitation{}).Error
	//sqlStr := "DELETE FROM invitations WHERE id = ? AND des_username = ?"
	//err = dao.DB.Debug().Raw(sqlStr,id,username).Error
	return
}

// QueryInvitationInfos 邀请投递页面
func (is *InvitationService) QueryInvitationInfos(desUsername, host string, pageNum int) ([]mysqlModel.InvitationArt, error) {
	sqlPage := sqlUtil.PageNumToSqlPage(pageNum, 10)
	var invInfos []mysqlModel.InvitationArt
	err := dao.DB.Table("articles").
		Select("invitations.id,job_label,tags,articles.art_id,articles.title,articles.salary_min,"+
			"articles.salary_max,articles.region,companies.company_name,articles.view,"+
			"companies.person_scale,users.nickname,users.head_pic,invitations.invite_time ").
		Joins("INNER JOIN invitations on articles.art_id = invitations.art_id "+
			"INNER JOIN companies on articles.company_id = companies.com_id "+
			"INNER JOIN users on users.user_id = articles.boss_id").
		Where("invitations.des_username = ?", desUsername).
		Order("invitations.id DESC").Limit(10).Offset(sqlPage).Find(&invInfos).Error
	for i, n := 0, len(invInfos); i < n; i++ {
		invInfos[i].InviteTimeOut = timeUtil.TimeFormatToStr(invInfos[i].InviteTime)
		invInfos[i].TagsOut = sqlUtil.SqlStringToSli(invInfos[i].Tags)
		invInfos[i].HeadPic = formatUtil.GetPicHeaderBody(host, invInfos[i].HeadPic)
	}
	return invInfos, err
}
