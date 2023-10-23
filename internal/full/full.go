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
	"github.com/openimsdk/openim-sdk-core/v3/internal/friend"
	"github.com/openimsdk/openim-sdk-core/v3/internal/group"
	"github.com/openimsdk/openim-sdk-core/v3/internal/user"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

type Full struct {
	user   *user.User
	friend *friend.Friend
	group  *group.Group
	ch     chan common.Cmd2Value
	db     db_interface.DataBase
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group, ch chan common.Cmd2Value,
	db db_interface.DataBase) *Full {
	return &Full{user: user, friend: friend, group: group, ch: ch, db: db}
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
