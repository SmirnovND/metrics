package domain

type MemStorage struct {
	collection []Metric
}

func (m *MemStorage) AddMetric(metric Metric) {
	m.collection = append(m.collection, metric)
}
