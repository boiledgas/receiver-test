package cache

import (
	"log"
	"receiver/data"
	"receiver/repository"
	"sync"
	"time"
)

const BATCH_SIZE uint16 = 3

type UpdateFunc func(data.Configuration)

type Configuration struct {
	sync.RWMutex
	UpdateFunc UpdateFunc                               // функция обновления конфигурации при обновлении в кэше
	Repository *repository.Configuration                //(надо вынести в микросервис работы с хранилищем)
	Cache      map[data.CodeIdentity]data.Configuration // кэшированные конфигурации
}

// вызывать функцию обновления для всех устройств
func (p *Configuration) Get(code data.CodeIdentity, configuration *data.Configuration) (err error) {
	p.RLock()
	defer p.RUnlock()

	if conf, ok := p.Cache[code]; !ok {
		if err = p.Repository.FindByCode(code, configuration); err != nil {
			return
		}
		p.RUnlock()
		defer p.RLock()
		p.Lock()
		defer p.Unlock()
		p.Cache[code] = *configuration
	} else {
		*configuration = conf
	}
	return
}

func (c *Configuration) watch() {
	defer time.AfterFunc(time.Second*1, c.watch)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] watch: %v", err)
		}
	}()
	c.ReloadCache()
}

func (c *Configuration) ReloadCache() (count uint16, err error) {
	c.RLock()
	defer c.RUnlock()

	// escapes to heap
	//keys := make([]data.CodeIdentity, BATCH_SIZE)
	//configurations := make([]data.Configuration, BATCH_SIZE)
	//updates := make([]bool, BATCH_SIZE)
	//keys []data.CodeIdentity, configurations []data.Configuration, updates []bool,

	var keys [BATCH_SIZE]data.CodeIdentity
	var configurations [BATCH_SIZE]data.Configuration
	var updates [BATCH_SIZE]bool
	batchFunc := func(currentBatchSize uint16) {
		if err := c.Repository.FindByCodes(keys[0:currentBatchSize], configurations[:]); err != nil {
			heapErr := err
			log.Printf("update Devices: %v", heapErr)
			return
		}
		for i, configuration := range configurations {
			if w, ok := c.Cache[configuration.Code]; ok && w.ETag != configuration.ETag {
				updates[i] = true
			}
		}
		for i := 0; i < len(keys); i++ {
			if code, ok := keys[i], updates[i]; ok {
				if _, ok := c.Cache[code]; ok {
					c.RUnlock()
					c.Lock()
					count++
					c.Cache[code] = configurations[i]
					c.UpdateFunc(c.Cache[code])
					c.Unlock()
					c.RLock()
				}
			}
		}
	}

	var currentBatchSize uint16 = 0
	for code, _ := range c.Cache {
		if currentBatchSize < BATCH_SIZE {
			keys[currentBatchSize] = code
			updates[currentBatchSize] = false
			currentBatchSize++
		} else {
			batchFunc(currentBatchSize)
			currentBatchSize = 0
		}
	}

	if currentBatchSize > 0 {
		batchFunc(currentBatchSize)
	}
	return
}
