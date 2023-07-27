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
	"open_im_sdk/open_im_sdk"
	"testing"
)

//func Test_SyncJoinGroup(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Group().SyncJoinedGroup(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func Test_SyncGroupMember(t *testing.T) {
	//groups, err := open_im_sdk.UserForSDK.Group().GetServerJoinGroup(ctx)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//for _, group := range groups {
	//	err := open_im_sdk.UserForSDK.Group().SyncGroupMember(ctx, group.GroupID)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}
}

//func Test_SyncSelfGroupApplication(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Group().SyncSelfGroupApplication(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

//func Test_SyncAdminGroupApplication(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Group().SyncAdminGroupApplication(ctx)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func Test_SyncSelfFriendApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncAllSelfFriendApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncFriendApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncAllFriendApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncFriend(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncAllFriendList(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncBlack(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncAllBlackList(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncAllConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SyncAllConversations(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
