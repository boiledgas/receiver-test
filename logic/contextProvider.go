package logic

import (
	"receiver/cache"
	"receiver/data"
	"sync"
)

type ContextProvider struct {
	sync.RWMutex
	contexts map[data.CodeId]*Context // словарь подключений

	Cache   *cache.Configuration      // кэш конфигураций
	Storage *Storage                  // хранилище данных
}

// Регистрация устройства
func (p *ContextProvider) Register(code data.CodeId, context *Context) {
	if p.contexts == nil {
		p.contexts = make(map[data.CodeId]*Context)
	}

	if !context.static || !context.init {
		if err := p.Cache.GetByCode(code, &context.Config); err != nil {
			panic(err)
		}
		context.init = true
	}

	p.RLock()
	defer p.RUnlock()
	if _, ok := p.contexts[code]; !ok {
		p.RUnlock()
		defer p.RLock()
		p.Lock()
		defer p.Unlock()
		p.contexts[code] = context
		context.Out = p.Storage.GetChan(code)
	}
}

// Отключение устройства
func (p *ContextProvider) Close(code data.CodeId) {
	if p.contexts == nil {
		return
	}

	p.RLock()
	defer p.RUnlock()

	if context, ok := p.contexts[code]; ok {
		context.Closer.Close()
		p.Storage.Free(code)
	}
}

func (p *ContextProvider) UpdateConfiguration(newConfig data.Configuration) {
	p.RLock()
	defer p.RUnlock()

	if context, ok := p.contexts[newConfig.Code]; ok {
		if context.Config.ETag != newConfig.ETag {
			context.Lock()
			context.newConfig = newConfig
			context.updated = true
			context.Unlock()
		}
	}
}
