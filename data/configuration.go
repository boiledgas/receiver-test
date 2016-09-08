package data

import (
	"errors"
	"fmt"
	"github.com/boiledgas/receiver-test/data/values"
)

type CodeId string

type ConfId interface{}

type Conf struct {
	Id   ConfId // configuration identity
	ETag uint64 // modified hash

	Code       CodeId              // code
	Modules    map[uint16]Module   // modules
	Properties map[uint16]Property // properties
}

type Module struct {
	Id   uint16
	Code CodeId
	Name string
}

func (d *Conf) Module(m Module) (id uint16, err error) {
	if d.Modules == nil {
		d.Modules = make(map[uint16]Module)
	}
	if _, ok := d.GetModule(m.Code); ok {
		err = errors.New("module exists")
	} else {
		id = uint16(len(d.Modules))
		m.Id = id
		d.Modules[id] = m
	}
	return
}

func (d *Conf) Property(name CodeId, p Property) (id uint16, err error) {
	if d.Modules == nil {
		err = errors.New("module not exist")
		return
	}
	if d.Properties == nil {
		d.Properties = make(map[uint16]Property)
	}
	if m, ok := d.GetModule(name); ok {
		if _, ok := d.GetProperty(name, p.Code); ok {
			err = errors.New("property exist")
		} else {
			id = uint16(len(d.Properties))
			p.Id = id
			p.ModuleId = m.Id
			d.Properties[id] = p
		}
	} else {
		err = errors.New("module not exist")
	}
	return
}

func (d *Conf) GetProperty(module CodeId, name CodeId) (result Property, ok bool) {
	ok = false
	if m, mok := d.GetModule(module); mok {
		result, ok = d.getProperty(m.Id, name)
	}
	return
}

func (d *Conf) GetPropertyByType(dataType values.DataType, p *Property) (err error) {
	for _, property := range d.Properties {
		if property.Type == dataType {
			*p = property
			return
		}
	}
	err = errors.New(fmt.Sprintf("property (%v) not found", dataType))
	return
}

func (d *Conf) GetModule(name CodeId) (result Module, ok bool) {
	ok = false
	if d.Modules == nil {
		return
	}
	for _, m := range d.Modules {
		if m.Code == name {
			result, ok = m, true
			break
		}
	}
	return
}

func (d *Conf) getProperty(id uint16, name CodeId) (result Property, ok bool) {
	ok = false
	if d.Properties == nil {
		return
	}
	for _, p := range d.Properties {
		if p.ModuleId == id && p.Code == name {
			result, ok = p, true
			break
		}
	}
	return
}
