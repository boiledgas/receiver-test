package receiver

import (
	"receiver/metrics"
	"receiver/receiver"
	"receiver/source"
	"receiver/transmitter"
)

type Config struct {
	InstanceId  string `yaml:"instance"`
	Receiver    map[string]receiver.Config
	Transmitter map[string]transmitter.Config
	Source      map[string]source.Config
	Metrics     metrics.InfluxDbConfig
}
