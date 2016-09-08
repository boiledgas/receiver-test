package parser

import (
	"errors"
	"receiver/data"
	"receiver/data/values"
	"receiver/receiver"
	"receiver/transmitter"
	"io"
	"log"
	"github.com/boiledgas/protocol/telematics"
	"github.com/boiledgas/protocol/telematics/section"
	"github.com/boiledgas/protocol/telematics/value"
	"reflect"
	"time"
)

func init() {
	if err := receiver.Factory.Register("Telematics", reflect.TypeOf(&TelematicsReceiver{})); err != nil {
		log.Printf("telematics: receiver init %v", err)
	}
	if err := transmitter.Factory.Register("Telematics", reflect.TypeOf(&TelematicsTransmitter{})); err != nil {
		log.Printf("telematics: transmitter init %v", err)
	}
}

type TelematicsReceiver struct {
	Confirm map[byte]time.Time
	Reader  *telematics.TelematicsReader
	Writer  *telematics.TelematicsWriter
}

func (p *TelematicsReceiver) Parse(rw io.ReadWriter, device receiver.Device) (err error) {
	if p.Reader == nil {
		p.Reader = telematics.NewReader(rw)
	}
	if p.Writer == nil {
		p.Writer = telematics.NewWriter(rw)
	}

	var packet telematics.Packet
	if err = p.Reader.Read(&packet); err != nil {
		return
	}
	if packet.Has(telematics.FLAG_REQUEST) {
		if packet.Request.HasConfiguration() {
			p.Reader.Configuration = &packet.Request.Conf
			p.Writer.Configuration = &packet.Request.Conf
		}
		if packet.Request.Has(section.FLAG_MODULE_PROPERTY_VALUE) {
			for _, v := range packet.Request.Values {
				log.Printf("TelematicsReceiver: value: %v => %v", v.ModuleId, v.Values)
			}
		}
		resp := telematics.Response{
			Sequence: packet.Request.Sequence,
			Flags:    telematics.RESPONSE_OK,
		}
		if packet.Request.Has(section.SECTION_IDENTIFICATION.Flag()) {
			if packet.Request.Id.Has(section.IDENTIFICATION_FLAGS_CODETEXT) {
				device.Id(data.CodeId(packet.Request.Id.CodeText))
			} else {
				panic("no codetext")
			}
			resp.Flags |= telematics.RESPONSE_DESCRIPTION
		}
		if err = p.Writer.WriteResponse(&resp); err != nil {
			return
		}
	}
	if packet.Has(telematics.FLAG_RESPONSE) {
	}

	return
}

type TelematicsTransmitter struct {
	Reader *telematics.TelematicsReader
	Writer *telematics.TelematicsWriter
	Conf   *telematics.Configuration
}

