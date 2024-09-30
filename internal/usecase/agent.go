package usecase

import (
	"github.com/SmirnovND/metrics/internal/pkg/config"
	"github.com/SmirnovND/metrics/internal/services/metricscollector"
	"time"
)

func TrackingMetrics(cf config.Config) {
	metrics := metricscollector.NewMetrics()
	// Тикер для обновления метрик
	updateTicker := time.NewTicker(time.Second * time.Duration(cf.PollInterval))
	defer updateTicker.Stop()

	// Тикер для отправки метрик
	sendTicker := time.NewTicker(time.Second * time.Duration(cf.ReportInterval))
	defer sendTicker.Stop()

	go func() {
		for range updateTicker.C {
			metrics.Update() // Обновляем метрики
		}
	}()

	go func() {
		for range sendTicker.C {
			metrics.Send(cf.ServerHost) // Отправляем метрики
		}
	}()

	select {}
}
