package controllerBiz

import (
	"sanHeRecruitment/models/mysqlModel"
	"sort"
)

type UserControllerModule struct {
}

func (u *UserControllerModule) MsgListSortByStartTime(qelm []mysqlModel.MsgObjUserOut) []mysqlModel.MsgObjUserOut {
	sort.Slice(qelm, func(i, j int) bool { // desc
		return qelm[i].MsgTime > qelm[j].MsgTime
	})
	return qelm
}
