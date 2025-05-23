package user

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	userPb "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

func (u *User) SyncLoginUserInfo(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromServer(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && (!errs.ErrRecordNotFound.Is(errs.Unwrap(err))) {
		return err
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil)
}

func (u *User) SyncLoginUserInfoWithoutNotice(ctx context.Context) error {
	remoteUser, err := u.GetSingleUserFromServer(ctx, u.loginUserID)
	if err != nil {
		return err
	}
	localUser, err := u.GetLoginUser(ctx, u.loginUserID)
	if err != nil && (!errs.ErrRecordNotFound.Is(errs.Unwrap(err))) {
		return err
	}
	var localUsers []*model_struct.LocalUser
	if err == nil {
		localUsers = []*model_struct.LocalUser{localUser}
	}
	log.ZDebug(ctx, "SyncLoginUserInfo", "remoteUser", remoteUser, "localUser", localUser)
	return u.userSyncer.Sync(ctx, []*model_struct.LocalUser{remoteUser}, localUsers, nil, false, true)
}

func (u *User) SyncAllCommand(ctx context.Context) error {
	return u.syncAllCommand(ctx, true)
}

func (u *User) SyncAllCommandWithoutNotice(ctx context.Context) error {
	return u.syncAllCommand(ctx, false)
}

func (u *User) syncAllCommand(ctx context.Context, withNotice bool) error {
	resp, err := u.processUserCommandGetAll(ctx, &userPb.ProcessUserCommandGetAllReq{UserID: u.loginUserID})
	if err != nil {
		return err
	}
	localData, err := u.DataBase.ProcessUserCommandGetAll(ctx)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "sync command", "data from server", resp, "data from local", localData)
	if withNotice {
		return u.commandSyncer.Sync(ctx, datautil.Batch(ServerCommandToLocalCommand, resp.CommandResp), localData, nil)
	} else {
		return u.commandSyncer.Sync(ctx, datautil.Batch(ServerCommandToLocalCommand, resp.CommandResp), localData, nil, false, true)
	}
}
