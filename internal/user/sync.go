package user

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, []*model_struct.LocalUser{localUser}, nil)
}
