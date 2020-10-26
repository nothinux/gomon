package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	influxclient "github.com/influxdata/influxdb-client-go"
	influxapi "github.com/influxdata/influxdb-client-go/api"
	"github.com/nothinux/gomon/collector"
	"github.com/nothinux/gomon/pkg/models/influxdb"
)

func main() {
	config, err := getConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}
	db := InfluxConnect(config)
	defer db.Close()

	cron := gocron.NewScheduler(time.Local)

	// load
	cron.Every(15).Second().Do(func() {
		loads, err := collector.GetLoad()
		if err != nil {
			log.Fatal(err)
		}

		datapoint := fmt.Sprintf("loadavg,hostname=server-1 load1=%.2f,load5=%.2f,load15=%.2f", loads["load1"], loads["load5"], loads["load15"])

		influxdb.Write(db, datapoint)
		log.Println("load average datapoint has been written")
	})

	// storage
	cron.Every(15).Second().Do(func() {
		storages, err := collector.GetStorage()
		if err != nil {
			log.Fatal(err)
		}

		for storage := range storages {
			var s collector.StorageInfo

			if err := json.Unmarshal([]byte(storages[storage]), &s); err != nil {
				log.Fatal(err)
			}

			datapoint := fmt.Sprintf("diskusage,hostname=server-1,mountpoint=%s size=%d,avail=%d,used=%d", storage, s.Size, s.Avail, s.Used)
			influxdb.Write(db, datapoint)

		}
		log.Println("storage datapoint has been written")
	})

	// memory
	cron.Every(15).Second().Do(func() {
		memory, err := collector.GetMemory()
		if err != nil {
			log.Fatal(err)
		}

		datapoint := fmt.Sprintf("memory,hostname=server-1 memtotal=%d,memused=%d", memory["MemTotal"], memory["Used"])
		influxdb.Write(db, datapoint)
		log.Println("memory datapoint has been written")
	})

	// cpu
	cpustat := &collector.CPUStats{}
	cron.Every(1).Second().Do(func() {
		cputotal := collector.GetCPU(cpustat)

		datapoint := fmt.Sprintf("cpu,hostname=server-1 totaltime=%.1f", cputotal)
		influxdb.Write(db, datapoint)
		log.Println("cpu datapoint has been written")
	})

	cron.StartBlocking()
}

func InfluxConnect(config *Config) influxapi.WriteApi {
	conf := config.Database.Influxdb
	client := influxclient.NewClient(fmt.Sprintf("http://%s:%s", conf.Host, conf.Port), fmt.Sprintf("%s:%s", conf.Username, conf.Password))

	return client.WriteAPI("", fmt.Sprintf("%s/%s", conf.Database, conf.RetentionPolicy))
}
