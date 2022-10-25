package uploadUtil

import (
	"sanHeRecruitment/config"
	"sanHeRecruitment/util"
	"strconv"
	"time"
)

// SaveFormat 获取存储文件url和存储地址
func SaveFormat(fileFormat, host string) (fileUrl, fileAddr string) {
	uuid := util.GetUUID() + strconv.Itoa(int(time.Now().Unix()))
	nowTimeUnix := strconv.Itoa(int(time.Now().Unix()))
	fileName := uuid + "-" + nowTimeUnix + fileFormat
	fileUrl = "uploadPic/" + fileName
	fileAddr = config.PicSaverPath + "/" + fileName
	return fileUrl, fileAddr
}

func FormatJudge(fileFormat string, FormatDepends ...string) bool {
	for _, item := range FormatDepends {
		if fileFormat == item {
			return true
		}
	}
	return false
}
