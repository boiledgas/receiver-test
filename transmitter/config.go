package transmitter

import "errors"

type Config struct {
	Parser   string // протокол кодирования данных
	Protocol string // протокол передачи данных
	Server   string // адрес назначения
	Source   string // источник данных
}

func (t *Config) Validate() (err error) {
	if len(t.Server) == 0 {
		return errors.New("server is not defined")
	}
	if len(t.Protocol) == 0 {
		return errors.New("github.com/boiledgas/protocol is not defined")
	}
	if len(t.Source) == 0 {
		return errors.New("source is not defined")
	}
	return
}
