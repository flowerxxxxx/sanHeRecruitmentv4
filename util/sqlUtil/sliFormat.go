package sqlUtil

import (
	"fmt"
	"strings"
)

// JudgeFormatLegal 判断标签是否合法
func JudgeFormatLegal(upLoadSli []string) bool {
	for _, item := range upLoadSli {
		judgeFlag := strings.Index(item, ",")
		if judgeFlag != -1 {
			return false
		}
	}
	return true
}

// SliToSqlString make sli to string to save in mysql
func SliToSqlString(upLoadSli []string) string {
	return strings.Join(upLoadSli, ",")
}

// SqlStringToSli make sqlString to slice for return
func SqlStringToSli(LoadString string) []string {
	if LoadString == "" {
		return []string{}
	} else {
		return strings.Split(LoadString, ",")
	}
}

// MakeDisappearFormat make **format for Name
func MakeDisappearFormat(inStr string) string {
	return fmt.Sprintf("%v*****", string([]rune(inStr)))
}
