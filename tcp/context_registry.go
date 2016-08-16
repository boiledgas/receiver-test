package tcp

import (
	"sync"
	"sync/atomic"
)

type ContextRegistry struct {
	next_id       uint32              // идентификатор следующего клиента
	contexts_lock sync.RWMutex        // блокировка словаря клиентов
	contexts      map[uint32]*Context // словарь клиентов
}

// Получение идентификатора tcp сессии
func (e *ContextRegistry) getId() uint32 {
	e.contexts_lock.RLock()
	defer e.contexts_lock.RUnlock()

	id := atomic.AddUint32(&e.next_id, 1)
	if id == 0 {
		id = atomic.AddUint32(&e.next_id, 1)
	}
	for {
		if _, ok := e.contexts[id]; !ok {
			return id
		}
		id = atomic.AddUint32(&e.next_id, 1)
	}
}

func (e *ContextRegistry) Remove(id uint32) {
	e.contexts_lock.Lock()
	delete(e.contexts, id)
	e.contexts_lock.Unlock()
}

func (e *ContextRegistry) Update(c *Context) {
	e.contexts_lock.Lock()
	e.contexts[c.id] = c
	e.contexts_lock.Unlock()
}

func (e *ContextRegistry) Get(id uint32) (ctx *Context, ok bool) {
	e.contexts_lock.RLock()
	defer e.contexts_lock.RUnlock()

	ctx, ok = e.contexts[id]
	return
}

func (e *ContextRegistry) Each(fun func(*Context)) {
	e.contexts_lock.RLock()
	defer e.contexts_lock.RUnlock()

	for _, c := range e.contexts {
		fun(c)
	}
}