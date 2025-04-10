//go:build js && wasm
// +build js,wasm

package wasm_wrapper

// ------------------------------------group---------------------------
type WrapperSignaling struct {
	*WrapperCommon
}

func NewWrapperSignaling(wrapperCommon *WrapperCommon) *WrapperSignaling {
	return &WrapperSignaling{WrapperCommon: wrapperCommon}
}

//func (w *WrapperSignaling) SignalingInviteInGroup(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingInviteInGroup, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingInvite(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingInvite, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingAccept(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingAccept, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingReject(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingReject, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingCancel(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingCancel, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingHungUp(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingHungUp, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingGetRoomByGroupID(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingGetRoomByGroupID, callback, &args).AsyncCallWithCallback()
//}
//
//func (w *WrapperSignaling) SignalingGetTokenByRoomID(_ js.Value, args []js.Value) interface{} {
//	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
//	return event_listener.NewCaller(open_im_sdk.SignalingGetTokenByRoomID, callback, &args).AsyncCallWithCallback()
//}
