package agent

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/flags"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

func SendGaugeOnServer(reportInterval, pollInterval time.Duration) {
	var m storage.MetricService = storage.New()

	c := make(chan os.Signal, 1)
	pollTicker := time.NewTicker(pollInterval * time.Second)

	d := make(chan os.Signal, 1)
	reportTicker := time.NewTicker(reportInterval * time.Second)

	go func() {
		for {
			select {
			case <-pollTicker.C:
				updateMetrics(m)
			case <-reportTicker.C:
				mapSender(m.MakeStorageCopy())
			}
		}
	}()

	<-c
	<-d
}

func mapSender(gauge *map[string]float64, counter *map[string]int64) {
	for k, v := range *gauge {
		resURL := fmt.Sprintf("http://%s/update/gauge/%s/%.2f", flags.DefServerAddress, k, v)
		makeRequest(resURL)
	}

	for k, v := range *counter {
		resURL := fmt.Sprintf("http://%s/update/counter/%s/%d", flags.DefServerAddress, k, v)
		makeRequest(resURL)
	}
}

func updateMetrics(s storage.MetricService) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	s.UpdateGauge("Alloc", float64(rtm.Alloc))
	s.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	s.UpdateGauge("Frees", float64(rtm.Frees))
	s.UpdateGauge("GCCPUFraction", rtm.GCCPUFraction)
	s.UpdateGauge("GCSys", float64(rtm.GCSys))
	s.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	s.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	s.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	s.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	s.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	s.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	s.UpdateGauge("LastGC", float64(rtm.LastGC))
	s.UpdateGauge("Lookups", float64(rtm.Lookups))
	s.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	s.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	s.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	s.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
	s.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	s.UpdateGauge("NextGC", float64(rtm.NextGC))
	s.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	s.UpdateGauge("NumGC", float64(rtm.NumGC))
	s.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	s.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	s.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	s.UpdateGauge("StackSys", float64(rtm.StackSys))
	s.UpdateGauge("Sys", float64(rtm.Sys))
	s.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	s.UpdateGauge("RandomValue", randRange(-999, 999))

	s.UpdateCounter("PollCount", int64((rand.IntN(20-1) + 1)))
}

func makeRequest(resURL string) {
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

func randRange(min, max int) float64 {
	return float64(rand.IntN(max-min) + min)
}
