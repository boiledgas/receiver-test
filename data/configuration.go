package data

import (
	"receiver/errors"
)

type CodeId string

type ConfigurationId interface{}

type Configuration struct {
	Id         ConfigurationId     // configuration identity
	ETag       uint64              // modified hash

	Code       CodeId              // code
	Modules    map[uint16]Module   // modules
	Properties map[uint16]Property // properties
}

type Module struct {
	Id   uint16
	Code CodeId
	Name string
}

func (d *Configuration) Module(m Module) (id uint16, err error) {
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

func (d *Configuration) Property(name CodeId, p Property) (id uint16, err error) {
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

func (d *Configuration) GetProperty(module CodeId, name CodeId) (result Property, ok bool) {
	ok = false
	if m, mok := d.GetModule(module); mok {
		result, ok = d.getProperty(m.Id, name)
	}
	return
}

func (d *Configuration) GetModule(name CodeId) (result Module, ok bool) {
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

func (d *Configuration) getProperty(id uint16, name CodeId) (result Property, ok bool) {
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
