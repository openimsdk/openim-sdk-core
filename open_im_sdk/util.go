package open_im_sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"reflect"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func isNumeric(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func setNumeric(in interface{}, out interface{}) {
	inValue := reflect.ValueOf(in)
	outValue := reflect.ValueOf(out)
	outElem := outValue.Elem()
	outType := outElem.Type()
	inType := inValue.Type()
	if outType.AssignableTo(inType) {
		outElem.Set(inValue)
		return
	}
	inKind := inValue.Kind()
	outKind := outElem.Kind()
	switch inKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(inValue.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Int()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Int()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Uint()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(inValue.Uint())
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(float64(inValue.Uint()))
		}
	case reflect.Float32, reflect.Float64:
		switch outKind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			outElem.SetInt(int64(inValue.Float()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			outElem.SetUint(uint64(inValue.Float()))
		case reflect.Float32, reflect.Float64:
			outElem.SetFloat(inValue.Float())
		}
	}
}

func call_(operationID string, fn any, args ...any) (res any, err error) {
	//defer func() {
	//	if r := recover(); r != nil {
	//		err = fmt.Errorf("call panic: %+v", r)
	//	}
	//}()
	if operationID == "" {
		return nil, errs.ErrArgs.Wrap("call func operationID is empty")
	}
	//if err := CheckResourceLoad(UserForSDK); err != nil {
	//	return nil, err
	//}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		return nil, fmt.Errorf("call func fn is not func, is %T", fn)
	}
	fnt := fnv.Type()
	nin := fnt.NumIn()
	if len(args)+1 != nin {
		return nil, fmt.Errorf("go code error: fn in args num is not match")
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "operationID", operationID)
	ctx = context.WithValue(ctx, "token", UserForSDK.GetToken())
	ctx = context.WithValue(ctx, "apiHost", UserForSDK.GetConfig().ApiAddr)

	ins := make([]reflect.Value, 0, nin)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		inFnField := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == inFnField.String() || inFnField.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			var ptr int
			for inFnField.Kind() == reflect.Ptr {
				inFnField = inFnField.Elem()
				ptr++
			}
			switch inFnField.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
				v := reflect.New(inFnField)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					return nil, fmt.Errorf("go call json.Unmarshal error: %w", err)
				}
				if ptr == 0 {
					v = v.Elem()
				} else if ptr != 1 {
					for i := ptr - 1; i > 0; i-- {
						temp := reflect.New(v.Type())
						temp.Elem().Set(v)
						v = temp
					}
				}
				ins = append(ins, v)
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(inFnField.Kind()) {
			v := reflect.Zero(inFnField).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		return nil, fmt.Errorf("go code error: fn in args type is not match")
	}
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		return "", nil
	}
	if fnt.Out(len(outs) - 1).Implements(reflect.ValueOf(new(error)).Elem().Type()) {
		if errValueOf := outs[len(outs)-1]; !errValueOf.IsNil() {
			if err := errValueOf.Interface().(error); err != nil {
				return nil, err
			}
		}
		switch len(outs) {
		case 1:
			return "", nil
		case 2:
			return outs[0].Interface(), nil
		}
		outs = outs[:len(outs)-1]
	}
	val := make([]any, 0, len(outs))
	for i := range outs {
		val = append(val, outs[i].Interface())
	}
	return val, nil
}

func call(callback open_im_sdk_callback.Base, operationID string, fn any, args ...any) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	go func() {
		res, err := call_(operationID, fn, args...)
		if err != nil {
			callback.OnError(10000, err.Error())
			return
		}
		data, err := json.Marshal(res)
		if err != nil {
			callback.OnError(10000, err.Error())
			return
		}
		callback.OnSuccess(string(data))
	}()
}

