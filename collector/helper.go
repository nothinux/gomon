package collector

import "strconv"

func toFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func toUint64(s string) uint64 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return u
}
