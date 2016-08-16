package client

import (
	"encoding/binary"
	"fmt"
	"net"
	"receiver/errors"
	"sync"
	"time"
)

type Client struct {
	sync.RWMutex

	Host string
	Port int32

	conn net.Conn
}

func (c *Client) Start() (err error) {
	c.Lock()
	if c.conn != nil {
		return errors.New("client is started")
	}
	c.conn, err = net.Dial("tcp", fmt.Sprintf("%v:%d", c.Host, c.Port))
	if err != nil {
		return err
	}
	c.Unlock()

	go func() {
		var value int32
	loop:
		for i := 1; i < 100; i++ {
			if err := binary.Write(c.conn, binary.BigEndian, value); err != nil {
				break loop
			}
			if err := binary.Write(c.conn, binary.BigEndian, value); err != nil {
				break loop
			}
			if err := binary.Read(c.conn, binary.BigEndian, &value); err != nil {
				break loop
			}
			time.Sleep(time.Second * 10)
		}
	}()
	return
}

func (c *Client) Stop() {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		return
	}

	c.conn.Close()
	c.conn = nil
}
