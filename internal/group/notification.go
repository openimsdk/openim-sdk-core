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
	"fmt"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
)

const (
	groupSortIDUnchanged = 0
	groupSortIDChanged   = 1
)

func (g *Group) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	if err := g.doNotification(ctx, msg); err != nil {
		log.ZError(ctx, "DoGroupNotification failed", err)
	}
}

func (g *Group) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	switch msg.ContentType {
	case constant.GroupApplicationAcceptedNotification: // 1505
		var detail sdkws.GroupApplicationAcceptedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if g.filter.ShouldExecute(detail.Uuid) {
			g.listener().OnGroupApplicationAccepted(utils.StructToJsonString(
				ServerGroupRequestToLocalGroupRequestForNotification(detail.GetGroup(), detail.GetRequest())))
		}

	case constant.GroupApplicationRejectedNotification: // 1506
		var detail sdkws.GroupApplicationRejectedTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if g.filter.ShouldExecute(detail.Uuid) {
			g.listener().OnGroupApplicationRejected(utils.StructToJsonString(
				ServerGroupRequestToLocalGroupRequestForNotification(detail.GetGroup(), detail.GetRequest())))
		}
	case constant.JoinGroupApplicationNotification: // 1503
		var detail sdkws.JoinGroupApplicationTips
		if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
			return err
		}
		if g.filter.ShouldExecute(detail.Uuid) {
			g.listener().OnGroupApplicationAdded(utils.StructToJsonString(
				ServerGroupRequestToLocalGroupRequestForNotification(detail.GetGroup(), detail.GetRequest())))
		}
	default:
		g.groupSyncMutex.Lock()
		defer g.groupSyncMutex.Unlock()
		switch msg.ContentType {
		case constant.GroupCreatedNotification: // 1501
			var detail sdkws.GroupCreatedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}

			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			if err := g.db.InsertGroupReadCursorState(ctx, &model_struct.LocalGroupReadCursorState{
				ConversationID: conversationID,
				CursorVersion:  1,
			}); err != nil {
				log.ZError(ctx, "InsertGroupReadCursorState on GroupCreatedNotification failed", err, "groupID", detail.Group.GroupID)
			}
			for _, member := range detail.MemberList {
				if err := g.db.InsertGroupReadCursor(ctx, &model_struct.LocalGroupReadCursor{
					ConversationID: conversationID,
					UserID:         member.UserID,
					MaxReadSeq:     0,
				}); err != nil {
					log.ZError(ctx, "InsertGroupReadCursor on GroupCreatedNotification failed", err, "groupID", detail.Group.GroupID, "userID", member.UserID)
				}
			}

			if err := g.IncrSyncJoinGroup(ctx); err != nil {
				return err
			}
			return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)

		case constant.GroupInfoSetNotification: // 1502
			var detail sdkws.GroupInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				nil, nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)

		case constant.MemberQuitNotification: // 1504
			var detail sdkws.MemberQuitTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			if detail.QuitUser.UserID == g.loginUserID {
				if err := g.db.DeleteGroupReadCursorsByConversationID(ctx, conversationID); err != nil {
					log.ZWarn(ctx, "DeleteGroupReadCursorsByConversationID err", err, "conversationID", conversationID)
				}
				if err := g.db.DeleteGroupReadCursorState(ctx, conversationID); err != nil {
					log.ZWarn(ctx, "DeleteGroupReadCursorState err", err, "conversationID", conversationID)
				}
				return g.IncrSyncJoinGroup(ctx)
			} else {
				if err := g.db.DeleteGroupReadCursor(ctx, conversationID, detail.QuitUser.UserID); err != nil {
					log.ZWarn(ctx, "DeleteGroupReadCursor err", err, "conversationID", conversationID, "userID", detail.QuitUser.UserID)
				} else {
					if err := g.db.IncrementGroupReadCursorVersion(ctx, conversationID); err != nil {
						log.ZWarn(ctx, "IncrementGroupReadCursorVersion err", err, "conversationID", conversationID)
					}
				}
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, []*sdkws.GroupMemberFullInfo{detail.QuitUser}, nil,
					nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.GroupOwnerTransferredNotification: // 1507
			var detail sdkws.GroupOwnerTransferredTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			if detail.Group == nil {
				return errs.New(fmt.Sprintf("group is nil, groupID: %s", detail.Group.GroupID)).Wrap()
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.NewGroupOwner, detail.OldGroupOwnerInfo}, nil,
				detail.Group, groupSortIDChanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
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
			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			if self {
				if err := g.db.DeleteGroupReadCursorsByConversationID(ctx, conversationID); err != nil {
					log.ZWarn(ctx, "DeleteGroupReadCursorsByConversationID err", err, "conversationID", conversationID)
				}
				if err := g.db.DeleteGroupReadCursorState(ctx, conversationID); err != nil {
					log.ZWarn(ctx, "DeleteGroupReadCursorState err", err, "conversationID", conversationID)
				}
				return g.IncrSyncJoinGroup(ctx)
			} else {
				deleted := false
				for _, info := range detail.KickedUserList {
					if err := g.db.DeleteGroupReadCursor(ctx, conversationID, info.UserID); err != nil {
						log.ZWarn(ctx, "DeleteGroupReadCursor err", err, "conversationID", conversationID, "userID", info.UserID)
					} else {
						deleted = true
					}
				}
				if deleted {
					if err := g.db.IncrementGroupReadCursorVersion(ctx, conversationID); err != nil {
						log.ZWarn(ctx, "IncrementGroupReadCursorVersion err", err, "conversationID", conversationID)
					}
				}
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, detail.KickedUserList, nil,
					nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.MemberInvitedNotification: // 1509
			var detail sdkws.MemberInvitedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			userIDMap := datautil.SliceSetAny(detail.InvitedUserList, func(e *sdkws.GroupMemberFullInfo) string {
				return e.UserID
			})

			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			for _, member := range detail.InvitedUserList {
				if err := g.db.InsertGroupReadCursor(ctx, &model_struct.LocalGroupReadCursor{
					ConversationID: conversationID,
					UserID:         member.UserID,
					MaxReadSeq:     0,
				}); err != nil {
					log.ZError(ctx, "InsertGroupReadCursor on MemberInvitedNotification failed", err, "groupID", detail.Group.GroupID, "userID", member.UserID)
				}
			}
			if err := g.db.IncrementGroupReadCursorVersion(ctx, conversationID); err != nil {
				log.ZError(ctx, "IncrementGroupReadCursorVersion err", err, "conversationID", conversationID)
			}

			//Also invited as a member
			if _, ok := userIDMap[g.loginUserID]; ok {
				if err := g.IncrSyncJoinGroup(ctx); err != nil {
					return err
				}
				return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
					detail.InvitedUserList, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.MemberEnterNotification: // 1510
			var detail sdkws.MemberEnterTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}

			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			if err := g.db.InsertGroupReadCursor(ctx, &model_struct.LocalGroupReadCursor{
				ConversationID: conversationID,
				UserID:         detail.EntrantUser.UserID,
				MaxReadSeq:     0,
			}); err != nil {
				log.ZError(ctx, "InsertGroupReadCursor on MemberEnterNotification failed", err, "groupID", detail.Group.GroupID, "userID", detail.EntrantUser.UserID)
			}
			if err := g.db.IncrementGroupReadCursorVersion(ctx, conversationID); err != nil {
				log.ZError(ctx, "IncrementGroupReadCursorVersion err", err, "conversationID", conversationID)
			}

			if detail.EntrantUser.UserID == g.loginUserID {
				if err := g.IncrSyncJoinGroup(ctx); err != nil {
					return err
				}
				return g.IncrSyncGroupAndMember(ctx, detail.Group.GroupID)
			} else {
				return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
					[]*sdkws.GroupMemberFullInfo{detail.EntrantUser}, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
			}
		case constant.GroupDismissedNotification: // 1511
			var detail sdkws.GroupDismissedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			g.listener().OnGroupDismissed(utils.StructToJsonString(detail.Group))

			conversationID := utils.GetConversationIDByGroupID(detail.Group.GroupID)
			if err := g.db.DeleteGroupReadCursorsByConversationID(ctx, conversationID); err != nil {
				log.ZWarn(ctx, "DeleteGroupReadCursorsByConversationID err", err, "conversationID", conversationID)
			}
			if err := g.db.DeleteGroupReadCursorState(ctx, conversationID); err != nil {
				log.ZWarn(ctx, "DeleteGroupReadCursorState err", err, "conversationID", conversationID)
			}

			return g.IncrSyncJoinGroup(ctx)
		case constant.GroupMemberMutedNotification: // 1512
			var detail sdkws.GroupMemberMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.MutedUser}, nil, nil,
				groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberCancelMutedNotification: // 1513
			var detail sdkws.GroupMemberCancelMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.MutedUser}, nil, nil,
				groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMutedNotification: // 1514
			var detail sdkws.GroupMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupCancelMutedNotification: // 1515
			var detail sdkws.GroupCancelMutedTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberInfoSetNotification: // 1516
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberSetToAdminNotification: // 1517
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupMemberSetToOrdinaryUserNotification: // 1518
			var detail sdkws.GroupMemberInfoSetTips
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil,
				[]*sdkws.GroupMemberFullInfo{detail.ChangedUser}, nil, nil,
				detail.GroupSortVersion, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupInfoSetAnnouncementNotification: // 1519
			var detail sdkws.GroupInfoSetAnnouncementTips //
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		case constant.GroupInfoSetNameNotification: // 1520
			var detail sdkws.GroupInfoSetNameTips //
			if err := utils.UnmarshalNotificationElem(msg.Content, &detail); err != nil {
				return err
			}
			return g.onlineSyncGroupAndMember(ctx, detail.Group.GroupID, nil, nil,
				nil, detail.Group, groupSortIDUnchanged, detail.GroupMemberVersion, detail.GroupMemberVersionID)
		default:
			return errs.New("unknown tips type", "contentType", msg.ContentType).Wrap()
		}
	}
	return nil
}
