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
	"fmt"
	"github.com/openimsdk/tools/utils/datautil"
	"time"

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	friend "github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
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
		return f.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
	} else if fromUserID == f.loginUserID {
		return f.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
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
	return f.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
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
	return f.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

func (f *Friend) SyncAllFriendList(ctx context.Context) error {
	t := time.Now()
	defer func(start time.Time) {

		elapsed := time.Since(start).Milliseconds()
		log.ZDebug(ctx, "SyncAllFriendList fn call end", "cost time", fmt.Sprintf("%d ms", elapsed))

	}(t)
	return f.IncrSyncFriends(ctx)
	//req := &friend.GetPaginationFriendsReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	//fn := func(resp *friend.GetPaginationFriendsResp) []*sdkws.FriendInfo { return resp.FriendsInfo }
	//friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, fn)
	//if err != nil {
	//	return err
	//}
	//localData, err := f.db.GetAllFriendList(ctx)
	//if err != nil {
	//	return err
	//}
	//log.ZDebug(ctx, "sync friend", "data from server", friends, "data from local", localData)
	//return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil)
}

func (f *Friend) deleteFriend(ctx context.Context, friendUserID string) error {
	return f.IncrSyncFriends(ctx)
	//friends, err := f.db.GetFriendInfoList(ctx, []string{friendUserID})
	//if err != nil {
	//	return err
	//}
	//if len(friends) == 0 {
	//	return sdkerrs.ErrUserIDNotFound.WrapMsg("friendUserID not found")
	//}
	//if err := f.db.DeleteFriendDB(ctx, friendUserID); err != nil {
	//	return err
	//}
	//f.friendListener.OnFriendDeleted(*friends[0])
	//return nil
}

func (f *Friend) SyncFriends(ctx context.Context, friendIDs []string) error {
	return f.IncrSyncFriends(ctx)
	//var resp friend.GetDesignatedFriendsResp
	//if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &friend.GetDesignatedFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
	//	return err
	//}
	//localData, err := f.db.GetFriendInfoList(ctx, friendIDs)
	//if err != nil {
	//	return err
	//}
	//log.ZDebug(ctx, "sync friend", "data from server", resp.FriendsInfo, "data from local", localData)
	//return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, resp.FriendsInfo), localData, nil)
}

//func (f *Friend) SyncFriendPart(ctx context.Context) error {
//	hashResp, err := util.CallApi[friend.GetFriendHashResp](ctx, constant.GetFriendHash, &friend.GetFriendHashReq{UserID: f.loginUserID})
//	if err != nil {
//		return err
//	}
//	friends, err := f.db.GetAllFriendList(ctx)
//	if err != nil {
//		return err
//	}
//	hashCode := f.CalculateHash(friends)
//	log.ZDebug(ctx, "SyncFriendPart", "serverHash", hashResp.Hash, "serverTotal", hashResp.Total, "localHash", hashCode, "localTotal", len(friends))
//	if hashCode == hashResp.Hash {
//		return nil
//	}
//	req := &friend.GetPaginationFriendsReq{
//		UserID:     f.loginUserID,
//		Pagination: &sdkws.RequestPagination{PageNumber: pconstant.FirstPageNumber, ShowNumber: pconstant.MaxSyncPullNumber},
//	}
//	resp, err := util.CallApi[friend.GetPaginationFriendsResp](ctx, constant.GetFriendListRouter, req)
//	if err != nil {
//		return err
//	}
//	serverFriends := util.Batch(ServerFriendToLocalFriend, resp.FriendsInfo)
//	return f.friendSyncer.Sync(ctx, serverFriends, friends, nil)
//}

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
	return f.blockSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}

func (f *Friend) GetDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	resp := &friend.GetDesignatedFriendsResp{}
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &friend.GetDesignatedFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
		return nil, err
	}
	return resp.FriendsInfo, nil
}
