package transmitter

import (
	"receiver/cache"
	"receiver/config"
	"receiver/data"
	"sync/atomic"
	"receiver/errors"
	"sync"
)

type Transmitter struct {
	sync.WaitGroup
	Config  config.Transmitter     // transmitter configuration
	Cache   *cache.Configuration   // configuration cache
	Clients map[data.CodeId]Client // clients
	Source  Source                 // source of client records
}

func (t *Transmitter) Start() {
	for _, client := range t.Clients {
		client.Start()
	}
}

func (t *Transmitter) Stop() {
	for _, client := range t.Clients {
		client.Stop()
	}

}

func (t *Transmitter) Bind(code data.CodeId, ch chan []data.Record) (err error) {
	var client Client
	var ok bool
	if client, ok = t.Clients[code]; !ok {
		client = Client{Transmitter: t}
		if err = client.Connect(); err != nil {
			return
		}
	}
	atomic.SwapPointer(client.Records, ch)
	go client.Start()
	return
}

func (t *Transmitter) Unbind(code data.CodeId) (err error) {
	if client, ok := t.Clients[code]; ok {
		client.Disconnect()
	} else {
		err = errors.New("not found")
	}
}
