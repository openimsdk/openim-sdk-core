package user

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/protocol/sdkws"
	userPb "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
)

// GetSingleUserFromSvr retrieves user information from the server.
func (u *User) GetSingleUserFromSvr(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	users, err := u.GetUsersInfoFromSvr(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return users[0], nil
	}
	return nil, sdkerrs.ErrUserIDNotFound.WrapMsg(fmt.Sprintf("getSelfUserInfo failed, userID: %s not exist", userID))
}

// ProcessUserCommandGetAll get user's choice
func (u *User) ProcessUserCommandGetAll(ctx context.Context) ([]*userPb.CommandInfoResp, error) {
	localCommands, err := u.DataBase.ProcessUserCommandGetAll(ctx)
	if err != nil {
		return nil, err // Handle the error appropriately
	}
	var result []*userPb.CommandInfoResp
	for _, localCommand := range localCommands {
		result = append(result, &userPb.CommandInfoResp{
			Type:       localCommand.Type,
			CreateTime: localCommand.CreateTime,
			Uuid:       localCommand.Uuid,
			Value:      localCommand.Value,
		})
	}
	return result, nil
}

func (u *User) UserOnlineStatusChange(users map[string][]int32) {
	for userID, onlinePlatformIDs := range users {
		status := userPb.OnlineStatus{
			UserID:      userID,
			PlatformIDs: onlinePlatformIDs,
		}
		if len(status.PlatformIDs) == 0 {
			status.Status = constant.Offline
		} else {
			status.Status = constant.Online
		}
		u.listener().OnUserStatusChanged(utils.StructToJsonString(&status))
	}
}

func (u *User) GetSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	return u.getSelfUserInfo(ctx)
}

func (u *User) SetSelfInfo(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	return u.updateSelfUserInfo(ctx, userInfo)
}
func (u *User) SetGlobalRecvMessageOpt(ctx context.Context, opt int) error {
	if err := u.setGlobalRecvMessageOpt(ctx, int32(opt)); err != nil {
		return err
	}
	err := u.SyncLoginUserInfo(ctx)
	if err != nil {
		log.ZWarn(ctx, "SyncLoginUserInfo", err)
	}
	return nil
}

// ProcessUserCommandAdd CRUD user command
func (u *User) ProcessUserCommandAdd(ctx context.Context, userCommand *userPb.ProcessUserCommandAddReq) error {
	req := &userPb.ProcessUserCommandAddReq{UserID: u.loginUserID, Type: userCommand.Type, Uuid: userCommand.Uuid, Value: userCommand.Value}
	if err := u.processUserCommandAdd(ctx, req); err != nil {
		return err
	}
	return u.SyncAllCommand(ctx)
}

// ProcessUserCommandDelete delete user's choice
func (u *User) ProcessUserCommandDelete(ctx context.Context, userCommand *userPb.ProcessUserCommandDeleteReq) error {
	req := &userPb.ProcessUserCommandDeleteReq{UserID: u.loginUserID, Type: userCommand.Type, Uuid: userCommand.Uuid}
	if err := u.processUserCommandDelete(ctx, req); err != nil {
		return err
	}
	return u.SyncAllCommand(ctx)
}

// ProcessUserCommandUpdate update user's choice
func (u *User) ProcessUserCommandUpdate(ctx context.Context, userCommand *userPb.ProcessUserCommandUpdateReq) error {
	req := &userPb.ProcessUserCommandUpdateReq{UserID: u.loginUserID, Type: userCommand.Type, Uuid: userCommand.Uuid, Value: userCommand.Value}
	if err := u.processUserCommandUpdate(ctx, req); err != nil {
		return err
	}
	return u.SyncAllCommand(ctx)
}

// GetServerUserInfo retrieves user information from the server.
func (u *User) GetServerUserInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	return u.getUsersInfo(ctx, userIDs)
}
