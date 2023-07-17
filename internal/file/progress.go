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

package file

import (
	"io"
)

func NewProgressReader(r io.Reader, fn func(current int64)) io.Reader {
	if r == nil || fn == nil {
		return r
	}
	return &Reader{
		r:  r,
		fn: fn,
	}
}

type Reader struct {
	r    io.Reader
	read int64
	fn   func(current int64)
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == nil && n > 0 {
		r.read += int64(n)
		r.fn(r.read)
	}
	return n, err
}
