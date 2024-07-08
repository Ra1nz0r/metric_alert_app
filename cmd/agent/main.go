package main

import "github.com/ra1nz0r/metric_alert_app/internal/agent"

func main() {
	agent.SendGaugeOnServer(10, 2)
}
