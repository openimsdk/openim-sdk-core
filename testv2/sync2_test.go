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
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"testing"
	"time"
)

func TestSyncFriend2(t *testing.T) {
	for i := 0; ; i++ {
		if err := open_im_sdk.UserForSDK.Friend().IncrSyncFriends(ctx); err != nil {
			t.Log("IncrSyncFriends error-->", err)
			continue
		}
		time.Sleep(time.Second)
	}
}

func TestSyncJoinGroup2(t *testing.T) {
	for i := 0; ; i++ {
		if err := open_im_sdk.UserForSDK.Group().IncrSyncJoinGroup(ctx); err != nil {
			t.Log("IncrSyncJoinGroup error-->", err)
			continue
		}
		time.Sleep(time.Second)
	}
}

func TestSyncGroupMember2(t *testing.T) {
	for i := 0; ; i++ {
		if err := open_im_sdk.UserForSDK.Group().IncrSyncJoinGroupMember(ctx); err != nil {
			t.Log("IncrSyncGroupAndMember error-->", err)
			continue
		}
		time.Sleep(time.Second)
	}
}

func TestName(t *testing.T) {
	for i := 1; i <= 600; i++ {
		_, err := open_im_sdk.UserForSDK.Group().CreateGroup(ctx, &group.CreateGroupReq{
			GroupInfo: &sdkws.GroupInfo{
				GroupType: constant.WorkingGroup,
				GroupName: fmt.Sprintf("group_%d", i),
			},
			MemberUserIDs: []string{"9556972319", "9719689061", "3872159645"},
		})
		if err != nil {
			log.ZError(ctx, "group create failed", err, "index", i)
		}
	}
}
