package tcp

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Контекст tcp клиента
type Context struct {
	id         uint32   // идентификатор сессии
	time       int64    // время последнего пакета
	connect    int64    // время подключения
	connection net.Conn // подключение
}

func (c *Context) String() string {
	return fmt.Sprintf("%v", c.id)
}

func (c *Context) UpdateTime() {
	atomic.StoreInt64(&c.time, time.Now().Unix())
}

func (c *Context) Close() error {
	return c.connection.Close()
}

func (c *Context) ConnectionTime() int64 {
	return c.connect
}

func (c *Context) LastPacketDate() int64 {
	return atomic.LoadInt64(&c.time)
}

func (c *Context) PacketCount() uint32 {
	return 0
}
