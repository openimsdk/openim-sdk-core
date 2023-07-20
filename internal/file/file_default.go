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
	"io"
	"os"
)

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
