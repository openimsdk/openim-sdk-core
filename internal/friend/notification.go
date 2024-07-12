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

package friend

import (
	"context"
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (f *Friend) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		f.friendSyncMutex.Lock()
		defer f.friendSyncMutex.Unlock()

		if err := f.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "doNotification error", err, "msg", msg)
		}
	}()
}

func (f *Friend) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	switch msg.ContentType {
	case constant.FriendApplicationNotification:
		tips := sdkws.FriendApplicationTips{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.SyncBothFriendRequest(ctx,
			tips.FromToUserID.FromUserID, tips.FromToUserID.ToUserID)
	case constant.FriendApplicationApprovedNotification:
		var tips sdkws.FriendApplicationApprovedTips
		err := utils.UnmarshalNotificationElem(msg.Content, &tips)
		if err != nil {
			return err
		}

		if tips.FromToUserID.FromUserID == f.loginUserID {
			err = f.IncrSyncFriends(ctx)
		} else if tips.FromToUserID.ToUserID == f.loginUserID {
			err = f.IncrSyncFriends(ctx)
		}
		if err != nil {
			return err
		}
		return f.SyncBothFriendRequest(ctx, tips.FromToUserID.FromUserID, tips.FromToUserID.ToUserID)
	case constant.FriendApplicationRejectedNotification:
		var tips sdkws.FriendApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.SyncBothFriendRequest(ctx, tips.FromToUserID.FromUserID, tips.FromToUserID.ToUserID)
	case constant.FriendAddedNotification:
		var tips sdkws.FriendAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.Friend != nil && tips.Friend.FriendUser != nil {
			if tips.Friend.FriendUser.UserID == f.loginUserID {
				return f.IncrSyncFriends(ctx)
			} else if tips.Friend.OwnerUserID == f.loginUserID {
				return f.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendDeletedNotification:
		var tips sdkws.FriendDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID != nil {
			if tips.FromToUserID.FromUserID == f.loginUserID {
				return f.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendRemarkSetNotification:
		var tips sdkws.FriendInfoChangedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID != nil {
			if tips.FromToUserID.FromUserID == f.loginUserID {
				return f.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendInfoUpdatedNotification:
		var tips sdkws.UserInfoUpdatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.UserID != f.loginUserID {
			return f.IncrSyncFriends(ctx)
		}
	case constant.BlackAddedNotification:
		var tips sdkws.BlackAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncAllBlackList(ctx)
		}
	case constant.BlackDeletedNotification:
		var tips sdkws.BlackDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncAllBlackList(ctx)
		}
	case constant.FriendsInfoUpdateNotification:

		var tips sdkws.FriendsInfoUpdateTips

		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.ToUserID == f.loginUserID {
			return f.IncrSyncFriends(ctx)
		}
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
	return nil
}
