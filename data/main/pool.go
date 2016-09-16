package main

import (
	"errors"
	"reflect"
	"unsafe"
)

type Pool struct {
	buffer           []byte
	startOffset      uintptr
	allocationOffset uintptr
	availableSize    int
}

func NewPool(bufferSize int) *Pool {
	buffer := make([]byte, bufferSize)
	return &Pool{
		buffer:           buffer,
		startOffset:      uintptr(unsafe.Pointer(&buffer[0])),
		availableSize:    len(buffer),
		allocationOffset: 0,
	}
}

const OFFSET_HEADER uintptr = 2

var ErrorNotEnoughSpace error = errors.New("not enough space for allocation")
var ErrorNotBelongToPool error = errors.New("object not belong to pool")
var ErrorNotImplemented error = errors.New("not implemented")

func (p *Pool) Allocate(size int) (result []byte, err error) {
	if p.availableSize < size {
		err = ErrorNotEnoughSpace
		return
	}
	if size > 0x7fff {
		err = ErrorNotImplemented
		return
	}

	p.buffer[p.allocationOffset] = 1
	p.buffer[p.allocationOffset+1] = byte(size)
	offset := uintptr(size)
	result = p.buffer[p.allocationOffset+2 : p.allocationOffset+2+offset]
	p.allocationOffset += offset
	p.availableSize -= size
	return
}

func (p *Pool) Deallocate(ptr interface{}) error {
	offset := int(reflect.ValueOf(ptr).Pointer()-p.startOffset) - 2
	if offset >= len(p.buffer) || offset < 0 {
		return ErrorNotBelongToPool
	}
	len := int(p.buffer[offset+1])
	for i := int(0); i < len; i++ {
		p.buffer[offset+2+i] = 0
	}
	p.buffer[offset] = 0
	p.buffer[offset+1] = 0
	return nil
}
