package relation

import (
	"context"

	friend "github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
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

func (r *Relation) AddFriend(ctx context.Context, userIDReqMsg *friend.ApplyToAddFriendReq) error {
	if userIDReqMsg.FromUserID == "" {
		userIDReqMsg.FromUserID = r.loginUserID
	}
	if err := util.ApiPost(ctx, constant.AddFriendRouter, userIDReqMsg, nil); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.SyncAllFriendApplication(ctx)
}

func (r *Relation) GetFriendApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return r.db.GetRecvFriendApplication(ctx)
}

func (r *Relation) GetFriendApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return r.db.GetSendFriendApplication(ctx)
}

func (r *Relation) AcceptFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return r.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: r.loginUserID, HandleResult: constant.FriendResponseAgree, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (r *Relation) RefuseFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return r.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: r.loginUserID, HandleResult: constant.FriendResponseRefuse, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (r *Relation) RespondFriendApply(ctx context.Context, req *friend.RespondFriendApplyReq) error {
	if req.ToUserID == "" {
		req.ToUserID = r.loginUserID
	}
	if err := util.ApiPost(ctx, constant.AddFriendResponse, req, nil); err != nil {
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
	if err := util.ApiPost(ctx, constant.DeleteFriendRouter, &friend.DeleteFriendReq{OwnerUserID: r.loginUserID, FriendUserID: friendUserID}, nil); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}

// Full GetFriendList
func (r *Relation) GetFriendList(ctx context.Context) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := r.db.GetAllFriendList(ctx)
	if err != nil {
		return nil, err
	}
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlockUserID] = localBlackList[i]
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

func (r *Relation) GetFriendListPage(ctx context.Context, offset, count int32) ([]*server_api_params.FullUserInfo, error) {
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
			serverFriend, err := r.GetDesignatedFriends(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerFriendToLocalFriend, serverFriend), nil
		},
	)

	localFriendList, err := dataFetcher.FetchWithPagination(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}

	// don't need extra handle. only full pull.
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlockUserID] = localBlackList[i]
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

func (r *Relation) GetFriendListPageV2(ctx context.Context, offset, count int32) (*GetFriendInfoListV2, error) {
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
			localFriendList, err := r.db.GetFriendInfoList(ctx, userIDs)
			return localFriendList, true, err
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			serverFriend, err := r.GetDesignatedFriends(ctx, userIDs)
			if err != nil {
				return nil, err
			}
			return datautil.Batch(ServerFriendToLocalFriend, serverFriend), nil
		},
	)

	localFriendList, isEnd, err := datafetcher.FetchWithPaginationV2(ctx, int(offset), int(count))
	if err != nil {
		return nil, err
	}

	// don't need extra handle. only full pull.
	localBlackList, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*model_struct.LocalBlack)
	for i, black := range localBlackList {
		m[black.BlockUserID] = localBlackList[i]
	}
	fullUserInfo := make([]*server_api_params.FullUserInfo, 0, len(localFriendList))
	for _, localFriend := range localFriendList {
		fullUserInfo = append(fullUserInfo, &server_api_params.FullUserInfo{
			PublicInfo: nil,
			FriendInfo: localFriend,
			BlackInfo:  m[localFriend.FriendUserID],
		})
	}
	response := &GetFriendInfoListV2{
		FullUserInfoList: fullUserInfo,
		IsEnd:            isEnd,
	}
	return response, nil
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

func (r *Relation) SetFriendRemark(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	if err := util.ApiPost(ctx, constant.SetFriendRemark, &friend.SetFriendRemarkReq{OwnerUserID: r.loginUserID, FriendUserID: userIDRemark.ToUserID, Remark: userIDRemark.Remark}, nil); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}

func (r *Relation) PinFriends(ctx context.Context, friends *sdk.SetFriendPinParams) error {
	if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: r.loginUserID, FriendUserIDs: friends.ToUserIDs, IsPinned: friends.IsPinned}, nil); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.IncrSyncFriends(ctx)
}

func (r *Relation) AddBlack(ctx context.Context, blackUserID string, ex string) error {
	if err := util.ApiPost(ctx, constant.AddBlackRouter, &friend.AddBlackReq{OwnerUserID: r.loginUserID, BlackUserID: blackUserID, Ex: ex}, nil); err != nil {
		return err
	}
	return r.SyncAllBlackList(ctx)
}

func (r *Relation) RemoveBlack(ctx context.Context, blackUserID string) error {
	if err := util.ApiPost(ctx, constant.RemoveBlackRouter, &friend.RemoveBlackReq{OwnerUserID: r.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
		return err
	}

	r.relationSyncMutex.Lock()
	defer r.relationSyncMutex.Unlock()

	return r.SyncAllBlackList(ctx)
}

func (r *Relation) GetBlackList(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	return r.db.GetBlackListDB(ctx)
}

func (r *Relation) SetFriendsEx(ctx context.Context, friendIDs []string, ex string) error {
	if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: r.loginUserID, FriendUserIDs: friendIDs, Ex: &wrapperspb.StringValue{
		Value: ex,
	}}, nil); err != nil {
		return err
	}
	// Check if the specified ID is a friend
	friendResults, err := r.CheckFriend(ctx, friendIDs)
	if err != nil {
		return errs.WrapMsg(err, "Error checking friend status")
	}

	// Determine if friendID is indeed a friend
	// Iterate over each friendID
	for _, friendID := range friendIDs {
		isFriend := false

		// Check if this friendID is in the friendResults
		for _, result := range friendResults {
			if result.UserID == friendID && result.Result == 1 { // Assuming result 1 means they are friends
				isFriend = true
				break
			}
		}

		// If this friendID is not a friend, return an error
		if !isFriend {
			return errs.ErrRecordNotFound.WrapMsg("Not friend")
		}
	}

	// If the code reaches here, all friendIDs are confirmed as friends
	// Update friend information if they are friends

	updateErr := r.db.UpdateColumnsFriend(ctx, friendIDs, map[string]interface{}{"Ex": ex})
	if updateErr != nil {
		return errs.WrapMsg(updateErr, "Error updating friend information")
	}
	return nil
}
