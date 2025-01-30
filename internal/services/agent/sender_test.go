package agent

import (
	"fmt"
	"github.com/SmirnovND/metrics/internal/domain"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
)

// MockMetric implements the Metric interface for testing purposes.
type MockMetric struct {
	name  string
	value interface{}
	typ   string
}

func (m *MockMetric) GetValue() interface{} {
	return m.value
}

func (m *MockMetric) GetName() string {
	return m.name
}

func (m *MockMetric) GetType() string {
	return m.typ
}

func (m *MockMetric) SetValue(value interface{}) domain.MetricInterface {
	m.value = value
	return m
}

func (m *MockMetric) SetType(_ string) domain.MetricInterface {
	return m
}

func (m *MockMetric) SetName(name string) domain.MetricInterface {
	m.name = name
	return m
}

func ToMetric(m *MockMetric) *domain.Metric {
	met := &domain.Metric{}
	met.SetName(m.name).SetType(m.typ).SetValue(m.value)
	return met
}

// TestSend tests the Send method of Metrics.
func TestSend(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("expected Content-Type to be text/plain, got %s", r.Header.Get("Content-Type"))
		}

		expectedPath := fmt.Sprintf("/update/%s/%s/%v", "gauge", "testMetric", "123.45")
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}))
	defer ts.Close() // Закрываем сервер после теста

	v := 123.45
	// Создаем объект Metrics
	m := &domain.Metrics{
		Data: map[string]*domain.Metric{
			"testMetric": ToMetric(&MockMetric{name: "testMetric", value: v, typ: "gauge"}),
		},
		Mu: sync.RWMutex{},
	}

	// Вызываем метод Send
	Send(m, ts.URL)
}

// MockMetricDefinition представляет собой мок для определения метрики.
type MockMetricDefinition struct {
	Type  string
	Value interface{}
}

// Update возвращает значение метрики.
func (m *MockMetricDefinition) Update(*runtime.MemStats) interface{} {
	return m.Value
}
