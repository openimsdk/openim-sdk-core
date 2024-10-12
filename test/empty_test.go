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

package test

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/tools/log"
	"testing"
	"time"
)

func Test_Empty(t *testing.T) {
	for {
		time.Sleep(time.Second * 1)
	}
}

func Test_ChangeInputState(t *testing.T) {
	for {
		err := open_im_sdk.UserForSDK.Conversation().ChangeInputStates(ctx, "sg_2309860938", true)
		if err != nil {
			log.ZError(ctx, "ChangeInputStates", err)
		}
		time.Sleep(time.Second * 1)
	}
}

func Test_RunWait(t *testing.T) {
	time.Sleep(time.Second * 10)
}

func Test_OnlineState(t *testing.T) {
	defer func() { select {} }()
	userIDs := []string{
		//"3611802798",
		"2110910952",
	}
	for i := 1; ; i++ {
		time.Sleep(time.Second)
		//open_im_sdk.UserForSDK.LongConnMgr().UnsubscribeUserOnlinePlatformIDs(ctx, userIDs)
		res, err := open_im_sdk.UserForSDK.LongConnMgr().GetUserOnlinePlatformIDs(ctx, userIDs)
		if err != nil {
			t.Logf("@@@@@@@@@@@@=====> <%d> error %s", i, err)
			continue
		}
		t.Logf("@@@@@@@@@@@@=====> <%d> success %+v", i, res)
	}
}
