package repository

import (
	"receiver/data"
	"receiver/errors"
	"sync"
)

type Configuration struct {
	sync.RWMutex
	Data map[data.CodeId]data.Configuration
}

var increment uint64 = 1

func (c *Configuration) Init() {
	c.Data = make(map[data.CodeId]data.Configuration)
}

func (c *Configuration) TestData() {
	c1 := data.Configuration{Code: "test1"}
	c2 := data.Configuration{Code: "test2"}
	c3 := data.Configuration{Code: "test3"}
	c.Update(&c1)
	c.Update(&c2)
	c.Update(&c3)
}

func (c *Configuration) Update(conf *data.Configuration) {
	c.Lock()
	defer c.Unlock()

	if conf.Id == nil {
		conf.Id = increment
		increment++
	}
	conf.ETag = conf.ETag + 1
	c.Data[conf.Code] = *conf
}

func (c *Configuration) GetById(id data.ConfigurationId, conf *data.Configuration) (err error) {
	c.RLock()
	defer c.RUnlock()

	for _, c := range c.Data {
		if c.Id == id {
			*conf = c
			return
		}
	}
	err = errors.New("not found")
	return
}

func (c *Configuration) GetByCode(code data.CodeId, conf *data.Configuration) (err error) {
	c.RLock()
	defer c.RUnlock()

	var ok bool
	if *conf, ok = c.Data[code]; !ok {
		err = errors.New("not exists")
	}
	return
}

func (c *Configuration) GetByCodes(codes []data.CodeId, buf []data.Configuration) (err error) {
	c.RLock()
	defer c.RUnlock()

	if len(codes) > len(buf) {
		return errors.New("buf size not enough")
	}
	for i := 0; i < len(codes); i++ {
		buf[i] = c.Data[codes[i]]
	}
	return
}
