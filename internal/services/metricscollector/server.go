package metricscollector

import "github.com/SmirnovND/metrics/domain"

// а тут реализация функциональности, без бизнес логики
func SaveMetric(m domain.Metric) {
	memStorage := &domain.MemStorage{}
	memStorage.AddMetric(m)
}