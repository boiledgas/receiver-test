package receiver

import (
	"errors"
	"github.com/boiledgas/receiver-test/data"
	"github.com/boiledgas/receiver-test/data/cache"
	"github.com/boiledgas/receiver-test/metrics"
	"github.com/boiledgas/receiver-test/receiver"
	"github.com/boiledgas/receiver-test/receiver/tcp"
	"github.com/boiledgas/receiver-test/source"
	"github.com/boiledgas/receiver-test/transmitter"
	"log"
	"os"
	"sync"
	"time"
)

type client struct {
	id   uint   // внутренний идентификатор клиента
	code string // идентификатор клиента
}

type Service struct {
	Config   Config
	Provider *receiver.ContextProvider
	Cache    *cache.Configuration

	receivers_lock sync.RWMutex                  // блокировка словаря слушателей
	receivers      map[string]Receiver           // словарь слушателей
	transmitters   map[string]Transmitter        // словарь передатчиков
	sources        map[string]transmitter.Source // словарь источников для передатчиков
}

func (s *Service) ListenAndServe() (err error) {
	log.Printf("service: starting serve")
	defer log.Printf("service: started serve")

	go metrics.InfluxDb(time.Second*1, s.Config.Metrics)
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()
	s.receivers_lock.Lock()
	defer s.receivers_lock.Unlock()

	s.receivers = make(map[string]Receiver)
	var r Receiver
	for name, cfg := range s.Config.Receiver {
		if r, err = s.initReceiver(name, cfg); err == nil {
			if err = r.Start(); err != nil {
				return
			}
			s.receivers[name] = r
		} else {
			return err
		}
	}

	s.sources = make(map[string]transmitter.Source)
	var source transmitter.Source
	for name, cfg := range s.Config.Source {
		if source, err = s.initSource(name, cfg); err == nil {
			s.sources[name] = source
		} else {
			return err
		}
	}

	s.transmitters = make(map[string]Transmitter)
	var t Transmitter
	for name, cfg := range s.Config.Transmitter {
		if t, err = s.initTransmitter(name, cfg); err == nil {
			if err = t.Start(); err != nil {
				return err
			}

			s.transmitters[name] = t
		} else {
			return err
		}
	}

	return
}

func (s *Service) Start(name string) (err error) {
	log.Printf("service: starting %v", name)
	defer log.Printf("service: started %v", name)

	s.receivers_lock.Lock()
	defer s.receivers_lock.Unlock()

	var receiver Receiver
	var ok bool
	if receiver, ok = s.receivers[name]; !ok {
		if cfg, ok := s.Config.Receiver[name]; ok {
			receiver, err = s.initReceiver(name, cfg)
		} else {
			err = errors.New("config not founc")
		}
	}

	if err == nil {
		s.receivers[name] = receiver
		if !receiver.IsActive() {
			receiver.Start()
		}
	}
	return
}

func (s *Service) Stop(name string) (err error) {
	log.Printf("service: stoping %v", name)
	defer log.Printf("service: stoped %v", name)
	if receiver, ok := s.receivers[name]; ok {
		if receiver.IsActive() {
			receiver.Stop()
		}
	} else {
		err = errors.New("not found")
	}
	return
}

func (s *Service) initReceiver(name string, config receiver.Config) (result Receiver, err error) {
	var factory receiver.ParserFactory
	if factory, err = receiver.Factory.Create(config.Parser); err != nil {
		return
	}
	var metrics *receiver.ConnectionsMetric
	if metrics, err = receiver.NewConnectionsMetric(name, config.Parser); err != nil {
		return
	}
	switch config.Protocol {
	case "tcp":
		result = &tcp.Receiver{
			Config:   config,
			Factory:  factory,
			Provider: s.Provider,
			Metrics:  metrics,
		}
	case "udp":
		err = errors.New("not found")
	}
	return
}

func (s *Service) initSource(name string, config source.Config) (result transmitter.Source, err error) {
	switch config.Type {
	case "random":
		result = &source.Random{
			Cache:   s.Cache,
			Clients: make(map[data.CodeId]source.RandomClient),
			Config:  config,
		}
	default:
		err = errors.New("not founc source type")
	}
	return
}

func (s *Service) initTransmitter(name string, config transmitter.Config) (result Transmitter, err error) {
	var factory transmitter.ParserFactory
	if factory, err = transmitter.Factory.Create(config.Parser); err != nil {
		return
	}
	result = &transmitter.SingleClientTransmitter{
		Config:  config,
		Cache:   s.Cache,
		Clients: make(map[data.CodeId]*transmitter.Client),
		Factory: factory,
	}
	if s, ok := s.sources[config.Source]; !ok {
		err = errors.New("source not found")
		return
	} else {
		result.SetSource(s)
	}
	return
}
