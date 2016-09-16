package main

import "unsafe"

type Allocator struct {
	pool    *Pool
	gpsSize int
}

func NewAllocator(size int) *Allocator {
	return &Allocator{
		pool:    NewPool(size),
		gpsSize: int(unsafe.Sizeof(GpsData{})),
	}
}

func (a *Allocator) NewGpsData() (result *GpsData, err error) {
	var buf []byte
	if buf, err = a.pool.Allocate(a.gpsSize); err != nil {
		return
	}
	result = (*GpsData)(unsafe.Pointer(&buf[0]))
	return
}

func (a *Allocator) Allocate(size int) ([]byte, error) {
	return a.pool.Allocate(size)
}

func (a *Allocator) Deallocate(value interface{}) error {
	return a.pool.Deallocate(value)
}
