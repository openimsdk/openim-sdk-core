package open_im_sdk

import (
	"errors"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"reflect"
	"runtime"
	"sync"
)

func init() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*login.LoginMgr, 0)
}

var UserSDKRwLock sync.RWMutex

//用于web和pc的userMap
var UserRouterMap map[string]*login.LoginMgr

//客户端独立的user类
var userForSDK *login.LoginMgr

func GetUserWorker(uid string) *login.LoginMgr {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	v, ok := UserRouterMap[uid]
	if ok {
		return v
	}
	UserRouterMap[uid] = new(login.LoginMgr)

	return UserRouterMap[uid]
}

type Caller interface {
	BaseCaller(funcName interface{}, base open_im_sdk_callback.Base, args ...interface{})
	SendMessageCaller(funcName interface{}, messageCallback open_im_sdk_callback.SendMsgCallBack, args ...interface{})
}

type name struct {
}

var ErrNotSetCallback = errors.New("not set callback to call")
var ErrNotSetFunc = errors.New("not set func to call")

func BaseCaller(funcName interface{}, callback open_im_sdk_callback.Base, args ...interface{}) {
	var operationID string
	if len(args) <= 0 {
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		return
	}
	if v, ok := args[len(args)-1].(string); ok {
		operationID = v
	} else {
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	defer func() {
		if rc := recover(); rc != nil {
			log.Error(operationID, "err:", rc)
		}
	}()
	if funcName == nil {
		panic(utils.Wrap(ErrNotSetFunc, ""))
	}
	var refFuncName reflect.Value
	var values []reflect.Value
	refFuncName = reflect.ValueOf(funcName)
	if callback != nil {
		values = append(values, reflect.ValueOf(callback))
	} else {
		log.Error("AsyncCallWithCallback", "not set callback")
		panic(ErrNotSetCallback)
	}
	for i := 0; i < len(args); i++ {
		values = append(values, reflect.ValueOf(args[i]))
	}
	pc, _, _, _ := runtime.Caller(1)
	funcNameString := utils.CleanUpfuncName(runtime.FuncForPC(pc).Name())
	log.Debug(operationID, funcNameString, "input args:", args)
	go refFuncName.Call(values)
}
