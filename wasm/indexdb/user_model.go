//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"runtime"
	"syscall/js"
	"time"
)

var ErrType = errors.New("from javascript data type err")
var PrimaryKeyNull = errors.New("primary key is null err")

const TIMEOUT = 5

var ErrTimoutFromJavaScript = errors.New("invoke javascript timeout，maybe should check  function from javascript")
var jsErr = js.Global().Get("Error")

type LocalUsers struct {
}
type CallbackData struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    interface{} `json:"data"`
}

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
		log.Debug("js then func", "=> (main go context) "+funcName+" with respone ", args[0].String())
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
		log.Debug("js catch func", "=> (main go context) "+funcName+" with respone ", args[0].String())
		catchChannel <- args
		return nil
	})
	defer catchFunc.Release()
	js.Global().Call(utils.FirstLower(funcName), args...).Call("then", thenFunc).Call("catch", catchFunc)
	select {
	case result := <-thenChannel:
		interErr := utils.JsonStringToStruct(result[0].String(), &data)
		if interErr != nil {
			err = utils.Wrap(err, "return json unmarshal err from javascript")
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

func (l *LocalUsers) GetLoginUser(userID string) (*model_struct.LocalUser, error) {
	user, err := Exec(userID)
	if err != nil {
		return nil, err
	} else {
		if v, ok := user.(string); ok {
			result := model_struct.LocalUser{}
			err := utils.JsonStringToStruct(v, &result)
			if err != nil {
				return nil, err
			}
			return &result, err
		} else {
			return nil, ErrType
		}

	}
}

func (l *LocalUsers) UpdateLoginUser(user *model_struct.LocalUser) error {
	_, err := Exec(user)
	return err

}
func (l *LocalUsers) UpdateLoginUserByMap(user *model_struct.LocalUser, args map[string]interface{}) error {
	_, err := Exec(user.UserID, args)
	return err
}
func (l *LocalUsers) InsertLoginUser(user *model_struct.LocalUser) error {
	_, err := Exec(utils.StructToJsonString(user))
	return err
}
