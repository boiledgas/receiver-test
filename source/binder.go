package source

import "github.com/boiledgas/receiver-test/data"

type Binder interface {
	Bind(data.CodeId) (chan []data.Record, error)
}
