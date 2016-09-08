package metrics

import "time"

var registry Registry = Registry{
	Measures: make(map[string]*Measure),
}

func NewMeasure(measure string, tags map[string]string) (*Measure, error) {
	return registry.NewMeasure(measure, tags)
}

func Release(m *Measure) (err error) {
	return registry.Release(m)
}

func InfluxDb(interval time.Duration, cfg InfluxDbConfig) {
	registry.InfluxDb(interval, cfg)
}
