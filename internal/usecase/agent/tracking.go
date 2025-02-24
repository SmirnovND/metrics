package agent

import (
	"context"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/agent"
	"sync"
	"time"
)

// MetricsTracking запускает процесс отслеживания и отправки метрик.
// - Метрики обновляются с интервалом `PollInterval`.
// - Отправка метрик на сервер выполняется с интервалом `ReportInterval`.
// - Функция работает бесконечно, используя `select {}` для блокировки выполнения.
func MetricsTracking(ctx context.Context, cf interfaces.ConfigAgent) {
	metrics := domain.NewMetrics()
	var wg sync.WaitGroup

	// Тикер для обновления метрик
	updateTicker := time.NewTicker(time.Second * time.Duration(cf.GetPollInterval()))
	defer updateTicker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершаем обновление базовых метрик")
				return
			case <-updateTicker.C:
				agent.Update(metrics, agent.BaseMetric)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершаем обновление расширенных метрик")
				return
			case <-updateTicker.C:
				agent.Update(metrics, agent.AdvancedMetricsDefinitions)
			}
		}
	}()

	// Тикер для отправки метрик
	sendTicker := time.NewTicker(time.Second * time.Duration(cf.GetReportInterval()))
	defer sendTicker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Завершаем отправку метрик")
				return
			case <-sendTicker.C:
				agent.SendJSON(metrics, cf.GetServerHost(), cf.GetCryptoKey())
			}
		}
	}()

	<-ctx.Done() // Ожидаем завершения контекста
	fmt.Println("Все процессы завершены. Агент остановлен.")

	wg.Wait() // Ждем завершения всех горутин
}
