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
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdkerrs"

	"github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
)

func (f *Friend) SyncBothFriendRequest(ctx context.Context, fromUserID, toUserID string) error {
	var resp friend.GetDesignatedFriendsApplyResp
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsApplyRouter, &friend.GetDesignatedFriendsApplyReq{FromUserID: fromUserID, ToUserID: toUserID}, &resp); err != nil {
		return nil
	}
	localData, err := f.db.GetBothFriendReq(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}
	if toUserID == f.loginUserID {
		return f.requestRecvSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
	} else if fromUserID == f.loginUserID {
		return f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
	}
	return nil
}

// send
func (f *Friend) SyncAllSelfFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyFromReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friend.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestSendSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

// recv
func (f *Friend) SyncAllFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friend.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest { return resp.FriendRequests }
	requests, err := util.GetPageAll(ctx, constant.GetFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestRecvSyncer.Sync(ctx, util.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

func (f *Friend) SyncAllFriendList(ctx context.Context) error {
	req := &friend.GetPaginationFriendsReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friend.GetPaginationFriendsResp) []*sdkws.FriendInfo { return resp.FriendsInfo }
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil)
}

func (f *Friend) deleteFriend(ctx context.Context, friendUserID string) error {
	friends, err := f.db.GetFriendInfoList(ctx, []string{friendUserID})
	if err != nil {
		return err
	}
	if len(friends) == 0 {
		return sdkerrs.ErrUserIDNotFound.Wrap("friendUserID not found")
	}
	if err := f.db.DeleteFriendDB(ctx, friendUserID); err != nil {
		return err
	}
	f.friendListener.OnFriendDeleted(*friends[0])
	return nil
}

func (f *Friend) SyncFriends(ctx context.Context, friendIDs []string) error {
	var resp friend.GetDesignatedFriendsResp
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &friend.GetDesignatedFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
		return err
	}
	localData, err := f.db.GetFriendInfoList(ctx, friendIDs)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync friend", "data from server", resp.FriendsInfo, "data from local", localData)
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, resp.FriendsInfo), localData, nil)
}

func (f *Friend) SyncAllBlackList(ctx context.Context) error {
	req := &friend.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *friend.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return f.blockSyncer.Sync(ctx, util.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}
