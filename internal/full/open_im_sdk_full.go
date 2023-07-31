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
	"open_im_sdk/pkg/db/model_struct"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/tools/log"
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
		if err != nil || conversation.ConversationID == "" {
			log.ZWarn(ctx, "GetUsersInfo GetConversationByUserID failed", err, "userID", userID)
		} else {
			log.ZDebug(ctx, "GetConversationByUserID", "conversation", conversation)
			if conversation.ShowName != info.Nickname || conversation.FaceURL != info.FaceURL {
				friends, err := u.db.GetFriendInfoList(ctx, []string{userID})
				if err != nil {
					log.ZError(ctx, "GetUsersInfo GetFriendInfoList failed", err, "userID", userID)
				} else {
					// not friend
					if len(friends) == 0 {
						conversation.ShowName = info.Nickname
					} else {
						// friend and have remark
						if friends[0].Remark != "" {
							conversation.ShowName = friends[0].Remark
						}
					}
				}
				conversation.FaceURL = info.FaceURL
				if err := u.db.UpdateConversation(ctx, conversation); err != nil {
					log.ZError(ctx, "GetUsersInfo UpdateMsgSenderFaceURLAndSenderNickname failed", err, "userID", userID)
				} else {
					u.conversationListner.OnConversationChanged(utils.StructToJsonString([]*model_struct.LocalConversation{conversation}))
				}
			}
		}
		// update joined groups
		groupIDs, err := u.db.GetUserJoinedGroupIDs(ctx, userID)
		if err != nil {
			log.ZError(ctx, "GetUsersInfo GetUserJoinedGroupIDs failed", err, "userID", userID)
		} else {
			for _, groupID := range groupIDs {
				localGroupMember, err := u.db.GetGroupMemberInfoByGroupIDUserID(ctx, groupID, userID)
				if err != nil {
					log.ZError(ctx, "GetUsersInfo GetGroupMemberInfoByGroupIDUserID failed", err, "userID", userID, "groupID", groupID)
				} else {
					if localGroupMember.Nickname != info.Nickname || localGroupMember.FaceURL != info.FaceURL {
						if err := u.db.UpdateMsgSenderFaceURLAndSenderNickname(ctx, utils.GetConversationIDByGroupID(groupID), userID, info.FaceURL, info.Nickname); err != nil {
							log.ZError(ctx, "UpdateMsgSenderFaceURLAndSenderNickname failed", err, "conversationID", utils.GetConversationIDByGroupID(groupID), "userID", userID, "nickname", info.Nickname, "faceURL", info.FaceURL)
						}
						localGroupMember.FaceURL = info.FaceURL
						localGroupMember.Nickname = info.Nickname
						if err := u.db.UpdateGroupMember(ctx, localGroupMember); err != nil {
							log.ZError(ctx, "UpdateGroupMember failed", err, "userID", userID)
						}
					}
				}

			}
		}
	}
	return res, nil
}
