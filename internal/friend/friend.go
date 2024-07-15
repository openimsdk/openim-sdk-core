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
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/page"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	friendSyncLimit = 100
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
	friendSyncer       *syncer.Syncer[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string]
	blockSyncer        *syncer.Syncer[*model_struct.LocalBlack, syncer.NoResp, [2]string]
	requestRecvSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, syncer.NoResp, [2]string]
	requestSendSyncer  *syncer.Syncer[*model_struct.LocalFriendRequest, syncer.NoResp, [2]string]
	conversationCh     chan common.Cmd2Value
	listenerForService open_im_sdk_callback.OnListenerForService
	friendSyncMutex    sync.Mutex
}

func (f *Friend) initSyncer() {
	f.friendSyncer = syncer.New2[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](
		syncer.WithInsert[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriend) error {
			return f.db.InsertFriend(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriend) error {
			return f.db.DeleteFriendDB(ctx, value.FriendUserID)
		}),
		syncer.WithUpdate[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, server, local *model_struct.LocalFriend) error {
			return f.db.UpdateFriend(ctx, server)
		}),
		syncer.WithUUID[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(value *model_struct.LocalFriend) [2]string {
			return [...]string{value.OwnerUserID, value.FriendUserID}
		}),
		syncer.WithNotice[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, state int, server, local *model_struct.LocalFriend) error {
			switch state {
			case syncer.Insert:
				f.friendListener.OnFriendAdded(*server)
			case syncer.Delete:
				log.ZDebug(ctx, "syncer OnFriendDeleted", "local", local)
				f.friendListener.OnFriendDeleted(*local)
			case syncer.Update:
				f.friendListener.OnFriendInfoChanged(*server)
				if local.Nickname != server.Nickname || local.FaceURL != server.FaceURL || local.Remark != server.Remark {
					if server.Remark != "" {
						server.Nickname = server.Remark
					}
					_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{
						Action: constant.UpdateConFaceUrlAndNickName,
						Args: common.SourceIDAndSessionType{
							SourceID:    server.FriendUserID,
							SessionType: constant.SingleChatType,
							FaceURL:     server.FaceURL,
							Nickname:    server.Nickname,
						},
					}, f.conversationCh)
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{
						Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{
							SessionType: constant.SingleChatType,
							UserID:      server.FriendUserID,
							FaceURL:     server.FaceURL,
							Nickname:    server.Nickname,
						},
					}, f.conversationCh)
				}
			}
			return nil
		}),
		syncer.WithBatchInsert[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, values []*model_struct.LocalFriend) error {
			log.ZDebug(ctx, "BatchInsertFriend", "length", len(values))
			return f.db.BatchInsertFriend(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, _ string) error {
			return f.db.DeleteAllFriend(ctx)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(entityID string) page.PageReq {
			return &relation.GetPaginationFriendsReq{UserID: entityID,
				Pagination: &sdkws.RequestPagination{ShowNumber: 100}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(resp *relation.GetPaginationFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.FriendsInfo)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](constant.GetFriendListRouter),
		syncer.WithFullSyncLimit[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](friendSyncLimit),
	)

	f.blockSyncer = syncer.New[*model_struct.LocalBlack, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return f.db.DeleteBlack(ctx, value.BlockUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return f.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlockUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalBlack) error {
		switch state {
		case syncer.Insert:
			f.friendListener.OnBlackAdded(*server)
		case syncer.Delete:
			f.friendListener.OnBlackDeleted(*local)
		}
		return nil
	})
	f.requestRecvSyncer = syncer.New[*model_struct.LocalFriendRequest, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalFriendRequest) error {
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*server)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*local)
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*server)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*server)
			case constant.FriendResponseDefault:
				f.friendListener.OnFriendApplicationAdded(*server)
			}
		}
		return nil
	})
	f.requestSendSyncer = syncer.New[*model_struct.LocalFriendRequest, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.InsertFriendRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalFriendRequest) error {
		return f.db.DeleteFriendRequestBothUserID(ctx, value.FromUserID, value.ToUserID)
	}, func(ctx context.Context, server *model_struct.LocalFriendRequest, local *model_struct.LocalFriendRequest) error {
		return f.db.UpdateFriendRequest(ctx, server)
	}, func(value *model_struct.LocalFriendRequest) [2]string {
		return [...]string{value.FromUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalFriendRequest) error {
		switch state {
		case syncer.Insert:
			f.friendListener.OnFriendApplicationAdded(*server)
		case syncer.Delete:
			f.friendListener.OnFriendApplicationDeleted(*local)
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				f.friendListener.OnFriendApplicationAccepted(*server)
			case constant.FriendResponseRefuse:
				f.friendListener.OnFriendApplicationRejected(*server)
			}
		}
		return nil
	})
}

func (f *Friend) Db() db_interface.DataBase {
	return f.db
}

func (f *Friend) SetListener(listener func() open_im_sdk_callback.OnFriendshipListener) {
	f.friendListener = open_im_sdk_callback.NewOnFriendshipListenerSdk(listener)
}

func (f *Friend) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	f.listenerForService = listener
}
