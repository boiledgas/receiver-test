package logic

import (
	"io"
	"receiver/data"
	"reflect"
)

var ParserRegistry map[string]reflect.Type = make(map[string]reflect.Type)

type ParserFactory func() ReadParser

type ReadParser interface {
	Parse(io.ReadWriter, Device) (err error)
}

type WriteParser interface {
	Parse(io.ReadWriter, data.Configuration, []*data.Record) (err error)
}
