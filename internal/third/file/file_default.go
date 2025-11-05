//go:build !js

package file

import (
	"bufio"
	"io"
	"os"
)

const readBufferSize = 1024 * 1024 * 5 // 5mb

func Open(req *UploadFileReq) (ReadFile, error) {
	file, err := os.Open(req.Filepath)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	df := &defaultFile{
		file: file,
		info: info,
	}
	df.resetReaderBuffer()
	return df, nil
}

type defaultFile struct {
	file   *os.File
	info   os.FileInfo
	reader io.Reader
}

func (d *defaultFile) resetReaderBuffer() {
	d.reader = bufio.NewReaderSize(d.file, readBufferSize)
}

func (d *defaultFile) Read(p []byte) (n int, err error) {
	return d.reader.Read(p)
}

func (d *defaultFile) Close() error {
	return d.file.Close()
}

func (d *defaultFile) StartSeek(whence int) error {
	if _, err := d.file.Seek(io.SeekStart, whence); err != nil {
		return err
	}
	d.resetReaderBuffer()
	return nil
}

func (d *defaultFile) Size() int64 {
	return d.info.Size()
}
