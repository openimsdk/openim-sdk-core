package relation

import (
	"context"

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

func (r *Relation) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string) ([]*server_api_params.FullUserInfo, error) {
	datafetcher := datafetcher.NewDataFetcher(
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
			serverFriend, err := r.GetDesignatedFriends(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerFriendToLocalFriend, serverFriend), nil
		},
	)
	localFriendList, err := datafetcher.FetchMissingAndFillLocal(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}

	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "localFriendList", localFriendList)
	blackList, err := r.db.GetBlackInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	log.ZDebug(ctx, "GetDesignatedFriendsInfo", "blackList", blackList)
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range blackList {
		m[black.BlockUserID] = blackList[i]
	}
	res := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		res = append(res, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	return res, nil
}

func (r *Relation) AddFriend(ctx context.Context, req *relation.ApplyToAddFriendReq) error {
	if err := r.addFriend(ctx, req); err != nil {
		return err
	}
	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.SyncAllSelfFriendApplication(ctx)
}

func (r *Relation) GetFriendApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return r.db.GetRecvFriendApplication(ctx)
}

func (r *Relation) GetFriendApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return r.db.GetSendFriendApplication(ctx)
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
	_ = r.SyncAllFriendApplication(ctx)
	return nil
	// return r.SyncFriendApplication(ctx)
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
	if localFriendList == nil {
		localFriendList = []*model_struct.LocalFriend{}
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
	blackSet := make(map[string]struct{})
	for _, black := range localBlackList {
		blackSet[black.BlockUserID] = struct{}{}
	}
	res := localFriendList[:0]
	for i, friend := range localFriendList {
		if _, ok := blackSet[friend.FriendUserID]; !ok {
			res = append(res, localFriendList[i])
		}
	}
	return res, nil
}

func (r *Relation) GetFriendListPage(ctx context.Context, offset, count int, filterBlack bool) ([]*model_struct.LocalFriend, error) {
	friends, err := r.GetFriendList(ctx, filterBlack)
	if err != nil {
		return nil, err
	}
	if offset >= len(friends) {
		return friends[:0], nil
	}
	friends = friends[offset:]
	if len(friends) > count {
		return friends[:count], nil
	} else {
		return friends, nil
	}
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
	if err := r.updateFriends(ctx, req); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}