func (p *TelematicsTransmitter) Parse(writer io.ReadWriter, conf *data.Conf, recs []data.Record) (err error) {
	if p.Conf == nil {
		p.Writer = telematics.NewWriter(writer)
		p.Reader = telematics.NewReader(writer)
		p.Conf = &telematics.Configuration{}
		if err = ToTelematicsConf(conf, p.Conf); err != nil {
			return
		}
		p.Writer.Configuration = p.Conf
		p.Reader.Configuration = p.Conf
		req := telematics.Request{
			Sequence:  telematics.Sequence(),
			Timestamp: int32(time.Now().Unix()),
			Id: section.Identification{
				Hash:     p.Conf.Hash,
				CodeText: string(conf.Code),
			},
		}
		req.Set(section.SECTION_IDENTIFICATION.Flag(), true)
		req.Id.Set(section.IDENTIFICATION_FLAGS_CODETEXT, true)
		req.Id.Set(section.IDENTIFICATION_FLAGS_DEVICEHASH, true)
		resp := telematics.Response{}

		log.Printf("TelematicsTransmitter: %v write identification", conf.Code)
		if err = p.Writer.WriteRequest(&req); err != nil {
			return
		}
		if err = p.Reader.ReadResponse(&resp); err != nil {
			return
		}
		log.Printf("TelematicsTransmitter: resp %v", resp)
		if resp.Flags&telematics.RESPONSE_DESCRIPTION > 0 {
			log.Printf("TelematicsTransmitter: %v write configuration", conf.Code)
			req = telematics.Request{
				Sequence:  telematics.Sequence(),
				Timestamp: int32(time.Now().Unix()),
				Conf:      *p.Conf,
			}
			req.Set(section.FLAG_MODULE, true)
			req.Set(section.FLAG_MODULE_PROPERTY, true)
			req.Set(section.FLAG_COMMAND, true)
			req.Set(section.FLAG_COMMAND_ARGUMENT, true)
			resp := telematics.Response{}
			if err = p.Writer.WriteRequest(&req); err != nil {
				return
			}
			if err = p.Reader.ReadResponse(&resp); err != nil {
				return
			}
		}
	}
	for _, rec := range recs {
		valuesCount := len(rec.Values)
		req := telematics.Request{
			Sequence:  telematics.Sequence(),
			Timestamp: int32(rec.Time.Unix()),
			Values:    make([]section.ModulePropertyValue, valuesCount),
		}

		modulesValues := make(map[uint16]section.ModulePropertyValue)
		var moduleValues section.ModulePropertyValue
		var ok bool
		for _, propertyValue := range rec.Values {
			if moduleValues, ok = modulesValues[propertyValue.ModuleId]; !ok {
				moduleValues = section.ModulePropertyValue{
					ModuleId: byte(propertyValue.ModuleId),
					Values:   make(map[byte]interface{}),
				}
			}
			if moduleValues.Values[byte(propertyValue.PropertyId)], err = ToTelematicsValue(propertyValue.Value); err != nil {
				return
			}
			modulesValues[propertyValue.ModuleId] = moduleValues
		}
		for _, v := range modulesValues {
			req.Values = append(req.Values, v)
		}
		req.Set(section.SECTION_MODULE_PROPERTY_VALUE.Flag(), true)
		resp := telematics.Response{}
		log.Printf("TelematicsTransmitter: client %v write record %v", conf.Code, rec.Values)
		if err = p.Writer.WriteRequest(&req); err != nil {
			return
		}
		if err = p.Reader.ReadResponse(&resp); err != nil {
			return
		}
		if resp.Sequence != req.Sequence {
			err = errors.New("Sequence doesnt equals")
		}

	}
	return
}

func ToTelematicsConf(conf *data.Conf, tConf *telematics.Configuration) (err error) {
	tConf.Hash = byte(13)
	for _, m := range conf.Modules {
		tModule := section.Module{
			Id:          byte(m.Id),
			Name:        string(m.Code),
			Description: m.Name,
		}
		tModule.Set(section.MODULE_FLAGS_NAME|section.MODULE_FLAGS_DESCRIPTION, true)
		tConf.Modules = append(tConf.Modules, tModule)
	}
	for _, p := range conf.Properties {
		tProperty := section.ModuleProperty{
			ModuleId: byte(p.ModuleId),
			Id:       byte(p.Id),
			Name:     string(p.Code),
			Type:     ToTelematicsType(p.Type),
		}
		tProperty.Set(section.MODULE_PROPERTY_FLAGS_NAME|section.MODULE_PROPERTY_FLAGS_DESCRIPTION, true)
		tConf.Properties = append(tConf.Properties, tProperty)
	}
	// TODO: commands for test
	return
}

func ToTelematicsType(dataType values.DataType) value.DataType {
	switch dataType {
	case values.DATATYPE_GPS:
		return value.GPS
	}

	return value.NotSet
}

func ToTelematicsValue(data interface{}) (result interface{}, err error) {
	switch v := data.(type) {
	case values.Gps:
		gps := value.Gps{
			Latitude:  float64(v.Latitude),
			Longitude: float64(v.Longitude),
			Altitude:  int16(v.Altitude),
			Speed:     byte(v.Speed),
			Course:    byte(v.Course),
			Sat:       v.Sat,
		}
		gps.Set(value.GPS_FLAG_LATLNG|value.GPS_FLAG_ALTITUDE|value.GPS_FLAG_SPEED|value.GPS_FLAG_COURSE|value.GPS_FLAG_SATELLITES, true)
		result = gps
	default:
		err = errors.New("no default convert to telematicsValue")
	}
	return
}
