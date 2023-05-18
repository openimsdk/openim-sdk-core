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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/pkg/constant"
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
	case constant.GroupCreatedNotification:
		var detail sdkws.GroupCreatedTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupInfoSetNotification:
		var detail sdkws.GroupInfoSetTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncJoinedGroup(ctx)
	case constant.JoinGroupApplicationNotification:
		var detail sdkws.JoinGroupApplicationTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if detail.Applicant.UserID == g.loginUserID {
			return g.SyncSelfGroupApplication(ctx)
		} else {
			return g.SyncAdminGroupApplication(ctx)
		}
	case constant.MemberQuitNotification:
		var detail sdkws.MemberQuitTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if detail.QuitUser.UserID == g.loginUserID {
			if err := g.SyncJoinedGroup(ctx); err != nil {
				return err
			}
			return g.SyncGroupMember(ctx, detail.Group.GroupID)
		} else {
			return g.SyncGroupMember(ctx, detail.Group.GroupID)
		}
	case constant.GroupApplicationAcceptedNotification:
		var detail sdkws.GroupApplicationAcceptedTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			return g.SyncAdminGroupApplication(ctx)
		}
		if detail.ReceiverAs == 1 {
			return g.SyncAdminGroupApplication(ctx)
		}
		return g.SyncJoinedGroup(ctx)
	case constant.GroupApplicationRejectedNotification:
		var detail sdkws.GroupApplicationRejectedTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if detail.OpUser.UserID == g.loginUserID {
			return g.SyncAdminGroupApplication(ctx)
		}
		if detail.ReceiverAs == 1 {
			return g.SyncAdminGroupApplication(ctx)
		}
		return g.SyncSelfGroupApplication(ctx)
	case constant.GroupOwnerTransferredNotification:
		var detail sdkws.GroupOwnerTransferredTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberKickedNotification:
		var detail sdkws.MemberKickedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberInvitedNotification:
		var detail sdkws.MemberInvitedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberEnterNotification:
		var detail sdkws.MemberEnterTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupDismissedNotification:
		var detail sdkws.GroupDismissedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberMutedNotification:
		var detail sdkws.GroupMemberMutedTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberCancelMutedNotification:
		var detail sdkws.GroupMemberCancelMutedTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMutedNotification:
		return g.SyncJoinedGroup(ctx)
	case constant.GroupCancelMutedNotification:
		return g.SyncJoinedGroup(ctx)
	case constant.GroupMemberInfoSetNotification:
		var detail sdkws.GroupMemberInfoSetTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberSetToAdminNotification:
		var detail sdkws.GroupMemberInfoSetTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberSetToOrdinaryUserNotification:
		var detail sdkws.GroupMemberInfoSetTips
		if err := json.Unmarshal(msg.Content, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	default:
		return fmt.Errorf("unknown tips type: %d", msg.ContentType)
	}
}
