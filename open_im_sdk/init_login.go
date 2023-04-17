package open_im_sdk

import (
	"encoding/json"
	"fmt"
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
	if userForSDK != nil {
		fmt.Println(operationID, "Initialize multiple times, use the existing ", userForSDK, " Previous configuration ", userForSDK.ImConfig(), " now configuration: ", config)
		return true
	}
	if err := json.Unmarshal([]byte(config), &sdk_struct.SvrConf); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	log.NewPrivateLog("", sdk_struct.SvrConf.LogLevel)
	if !strings.Contains(sdk_struct.SvrConf.ApiAddr, "http") {
		log.Error(operationID, "api is http protocol", sdk_struct.SvrConf.ApiAddr)
		return false
	}
	if !strings.Contains(sdk_struct.SvrConf.WsAddr, "ws") {
		log.Error(operationID, "ws is ws protocol", sdk_struct.SvrConf.ApiAddr)
		return false
	}

	log.Info(operationID, "config ", config, sdk_struct.SvrConf)
	log.NewInfo(operationID, utils.GetSelfFuncName(), config, SdkVersion())
	if listener == nil || config == "" {
		log.Error(operationID, "listener or config is nil")
		return false
	}

	userForSDK = new(login.LoginMgr)

	return userForSDK.InitSDK(sdk_struct.SvrConf, listener, operationID)
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, userForSDK.WakeUp, userID, token)
}

func WakeUp(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.WakeUp)
}

func NetworkChanged(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.WakeUp)
}

func UploadImage(callback open_im_sdk_callback.Base, operationID string, filePath string, token, obj string) string {
	//return userForSDK.UploadImage(callback, filePath, token, obj, operationID)
	return ""
}

func UploadFile(callback open_im_sdk_callback.SendMsgCallBack, operationID string, filePath string) {
	//userForSDK.UploadFile(callback, filePath, operationID)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Logout)
}

func SetAppBackgroundStatus(callback open_im_sdk_callback.Base, operationID string, isBackground bool) {
	BaseCaller(userForSDK.SetAppBackgroundStatus, callback, isBackground, operationID)
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
