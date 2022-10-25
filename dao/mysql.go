package dao

import (
	"github.com/jinzhu/gorm"
	//"gorm.io/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sanHeRecruitment/config"
)

var (
	DB *gorm.DB
)

func InitMySQL() (err error) {
	dsn := config.MysqlConfig.Dsn
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	//设置连接最大存活时间连接最大存活时间
	DB.DB().SetConnMaxLifetime(config.MysqlConnMaxLivingTime)
	//DB.DB().SetMaxIdleConns(0)
	//DB.DB().SetMaxOpenConns(0)
	return DB.DB().Ping()
}

func Close() {
	DB.Close()
}
