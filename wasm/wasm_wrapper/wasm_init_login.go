package wasm_wrapper

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"runtime"
	"syscall/js"
)

const COMMONEVENTFUNC = "commonEventFunc"

var commonFunc *js.Value
var ErrArgsLength = errors.New("from javascript args length err")

func checker(callback open_im_sdk_callback.Base, args *[]js.Value, count int) {
	if len(*args) != count {
		callback.OnError(100, ErrArgsLength.Error())
		runtime.Goexit()
	}
}

func CommonEventFunc(_ js.Value, args []js.Value) interface{} {
	if len(args) >= 1 {
		commonFunc = &args[len(args)-1]
		return js.ValueOf(true)
	} else {
		return js.ValueOf(false)
	}
}

func InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewInitCallback(commonFunc)
	return js.ValueOf(open_im_sdk.InitSDK(callback, args[0].String(), args[1].String()))
}
func Login(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), commonFunc)
	checker(callback, &args, 3)
	callback.EventData().SetOperationID(args[0].String())
	open_im_sdk.Login(callback, args[0].String(), args[1].String(), args[2].String())
	return nil
}
