package file

import (
	"context"
	"io"
	"time"
)

var temp int

func NewReader(ctx context.Context, r io.Reader, totalSize int64, fn func(current, total int64)) io.Reader {
	temp++
	if r == nil || fn == nil {
		return r
	}
	return &Reader{
		done:      ctx.Done(),
		r:         r,
		totalSize: totalSize,
		fn:        fn,
	}
}

type Reader struct {
	done      <-chan struct{}
	r         io.Reader
	totalSize int64
	read      int64
	fn        func(current, total int64)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	defer func() {
		if temp == 2 {
			time.Sleep(time.Second / 100)
		}
	}()
	select {
	case <-r.done:
		return 0, context.Canceled
	default:
		n, err = r.r.Read(p)
		if err == nil && n > 0 {
			r.read += int64(n)
			r.fn(r.read, r.totalSize)
		}
		return n, err
	}
}
