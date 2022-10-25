package formatUtil

import "strings"

func GetPicHeaderBody(host, storeUrl string) string {
	if strings.Index(storeUrl, "uploadPic/") != -1 {
		return "https://" + host + "/" + storeUrl
	} else {
		return storeUrl
	}
}

func SavePicHeaderCutter(sourceUrl string) string {
	saveFlag := strings.Index(sourceUrl, "uploadPic/")
	return sourceUrl[saveFlag:]
}
