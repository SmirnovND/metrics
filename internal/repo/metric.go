package repo

import (
	"errors"
	"github.com/SmirnovND/metrics/internal/domain"
	"sync"
)

type MemStorage struct {
	collection map[string]domain.Metric
	mu         sync.RWMutex
}

func NewMetricRepo() *MemStorage {
	return &MemStorage{
		collection: make(map[string]domain.Metric),
	}
}

func (m *MemStorage) UpdateMetric(metricNew domain.Metric) {
	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления
	metricCurrent, ok := m.collection[metricNew.GetName()]
	if ok && metricCurrent.GetType() == domain.MetricTypeCounter {
		m.collection[metricCurrent.GetName()].SetValue(metricCurrent.GetValue().(int64) + metricNew.GetValue().(int64))
	} else {
		m.collection[metricNew.GetName()] = metricNew
	}
}

func (m *MemStorage) GetMetric(name string) (domain.Metric, error) {
	m.mu.RLock()         // Блокируем на чтение, так как могут быть конкурентные записи
	defer m.mu.RUnlock() // Освобождаем доступ после чтения
	v, ok := m.collection[name]
	if !ok {
		return nil, errors.New("no data for the metric")
	}
	return v, nil
}
