package wasm_wrapper

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

const COMMONEVENTFUNC = "commonEventFunc"

var ErrArgsLength = errors.New("from javascript args length err")
var ErrFunNameNotSet = errors.New("reflect func not to set")

type SetListener struct {
	*WrapperCommon
}

func NewSetListener(wrapperCommon *WrapperCommon) *SetListener {
	return &SetListener{WrapperCommon: wrapperCommon}
}

func (s *SetListener) setConversationListener() {
	callback := event_listener.NewConversationCallback(s.commonFunc)
	open_im_sdk.SetConversationListener(callback)
}
func (s *SetListener) setAdvancedMsgListener() {
	callback := event_listener.NewAdvancedMsgCallback(s.commonFunc)
	open_im_sdk.SetAdvancedMsgListener(callback)
}

//func (s *SetListener) setAdvancedMsgListener() {
//	callback := event_listener.NewAdvancedMsgCallback(s.commonFunc)
//	open_im_sdk.SetAdvancedMsgListener(callback)
//}

func (s *SetListener) SetAllListener() {
	s.setConversationListener()
	s.setAdvancedMsgListener()
}

type WrapperCommon struct {
	commonFunc *js.Value
}

func NewWrapperCommon() *WrapperCommon {
	return &WrapperCommon{}
}
func (w *WrapperCommon) CommonEventFunc(_ js.Value, args []js.Value) interface{} {
	log.NewDebug("CommonEventFunc", "js com here")

	if len(args) >= 1 {
		w.commonFunc = &args[len(args)-1]
		return js.ValueOf(true)
	} else {
		return js.ValueOf(false)
	}
}

type WrapperInitLogin struct {
	*WrapperCommon
	caller event_listener.Caller
}

func NewWrapperInitLogin(wrapperCommon *WrapperCommon) *WrapperInitLogin {
	return &WrapperInitLogin{WrapperCommon: wrapperCommon, caller: &event_listener.ReflectCall{}}
}
func (w *WrapperInitLogin) InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewConnCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return js.ValueOf(w.caller.NewCaller(open_im_sdk.InitSDK, callback, &args).AsyncCallWithCallback())
}
func (w *WrapperInitLogin) Login(_ js.Value, args []js.Value) interface{} {
	listener := NewSetListener(w.WrapperCommon)
	listener.SetAllListener()
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.NewCaller(open_im_sdk.Login, callback, &args).AsyncCallWithCallback()
	return callback.HandlerFunc()
}
