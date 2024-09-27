package metricscollector

import "github.com/SmirnovND/metrics/domain"

func ProcessMetrics(m domain.Metric) {
	memStorage := &domain.MemStorage{}
	memStorage.AddMetric(m)
}
