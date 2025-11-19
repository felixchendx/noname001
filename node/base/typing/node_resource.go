package typing

type TempNodeSystemResourceSummary struct {
	CPUPercent float64

	MemoryTotal       uint64
	MemoryAvailable   uint64
	MemoryUsed        uint64
	MemoryUsedPercent float64
}
