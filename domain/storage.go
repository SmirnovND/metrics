package domain

type MemStorage struct {
	collection []Metric
}

func (m *MemStorage) AddMetric(metric Metric) {
	//fmt.Println(metric)
	m.collection = append(m.collection, metric)
}
