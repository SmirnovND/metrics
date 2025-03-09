package server

import (
	"encoding/json"
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockServiceCollector struct {
	mock.Mock
}

func (ms *MockServiceCollector) SaveMetric(m domain.MetricInterface) {
	ms.Called(m)
}

func (ms *MockServiceCollector) GetMetricValue(nameMetric string, typeMetric string) (string, error) {
	args := ms.Called(nameMetric, typeMetric)
	return args.String(0), args.Error(1)
}

func (ms *MockServiceCollector) FindMetric(nameMetric string, typeMetric string) (domain.MetricInterface, error) {
	args := ms.Called(nameMetric, typeMetric)
	return args.Get(0).(domain.MetricInterface), args.Error(1)
}

func TestSaveAndFind(t *testing.T) {
	// Подготавливаем mock-объекты
	mockService := new(MockServiceCollector)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1")

	// Мокаем методы
	mockService.On("SaveMetric", metric).Return(nil)
	mockService.On("FindMetric", metric.GetName(), metric.GetType()).Return(metric, nil)

	// Вызов функции SaveAndFind
	result, err := SaveAndFind(metric, mockService, nil)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что результат содержит имя метрики
	assert.Contains(t, string(result), "metric1")

	// Проверяем, что методы SaveMetric и FindMetric были вызваны с нужными аргументами
	mockService.AssertExpectations(t)
}

func TestSaveAndFindArr(t *testing.T) {
	// Подготавливаем mock-объекты
	mockService := new(MockServiceCollector)

	metrics := []*domain.Metric{
		(&domain.Metric{}).SetType("counter").SetName("metric1").(*domain.Metric),
		(&domain.Metric{}).SetType("gauge").SetName("metric2").(*domain.Metric),
	}

	// Мокаем методы
	mockService.On("SaveMetric", metrics[0]).Return(nil)
	mockService.On("SaveMetric", metrics[1]).Return(nil)
	mockService.On("FindMetric", metrics[0].GetName(), metrics[0].GetType()).Return(metrics[0], nil)
	mockService.On("FindMetric", metrics[1].GetName(), metrics[1].GetType()).Return(metrics[1], nil)

	// Вызов функции SaveAndFindArrHTTP
	// Мы не будем вызывать HTTP-запрос, а просто проверим вызовы функций
	result, err := SaveAndFindArrHTTP(metrics, mockService, nil)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что результат содержит имена метрик
	var response []domain.Metric
	err = json.Unmarshal(result, &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Contains(t, response[0].GetName(), "metric1")
	assert.Contains(t, response[1].GetName(), "metric2")

	// Проверяем, что методы SaveMetric и FindMetric были вызваны с нужными аргументами
	mockService.AssertExpectations(t)
}

func TestFindAndResponseAsJSON(t *testing.T) {
	// Подготавливаем mock-объекты
	mockService := new(MockServiceCollector)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1")

	// Мокаем методы
	mockService.On("FindMetric", metric.GetName(), metric.GetType()).Return(metric, nil)

	// Вызов функции FindAndResponseAsJSON
	// Мы не будем вызывать HTTP-запрос, а просто проверим вызовы функций
	result, err := FindAndResponseAsJSON(metric, mockService, nil)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что результат содержит имя метрики
	assert.Contains(t, string(result), "metric1")

	// Проверяем, что метод FindMetric был вызван с нужными аргументами
	mockService.AssertExpectations(t)
}

func TestFind(t *testing.T) {
	// Подготавливаем mock-объекты
	mockService := new(MockServiceCollector)
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1")

	// Мокаем методы
	mockService.On("FindMetric", metric.GetName(), metric.GetType()).Return(metric, nil)

	// Вызов функции Find
	// Мы не будем вызывать HTTP-запрос, а просто проверим вызовы функций
	result, err := Find(metric, mockService, nil)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что результат содержит имя метрики
	assert.Equal(t, result.GetName(), "metric1")

	// Проверяем, что метод FindMetric был вызван с нужными аргументами
	mockService.AssertExpectations(t)
}

func TestJSONResponse(t *testing.T) {
	// Подготавливаем данные для JSON-сериализации
	metric := &domain.Metric{}
	metric.SetType("counter").SetName("metric1")

	// Создаем ResponseWriter (не используется в этом тесте, но требуется функцией)
	rr := httptest.NewRecorder()

	// Вызов функции JSONResponse
	result, err := JSONResponse(metric, rr)

	// Проверяем, что ошибок не возникло
	assert.NoError(t, err)

	// Проверяем, что результат сериализован в JSON
	assert.Contains(t, string(result), "metric1")

	// Проверяем, что в ответе есть корректный код
	assert.Equal(t, http.StatusOK, rr.Code)
}
