package controllers

import (
	"context"
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
	"github.com/SmirnovND/metrics/internal/pkg/paramsparser"
	"github.com/SmirnovND/metrics/internal/repo"
	serverSaver "github.com/SmirnovND/metrics/internal/usecase/server"
	"github.com/SmirnovND/metrics/pb"
	"github.com/jmoiron/sqlx"
)

func NewServiceServer(
	Storage *repo.MemStorage,
	DB *sqlx.DB,
	Cf interfaces.ConfigServerInterface,
	ServiceCollector interfaces.ServiceCollectorInterface,
) *ServiceServer {
	return &ServiceServer{
		DB:               DB,
		Storage:          Storage,
		Cf:               Cf,
		ServiceCollector: ServiceCollector,
	}
}

type ServiceServer struct {
	pb.UnimplementedMetricsServiceServer
	Storage          *repo.MemStorage
	DB               *sqlx.DB
	Cf               interfaces.ConfigServerInterface
	ServiceCollector interfaces.ServiceCollectorInterface
}

func (s *ServiceServer) SendMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
	// Разбор метрик из запроса
	metrics := paramsparser.ConvertPbToDomain(req.Metrics)

	// Сохранение и поиск обновленных метрик
	err := s.saveAndFindMetrics(metrics, s.ServiceCollector)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении и поиске метрик: %v", err)
	}
	fmt.Println("Метрики успешно сохранены")
	// Формирование ответа с обновленными метриками
	return &pb.MetricsResponse{}, nil
}

// saveAndFindMetrics сохраняет метрики и находит их в хранилище
func (s *ServiceServer) saveAndFindMetrics(metrics []*domain.Metric, ServiceCollector interfaces.ServiceCollectorInterface) error {
	err := serverSaver.SaveAndFindArrGRPC(metrics, ServiceCollector)
	if err != nil {
		return fmt.Errorf("ошибка при сохранении метрик: %v", err)
	}

	return nil
}
