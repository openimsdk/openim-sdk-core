//go:build !js

package file

import (
	"io"
	"os"
)

func Open(path string) (ReadFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	return &defaultFile{
		file: file,
		info: info,
	}, nil
}

type defaultFile struct {
	file *os.File
	info os.FileInfo
}

func (d *defaultFile) Read(p []byte) (n int, err error) {
	return d.file.Read(p)
}

func (d *defaultFile) Close() error {
	return d.file.Close()
}

func (d *defaultFile) StartSeek(whence int) error {
	_, err := d.file.Seek(io.SeekStart, whence)
	return err
}

func (d *defaultFile) Size() int64 {
	return d.info.Size()
}
