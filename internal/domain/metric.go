package domain

const MetricTypeGauge = "gauge"
const MetricTypeCounter = "counter"

type Metric interface {
	GetValue() interface{}
	GetName() string
	GetType() string
	SetValue(value interface{})
}

type Gauge struct {
	Value float64
	Name  string
}

func (g *Gauge) GetValue() interface{} {
	return g.Value
}

func (g *Gauge) SetValue(value interface{}) {
	g.Value = value.(float64)
}

func (g *Gauge) GetName() string {
	return g.Name
}

func (g *Gauge) GetType() string {
	return MetricTypeGauge
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

func (c *Counter) SetValue(value interface{}) {
	c.Value = value.(int64)
}
