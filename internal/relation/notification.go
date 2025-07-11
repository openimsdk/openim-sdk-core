package relation

import (
	"context"
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (r *Relation) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	if err := r.doNotification(ctx, msg); err != nil {
		log.ZError(ctx, "doNotification error", err, "msg", msg)
	}
}

func (r *Relation) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	switch msg.ContentType {
	case constant.FriendApplicationNotification:
		tips := sdkws.FriendApplicationTips{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		r.friendshipListener.OnFriendApplicationAdded(*ServerFriendRequestToLocalFriendRequest(tips.Request))
	case constant.FriendApplicationApprovedNotification:
		var tips sdkws.FriendApplicationApprovedTips
		err := utils.UnmarshalNotificationElem(msg.Content, &tips)
		if err != nil {
			return err
		}
		if tips.Request != nil {
			r.friendshipListener.OnFriendApplicationAccepted(*ServerFriendRequestToLocalFriendRequest(tips.Request))
		}
		return r.IncrSyncFriends(ctx)
	case constant.FriendApplicationRejectedNotification:
		var tips sdkws.FriendApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		r.friendshipListener.OnFriendApplicationRejected(*ServerFriendRequestToLocalFriendRequest(tips.Request))
	case constant.FriendAddedNotification:
		var tips sdkws.FriendAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.Friend != nil && tips.Friend.FriendUser != nil {
			if tips.Friend.FriendUser.UserID == r.loginUserID {
				return r.IncrSyncFriends(ctx)
			} else if tips.Friend.OwnerUserID == r.loginUserID {
				return r.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendDeletedNotification:
		var tips sdkws.FriendDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID != nil {
			if tips.FromToUserID.FromUserID == r.loginUserID {
				return r.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendRemarkSetNotification:
		var tips sdkws.FriendInfoChangedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID != nil {
			if tips.FromToUserID.FromUserID == r.loginUserID {
				return r.IncrSyncFriends(ctx)
			}
		}
	case constant.FriendInfoUpdatedNotification:
		var tips sdkws.UserInfoUpdatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.UserID != r.loginUserID {
			return r.IncrSyncFriends(ctx)
		}
	case constant.BlackAddedNotification:
		var tips sdkws.BlackAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == r.loginUserID {
			return r.SyncAllBlackList(ctx)
		}
	case constant.BlackDeletedNotification:
		var tips sdkws.BlackDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == r.loginUserID {
			return r.SyncAllBlackList(ctx)
		}
	case constant.FriendsInfoUpdateNotification:
		var tips sdkws.FriendsInfoUpdateTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.ToUserID == r.loginUserID {
			return r.IncrSyncFriends(ctx)
		}
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
	return nil
}
