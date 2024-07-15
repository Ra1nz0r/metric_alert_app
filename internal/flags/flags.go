package flags

import (
	"flag"
)

var (
	DefServerAddress  = "0.0.0.0:8080" // стандартный адрес для агента и сервера
	DefReportInterval = 10             // стандартная частота отправки метрик на сервер для агента
	DefPollInterval   = 2              // стандартная частоты опроса метрик для агента
)

func AgentFlags() {
	flag.StringVar(&DefServerAddress, "a", DefServerAddress, "address and port to run server/agent")
	flag.IntVar(&DefReportInterval, "r", DefReportInterval, "changing the frequency of sending metrics to the server for the agent")
	flag.IntVar(&DefPollInterval, "p", DefPollInterval, "changing the metric polling frequency for the agent")
	flag.Parse()
}

func ServerFlags() {
	flag.StringVar(&DefServerAddress, "a", DefServerAddress, "address and port to run server/agent")
	flag.Parse()
}
