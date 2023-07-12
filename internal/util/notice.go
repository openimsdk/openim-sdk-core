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

package util

import (
	"context"
	"encoding/json"
	"open_im_sdk/pkg/syncer"
)

func NoticeChange[T any](fn func(data string)) func(ctx context.Context, state int, value T) error {
	return func(ctx context.Context, state int, value T) error {
		if state != syncer.Unchanged {
			data, err := json.Marshal(value)
			if err != nil {
				return err
			}
			fn(string(data))
		}
		return nil
	}
}
