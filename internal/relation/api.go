package relation

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/datafetcher"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	sdk "github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"

	"github.com/openimsdk/tools/log"
)

func (r *Relation) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string, filterBlack bool) ([]*model_struct.LocalFriend, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		r.db,
		r.friendListTableName(),
		r.loginUserID,
		func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		func(ctx context.Context, values []*model_struct.LocalFriend) error {
			return r.db.BatchInsertFriend(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, bool, error) {
			localFriends, err := r.db.GetFriendInfoList(ctx, userIDs)
			return localFriends, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			serverFriend, err := r.getDesignatedFriends(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerFriendToLocalFriend, serverFriend), nil
		},
	)
	localFriendList, err := dataFetcher.FetchMissingAndFillLocal(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	if !filterBlack {
		return localFriendList, nil
	}
	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "localFriendList", localFriendList)
	blackList, err := r.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	if len(blackList) == 0 {
		return localFriendList, nil
	}

	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "blackList", blackList)
	m := datautil.SliceSetAny(blackList, func(e *model_struct.LocalBlack) string {
		return e.BlockUserID
	})
	var res []*model_struct.LocalFriend
	for _, localFriend := range localFriendList {
		if _, ok := m[localFriend.FriendUserID]; !ok {
			res = append(res, localFriend)
		}
	}
	return res, nil
}

func (r *Relation) AddFriend(ctx context.Context, req *relation.ApplyToAddFriendReq) error {
	return r.addFriend(ctx, req)
}

func (r *Relation) GetFriendApplicationListAsRecipient(ctx context.Context, req *sdk.GetFriendApplicationListAsRecipientReq) ([]*model_struct.LocalFriendRequest, error) {
	friendRequests, err := r.getRecvFriendApplicationList(ctx, req.HandleResults, utils.GetPageNumber(req.Offset, req.Count), req.Count)
	if err != nil {
		return nil, err
	}
	return datautil.Batch(ServerFriendRequestToLocalFriendRequest, friendRequests), nil
}

func (r *Relation) GetFriendApplicationListAsApplicant(ctx context.Context, req *sdk.GetFriendApplicationListAsApplicantReq) ([]*model_struct.LocalFriendRequest, error) {
	friendRequests, err := r.getSelfFriendApplicationList(ctx, utils.GetPageNumber(req.Offset, req.Count), req.Count)
	if err != nil {
		return nil, err
	}
	return datautil.Batch(ServerFriendRequestToLocalFriendRequest, friendRequests), nil
}

