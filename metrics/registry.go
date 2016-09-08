package metrics

import (
	"errors"
	"sync"
	"sync/atomic"
)

type CounterType byte

const (
	RESET_COUNTER  CounterType = 0x01
	SIMPLE_COUNTER CounterType = 0x02
)

var ErrorAlreadyExists error = errors.New("already exists")
var ErrorNotExists error = errors.New("bot exists")

type Counter struct {
	value *int64
}

func (c *Counter) Inc() {
	atomic.AddInt64(c.value, 1)
}

func (c *Counter) IncBy(val int64) {
	atomic.AddInt64(c.value, val)
}

func (c *Counter) Dec() {
	atomic.AddInt64(c.value, -1)
}

type Measure struct {
	Name     string                  // measure name
	Counters map[string]*CounterMeta // counters of measure
	// tags of counter (eg
	// system:=receiver|transmitter,
	// protocol:=tcp|udp,
	// parser:=test|telematics,
	// client_code:=test1|test2|test3, client_ip)
	Tags map[string]string
}

func (m *Measure) counter(name string, counterType CounterType) (c *Counter, err error) {
	if _, ok := m.Counters[name]; ok {
		err = ErrorAlreadyExists
		return
	}

	counterMeta := CounterMeta{
		Type: counterType,
	}
	m.Counters[name] = &counterMeta
	c = &Counter{
		value: &counterMeta.value,
	}
	return
}

func (m *Measure) ResetCounter(name string) (*Counter, error) {
	return m.counter(name, RESET_COUNTER)
}

func (m *Measure) SimpleCounter(name string) (*Counter, error) {
	return m.counter(name, SIMPLE_COUNTER)
}

type CounterMeta struct {
	Type  CounterType // type of counter
	value int64       // value of counter
}

func (m *CounterMeta) Value() int64 {
	switch m.Type {
	case RESET_COUNTER:
		return atomic.SwapInt64(&m.value, 0)
	case SIMPLE_COUNTER:
		return atomic.LoadInt64(&m.value)
	default:
		panic("not defined counter type")
	}
}

type Registry struct {
	sync.RWMutex
	Measures map[string]*Measure
}

func (r *Registry) NewMeasure(measure string, tags map[string]string) (result *Measure, err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Measures[measure]; ok {
		err = ErrorAlreadyExists
		return
	}
	result = &Measure{
		Name:     measure,
		Tags:     tags,
		Counters: make(map[string]*CounterMeta),
	}
	r.Measures[measure] = result
	return
}

func (r *Registry) Release(m *Measure) (err error) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.Measures[m.Name]; !ok {
		err = ErrorNotExists
	}
	delete(r.Measures, m.Name)
	return
}
