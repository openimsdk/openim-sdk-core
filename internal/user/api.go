package user

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/sdkws"
	userPb "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
)

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
	return u.GetUserInfoWithCache(ctx, u.loginUserID)
}

func (u *User) SetSelfInfo(ctx context.Context, userInfo *sdkws.UserInfoWithEx) error {
	// updateSelfUserInfo updates the user's information with Ex field.
	userInfo.UserID = u.loginUserID
	if err := u.updateUserInfo(ctx, userInfo); err != nil {
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

func (u *User) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdk_struct.PublicUser, error) {
	usersInfo, err := u.GetUsersInfoWithCache(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	res := datautil.Batch(LocalUserToPublicUser, usersInfo)

	friendList, err := u.GetFriendInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	friendMap := datautil.SliceToMap(friendList, func(friend *model_struct.LocalFriend) string {
		return friend.FriendUserID
	})

	for _, userInfo := range res {

		// update single conversation

		conversation, err := u.GetConversationByUserID(ctx, userInfo.UserID)
		if err != nil {
			log.ZWarn(ctx, "GetConversationByUserID failed", err, "userInfo", usersInfo)
		} else {
			if _, ok := friendMap[userInfo.UserID]; ok {
				continue
			}
			log.ZDebug(ctx, "GetConversationByUserID", "conversation", conversation)
			if conversation.ShowName != userInfo.Nickname || conversation.FaceURL != userInfo.FaceURL {
				_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
					Args: common.SourceIDAndSessionType{SourceID: userInfo.UserID, SessionType: conversation.ConversationType, FaceURL: userInfo.FaceURL, Nickname: userInfo.Nickname}}, u.conversationCh)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{SessionType: conversation.ConversationType, UserID: userInfo.UserID, FaceURL: userInfo.FaceURL, Nickname: userInfo.Nickname}}, u.conversationCh)
			}
		}
	}
	return res, nil
}
