package interfaces

import "github.com/SmirnovND/metrics/internal/domain"

// MemStorageInterface определяет интерфейс для работы с хранилищем метрик.
type MemStorageInterface interface {
	// Обновляет метрику в хранилище.
	UpdateMetric(metricNew domain.MetricInterface)

	// Получает метрику по имени и типу.
	GetMetric(name string, typeMetric string) (domain.MetricInterface, error)

	// Выполняет переданную функцию с блокировкой для чтения коллекции метрик.
	ExecuteWithLock(fn func(collection map[string]domain.MetricInterface))
}
