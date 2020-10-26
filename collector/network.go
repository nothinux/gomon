package collector

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type NetworkInfo struct {
	TxBytes uint64 `json:"txbytes"`
	RxBytes uint64 `json:"rxbytes"`
}

func GetNetworks() map[string]string {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return ParseNetwork(file)

}

func ParseNetwork(file *os.File) map[string]string {
	netstats := make(map[string]string)

	scanner := bufio.NewScanner(file)
	// skip header
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())

		key := strings.TrimRight(parts[0], ":")
		rxBytes := toUint64(parts[1])
		txBytes := toUint64(parts[9])

		nets := NetworkInfo{
			TxBytes: txBytes,
			RxBytes: rxBytes,
		}

		n, _ := json.Marshal(nets)

		netstats[key] = string(n)
	}

	return netstats
}
