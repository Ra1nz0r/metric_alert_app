package agent

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

type SenderStorage struct {
	sMS          storage.MetricService
	wg           sync.WaitGroup
	pollTicker   *time.Ticker
	reportTicker *time.Ticker
}

func NewSender(sMS storage.MetricService) *SenderStorage {
	return &SenderStorage{sMS: sMS}
}

// По указанному хосту, отправляет через POST запрос все метрики из локального хранилища на сервер.
func MapSender(host string, gauge *map[string]float64, counter *map[string]int64) {
	for k, v := range *gauge {
		resURL := fmt.Sprintf("http://%s/update/gauge/%s/%.2f", host, k, v)
		MakeRequest(resURL)
	}

	for k, v := range *counter {
		resURL := fmt.Sprintf("http://%s/update/counter/%s/%d", host, k, v)
		MakeRequest(resURL)
	}
}

// Создает POST запрос для отправки метрик на сервер по указанной ссылке.
func MakeRequest(resURL string) {
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

// Обновляет метрики из пакета runtime.
// Вызывает методы интерфейса хранилища, где реализовано взаимодействие и работа с ним.
func (s *SenderStorage) UpdateMetrics() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	s.sMS.UpdateGauge("Alloc", float64(rtm.Alloc))
	s.sMS.UpdateGauge("BuckHashSys", float64(rtm.BuckHashSys))
	s.sMS.UpdateGauge("Frees", float64(rtm.Frees))
	s.sMS.UpdateGauge("GCCPUFraction", rtm.GCCPUFraction)
	s.sMS.UpdateGauge("GCSys", float64(rtm.GCSys))
	s.sMS.UpdateGauge("HeapAlloc", float64(rtm.HeapAlloc))
	s.sMS.UpdateGauge("HeapIdle", float64(rtm.HeapIdle))
	s.sMS.UpdateGauge("HeapInuse", float64(rtm.HeapInuse))
	s.sMS.UpdateGauge("HeapObjects", float64(rtm.HeapObjects))
	s.sMS.UpdateGauge("HeapReleased", float64(rtm.HeapReleased))
	s.sMS.UpdateGauge("HeapSys", float64(rtm.HeapSys))
	s.sMS.UpdateGauge("LastGC", float64(rtm.LastGC))
	s.sMS.UpdateGauge("Lookups", float64(rtm.Lookups))
	s.sMS.UpdateGauge("MCacheInuse", float64(rtm.MCacheInuse))
	s.sMS.UpdateGauge("MCacheSys", float64(rtm.MCacheSys))
	s.sMS.UpdateGauge("MSpanInuse", float64(rtm.MSpanInuse))
	s.sMS.UpdateGauge("MSpanSys", float64(rtm.MSpanSys))
	s.sMS.UpdateGauge("Mallocs", float64(rtm.Mallocs))
	s.sMS.UpdateGauge("NextGC", float64(rtm.NextGC))
	s.sMS.UpdateGauge("NumForcedGC", float64(rtm.NumForcedGC))
	s.sMS.UpdateGauge("NumGC", float64(rtm.NumGC))
	s.sMS.UpdateGauge("OtherSys", float64(rtm.OtherSys))
	s.sMS.UpdateGauge("PauseTotalNs", float64(rtm.PauseTotalNs))
	s.sMS.UpdateGauge("StackInuse", float64(rtm.StackInuse))
	s.sMS.UpdateGauge("StackSys", float64(rtm.StackSys))
	s.sMS.UpdateGauge("Sys", float64(rtm.Sys))
	s.sMS.UpdateGauge("TotalAlloc", float64(rtm.TotalAlloc))
	s.sMS.UpdateGauge("RandomValue", randRange(-999, 999))

	s.sMS.UpdateCounter("PollCount", int64((rand.IntN(20-1) + 1)))
}

// Простой генератор рандомного числа.
func randRange(min, max int) float64 {
	return float64(rand.IntN(max-min) + min)
}