func (r *Relation) AcceptFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return r.RespondFriendApply(ctx, &relation.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: r.loginUserID, HandleResult: constant.FriendResponseAgree, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (r *Relation) RefuseFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return r.RespondFriendApply(ctx, &relation.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: r.loginUserID, HandleResult: constant.FriendResponseRefuse, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (r *Relation) RespondFriendApply(ctx context.Context, req *relation.RespondFriendApplyReq) error {
	if err := r.addFriendResponse(ctx, req); err != nil {
		return err
	}
	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	if req.HandleResult == constant.FriendResponseAgree {
		_ = r.IncrSyncFriends(ctx)
	}
	return nil
}

func (r *Relation) CheckFriend(ctx context.Context, friendUserIDList []string) ([]*server_api_params.UserIDResult, error) {
	friendList, err := r.db.GetFriendInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	blackList, err := r.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	res := make([]*server_api_params.UserIDResult, 0, len(friendUserIDList))
	for _, v := range friendUserIDList {
		var r server_api_params.UserIDResult
		isBlack := false
		isFriend := false
		for _, b := range blackList {
			if v == b.BlockUserID {
				isBlack = true
				break
			}
		}
		for _, r := range friendList {
			if v == r.FriendUserID {
				isFriend = true
				break
			}
		}
		r.UserID = v
		if isFriend && !isBlack {
			r.Result = 1
		} else {
			r.Result = 0
		}
		res = append(res, &r)
	}
	return res, nil
}

func (r *Relation) DeleteFriend(ctx context.Context, friendUserID string) error {
	if err := r.deleteFriend(ctx, friendUserID); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}

func (r *Relation) GetFriendList(ctx context.Context, filterBlack bool) ([]*model_struct.LocalFriend, error) {
	localFriendList, err := r.db.GetAllFriendList(ctx)
	if err != nil {
		return nil, err
	}
	if len(localFriendList) == 0 || !filterBlack {
		return localFriendList, nil
	}
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	if len(localBlackList) == 0 {
		return localFriendList, nil
	}
	blackSet := datautil.SliceSetAny(localBlackList, func(e *model_struct.LocalBlack) string {
		return e.BlockUserID
	})
	var res []*model_struct.LocalFriend
	for _, friend := range localFriendList {
		if _, ok := blackSet[friend.FriendUserID]; !ok {
			res = append(res, friend)
		}
	}
	return res, nil
}

func (r *Relation) GetFriendListPage(ctx context.Context, offset, count int32, filterBlack bool) ([]*model_struct.LocalFriend, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		r.db,
		r.friendListTableName(),
		r.loginUserID,
		func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		func(ctx context.Context, values []*model_struct.LocalFriend) error {
			return r.db.BatchInsertFriend(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, bool, error) {
			localFriendList, err := r.db.GetFriendInfoList(ctx, userIDs)
			return localFriendList, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			serverFriend, err := r.getDesignatedFriends(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerFriendToLocalFriend, serverFriend), nil
		},
	)
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	if (!filterBlack) || len(localBlackList) == 0 {
		return dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
	}
	localFriendList, err := dataFetcher.FetchWithPagination(ctx, int(offset), int(count*2))
	if err != nil {
		return nil, err
	}
	blackUserIDs := datautil.SliceSetAny(localBlackList, func(e *model_struct.LocalBlack) string {
		return e.BlockUserID
	})
	res := localFriendList[:0]
	for _, friend := range localFriendList {
		if _, ok := blackUserIDs[friend.FriendUserID]; !ok {
			res = append(res, friend)
		}
		if len(res) == int(count) {
			break
		}
	}
	return res, nil
}

func (r *Relation) SearchFriends(ctx context.Context, param *sdk.SearchFriendsParam) ([]*sdk.SearchFriendItem, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchNickname && !param.IsSearchUserID && !param.IsSearchRemark) {
		return nil, sdkerrs.ErrArgs.WrapMsg("keyword is null or search field all false")
	}
	localFriendList, err := r.db.SearchFriendList(ctx, param.KeywordList[0], param.IsSearchUserID, param.IsSearchNickname, param.IsSearchRemark)
	if err != nil {
		return nil, err
	}
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]struct{})
	for _, black := range localBlackList {
		m[black.BlockUserID] = struct{}{}
	}
	res := make([]*sdk.SearchFriendItem, 0, len(localFriendList))
	for i, localFriend := range localFriendList {
		var relationship int
		if _, ok := m[localFriend.FriendUserID]; ok {
			relationship = constant.BlackRelationship
		} else {
			relationship = constant.FriendRelationship
		}
		res = append(res, &sdk.SearchFriendItem{
			LocalFriend:  *localFriendList[i],
			Relationship: relationship,
		})
	}
	return res, nil
}

func (r *Relation) AddBlack(ctx context.Context, blackUserID string, ex string) error {
	if err := r.addBlack(ctx, &relation.AddBlackReq{BlackUserID: blackUserID, Ex: ex}); err != nil {
		return err
	}
	return r.SyncAllBlackList(ctx)
}

func (r *Relation) RemoveBlack(ctx context.Context, blackUserID string) error {
	if err := r.removeBlack(ctx, blackUserID); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.SyncAllBlackList(ctx)
}

func (r *Relation) GetBlackList(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	return r.db.GetBlackListDB(ctx)
}

func (r *Relation) UpdateFriends(ctx context.Context, req *relation.UpdateFriendsReq) error {
	req.OwnerUserID = r.loginUserID
	if err := r.updateFriends(ctx, req); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}

func (r *Relation) GetFriendApplicationUnhandledCount(ctx context.Context, req *sdk.GetSelfUnhandledApplyCountReq) (int32, error) {
	return r.getSelfUnhandledApplyCount(ctx, req.Time)
}
