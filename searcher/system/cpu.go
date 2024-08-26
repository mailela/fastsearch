package system

import (
	"math/rand"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUStatus struct {
	Cores       int     `json:"cores"`
	UsedPercent float64 `json:"usedPercent"`
	ModelName   string  `json:"modelName"`
}

func GetCPUStatus() CPUStatus {
	// percent, _ := cpu.Percent(time.Second, false)
	info, _ := cpu.Info()
	c := CPUStatus{
		// UsedPercent: GetPercent(percent[0]),
		UsedPercent: GetPercent(rand.Float64() * 10),
		Cores:       runtime.NumCPU(),
		ModelName:   info[0].ModelName,
	}

	return c
}
