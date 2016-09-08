package transmitter

import (
	"github.com/boiledgas/receiver-test/data"
	"io"
)

type Parser interface {
	Parse(io.ReadWriter, *data.Conf, []data.Record) (err error)
}
