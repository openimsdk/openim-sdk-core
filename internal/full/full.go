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
	"fmt"
	"open_im_sdk/internal/cache"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
)

type Full struct {
	user                *user.User
	friend              *friend.Friend
	group               *group.Group
	ch                  chan common.Cmd2Value
	userCache           *cache.Cache
	db                  db_interface.DataBase
	conversationListner open_im_sdk_callback.OnConversationListener
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group, ch chan common.Cmd2Value,
	userCache *cache.Cache, db db_interface.DataBase, conversationListner open_im_sdk_callback.OnConversationListener) *Full {
	return &Full{user: user, friend: friend, group: group, ch: ch, userCache: userCache, db: db, conversationListner: conversationListner}
}

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
