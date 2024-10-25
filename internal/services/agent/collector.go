package agent

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"math/rand"
	"runtime"
)

var metricDefinitions = map[string]struct {
	Type   string
	Update func(rtm *runtime.MemStats) interface{}
}{
	"Alloc":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Alloc) }},
	"BuckHashSys":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.BuckHashSys) }},
	"Frees":         {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.Frees) }},
	"GCCPUFraction": {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return rtm.GCCPUFraction }},
	"GCSys":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.GCSys) }},
	"HeapAlloc":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapAlloc) }},
	"HeapIdle":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapIdle) }},
	"HeapInuse":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapInuse) }},
	"HeapObjects":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapObjects) }},
	"HeapReleased":  {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapReleased) }},
	"HeapSys":       {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapSys) }},
	"LastGC":        {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.LastGC) }},
	"Lookups":       {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.Lookups) }},
	"MCacheInuse":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheInuse) }},
	"MCacheSys":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheSys) }},
	"MSpanInuse":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanInuse) }},
	"MSpanSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanSys) }},
	"Mallocs":       {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.Mallocs) }},
	"NextGC":        {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.NextGC) }},
	"NumForcedGC":   {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.NumForcedGC) }},
	"NumGC":         {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.NumGC) }},
	"OtherSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.OtherSys) }},
	"PauseTotalNs":  {domain.MetricTypeCounter, func(rtm *runtime.MemStats) interface{} { return int64(rtm.PauseTotalNs) }},
	"StackInuse":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackInuse) }},
	"StackSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackSys) }},
	"Sys":           {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Sys) }},
	"TotalAlloc":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.TotalAlloc) }},
}

func Update(m *domain.Metrics) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Mu.Lock()         // Блокируем доступ к мапе
	defer m.Mu.Unlock() // Освобождаем доступ после обновления

	// Проходим по метрикам и обновляем значения
	for name, def := range metricDefinitions {
		value := def.Update(&rtm)

		switch def.Type {
		case domain.MetricTypeGauge:
		case domain.MetricTypeCounter:
			m.Data[name] = (&domain.Metric{}).SetValue(value).SetName(name).SetType(def.Type)
		}
	}

	if pollCount, ok := m.Data["PollCount"].(domain.MetricInterface); ok {
		if pollCount.GetType() == domain.MetricTypeCounter {
			value := m.Data["PollCount"].GetValue().(int64) + 1
			m.Data["PollCount"].SetValue(value).SetName("PollCount")
		}
	} else {
		m.Data["PollCount"] = (&domain.Metric{}).SetValue(1).SetName("PollCount").SetType(domain.MetricTypeCounter)
	}

	// Обновляем RandomValue
	m.Data["RandomValue"] = (&domain.Metric{}).SetValue(rand.Float64() * 100).SetName("RandomValue").SetType(domain.MetricTypeGauge)
}
