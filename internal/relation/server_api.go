package relation

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
)

func (r *Relation) getDesignatedFriendsApply(ctx context.Context, req *relation.GetDesignatedFriendsApplyReq) ([]*sdkws.FriendRequest, error) {
	return api.ExtractField(ctx, api.GetDesignatedFriendsApply.Invoke, req, (*relation.GetDesignatedFriendsApplyResp).GetFriendRequests)
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
	return api.ExtractField(ctx, api.GetDesignatedFriends.Invoke, req, (*relation.GetDesignatedFriendsResp).GetFriendsInfo)
}

func (r *Relation) getIncrementalFriends(ctx context.Context, req *relation.GetIncrementalFriendsReq) (*relation.GetIncrementalFriendsResp, error) {
	return api.GetIncrementalFriends.Invoke(ctx, req)
}

func (r *Relation) getFullFriendUserIDs(ctx context.Context, req *relation.GetFullFriendUserIDsReq) (*relation.GetFullFriendUserIDsResp, error) {
	return api.GetFullFriendUserIDs.Invoke(ctx, req)
}

func (r *Relation) addFriend(ctx context.Context, req *relation.ApplyToAddFriendReq) error {
	req.FromUserID = r.loginUserID
	return api.AddFriend.Execute(ctx, req)
}

func (r *Relation) deleteFriend(ctx context.Context, friendUserID string) error {
	req := &relation.DeleteFriendReq{OwnerUserID: r.loginUserID, FriendUserID: friendUserID}
	return api.DeleteFriend.Execute(ctx, req)
}

func (r *Relation) addFriendResponse(ctx context.Context, req *relation.RespondFriendApplyReq) error {
	req.ToUserID = r.loginUserID
	return api.AddFriendResponse.Execute(ctx, req)
}

func (r *Relation) updateFriends(ctx context.Context, req *relation.UpdateFriendsReq) error {
	req.OwnerUserID = r.loginUserID
	return api.UpdateFriends.Execute(ctx, req)
}

func (r *Relation) addBlack(ctx context.Context, req *relation.AddBlackReq) error {
	req.OwnerUserID = r.loginUserID
	return api.AddBlack.Execute(ctx, req)
}

func (r *Relation) removeBlack(ctx context.Context, userID string) error {
	return api.RemoveBlack.Execute(ctx, &relation.RemoveBlackReq{
		OwnerUserID: r.loginUserID,
		BlackUserID: userID,
	})
}
