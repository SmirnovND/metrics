package collector

import (
	"errors"
	"fmt"
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
	s.storage.UpdateMetric(m)
}

func (s *ServiceCollector) GetMetricValue(nameMetric string, typeMetric string) (string, error) {
	metric, err := s.FindMetric(nameMetric, typeMetric)
	if err != nil {
		return "", err
	}

	return s.formatValue(metric), nil
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

func (s *ServiceCollector) formatValue(metric domain.Metric) string {
	switch value := metric.GetValue().(type) {
	case int64:
		return s.formatInt(value)
	case float64:
		return s.formatFloat(value)
	default:
		return ""
	}
}

func (s *ServiceCollector) formatFloat(value float64) string {
	strValue := fmt.Sprintf("%f", value)

	// Убираем все нули в конце и, если необходимо, убираем точку
	trimmedValue := strValue
	for trimmedValue[len(trimmedValue)-1] == '0' {
		trimmedValue = trimmedValue[:len(trimmedValue)-1]
	}

	// Убираем точку, если она в конце
	if trimmedValue[len(trimmedValue)-1] == '.' {
		trimmedValue = trimmedValue[:len(trimmedValue)-1]
	}

	return trimmedValue
}

func (s *ServiceCollector) formatInt(value int64) string {
	return fmt.Sprintf("%d", value)
}