func syncCall(operationID string, fn any, args ...any) string {
	//defer func() {
	//	if r := recover(); r != nil {
	//		//callback.OnError(10001, fmt.Sprintf("%+v", r))
	//		return
	//	}
	//}()
	if operationID == "" {
		//callback.OnError(constant.ErrArgs.ErrCode, errs.ErrArgs.Wrap("operationID is empty").Error())
		return ""
	}
	if err := CheckResourceLoad(UserForSDK); err != nil {
		//callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return ""
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		//callback.OnError(10000, "go code error: fn is not func")
		return ""
	}
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	if len(args)+1 != numIn {
		//callback.OnError(10000, "go code error: fn in args num is not match")
		return ""
	}
	ins := make([]reflect.Value, 0, numIn)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "operationID", operationID)
	ctx = context.WithValue(ctx, "token", UserForSDK.GetToken())
	ctx = context.WithValue(ctx, "apiHost", UserForSDK.GetConfig().ApiAddr)

	ins = append(ins, reflect.ValueOf(ctx))
	for i := 0; i < len(args); i++ {
		tag := fnt.In(i + 1)
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
				v := reflect.New(tag)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					//callback.OnError(constant.ErrArgs.ErrCode, err.Error())
					return ""
				}
				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		//callback.OnError(constant.ErrArgs.ErrCode, "go code error: fn in args type is not match")
		return ""
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.TypeOf((*error)(nil)).Elem())
	}
	fmt.Println("fnv:", fnv.Interface(), "ins:", ins)
	outs := fnv.Call(ins)
	if len(outs) == 0 {
		//callback.OnSuccess("")
		return ""
	}
	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			//callback.OnError(10000, last.(error).Error())
			return ""
		}
		if len(outs) == 1 {
			//callback.OnSuccess("") // 只有一个返回值为error，且error == nil
			return ""
		}
		outVals = outVals[:len(outVals)-1]
	}
	// 将map和slice的nil转换为非nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		//callback.OnError(constant.ErrArgs.ErrCode, err.Error())
		return ""
	}
	return string(jsonData)
}
func messageCall(callback open_im_sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	go messageCall_(callback, operationID, fn, args...)
}
func messageCall_(callback open_im_sdk_callback.SendMsgCallBack, operationID string, fn any, args ...any) {
	defer func() {
		if r := recover(); r != nil {
			callback.OnError(10001, fmt.Sprintf("%+v", r))
			return
		}
	}()
	if operationID == "" {
		callback.OnError(constant.ErrArgs.ErrCode, errs.ErrArgs.Wrap("operationID is empty").Error())
		return
	}
	if err := CheckResourceLoad(UserForSDK); err != nil {
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		callback.OnError(10000, "go code error: fn is not func")
		return
	}
	fnt := fnv.Type()
	numIn := fnt.NumIn()
	if len(args) != numIn+1 {
		callback.OnError(10000, "go code error: fn in args num is not match")
		return
	}
	ins := make([]reflect.Value, 0, numIn)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "operationID", operationID)
	ctx = context.WithValue(ctx, "callback", callback)
	ins = append(ins, reflect.ValueOf(ctx))
	for i := 2; i < len(args); i++ { // callback open_im_sdk_callback.Base, operationID string, ...
		tag := fnt.In(i - 1) // ctx context.Context, ...
		arg := reflect.TypeOf(args[i])
		if arg.String() == tag.String() || tag.Kind() == reflect.Interface {
			ins = append(ins, reflect.ValueOf(args[i]))
			continue
		}
		if arg.Kind() == reflect.String { // json
			switch tag.Kind() {
			case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
				v := reflect.New(tag)
				if err := json.Unmarshal([]byte(args[i].(string)), v.Interface()); err != nil {
					callback.OnError(constant.ErrArgs.ErrCode, err.Error())
					return
				}
				ins = append(ins, v.Elem())
				continue
			}
		}
		if isNumeric(arg.Kind()) && isNumeric(tag.Kind()) {
			v := reflect.Zero(tag).Interface()
			setNumeric(args[i], &v)
			ins = append(ins, reflect.ValueOf(v))
			continue
		}
		callback.OnError(constant.ErrArgs.ErrCode, "go code error: fn in args type is not match")
		return
	}
	var lastErr bool
	if numOut := fnt.NumOut(); numOut > 0 {
		lastErr = fnt.Out(numOut - 1).Implements(reflect.TypeOf(error(nil)).Elem())
	}
	//fmt.Println("fnv:", fnv.Interface(), "ins:", ins)
	outs := fnv.Call(ins)

	outVals := make([]any, 0, len(outs))
	for i := 0; i < len(outs); i++ {
		outVals = append(outVals, outs[i].Interface())
	}
	if lastErr {
		if last := outVals[len(outVals)-1]; last != nil {
			callback.OnError(10000, last.(error).Error())
			return
		}

		outVals = outVals[:len(outVals)-1]
	}
	// 将map和slice的nil转换为非nil
	for i := 0; i < len(outVals); i++ {
		switch outs[i].Kind() {
		case reflect.Map:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeMap(outs[i].Type()).Interface()
			}
		case reflect.Slice:
			if outs[i].IsNil() {
				outVals[i] = reflect.MakeSlice(outs[i].Type(), 0, 0).Interface()
			}
		}
	}
	var jsonVal any
	if len(outVals) == 1 {
		jsonVal = outVals[0]
	} else {
		jsonVal = outVals
	}
	jsonData, err := json.Marshal(jsonVal)
	if err != nil {
		callback.OnError(constant.ErrArgs.ErrCode, err.Error())
		return
	}
	callback.OnSuccess(string(jsonData))
}
