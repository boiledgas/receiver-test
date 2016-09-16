package main

import (
	"log"
)

type GpsData struct {
	Lat   float32
	Lon   float32
	Speed float32
	Sat   byte
}

func main() {
	alloc1 := NewAllocator(1024)
	alloc2 := NewAllocator(1024)

	d1, _ := alloc1.NewGpsData()
	d1.Lat, d1.Lon = 2.4, 3.5
	alloc1.Deallocate(d1)

	d2, _ := alloc2.NewGpsData()
	d2.Lat = 2.4
	d2.Lon = 3.5
	d2.Speed = 4.6
	if err := alloc1.Deallocate(d2); err != nil {
		log.Printf("%v", err)
	}
}
