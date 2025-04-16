package user

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/cliconf"
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
					Args: common.SourceIDAndSessionType{SourceID: userInfo.UserID, SessionType: conversation.ConversationType, FaceURL: userInfo.FaceURL, Nickname: userInfo.Nickname}}, u.conversationEventQueue)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{SessionType: conversation.ConversationType, UserID: userInfo.UserID, FaceURL: userInfo.FaceURL, Nickname: userInfo.Nickname}}, u.conversationEventQueue)
			}
		}
	}
	return res, nil
}

func (u *User) GetUserClientConfig(ctx context.Context) (map[string]string, error) {
	res, err := cliconf.GetClientConfig(ctx)
	if err != nil {
		return nil, err
	}
	return res.RawConfig, nil
}
