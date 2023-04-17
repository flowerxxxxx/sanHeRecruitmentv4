package backupBiz

import (
	"fmt"
	"log"
	"runtime"
	"sanHeRecruitment/config"
	"time"
)

func Backer() (zipName string, errAny error) {
	sysType := runtime.GOOS
	timeId := time.Now().Format("20060102150405")
	DirName := fmt.Sprintf("backer_%s", timeId)
	backUpDirPath := config.BackUpConfig.SavePath + "/" + DirName
	picBackerDirPath := backUpDirPath + "/" + "picBackup"
	errPar := createDir(backUpDirPath)
	if errPar != nil {
		return "", errPar
	}
	errSon := createDir(picBackerDirPath)
	if errSon != nil {
		return "", errSon
	}
	if sysType == "windows" {
		errCo := fileCopy(config.PicSaverPath, picBackerDirPath)
		if errCo != nil {
			log.Println("windows backup fileCopy failed err,", errCo)
			return "", errCo
		}
	} else {
		errCo := linuxFileCp(config.PicSaverPath, picBackerDirPath)
		if errCo != nil {
			log.Println("linux backup fileCopy failed err,", errCo)
			return "", errCo
		}
	}
	errDump, _ := BackupMySqlDb(
		config.MysqlConfig.Host,
		config.MysqlConfig.Port,
		config.MysqlConfig.User,
		config.MysqlConfig.Password,
		config.MysqlConfig.DataBaseName,
		"",
		backUpDirPath+"/",
	)
	if errDump != nil {
		log.Println("backup BackupMySqlDb failed err,", errDump)
		return "", errDump
	}
	ziperName := DirName + ".zip"
	//dirZip(backUpDirPath, config.BackUpConfig.SavePath+"/"+ziperName)
	if sysType == "windows" {
		dirZip(backUpDirPath, config.BackUpConfig.SavePath+"/"+ziperName)
	} else {
		errZip := linuxZiper(config.BackUpConfig.SavePath+"/"+ziperName, backUpDirPath)
		if errZip != nil {
			log.Println("backup linuxZiper failed err,", errZip)
			return "", errZip
		}
	}

	go func(sysT string) {
		if sysT == "windows" {
			errDelete := removeDir(backUpDirPath)
			if errDelete != nil {
				log.Println("windows backup removeDir failed err,", errDelete)
			}
		} else {
			errDelete := linuxRemoveDir(backUpDirPath)
			if errDelete != nil {
				log.Println("linux backup removeDir failed err,", errDelete)
			}
		}
	}(sysType)
	return ziperName, nil
}
