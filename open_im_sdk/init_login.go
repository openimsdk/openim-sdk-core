package open_im_sdk

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/internal/login"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	localLog "open_im_sdk/pkg/log"
	"open_im_sdk/sdk_struct"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
)

func SdkVersion() string {
	return constant.SdkVersion + constant.BigVersion + constant.UpdateVersion
}

func SetHeartbeatInterval(heartbeatInterval int) {
	constant.HeartbeatInterval = heartbeatInterval
}

func InitSDK(listener open_im_sdk_callback.OnConnListener, operationID string, config string) bool {
	if UserForSDK != nil {
		fmt.Println(operationID, "Initialize multiple times, use the existing ", UserForSDK, " Previous configuration ", UserForSDK.ImConfig(), " now configuration: ", config)
		return true
	}
	if err := json.Unmarshal([]byte(config), &sdk_struct.SvrConf); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if err := log.InitFromConfig("", int(sdk_struct.SvrConf.LogLevel), true, false, "", 0); err != nil {
		fmt.Println(operationID, "log init failed ", err.Error())
		return false
	}

	localLog.NewPrivateLog("", sdk_struct.SvrConf.LogLevel)
	ctx := mcontext.NewCtx(operationID)
	if !strings.Contains(sdk_struct.SvrConf.ApiAddr, "http") {
		log.ZError(ctx, "api is http protocol, api format is invalid", nil)
		return false
	}
	if !strings.Contains(sdk_struct.SvrConf.WsAddr, "ws") {
		log.ZError(ctx, "ws is ws protocol, ws format is invalid", nil)
		return false
	}

	log.ZInfo(ctx, "InitSDK info", "config", sdk_struct.SvrConf, "sdkVersion", SdkVersion())
	if listener == nil || config == "" {
		log.ZError(ctx, "listener or config is nil", nil)
		return false
	}
	UserForSDK = new(login.LoginMgr)
	return UserForSDK.InitSDK(sdk_struct.SvrConf, listener, operationID)
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, UserForSDK.Login, userID, token)
}

func WakeUp(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.WakeUp)
}

func NetworkChanged(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.WakeUp)
}

func UploadImage(callback open_im_sdk_callback.Base, operationID string, filePath string, token, obj string) string {
	//return UserForSDK.UploadImage(callback, filePath, token, obj, operationID)
	return ""
}

func UploadFile(callback open_im_sdk_callback.SendMsgCallBack, operationID string, filePath string) {
	//UserForSDK.UploadFile(callback, filePath, operationID)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, UserForSDK.Logout)
}

func SetAppBackgroundStatus(callback open_im_sdk_callback.Base, operationID string, isBackground bool) {
	BaseCaller(UserForSDK.SetAppBackgroundStatus, callback, isBackground, operationID)
}

func GetLoginStatus() int32 {
	if UserForSDK == nil {
		log.Error("", "UserForSDK == nil")
		return -1
	}
	if UserForSDK.Ws() == nil {
		log.Error("", "UserForSDK.Ws() == nil")
		return -2
	}
	return UserForSDK.GetLoginStatus()
}

func GetLoginUser() string {
	if UserForSDK == nil {
		log.Error("", "UserForSDK == nil")
		return ""
	}
	return UserForSDK.GetLoginUser()
}
