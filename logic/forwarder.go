package logic

import (
	"io"
)

type Forwarder struct {
	ReaderSrc io.Reader
	WriterSrc io.Writer
	ReaderDst io.Writer
	WriteDst  io.Writer
}

func (r *Forwarder) Read(buf []byte) (n int, err error) {
	n, err = r.ReaderSrc.Read(buf)
	if err != nil {
		panic(err)
	}
	r.ReaderDst.Write(buf[:n])
	return
}

func (r *Forwarder) Write(buf []byte) (n int, err error) {
	if n, err = r.WriterSrc.Write(buf); err != nil {
		panic(err)
	}
	r.WriteDst.Write(buf[:n])
	return
}
