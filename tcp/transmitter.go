package tcp

import (
	"receiver/cache"
	"receiver/config"
	"receiver/data"
	"receiver/logic"
)

type input struct {
	channel chan *data.Record
	conf    *data.Configuration
}

type Transmitter struct {
	Config   config.Transmitter     // конфигурация передатчика
	Factory  logic.ParserFactory    // фабрика парсеров
	Provider *logic.ContextProvider // провайдер контекста

	channel  chan *data.Record      // канал записей для передачи
	cache    cache.Configuration    // кэш конфигураций устройств
	contexts map[uint32]*Context    // словарь клиентов
}

func (t *Transmitter) handleClient() {
	//conn, err := net.Dial("tcp", t.Config.Server)
	//if err != nil {
	//	log.Printf("dial: %v", err)
	//}
	//context := Context{
	//	connection: conn,
	//}
}

func (t *Transmitter) processClient(context *Context) {
}
