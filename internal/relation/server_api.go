package relation

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
)

func (r *Relation) getDesignatedFriendsApply(ctx context.Context, req *relation.GetDesignatedFriendsApplyReq) ([]*sdkws.FriendRequest, error) {
	resp, err := api.GetDesignatedFriendsApply.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.FriendRequests, nil
}

func (r *Relation) getSelfFriendApplicationList(ctx context.Context, req *relation.GetPaginationFriendsApplyFromReq) ([]*sdkws.FriendRequest, error) {
	return util.PageNext(ctx, req, api.GetSelfFriendApplicationList.Invoke, func(resp *relation.GetPaginationFriendsApplyFromResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	})
}

func (r *Relation) getFriendApplicationList(ctx context.Context, req *relation.GetPaginationFriendsApplyToReq) ([]*sdkws.FriendRequest, error) {
	return util.PageNext(ctx, req, api.GetFriendApplicationList.Invoke, func(resp *relation.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest {
		return resp.FriendRequests
	})
}

func (r *Relation) getBlackList(ctx context.Context) ([]*sdkws.BlackInfo, error) {
	req := &relation.GetPaginationBlacksReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return util.PageNext(ctx, req, api.GetBlackList.Invoke, func(resp *relation.GetPaginationBlacksResp) []*sdkws.BlackInfo {
		return resp.Blacks
	})
}

func (r *Relation) getDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	req := &relation.GetDesignatedFriendsReq{OwnerUserID: r.loginUserID, FriendUserIDs: friendIDs}
	resp, err := api.GetDesignatedFriends.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.FriendsInfo, nil
}

func (r *Relation) getIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	return api.GetIncrementalFriends.Invoke(ctx, req)
}

func (r *Relation) getFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	return api.GetFullFriendUserIDs.Invoke(ctx, req)
}
