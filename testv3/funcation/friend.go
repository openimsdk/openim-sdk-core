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

package funcation

import (
	"open_im_sdk/internal/login"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

type testInitLister struct {
}

func (t *testInitLister) OnUserTokenExpired() {
	log.Info("", utils.GetSelfFuncName())
}
func (t *testInitLister) OnConnecting() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testInitLister) OnConnectSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testInitLister) OnConnectFailed(ErrCode int32, ErrMsg string) {
	log.Info("", utils.GetSelfFuncName(), ErrCode, ErrMsg)
}

func (t *testInitLister) OnKickedOffline() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testInitLister) OnSelfInfoUpdated(info string) {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testInitLister) OnSuccess() {
	log.Info("", utils.GetSelfFuncName())
}

func (t *testInitLister) OnError(code int32, msg string) {
	log.Info("", utils.GetSelfFuncName(), code, msg)
}

func ReliabilityInitAndLogin(index int, uid, token string) {
	cf := sdk_struct.IMConfig{
		ApiAddr:    APIADDR,
		WsAddr:     WSADDR,
		PlatformID: PlatformID,
		DataDir:    "./../",
		LogLevel:   LogLevel,
	}

	log.Info("", "DoReliabilityTest", uid, token, WSADDR, APIADDR)

	var testinit testInitLister

	lg := new(login.LoginMgr)
	lg.InitSDK(cf, &testinit)
	log.Info(uid, "new login ", lg)
	allLoginMgr[index].mgr = lg
	log.Info(uid, "InitSDK ", cf)
	lg.Login(ctx, uid, token)
}
