package file

import "io"

type ReadFile interface {
	io.Reader
	io.Closer
	Size() int64
	StartSeek(whence int) error
}
