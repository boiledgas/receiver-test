package receiver

import (
	"io"
)

type ReadWriter struct {
	ReaderSrc io.Reader
	WriterSrc io.Writer
	ReaderDst io.Writer
	WriteDst  io.Writer
}

func (r *ReadWriter) Read(buf []byte) (n int, err error) {
	if n, err = r.ReaderSrc.Read(buf); err != nil {
		return
	}
	_, err = r.ReaderDst.Write(buf[:n])
	return
}

func (r *ReadWriter) Write(buf []byte) (n int, err error) {
	if n, err = r.WriterSrc.Write(buf); err != nil {
		return
	}
	r.WriteDst.Write(buf[:n])
	return
}
