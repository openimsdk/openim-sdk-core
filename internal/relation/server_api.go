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

func (r *Relation) getSelfFriendApplicationList(ctx context.Context, pageNumber, showNumber int32) ([]*sdkws.FriendRequest, error) {
	req := &relation.GetPaginationFriendsApplyFromReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: showNumber}}
	if showNumber <= 0 {
		return api.Page(ctx, req, api.GetSelfFriendApplicationList.Invoke, (*relation.GetPaginationFriendsApplyFromResp).GetFriendRequests)
	}
	resp, err := api.GetSelfFriendApplicationList.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetFriendRequests(), nil
}

func (r *Relation) getRecvFriendApplicationList(ctx context.Context, handleResults []int32, pageNumber, showNumber int32) ([]*sdkws.FriendRequest, error) {
	req := &relation.GetPaginationFriendsApplyToReq{UserID: r.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: pageNumber, ShowNumber: showNumber},
		HandleResults: handleResults}
	if showNumber <= 0 {
		return api.Page(ctx, req, api.GetRecvFriendApplicationList.Invoke, (*relation.GetPaginationFriendsApplyToResp).GetFriendRequests)
	}
	resp, err := api.GetRecvFriendApplicationList.Invoke(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetFriendRequests(), nil
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

func (r *Relation) getSelfUnhandledApplyCount(ctx context.Context, time int64) (int32, error) {
	resp, err := api.GetSelfUnhandledApplyCount.Invoke(ctx, &relation.GetSelfUnhandledApplyCountReq{UserID: r.loginUserID, Time: time})
	if err != nil {
		return 0, err
	}
	return int32(resp.GetCount()), nil
}
