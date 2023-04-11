package friend

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/internal/util"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
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
	ls := make([]*model_struct.LocalFriendRequest, 0, len(requests))
	for _, info := range requests {
		ls = append(ls, &model_struct.LocalFriendRequest{
			FromUserID:    info.FromUserID,
			FromNickname:  info.FromNickname,
			FromFaceURL:   info.FromFaceURL,
			FromGender:    info.FromGender,
			ToUserID:      info.ToUserID,
			ToNickname:    info.ToNickname,
			ToFaceURL:     info.ToFaceURL,
			ToGender:      info.ToGender,
			HandleResult:  info.HandleResult,
			ReqMsg:        info.ReqMsg,
			CreateTime:    info.CreateTime,
			HandlerUserID: info.HandlerUserID,
			HandleMsg:     info.HandleMsg,
			HandleTime:    info.HandleTime,
			Ex:            info.Ex,
			//AttachedInfo:  info.AttachedInfo,
		})
	}
	return syncer.New(nil).AddGlobal(map[string]any{"friend_user_id": f.loginUserID}, ls).Start()
}

// recv
func (f *Friend) SyncFriendApplication(ctx context.Context) error {
	req := &friend.GetPaginationFriendsApplyToReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
	fn := func(resp *friend.GetPaginationFriendsApplyToResp) []*sdkws.FriendRequest { return resp.FriendRequests }
	requests, err := util.GetPageAll(ctx, constant.GetFriendApplicationListRouter, req, fn)
	if err != nil {
		return err
	}
	ls := make([]*model_struct.LocalFriendRequest, 0, len(requests))
	for _, info := range requests {
		ls = append(ls, &model_struct.LocalFriendRequest{
			FromUserID:    info.FromUserID,
			FromNickname:  info.FromNickname,
			FromFaceURL:   info.FromFaceURL,
			FromGender:    info.FromGender,
			ToUserID:      info.ToUserID,
			ToNickname:    info.ToNickname,
			ToFaceURL:     info.ToFaceURL,
			ToGender:      info.ToGender,
			HandleResult:  info.HandleResult,
			ReqMsg:        info.ReqMsg,
			CreateTime:    info.CreateTime,
			HandlerUserID: info.HandlerUserID,
			HandleMsg:     info.HandleMsg,
			HandleTime:    info.HandleTime,
			Ex:            info.Ex,
			//AttachedInfo:  info.AttachedInfo,
		})
	}
	return syncer.New(nil).AddGlobal(map[string]any{"owner_user_id": f.loginUserID}, ls).Start()
}

func (f *Friend) SyncFriendList(ctx context.Context) error {
	req := &friend.GetPaginationFriendsReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
	fn := func(resp *friend.GetPaginationFriendsResp) []*sdkws.FriendInfo { return resp.FriendsInfo }
	friends, err := util.GetPageAll(ctx, constant.GetFriendListRouter, req, fn)
	if err != nil {
		return err
	}
	ls := make([]*model_struct.LocalFriend, 0, len(friends))
	for _, info := range friends {
		ls = append(ls, &model_struct.LocalFriend{
			OwnerUserID:    info.OwnerUserID,
			FriendUserID:   info.FriendUser.UserID,
			Remark:         info.Remark,
			CreateTime:     info.CreateTime,
			AddSource:      info.AddSource,
			OperatorUserID: info.OperatorUserID,
			Nickname:       info.FriendUser.Nickname,
			FaceURL:        info.FriendUser.FaceURL,
			Ex:             info.Ex,
			//AttachedInfo:   info.FriendUser.AttachedInfo,
		})
	}
	return syncer.New(nil).AddGlobal(map[string]any{"owner_user_id": f.loginUserID}, ls).Start()
}

func (f *Friend) SyncBlackList(ctx context.Context) error {
	req := &friend.GetPaginationBlacksReq{UserID: f.loginUserID, Pagination: &sdkws.RequestPagination{PageNumber: 0, ShowNumber: 30}}
	fn := func(resp *friend.GetPaginationBlacksResp) []*sdkws.BlackInfo { return resp.Blacks }
	blacks, err := util.GetPageAll(ctx, constant.GetBlackListRouter, req, fn)
	if err != nil {
		return err
	}
	ls := make([]*model_struct.LocalBlack, 0, len(blacks))
	for _, info := range blacks {
		ls = append(ls, &model_struct.LocalBlack{
			OwnerUserID:    info.OwnerUserID,
			BlockUserID:    info.BlackUserInfo.UserID,
			CreateTime:     info.CreateTime,
			AddSource:      info.AddSource,
			OperatorUserID: info.OperatorUserID,
			Nickname:       info.BlackUserInfo.Nickname,
			FaceURL:        info.BlackUserInfo.FaceURL,
			Ex:             info.Ex,
			//AttachedInfo:   info.FriendUser.AttachedInfo,
		})
	}
	return syncer.New(nil).AddGlobal(map[string]any{"owner_user_id": f.loginUserID}, ls).Start()
}
