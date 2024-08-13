package config

import (
	"flag"
	"os"
	"strconv"
)

var (
	DefServerHost     = "0.0.0.0:8080" // стандартный адрес для агента и сервера
	DefReportInterval = 10             // стандартная частота отправки метрик на сервер для агента в секундах
	DefPollInterval   = 2              // стандартная частоты опроса метрик для агента в секундах
	DefLogLevel       = "info"         // стандартный уровень логирования
)

// Создаёт флаги для запуска агента, если в терминале переданы переменные окружения,
// то приоритет будет отдаваться им.
func AgentFlags() {
	flag.StringVar(&DefServerHost, "a", DefServerHost, "address and port to run server/agent")
	flag.IntVar(&DefReportInterval, "r", DefReportInterval, "changing the frequency of sending metrics to the server for the agent (in seconds)")
	flag.IntVar(&DefPollInterval, "p", DefPollInterval, "changing the metric polling frequency for the agent (in seconds)")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		DefServerHost = envServerAddress
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		if value, err := strconv.Atoi(envReportInterval); err == nil {
			DefReportInterval = value
		}
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		if value, err := strconv.Atoi(envPollInterval); err == nil {
			DefPollInterval = value
		}
	}
}

// Создаёт флаги для запуска сервера, если в терминале переданы переменные окружения,
// то приоритет будет отдаваться им.
func ServerFlags() {
	flag.StringVar(&DefServerHost, "a", DefServerHost, "address and port to run server/agent")
	flag.StringVar(&DefLogLevel, "l", DefLogLevel, "set log level")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		DefServerHost = envServerAddress
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		DefLogLevel = envLogLevel
	}
}
