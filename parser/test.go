package parser

import (
	"encoding/binary"
	"github.com/boiledgas/receiver-test/data/values"
	"github.com/boiledgas/receiver-test/receiver"
	"io"
	"reflect"
	"time"
)

func init() {
	receiver.Factory.Register("Test", reflect.TypeOf(&TestParser{}))
}

type TestParser struct {
}

func (p *TestParser) Parse(rw io.ReadWriter, device receiver.Device) (err error) {
	device.Id("test1")
	gps := device.Value("module", "property", time.Now()).(values.Gps)
	gps.Latitude = 54.45
	gps.Longitude = 47.15
	gps.Altitude = 12
	gps.Speed = 32
	gps.Course = 45
	var v int32
	binary.Read(rw, binary.BigEndian, &v)
	binary.Read(rw, binary.BigEndian, &v)
	binary.Write(rw, binary.BigEndian, v)
	return
}
