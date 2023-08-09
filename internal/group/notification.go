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
	"encoding/json"
	"errors"
	"fmt"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
)

func (g *Group) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := g.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "DoGroupNotification failed", err)
		}
	}()
}

func (g *Group) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if g.listener == nil {
		return errors.New("listener is nil")
	}
	switch msg.ContentType {
	case constant.GroupCreatedNotification: // 1501
		var detail sdkws.GroupCreatedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupInfoSetNotification: // 1502
		var detail sdkws.GroupInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	case constant.JoinGroupApplicationNotification: // 1503
		var detail sdkws.JoinGroupApplicationTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.Applicant.UserID == g.loginUserID {
			return g.SyncSelfGroupApplications(ctx, detail.Group.GroupID)
		} else {
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
	case constant.GroupApplicationAcceptedNotification: // 1505
		var detail sdkws.GroupApplicationAcceptedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
		if detail.ReceiverAs == 1 {
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	case constant.GroupApplicationRejectedNotification: // 1506
		var detail sdkws.GroupApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
		if detail.ReceiverAs == 1 {
			return g.SyncAdminGroupApplications(ctx, detail.Group.GroupID)
		}
		return g.SyncSelfGroupApplications(ctx, detail.Group.GroupID)
	case constant.GroupOwnerTransferredNotification: // 1507
		var detail sdkws.GroupOwnerTransferredTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		if detail.Group == nil {
			return errors.New(fmt.Sprintf("group is nil, groupID: %s", detail.Group.GroupID))
		}
		return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberKickedNotification: // 1508
		var detail sdkws.MemberKickedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		var self bool
		for _, info := range detail.KickedUserList {
			if info.UserID == g.loginUserID {
				self = true
				break
			}
		}
		if self {
			members, err := g.db.GetGroupMemberListSplit(ctx, detail.Group.GroupID, 0, 0, 999999)
			if err != nil {
				return err
			}
			if err := g.db.DeleteGroupAllMembers(ctx, detail.Group.GroupID); err != nil {
				return err
			}
			for _, member := range members {
				data, err := json.Marshal(member)
				if err != nil {
					return err
				}
				g.listener.OnGroupMemberDeleted(string(data))
			}
			//for _, member := range util.Batch(ServerGroupMemberToLocalGroupMember, detail.KickedUserList) {
			//	data, err := json.Marshal(member)
			//	if err != nil {
			//		return err
			//	}
			//	g.listener.OnGroupMemberDeleted(string(data))
			//}
			group, err := g.db.GetGroupInfoByGroupID(ctx, detail.Group.GroupID)
			if err != nil {
				return err
			}
			group.MemberCount = 0
			data, err := json.Marshal(group)
			if err != nil {
				return err
			}
			if err := g.db.DeleteGroup(ctx, detail.Group.GroupID); err != nil {
				return err
			}
			g.listener.OnGroupInfoChanged(string(data))
			return nil
		} else {
			var userIDs []string
			for _, info := range detail.KickedUserList {
				userIDs = append(userIDs, info.UserID)
			}
			return g.SyncGroupMembers(ctx, detail.Group.GroupID, userIDs...)
		}
	case constant.MemberQuitNotification: // 1504
		var detail sdkws.MemberQuitTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.QuitUser.UserID == g.loginUserID {
			//if err := g.db.DeleteGroupAllMembers(ctx, detail.Group.GroupID); err != nil {
			//	return err
			//}
			//return g.SyncJoinedGroup(ctx)
			//group, err := g.db.GetGroupInfoByGroupID(ctx, detail.Group.GroupID)
			//if err != nil {
			//	return err
			//}
			//group.MemberCount = 0
			//data, err := json.Marshal(group)
			//if err != nil {
			//	return err
			//}
			//if err := g.db.DeleteGroup(ctx, detail.Group.GroupID); err != nil {
			//	return err
			//}
			//g.listener.OnGroupInfoChanged(string(data))
			//return nil
			members, err := g.db.GetGroupMemberListSplit(ctx, detail.Group.GroupID, 0, 0, 999999)
			if err != nil {
				return err
			}
			if err := g.db.DeleteGroupAllMembers(ctx, detail.Group.GroupID); err != nil {
				return err
			}
			for _, member := range members {
				data, err := json.Marshal(member)
				if err != nil {
					return err
				}
				g.listener.OnGroupMemberDeleted(string(data))
			}
			//for _, member := range util.Batch(ServerGroupMemberToLocalGroupMember, detail.KickedUserList) {
			//	data, err := json.Marshal(member)
			//	if err != nil {
			//		return err
			//	}
			//	g.listener.OnGroupMemberDeleted(string(data))
			//}
			group, err := g.db.GetGroupInfoByGroupID(ctx, detail.Group.GroupID)
			if err != nil {
				return err
			}
			group.MemberCount = 0
			data, err := json.Marshal(group)
			if err != nil {
				return err
			}
			if err := g.db.DeleteGroup(ctx, detail.Group.GroupID); err != nil {
				return err
			}
			g.listener.OnGroupInfoChanged(string(data))
			return nil
		} else {
			return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.QuitUser.UserID)
		}
	case constant.MemberInvitedNotification: // 1509
		var detail sdkws.MemberInvitedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		var userIDs []string
		for _, info := range detail.InvitedUserList {
			userIDs = append(userIDs, info.UserID)
		}

		if utils.IsContain(g.loginUserID, userIDs) {
			return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
		} else {
			return g.SyncGroupMembers(ctx, detail.Group.GroupID, userIDs...)
		}
	case constant.MemberEnterNotification: // 1510
		var detail sdkws.MemberEnterTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if detail.EntrantUser.UserID == g.loginUserID {
			if err := g.SyncGroups(ctx, detail.Group.GroupID); err != nil {
				return err
			}
			return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
		} else {
			return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.EntrantUser.UserID)
		}
	case constant.GroupDismissedNotification: // 1511
		var detail sdkws.GroupDismissedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.db.DeleteGroupAllMembers(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		if err := g.db.DeleteGroup(ctx, detail.Group.GroupID); err != nil {
			return err
		}
		return g.SyncAllGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberMutedNotification: // 1512
		var detail sdkws.GroupMemberMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.MutedUser.UserID)
	case constant.GroupMemberCancelMutedNotification: // 1513
		var detail sdkws.GroupMemberCancelMutedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.MutedUser.UserID)
	case constant.GroupMutedNotification: // 1514
		return g.SyncGroups(ctx, msg.GroupID)
	case constant.GroupCancelMutedNotification: // 1515
		return g.SyncGroups(ctx, msg.GroupID)
	case constant.GroupMemberInfoSetNotification: // 1516
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}

		return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID) //detail.ChangedUser.UserID
	case constant.GroupMemberSetToAdminNotification: // 1517
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID)
	case constant.GroupMemberSetToOrdinaryUserNotification: // 1518
		var detail sdkws.GroupMemberInfoSetTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMembers(ctx, detail.Group.GroupID, detail.ChangedUser.UserID)
	case 1519: // 1519  constant.GroupInfoSetAnnouncementNotification
		var detail sdkws.GroupInfoSetTips // sdkws.GroupInfoSetAnnouncementTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	case 1520: // 1520  constant.GroupInfoSetNameNotification
		var detail sdkws.GroupInfoSetTips // sdkws.GroupInfoSetNameTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroups(ctx, detail.Group.GroupID)
	default:
		return fmt.Errorf("unknown tips type: %d", msg.ContentType)
	}
}
