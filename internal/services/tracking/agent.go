package tracking

import (
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/collector"
	"time"
)

func MetricsTracking(cf interfaces.ConfigAgent) {
	metrics := collector.NewMetrics()
	// Тикер для обновления метрик
	updateTicker := time.NewTicker(time.Second * time.Duration(cf.GetPollInterval()))
	defer updateTicker.Stop()

	// Тикер для отправки метрик
	sendTicker := time.NewTicker(time.Second * time.Duration(cf.GetReportInterval()))
	defer sendTicker.Stop()

	go func() {
		for range updateTicker.C {
			metrics.Update() // Обновляем метрики
		}
	}()

	go func() {
		for range sendTicker.C {
			metrics.Send(cf.GetServerHost()) // Отправляем метрики
		}
	}()

	select {}
}
