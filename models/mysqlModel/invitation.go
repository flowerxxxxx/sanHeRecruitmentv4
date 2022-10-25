package mysqlModel

import "time"

type Invitation struct {
	ID           int       `json:"id" gorm:"primary_key"`
	ArtId        int       `json:"art_id"`
	BossUsername string    `json:"boss_username"`
	DesUsername  string    `json:"des_username"` //目标邀请用户用户名
	InviteTime   time.Time `json:"invite_time"`
}

// InvitationArt 邀请投递信息展示结构体
type InvitationArt struct {
	UserComArticle
	InviteTime    time.Time `json:"-"`
	InviteTimeOut string    `json:"invite_time"`
	ID            int       `json:"invite_id"` //邀请记录id
}
