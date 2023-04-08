package controllerModule

import (
	"github.com/pkg/errors"
	"sanHeRecruitment/models/moduleModel"
	"sanHeRecruitment/service/mysqlService"
	"sort"
)

type DataControlModule struct {
	*mysqlService.DailySaverService
}

// GetDailyHotLabel 获取单日标签热度排序
func (dcm DataControlModule) GetDailyHotLabel(queryDate, artType, queryType string) []moduleModel.KVPair {
	dailyLabel, _ := dcm.DailySaverService.QueryDailyLabel(queryDate, artType, queryType)
	labelMap := map[string]int{}
	for _, v := range dailyLabel {
		if _, ok := labelMap[v.Label]; ok {
			if queryType == "day_delivery" {
				labelMap[v.Label] += v.DayDelivery
			} else {
				labelMap[v.Label] += v.DayView
			}
		} else {
			if queryType == "day_delivery" {
				labelMap[v.Label] = v.DayDelivery
			} else {
				labelMap[v.Label] = v.DayView
			}
		}
	}
	tmpList := make([]moduleModel.KVPair, 0)
	for k, v := range labelMap {
		tmpList = append(tmpList, moduleModel.KVPair{Key: k, Value: v})
	}
	sort.Slice(tmpList, func(i, j int) bool {
		return tmpList[i].Value > tmpList[j].Value
	})
	return tmpList
}

func (dcm DataControlModule) HotLabelCutSliByPageNum(desSli []moduleModel.KVPair, pageNum, pageSize int) (data interface{}, err error) {
	desSliLen := len(desSli)
	if desSliLen == 0 || pageNum <= 0 {
		err = errors.New("no Data")
		return
	}
	if pageNum*pageSize >= desSliLen && pageSize*(pageNum-1) <= desSliLen {
		return desSli[(pageNum-1)*pageSize:], nil
	} else if pageNum*pageSize <= desSliLen {
		return desSli[(pageNum-1)*pageSize : pageNum*pageSize], nil
	} else {
		err = errors.New("pageNum err")
	}
	return
}
