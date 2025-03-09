package agent

import (
	"context"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/services/agent"
	"github.com/SmirnovND/metrics/pb"
	"google.golang.org/grpc"
	"log"
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
				var sender agent.Sender
				if cf.IsUseGRPC() {
					conn, err := grpc.Dial(cf.GetGRPCServerHost(), grpc.WithInsecure()) // Подключение к gRPC серверу
					if err != nil {
						log.Fatalf("Ошибка при подключении к gRPC: %v", err)
					}
					defer conn.Close()
					sender = &agent.GRPCSender{Client: pb.NewMetricsServiceClient(conn)}
				} else {
					sender = &agent.HTTPSender{ServerHost: cf.GetServerHost(), Key: cf.GetKey()}
				}
				agent.SendJSON(metrics, sender)
			}
		}
	}()

	<-ctx.Done() // Ожидаем завершения контекста
	fmt.Println("Все процессы завершены. Агент остановлен.")

	wg.Wait() // Ждем завершения всех горутин
}
