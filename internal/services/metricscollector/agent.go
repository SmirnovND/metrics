package metricscollector

import (
	"bytes"
	"fmt"
	"github.com/SmirnovND/metrics/domain"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
)

type Metrics struct {
	data map[string]domain.Metric
	mu   sync.Mutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		data: make(map[string]domain.Metric),
	}
}

var metricDefinitions = map[string]struct {
	Type   string
	Update func(rtm *runtime.MemStats) interface{}
}{
	"Alloc":         {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.Alloc) }},
	"BuckHashSys":   {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.BuckHashSys) }},
	"Frees":         {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.Frees) }},
	"GCCPUFraction": {"gauge", func(rtm *runtime.MemStats) interface{} { return rtm.GCCPUFraction }},
	"GCSys":         {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.GCSys) }},
	"HeapAlloc":     {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapAlloc) }},
	"HeapIdle":      {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapIdle) }},
	"HeapInuse":     {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapInuse) }},
	"HeapObjects":   {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapObjects) }},
	"HeapReleased":  {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapReleased) }},
	"HeapSys":       {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.HeapSys) }},
	"LastGC":        {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.LastGC) }},
	"Lookups":       {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.Lookups) }},
	"MCacheInuse":   {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheInuse) }},
	"MCacheSys":     {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.MCacheSys) }},
	"MSpanInuse":    {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanInuse) }},
	"MSpanSys":      {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.MSpanSys) }},
	"Mallocs":       {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.Mallocs) }},
	"NextGC":        {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.NextGC) }},
	"NumForcedGC":   {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.NumForcedGC) }},
	"NumGC":         {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.NumGC) }},
	"OtherSys":      {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.OtherSys) }},
	"PauseTotalNs":  {"counter", func(rtm *runtime.MemStats) interface{} { return int64(rtm.PauseTotalNs) }},
	"StackInuse":    {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackInuse) }},
	"StackSys":      {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.StackSys) }},
	"Sys":           {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.Sys) }},
	"TotalAlloc":    {"gauge", func(rtm *runtime.MemStats) interface{} { return float64(rtm.TotalAlloc) }},
}

func (m *Metrics) Update() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления

	// Проходим по метрикам и обновляем значения
	for name, def := range metricDefinitions {
		value := def.Update(&rtm)

		switch def.Type {
		case "gauge":
			m.data[name] = &domain.Gauge{Value: value.(float64), Name: name}
		case "counter":
			m.data[name] = &domain.Counter{Value: value.(int64), Name: name}
		}
	}

	// Обновляем PollCount
	if pollCount, ok := m.data["PollCount"].(*domain.Counter); ok {
		m.data["PollCount"] = &domain.Counter{Value: pollCount.Value + 1, Name: "PollCount"}
	} else {
		m.data["PollCount"] = &domain.Counter{Value: 1, Name: "PollCount"}
	}

	// Обновляем RandomValue
	m.data["RandomValue"] = &domain.Gauge{Value: rand.Float64() * 100, Name: "RandomValue"}
}

// Метод для отправки метрик
func (m *Metrics) Send(serverHost string) {
	m.mu.Lock()         // Блокируем доступ к мапе
	defer m.mu.Unlock() // Освобождаем доступ после обновления

	for _, metric := range m.data {
		url := fmt.Sprintf("%s/update/%s/%s/%v", serverHost, metric.GetType(), metric.GetName(), metric.GetValue())

		// Создание HTTP-запроса
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
		if err != nil {
			fmt.Println("Error creating request:", err)
			continue
		}

		req.Header.Set("Content-Type", "text/plain")

		// Отправка запроса
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request:", err)
			continue
		}
		defer resp.Body.Close()

		// Обработка ответа
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Metric sent successfully:", metric.GetName())
		} else {
			fmt.Printf("Failed to send metric %s: %s\n", metric.GetName(), resp.Status)
		}
	}
}
