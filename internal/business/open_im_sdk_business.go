// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package business

import (
	"open_im_sdk/open_im_sdk_callback"
)

func (w *Business) SetListener(callback open_im_sdk_callback.OnCustomBusinessListener) {
	if callback == nil {
		return
	}
	w.listener = callback
}
