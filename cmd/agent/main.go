package main

import (
	config "github.com/SmirnovND/metrics/internal/pkg/config/agent"
	"github.com/SmirnovND/metrics/internal/usecase/agent"
)

func main() {
	cf := config.NewConfigCommand()
	agent.MetricsTracking(cf)
}
