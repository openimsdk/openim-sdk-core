package user

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/api"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
)

func (u *User) getUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	req := &user.GetDesignateUsersReq{UserIDs: userIDs}
	return api.ExtractField(ctx, api.GetUsersInfo.Invoke, req, (*user.GetDesignateUsersResp).GetUsersInfo)
}

func (u *User) updateUserInfo(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	userInfo.UserID = u.loginUserID
	return api.UpdateUserInfoEx.Execute(ctx, &user.UpdateUserInfoExReq{UserInfo: userInfo})
}
