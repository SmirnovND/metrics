package repo

import (
	"errors"
	"github.com/SmirnovND/metrics/internal/domain"
	"sync"
)

type MemStorage struct {
	collection map[string]domain.MetricInterface
	mu         sync.RWMutex
}

func NewMetricRepo(collection map[string]domain.MetricInterface) *MemStorage {
	return &MemStorage{
		collection: collection,
	}
}

func (m *MemStorage) UpdateMetric(metricNew domain.MetricInterface) {
	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления

	metricCurrent, ok := m.collection[metricNew.GetName()+metricNew.GetType()]
	if ok && metricCurrent.GetType() == domain.MetricTypeCounter {
		var currentValue int64
		if val, ok := metricCurrent.GetValue().(int64); ok {
			currentValue = val
		}

		// Обработка значения metricNew
		var newValue int64
		if val, ok := metricNew.GetValue().(int64); ok {
			newValue = val
		}
		setValue := currentValue + newValue
		// Устанавливаем новое значение, если оба значения обработаны
		metricCurrent.SetValue(setValue)
	} else {
		m.collection[metricNew.GetName()+metricNew.GetType()] = metricNew
	}
}

func (m *MemStorage) GetMetric(name string, typeMetric string) (domain.MetricInterface, error) {
	m.mu.RLock()         // Блокируем на чтение, так как могут быть конкурентные записи
	defer m.mu.RUnlock() // Освобождаем доступ после чтения
	v, ok := m.collection[name+typeMetric]
	if !ok {
		return nil, errors.New("no data for the metric")
	}
	return v, nil
}

func (m *MemStorage) ExecuteWithLock(fn func(collection map[string]domain.MetricInterface)) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	fn(m.collection)
}
