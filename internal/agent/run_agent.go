package agent

import (
	"time"

	"github.com/ra1nz0r/metric_alert_app/internal/flags"
)

func RunAgent() {
	flags.AgentFlags()
	SendGaugeOnServer(time.Duration(flags.DefReportInterval), time.Duration(flags.DefPollInterval))
}
