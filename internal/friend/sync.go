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
	"time"

	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (f *Friend) SyncBothFriendRequest(ctx context.Context, fromUserID, toUserID string) error {
	var resp relation.GetDesignatedFriendsApplyResp
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsApplyRouter, &relation.GetDesignatedFriendsApplyReq{FromUserID: fromUserID, ToUserID: toUserID}, &resp); err != nil {
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
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
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

func (f *Friend) SyncAllSelfFriendApplicationWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
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
	return f.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

// recv
func (f *Friend) SyncAllFriendApplication(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
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
func (f *Friend) SyncAllFriendApplicationWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return f.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

func (f *Friend) SyncAllFriendList(ctx context.Context) error {
	t := time.Now()
	defer func(start time.Time) {

		elapsed := time.Since(start).Milliseconds()
		log.ZDebug(ctx, "SyncAllFriendList fn call end", "cost time", fmt.Sprintf("%d ms", elapsed))

	}(t)
	return f.IncrSyncFriends(ctx)
}

func (f *Friend) SyncAllBlackList(ctx context.Context) error {
	req := &relation.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
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

func (f *Friend) SyncAllBlackListWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
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
	return f.blockSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil, false, true)
}

func (f *Friend) GetDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	resp := &relation.GetDesignatedFriendsResp{}
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &relation.GetDesignatedFriendsReq{OwnerUserID: f.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
		return nil, err
	}
	return resp.FriendsInfo, nil
}
