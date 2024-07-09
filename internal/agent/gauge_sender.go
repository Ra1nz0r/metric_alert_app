package agent

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

var metricURL = "http://localhost:8080/update/%s/%s/%d"

var mu sync.RWMutex

func SendGaugeOnServer(reportInterval, pollInterval time.Duration) {
	var m storage.MetricService = storage.New()

	_, counterMap := m.GetMap()

	cnt := 1

	c := make(chan os.Signal, 1)
	pollTicker := time.NewTicker(pollInterval * time.Second)

	d := make(chan os.Signal, 1)
	reportTicker := time.NewTicker(reportInterval * time.Second)

	go func() {
		for {
			select {
			case <-pollTicker.C:
				updateMetrics(*counterMap, int64(cnt))
				cnt++
			case <-reportTicker.C:
				mu.RLock()
				mapSender(metricURL, *counterMap)
				mu.RUnlock()
			}
		}
	}()

	<-c
	<-d
}

func mapSender(url string, m map[string]int64) {
	for k, v := range m {
		metType := "gauge"
		if k == "PollCount" {
			metType = "counter"
		}

		resURL := fmt.Sprintf(url, metType, k, v)

		req, err := http.NewRequest("POST", resURL, nil)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
	}
}

func updateMetrics(nameMetric map[string]int64, cnt int64) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	mu.Lock()

	nameMetric["Alloc"] = int64(rtm.Alloc)
	nameMetric["BuckHashSys"] = int64(rtm.BuckHashSys)
	nameMetric["Frees"] = int64(rtm.Frees)
	nameMetric["GCCPUFraction"] = int64(rtm.GCCPUFraction)
	nameMetric["GCSys"] = int64(rtm.GCSys)
	nameMetric["HeapAlloc"] = int64(rtm.HeapAlloc)
	nameMetric["HeapIdle"] = int64(rtm.HeapIdle)
	nameMetric["HeapInuse"] = int64(rtm.HeapInuse)
	nameMetric["HeapObjects"] = int64(rtm.HeapObjects)
	nameMetric["HeapReleased"] = int64(rtm.HeapReleased)
	nameMetric["HeapSys"] = int64(rtm.HeapSys)
	nameMetric["LastGC"] = int64(rtm.LastGC)
	nameMetric["Lookups"] = int64(rtm.Lookups)
	nameMetric["MCacheInuse"] = int64(rtm.MCacheInuse)
	nameMetric["MCacheSys"] = int64(rtm.MCacheSys)
	nameMetric["MSpanInuse"] = int64(rtm.MSpanInuse)
	nameMetric["MSpanSys"] = int64(rtm.MSpanSys)
	nameMetric["Mallocs"] = int64(rtm.Mallocs)
	nameMetric["NextGC"] = int64(rtm.NextGC)
	nameMetric["NumForcedGC"] = int64(rtm.NumForcedGC)
	nameMetric["NumGC"] = int64(rtm.NumGC)
	nameMetric["OtherSys"] = int64(rtm.OtherSys)
	nameMetric["PauseTotalNs"] = int64(rtm.PauseTotalNs)
	nameMetric["StackInuse"] = int64(rtm.StackInuse)
	nameMetric["StackSys"] = int64(rtm.StackSys)
	nameMetric["Sys"] = int64(rtm.Sys)
	nameMetric["TotalAlloc"] = int64(rtm.TotalAlloc)

	nameMetric["PollCount"] = cnt
	nameMetric["RandomValue"] = randRange(-999, 999)

	mu.Unlock()
}

func randRange(min, max int) int64 {
	return int64(rand.IntN(max-min) + min)
}
