package config

type Receiver struct {
	Parser   string
	Protocol string
	Host     string
	Port     int32
	Timeout  int64
	Static   bool
}

type Transmitter struct {
	Parser   string // протокол кодирования данных
	Protocol string // протокол передачи данных
	Server   string // адрес назначения
	Count    byte   // количество потоков передачи
}

type Service struct {
	InstanceId string `yaml:"instance"`
	Receiver   map[string]Receiver
}
