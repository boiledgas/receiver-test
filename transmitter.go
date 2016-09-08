package receiver

import (
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/transmitter"
)

type Transmitter interface {
	Start() error
	Stop() error
	Bind(data.CodeId) (chan []data.Record, error)
	SetSource(transmitter.Source) error
}
