package source

import "receiver/data"

type Binder interface {
	Bind(data.CodeId) (chan []data.Record, error)
}
