package agent

import (
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/config"
	"github.com/ra1nz0r/metric_alert_app/internal/storage"
)

// Запускает агент, который будет через указанное время обновлять метрики
// в локальном хранилище и отправлять их на сервер.
func RunAgent() {
	// Создаем интерфейс и новое хранилище.
	var ms storage.MetricService = storage.New()

	// Запускаем флаги и переменные окружения для агента.
	config.AgentFlags()

	ss := NewSender(ms)

	// Обновляем и отправляем метрики на сервер.
	ss.SendMetricsOnServer(time.Duration(config.DefReportInterval), time.Duration(config.DefPollInterval))
}
