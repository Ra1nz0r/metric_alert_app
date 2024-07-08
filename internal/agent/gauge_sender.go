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
)

var metricURL = "http://localhost:8080/update/%s/%s/%d"

func SendGaugeOnServer(reportInterval, pollInterval time.Duration) {
	var nameMetric sync.Map
	cnt := 1

	c := make(chan os.Signal, 1)
	ticker := time.NewTicker(pollInterval * time.Second)

	d := make(chan os.Signal, 1)
	ticker1 := time.NewTicker(reportInterval * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				countGauge(&nameMetric, cnt)
				cnt++
			case <-ticker1.C:
				mapPostSender(&nameMetric, metricURL)
			}
		}
	}()

	<-c
	<-d
}

func mapPostSender(s *sync.Map, url string) {

	//var res *http.Response

	s.Range(func(k, v any) bool {
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
			return false
		}
		defer res.Body.Close()

		//data, err := io.ReadAll(res.Body) // удалить потом
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println(string(data))

		return true
	})

	//data, err := io.ReadAll(res.Body) // удалить потом
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(string(data)) // <- до сюда
	//defer res.Body.Close()
}

func countGauge(nameMetric *sync.Map, cnt int) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	nameMetric.Store("BuckHashSys", rtm.BuckHashSys)
	nameMetric.Store("Frees", rtm.Frees)
	nameMetric.Store("GCCPUFraction", int(rtm.GCCPUFraction))
	nameMetric.Store("GCSys", rtm.GCSys)
	nameMetric.Store("HeapAlloc", rtm.HeapAlloc)
	nameMetric.Store("HeapIdle", rtm.HeapIdle)
	nameMetric.Store("HeapInuse", rtm.HeapInuse)
	nameMetric.Store("HeapObjects", rtm.HeapObjects)
	nameMetric.Store("HeapReleased", rtm.HeapReleased)
	nameMetric.Store("HeapSys", rtm.HeapSys)
	nameMetric.Store("LastGC", rtm.LastGC)
	nameMetric.Store("Lookups", rtm.Lookups)
	nameMetric.Store("MCacheInuse", rtm.MCacheInuse)
	nameMetric.Store("MCacheSys", rtm.MCacheSys)
	nameMetric.Store("MSpanInuse", rtm.MSpanInuse)
	nameMetric.Store("MSpanSys", rtm.MSpanSys)
	nameMetric.Store("Mallocs", rtm.Mallocs)
	nameMetric.Store("NextGC", rtm.NextGC)
	nameMetric.Store("NumForcedGC", rtm.NumForcedGC)
	nameMetric.Store("NumGC", rtm.NumGC)
	nameMetric.Store("OtherSys", rtm.OtherSys)
	nameMetric.Store("PauseTotalNs", rtm.PauseTotalNs)
	nameMetric.Store("StackInuse", rtm.StackInuse)
	nameMetric.Store("StackSys", rtm.StackSys)
	nameMetric.Store("Sys", rtm.Sys)
	nameMetric.Store("TotalAlloc", rtm.TotalAlloc)

	nameMetric.Store("PollCount", cnt)
	nameMetric.Store("RandomValue", randRange(-999, 999))
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
