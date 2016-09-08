package transmitter

import (
	"receiver/data"
	"io"
)

type Parser interface {
	Parse(io.ReadWriter, *data.Conf, []data.Record) (err error)
}
