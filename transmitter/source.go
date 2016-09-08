package transmitter

import "receiver/source"

type Source interface {
	SetBinder(source.Binder) (err error)
}
