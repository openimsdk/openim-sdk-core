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

func (f *Friend) GetSpecifiedFriendsInfo(ctx context.Context, friendUserIDList []string) ([]*server_api_params.FullUserInfo, error) {
	datafetcher := datafetcher.NewDataFetcher(
		f.db,
		f.friendListTableName(),
		f.loginUserID,
		func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		func(ctx context.Context, values []*model_struct.LocalFriend) error {
			return f.db.BatchInsertFriend(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			return f.db.GetFriendInfoList(ctx, userIDs)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			serverFriend, err := f.GetDesignatedFriends(ctx, userIDs)
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
	blackList, err := f.db.GetBlackInfoList(ctx, friendUserIDList)
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

func (f *Friend) AddFriend(ctx context.Context, userIDReqMsg *friend.ApplyToAddFriendReq) error {
	if userIDReqMsg.FromUserID == "" {
		userIDReqMsg.FromUserID = f.loginUserID
	}
	if err := util.ApiPost(ctx, constant.AddFriendRouter, userIDReqMsg, nil); err != nil {
		return err
	}
	return f.SyncAllFriendApplication(ctx)
}

func (f *Friend) GetFriendApplicationListAsRecipient(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetRecvFriendApplication(ctx)
}

func (f *Friend) GetFriendApplicationListAsApplicant(ctx context.Context) ([]*model_struct.LocalFriendRequest, error) {
	return f.db.GetSendFriendApplication(ctx)
}

func (f *Friend) AcceptFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, HandleResult: constant.FriendResponseAgree, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RefuseFriendApplication(ctx context.Context, userIDHandleMsg *sdk.ProcessFriendApplicationParams) error {
	return f.RespondFriendApply(ctx, &friend.RespondFriendApplyReq{FromUserID: userIDHandleMsg.ToUserID, ToUserID: f.loginUserID, HandleResult: constant.FriendResponseRefuse, HandleMsg: userIDHandleMsg.HandleMsg})
}

func (f *Friend) RespondFriendApply(ctx context.Context, req *friend.RespondFriendApplyReq) error {
	if req.ToUserID == "" {
		req.ToUserID = f.loginUserID
	}
	if err := util.ApiPost(ctx, constant.AddFriendResponse, req, nil); err != nil {
		return err
	}
	if req.HandleResult == constant.FriendResponseAgree {
		_ = f.SyncFriends(ctx, []string{req.FromUserID})
	}
	_ = f.SyncAllFriendApplication(ctx)
	return nil
	// return f.SyncFriendApplication(ctx)
}

func (f *Friend) CheckFriend(ctx context.Context, friendUserIDList []string) ([]*server_api_params.UserIDResult, error) {
	friendList, err := f.db.GetFriendInfoList(ctx, friendUserIDList)
	if err != nil {
		return nil, err
	}
	blackList, err := f.db.GetBlackInfoList(ctx, friendUserIDList)
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
		for _, f := range friendList {
			if v == f.FriendUserID {
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

func (f *Friend) DeleteFriend(ctx context.Context, friendUserID string) error {
	if err := util.ApiPost(ctx, constant.DeleteFriendRouter, &friend.DeleteFriendReq{OwnerUserID: f.loginUserID, FriendUserID: friendUserID}, nil); err != nil {
		return err
	}
	return f.deleteFriend(ctx, friendUserID)
}

func (f *Friend) GetFriendList(ctx context.Context) ([]*server_api_params.FullUserInfo, error) {
	localFriendList, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
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

func (f *Friend) GetFriendListPage(ctx context.Context, offset, count int32) ([]*server_api_params.FullUserInfo, error) {
	dataFetcher := datafetcher.NewDataFetcher(
		f.db,
		f.friendListTableName(),
		f.loginUserID,
		func(localFriend *model_struct.LocalFriend) string {
			return localFriend.FriendUserID
		},
		func(ctx context.Context, values []*model_struct.LocalFriend) error {
			return f.db.BatchInsertFriend(ctx, values)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			return f.db.GetFriendInfoList(ctx, userIDs)
		},
		func(ctx context.Context, userIDs []string) ([]*model_struct.LocalFriend, error) {
			serverFriend, err := f.GetDesignatedFriends(ctx, userIDs)
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
	localBlackList, err := f.db.GetBlackListDB(ctx)
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

func (f *Friend) SearchFriends(ctx context.Context, param *sdk.SearchFriendsParam) ([]*sdk.SearchFriendItem, error) {
	if len(param.KeywordList) == 0 || (!param.IsSearchNickname && !param.IsSearchUserID && !param.IsSearchRemark) {
		return nil, sdkerrs.ErrArgs.WrapMsg("keyword is null or search field all false")
	}
	localFriendList, err := f.db.SearchFriendList(ctx, param.KeywordList[0], param.IsSearchUserID, param.IsSearchNickname, param.IsSearchRemark)
	if err != nil {
		return nil, err
	}
	localBlackList, err := f.db.GetBlackListDB(ctx)
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

func (f *Friend) SetFriendRemark(ctx context.Context, userIDRemark *sdk.SetFriendRemarkParams) error {
	if err := util.ApiPost(ctx, constant.SetFriendRemark, &friend.SetFriendRemarkReq{OwnerUserID: f.loginUserID, FriendUserID: userIDRemark.ToUserID, Remark: userIDRemark.Remark}, nil); err != nil {
		return err
	}
	return f.SyncFriends(ctx, []string{userIDRemark.ToUserID})
}

func (f *Friend) PinFriends(ctx context.Context, friends *sdk.SetFriendPinParams) error {
	if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friends.ToUserIDs, IsPinned: friends.IsPinned}, nil); err != nil {
		return err
	}
	return f.SyncFriends(ctx, friends.ToUserIDs)
}

func (f *Friend) AddBlack(ctx context.Context, blackUserID string, ex string) error {
	if err := util.ApiPost(ctx, constant.AddBlackRouter, &friend.AddBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID, Ex: ex}, nil); err != nil {
		return err
	}
	return f.SyncAllBlackList(ctx)
}

func (f *Friend) RemoveBlack(ctx context.Context, blackUserID string) error {
	if err := util.ApiPost(ctx, constant.RemoveBlackRouter, &friend.RemoveBlackReq{OwnerUserID: f.loginUserID, BlackUserID: blackUserID}, nil); err != nil {
		return err
	}
	return f.SyncAllBlackList(ctx)
}

func (f *Friend) GetBlackList(ctx context.Context) ([]*model_struct.LocalBlack, error) {
	return f.db.GetBlackListDB(ctx)
}

func (f *Friend) SetFriendsEx(ctx context.Context, friendIDs []string, ex string) error {
	if err := util.ApiPost(ctx, constant.UpdateFriends, &friend.UpdateFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs, Ex: &wrapperspb.StringValue{
		Value: ex,
	}}, nil); err != nil {
		return err
	}
	// Check if the specified ID is a friend
	friendResults, err := f.CheckFriend(ctx, friendIDs)
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

	updateErr := f.db.UpdateColumnsFriend(ctx, friendIDs, map[string]interface{}{"Ex": ex})
	if updateErr != nil {
		return errs.WrapMsg(updateErr, "Error updating friend information")
	}
	return nil
}
