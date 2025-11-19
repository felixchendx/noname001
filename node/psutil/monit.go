package psutil

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	nodeTyping "noname001/node/base/typing"
)

func TempSystemResourceSummary() (*nodeTyping.TempNodeSystemResourceSummary) {
	var summ = &nodeTyping.TempNodeSystemResourceSummary{}

	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		// logge

	} else {
		if len(cpuPercent) > 0 {
			summ.CPUPercent = cpuPercent[0]
		}
	}

	vmem, err := mem.VirtualMemory()
	if err != nil {
		// logge

	} else {
		summ.MemoryTotal       = vmem.Total
		summ.MemoryAvailable   = vmem.Available
		summ.MemoryUsed        = vmem.Used
		summ.MemoryUsedPercent = vmem.UsedPercent
	}

	return summ
}

// type t_monitor struct {
// 	ctx       context.Context
// 	ctxCancel context.CancelFunc
// 	logger    

// 	cpuMonitor *t_cpuMonitor
// 	memMonitor *t_memMonitor
// 	// TODO: disk and network
// }
