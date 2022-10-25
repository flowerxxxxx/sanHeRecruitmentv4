package sqlUtil

func PageNumToSqlPage(pageNum, pageSize int) (sqlLimitPage int) {
	sqlLimitPage = (pageNum - 1) * pageSize
	return
}
