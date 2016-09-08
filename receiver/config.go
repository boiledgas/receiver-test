package receiver

type Config struct {
	Parser    string // тип парсера
	Protocol  string // протокол передачи данных
	Host      string // хост передачи
	Port      int32  // порт
	Timeout   int64  // максимальный интервал времени между двумя пакетами
	Listeners int    // count of accept goroutines
	Static    bool   // конфигурация является неизменяемой
}
