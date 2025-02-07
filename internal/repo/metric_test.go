package repo

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage_UpdateMetric(t *testing.T) {
	// Создаем исходные данные для теста
	metricRepo := NewMetricRepo(make(map[string]domain.MetricInterface))

	metric := &domain.Metric{}
	metric.SetType(domain.MetricTypeCounter).SetValue(int64(5)).SetName("test_metric")

	metricRepo.UpdateMetric(metric)

	// Получаем метрику из хранилища
	updatedMetric, err := metricRepo.GetMetric("test_metric", domain.MetricTypeCounter)
	assert.NoError(t, err)

	// Проверяем значение метрики
	assert.Equal(t, int64(5), updatedMetric.GetValue())

	metricNew := &domain.Metric{}
	metricNew.SetType(domain.MetricTypeCounter).SetValue(int64(3)).SetName("test_metric")

	metricRepo.UpdateMetric(metricNew)

	// Получаем обновленную метрику
	updatedMetric, err = metricRepo.GetMetric("test_metric", domain.MetricTypeCounter)
	assert.NoError(t, err)

	// Проверяем, что значение метрики обновилось правильно (5 + 3 = 8)
	assert.Equal(t, int64(8), updatedMetric.GetValue())
}

func TestMemStorage_GetMetric(t *testing.T) {
	metricRepo := NewMetricRepo(make(map[string]domain.MetricInterface))

	// Попытка получить метрику, которая не существует
	metric, err := metricRepo.GetMetric("non_existent_metric", domain.MetricTypeCounter)
	assert.Error(t, err)
	assert.Nil(t, metric)

	metricNew := &domain.Metric{}
	metricNew.SetType(domain.MetricTypeCounter).SetValue(int64(10)).SetName("existing_metric")
	metricRepo.UpdateMetric(metricNew)

	// Получаем метрику из хранилища
	metric, err = metricRepo.GetMetric("existing_metric", domain.MetricTypeCounter)
	assert.NoError(t, err)
	assert.Equal(t, int64(10), metric.GetValue())
}

//func TestMemStorage_ExecuteWithLock(t *testing.T) {
//	metricRepo := NewMetricRepo(make(map[string]domain.MetricInterface))
//
//	metric := &domain.Metric{}
//	metric.SetType(domain.MetricTypeCounter).SetValue(10).SetName("metric_for_lock_test")
//	metricRepo.UpdateMetric(metric)
//
//	// Выполняем функцию с блокировкой
//	metricRepo.ExecuteWithLock(func(collection map[string]domain.MetricInterface) {
//		// Проверяем наличие метрики в коллекции внутри заблокированного кода
//		assert.Contains(t, collection, "metric_for_lock_test"+domain.MetricTypeCounter)
//		// Проверяем значение метрики
//		assert.Equal(t, int64(10), collection["metric_for_lock_test"+domain.MetricTypeCounter].GetValue())
//	})
//
//	// После выполнения блока функции, метрика должна остаться
//	metric, err := metricRepo.GetMetric("metric_for_lock_test", domain.MetricTypeCounter)
//	assert.NoError(t, err)
//	assert.Equal(t, int64(10), metric.GetValue())
//}
