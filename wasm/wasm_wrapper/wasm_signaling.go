package wasm_wrapper

import (
"open_im_sdk/wasm/event_listener"
)

//------------------------------------group---------------------------
type WrapperSignaling struct {
	*WrapperCommon
	caller event_listener.Caller
}

func NewWrapperSignaling(wrapperCommon *WrapperCommon) *WrapperSignaling {
	return &WrapperSignaling{WrapperCommon: wrapperCommon, caller: &event_listener.ReflectCall{}}
}

