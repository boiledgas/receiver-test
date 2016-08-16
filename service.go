package receiver

import (
	"receiver/logic"
	"receiver/errors"
	"receiver/tcp"
	"reflect"
	"sync"
	"receiver/config"
)

type client struct {
	id   uint   // внутренний идентификатор клиента
	code string // идентификатор клиента
}

type Service struct {
	Config   config.Service
	Provider *logic.ContextProvider

	receivers_lock sync.RWMutex              // блокировка словаря слушателей
	receivers      map[string]logic.Receiver // словарь слушателей
}

func (s *Service) ListenAndServe() {
	s.receivers_lock.Lock()
	defer s.receivers_lock.Unlock()
	s.receivers = make(map[string]logic.Receiver)
	for name, cfg := range s.Config.Receiver {
		if receiver, err := s.initReceiver(name, cfg); err == nil {
			receiver.Start()
			s.receivers[name] = receiver
		} else {
			panic(err)
		}
	}
}

func (s *Service) Start(name string) (err error) {
	s.receivers_lock.Lock()
	defer s.receivers_lock.Unlock()

	var receiver logic.Receiver
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
	if receiver, ok := s.receivers[name]; ok {
		if receiver.IsActive() {
			receiver.Stop()
		}
	} else {
		err = errors.New("not found")
	}
	return
}

func (s *Service) initReceiver(name string, config config.Receiver) (receiver logic.Receiver, err error) {
	var factory logic.ParserFactory
	switch config.Protocol {
	case "tcp":
		if factory, err = createParserFactory(config.Parser); err != nil {
			return
		}
		receiver = &tcp.Receiver{
			Config:   config,
			Factory:  factory,
			Provider: s.Provider,
		}
	case "udp":
		err = errors.New("not found")
	}
	return
}

func createParserFactory(name string) (factory logic.ParserFactory, err error) {
	if _, ok := logic.ParserRegistry[name]; !ok {
		err = errors.New("not found")
		return
	}
	factory = func() logic.Parser {
		instance := reflect.New(logic.ParserRegistry[name]).Interface()
		return instance.(logic.Parser)
	}
	return
}
