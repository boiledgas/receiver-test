package transmitter

import "github.com/boiledgas/receiver-test/source"

type Source interface {
	SetBinder(source.Binder) (err error)
}
