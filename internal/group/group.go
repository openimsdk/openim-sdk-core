// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package group

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/syncer"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func NewGroup(loginUserID string, db db_interface.DataBase,
	joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value,
	conversationCh chan common.Cmd2Value) *Group {
	g := &Group{
		loginUserID:        loginUserID,
		db:                 db,
		joinedSuperGroupCh: joinedSuperGroupCh,
		heartbeatCmdCh:     heartbeatCmdCh,
		conversationCh:     conversationCh,
	}
	g.initSyncer()
	return g
}

// //utils.GetCurrentTimestampByMill()
type Group struct {
	listener                open_im_sdk_callback.OnGroupListener
	loginUserID             string
	db                      db_interface.DataBase
	groupSyncer             *syncer.Syncer[*model_struct.LocalGroup, string]
	groupMemberSyncer       *syncer.Syncer[*model_struct.LocalGroupMember, [2]string]
	groupRequestSyncer      *syncer.Syncer[*model_struct.LocalGroupRequest, [2]string]
	groupAdminRequestSyncer *syncer.Syncer[*model_struct.LocalAdminGroupRequest, [2]string]
	loginTime               int64
	joinedSuperGroupCh      chan common.Cmd2Value
	heartbeatCmdCh          chan common.Cmd2Value

	conversationCh chan common.Cmd2Value
	//	memberSyncMutex sync.RWMutex

	listenerForService open_im_sdk_callback.OnListenerForService
}

func (g *Group) initSyncer() {
	g.groupSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroup) error {
		return g.db.InsertGroup(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroup) error {
		return g.db.DeleteGroup(ctx, value.GroupID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroup) error {
		return g.db.UpdateGroup(ctx, server)
	}, func(value *model_struct.LocalGroup) string {
		return value.GroupID
	}, nil, nil)

	g.groupMemberSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.InsertGroupMember(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.DeleteGroupMember(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupMember) error {
		return g.db.UpdateGroupMember(ctx, server)
	}, func(value *model_struct.LocalGroupMember) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

	g.groupRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.InsertGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.DeleteGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupRequest) error {
		return g.db.UpdateGroupRequest(ctx, server)
	}, func(value *model_struct.LocalGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

	g.groupAdminRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.InsertAdminGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.DeleteAdminGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalAdminGroupRequest) error {
		return g.db.UpdateAdminGroupRequest(ctx, server)
	}, func(value *model_struct.LocalAdminGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, nil)

}

func (g *Group) SetGroupListener(callback open_im_sdk_callback.OnGroupListener) {
	if callback == nil {
		return
	}
	g.listener = callback
}

func (g *Group) LoginTime() int64 {
	return g.loginTime
}

func (g *Group) SetLoginTime(loginTime int64) {
	g.loginTime = loginTime
}

func (g *Group) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	g.listenerForService = listener
}

func (g *Group) GetGroupOwnerIDAndAdminIDList(ctx context.Context, groupID string) (ownerID string, adminIDList []string, err error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	adminIDList, err = g.db.GetGroupAdminID(ctx, groupID)
	if err != nil {
		return "", nil, err
	}
	return localGroup.OwnerUserID, adminIDList, nil
}

func (g *Group) GetGroupInfoFromLocal2Svr(ctx context.Context, groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(ctx, groupID)
	if err == nil {
		return localGroup, nil
	}
	svrGroup, err := g.getGroupsInfoFromSvr(ctx, []string{groupID})
	if err != nil {
		return nil, err
	}
	if len(svrGroup) == 0 {
		return nil, errs.ErrGroupIDNotFound.Wrap("server not this group")
	}
	return ServerGroupToLocalGroup(svrGroup[0]), nil
}

func (g *Group) getGroupsInfoFromSvr(ctx context.Context, groupIDs []string) ([]*sdkws.GroupInfo, error) {
	resp, err := util.CallApi[group.GetGroupsInfoResp](ctx, constant.GetGroupsInfoRouter, &group.GetGroupsInfoReq{GroupIDs: groupIDs})
	if err != nil {
		return nil, err
	}
	return resp.GroupInfos, nil
}

func (g *Group) getGroupAbstractInfoFromSvr(ctx context.Context, groupIDs []string) (*group.GetGroupAbstractInfoResp, error) {
	return util.CallApi[group.GetGroupAbstractInfoResp](ctx, constant.GetGroupAbstractInfoRouter, &group.GetGroupAbstractInfoReq{GroupIDs: groupIDs})
}

func (g *Group) GetJoinedDiffusionGroupIDListFromSvr(ctx context.Context) ([]string, error) {
	groups, err := g.GetServerJoinGroup(ctx)
	if err != nil {
		return nil, err
	}
	var groupIDs []string
	for _, g := range groups {
		if g.GroupType == constant.WorkingGroup {
			groupIDs = append(groupIDs, g.GroupID)
		}
	}
	return groupIDs, nil
}

func (g *Group) SyncJoinedGroupMemberForFirstLogin(ctx context.Context) error {
	groups, err := g.syncJoinedGroup(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			if err := g.SyncGroupMember(ctx, groupID); err != nil {
				log.ZError(ctx, "SyncGroupMember failed", err)
			}
		}(group.GroupID)
	}
	wg.Wait()
	return nil
}
