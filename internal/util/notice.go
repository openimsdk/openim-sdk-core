// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

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
