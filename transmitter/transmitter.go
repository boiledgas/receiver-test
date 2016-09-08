package transmitter

import (
	"errors"
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/data/cache"
	"log"
	"sync"
)

var ErrorClientExist error = errors.New("Client exists")
var ErrorWrongState error = errors.New("Wrong state operation")

type SingleClientTransmitter struct {
	sync.RWMutex
	sync.WaitGroup

	Config  Config        // transmitter configuration
	Factory ParserFactory // parser factory
	State   bool          // run or stopped state

	Sources     []Source                // sources of client records
	clientsLock sync.RWMutex            // lock for clients map
	Clients     map[data.CodeId]*Client // clients
	Cache       *cache.Configuration    // configuration cache
}

func (t *SingleClientTransmitter) Start() (err error) {
	log.Printf("transmitter: starting transmitter")
	defer log.Printf("transmitter: started transmitter")
	t.Lock()
	defer t.Unlock()

	if t.State {
		err = ErrorWrongState
		return
	}
	if err = t.Config.Validate(); err != nil {
		return
	}
	t.State = true
	return
}

func (t *SingleClientTransmitter) Stop() (err error) {
	log.Printf("transmitter: stopping transmitter")
	defer log.Printf("transmitter: stoped transmitter")
	t.Lock()
	defer t.Unlock()

	if !t.State {
		err = ErrorWrongState
		return
	}
	for _, client := range t.Clients {
		client.Stop()
	}
	t.State = false
	return
}

func (t *SingleClientTransmitter) SetSource(source Source) (err error) {
	if err = source.SetBinder(t); err != nil {
		return
	}
	t.Sources = append(t.Sources, source)
	return
}

func (t *SingleClientTransmitter) Bind(code data.CodeId) (ch chan []data.Record, err error) {
	t.clientsLock.Lock()
	defer t.clientsLock.Unlock()
	log.Printf("transmitter: bind %v", code)

	if _, ok := t.Clients[code]; ok {
		err = ErrorClientExist
		return
	}

	ch = make(chan []data.Record)
	client := &Client{
		Config:  t.Config,
		Factory: t.Factory,
		Records: ch,
	}
	if err = t.Cache.GetByCode(code, &client.Conf); err != nil {
		return
	}
	if err = client.Start(); err != nil {
		return
	}
	t.Clients[code] = client
	return
}

func (t *SingleClientTransmitter) Unbind(code data.CodeId) (err error) {
	log.Printf("transmitter: unbind %v", code)

	if client, ok := t.Clients[code]; ok {
		client.Stop()
		close(client.Records)
		delete(t.Clients, code)
	} else {
		err = errors.New("not found")
	}
	return
}
