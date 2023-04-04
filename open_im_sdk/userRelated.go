package open_im_sdk

import (
	"errors"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
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

// 用于web和pc的userMap
var UserRouterMap map[string]*login.LoginMgr

// 客户端独立的user类
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
func InitOnce(config *sdk_struct.IMConfig) bool {
	sdk_struct.SvrConf = *config
	return true
}

func CheckToken(userID, token string, operationID string) error {
	err, _ := login.CheckToken(userID, token, operationID)
	return err
}

func CheckResourceLoad(uSDK *login.LoginMgr) error {
	if uSDK == nil {
		//	callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return utils.Wrap(errors.New("CheckResourceLoad failed uSDK == nil "), "")
	}
	if uSDK.Friend() == nil || uSDK.User() == nil || uSDK.Group() == nil || uSDK.Conversation() == nil ||
		uSDK.Full() == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed, resource nil "), "")
	}
	return nil
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
			var temp string
			switch x := rc.(type) {
			case string:
				temp = errors.New(x).Error()
			case error:
				buf := make([]byte, 1<<20)
				runtime.Stack(buf, true)
				temp = x.Error()
			default:
				temp = errors.New("unknown panic").Error()
			}
			callback.OnError(constant.ErrArgs.ErrCode, temp)
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
func SendMessageCaller(funcName interface{}, callback open_im_sdk_callback.SendMsgCallBack, args ...interface{}) {
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
