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

func NewMetricRepo() *MemStorage {
	return &MemStorage{
		collection: make(map[string]domain.MetricInterface),
	}
}

func (m *MemStorage) UpdateMetric(metricNew domain.MetricInterface) {
	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления

	metricCurrent, ok := m.collection[metricNew.GetName()]
	if ok && metricCurrent.GetType() == domain.MetricTypeCounter {
		if currentValue, ok := metricCurrent.GetValue().(int64); ok {
			if newValue, ok := metricNew.GetValue().(int64); ok {
				metricCurrent.SetValue(currentValue + newValue)
			} else {
				m.collection[metricNew.GetName()] = metricNew
			}
		} else {
			m.collection[metricNew.GetName()] = metricNew
		}
	} else {
		m.collection[metricNew.GetName()] = metricNew
	}
}

func (m *MemStorage) GetMetric(name string) (domain.MetricInterface, error) {
	m.mu.RLock()         // Блокируем на чтение, так как могут быть конкурентные записи
	defer m.mu.RUnlock() // Освобождаем доступ после чтения
	v, ok := m.collection[name]
	if !ok {
		return nil, errors.New("no data for the metric")
	}
	return v, nil
}
