package user

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

// NewUser creates a new User object.
func NewUser(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *User {
	user := &User{DataBase: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	user.initSyncer()
	user.UserBasicCache = cache.NewCache[string, *sdk_struct.BasicInfo]()
	//user.OnlineStatusCache = cache.NewCache[string, *userPb.OnlineStatus]()
	return user
}

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID    string
	listener       func() open_im_sdk_callback.OnUserListener
	userSyncer     *syncer.Syncer[*model_struct.LocalUser, syncer.NoResp, string]
	commandSyncer  *syncer.Syncer[*model_struct.LocalUserCommand, syncer.NoResp, string]
	conversationCh chan common.Cmd2Value
	UserBasicCache *cache.Cache[string, *sdk_struct.BasicInfo]

	//OnlineStatusCache *cache.Cache[string, *userPb.OnlineStatus]
}

func (u *User) initSyncer() {
	u.userSyncer = syncer.New[*model_struct.LocalUser, syncer.NoResp, string](
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return u.InsertLoginUser(ctx, value)
		},
		func(ctx context.Context, value *model_struct.LocalUser) error {
			return fmt.Errorf("not support delete user %s", value.UserID)
		},
		func(ctx context.Context, serverUser, localUser *model_struct.LocalUser) error {
			return u.DataBase.UpdateLoginUser(context.Background(), serverUser)
		},
		func(user *model_struct.LocalUser) string {
			return user.UserID
		},
		nil,
		func(ctx context.Context, state int, server, local *model_struct.LocalUser) error {
			switch state {
			case syncer.Update:
				u.listener().OnSelfInfoUpdated(utils.StructToJsonString(server))
				if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
					_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{SessionType: constant.SingleChatType, UserID: server.UserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationCh)
				}
			}
			return nil
		},
	)
	u.commandSyncer = syncer.New[*model_struct.LocalUserCommand, syncer.NoResp, string](
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to insert a command
			return u.DataBase.ProcessUserCommandAdd(ctx, command)
		},
		func(ctx context.Context, command *model_struct.LocalUserCommand) error {
			// Logic to delete a command
			return u.DataBase.ProcessUserCommandDelete(ctx, command)
		},
		func(ctx context.Context, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			// Logic to update a command
			if serverCommand == nil || localCommand == nil {
				return fmt.Errorf("nil command reference")
			}
			return u.DataBase.ProcessUserCommandUpdate(ctx, serverCommand)
		},
		func(command *model_struct.LocalUserCommand) string {
			// Return a unique identifier for the command
			if command == nil {
				return ""
			}
			return command.Uuid
		},
		func(a *model_struct.LocalUserCommand, b *model_struct.LocalUserCommand) bool {
			// Compare two commands to check if they are equal
			if a == nil || b == nil {
				return false
			}
			return a.Uuid == b.Uuid && a.Type == b.Type && a.Value == b.Value
		},
		func(ctx context.Context, state int, serverCommand *model_struct.LocalUserCommand, localCommand *model_struct.LocalUserCommand) error {
			if u.listener == nil {
				return nil
			}
			switch state {
			case syncer.Delete:
				u.listener().OnUserCommandDelete(utils.StructToJsonString(serverCommand))
			case syncer.Update:
				u.listener().OnUserCommandUpdate(utils.StructToJsonString(serverCommand))
			case syncer.Insert:
				u.listener().OnUserCommandAdd(utils.StructToJsonString(serverCommand))
			}
			return nil
		},
	)
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener func() open_im_sdk_callback.OnUserListener) {
	u.listener = listener
}

// getSelfUserInfo retrieves the user's information.
func (u *User) getSelfUserInfo(ctx context.Context) (*model_struct.LocalUser, error) {
	userInfo, errLocal := u.GetLoginUser(ctx, u.loginUserID)
	if errLocal != nil {
		srvUserInfo, errServer := u.GetServerUserInfo(ctx, []string{u.loginUserID})
		if errServer != nil {
			return nil, errServer
		}
		if len(srvUserInfo) == 0 {
			return nil, sdkerrs.ErrUserIDNotFound
		}
		userInfo = ServerUserToLocalUser(srvUserInfo[0])
		_ = u.InsertLoginUser(ctx, userInfo)
	}
	return userInfo, nil
}
