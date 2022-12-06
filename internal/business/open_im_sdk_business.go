package business

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
)

func (w *Business) SetListener(callback open_im_sdk_callback.OnCustomBusinessListener) {
	if callback == nil {
		log.NewError("", "callback is null")
		return
	}
	log.NewDebug("", "callback set success")
	w.listener = callback
}
