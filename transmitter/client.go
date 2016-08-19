package transmitter

import (
	"log"
	"net"
	"receiver/cache"
	"receiver/config"
	"receiver/data"
	"receiver/logic"
	"sync"
)

type Client struct {
	sync.WaitGroup

	Config config.Transmitter
	Parser logic.WriteParser
	Cache  *cache.Configuration

	Conn    net.Conn
	Records chan []*data.Record
}

func (c *Client) Disconnect() (err error) {
	if c.Conn != nil {
		err = c.Conn.Close()
		c.Conn = nil
	}
	return
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
	if err = c.Connect(); err != nil {
		return
	}
	go c.HandleClient()
	return
}

func (c *Client) Stop() (err error) {
	if c.Conn != nil {
		err = c.Conn.Close()
	}
	return
}

func (c *Client) HandleClient() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err) // fatal error
			c.Stop()
			c.Start()
		}
	}()

	for recs := range c.Records {
		if len(recs) == 0 {
			continue
		}

		var configuration data.Configuration
		if err := c.Cache.GetById(recs[0].DeviceId, &configuration); err != nil {
			log.Println(err)
			continue
		}

		if err := c.Parser.Parse(c.Conn, configuration, recs); err != nil {
			log.Println(err)
		}
	}
}
