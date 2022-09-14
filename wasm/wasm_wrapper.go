package wasm

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"syscall/js"
)

const COMMONEVENTFUNC = "commonEventFunc"

var commonFunc js.Value

func CommonEventFunc(_ js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		commonFunc = args[len(args)-1]
		return js.ValueOf(true)
	} else {
		return js.ValueOf(false)
	}
}
func InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := NewInitCallback(commonFunc)
	return js.ValueOf(open_im_sdk.InitSDK(callback, args[0].String(), args[1].String()))
}
func Login(_ js.Value, args []js.Value) interface{} {
	callback := NewBaseCallback(utils.GetSelfFuncName(), args[0].String(), commonFunc)
	if len(args) <3 {
		callback.OnError(100,"args err")
		return nil
	}
	open_im_sdk.Login(callback, args[0].String(), args[1].String(), args[2].String())
	return nil
}
