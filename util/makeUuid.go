package util

import (
	"github.com/satori/go.uuid"
)

// 生成uuid
func GetUUID() string {
	u2 := uuid.NewV4()
	return u2.String()
}
