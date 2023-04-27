// Copyright © 2023 OpenIM SDK.
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
	"context"
	"io"
)

func NewReader(ctx context.Context, r io.Reader, totalSize int64, fn func(current, total int64)) io.Reader {
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
