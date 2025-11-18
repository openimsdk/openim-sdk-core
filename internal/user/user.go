package user

import (
	"context"
	"fmt"
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdkerrs"
	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/cache"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/utils/datautil"
)

// NewUser creates a new User object.
func NewUser(conversationEventQueue chan common.Cmd2Value) *User {
	user := &User{conversationEventQueue: conversationEventQueue}
	user.initSyncer()
	return user
}

// User is a struct that represents a user in the system.
type User struct {
	db_interface.DataBase
	loginUserID            string
	listener               func() open_im_sdk_callback.OnUserListener
	userSyncer             *syncer.Syncer[*model_struct.LocalUser, syncer.NoResp, string]
	conversationEventQueue chan common.Cmd2Value
	userCache              *cache.UserCache[string, *model_struct.LocalUser]
	once                   sync.Once

	//OnlineStatusCache *cache.Cache[string, *userPb.OnlineStatus]
}

func (u *User) UserCache() *cache.UserCache[string, *model_struct.LocalUser] {
	u.once.Do(func() {
		u.userCache = cache.NewUserCache[string, *model_struct.LocalUser](
			func(value *model_struct.LocalUser) string { return value.UserID },
			nil,
			u.GetLoginUser,
			u.GetUsersInfoFromServer,
		)
	})
	return u.userCache
}

// SetDataBase sets the DataBase field in User struct
func (u *User) SetDataBase(db db_interface.DataBase) {
	u.DataBase = db
}

// SetLoginUserID sets the loginUserID field in User struct
func (u *User) SetLoginUserID(loginUserID string) {
	u.loginUserID = loginUserID
}

// SetListener sets the user's listener.
func (u *User) SetListener(listener func() open_im_sdk_callback.OnUserListener) {
	u.listener = listener
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
			u.UserCache().Delete(localUser.UserID)
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
					_ = common.DispatchUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
						Args: common.SourceIDAndSessionType{SourceID: server.UserID, SessionType: constant.SingleChatType, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationEventQueue)
					_ = common.DispatchUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
						Args: common.UpdateMessageInfo{SessionType: constant.SingleChatType, UserID: server.UserID, FaceURL: server.FaceURL, Nickname: server.Nickname}}, u.conversationEventQueue)
				}
			}
			return nil
		},
	)
}

func (u *User) GetUserInfoWithCache(ctx context.Context, cacheKey string) (*model_struct.LocalUser, error) {
	return u.UserCache().Fetch(ctx, cacheKey)
}

func (u *User) GetUsersInfoWithCache(ctx context.Context, cacheKeys []string) ([]*model_struct.LocalUser, error) {
	m, err := u.UserCache().BatchFetch(ctx, cacheKeys)
	if err != nil {
		return nil, err
	}
	return datautil.Values(m), nil
}

// GetSingleUserFromServer retrieves user information from the server.
func (u *User) GetSingleUserFromServer(ctx context.Context, userID string) (*model_struct.LocalUser, error) {
	users, err := u.getUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return ServerUserToLocalUser(users[0]), nil
	}
	return nil, sdkerrs.ErrUserIDNotFound.WrapMsg(fmt.Sprintf("getSelfUserInfo failed, userID: %s not exist", userID))
}

// GetUsersInfoFromServer retrieves user information from the server.
func (u *User) GetUsersInfoFromServer(ctx context.Context, userIDs []string) ([]*model_struct.LocalUser, error) {
	var err error

	serverUsersInfo, err := u.getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	if len(serverUsersInfo) == 0 {
		log.ZError(ctx, "serverUsersInfo is empty", err, "userIDs", userIDs)
		return nil, err
	}

	return datautil.Batch(ServerUserToLocalUser, serverUsersInfo), nil
}
