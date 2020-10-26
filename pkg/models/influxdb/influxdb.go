package influxdb

import (
	"github.com/influxdata/influxdb-client-go/api"
)

func Write(writeApi api.WriteAPI, line string) {
	writeApi.WriteRecord(line)
}
