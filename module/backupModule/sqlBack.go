package backupModule

import (
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

/**
 *
 * 备份MySql数据库
 * @param 	host: 			数据库地址: localhost
 * @param 	port:			端口: 3306
 * @param 	user:			用户名: root
 * @param 	password:		密码: root
 * @param 	databaseName:	需要被分的数据库名: test
 * @param 	tableName:		需要备份的表名: user
 * @param 	sqlPath:		备份SQL存储路径: D:/backup/test/
 * @return 	backupPath
 *
 */

func BackupMySqlDb(host, port, user, password, databaseName, tableName, sqlPath string) (error, string) {
	//定义Cmd结构体对象指针
	var cmd *exec.Cmd

	//在这里如果没有传输表名，那么将会备份整个数据库，否则将只备份自己传入的表
	if tableName == "" {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName)
	} else {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName, tableName)
	}

	//StdinPipe方法返回一个在命令Start后与命令标准输入关联的管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return err, ""
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err, ""
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return err, ""
	}
	//获得一个当前的时间戳
	now := time.Now().Format("20060102150405")
	var backupPath string

	//设置我们备份文件的名字
	if tableName == "" {
		backupPath = sqlPath + databaseName + "_" + now + ".sql"
	} else {
		backupPath = sqlPath + databaseName + "_" + tableName + "_" + now + ".sql"
	}
	//写入文件并设置文件权限
	err = ioutil.WriteFile(backupPath, bytes, 0644)

	if err != nil {
		log.Println(err)
		return err, ""
	}
	return nil, backupPath
}
