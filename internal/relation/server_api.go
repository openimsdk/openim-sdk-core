package relation

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
)

func (r *Relation) getDesignatedFriendsApply(ctx context.Context, req *relation.GetDesignatedFriendsApplyReq) ([]*sdkws.FriendRequest, error) {
	return api.Field(ctx, api.GetDesignatedFriendsApply.Invoke, req, (*relation.GetDesignatedFriendsApplyResp).GetFriendRequests)
}

func (r *Relation) getSelfFriendApplicationList(ctx context.Context, req *relation.GetPaginationFriendsApplyFromReq) ([]*sdkws.FriendRequest, error) {
	return api.Page(ctx, req, api.GetSelfFriendApplicationList.Invoke, (*relation.GetPaginationFriendsApplyFromResp).GetFriendRequests)
}

func (r *Relation) getFriendApplicationList(ctx context.Context, req *relation.GetPaginationFriendsApplyToReq) ([]*sdkws.FriendRequest, error) {
	return api.Page(ctx, req, api.GetFriendApplicationList.Invoke, (*relation.GetPaginationFriendsApplyToResp).GetFriendRequests)
}

func (r *Relation) getBlackList(ctx context.Context) ([]*sdkws.BlackInfo, error) {
	req := &relation.GetPaginationBlacksReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{}}
	return api.Page(ctx, req, api.GetBlackList.Invoke, (*relation.GetPaginationBlacksResp).GetBlacks)
}

func (r *Relation) getDesignatedFriends(ctx context.Context, friendIDs []string) ([]*sdkws.FriendInfo, error) {
	req := &relation.GetDesignatedFriendsReq{OwnerUserID: r.loginUserID, FriendUserIDs: friendIDs}
	return api.Field(ctx, api.GetDesignatedFriends.Invoke, req, (*relation.GetDesignatedFriendsResp).GetFriendsInfo)
}

func (r *Relation) getIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	return api.GetIncrementalFriends.Invoke(ctx, req)
}

func (r *Relation) getFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	return api.GetFullFriendUserIDs.Invoke(ctx, req)
}
