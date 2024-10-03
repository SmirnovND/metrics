package collector

import (
	"errors"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/repo"
)

type ServiceCollector struct {
	storage *repo.MemStorage
}

func NewCollectorService(storage *repo.MemStorage) *ServiceCollector {
	return &ServiceCollector{
		storage: storage,
	}
}

// а тут реализация функциональности, без бизнес логики
func (s *ServiceCollector) SaveMetric(m domain.Metric) {
	s.storage.AddMetric(m)
}

func (s *ServiceCollector) FindMetric(nameMetric string, typeMetric string) (domain.Metric, error) {
	metric, err := s.storage.GetMetric(nameMetric)
	if err != nil {
		return nil, err
	}

	if metric.GetType() != typeMetric {
		return nil, errors.New("not found metric with this type")
	}

	return metric, nil
}
