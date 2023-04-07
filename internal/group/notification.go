package group

import (
	comm "open_im_sdk/internal/common"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

func (g *Group) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
	if g.listener == nil {
		return
	}
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID, msg.ContentType)
	if msg.SendTime < g.loginTime || g.loginTime == 0 {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType, "msg.SendTime: ", msg.SendTime, "g.loginTime: ", g.loginTime)
		return
	}
	go func() {
		switch msg.ContentType {
		case constant.GroupCreatedNotification:
			g.groupCreatedNotification(msg, operationID)
		case constant.GroupInfoSetNotification:
			g.groupInfoSetNotification(msg, conversationCh, operationID)
		case constant.JoinGroupApplicationNotification:
			g.joinGroupApplicationNotification(msg, operationID)
		case constant.MemberQuitNotification:
			g.memberQuitNotification(msg, operationID) //1
		case constant.GroupApplicationAcceptedNotification:
			g.groupApplicationAcceptedNotification(msg, operationID)
		case constant.GroupApplicationRejectedNotification:
			g.groupApplicationRejectedNotification(msg, operationID)
		case constant.GroupOwnerTransferredNotification:
			g.groupOwnerTransferredNotification(msg, operationID) //ok
		case constant.MemberKickedNotification:
			g.memberKickedNotification(msg, operationID) //1
		case constant.MemberInvitedNotification:
			g.memberInvitedNotification(msg, operationID) //1
		case constant.MemberEnterNotification:
			g.memberEnterNotification(msg, operationID) //1
		case constant.GroupDismissedNotification:
			g.groupDismissNotification(msg, operationID)
		case constant.GroupMemberMutedNotification:
			g.groupMemberMuteChangedNotification(msg, false, operationID) //ok
		case constant.GroupMemberCancelMutedNotification:
			g.groupMemberMuteChangedNotification(msg, true, operationID) //ok
		case constant.GroupMutedNotification:
			fallthrough
		case constant.GroupCancelMutedNotification:
			g.groupMuteChangedNotification(msg, operationID)
		case constant.GroupMemberInfoSetNotification:
			g.groupMemberInfoSetNotification(msg, operationID) //ok
		case constant.GroupMemberSetToAdminNotification:
			g.groupMemberInfoSetNotification(msg, operationID) //ok
		case constant.GroupMemberSetToOrdinaryUserNotification:
			g.groupMemberInfoSetNotification(msg, operationID) //ok
		default:
			log.Error(operationID, "ContentType tip failed ", msg.ContentType)
		}
	}()
}

func (g *Group) groupCreatedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupCreatedTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, false)
}

func (g *Group) groupInfoSetNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupInfoSetTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.SyncJoinedGroupList(operationID) //todo,  sync some group info

}

func (g *Group) joinGroupApplicationNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.JoinGroupApplicationTips{Group: &api.GroupInfo{}, Applicant: &api.PublicUserInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.Applicant.UserID == g.loginUserID {
		g.SyncSelfGroupApplication(operationID)
	} else {
		g.SyncAdminGroupApplication(operationID)
	}
}

func (g *Group) memberQuitNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberQuitTips{Group: &api.GroupInfo{}, QuitUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.QuitUser.UserID == g.loginUserID {
		g.SyncJoinedGroupList(operationID)
		g.db.DeleteGroupAllMembers(detail.Group.GroupID)
	} else {
		log.Info(operationID, "deleteMemberImmediately ", detail.Group.GroupID, detail.QuitUser.UserID)
		g.deleteMemberImmediately(detail.Group.GroupID, detail.QuitUser.UserID, operationID)
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
		g.updateMemberCount(detail.Group.GroupID, operationID)
	}
}

func (g *Group) groupApplicationAcceptedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationAcceptedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg.String())
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "detail: ", msg.String())
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncAdminGroupApplication(operationID)
		return
	}
	if detail.ReceiverAs == 1 {
		g.SyncAdminGroupApplication(operationID)
		return
	}
	g.SyncSelfGroupApplication(operationID)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) groupApplicationRejectedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationRejectedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg.String())
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "detail: ", msg.String())
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncAdminGroupApplication(operationID)
		return
	}
	if detail.ReceiverAs == 1 {
		g.SyncAdminGroupApplication(operationID)
		return
	}
	g.SyncSelfGroupApplication(operationID)
}

