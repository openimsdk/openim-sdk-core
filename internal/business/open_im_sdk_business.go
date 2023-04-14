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
