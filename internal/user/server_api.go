package user

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	authPb "github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/sdkws"
	userPb "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/utils/datautil"
)

// ParseTokenFromSvr parses a token from the server.
func (u *User) ParseTokenFromSvr(ctx context.Context) (int64, error) {
	resp, err := api.ParseToken.Invoke(ctx, &authPb.ParseTokenReq{})
	if err != nil {
		return 0, err
	}
	return resp.ExpireTimeSeconds, nil
}

// GetServerUserInfo retrieves user information from the server.
func (u *User) GetServerUserInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := api.GetUsersInfo.Invoke(ctx, &userPb.GetDesignateUsersReq{UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.UsersInfo, nil
}

// updateSelfUserInfo updates the user's information.
func (u *User) updateSelfUserInfo(ctx context.Context, userInfo *sdkws.UserInfo) error {
	_, err := api.UpdateSelfUserInfo.Invoke(ctx, &userPb.UpdateUserInfoReq{UserInfo: userInfo})
	return err
}

// updateSelfUserInfoEx updates the user's information with Ex field.
func (u *User) updateSelfUserInfoEx(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	_, err := api.UpdateSelfUserInfoEx.Invoke(ctx, &userPb.UpdateUserInfoExReq{UserInfo: userInfo})
	return err
}

func (u *User) processUserCommandAdd(ctx context.Context, userCommand *userPb.ProcessUserCommandAddReq) error {
	_, err := api.ProcessUserCommandAdd.Invoke(ctx, &userPb.ProcessUserCommandAddReq{
		UserID: u.loginUserID,
		Type:   userCommand.Type,
		Uuid:   userCommand.Uuid,
		Value:  userCommand.Value,
	})
	return err
}

// processUserCommandDelete is a private method to handle the actual delete command API call.
func (u *User) processUserCommandDelete(ctx context.Context, userCommand *userPb.ProcessUserCommandDeleteReq) error {
	_, err := api.ProcessUserCommandDelete.Invoke(ctx, &userPb.ProcessUserCommandDeleteReq{
		UserID: u.loginUserID,
		Type:   userCommand.Type,
		Uuid:   userCommand.Uuid,
	})
	return err
}

// processUserCommandUpdate is a private method to handle the actual update command API call.
func (u *User) processUserCommandUpdate(ctx context.Context, userCommand *userPb.ProcessUserCommandUpdateReq) error {
	_, err := api.ProcessUserCommandUpdate.Invoke(ctx, &userPb.ProcessUserCommandUpdateReq{
		UserID: u.loginUserID,
		Type:   userCommand.Type,
		Uuid:   userCommand.Uuid,
		Value:  userCommand.Value,
	})
	return err
}

// GetUsersInfoFromSvr retrieves user information from the server.
func (u *User) GetUsersInfoFromSvr(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	resp, err := api.GetUsersInfo.Invoke(ctx, &userPb.GetDesignateUsersReq{UserIDs: userIDs})
	if err != nil {
		return nil, sdkerrs.WrapMsg(err, "GetUsersInfoFromSvr failed")
	}
	return datautil.Batch(ServerUserToLocalUser, resp.UsersInfo), nil
}

// processUserCommandGetAll is a private method that requests all user commands from the server.
func (u *User) processUserCommandGetAll(ctx context.Context) ([]*userPb.AllCommandInfoResp, error) {
	resp, err := api.ProcessUserCommandGetAll.Invoke(ctx, &userPb.ProcessUserCommandGetAllReq{
		UserID: u.loginUserID,
	})
	if err != nil {
		return nil, err
	}
	return resp.CommandResp, nil
}
