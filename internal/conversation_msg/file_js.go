// Copyright © 2023 OpenIM SDK. All rights reserved.
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

//go:build js && wasm
// +build js,wasm

package conversation_msg

import (
	"open_im_sdk/wasm/exec"
	"syscall/js"
)

type JSFile struct {
}

func NewFile() *JSFile {
	return &JSFile{}
}
func (j *JSFile) Open(uuid string) (int64, error) {
	return WasmOpen(uuid)
}

func (j *JSFile) Read(uuid string, offset int64, length int64) ([]byte, error) {
	return WasmRead(uuid, offset, length)
}

func (j *JSFile) Close(uuid string) error {
	return WasmClose(uuid)
}

func WasmOpen(uuid string) (int64, error) {
	result, err := exec.Exec(uuid)
	if err != nil {
		return 0, err
	}
	if v, ok := result.(float64); ok {
		return int64(v), nil
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
