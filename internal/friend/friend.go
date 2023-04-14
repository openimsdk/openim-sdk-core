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
	"github.com/golang/protobuf/proto"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/syncer"
)

func NewFriend(loginUserID string, db db_interface.DataBase, user *user.User, p *ws.PostApi, conversationCh chan common.Cmd2Value) *Friend {
	f := &Friend{loginUserID: loginUserID, db: db, user: user, p: p, conversationCh: conversationCh}
	f.initSyncer()
	return f
}

type Friend struct {
	friendListener     open_im_sdk_callback.OnFriendshipListener
	loginUserID        string
	db                 db_interface.DataBase
	user               *user.User
	p                  *ws.PostApi
	friendSyncer       *syncer.Syncer[*model_struct.LocalFriend, [2]string]
	blockSyncer        *syncer.Syncer[*model_struct.LocalBlack, [2]string]
	requestRecvSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, string]
	requestSendSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, string]
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
	}, nil, nil)

	f.blockSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.DeleteBlack(ctx, value.BlockUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return f.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlockUserID}
	}, nil, nil)

	f.requestRecvSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) string {
		return value.FromUserID
	}, nil, nil)

	f.requestSendSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) string {
		return value.ToUserID
	}, nil, nil)

}

func (f *Friend) LoginTime() int64 {
	return f.loginTime
}

func (f *Friend) SetLoginTime(loginTime int64) {
	f.loginTime = loginTime
}

func (f *Friend) SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	if listener == nil {
		return
	}
	f.friendListener = listener
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func (f *Friend) SetListener(listener open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = listener
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}

func (f *Friend) DoNotification(ctx context.Context, msg *api.MsgData) {
	go func() {
		if err := f.doNotification(ctx, msg); err != nil {
			// todo log
		}
	}()
}

func (f *Friend) doNotification(ctx context.Context, msg *api.MsgData) error {
	if f.friendListener == nil {
		return errors.New("f.friendListener == nil")
	}
	if msg.SendTime < f.loginTime || f.loginTime == 0 {
		return errors.New("ignore notification")
	}
	var tips api.TipsComm
	if err := proto.Unmarshal(msg.Content, &tips); err != nil {
		return err
	}
	switch msg.ContentType {
	case constant.FriendApplicationNotification:
		var detail api.FriendApplicationTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncSelfFriendApplication(ctx)
		} else if detail.FromToUserID.ToUserID == f.loginUserID {
			return f.SyncFriendApplication(ctx)
		} else {
			return fmt.Errorf("friend application notification error, fromUserID: %s, toUserID: %s", detail.FromToUserID.FromUserID, detail.FromToUserID.ToUserID)
		}
	case constant.FriendApplicationApprovedNotification:
		var detail api.FriendApplicationApprovedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendApplication(ctx)
		} else if detail.FromToUserID.ToUserID == f.loginUserID {
			return f.SyncSelfFriendApplication(ctx)
		} else {
			return fmt.Errorf("friend application notification error, fromUserID: %s, toUserID: %s", detail.FromToUserID.FromUserID, detail.FromToUserID.ToUserID)
		}
	case constant.FriendApplicationRejectedNotification:
		var detail api.FriendApplicationRejectedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendApplication(ctx)
		} else if detail.FromToUserID.ToUserID == f.loginUserID {
			return f.SyncSelfFriendApplication(ctx)
		} else {
			return fmt.Errorf("friend application notification error, fromUserID: %s, toUserID: %s", detail.FromToUserID.FromUserID, detail.FromToUserID.ToUserID)
		}
	case constant.FriendAddedNotification:
		//var detail api.FriendAddedTips
		//if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
		//	return err
		//}
		return f.SyncFriendList(ctx)
	case constant.FriendDeletedNotification:
		var detail api.FriendDeletedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendList(ctx)
		}
		return nil
	case constant.FriendRemarkSetNotification:
		var detail api.FriendInfoChangedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncFriendList(ctx)
		}
		return nil
	case constant.FriendInfoUpdatedNotification:
		//var detail api.UserInfoUpdatedTips
		//if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
		//	return err
		//}
		return f.SyncFriendList(ctx)
	case constant.BlackAddedNotification:
		var detail api.BlackAddedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	case constant.BlackDeletedNotification:
		var detail api.BlackDeletedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.FromToUserID.FromUserID == f.loginUserID {
			return f.SyncBlackList(ctx)
		}
		return nil
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
}
