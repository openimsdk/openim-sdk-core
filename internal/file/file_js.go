//go:build js && wasm
// +build js,wasm

package file

import (
	"errors"
	"io"
	"open_im_sdk/wasm/exec"
)

func Open(uuid string) (ReadFile, error) {
	file := &jsCallFile{uuid: uuid}
	size, err := file.Open()
	if err != nil {
		return nil, err
	}
	return &jsFile{
		size: size,
		file: file,
	}
}

type jsFile struct {
	size   int64
	file   *jsCallFile
	whence int
}

func (j *jsFile) Read(p []byte) (n int, err error) {
	length := len(p)
	if length == 0 {
		return 0, errors.New("read buffer is empty")
	}
	if j.whence >= j.size {
		return 0, io.EOF
	}
	if j.whence+length > j.size {
		length = int(j.size - j.whence)
	}
	data, err := j.file.Read(j.whence, length)
	if err != nil {
		return 0, err
	}
	if len(data) > len(p) {
		return 0, errors.New("js read data > length")
	}
	j.whence += len(data)
	copy(p, data)
	return len(data), nil
}

func (j *jsFile) Close() error {
	return j.file.Close()
}

func (j *jsFile) Size() int64 {
	return j.size
}

func (j *jsFile) StartSeek(whence int) error {
	if whence < 0 || whence > j.size {
		return errors.New("seek whence is out of range")
	}
	j.whence = whence
	return nil
}

type jsCallFile struct {
	uuid string
}

func (j *jsCallFile) Open() (int64, error) {
	return WasmOpen(j.uuid)
}

func (j *jsCallFile) Read(offset int64, length int64) ([]byte, error) {
	return WasmRead(j.uuid, offset, length)
}

func (j *jsCallFile) Close() error {
	return WasmClose(j.uuid)
}

func WasmOpen(uuid string) (int64, error) {
	result, err := exec.Exec(uuid)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		size := int64(v)
		if size < 0 {
			return 0, errors.New("file size < 0")
		}
		return size, nil
	}
	return 0, exec.ErrType
}
func WasmRead(uuid string, offset int64, length int64) ([]byte, error) {
	result, err := exec.Exec(uuid, offset, length)
	if err != nil {
		return nil, err
	} else {
		if v, ok := result.(js.Value); ok {
			return exec.ExtractArrayBuffer(v), nil
		} else {
			return nil, exec.ErrType
		}
	}
}
func WasmClose(uuid string) error {
	_, err := exec.Exec(uuid)
	return err
}
