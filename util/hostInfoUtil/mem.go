package hostInfoUtil

import (
	"github.com/shirou/gopsutil/mem"
)

type memInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

func GetMemInfo() *memInfo {
	sInfo, _ := mem.VirtualMemory()
	return &memInfo{
		Total:       sInfo.Total / 1024 / 1024 / 1024,
		Available:   sInfo.Available / 1024 / 1024 / 1024,
		Used:        sInfo.Total/1024/1024/1024 - sInfo.Available/1024/1024/1024,
		UsedPercent: roundTo(sInfo.UsedPercent, 1),
	}
}
