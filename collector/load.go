package collector

import (
	"io/ioutil"
	"strings"
)

func GetLoad() (map[string]float64, error) {
	file, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return nil, err
	}

	return ParseLoad(string(file)), nil
}

func ParseLoad(data string) map[string]float64 {
	loads := make(map[string]float64)
	parts := strings.Fields(data)

	loads["load1"] = toFloat(parts[0])
	loads["load5"] = toFloat(parts[1])
	loads["load15"] = toFloat(parts[2])
	return loads
}
