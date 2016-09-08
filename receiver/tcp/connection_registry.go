package tcp

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ConnectionRegistry struct {
	next_id          uint32                 // идентификатор следующего клиента
	connections_lock sync.RWMutex           // блокировка словаря клиентов
	connections      map[uint32]*Connection // словарь клиентов
}

// Получение идентификатора tcp сессии
func (e *ConnectionRegistry) getId() uint32 {
	e.connections_lock.RLock()
	defer e.connections_lock.RUnlock()

	id := atomic.AddUint32(&e.next_id, 1)
	if id == 0 {
		id = atomic.AddUint32(&e.next_id, 1)
	}
	for {
		if _, ok := e.connections[id]; !ok {
			return id
		}
		id = atomic.AddUint32(&e.next_id, 1)
	}
}

func (e *ConnectionRegistry) Disconnect(id uint32) {
	e.connections_lock.Lock()
	delete(e.connections, id)
	e.connections_lock.Unlock()
}

func (e *ConnectionRegistry) Connect(conn net.Conn) (c *Connection) {
	c = &Connection{
		id:         e.getId(),
		connection: conn,
		time:       time.Now().Unix(),
	}
	e.connections_lock.Lock()
	e.connections[c.id] = c
	e.connections_lock.Unlock()
	return
}

func (e *ConnectionRegistry) Get(id uint32) (ctx *Connection, ok bool) {
	e.connections_lock.RLock()
	defer e.connections_lock.RUnlock()

	ctx, ok = e.connections[id]
	return
}

func (e *ConnectionRegistry) Each(fun func(*Connection)) {
	e.connections_lock.RLock()
	defer e.connections_lock.RUnlock()

	for _, c := range e.connections {
		fun(c)
	}
}
