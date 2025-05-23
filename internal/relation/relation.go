package relation

import (
	"context"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
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
	friendSyncLimit int64 = 10000
)

func NewFriend(loginUserID string, db db_interface.DataBase, user *user.User, conversationCh chan common.Cmd2Value) *Relation {
	r := &Relation{loginUserID: loginUserID, db: db, user: user, conversationCh: conversationCh}
	r.initSyncer()
	return r
}

type Relation struct {
	friendshipListener open_im_sdk_callback.OnFriendshipListenerSdk
	loginUserID        string
	db                 db_interface.DataBase
	user               *user.User
	friendSyncer       *syncer.Syncer[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string]
	blackSyncer        *syncer.Syncer[*model_struct.LocalBlack, syncer.NoResp, [2]string]
	conversationCh     chan common.Cmd2Value
	listenerForService open_im_sdk_callback.OnListenerForService
	relationSyncMutex  sync.Mutex
}

func (r *Relation) initSyncer() {
	r.friendSyncer = syncer.New2[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](
		syncer.WithInsert[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriend) error {
			return r.db.InsertFriend(ctx, value)
		}),
		syncer.WithDelete[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, value *model_struct.LocalFriend) error {
			return r.db.DeleteFriendDB(ctx, value.FriendUserID)
		}),
		syncer.WithUpdate[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, server, local *model_struct.LocalFriend) error {
			r.user.UserCache.Delete(server.FriendUserID)
			return r.db.UpdateFriend(ctx, server)
		}),
		syncer.WithUUID[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(value *model_struct.LocalFriend) [2]string {
			return [...]string{value.OwnerUserID, value.FriendUserID}
		}),
		syncer.WithNotice[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, state int, server, local *model_struct.LocalFriend) error {
			switch state {
			case syncer.Insert:
				r.friendshipListener.OnFriendAdded(*server)
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
				}, r.conversationCh)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{
					Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{
						SessionType: constant.SingleChatType,
						UserID:      server.FriendUserID,
						FaceURL:     server.FaceURL,
						Nickname:    server.Nickname,
					},
				}, r.conversationCh)
			case syncer.Delete:
				log.ZDebug(ctx, "syncer OnFriendDeleted", "local", local)
				r.friendshipListener.OnFriendDeleted(*local)
			case syncer.Update:
				r.friendshipListener.OnFriendInfoChanged(*server)
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
					}, r.conversationCh)
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{
						Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{
							SessionType: constant.SingleChatType,
							UserID:      server.FriendUserID,
							FaceURL:     server.FaceURL,
							Nickname:    server.Nickname,
						},
					}, r.conversationCh)
				}
			}
			return nil
		}),
		syncer.WithBatchInsert[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, values []*model_struct.LocalFriend) error {
			log.ZDebug(ctx, "BatchInsertFriend", "length", len(values))
			return r.db.BatchInsertFriend(ctx, values)
		}),
		syncer.WithDeleteAll[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(ctx context.Context, _ string) error {
			return r.db.DeleteAllFriend(ctx)
		}),
		syncer.WithBatchPageReq[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(entityID string) page.PageReq {
			return &relation.GetPaginationFriendsReq{UserID: entityID,
				Pagination: &sdkws.RequestPagination{ShowNumber: 100}}
		}),
		syncer.WithBatchPageRespConvertFunc[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](func(resp *relation.GetPaginationFriendsResp) []*model_struct.LocalFriend {
			return datautil.Batch(ServerFriendToLocalFriend, resp.FriendsInfo)
		}),
		syncer.WithReqApiRouter[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](api.GetFriendList.Route()),
		syncer.WithFullSyncLimit[*model_struct.LocalFriend, relation.GetPaginationFriendsResp, [2]string](friendSyncLimit),
	)

	r.blackSyncer = syncer.New[*model_struct.LocalBlack, syncer.NoResp, [2]string](func(ctx context.Context, value *model_struct.LocalBlack) error {
		return r.db.InsertBlack(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalBlack) error {
		return r.db.DeleteBlack(ctx, value.BlockUserID)
	}, func(ctx context.Context, server *model_struct.LocalBlack, local *model_struct.LocalBlack) error {
		return r.db.UpdateBlack(ctx, server)
	}, func(value *model_struct.LocalBlack) [2]string {
		return [...]string{value.OwnerUserID, value.BlockUserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalBlack) error {
		switch state {
		case syncer.Insert:
			r.friendshipListener.OnBlackAdded(*server)
		case syncer.Delete:
			r.friendshipListener.OnBlackDeleted(*local)
		}
		return nil
	})
}

func (r *Relation) Db() db_interface.DataBase {
	return r.db
}

func (r *Relation) SetListener(listener func() open_im_sdk_callback.OnFriendshipListener) {
	r.friendshipListener = open_im_sdk_callback.NewOnFriendshipListenerSdk(listener)
}

func (r *Relation) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	r.listenerForService = listener
}
