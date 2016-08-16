package logic

import (
	"io"
	"receiver/config"
	"receiver/data"
	"sync"
)

type Device interface {
	Id(code data.CodeIdentity)
	Value(moduleCode data.CodeIdentity, propertyCode data.CodeIdentity, time int64) (value interface{})
}

// Содержит конфигурацию устройства и обеспечивает отправку записей
type Context struct {
	sync.RWMutex
	Provider *ContextProvider // провайдер контекста

	Receiver *config.Receiver   // конфигурация приемщика
	Info     ConnectionInfo     // информация о подключении
	Closer   io.Closer          // закрытие подключения
	Config   data.Configuration // текущая конфигурация

	Out chan []*data.Record // канал отправки данных пакета

	static    bool               // статичная конфигурация
	init      bool               // конфигурация установлена
	updated   bool               // конфигурация была обновлена
	newConfig data.Configuration // обновленная конфигурация

	records map[int64]*data.Record // кэш записей для отправки
}

func (c *Context) Id(code data.CodeIdentity) {
	c.Provider.Register(code, c)
}

func (c *Context) Value(moduleCode data.CodeIdentity, propertyCode data.CodeIdentity, time int64) (value interface{}) {
	if !c.init {
		return
	}
	if c.records == nil {
		c.records = make(map[int64]*data.Record)
	}

	c.Lock()
	defer c.Unlock()
	if p, ok := c.Config.GetProperty(moduleCode, propertyCode); ok {
		var record *data.Record
		var ok bool
		if record, ok = c.records[time]; !ok {
			record = &data.Record{
				DeviceId: c.Config.Id,
				Time:     time,
				Values:   make([]data.PropertyValue, 0, len(c.Config.Properties)),
			}
			c.records[time] = record
		}
		record.DeviceId = c.Config.Id
		record.Time = time

		value = p.Type.GetValue()
		propertyValue := data.PropertyValue{
			ModuleId:   p.ModuleId,
			PropertyId: p.Id,
			Value:      value,
		}
		record.Values = append(record.Values, propertyValue)
	}
	return
}

func (c *Context) Flush() {
	recs := make([]*data.Record, len(c.records))
	i := 0
	for _, rec := range c.records {
		recs[i] = rec
		i++
	}
	c.records = nil
	if c.Out != nil {
		c.Out <- recs
	}
	if c.updated {
		c.newConfig = c.newConfig
		c.updated = false
	}
}
