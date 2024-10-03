package repo

import (
	"errors"
	"github.com/SmirnovND/metrics/internal/domain"
	"sync"
)

type MemStorage struct {
	collection map[string]domain.Metric
	mu         sync.Mutex
}

func NewMetricRepo() *MemStorage {
	return &MemStorage{
		collection: make(map[string]domain.Metric),
	}
}

func (m *MemStorage) AddMetric(metric domain.Metric) {
	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления
	m.collection[metric.GetName()] = metric
}

func (m *MemStorage) GetMetric(name string) (domain.Metric, error) {
	m.mu.Lock()         // Блокируем на чтение, так как могут быть конкурентные записи
	defer m.mu.Unlock() // Освобождаем доступ после чтения
	v, ok := m.collection[name]
	if !ok {
		return nil, errors.New("no data for the metric")
	}
	return v, nil
}
