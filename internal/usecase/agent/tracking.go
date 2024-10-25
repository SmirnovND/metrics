package agent

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/agent"
	"time"
)

func MetricsTracking(cf interfaces.ConfigAgent) {
	metrics := domain.NewMetrics()
	// Тикер для обновления метрик
	updateTicker := time.NewTicker(time.Second * time.Duration(cf.GetPollInterval()))
	defer updateTicker.Stop()

	go func() {
		for range updateTicker.C {
			agent.Update(metrics) // Обновляем метрики
		}
	}()

	// Тикер для отправки метрик
	sendTicker := time.NewTicker(time.Second * time.Duration(cf.GetReportInterval()))
	defer sendTicker.Stop()

	go func() {
		for range sendTicker.C {
			agent.SendJSON(metrics, cf.GetServerHost())
		}
	}()

	select {}
}
