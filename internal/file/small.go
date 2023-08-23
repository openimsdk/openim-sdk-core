package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

const smallFile = 1024 * 1024 * 5 // 5MB

func NewSmallBuffer(file ReadFile) (ReadFile, error) {
	if file == nil {
		return nil, errors.New("file is nil")
	}
	if buf, ok := file.(*smallBuffer); ok {
		return buf, nil
	}
	size := file.Size()
	if size > smallFile {
		return file, nil
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if len(data) != int(size) {
		return nil, fmt.Errorf("read file error, size: %d, read: %d", size, len(data))
	}
	return &smallBuffer{
		buf: bytes.NewReader(data),
	}, nil
}

type smallBuffer struct {
	buf *bytes.Reader
}

func (s *smallBuffer) Read(p []byte) (n int, err error) {
	if s.buf == nil {
		return 0, io.EOF
	}
	return s.buf.Read(p)
}

func (s *smallBuffer) Close() error {
	s.buf = nil
	return nil
}

func (s *smallBuffer) Size() int64 {
	return s.buf.Size()
}

func (s *smallBuffer) StartSeek(whence int) error {
	if s.buf == nil {
		return io.EOF
	}
	_, err := s.buf.Seek(io.SeekStart, whence)
	return err
}
