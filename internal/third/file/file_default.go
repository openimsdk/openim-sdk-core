// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
