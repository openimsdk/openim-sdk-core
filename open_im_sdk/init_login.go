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
	var configArgs sdk_struct.IMConfig
	if err := json.Unmarshal([]byte(config), &configArgs); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if err := log.InitFromConfig("", int(configArgs.LogLevel), true, false, "", 0); err != nil {
		fmt.Println(operationID, "log init failed ", err.Error())
		return false
	}

	localLog.NewPrivateLog("", configArgs.LogLevel)
	ctx := mcontext.NewCtx(operationID)
	if !strings.Contains(configArgs.ApiAddr, "http") {
		log.ZError(ctx, "api is http protocol, api format is invalid", nil)
		return false
	}
	if !strings.Contains(configArgs.WsAddr, "ws") {
		log.ZError(ctx, "ws is ws protocol, ws format is invalid", nil)
		return false
	}

	log.ZInfo(ctx, "InitSDK info", "config", configArgs, "sdkVersion", SdkVersion())
	if listener == nil || config == "" {
		log.ZError(ctx, "listener or config is nil", nil)
		return false
	}
	UserForSDK = new(login.LoginMgr)
	return UserForSDK.InitSDK(configArgs, listener, operationID)
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

func GetLoginStatus() int {
	if UserForSDK == nil {
		log.Error("", "UserForSDK == nil")
		return -1
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
