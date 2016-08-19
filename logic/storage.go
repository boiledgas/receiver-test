package logic

import (
	"receiver/data"
	"sync/atomic"
)

type Buffer struct {
	Locked int32
	Delete int32
	Out    chan []*data.Record
}

type Storage struct {
	BufferSize int32
	Buffers    map[data.CodeId]*Buffer // словарь буферов
}

func (s *Storage) GetChan(code data.CodeId) (ch chan []*data.Record) {
	if _, ok := s.Buffers[code]; !ok {
		s.Buffers[code] = &Buffer{Out: make(chan []*data.Record)}
	}
	buf := s.Buffers[code]
	atomic.SwapInt32(&buf.Delete, 0)
	ch = buf.Out
	return
}

func (s *Storage) Free(code data.CodeId) {
	if buf, ok := s.Buffers[code]; ok {
		atomic.SwapInt32(&buf.Locked, 0)
	}
}

func (s *Storage) Drain() (ch chan []*data.Record, code data.CodeId, ok bool) {
	var buf *Buffer
	for code, buf = range s.Buffers {
		if atomic.CompareAndSwapInt32(&buf.Locked, 0, 1) {
			ok = true
			ch = buf.Out
			return
		}
	}
	return
}
