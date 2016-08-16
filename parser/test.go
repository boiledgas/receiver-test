package parser

import (
	"encoding/binary"
	"io"
	"log"
	"receiver/logic"
	"receiver/data/values"
	"reflect"
	"time"
)

func init() {
	logic.ParserRegistry["Test"] = reflect.TypeOf(TestParser{})
}

type TestParser struct {
}

func (p *TestParser) Parse(rw io.ReadWriter, device logic.Device) (err error) {
	device.Id("test1")
	gps := device.Value("module", "property", time.Now().Unix()).(values.GpsValue)
	if gps != nil {
		gps.SetLatitude(54.45)
		gps.SetLongitude(47.15)
		gps.SetAltitude(12)
		gps.SetSpeed(32)
		gps.SetCourse(45)
	} else {
		log.Println("property not found")
	}
	var v int32
	binary.Read(rw, binary.BigEndian, &v)
	binary.Read(rw, binary.BigEndian, &v)
	binary.Write(rw, binary.BigEndian, v)
	return
}
