package event_listener

import (
	"errors"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"reflect"
	"runtime"
	"strings"
	"syscall/js"
)

type Caller interface {
	NewCaller(funcName interface{}, callback CallbackWriter, arguments *[]js.Value) Caller
	// AsyncCallWithCallback has promise object return
	AsyncCallWithCallback() interface{}
	// AsyncCallWithOutCallback has promise object return
	AsyncCallWithOutCallback() interface{}
	// SyncCall has not promise
	SyncCall() (result []interface{})
}

var ErrNotSetCallback = errors.New("not set callback to call")

type ReflectCall struct {
	funcName  interface{}
	callback  CallbackWriter
	arguments []js.Value
}

func (r *ReflectCall) AsyncCallWithCallback() interface{} {
	defer func() {
		if rc := recover(); rc != nil {
			r.ErrHandle(rc)
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
	} else {
		panic(ErrNotSetCallback)
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
	funcName.Call(values)
	return r.callback.HandlerFunc()

}
func (r *ReflectCall) NewCaller(funcName interface{}, callback CallbackWriter, arguments *[]js.Value) Caller {
	r.funcName = funcName
	r.callback = callback
	r.arguments = *arguments
	return r
}

type fn func(this js.Value, args []js.Value) interface{}

func (r *ReflectCall) SyncCall() (result []interface{}) {
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
func (r *ReflectCall) AsyncCallWithOutCallback() interface{} {
	defer func() {
		if rc := recover(); rc != nil {
			r.ErrHandle(rc)
		}
	}()
	var funcName reflect.Value
	var typeFuncName reflect.Type
	var temp int
	if r.funcName == nil {
		return nil
	} else {
		funcName = reflect.ValueOf(r.funcName)
		typeFuncName = reflect.TypeOf(r.funcName)
	}
	var values []reflect.Value
	if r.callback == nil {
		r.callback = NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), nil)
	}
	r.callback.SetOperationID(r.arguments[0].String())
	for i := 0; i < len(r.arguments); i++ {
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
	go func() {
		returnValues := funcName.Call(values)
		if len(returnValues) != 0 {
			var result []interface{}
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
			r.callback.SetData(result).SendMessage()
		} else {
			r.callback.SetErrCode(200).SetErrMsg(errors.New("null string").Error()).SendMessage()
		}
	}()
	return r.callback.HandlerFunc()

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
