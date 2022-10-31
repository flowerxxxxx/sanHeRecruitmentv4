package formatUtil

import (
	"errors"
	"strings"
)

func GetPicHeaderBody(host, storeUrl string) string {
	if strings.Index(storeUrl, "uploadPic/") != -1 {
		return "https://" + host + "/" + storeUrl
	} else {
		return storeUrl
	}
}

func SavePicHeaderCutter(sourceUrl string) (string, error) {
	saveFlag := strings.Index(sourceUrl, "uploadPic/")
	if saveFlag == -1 {
		return "", errors.New("url err,cut failed")
	}
	return sourceUrl[saveFlag:], nil
}
