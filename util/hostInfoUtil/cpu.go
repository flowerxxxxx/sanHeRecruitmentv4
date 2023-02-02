package hostInfoUtil

import (
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return roundTo(percent[0], 1)
}
