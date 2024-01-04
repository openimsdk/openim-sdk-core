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

package testv2

import (
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"testing"
	"time"
)

func Test_ChangeInputState(t *testing.T) {
	for {
		err := open_im_sdk.UserForSDK.Conversation().ChangeInputStates(ctx, "sg_2309860938", true)
		if err != nil {
			log.ZError(ctx, "ChangeInputStates", err)
		}
		time.Sleep(time.Second * 1)
	}
}

func Test_Empty(t *testing.T) {
	for {
		time.Sleep(time.Second * 1)
	}
}

func Test_RunWait(t *testing.T) {
	time.Sleep(time.Second * 10)
}
