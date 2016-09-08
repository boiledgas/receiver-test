package tcp

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"receiver/receiver"
	"io"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"time"
)

type Receiver struct {
	sync.RWMutex       // запуска/останова приемщика
	sync.WaitGroup     // блокировка до завершения активных горутин
	ConnectionRegistry // реестр контекстов

	Name     string                      // receiver name
	Config   receiver.Config             // конфигурация приемника
	Factory  receiver.ParserFactory      // фабрика парсеров
	Provider *receiver.ContextProvider   // провайдер контекста
	Metrics  *receiver.ConnectionsMetric // метрика для статистики по подключениям

	listener    net.Listener // слушатель порта
	exitTimeout chan bool    // канал завершения проверки таймаутов
}

func (r *Receiver) Start() (err error) {
	r.Lock()
	defer r.Unlock()

	host := fmt.Sprintf("%v:%d", r.Config.Host, r.Config.Port)
	log.Printf("receiver: starting %v", host)
	defer log.Printf("receiver: started %v", host)
	if r.listener != nil {
		err = errors.New("listener is active")
		return
	}

	if r.listener, err = net.Listen("tcp", host); err != nil {
		return
	}

	r.connections = make(map[uint32]*Connection)
	if r.Config.Timeout > 0 {
		r.exitTimeout = make(chan bool)
		r.Add(1)
		go r.checkTimeout()
	}

	if r.Config.Listeners == 0 {
		r.Config.Listeners = 1
	}
	for i := 0; i < r.Config.Listeners; i++ {
		r.Add(1)
		go r.listen(i)
	}
	return
}

func (e *Receiver) Stop() (err error) {
	e.Lock()
	defer e.Unlock()

	host := fmt.Sprintf("%v:%d", e.Config.Host, e.Config.Port)
	log.Printf("receiver: stoping %v", host)
	defer log.Printf("stoped %v", host)
	if e.Config.Timeout > 0 {
		close(e.exitTimeout)
	}
	e.listener.Close()
	e.disconnectAll()
	e.Wait()
	e.connections = nil
	e.listener = nil
	return
}

func (r *Receiver) IsActive() bool {
	r.RLock()
	defer r.RUnlock()
	return r.listener != nil
}

func (r *Receiver) Disconnect(id uint32) {
	if c, ok := r.Get(id); ok {
		c.Close()
	}
}

func (r *Receiver) listen(id int) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("receiver: FATAL %v", err)
			debug.PrintStack()
			r.Metrics.Error()
		}
		r.Done()
	}()

	log.Printf("receiver: listener %v started", id)
	defer log.Printf("receiver: listener %v stopped", id)

	for {
		conn, err := r.listener.Accept()
		if err != nil {
			r.Metrics.Error()
			break
		}

		r.Metrics.Accept()
		if conn != nil {
			r.Add(1)
			go r.handleConnection(conn)
		}
	}
}

func (r *Receiver) handleConnection(conn net.Conn) {
	defer func() {
		if state := recover(); state != nil {
			log.Printf("receiver: FATAL %v", state)
			debug.PrintStack()
			r.Metrics.Error()
		}
		r.Done()
	}()

	r.Metrics.Connect()
	defer r.Metrics.Disconnect()
	c := r.Connect(conn)
	defer r.Disconnect(c.id)
	log.Printf("receiver: [%v] connection kept", c)
	defer log.Printf("receiver: [%v] connection dropped", c)

	readBuf := bytes.Buffer{}
	writeBuf := bytes.Buffer{}
	defer func() {
		if err := recover(); err != nil {
			if err == io.EOF {
				log.Printf("receiver: [%v] client close connection", c)
			} else if op, ok := err.(net.OpError); ok {
				if op.Op != "read" {
					log.Printf("receiver: [%v] parse error: %v", c, err)
					log.Printf("receiver: error receive: %v", hex.EncodeToString(readBuf.Bytes()))
					log.Printf("receiver: error response: %v", hex.EncodeToString(writeBuf.Bytes()))
				}
			} else {
				log.Printf("receiver: [%v] parse error: %v", c, err)
				debug.PrintStack()
			}
			r.Metrics.Error()
		}
	}()

	reader := receiver.ReadWriter{
		ReaderSrc: conn,
		WriterSrc: conn,
		ReaderDst: &readBuf,
		WriteDst:  &writeBuf,
	}
	parser, _ := r.Factory()
	context := receiver.Context{
		Provider: r.Provider,
		Receiver: &r.Config,
		Info:     c,
		Closer:   conn,
	}
	for {
		if err := parser.Parse(&reader, &context); err != nil {
			panic(err)
		}
		log.Printf("receiver: receive: %v", hex.EncodeToString(readBuf.Bytes()))
		log.Printf("receiver: response: %v", hex.EncodeToString(writeBuf.Bytes()))
		r.Metrics.Bytes(readBuf.Len())
		readBuf.Reset()
		writeBuf.Reset()
		c.UpdateTime()
		context.Flush()
	}
}

func (e *Receiver) disconnectAll() {
	defer log.Printf("receiver: all clients disconnect")

	close := func(c *Connection) {
		c.Close()
	}
	e.Each(close)
}

func (e *Receiver) disconnectByTimeout() {
	t := time.Now().Unix()
	checkDisconnect := func(c *Connection) {
		if t-c.LastPacketDate() > e.Config.Timeout {
			log.Printf("receiver: timeout client: %v", c.id)
			c.Close()
		}
	}
	e.Each(checkDisconnect)
}

func (e *Receiver) checkTimeout() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("receiver: FATAL %v", err)
			debug.PrintStack()
		}
		e.Done()
	}()

	log.Println("receiver: start check for shutdown")
	defer log.Println("receiver: stop check for shutdown")
loop:
	for {
		select {
		case <-e.exitTimeout:
			break loop
		case <-time.After(time.Second * 1):
			e.disconnectByTimeout()
		}
	}
}
