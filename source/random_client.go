package source

import (
	"receiver/data"
	"log"
	"runtime"
	"sync"
)

type RandomClient struct {
	Group *sync.WaitGroup
	Exit  chan struct{}

	Channel       chan []data.Record
	Configuration data.Conf
	Binder        Binder
}

func (c *RandomClient) Start() (err error) {
	log.Printf("random client: starting %v", c.Configuration.Id)
	defer log.Printf("random client: started %v", c.Configuration.Id)

	c.Channel, err = c.Binder.Bind(c.Configuration.Code)
	if err != nil {
		return
	}

	c.Group.Add(1)
	go c.Handle()
	return
}

func (c *RandomClient) Stop() (err error) {
	c.Exit <- struct{}{}
	return
}

func (c *RandomClient) Handle() {
	defer c.Group.Done()
	log.Printf("random client: start handle %v", c.Configuration.Id)
	defer log.Printf("random client: stop handle %v", c.Configuration.Id)
loop:
	for {
		select {
		case <-c.Exit:
			break loop
		default:
			if err := c.GenerateRecords(); err != nil {
				log.Printf("random client: (%v): %v", c.Configuration.Id, err)
				break loop
			}
		}
		runtime.Gosched()
	}
}

func (c *RandomClient) GenerateRecords() (err error) {
	recs := make([]data.Record, 5)
	for i := 0; i < len(recs); i++ {
		recs[i] = data.Record{}
		if err = GenerateRecord(&c.Configuration, &recs[i]); err != nil {
			return
		}
	}
	c.Channel <- recs
	return
}
