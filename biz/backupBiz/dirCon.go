package backupBiz

import (
	"log"
	"os"
	"os/exec"
)

// 判断文件夹是否存在
func hasDir(path string) (bool, error) {
	_, _err := os.Stat(path)
	if _err == nil {
		return true, nil
	}
	if os.IsNotExist(_err) {
		return false, nil
	}
	return false, _err
}

// CreateDir 创建文件夹
func createDir(path string) (err error) {
	exist, err := hasDir(path)
	if err != nil {
		log.Printf("createDir %v 获取文件夹异常 -> %v\n", path, err)
		return
	}
	if exist {
		log.Println("createDir failed,existed")
	} else {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Printf("%v 创建目录异常 -> %v\n", path, err)
			return err
		} else {
			return nil
		}
	}
	return
}

// 删除文件夹
func removeDir(path string) error {
	_err := os.RemoveAll(path)
	return _err
}

func linuxRemoveDir(deleteFolder string) (err error) {
	cpCmd := exec.Command("rm", "-rf", deleteFolder)
	err = cpCmd.Run()
	return
}
