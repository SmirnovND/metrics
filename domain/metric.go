package domain

type Metric interface {
	GetValue() interface{}
	GetName() string
	GetType() string
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
	return "gauge"
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
	return "counter"
}
