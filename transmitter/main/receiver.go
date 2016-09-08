package main

import (
	"fmt"
	"log"
	"net"
	"github.com/boiledgas/protocol/telematics"
	"github.com/boiledgas/protocol/telematics/section"
	"sync/atomic"
)

type Receiver struct {
	Port     int
	Listener net.Listener
	ClientId int32
}

func (r *Receiver) Start() (err error) {
	if r.Listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%v", r.Port)); err != nil {
		return
	}
	for i := int32(0); i < 4; i++ {
		go r.Listen(i)
	}
	return
}

func (r *Receiver) Listen(id int32) {
	for {
		if conn, err := r.Listener.Accept(); err != nil {
			log.Printf("->listen (%v): %v", id, err)
		} else {
			log.Printf("-> accept (%v): client: %v", id, conn.RemoteAddr().String())
			clientId := atomic.AddInt32(&r.ClientId, 1)
			go r.HandleClient(id, clientId, conn)
		}
	}
}

func (r *Receiver) HandleClient(id int32, clientId int32, conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("-> receiver: %v", err)
		}
	}()
	reader := telematics.NewReader(conn)
	writer := telematics.NewWriter(conn)
	for {
		req := telematics.Request{}
		if err := reader.ReadRequest(&req); err != nil {
			log.Printf("-> err: %v", err)
		}
		if req.HasConfiguration() {
			reader.Configuration = &req.Conf
			writer.Configuration = &req.Conf
		}
		if req.Has(section.FLAG_MODULE_PROPERTY_VALUE) {
			for _, v := range req.Values {
				log.Printf("-> value: %v => %v", v.ModuleId, v.Values)
			}
		}
		resp := telematics.Response{
			Sequence: req.Sequence,
			Flags:    telematics.RESPONSE_OK,
		}
		if req.Has(section.SECTION_IDENTIFICATION.Flag()) {
			resp.Flags |= telematics.RESPONSE_DESCRIPTION
		}
		writer.WriteResponse(&resp)
	}
}
