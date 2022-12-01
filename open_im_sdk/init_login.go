package open_im_sdk

import (
	"encoding/json"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"strings"
)

func SdkVersion() string {
	return constant.SdkVersion + constant.BigVersion + constant.UpdateVersion
}

func SetHeartbeatInterval(heartbeatInterval int) {
	constant.HeartbeatInterval = heartbeatInterval
}

func InitSDK(listener open_im_sdk_callback.OnConnListener, operationID string, config string) bool {
	log.Debug(operationID, "init sdk args:", config)
	if err := json.Unmarshal([]byte(config), &sdk_struct.SvrConf); err != nil {
		log.Error(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if !strings.Contains(sdk_struct.SvrConf.ApiAddr, "http") {
		log.Error(operationID, "api is http protocol", sdk_struct.SvrConf.ApiAddr)
		return false
	}
	if !strings.Contains(sdk_struct.SvrConf.WsAddr, "ws") {
		log.Error(operationID, "ws is ws protocol", sdk_struct.SvrConf.ApiAddr)
		return false
	}
	log.NewPrivateLog("", sdk_struct.SvrConf.LogLevel)
	log.Info(operationID, "config ", config, sdk_struct.SvrConf)
	log.NewInfo(operationID, utils.GetSelfFuncName(), config, SdkVersion())
	if listener == nil || config == "" {
		log.Error(operationID, "listener or config is nil")
		return false
	}
	if userForSDK != nil {
		log.Warn(operationID, "Initialize multiple times, call logout")
		userForSDK.Logout(nil, utils.OperationIDGenerator())
	}
	userForSDK = new(login.LoginMgr)

	return userForSDK.InitSDK(sdk_struct.SvrConf, listener, operationID)
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	if userForSDK == nil {
		callback.OnError(constant.ErrArgs.ErrCode, constant.ErrArgs.ErrMsg)
		return
	}
	userForSDK.Login(callback, userID, token, operationID)
}

func WakeUp(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.WakeUp(callback, operationID)

}

func UploadImage(callback open_im_sdk_callback.Base, operationID string, filePath string, token, obj string) string {
	return userForSDK.UploadImage(callback, filePath, token, obj, operationID)
}

func UploadFile(callback open_im_sdk_callback.SendMsgCallBack, operationID string, filePath string) {
	userForSDK.UploadFile(callback, filePath, operationID)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		log.Error("callback is nil")
		return
	}
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}

	userForSDK.Logout(callback, operationID)
	userForSDK = nil
}

func GetLoginStatus() int32 {
	if userForSDK == nil {
		log.Error("", "userForSDK == nil")
		return -1
	}
	if userForSDK.Ws() == nil {
		log.Error("", "userForSDK.Ws() == nil")
		return -2
	}
	return userForSDK.GetLoginStatus()
}

func GetLoginUser() string {
	if userForSDK == nil {
		log.Error("", "userForSDK == nil")
		return ""
	}
	return userForSDK.GetLoginUser()
}
