package user

import (
	"context"
	"open_im_sdk/pkg/db/model_struct"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"gorm.io/gorm"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromSvr(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && errs.Unwrap(err) != gorm.ErrRecordNotFound {
		return err
	}
	var remoteUsers []*model_struct.LocalUser
	if err == nil {
		remoteUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, remoteUsers, nil)
}
