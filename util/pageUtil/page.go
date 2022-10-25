package pageUtil

func TotalPageByTotalNum(TotalNum, pageSize int) (totalPage int) {
	if TotalNum%pageSize != 0 {
		totalPage = TotalNum/pageSize + 1
	} else {
		totalPage = TotalNum / pageSize
	}
	if totalPage == 0 {
		totalPage = 1
	}
	return
}
