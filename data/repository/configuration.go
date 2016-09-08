package repository

import (
	"errors"
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/data/values"
	"sync"
)

type Configuration struct {
	sync.RWMutex
	Data map[data.CodeId]data.Conf
}

var increment uint64 = 1

func (c *Configuration) Init() {
	c.Data = make(map[data.CodeId]data.Conf)
}

func (c *Configuration) TestData() {
	c1 := data.Conf{
		Code:       "test1",
		Modules:    map[uint16]data.Module{1: data.Module{Id: 1, Code: "testcode"}},
		Properties: map[uint16]data.Property{1: data.Property{Id: 1, Code: "testProperty", Type: values.DATATYPE_GPS, ModuleId: 1}},
	}
	c2 := data.Conf{Code: "test2"}
	c3 := data.Conf{Code: "test3"}
	c.Update(&c1)
	c.Update(&c2)
	c.Update(&c3)
}

func (c *Configuration) Update(conf *data.Conf) {
	c.Lock()
	defer c.Unlock()

	if conf.Id == nil {
		conf.Id = increment
		increment++
	}
	conf.ETag = conf.ETag + 1
	c.Data[conf.Code] = *conf
}

func (c *Configuration) GetById(id data.ConfId, conf *data.Conf) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = state.(error)
		}
	}()
	c.RLock()
	defer c.RUnlock()

	for _, c := range c.Data {
		if c.Id == uint64(id.(int)) {
			*conf = c
			return
		}
	}
	err = errors.New("not found")
	return
}

func (c *Configuration) GetByCode(code data.CodeId, conf *data.Conf) (err error) {
	c.RLock()
	defer c.RUnlock()

	var ok bool
	if *conf, ok = c.Data[code]; !ok {
		err = errors.New("not exists")
	}
	return
}

func (c *Configuration) GetByCodes(codes []data.CodeId, buf []data.Conf) (err error) {
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

func (c *Configuration) GetAll(buf []data.Conf) (count int, err error) {
	c.RLock()
	defer c.RUnlock()

	for _, conf := range c.Data {
		buf[count] = conf
		count++
	}
	return
}
