package agent

import (
	"github.com/SmirnovND/metrics/internal/domain"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
)

type MetricDefinitions map[string]struct {
	Type   string
	Update func(rtm *runtime.MemStats) interface{}
}

var BaseMetric = MetricDefinitions{
	"Alloc":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Alloc) }},
	"BuckHashSys":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.BuckHashSys) }},
	"Frees":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Frees) }},
	"GCCPUFraction": {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return rtm.GCCPUFraction }},
	"GCSys":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.GCSys) }},
	"HeapAlloc":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapAlloc) }},
	"HeapIdle":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapIdle) }},
	"HeapInuse":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapInuse) }},
	"HeapObjects":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapObjects) }},
	"HeapReleased":  {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapReleased) }},
	"HeapSys":       {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapSys) }},
	"LastGC":        {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.LastGC) }},
	"Lookups":       {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Lookups) }},
	"MCacheInuse":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheInuse) }},
	"MCacheSys":     {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheSys) }},
	"MSpanInuse":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanInuse) }},
	"MSpanSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanSys) }},
	"Mallocs":       {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Mallocs) }},
	"NextGC":        {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.NextGC) }},
	"NumForcedGC":   {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.NumForcedGC) }},
	"NumGC":         {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.NumGC) }},
	"OtherSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.OtherSys) }},
	"PauseTotalNs":  {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.PauseTotalNs) }},
	"StackInuse":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackInuse) }},
	"StackSys":      {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackSys) }},
	"Sys":           {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.Sys) }},
	"TotalAlloc":    {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} { return float64(rtm.TotalAlloc) }},
}

var AdvancedMetricsDefinitions = MetricDefinitions{
	"TotalMemory": {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} {
		vmStat, _ := mem.VirtualMemory()
		return float64(vmStat.Total)
	}},
	"FreeMemory": {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} {
		vmStat, _ := mem.VirtualMemory()
		return float64(vmStat.Free)
	}},
	"CPUutilization1": {domain.MetricTypeGauge, func(rtm *runtime.MemStats) interface{} {
		cpuUtilization, _ := cpu.Percent(1*time.Second, false)
		return cpuUtilization
	}},
}

func Update(m *domain.Metrics, metrics MetricDefinitions) {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.Mu.Lock()         // Блокируем доступ к мапе
	defer m.Mu.Unlock() // Освобождаем доступ после обновления

	runtime.GC()
	debug.SetGCPercent(100)

	// Проходим по метрикам и обновляем значения
	for name, def := range metrics {

		value := def.Update(&rtm)

		switch def.Type {
		case domain.MetricTypeGauge:
			floatVal, _ := value.(float64)
			m.Data[name] = (&domain.Metric{}).SetType(def.Type).SetValue(&floatVal).SetName(name)
		case domain.MetricTypeCounter:
			intVal, _ := value.(int64)
			m.Data[name] = (&domain.Metric{}).SetType(def.Type).SetValue(&intVal).SetName(name)
		}
	}

	if pollCount, ok := m.Data["PollCount"]; ok {
		if pollCount.GetType() == domain.MetricTypeCounter {
			if value, ok := pollCount.GetValue().(*int64); ok {
				newValue := *value + 1
				// Создаем новую метрику и заменяем старую
				m.Data["PollCount"] = (&domain.Metric{}).SetType(domain.MetricTypeCounter).SetValue(&newValue).SetName("PollCount")
			}
		}
	} else {
		// Создаем новую метрику, если её нет в m.Data
		initialValue := int64(1)
		m.Data["PollCount"] = (&domain.Metric{}).SetType(domain.MetricTypeCounter).SetValue(&initialValue).SetName("PollCount")
	}

	// Обновляем RandomValue
	randomValue := rand.Float64() * 100
	m.Data["RandomValue"] = (&domain.Metric{}).SetType(domain.MetricTypeGauge).SetValue(&randomValue).SetName("RandomValue")
}
