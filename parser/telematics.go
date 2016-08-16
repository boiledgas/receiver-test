package parser

import (
	"github.com/boiledgas/protocol/telematics"
	"io"
	"receiver/logic"
	"reflect"
	"time"
)

func init() {
	logic.ParserRegistry["Telematics"] = reflect.TypeOf(Telematics{})
}

type Telematics struct {
	Confirm map[byte]time.Time
}

func (p *Telematics) Parse(rw io.ReadWriter, device logic.Device) (err error) {
	device.Id("code")
	reader := telematics.NewReader(rw)
	//writer := telematics.NewWriter(rw)

	var pt byte
	if err = reader.ReadByte(&pt); err != nil {
		return
	}
	switch pt {
	case telematics.PACKET_TYPE_REQUEST:
		req := telematics.Request{}
		reader.ReadRequest(&req)
	case telematics.PACKET_TYPE_RESPONSE:
		resp := telematics.Response{}
		reader.ReadResponse(&resp)
	}
	return
}