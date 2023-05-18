// Copyright 2021 OpenIM Corporation
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
	"errors"
	"fmt"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func NewFriend(loginUserID string, db db_interface.DataBase, user *user.User, conversationCh chan common.Cmd2Value) *Friend {
	f := &Friend{loginUserID: loginUserID, db: db, user: user, conversationCh: conversationCh}
	f.initSyncer()
	return f
}

type Friend struct {
	friendListener     open_im_sdk_callback.OnFriendshipListenerSdk
	loginUserID        string
	db                 db_interface.DataBase
	user               *user.User
	friendSyncer       *syncer.Syncer[*model_struct.LocalFriend, [2]string]
	blockSyncer        *syncer.Syncer[*model_struct.LocalBlack, [2]string]
	requestRecvSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, [2]string]
	requestSendSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, [2]string]
	loginTime          int64
	conversationCh     chan common.Cmd2Value
	listenerForService open_im_sdk_callback.OnListenerForService
}

func (f *Friend) initSyncer() {
	f.friendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.InsertFriend(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriend) error {
		return f.db.DeleteFriendDB(ctx, value.FriendUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriend, local *model_struct.LocalFriend) error {
		return f.db.UpdateFriend(ctx, server)
	}, func(value *model_struct.LocalFriend) [2]string {
		return [...]string{value.OwnerUserID, value.FriendUserID}
	}, nil, func(ctx context.Context, state int, value *model_struct.LocalFriend) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendAdded(*value)
		case syncer.Delete:
			f.friendListener.OnFriendDeleted(*value)
		case syncer.Update:
			f.friendListener.OnFriendInfoChanged(*value)
			_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: value.FriendUserID, SessionType: constant.SingleChatType}}, f.conversationCh)
			_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: value.FriendUserID, FaceURL: value.FaceURL, Nickname: value.Nickname}}, f.conversationCh)
		}
		return nil
	})
	f.blockSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.DeleteBlack(ctx, value.BlockUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return f.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlockUserID}
	}, nil, func(ctx context.Context, state int, value *model_struct.LocalBlack) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnBlackAdded(*value)
		case syncer.Delete:
			f.friendListener.OnBlackDeleted(*value)
		}
		return nil
	})
	f.requestRecvSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, value *model_struct.LocalFriendRequest) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*value)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*value)
		case syncer.Update:
			switch value.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*value)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*value)
			}
		}
		return nil
	})
	f.requestSendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, value *model_struct.LocalFriendRequest) error {
		if f.friendListener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*value)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*value)
		case syncer.Update:
			switch value.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*value)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*value)
			}
		}
		return nil
	})
}

func (f *Friend) LoginTime() int64 {
	return f.loginTime
}

func (f *Friend) SetLoginTime(loginTime int64) {
	f.loginTime = loginTime
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func (f *Friend) SetListener(listener open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = open_im_sdk_callback.NewOnFriendshipListenerSdk(listener)
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}

func (f *Friend) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := f.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "doNotification error", err, "msg", msg)
		}
	}()
}

func (f *Friend) syncApplication(ctx context.Context, from *sdkws.FromToUserID) error {
	if from.FromUserID == f.loginUserID {
		return f.SyncSelfFriendApplication(ctx)
		// send to me
	} else if from.ToUserID == f.loginUserID {
		return f.SyncFriendApplication(ctx)
	}
	return fmt.Errorf("friend application notification error, fromUserID: %s, toUserID: %s", from.FromUserID, from.ToUserID)
}

func (f *Friend) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if f.friendListener == nil {
		return errors.New("f.friendListener == nil")
	}
	if msg.SendTime < f.loginTime || f.loginTime == 0 {
		return errors.New("ignore notification")
	}
	switch msg.ContentType {
	case constant.FriendApplicationNotification:
		tips := sdkws.FriendApplicationTips{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncApplication(ctx, tips.FromToUserID)
	case constant.FriendApplicationApprovedNotification:
		var tips sdkws.FriendApplicationApprovedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if err := f.SyncFriendList(ctx); err != nil {
			return err
		}
		return f.syncApplication(ctx, tips.FromToUserID)
	case constant.FriendApplicationRejectedNotification:
		var tips sdkws.FriendApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return f.syncApplication(ctx, tips.FromToUserID)
	case constant.FriendAddedNotification:
		return f.SyncFriendList(ctx)
	case constant.FriendDeletedNotification:
		var tips sdkws.FriendDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendList(ctx)
		}
		return nil
	case constant.FriendRemarkSetNotification:
		var tips sdkws.FriendInfoChangedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendList(ctx)
		}
		return nil
	case constant.FriendInfoUpdatedNotification:
		return f.SyncFriendList(ctx)
	case constant.BlackAddedNotification:
		var tips sdkws.BlackAddedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	case constant.BlackDeletedNotification:
		var tips sdkws.BlackDeletedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		if tips.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
}
