package agent

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/agent"
	"time"
)

// MetricsTracking запускает процесс отслеживания и отправки метрик.
// - Метрики обновляются с интервалом `PollInterval`.
// - Отправка метрик на сервер выполняется с интервалом `ReportInterval`.
// - Функция работает бесконечно, используя `select {}` для блокировки выполнения.
func MetricsTracking(cf interfaces.ConfigAgent) {
	metrics := domain.NewMetrics()

	// Тикер для обновления метрик
	updateTicker := time.NewTicker(time.Second * time.Duration(cf.GetPollInterval()))
	defer updateTicker.Stop()

	go func() {
		for range updateTicker.C {
			agent.Update(metrics, agent.BaseMetric) // Обновляем метрики
		}
	}()

	go func() {
		for range updateTicker.C {
			agent.Update(metrics, agent.AdvancedMetricsDefinitions) // Обновляем метрики
		}
	}()

	// Тикер для отправки метрик
	sendTicker := time.NewTicker(time.Second * time.Duration(cf.GetReportInterval()))
	defer sendTicker.Stop()

	go func() {
		for range sendTicker.C {
			agent.SendJSON(metrics, cf.GetServerHost(), cf.GetKey(), cf.GetCryptoKey())
		}
	}()

	select {}
}
