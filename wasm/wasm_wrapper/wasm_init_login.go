package wasm_wrapper

import (
	"errors"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"reflect"
	"runtime"
	"strings"
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

//反射调用函数，实现javascript的参数js.Value到go语言参数转换，包括错误处理
type ReflectCall struct {
	funcName  interface{}
	callback  event_listener.CallbackWriter
	arguments []js.Value
}

func (r *ReflectCall) InitData(funcName interface{}, callback event_listener.CallbackWriter, arguments *[]js.Value) *ReflectCall {
	r.funcName = funcName
	r.callback = callback
	r.arguments = *arguments
	return r
}

type fn func(this js.Value, args []js.Value) interface{}

func (r *ReflectCall) Call() (result []interface{}) {
	defer func() {
		if rc := recover(); rc != nil {
			temp := r.ErrHandle(rc)
			for _, v := range temp {
				result = append(result, v)
			}
		}
	}()
	var funcName reflect.Value
	var typeFuncName reflect.Type
	var hasCallback bool
	var temp int
	if r.funcName == nil {
		return nil
	} else {
		funcName = reflect.ValueOf(r.funcName)
		typeFuncName = reflect.TypeOf(r.funcName)
	}
	var values []reflect.Value
	if r.callback != nil {
		hasCallback = true
		r.callback.SetOperationID(r.arguments[0].String())
		values = append(values, reflect.ValueOf(r.callback))
	}
	for i := 0; i < len(r.arguments); i++ {
		if hasCallback {
			temp++
		} else {
			temp = i
		}
		//log.NewDebug(r.callback.GetOperationID(), "type is ", typeFuncName.In(temp).Kind(), r.arguments[i].IsNaN())
		switch typeFuncName.In(temp).Kind() {
		case reflect.String:
			convertValue := r.arguments[i].String()
			if !strings.HasPrefix(convertValue, "<number: ") {
				values = append(values, reflect.ValueOf(convertValue))
			} else {
				panic("input args type err index:" + utils.IntToString(i))
			}
		case reflect.Int, reflect.Int32:
			log.NewDebug("", "type is ", r.arguments[i].Int())
			values = append(values, reflect.ValueOf(r.arguments[i].Int()))
		default:
			panic("implement me")
		}
	}
	returnValues := funcName.Call(values)
	if len(returnValues) != 0 {
		for _, v := range returnValues {
			switch v.Kind() {
			case reflect.String:
				result = append(result, v.String())
			case reflect.Bool:
				result = append(result, v.Bool())
			default:
				panic("not support type")
			}
		}
		return result

	} else {
		return nil
	}

}
func (r *ReflectCall) ErrHandle(recover interface{}) []string {
	var temp string
	switch x := recover.(type) {
	case string:
		temp = utils.Wrap(errors.New(x), "").Error()
	case error:
		buf := make([]byte, 1<<20)
		runtime.Stack(buf, true)
		temp = string(buf)
	default:
		temp = utils.Wrap(errors.New("unknown panic"), "").Error()
	}
	if r.callback != nil {
		r.callback.SetErrCode(100).SetErrMsg(temp).SendMessage()
		return []string{}
	} else {
		return []string{temp}
	}
}

type WrapperInitLogin struct {
	*WrapperCommon
	caller *ReflectCall
}

func NewWrapperInitLogin(wrapperCommon *WrapperCommon) *WrapperInitLogin {
	return &WrapperInitLogin{WrapperCommon: wrapperCommon, caller: &ReflectCall{}}
}
func (w *WrapperInitLogin) InitSDK(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewConnCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return js.ValueOf(w.caller.InitData(open_im_sdk.InitSDK, callback, &args).Call())
}
func (w *WrapperInitLogin) Login(_ js.Value, args []js.Value) interface{} {
	listener := NewSetListener(w.WrapperCommon)
	listener.SetAllListener()
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	w.caller.InitData(open_im_sdk.Login, callback, &args).Call()
	return callback.HandlerFunc()
}
