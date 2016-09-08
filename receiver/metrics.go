package receiver

import (
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/metrics"
)

const (
	TAG_NAME     string = "github.com/boiledgas/receiver-test"
	TAG_PROTOCOL string = "protocol"
	TAG_CLIENT   string = "client"
)

type ClientsMetric struct {
	measure *metrics.Measure
	packets *metrics.Counter
	records *metrics.Counter
}

func NewClientMetric(name string, code data.CodeId) (m *ClientsMetric, err error) {
	clientTags := map[string]string{
		TAG_NAME:   name,
		TAG_CLIENT: string(code),
	}
	m = &ClientsMetric{}
	if m.measure, err = metrics.NewMeasure("client", clientTags); err != nil {
		return
	}
	m.packets, _ = m.measure.ResetCounter("packets")
	m.records, _ = m.measure.ResetCounter("records")
	return
}

func (c *ClientsMetric) Packet() {
	c.packets.Inc()
}

func (c *ClientsMetric) Records(count int) {
	c.records.IncBy(int64(count))
}

func (c *ClientsMetric) Release() {
	metrics.Release(c.measure)
}

type ConnectionsMetric struct {
	measure     *metrics.Measure
	connections *metrics.Counter
	listener    *metrics.Counter
	errors      *metrics.Counter
	bytes       *metrics.Counter
}

func NewConnectionsMetric(name string, protocol string) (m *ConnectionsMetric, err error) {
	connectionsTags := map[string]string{
		TAG_NAME:     name,
		TAG_PROTOCOL: protocol,
	}
	m = &ConnectionsMetric{}
	if m.measure, err = metrics.NewMeasure("connection", connectionsTags); err != nil {
		return
	}
	m.connections, _ = m.measure.SimpleCounter("connection")
	m.listener, _ = m.measure.ResetCounter("accept")
	m.bytes, _ = m.measure.ResetCounter("bytes")
	m.errors, _ = m.measure.ResetCounter("error")
	return
}

func (c *ConnectionsMetric) Connect() {
	c.connections.Inc()
}

func (c *ConnectionsMetric) Disconnect() {
	c.connections.Dec()
}

func (c *ConnectionsMetric) Accept() {
	c.listener.Inc()
}

func (c *ConnectionsMetric) Bytes(count int) {
	c.bytes.IncBy(int64(count))
}

func (c *ConnectionsMetric) Error() {
	c.errors.Inc()
}

func (c *ConnectionsMetric) Release() {
	metrics.Release(c.measure)
}
