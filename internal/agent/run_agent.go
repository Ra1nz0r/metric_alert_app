package agent

import (
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

	ss.pollTicker = time.NewTicker(time.Duration(config.DefPollInterval) * time.Second)
	ss.reportTicker = time.NewTicker(time.Duration(config.DefReportInterval) * time.Second)

	ss.wg.Add(1)
	go func() {
		defer ss.wg.Done()
		for {
			select {
			case <-ss.pollTicker.C:
				ss.UpdateMetrics()
			case <-ss.reportTicker.C:
				g, c := ss.sMS.MakeStorageCopy()
				MapSender(config.DefServerHost, g, c)
			}
		}
	}()
	ss.wg.Wait()
}
