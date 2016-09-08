package transmitter

import (
	"errors"
	"reflect"
)

type ParserFactory func() (Parser, error)

var parserInterface reflect.Type

func init() {
	var parser Parser
	parserInterface = reflect.TypeOf(&parser).Elem()
}

type FactoryMap map[string]reflect.Type

var Factory FactoryMap = make(map[string]reflect.Type)

func (m FactoryMap) Exists(name string) (ok bool) {
	_, ok = m[name]
	return
}

func (m FactoryMap) Register(name string, parserType reflect.Type) (err error) {
	if _, exists := m[name]; exists {
		err = errors.New("name already register read parser")
		return
	}
	if parserType.Implements(parserInterface) {
		typeKind := parserType.Kind()
		switch typeKind {
		case reflect.Ptr:
			parserType = parserType.Elem()
		}
		m[name] = parserType
	} else {
		err = errors.New("type not implement receiver.Parser interface")
	}
	return
}

func (m FactoryMap) Create(name string) (factory ParserFactory, err error) {
	if _, ok := m[name]; !ok {
		err = errors.New("not found")
		return
	}
	factory = func() (p Parser, err error) {
		if t, ok := m[name]; !ok {
			err = errors.New("parser not found")
			return
		} else {
			instance := reflect.New(t).Interface()
			var ok bool
			if p, ok = instance.(Parser); !ok {
				err = errors.New("parser is not of transmitter.Parser type")
			}
		}
		return
	}
	return
}
