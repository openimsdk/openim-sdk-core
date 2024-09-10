package user

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	authPb "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
)

func (u *User) getUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := api.GetUsersInfo.Invoke(ctx, &user.GetDesignateUsersReq{UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.UsersInfo, nil
}

func (u *User) updateUserInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	userInfo.UserID = u.loginUserID
	_, err := api.UpdateUserInfo.Invoke(ctx, &user.UpdateUserInfoReq{UserInfo: userInfo})
	return err
}

func (u *User) updateUserInfoV2(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	userInfo.UserID = u.loginUserID
	_, err := api.UpdateUserInfoEx.Invoke(ctx, &user.UpdateUserInfoExReq{UserInfo: userInfo})
	return err
}

func (u *User) processUserCommandAdd(ctx context.Context, req *user.ProcessUserCommandAddReq) error {
	_, err := api.ProcessUserCommandAdd.Invoke(ctx, req)
	return err
}

func (u *User) processUserCommandDelete(ctx context.Context, req *user.ProcessUserCommandDeleteReq) error {
	_, err := api.ProcessUserCommandDelete.Invoke(ctx, req)
	return err
}

func (u *User) processUserCommandUpdate(ctx context.Context, req *user.ProcessUserCommandUpdateReq) error {
	_, err := api.ProcessUserCommandUpdate.Invoke(ctx, req)
	return err
}

func (u *User) parseToken(ctx context.Context) (*authPb.ParseTokenResp, error) {
	return api.ParseToken.Invoke(ctx, &authPb.ParseTokenReq{})
}
