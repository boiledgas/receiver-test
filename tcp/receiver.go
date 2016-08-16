package tcp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"receiver/config"
	"receiver/logic"
	"receiver/errors"
	"sync"
	"time"
)

type Receiver struct {
	ContextRegistry                    // реестр контекстов
	sync.RWMutex                       // запуска/останова приемщика
	sync.WaitGroup                     // блокировка до завершения активных горутин
	listener    net.Listener           // слушатель порта

	Config      config.Receiver        // конфигурация приемника
	Factory     logic.ParserFactory    // фабрика парсеров
	Provider    *logic.ContextProvider // провайдер контекста

	exitTimeout chan bool              // канал завершения проверки таймаутов
}

func (r *Receiver) Start() (err error) {
	r.Lock()
	defer r.Unlock()

	host := fmt.Sprintf("%v:%d", r.Config.Host, r.Config.Port)
	log.Printf("starting: %v", host)
	defer log.Printf("started: %v", host)
	if r.listener != nil {
		err = errors.New("listener is active")
		return
	}

	if r.listener, err = net.Listen("tcp", host); err != nil {
		return
	}

	r.contexts = make(map[uint32]*Context)
	if r.Config.Timeout > 0 {
		r.exitTimeout = make(chan bool)
		r.Add(1)
		go r.checkTimeout()
	}

	r.Add(1)
	go r.listen()
	return
}

func (e *Receiver) Stop() (err error) {
	e.Lock()
	defer e.Unlock()

	host := fmt.Sprintf("%v:%d", e.Config.Host, e.Config.Port)
	log.Printf("stoping %v", host)
	defer log.Printf("stoped %v", host)
	if e.Config.Timeout > 0 {
		close(e.exitTimeout)
	}
	e.listener.Close()
	e.disconnectAll()
	e.Wait()
	e.contexts = nil
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

func (r *Receiver) listen() {
	defer r.Done()

	log.Println("listener started")
	defer log.Println("listener stopped")

	for {
		conn, err := r.listener.Accept()
		if err != nil {
			break
		}

		if conn != nil {
			go r.handleConnection(r.getId(), conn)
		}
	}
}

func (r *Receiver) handleConnection(id uint32, conn net.Conn) {
	c := &Context{
		id:         id,
		connection: conn,
		time:       time.Now().Unix(),
	}
	r.Update(c)
	r.Add(1)
	go r.processClient(c)
}

func (r *Receiver) processClient(c *Context) {
	log.Printf("[%v] connection kept", c)
	defer log.Printf("[%v] connection dropped", c)
	defer r.Done()
	defer r.disconnect(c)

	readBuf := bytes.Buffer{}
	writeBuf := bytes.Buffer{}
	defer func() {
		if err := recover(); err != nil {
			if err == io.EOF {
				log.Printf("[%v] client close connection", c)
			} else if op, ok := err.(net.OpError); ok {
				if op.Op != "read" {
					log.Printf("[%v] parse error: %v", c, err)
					log.Printf("error receive: %v", hex.EncodeToString(readBuf.Bytes()))
					log.Printf("error response: %v", hex.EncodeToString(writeBuf.Bytes()))
				}
			} else {
				log.Printf("[%v] parse error: %v", c, err)
			}

		}
	}()

	reader := logic.Forwarder{
		ReaderSrc: c.connection,
		WriterSrc: c.connection,
		ReaderDst: &readBuf,
		WriteDst:  &writeBuf,
	}
	parser := r.Factory()
	context := logic.Context{
		Provider: r.Provider,
		Receiver: &r.Config,
		Info:     c,
		Closer:   c.connection,
	}
	for {
		parser.Parse(&reader, &context)
		log.Printf("receive: %v", hex.EncodeToString(readBuf.Bytes()))
		log.Printf("response: %v", hex.EncodeToString(writeBuf.Bytes()))
		readBuf.Reset()
		writeBuf.Reset()
		c.UpdateTime()
	}
}

func (e *Receiver) disconnect(c *Context) {
	e.Remove(c.id)
	c.connection.Close()
}

func (e *Receiver) disconnectAll() {
	defer log.Printf("all clients disconnect")

	close := func(c *Context) {
		c.Close()
	}
	e.Each(close)
}

func (e *Receiver) disconnectByTimeout() {
	t := time.Now().Unix()
	checkDisconnect := func(c *Context) {
		if t-c.LastPacketDate() > e.Config.Timeout {
			log.Printf("timeout client: %v", c.id)
			c.Close()
		}
	}
	e.Each(checkDisconnect)
}

func (e *Receiver) checkTimeout() {
	defer e.Done()

	log.Println("start check for shutdown")
	defer log.Println("stop check for shutdown")
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
