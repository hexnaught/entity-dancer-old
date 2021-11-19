package debugg

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var m runtime.MemStats

// PrintMemUsage ...
func PrintMemUsage() {
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v KiB", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KiB", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KiB", bToKb(m.Sys))
	fmt.Printf("\tMallocs = %v", m.Mallocs)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

	fmt.Printf("HeapInUse = %v KiB", bToKb(m.HeapInuse))
	fmt.Printf("\tHeapIdle = %v KiB", bToKb(m.HeapIdle))
	fmt.Printf("\tHeapReleased = %v KiB\n", bToKb(m.HeapReleased))

	debug.FreeOSMemory()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
