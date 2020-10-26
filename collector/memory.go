package collector

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

var MemStats = regexp.MustCompile("^(MemTotal|MemFree|MemAvailable|Buffers|Cached)$")

func GetMemory() (map[string]uint64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseMemory(file)
}

func ParseMemory(file *os.File) (map[string]uint64, error) {
	meminfo := make(map[string]uint64)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		key := strings.TrimRight(parts[0], ":")

		if MemStats.Match([]byte(key)) {
			meminfo[key] = toUint64(parts[1]) * 1024
		}

	}

	meminfo["Used"] = meminfo["MemTotal"] - meminfo["MemFree"] - meminfo["Buffers"] - meminfo["Cached"]

	return meminfo, nil
}
