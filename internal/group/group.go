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

package group

import (
	"context"
	"open_im_sdk/internal/util"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/sdkerrs"
	"open_im_sdk/pkg/syncer"
	"open_im_sdk/pkg/utils"
	"sync"

	"github.com/OpenIMSDK/protocol/group"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	utils2 "github.com/OpenIMSDK/tools/utils"
)

func NewGroup(loginUserID string, db db_interface.DataBase,
	conversationCh chan common.Cmd2Value) *Group {
	g := &Group{
		loginUserID:    loginUserID,
		db:             db,
		conversationCh: conversationCh,
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
		if err := g.db.DeleteGroupAllMembers(ctx, value.GroupID); err != nil {
			return err
		}
		return g.db.DeleteGroup(ctx, value.GroupID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroup) error {
		log.ZInfo(ctx, "groupSyncer trigger update funcation", "groupID", server.GroupID, "server", server, "local", local)
		return g.db.UpdateGroup(ctx, server)
	}, func(value *model_struct.LocalGroup) string {
		return value.GroupID
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroup) error {
		if g.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			g.listener.OnJoinedGroupAdded(utils.StructToJsonString(server))
		case syncer.Delete:
			g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(local))
		case syncer.Update:
			log.ZInfo(ctx, "groupSyncer trigger update", "groupID", server.GroupID, "data", server, "isDismissed", server.Status == constant.GroupStatusDismissed)
			if server.Status == constant.GroupStatusDismissed {
				if err := g.db.DeleteGroupAllMembers(ctx, server.GroupID); err != nil {
					log.ZError(ctx, "delete group all members failed", err)
				}
				g.listener.OnGroupDismissed(utils.StructToJsonString(server))
			} else {
				g.listener.OnGroupInfoChanged(utils.StructToJsonString(server))
				if server.GroupName != local.GroupName || local.FaceURL != server.FaceURL {
					_ = common.TriggerCmdUpdateConversation(ctx, common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: server.GroupID,
						SessionType: constant.SuperGroupChatType, FaceURL: server.FaceURL, Nickname: server.GroupName}}, g.conversationCh)
				}
			}
		}

		return nil
	})

	g.groupMemberSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.InsertGroupMember(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupMember) error {
		return g.db.DeleteGroupMember(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupMember) error {
		return g.db.UpdateGroupMember(ctx, server)
	}, func(value *model_struct.LocalGroupMember) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroupMember) error {
		if g.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			g.listener.OnGroupMemberAdded(utils.StructToJsonString(server))
		case syncer.Delete:
			g.listener.OnGroupMemberDeleted(utils.StructToJsonString(local))
		case syncer.Update:
			g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(server))
			if server.Nickname != local.Nickname || server.FaceURL != local.FaceURL {
				_ = common.TriggerCmdUpdateMessage(ctx, common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: server.UserID, FaceURL: server.FaceURL,
					Nickname: server.Nickname, GroupID: server.GroupID}}, g.conversationCh)
			}
		}
		return nil
	})

	g.groupRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.InsertGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalGroupRequest) error {
		return g.db.DeleteGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalGroupRequest) error {
		return g.db.UpdateGroupRequest(ctx, server)
	}, func(value *model_struct.LocalGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

	g.groupAdminRequestSyncer = syncer.New(func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.InsertAdminGroupRequest(ctx, value)
	}, func(ctx context.Context, value *model_struct.LocalAdminGroupRequest) error {
		return g.db.DeleteAdminGroupRequest(ctx, value.GroupID, value.UserID)
	}, func(ctx context.Context, server, local *model_struct.LocalAdminGroupRequest) error {
		return g.db.UpdateAdminGroupRequest(ctx, server)
	}, func(value *model_struct.LocalAdminGroupRequest) [2]string {
		return [...]string{value.GroupID, value.UserID}
	}, nil, func(ctx context.Context, state int, server, local *model_struct.LocalAdminGroupRequest) error {
		switch state {
		case syncer.Insert:
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
		case syncer.Update:
			switch server.HandleResult {
			case constant.FriendResponseAgree:
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(server))
			case constant.FriendResponseRefuse:
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(server))
			default:
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(server))
			}
		}
		return nil
	})

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
		return nil, sdkerrs.ErrGroupIDNotFound.Wrap("server not this group")
	}
	return ServerGroupToLocalGroup(svrGroup[0]), nil
}

func (g *Group) GetGroupsInfoFromLocal2Svr(ctx context.Context, groupIDs ...string) (map[string]*model_struct.LocalGroup, error) {
	groupMap := make(map[string]*model_struct.LocalGroup)
	if len(groupIDs) == 0 {
		return groupMap, nil
	}
	groups, err := g.db.GetGroups(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	var groupIDsNeedSync []string
	localGroupIDs := utils2.Slice(groups, func(group *model_struct.LocalGroup) string {
		return group.GroupID
	})
	for _, groupID := range groupIDs {
		if !utils2.Contain(groupID, localGroupIDs...) {
			groupIDsNeedSync = append(groupIDsNeedSync, groupID)
		}
	}

	if len(groupIDsNeedSync) > 0 {
		svrGroups, err := g.getGroupsInfoFromSvr(ctx, groupIDsNeedSync)
		if err != nil {
			return nil, err
		}
		for _, svrGroup := range svrGroups {
			groups = append(groups, ServerGroupToLocalGroup(svrGroup))
		}
	}
	for _, group := range groups {
		groupMap[group.GroupID] = group
	}
	return groupMap, nil
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

func (g *Group) SyncAllJoinedGroupMembers(ctx context.Context) error {
	groups, err := g.syncAllJoinedGroups(ctx)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, group := range groups {
		wg.Add(1)
		go func(groupID string) {
			defer wg.Done()
			if err := g.SyncAllGroupMember(ctx, groupID); err != nil {
				log.ZError(ctx, "SyncGroupMember failed", err)
			}
		}(group.GroupID)
	}
	wg.Wait()
	return nil
}
