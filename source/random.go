package source

import (
	"errors"
	"receiver/data"
	"receiver/data/cache"
	"receiver/data/values"
	"log"
	"sync"
	"time"
)

type Random struct {
	sync.WaitGroup
	exit    chan struct{}
	Config  Config
	Binder  Binder
	Cache   *cache.Configuration
	Clients map[data.CodeId]RandomClient
}

func (s *Random) start() (err error) {
	log.Printf("source: starting Random source")
	defer log.Printf("source: started Random source")

	if s.Binder == nil {
		err = errors.New("transmitter not set for source")
		return
	}

	var conf data.Conf
	idsCount := len(s.Config.Ids)
	for i := 0; i < idsCount; i++ {
		configurationId := s.Config.Ids[i]
		if err = s.Cache.GetById(configurationId, &conf); err != nil {
			log.Printf("source: client cannot be started: %v", configurationId)
			continue
		}
		client := RandomClient{
			Group:         &s.WaitGroup,
			Exit:          make(chan struct{}),
			Configuration: conf,
			Binder:        s.Binder,
		}
		if err = client.Start(); err != nil {
			log.Printf("source: client start: %v", err)
			continue
		}
		s.Clients[conf.Code] = client
	}
	return
}

func (s *Random) Stop() (err error) {
	s.exit <- struct{}{}
	s.Wait()
	return
}

func (s *Random) SetBinder(binder Binder) (err error) {
	if s.Binder != nil {
		err = errors.New("only single mode avaible")
		return
	}
	s.Binder = binder
	err = s.start()
	return
}

func GenerateRecord(conf *data.Conf, rec *data.Record) (err error) {
	var property data.Property
	if err = conf.GetPropertyByType(values.DATATYPE_GPS, &property); err != nil {
		return
	}
	gps := (property.Type.GetValue()).(values.Gps)
	gps.Latitude = 57.35
	gps.Longitude = 35.55
	value := data.PropertyValue{
		ModuleId:   property.ModuleId,
		PropertyId: property.Id,
		Value:      gps,
	}
	rec.ConfId = conf.Id
	rec.Time = time.Now()
	rec.Values = append(rec.Values, value)
	return
}
