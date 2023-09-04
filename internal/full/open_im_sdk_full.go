// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package full

import (
	"context"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
)

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
					Args: common.SourceIDAndSessionType{SourceID: userID, SessionType: constant.SingleChatType, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{UserID: userID, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
			}
		}
	}
	return res, nil
}

func (u *Full) GetUsersInfoStranger(ctx context.Context, userIDs []string, groupID string) ([]*api.FullUserInfoStranger, error) {
	friendList, err := u.db.GetFriendInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	blackList, err := u.db.GetBlackInfoList(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	strangerFlag := false
	users, err := u.user.GetServerUserInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if users == nil {
		strangerFlag = true
	}
	var groupMemberList []*model_struct.LocalGroupMember
	if groupID != "" {
		groupMemberList, err = u.db.GetGroupSomeMemberInfo(ctx, groupID, userIDs)
		if err != nil {
			return nil, err
		}
	}
	if !strangerFlag {
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
	res := make([]*api.FullUserInfoStranger, 0, len(users))
	for _, userID := range userIDs {
		info, ok := userMap[userID]
		if !ok {
			continue
		}
		res = append(res, &api.FullUserInfoStranger{
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
					Args: common.SourceIDAndSessionType{SourceID: userID, SessionType: constant.SingleChatType, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName,
					Args: common.UpdateMessageInfo{UserID: userID, FaceURL: info.FaceURL, Nickname: info.Nickname}}, u.ch)
			}
		}
	}
	return res, nil
}
