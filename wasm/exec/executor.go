// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package exec

import (
	"errors"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"runtime"
	"syscall/js"
	"time"
)

type CallbackData struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data"`
}

const TIMEOUT = 5

var ErrType = errors.New("from javascript data type err")
var PrimaryKeyNull = errors.New("primary key is null err")

var ErrTimoutFromJavaScript = errors.New("invoke javascript timeout，maybe should check  function from javascript")
var jsErr = js.Global().Get("Error")

func Exec(args ...interface{}) (output interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = utils.Wrap(errors.New(x), "")
			case error:
				err = x
			default:
				err = utils.Wrap(errors.New("unknown panic"), "")
			}
		}
	}()
	thenChannel := make(chan []js.Value)
	defer close(thenChannel)
	catchChannel := make(chan []js.Value)
	defer close(catchChannel)
	pc, _, _, _ := runtime.Caller(1)
	funcName := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	data := CallbackData{}
	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js then function", "=> (main go context) "+funcName+" with response ", args[0].String())
		thenChannel <- args
		return nil
	})
	defer thenFunc.Release()
	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = utils.Wrap(errors.New(x), "")
				case error:
					err = x
				default:
					err = utils.Wrap(errors.New("unknown panic"), "")
				}
			}
		}()
		log.Debug("js catch function", "=> (main go context) "+funcName+" with respone ", args[0].String())
		catchChannel <- args
		return nil
	})
	defer catchFunc.Release()
	js.Global().Call(utils.FirstLower(funcName), args...).Call("then", thenFunc).Call("catch", catchFunc)
	select {
	case result := <-thenChannel:
		if len(result) > 0 {
			switch result[0].Type() {
			case js.TypeString:
				interErr := utils.JsonStringToStruct(result[0].String(), &data)
				if interErr != nil {
					err = utils.Wrap(err, "return json unmarshal err from javascript")
				}
			case js.TypeObject:
				return result[0], nil

			default:
				err = errors.New("unknown return type from javascript")
			}

		} else {
			err = errors.New("args err,length is 0")
		}

	case catch := <-catchChannel:
		if catch[0].InstanceOf(jsErr) {
			return nil, js.Error{Value: catch[0]}
		} else {
			panic("unknown javascript exception")
		}
	case <-time.After(TIMEOUT * time.Second):
		panic(ErrTimoutFromJavaScript)
	}
	if data.ErrCode != 0 {
		return "", errors.New(data.ErrMsg)
	}
	return data.Data, err
}

func ExtractArrayBuffer(arrayBuffer js.Value) []byte {
	uint8Array := js.Global().Get("Uint8Array").New(arrayBuffer)
	dst := make([]byte, uint8Array.Length())
	js.CopyBytesToGo(dst, uint8Array)
	return dst
}
