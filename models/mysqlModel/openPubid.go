package mysqlModel

import "sanHeRecruitment/dao"

type OpenPubid struct {
	ID     int    `json:"id" gorm:"primary_key"`
	Openid string `json:"openid"`
	Pubid  string `json:"pubid"`
}

func OpenToPub(openid string) string {
	var openPub OpenPubid
	sql := "select pubid from open_pubids where openid = '" + openid + "'"
	dao.DB.Raw(sql).Scan(&openPub)
	return openPub.Pubid
}
