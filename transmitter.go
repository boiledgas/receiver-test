package receiver

import (
	"receiver/data"
	"receiver/transmitter"
)

type Transmitter interface {
	Start() error
	Stop() error
	Bind(data.CodeId) (chan []data.Record, error)
	SetSource(transmitter.Source) error
}
