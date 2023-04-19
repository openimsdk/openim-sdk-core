package group

import (
	"context"
	"errors"
	"fmt"
	comm "open_im_sdk/internal/common"
	"open_im_sdk/pkg/constant"
	api "open_im_sdk/pkg/server_api_params"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

func (g *Group) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		if err := g.doNotification(ctx, msg); err != nil {
			// todo log
		}
	}()
}

func (g *Group) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if g.listener == nil {
		return errors.New("listener is nil")
	}
	if msg.SendTime < g.loginTime || g.loginTime == 0 {
		return errors.New("ignore notification")
	}
	var tips api.TipsComm
	if err := proto.Unmarshal(msg.Content, &tips); err != nil {
		return err
	}
	switch msg.ContentType {
	case constant.GroupCreatedNotification:
		var detail api.GroupCreatedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupInfoSetNotification:
		var detail api.GroupInfoSetTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		return g.SyncJoinedGroup(ctx)
	case constant.JoinGroupApplicationNotification:
		var detail api.JoinGroupApplicationTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if detail.Applicant.UserID == g.loginUserID {
			return g.SyncSelfGroupApplication(ctx)
		} else {
			return g.SyncAdminGroupApplication(ctx)
		}
	case constant.MemberQuitNotification:
		var detail api.MemberQuitTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
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
		var detail api.GroupApplicationAcceptedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
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
		var detail api.GroupApplicationRejectedTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
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
		var detail api.GroupOwnerTransferredTips
		if err := proto.Unmarshal(tips.Detail, &detail); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberKickedNotification:
		var detail api.MemberKickedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberInvitedNotification:
		var detail api.MemberInvitedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.MemberEnterNotification:
		var detail api.MemberEnterTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupDismissedNotification:
		var detail api.GroupDismissedTips
		if err := g.SyncJoinedGroup(ctx); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberMutedNotification:
		var detail api.GroupMemberMutedTips
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberCancelMutedNotification:
		var detail api.GroupMemberCancelMutedTips
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMutedNotification:
		return g.SyncJoinedGroup(ctx)
	case constant.GroupCancelMutedNotification:
		return g.SyncJoinedGroup(ctx)
	case constant.GroupMemberInfoSetNotification:
		var detail api.GroupMemberInfoSetTips
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberSetToAdminNotification:
		var detail api.GroupMemberInfoSetTips
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	case constant.GroupMemberSetToOrdinaryUserNotification:
		var detail api.GroupMemberInfoSetTips
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			return err
		}
		return g.SyncGroupMember(ctx, detail.Group.GroupID)
	default:
		return fmt.Errorf("unknown tips type: %d", msg.ContentType)
	}
}
