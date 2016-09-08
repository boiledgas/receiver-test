package receiver

import (
	"receiver/data"
	"io"
	"sync"
	"time"
)

type Device interface {
	Id(code data.CodeId)
	Value(moduleCode data.CodeId, propertyCode data.CodeId, time time.Time) (value interface{})
}

// Содержит конфигурацию устройства и обеспечивает отправку записей
type Context struct {
	sync.RWMutex
	Provider *ContextProvider // провайдер контекста

	Receiver *Config        // конфигурация приемщика
	Info     ConnectionInfo // информация о подключении
	Closer   io.Closer      // закрытие подключения
	Config   data.Conf      // текущая конфигурация

	Out chan []*data.Record // канал отправки данных пакета

	static    bool      // статичная конфигурация
	init      bool      // конфигурация установлена
	updated   bool      // конфигурация была обновлена
	newConfig data.Conf // обновленная конфигурация

	records map[time.Time]*data.Record // кэш записей для отправки
	metrics *ClientsMetric             // client metrics
}

func (c *Context) Id(code data.CodeId) {
	c.Provider.Register(code, c)
}

func (c *Context) Value(moduleCode data.CodeId, propertyCode data.CodeId, valueTime time.Time) (value interface{}) {
	if !c.init {
		return
	}
	if c.records == nil {
		c.records = make(map[time.Time]*data.Record)
	}

	c.Lock()
	defer c.Unlock()
	if p, ok := c.Config.GetProperty(moduleCode, propertyCode); ok {
		var record *data.Record
		var ok bool
		if record, ok = c.records[valueTime]; !ok {
			record = &data.Record{
				ConfId: c.Config.Id,
				Time:   valueTime,
				Values: make([]data.PropertyValue, 0, len(c.Config.Properties)),
			}
			c.records[valueTime] = record
		}
		record.ConfId = c.Config.Id
		record.Time = valueTime

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
	rCount := len(c.records)
	c.metrics.Packet()
	c.metrics.Records(rCount)
	recs := make([]*data.Record, rCount)
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
