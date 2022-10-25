package mysqlModel

import (
	"sanHeRecruitment/util/timeUtil"
)

type MassSend struct {
	ID              int              `json:"id" gorm:"primary_key"`
	PublishUsername string           `json:"publish_username"`
	PublishTime     *timeUtil.MyTime `json:"publish_time" time_format:"2006-01-02 15:04:05"`
	Content         string           `json:"content"`
	DesRole         int              `json:"des_role"`
}
