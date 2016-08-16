package logic

import (
	"io"
	"reflect"
)

var ParserRegistry map[string]reflect.Type = make(map[string]reflect.Type)

type ParserFactory func() Parser

type Parser interface {
	Parse(io.ReadWriter, Device) (err error)
}
