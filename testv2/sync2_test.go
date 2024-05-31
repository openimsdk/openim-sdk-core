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
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"testing"
)

func TestSyncFriend2(t *testing.T) {
	for i := 0; ; i++ {
		if err := open_im_sdk.UserForSDK.Friend().IncrSyncFriends(ctx); err != nil {
			t.Log("error-->", err)
			//panic(err)
		}
	}
	//select {}
}
