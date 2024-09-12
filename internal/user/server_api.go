package user

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
)

func (u *User) getUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	req := &user.GetDesignateUsersReq{UserIDs: userIDs}
	return api.Field(ctx, api.GetUsersInfo.Invoke, req, (*user.GetDesignateUsersResp).GetUsersInfo)
}

func (u *User) updateUserInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	userInfo.UserID = u.loginUserID
	return api.UpdateUserInfo.Result(ctx, &user.UpdateUserInfoReq{UserInfo: userInfo})
}

func (u *User) updateUserInfoV2(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	userInfo.UserID = u.loginUserID
	return api.UpdateUserInfoEx.Result(ctx, &user.UpdateUserInfoExReq{UserInfo: userInfo})
}

func (u *User) processUserCommandAdd(ctx context.Context, req *user.ProcessUserCommandAddReq) error {
	return api.ProcessUserCommandAdd.Result(ctx, req)
}

func (u *User) processUserCommandDelete(ctx context.Context, req *user.ProcessUserCommandDeleteReq) error {
	return api.ProcessUserCommandDelete.Result(ctx, req)
}

func (u *User) processUserCommandUpdate(ctx context.Context, req *user.ProcessUserCommandUpdateReq) error {
	return api.ProcessUserCommandUpdate.Result(ctx, req)
}

func (u *User) parseToken(ctx context.Context) (*auth.ParseTokenResp, error) {
	return api.ParseToken.Invoke(ctx, &auth.ParseTokenReq{})
}

func (u *User) setGlobalRecvMessageOpt(ctx context.Context, opt int32) error {
	return api.SetGlobalRecvMessageOpt.Result(ctx, &user.SetGlobalRecvMessageOptReq{UserID: u.loginUserID, GlobalRecvMsgOpt: opt})
}
