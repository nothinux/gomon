package collector

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type CPUStats struct {
	PrevIdleTime  uint64
	PrevTotalTime uint64
}

func GetCPU(stat *CPUStats) float64 {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return ParseCPU(file, stat)
}

func ParseCPU(file *os.File, stat *CPUStats) float64 {
	//var prevIdleTime, prevTotalTime uint64
	var cpuStats float64

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	parts := strings.Fields(scanner.Text()[5:])

	idleTime := toUint64(parts[3])
	totalTime := uint64(0)

	for _, x := range parts {
		u := toUint64(x)
		totalTime = totalTime + u
	}

	deltaIdleTime := idleTime - stat.PrevIdleTime
	deltaTotalTime := totalTime - stat.PrevTotalTime

	stat.PrevIdleTime = idleTime
	stat.PrevTotalTime = totalTime

	cpuStats = (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0

	return cpuStats
}
