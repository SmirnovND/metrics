package interfaces

import "github.com/SmirnovND/metrics/internal/domain"

// ServiceCollectorInterface описывает методы для работы с метриками.
type ServiceCollectorInterface interface {
	// SaveMetric сохраняет метрику в хранилище.
	SaveMetric(m domain.MetricInterface)

	// GetMetricValue получает значение метрики по имени и типу.
	GetMetricValue(nameMetric string, typeMetric string) (string, error)

	// FindMetric ищет метрику по имени и типу.
	FindMetric(nameMetric string, typeMetric string) (domain.MetricInterface, error)
}
