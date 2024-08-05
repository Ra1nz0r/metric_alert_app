package agent

import (
	"sync"
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/config"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

// Запускает агент, который будет через указанное время обновлять метрики
// в локальном хранилище и отправлять их на сервер.
func RunAgent() {
	// Запускаем флаги и переменные окружения для агента.
	config.AgentFlags()

	// Создаем интерфейс и новое хранилище.
	ss := NewSender(storage.New())

	// Обновляем и отправляем метрики на сервер.
	var wg sync.WaitGroup

	pollTicker := time.NewTicker(time.Duration(config.DefPollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(config.DefReportInterval) * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			ss.SendMetricsOnServer(reportTicker, pollTicker)

		}

	}()
	wg.Wait()
}
