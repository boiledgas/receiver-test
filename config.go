package receiver

import (
	"github.com/boiledgas/receiver-test/metrics"
	"github.com/boiledgas/receiver-test/receiver"
	"github.com/boiledgas/receiver-test/source"
	"github.com/boiledgas/receiver-test/transmitter"
)

type Config struct {
	InstanceId  string `yaml:"instance"`
	Receiver    map[string]receiver.Config
	Transmitter map[string]transmitter.Config
	Source      map[string]source.Config
	Metrics     metrics.InfluxDbConfig
}
