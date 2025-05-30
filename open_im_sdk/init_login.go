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
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cliconf"
	pbConstant "github.com/openimsdk/protocol/constant"

	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/openim-sdk-core/v3/version"

	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
)

func GetSdkVersion() string {
	return version.Version
}

const (
	rotationTime uint = 24
)

func InitSDK(listener open_im_sdk_callback.OnConnListener, operationID string, config string) bool {
	var configArgs sdk_struct.IMConfig
	if err := json.Unmarshal([]byte(config), &configArgs); err != nil {
		fmt.Println(operationID, "Unmarshal failed ", err.Error(), config)
		return false
	}
	if configArgs.PlatformID == 0 {
		return false
	}
	var logRemainCount uint32
	if configArgs.LogRemainCount > 0 {
		logRemainCount = configArgs.LogRemainCount
	} else {
		logRemainCount = 1
	}
	if err := log.InitLoggerFromConfig("open-im-sdk-core", "", configArgs.SystemType, pbConstant.PlatformID2Name[int(configArgs.PlatformID)], int(configArgs.LogLevel), configArgs.IsLogStandardOutput, false, configArgs.LogFilePath, uint(logRemainCount), rotationTime, version.Version, true); err != nil {
		fmt.Println(operationID, "log init failed ", err.Error())
	}
	fmt.Println("init log success")
	ctx := mcontext.NewCtx(operationID)
	if !strings.Contains(configArgs.ApiAddr, "http") {
		log.ZError(ctx, "api is http protocol, api format is invalid", nil)
		return false
	}
	if !strings.Contains(configArgs.WsAddr, "ws") {
		log.ZError(ctx, "ws is ws protocol, ws format is invalid", nil)
		return false
	}

	log.ZInfo(ctx, "InitSDK info", "config", configArgs)
	if listener == nil || config == "" {
		log.ZError(ctx, "listener or config is nil", nil)
		return false
	}
	return IMUserContext.InitSDK(&configArgs, listener)
}
func UnInitSDK(_ string) {
	IMUserContext.UnInitSDK()
}

func GetLoginUserID() string {
	if IMUserContext == nil {
		return ""
	}
	return IMUserContext.GetLoginUserID()
}

func Login(callback open_im_sdk_callback.Base, operationID string, userID, token string) {
	call(callback, operationID, IMUserContext.Login, userID, token)
}

func Logout(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Logout)
}

func SetAppBackgroundStatus(callback open_im_sdk_callback.Base, operationID string, isBackground bool) {
	call(callback, operationID, IMUserContext.SetAppBackgroundStatus, isBackground)
}
func NetworkStatusChanged(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.NetworkStatusChanged)
}

func GetLoginStatus(operationID string) int {
	return IMUserContext.GetLoginStatus(ccontext.WithOperationID(context.Background(), operationID))
}

func (u *UserContext) Login(ctx context.Context, userID, token string) error {
	cliconf.SetLoginUserID(u.loginUserID)
	return u.login(ctx, userID, token)
}

func (u *UserContext) Logout(ctx context.Context) error {
	return u.logout(ctx, false)
}

func (u *UserContext) SetAppBackgroundStatus(ctx context.Context, isBackground bool) error {
	return u.setAppBackgroundStatus(ctx, isBackground)
}
func (u *UserContext) NetworkStatusChanged(ctx context.Context) {
	u.longConnMgr.Close(ctx)
}
func (u *UserContext) GetLoginStatus(ctx context.Context) int {
	return u.getLoginStatus(ctx)
}
