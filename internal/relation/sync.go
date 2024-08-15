package relation

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

func (r *Relation) SyncBothFriendRequest(ctx context.Context, fromUserID, toUserID string) error {
	var resp relation.GetDesignatedFriendsApplyResp
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsApplyRouter, &relation.GetDesignatedFriendsApplyReq{FromUserID: fromUserID, ToUserID: toUserID}, &resp); err != nil {
		return nil
	}
	localData, err := r.db.GetBothFriendReq(ctx, fromUserID, toUserID)
	if err != nil {
		return err
	}
	if toUserID == r.loginUserID {
		return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
	} else if fromUserID == r.loginUserID {
		return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, resp.FriendRequests), localData, nil)
	}
	return nil
}

// send
func (r *Relation) SyncAllSelfFriendApplication(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := r.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}

func (r *Relation) SyncAllSelfFriendApplicationWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetSelfFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := r.db.GetSendFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestSendSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

// recv
func (r *Relation) SyncAllFriendApplication(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyToReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := r.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil)
}
func (r *Relation) SyncAllFriendApplicationWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationFriendsApplyToReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	}
	requests, err := util.GetPageAll(ctx, constant.GetFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := r.db.GetRecvFriendApplication(ctx)
	if err != nil {
		return err
	}
	return r.requestRecvSyncer.Sync(ctx, datautil.Batch(ServerFriendRequestToLocalFriendRequest, requests), localData, nil, false, true)
}

func (r *Relation) SyncAllFriendList(ctx context.Context) error {
	t := time.Now()
	defer func(start time.Time) {

		elapsed := time.Since(start).Milliseconds()
		log.ZDebug(ctx, "SyncAllFriendList fn call end", "cost time", fmt.Sprintf("%d ms", elapsed))

	}(t)
	return r.IncrSyncFriends(ctx)
}

func (r *Relation) SyncAllBlackList(ctx context.Context) error {
	req := &relation.GetPaginationBlacksReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return r.blockSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil)
}

func (r *Relation) SyncAllBlackListWithoutNotice(ctx context.Context) error {
	req := &relation.GetPaginationBlacksReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	fn := func(resp *relation.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
	serverData, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from server", "data", serverData)
	localData, err := r.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "black from local", "data", localData)
	return r.blockSyncer.Sync(ctx, datautil.Batch(ServerBlackToLocalBlack, serverData), localData, nil, false, true)
}

func (r *Relation) GetDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	resp := &relation.GetDesignatedFriendsResp{}
	if err := util.ApiPost(ctx, constant.GetDesignatedFriendsRouter, &relation.GetDesignatedFriendsReq{OwnerUserID: r.loginUserID, FriendUserIDs: friendIDs}, &resp); err != nil {
		return nil, err
	}
	return resp.FriendsInfo, nil
}
