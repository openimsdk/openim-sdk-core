// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package full

import (
	"context"
	"fmt"
	"open_im_sdk/internal/cache"
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/super_group"
	"open_im_sdk/internal/user"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
)

type Full struct {
	user       *user.User
	friend     *friend.Friend
	group      *group.Group
	ch         chan common.Cmd2Value
	userCache  *cache.Cache
	db         db_interface.DataBase
	SuperGroup *super_group.SuperGroup
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group, ch chan common.Cmd2Value, userCache *cache.Cache, db db_interface.DataBase, superGroup *super_group.SuperGroup) *Full {
	return &Full{user: user, friend: friend, group: group, ch: ch, userCache: userCache, db: db, SuperGroup: superGroup}
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
	g1, err1 := u.group.GetJoinedDiffusionGroupIDListFromSvr(ctx)
	g2, err2 := u.SuperGroup.GetJoinedGroupIDListFromSvr(ctx)
	var groupIDList []string
	if err1 == nil {
		groupIDList = append(groupIDList, g1...)
	}
	if err2 == nil {
		groupIDList = append(groupIDList, g2...)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return groupIDList, err
}
