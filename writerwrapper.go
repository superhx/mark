package mark

import (
	"io"
)

//WriterWrapper ...
type WriterWrapper struct {
	w   io.Writer
	n   int64
	err error
}

func (wrapper *WriterWrapper) Write(p []byte) (n int, err error) {
	n, err = wrapper.w.Write(p)
	wrapper.n += int64(n)
	if err != nil {
		wrapper.err = err
	}
	return
}

//NumberOfBytes ...
func (wrapper WriterWrapper) NumberOfBytes() (n int64) {
	return wrapper.n
}

func (wrapper WriterWrapper) Error() (err error) {
	return wrapper.err
}
