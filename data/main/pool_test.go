package main

import "testing"

func Test_Allocate(t *testing.T) {
	allocator := NewAllocator(1000)
	if gps, err := allocator.NewGpsData(); err != nil {
		t.Error(err)
	} else {
		if err := allocator.Deallocate(gps); err != nil {
			t.Error(err)
		}
	}
}