package server_test

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/SmirnovND/metrics/internal/services/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockMemStorage struct {
	mock.Mock
}

func (m *MockMemStorage) UpdateMetric(metricNew domain.MetricInterface) {
	m.Called(metricNew)
}

func (m *MockMemStorage) GetMetric(name string, typeMetric string) (domain.MetricInterface, error) {
	args := m.Called(name, typeMetric)
	return args.Get(0).(domain.MetricInterface), args.Error(1)
}

func (m *MockMemStorage) ExecuteWithLock(f func(map[string]domain.MetricInterface)) {
	args := m.Called(f)
	f(args.Get(0).(map[string]domain.MetricInterface))
}

func TestServiceCollector_SaveMetric(t *testing.T) {
	// Мокируем хранилище
	mockStorage := new(MockMemStorage)
	metric := &domain.Metric{
		ID:    "metric1",
		MType: "counter",
		Value: 10,
	}

	// Ожидаем вызов метода UpdateMetric
	mockStorage.On("UpdateMetric", metric).Once()

	// Создаем сервис
	service := server.NewCollectorService(mockStorage)

	// Вызываем метод SaveMetric
	service.SaveMetric(metric)

	// Проверяем, что метод был вызван
	mockStorage.AssertExpectations(t)
}

func TestServiceCollector_GetMetricValue(t *testing.T) {
	// Мокируем хранилище
	mockStorage := new(MockMemStorage)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1").SetValue(int64(10))

	// Ожидаем вызов метода GetMetric
	mockStorage.On("GetMetric", "metric1", "counter").Return(metric, nil)

	// Создаем сервис
	service := server.NewCollectorService(mockStorage)

	// Проверяем метод GetMetricValue
	value, err := service.GetMetricValue("metric1", "counter")
	assert.NoError(t, err)
	assert.Equal(t, "10", value)

	// Проверяем, что метод был вызван
	mockStorage.AssertExpectations(t)
}

func TestServiceCollector_FindMetric(t *testing.T) {
	// Мокируем хранилище
	mockStorage := new(MockMemStorage)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1").SetValue(int64(10))

	// Ожидаем вызов метода GetMetric
	mockStorage.On("GetMetric", "metric1", "counter").Return(metric, nil)

	// Создаем сервис
	service := server.NewCollectorService(mockStorage)

	// Проверяем метод FindMetric
	result, err := service.FindMetric("metric1", "counter")
	assert.NoError(t, err)
	assert.Equal(t, "metric1", result.GetName())

	// Проверяем, что метод был вызван
	mockStorage.AssertExpectations(t)
}

func TestServiceCollector_formatValue(t *testing.T) {
	// Мокируем хранилище
	mockStorage := new(MockMemStorage)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1").SetValue(int64(10))

	// Ожидаем вызов метода GetMetric
	mockStorage.On("GetMetric", "metric1", "counter").Return(metric, nil)

	// Создаем сервис
	service := server.NewCollectorService(mockStorage)

	// Проверяем метод formatValue для int64
	formattedValue, _ := service.GetMetricValue("metric1", "counter")
	assert.Equal(t, "10", formattedValue)

	// Проверяем форматирование int64
	metric.SetValue(int64(123456))
	formattedValue, _ = service.GetMetricValue("metric1", "counter")
	assert.Equal(t, "123456", formattedValue)

	// Проверяем форматирование float64
	metric.SetValue(int64(10))
	formattedValue, _ = service.GetMetricValue("metric1", "counter")
	assert.Equal(t, "10", formattedValue)
}
