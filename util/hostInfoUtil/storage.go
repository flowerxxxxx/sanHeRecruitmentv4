package hostInfoUtil

import (
	"github.com/shirou/gopsutil/v3/disk"
)

type storageInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// GetStorageInfo 获取磁盘存储信息2
func GetStorageInfo() *storageInfo {
	sInfo, _ := disk.Usage(".")
	return &storageInfo{
		Total:       sInfo.Total / 1024 / 1024 / 1024,
		Free:        sInfo.Free / 1024 / 1024 / 1024,
		Used:        sInfo.Total/1024/1024/1024 - sInfo.Free/1024/1024/1024,
		UsedPercent: roundTo(sInfo.UsedPercent, 1),
	}
}
