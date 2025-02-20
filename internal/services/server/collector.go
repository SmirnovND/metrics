package server

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/interfaces"
)

type ServiceCollector struct {
	storage interfaces.MemStorageInterface
}

func NewCollectorService(storage interfaces.MemStorageInterface) interfaces.ServiceCollectorInterface {
	return &ServiceCollector{
		storage: storage,
	}
}

func (s *ServiceCollector) SaveMetric(m domain.MetricInterface) {
	s.storage.UpdateMetric(m)
}

func (s *ServiceCollector) GetMetricValue(nameMetric string, typeMetric string) (string, error) {
	metric, err := s.FindMetric(nameMetric, typeMetric)
	if err != nil {
		return "", err
	}

	return s.formatValue(metric), nil
}

func (s *ServiceCollector) FindMetric(nameMetric string, typeMetric string) (domain.MetricInterface, error) {
	metric, err := s.storage.GetMetric(nameMetric, typeMetric)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *ServiceCollector) formatValue(metric domain.MetricInterface) string {
	switch value := metric.GetValue().(type) {
	case *int64:
		return s.formatInt(*value)
	case *float64:
		return s.formatFloat(*value)
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
