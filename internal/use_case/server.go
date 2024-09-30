package use_case

import (
	"github.com/SmirnovND/metrics/domain"
	"github.com/SmirnovND/metrics/internal/services/metricscollector"
)

func ProcessMetrics(m domain.Metric) {
	//тут будет бизнес логика, если она будет и какие то сценарии
	metricscollector.SaveMetric(m)
}
