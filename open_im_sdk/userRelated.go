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

// Caller is an interface that defines the SDK's basic and message sending caller.
type Caller interface {
	BaseCaller(funcName interface{}, base open_im_sdk_callback.Base, args ...interface{})
	SendMessageCaller(funcName interface{}, messageCallback open_im_sdk_callback.SendMsgCallBack, args ...interface{})
}

var (
	UserSDKRwLock sync.RWMutex
	// userMap for web and pc
	UserRouterMap map[string]*login.LoginMgr
	// Client-independent user class
	userForSDK *login.LoginMgr
)

// init initializes the UserRouterMap to hold a map of string keys and *login.LoginMgr values.
func init() {
	UserSDKRwLock.Lock()
	defer UserSDKRwLock.Unlock()
	UserRouterMap = make(map[string]*login.LoginMgr)
}

// GetUserWorker returns a user's login manager by its ID.
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

// InitOnce initializes the SDK by setting the server configuration.
func InitOnce(config *sdk_struct.IMConfig) bool {
	sdk_struct.SvrConf = *config
	return true
}

// CheckToken checks user authentication token.
func CheckToken(userID, token string, operationID string) error {
	err, _ := login.CheckToken(userID, token, operationID)
	return err
}

// CheckResourceLoad checks the SDK's resource load status.
func CheckResourceLoad(uSDK *login.LoginMgr) error {
	if uSDK == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed uSDK == nil "), "")
	}
	if uSDK.Friend() == nil || uSDK.User() == nil || uSDK.Group() == nil || uSDK.Conversation() == nil ||
		uSDK.Full() == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed, resource nil "), "")
	}
	return nil
}

type name struct {
	
}

var ErrNotSetCallback = errors.New("not set callback to call")
var ErrNotSetFunc = errors.New("not set func to call")

// BaseCaller calls the SDK's basic caller by checking the arguments and verifying the callback.
// First, it checks that the number of arguments is correct and gets the operation ID.
// Then, it checks that the resources have been loaded, and returns an error if they have not.
// Finally, it uses reflection to call the function, passing in the callback and arguments, and runs the function in a different goroutine.
// If a panic occurs, it converts the panic to a string and returns its error through the callback.
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

// SendMessageCaller sends a message by calling the SDK's message sender.
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