func (g *Group) groupOwnerTransferredNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupOwnerTransferredTips{Group: &api.GroupInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	oldOwner, err := g.db.GetGroupMemberOwner(detail.Group.GroupID)
	if err == nil {
		g.updateLocalMemberImmediately(detail.Group.GroupID, oldOwner.UserID, map[string]interface{}{"role_level": constant.GroupOrdinaryUsers}, operationID)
	} else {
		log.Error(operationID, "GetGroupMemberOwner failed ", err.Error(), detail.Group.GroupID)
	}
	g.updateMemberImmediately(detail.NewGroupOwner, operationID)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
	g.SyncAdminGroupApplication(operationID)
}

func (g *Group) memberKickedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberKickedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}

	log.Info(operationID, "KickedUserList ", detail.KickedUserList)
	for _, v := range detail.KickedUserList {
		if v.UserID == g.loginUserID {
			g.SyncJoinedGroupList(operationID)
			g.db.DeleteGroupAllMembers(detail.Group.GroupID)
			return
		}
	}
	for _, v := range detail.KickedUserList {
		log.Info(operationID, "deleteMemberImmediately ", detail.Group.GroupID, v.UserID, operationID)
		g.deleteMemberImmediately(detail.Group.GroupID, v.UserID, operationID)
	}
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
	g.updateMemberCount(detail.Group.GroupID, operationID)
}

func (g *Group) memberInvitedNotification(msg *api.MsgData, operationID string) {
	log.NewWarn(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberInvitedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}

	log.Info(operationID, "detail InvitedUserList ", detail.InvitedUserList)
	for _, v := range detail.InvitedUserList {
		if v.UserID == g.loginUserID {
			g.SyncJoinedGroupList(operationID)
			g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, false)
			return
		}
	}
	for _, v := range detail.InvitedUserList {
		log.Info(operationID, "addMemberImmediately ", *v)
		g.addMemberImmediately(v, operationID)
	}
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
	g.updateMemberCount(detail.Group.GroupID, operationID)
}

func (g *Group) memberEnterNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberEnterTips{Group: &api.GroupInfo{}, EntrantUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.EntrantUser.UserID == g.loginUserID {
		g.SyncJoinedGroupList(operationID)
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, false)
	} else {
		log.Info(operationID, "addMemberImmediately ", *detail.EntrantUser)
		g.addMemberImmediately(detail.EntrantUser, operationID)
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
		g.updateMemberCount(detail.Group.GroupID, operationID)
	}
}

func (g *Group) groupDismissNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupDismissedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	g.SyncJoinedGroupList(operationID)
	g.db.DeleteGroupAllMembers(detail.Group.GroupID)
}

func (g *Group) groupMemberMuteChangedNotification(msg *api.MsgData, isCancel bool, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	var syncGroupID string
	MutedUser := &api.GroupMemberFullInfo{}
	if isCancel {
		detail := api.GroupMemberCancelMutedTips{Group: &api.GroupInfo{}}
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
			return
		}
		syncGroupID = detail.Group.GroupID
		MutedUser = detail.MutedUser
	} else {
		detail := api.GroupMemberMutedTips{Group: &api.GroupInfo{}}
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
			return
		}
		syncGroupID = detail.Group.GroupID
		MutedUser = detail.MutedUser
	}
	log.Info(operationID, "muted user ", *MutedUser)
	g.updateMemberImmediately(MutedUser, operationID)
	g.syncGroupMemberByGroupID(syncGroupID, operationID, true)
}

func (g *Group) groupMuteChangedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) groupMemberInfoSetNotification(msg *api.MsgData, operationID string) {
	detail := api.GroupMemberInfoSetTips{Group: &api.GroupInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID, msg.String(), "detail : ", detail.String())
	g.updateMemberImmediately(detail.ChangedUser, operationID)
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
}
