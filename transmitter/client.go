package transmitter

import (
	"errors"
	"receiver/data"
	"log"
	"net"
	"runtime/debug"
	"sync"
)

type Client struct {
	sync.Mutex
	sync.WaitGroup

	Config  Config
	Factory ParserFactory

	state   bool
	Conn    net.Conn
	Conf    data.Conf
	Records chan []data.Record
}

func (c *Client) Connect() (err error) {
	protocol := c.Config.Protocol
	server := c.Config.Server
	if c.Conn, err = net.Dial(protocol, server); err != nil {
		return
	}
	return
}

func (c *Client) Start() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.state {
		err = errors.New("already run")
		return
	}
	go c.handleClient()
	c.state = true
	return
}

func (c *Client) Stop() (err error) {
	c.Lock()
	defer c.Unlock()

	if !c.state {
		err = errors.New("not running")
		return
	}
	err = c.stop()
	c.state = false
	return
}

func (c *Client) Restart() (err error) {
	c.Lock()
	defer c.Unlock()

	if !c.state {
		return
	}
	if err = c.stop(); err != nil {
		return
	}
	go c.handleClient()
	return
}

func (c *Client) stop() (err error) {
	if c.Conn != nil {
		err = c.Conn.Close()
		c.Conn = nil
	}
	return
}

func (c *Client) handleClient() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("transmitter client: [%v] %v", c.Conf.Id, err)
			debug.PrintStack()
			if err = c.Restart(); err != nil {
				log.Printf("transmitter client: restart %v", err)
			}
		}
	}()

	log.Printf("transmitter client: started %v", c.Conf.Code)
	defer log.Printf("transmitter client: stoped %v", c.Conf.Code)

	if err := c.Connect(); err != nil {
		panic(err)
	}

	var parser Parser
	var err error
	if parser, err = c.Factory(); err != nil {
		panic(err)
	}
	for recs := range c.Records {
		if len(recs) == 0 {
			continue
		}

		if err = parser.Parse(c.Conn, &c.Conf, recs); err != nil {
			panic(err)
		}
	}
}
