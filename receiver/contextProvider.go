package receiver

import (
	"receiver/data"
	"receiver/data/cache"
	"sync"
)

type ContextProvider struct {
	sync.RWMutex

	Contexts map[data.CodeId]*Context // словарь подключений
	Cache    *cache.Configuration     // кэш конфигураций
}

// Регистрация устройства
func (p *ContextProvider) Register(code data.CodeId, context *Context) {
	if !context.static || !context.init {
		if err := p.Cache.GetByCode(code, &context.Config); err != nil {
			panic(err)
		}
		context.init = true
	}

	p.Lock()
	defer p.Unlock()
	if _, ok := p.Contexts[code]; !ok {
		context.metrics, _ = NewClientMetric("nothing", code)
		p.Contexts[code] = context
	}
}

// Отключение устройства
func (p *ContextProvider) Close(code data.CodeId) {
	if p.Contexts == nil {
		return
	}

	p.RLock()
	defer p.RUnlock()
	if context, ok := p.Contexts[code]; ok {
		context.metrics.Release()
		context.Closer.Close()
	}
}

func (p *ContextProvider) UpdateConfiguration(newConfig data.Conf) {
	p.RLock()
	defer p.RUnlock()

	if context, ok := p.Contexts[newConfig.Code]; ok {
		if context.Config.ETag != newConfig.ETag {
			context.Lock()
			context.newConfig = newConfig
			context.updated = true
			context.Unlock()
		}
	}
}
