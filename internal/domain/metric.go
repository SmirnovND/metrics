package domain

import (
	"encoding/json"
	"sync"
)

const MetricTypeGauge = "gauge"
const MetricTypeCounter = "counter"

func NewMetrics() *Metrics {
	return &Metrics{
		Data: make(map[string]MetricInterface),
	}
}

type MetricInterface interface {
	GetValue() interface{}
	GetName() string
	GetType() string
	SetValue(value interface{}) MetricInterface
	SetName(name string) MetricInterface
	SetType(mtype string) MetricInterface
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) GetValue() interface{} {
	if m.MType == MetricTypeGauge {
		return m.Value
	} else if m.MType == MetricTypeCounter {
		return m.Delta
	}
	return nil
}

func (m *Metric) SetValue(value interface{}) MetricInterface {
	if m.MType == MetricTypeGauge {
		m.Value = value.(*float64)
	} else if m.MType == MetricTypeCounter {
		m.Delta = value.(*int64)
	}
	return m
}

func (m *Metric) SetType(mtype string) MetricInterface {
	m.MType = mtype
	return m
}

func (m *Metric) GetName() string {
	return m.ID
}

func (m *Metric) SetName(name string) MetricInterface {
	m.ID = name
	return m
}

func (m *Metric) GetType() string {
	return m.MType
}

type Metrics struct {
	Data map[string]MetricInterface
	Mu   sync.RWMutex
}

type Gauge struct {
	Value float64
	Name  string
}

func (g *Gauge) GetValue() interface{} {
	return g.Value
}

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) GetType() string {
	return MetricTypeGauge
}

func (g *Gauge) SetType(_ string) MetricInterface {
	return g
}

func (g *Gauge) SetName(name string) MetricInterface {
	g.Name = name
	return g
}

func (g *Gauge) SetValue(value interface{}) MetricInterface {
	g.Value = value.(float64)
	return g
}

type Counter struct {
	Value int64
	Name  string
}

func (c *Counter) GetValue() interface{} {
	return c.Value
}

func (c *Counter) GetName() string {
	return c.Name
}

func (c *Counter) GetType() string {
	return MetricTypeCounter
}

func (c *Counter) SetValue(value interface{}) MetricInterface {
	c.Value = value.(int64)
	return c
}

func (c *Counter) SetType(_ string) MetricInterface {
	return c
}

func (c *Counter) SetName(name string) MetricInterface {
	c.Name = name
	return c
}

func (m *Metric) MarshalJSON() ([]byte, error) {
	// Создаем временную структуру для сериализации
	type Alias Metric
	metric := &struct {
		Delta int64   `json:"delta,omitempty"`
		Value float64 `json:"value,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	// Устанавливаем значения по умолчанию для Delta и Value
	if m.MType == MetricTypeCounter && m.Delta == nil {
		metric.Delta = 0 // Значение по умолчанию для Counter
	} else if m.Delta != nil {
		metric.Delta = *m.Delta
	}

	if m.MType == MetricTypeGauge && m.Value == nil {
		metric.Value = 0.0 // Значение по умолчанию для Gauge
	} else if m.Value != nil {
		metric.Value = *m.Value
	}

	// Выполняем стандартную сериализацию JSON
	return json.Marshal(metric)
}
