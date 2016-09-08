package cache

import (
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/data/repository"
	"log"
	"sync"
	"time"
)

const BATCH_SIZE uint16 = 3

type UpdateFunc func(data.Conf)

type Configuration struct {
	sync.RWMutex                             // mutex for cache
	UpdateFunc   UpdateFunc                  // function invoke than some configurations changed
	Repository   *repository.Configuration   // configuration repository
	Index        map[data.CodeId]data.ConfId // configuration index by CodeIdentity
	Cache        map[data.ConfId]data.Conf   // configuration cache
}

// вызывать функцию обновления для всех устройств
func (p *Configuration) GetByCode(code data.CodeId, configuration *data.Conf) (err error) {
	p.RLock()
	defer p.RUnlock()

	if configurationId, ok := p.Index[code]; !ok {
		if err = p.Repository.GetByCode(code, configuration); err != nil {
			return
		}
		p.RUnlock()
		defer p.RLock()
		p.put(configuration)
	} else {
		*configuration = p.Cache[configurationId]
	}
	return
}

func (c *Configuration) GetById(configurationId data.ConfId, configuration *data.Conf) (err error) {
	c.RLock()
	defer c.RUnlock()

	var ok bool
	if *configuration, ok = c.Cache[configurationId]; !ok {
		if err = c.Repository.GetById(configurationId, configuration); err != nil {
			return
		}
		c.RUnlock()
		defer c.RLock()
		c.put(configuration)
	}
	return
}

func (c *Configuration) put(configuration *data.Conf) {
	c.Lock()
	defer c.Unlock()

	existId := c.Index[configuration.Code]
	if existId != configuration.Id {
		c.Index[configuration.Code] = configuration.Id
		delete(c.Cache, existId)
	}
	c.Cache[configuration.Id] = *configuration
}

func (c *Configuration) watch() {
	var count uint16
	defer func() {
		if count == 0 {
			time.AfterFunc(time.Second*1, c.watch)
		} else {
			go c.watch()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] watch: %v", err)
		}
	}()
	var err error
	if count, err = c.ReloadCache(); err != nil {
		log.Printf("[ERROR] ReloadCache: %v", err)
	}
}

func (c *Configuration) ReloadCache() (count uint16, err error) {
	c.RLock()
	defer c.RUnlock()

	// escapes to heap
	//keys := make([]data.CodeIdentity, BATCH_SIZE)
	//configurations := make([]data.Configuration, BATCH_SIZE)
	//updates := make([]bool, BATCH_SIZE)
	//keys []data.CodeIdentity, configurations []data.Configuration, updates []bool,

	var keys [BATCH_SIZE]data.CodeId
	var configurations [BATCH_SIZE]data.Conf
	var updates [BATCH_SIZE]bool
	batchFunc := func(currentBatchSize uint16) {
		if err = c.Repository.GetByCodes(keys[0:currentBatchSize], configurations[:]); err != nil {
			log.Printf("update Devices: %v", err)
			return
		}
		for i, configuration := range configurations {
			if w, ok := c.Cache[configuration.Id]; ok && w.ETag != configuration.ETag {
				updates[i] = true
			}
		}
		for i := 0; i < len(keys); i++ {
			if code, ok := keys[i], updates[i]; ok {
				if _, ok := c.Index[code]; ok {
					c.RUnlock()
					count++
					c.put(&configurations[i])
					c.UpdateFunc(configurations[i])
					c.RLock()
				}
			}
		}
	}

	var currentBatchSize uint16 = 0
	for code, _ := range c.Index {
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
