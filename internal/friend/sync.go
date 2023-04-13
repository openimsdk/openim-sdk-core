package friend

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
)

func (f *Friend) SyncSelfFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyFromReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
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
func (f *Friend) SyncFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
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

func (f *Friend) SyncFriendList(ctx context.Context) error {
	req := &friend.GetPaginationFriendsReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
	fn := func(resp *friend.GetPaginationFriendsResp) []*sdkws.FriendInfo { return resp.FriendsInfo }
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetAllFriendList(ctx)
	if err != nil {
		return err
	}
	return f.friendSyncer.Sync(ctx, util.Batch(ServerFriendToLocalFriend, friends), localData, nil)
}

func (f *Friend) SyncBlackList(ctx context.Context) error {
	req := &friend.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
	fn := func(resp *friend.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
	blacks, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	localData, err := f.db.GetBlackListDB(ctx)
	if err != nil {
		return err
	}
	return f.blockSyncer.Sync(ctx, util.Batch(ServerBlackToLocalBlack, blacks), localData, nil)
}
