package mysqlService

import "errors"

var (
	MysqlErr   = errors.New("数据库错误")
	NoRecord   = errors.New("无结果")
	Existed    = errors.New("已存在")
	HasFound   = errors.New("该数据已存在")
	ServiceErr = errors.New("服务器错误")
)
