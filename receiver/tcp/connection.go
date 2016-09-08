package tcp

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Контекст tcp клиента
type Connection struct {
	id         uint32   // идентификатор сессии
	time       int64    // время последнего пакета
	connect    int64    // время подключения
	connection net.Conn // подключение
}

func (c *Connection) String() string {
	return fmt.Sprintf("%v", c.id)
}

func (c *Connection) UpdateTime() {
	atomic.StoreInt64(&c.time, time.Now().Unix())
}

func (c *Connection) Close() error {
	if c.connection == nil {
		return errors.New("connection not exists")
	}
	return c.connection.Close()
}

func (c *Connection) ConnectionTime() int64 {
	return c.connect
}

func (c *Connection) LastPacketDate() int64 {
	return atomic.LoadInt64(&c.time)
}

func (c *Connection) PacketCount() uint32 {
	return 0
}
