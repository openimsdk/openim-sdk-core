package full

import (
	"context"
	"fmt"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	api "github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

func (u *Full) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string, sessionType int32) (*model_struct.LocalGroup, error) {
	switch sessionType {
	case constant.GroupChatType:
		return u.group.GetGroupInfoFromLocal2Svr(ctx, groupID)
	case constant.SuperGroupChatType:
		return u.GetGroupInfoByGroupID(ctx, groupID)
	default:
		return nil, fmt.Errorf("sessionType is not support %d", sessionType)
	}
}
func (u *Full) GetReadDiffusionGroupIDList(ctx context.Context) ([]string, error) {
	g, err := u.group.GetJoinedDiffusionGroupIDListFromSvr(ctx)
	if err != nil {
		return nil, err
	}
	return g, err
}

func (u *Full) GetGroupInfoByGroupID(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	g2, err := u.group.GetGroupInfoFromLocal2Svr(ctx, groupID)
	return g2, err
}

func (u *Full) GetGroupsInfo(ctx context.Context, groupIDs ...string) (map[string]*model_struct.LocalGroup, error) {
	return u.group.GetGroupsInfoFromLocal2Svr(ctx, groupIDs...)
}

func (u *Full) GetUsersInfo(ctx context.Context, userIDs []string) ([]*api.FullUserInfo, error) {
	friendList, err := u.db.GetFriendInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	blackList, err := u.db.GetBlackInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	users, err := u.user.GetServerUserInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	friendMap := make(map[string]*model_struct.LocalFriend)
	for i, f := range friendList {
		friendMap[f.FriendUserID] = friendList[i]
	}
	blackMap := make(map[string]*model_struct.LocalBlack)
	for i, b := range blackList {
		blackMap[b.BlockUserID] = blackList[i]
	}
	userMap := make(map[string]*api.PublicUser)
	for _, info := range users {
		userMap[info.UserID] = &api.PublicUser{
			UserID:     info.UserID,
			Nickname:   info.Nickname,
			FaceURL:    info.FaceURL,
			Ex:         info.Ex,
			CreateTime: info.CreateTime,
		}
	}
	res := make([]*api.FullUserInfo, 0, len(users))
	for _, userID := range userIDs {
		info, ok := userMap[userID]
		if !ok {
			continue
		}
		res = append(res, &api.FullUserInfo{
			PublicInfo: info,
			FriendInfo: friendMap[userID],
			BlackInfo:  blackMap[userID],
		})

		// update single conversation

		conversation, err := u.db.GetConversationByUserID(ctx, userID)
		if err != nil {
			log.ZWarn(ctx, "GetConversationByUserID failed", err, "userID", userID)
		} else {
			if _, ok := friendMap[userID]; ok {
				continue
			}
			log.ZDebug(ctx, "GetConversationByUserID", "conversation", conversation)
			if conversation.ShowName != info.Nickname || conversation.FaceURL != info.FaceURL {
				_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
					Args: common.SourceIDAndSessionType{SourceID: userID, SessionType: conversation.ConversationType, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{SessionType: conversation.ConversationType, UserID: userID, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
			}
		}
	}
	return res, nil
}

func (u *Full) GetUsersInfoWithCache(ctx context.Context, userIDs []string, groupID string) ([]*api.FullUserInfoWithCache, error) {
	friendList, err := u.db.GetFriendInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	blackList, err := u.db.GetBlackInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	users, err := u.user.GetServerUserInfo(ctx, userIDs)
	if err == nil {
		var strangers []*model_struct.LocalStranger
		for _, val := range users {
			strangerTemp := &model_struct.LocalStranger{
				UserID:           val.UserID,
				Nickname:         val.Nickname,
				FaceURL:          val.FaceURL,
				CreateTime:       val.CreateTime,
				AppMangerLevel:   val.AppMangerLevel,
				Ex:               val.Ex,
				AttachedInfo:     val.Ex,
				GlobalRecvMsgOpt: val.GlobalRecvMsgOpt,
			}
			strangers = append(strangers, strangerTemp)
		}
		err := u.db.SetStrangerInfo(ctx, strangers)
		if err != nil {
			return nil, err
		}
	} else {
		strangerList, err := u.db.GetStrangerInfo(ctx, userIDs)
		if err != nil {
			return nil, err
		}
		for _, val := range strangerList {
			userTemp := &sdkws.UserInfo{
				UserID:           val.UserID,
				Nickname:         val.Nickname,
				FaceURL:          val.FaceURL,
				Ex:               val.Ex,
				CreateTime:       val.CreateTime,
				AppMangerLevel:   val.AppMangerLevel,
				GlobalRecvMsgOpt: val.GlobalRecvMsgOpt,
			}
			users = append(users, userTemp)
		}
	}
	var groupMemberList []*model_struct.LocalGroupMember
	if groupID != "" {
		groupMemberList, err = u.db.GetGroupSomeMemberInfo(ctx, groupID, userIDs)
		if err != nil {
			return nil, err
		}
	}
	friendMap := make(map[string]*model_struct.LocalFriend)
	for i, f := range friendList {
		friendMap[f.FriendUserID] = friendList[i]
	}
	blackMap := make(map[string]*model_struct.LocalBlack)
	for i, b := range blackList {
		blackMap[b.BlockUserID] = blackList[i]
	}
	groupMemberMap := make(map[string]*model_struct.LocalGroupMember)
	for i, b := range groupMemberList {
		groupMemberMap[b.UserID] = groupMemberList[i]
	}
	userMap := make(map[string]*api.PublicUser)
	for _, info := range users {
		userMap[info.UserID] = &api.PublicUser{
			UserID:     info.UserID,
			Nickname:   info.Nickname,
			FaceURL:    info.FaceURL,
			Ex:         info.Ex,
			CreateTime: info.CreateTime,
		}
	}
	res := make([]*api.FullUserInfoWithCache, 0, len(users))
	for _, userID := range userIDs {
		info, ok := userMap[userID]
		if !ok {
			continue
		}
		res = append(res, &api.FullUserInfoWithCache{
			PublicInfo:      info,
			FriendInfo:      friendMap[userID],
			BlackInfo:       blackMap[userID],
			GroupMemberInfo: groupMemberMap[userID],
		})

		// update single conversation

		conversation, err := u.db.GetConversationByUserID(ctx, userID)
		if err != nil {
			log.ZWarn(ctx, "GetConversationByUserID failed", err, "userID", userID)
		} else {
			if _, ok := friendMap[userID]; ok {
				continue
			}
			log.ZDebug(ctx, "GetConversationByUserID", "conversation", conversation)
			if conversation.ShowName != info.Nickname || conversation.FaceURL != info.FaceURL {
				_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName,
					Args: common.SourceIDAndSessionType{SourceID: userID, SessionType: conversation.ConversationType, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{SessionType: conversation.ConversationType, UserID: userID, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
			}
		}
	}
	return res, nil
}
