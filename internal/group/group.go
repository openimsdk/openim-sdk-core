package group

import (
	"errors"
	comm "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"sync"

	"github.com/jinzhu/copier"
)

//	//utils.GetCurrentTimestampByMill()
type Group struct {
	listener    open_im_sdk_callback.OnGroupListener
	loginUserID string
	db          *db.DataBase
	p           *ws.PostApi
	loginTime   int64
}

func (g *Group) LoginTime() int64 {
	return g.loginTime
}

func (g *Group) SetLoginTime(loginTime int64) {
	g.loginTime = loginTime
}

func NewGroup(loginUserID string, db *db.DataBase, p *ws.PostApi) *Group {
	return &Group{loginUserID: loginUserID, db: db, p: p}
}

func (g *Group) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
	if g.listener == nil {
		return
	}
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	if msg.SendTime < g.loginTime {
		log.Warn(operationID, "ignore notification ", msg.ClientMsgID, msg.ServerMsgID, msg.Seq, msg.ContentType)
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
			g.memberQuitNotification(msg, operationID)
		case constant.GroupApplicationAcceptedNotification:
			g.groupApplicationAcceptedNotification(msg, operationID)
		case constant.GroupApplicationRejectedNotification:
			g.groupApplicationRejectedNotification(msg, operationID)
		case constant.GroupOwnerTransferredNotification:
			g.groupOwnerTransferredNotification(msg, operationID)
		case constant.MemberKickedNotification:
			g.memberKickedNotification(msg, operationID)
		case constant.MemberInvitedNotification:
			g.memberInvitedNotification(msg, operationID)
		case constant.MemberEnterNotification:
			g.memberEnterNotification(msg, operationID)
		case constant.GroupDismissedNotification:
			g.groupDismissNotification(msg, operationID)
		case constant.GroupMemberMutedNotification:
			g.groupMemberMuteChangedNotification(msg, false, operationID)
		case constant.GroupMemberCancelMutedNotification:
			g.groupMemberMuteChangedNotification(msg, true, operationID)
		case constant.GroupMutedNotification:
			fallthrough
		case constant.GroupCancelMutedNotification:
			g.groupMuteChangedNotification(msg, operationID)
		case constant.GroupMemberInfoSetNotification:
			g.groupMemberInfoSetNotification(msg, operationID)
		case constant.GroupMemberSetToAdminNotification:
			g.groupMemberInfoSetNotification(msg, operationID)
		case constant.GroupMemberSetToOrdinaryUserNotification:
			g.groupMemberInfoSetNotification(msg, operationID)
		default:
			log.Error(operationID, "ContentType tip failed ", msg.ContentType)
		}
	}()
}

func (g *Group) groupCreatedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupCreatedTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, false)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) groupInfoSetNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupInfoSetTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.SyncJoinedGroupList(operationID) //todo,  sync some group info
	conversationID := utils.GetConversationIDBySessionType(detail.Group.GroupID, constant.GroupChatType)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UpdateFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: detail.Group.GroupID, SessionType: constant.GroupChatType}}, conversationCh)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, conversationCh)
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
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
		g.updateMemberCount(detail.Group.GroupID, operationID)
	}
}

func (g *Group) groupApplicationAcceptedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationAcceptedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncAdminGroupApplication(operationID)
	} else {
		g.SyncSelfGroupApplication(operationID)
		g.SyncJoinedGroupList(operationID)
	}
	//g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
}

func (g *Group) groupApplicationRejectedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationRejectedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncAdminGroupApplication(operationID)
	} else {
		g.SyncSelfGroupApplication(operationID)
	}
}

func (g *Group) groupOwnerTransferredNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupOwnerTransferredTips{Group: &api.GroupInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}
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
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
	g.updateMemberCount(detail.Group.GroupID, operationID)
}

func (g *Group) memberInvitedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberInvitedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
		return
	}

	for _, v := range detail.InvitedUserList {
		if v.UserID == g.loginUserID {
			g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, false)
			g.SyncJoinedGroupList(operationID)
			return
		}
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
	if isCancel {
		detail := api.GroupMemberCancelMutedTips{Group: &api.GroupInfo{}}
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
			return
		}
		syncGroupID = detail.Group.GroupID
	} else {
		detail := api.GroupMemberMutedTips{Group: &api.GroupInfo{}}
		if err := comm.UnmarshalTips(msg, &detail); err != nil {
			log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
			return
		}
		syncGroupID = detail.Group.GroupID
	}
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
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID, true)
	_ = g.db.UpdateMsgSenderFaceURLAndSenderNickname(detail.ChangedUser.UserID, detail.ChangedUser.FaceURL, detail.ChangedUser.Nickname, constant.GroupChatType)
}

func (g *Group) createGroup(callback open_im_sdk_callback.Base, group sdk.CreateGroupBaseInfoParam,
	memberList sdk.CreateGroupMemberRoleParam, operationID string) *sdk.CreateGroupCallback {
	apiReq := api.CreateGroupReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = g.loginUserID
	apiReq.MemberList = memberList
	for _, v := range apiReq.MemberList {
		if v.RoleLevel == 0 {
			v.RoleLevel = 1
		}
	}
	copier.Copy(&apiReq, &group)
	realData := api.CreateGroupResp{}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "api req args: ", apiReq)
	g.p.PostFatalCallback(callback, constant.CreateGroupRouter, apiReq, &realData.GroupInfo, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(realData.GroupInfo.GroupID, operationID, false)
	var temp sdk.CreateGroupCallback
	temp = sdk.CreateGroupCallback(realData.GroupInfo)
	return &temp
}

func (g *Group) joinGroup(groupID, reqMsg string, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.JoinGroupReq{}
	apiReq.OperationID = operationID
	apiReq.ReqMessage = reqMsg
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.JoinGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncSelfGroupApplication(operationID)
}

func (g *Group) quitGroup(groupID string, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.QuitGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.QuitGroupRouter, apiReq, nil, apiReq.OperationID)
	//	g.syncGroupMemberByGroupID(groupID, operationID, false) //todo
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) dismissGroup(groupID string, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.DismissGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.DismissGroupRouter, apiReq, nil, apiReq.OperationID)
	//g.SyncJoinedGroupList(operationID)
}

func (g *Group) changeGroupMute(groupID string, isMute bool, callback open_im_sdk_callback.Base, operationID string) {
	if isMute {
		apiReq := api.MuteGroupReq{}
		apiReq.OperationID = operationID
		apiReq.GroupID = groupID
		g.p.PostFatalCallback(callback, constant.MuteGroupRouter, apiReq, nil, apiReq.OperationID)
	} else {
		apiReq := api.CancelMuteGroupReq{}
		apiReq.OperationID = operationID
		apiReq.GroupID = groupID
		g.p.PostFatalCallback(callback, constant.CancelMuteGroupRouter, apiReq, nil, apiReq.OperationID)
	}
}

func (g *Group) changeGroupMemberMute(groupID, userID string, mutedSeconds uint32, callback open_im_sdk_callback.Base, operationID string) {
	if mutedSeconds == 0 {
		apiReq := api.CancelMuteGroupMemberReq{}
		apiReq.OperationID = operationID
		apiReq.GroupID = groupID
		apiReq.UserID = userID
		g.p.PostFatalCallback(callback, constant.CancelMuteGroupMemberRouter, apiReq, nil, apiReq.OperationID)
	} else {
		apiReq := api.MuteGroupMemberReq{}
		apiReq.OperationID = operationID
		apiReq.GroupID = groupID
		apiReq.UserID = userID
		apiReq.MutedSeconds = mutedSeconds
		g.p.PostFatalCallback(callback, constant.MuteGroupMemberRouter, apiReq, nil, apiReq.OperationID)
	}
}

func (g *Group) setGroupMemberRoleLevel(callback open_im_sdk_callback.Base, groupID, userID string, roleLevel int, operationID string) {
	apiReq := api.SetGroupMemberRoleLevelReq{
		SetGroupMemberInfoReq: api.SetGroupMemberInfoReq{
			OperationID: operationID,
			UserID:      userID,
			GroupID:     groupID,
		},
		RoleLevel: roleLevel,
	}
	g.p.PostFatalCallback(callback, constant.SetGroupMemberInfoRouter, apiReq, nil, apiReq.OperationID)
	//g.syncGroupMemberByGroupID(groupID, operationID, true)
}

func (g *Group) getJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) sdk.GetJoinedGroupListCallback {
	groupList, err := g.db.GetJoinedGroupList()
	log.Info("this is rpc", groupList)
	common.CheckDBErrCallback(callback, err, operationID)
	return groupList
}

func (g *Group) GetGroupInfoFromLocal2Svr(groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(groupID)
	if err == nil {
		return localGroup, nil
	}
	groupIDList := []string{groupID}
	operationID := utils.OperationIDGenerator()
	svrGroup, err := g.getGroupsInfoFromSvr(groupIDList, operationID)
	if err == nil && len(svrGroup) == 1 {
		transfer := common.TransferToLocalGroupInfo(svrGroup)
		return transfer[0], nil
	}
	if err != nil {
		return nil, utils.Wrap(err, "")
	} else {
		return nil, utils.Wrap(errors.New("no group"), "")
	}
}

func (g *Group) searchGroups(callback open_im_sdk_callback.Base, param sdk.SearchGroupsParam, operationID string) sdk.SearchGroupsCallback {
	if len(param.KeywordList) == 0 || (!param.IsSearchGroupName && !param.IsSearchGroupID) {
		common.CheckAnyErrCallback(callback, 201, errors.New("keyword is null or search field all false"), operationID)
	}
	localGroup, err := g.db.GetAllGroupInfoByGroupIDOrGroupName(param.KeywordList[0], param.IsSearchGroupID, param.IsSearchGroupName)
	common.CheckDBErrCallback(callback, err, operationID)
	return localGroup
}

func (g *Group) getGroupsInfo(groupIDList sdk.GetGroupsInfoParam, callback open_im_sdk_callback.Base, operationID string) sdk.GetGroupsInfoCallback {
	groupList, err := g.db.GetJoinedGroupList()
	common.CheckDBErrCallback(callback, err, operationID)
	var result sdk.GetGroupsInfoCallback
	var notInDB []string

	for _, v := range groupList {
		in := false
		for _, k := range groupIDList {
			if v.GroupID == k {
				in = true
				break
			}
		}
		if in {
			result = append(result, v)
		}
	}

	for _, v := range groupIDList {
		in := false
		for _, k := range result {
			if v == k.GroupID {
				in = true
				break
			}
		}
		if !in {
			notInDB = append(notInDB, v)
		}
	}
	if len(notInDB) > 0 {
		groupsInfoSvr, err := g.getGroupsInfoFromSvr(notInDB, operationID)
		log.Info(operationID, "getGroupsInfoFromSvr groupsInfoSvr", groupsInfoSvr)
		common.CheckArgsErrCallback(callback, err, operationID)
		transfer := common.TransferToLocalGroupInfo(groupsInfoSvr)
		result = append(result, transfer...)
	}

	return result
}

func (g *Group) getGroupsInfoFromSvr(groupIDList []string, operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetGroupInfoReq{}
	apiReq.GroupIDList = groupIDList
	apiReq.OperationID = operationID
	var groupInfoList []*api.GroupInfo
	err := g.p.PostReturn(constant.GetGroupsInfoRouter, apiReq, &groupInfoList)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return groupInfoList, nil
}

func (g *Group) setGroupInfo(callback open_im_sdk_callback.Base, groupInfo sdk.SetGroupInfoParam, groupID, operationID string) {
	apiReq := api.SetGroupInfoReq{}
	apiReq.GroupName = groupInfo.GroupName
	apiReq.FaceURL = groupInfo.FaceURL
	apiReq.Notification = groupInfo.Notification
	apiReq.Introduction = groupInfo.Introduction
	apiReq.Ex = groupInfo.Ex
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	apiReq.NeedVerification = groupInfo.NeedVerification
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", groupInfo, groupID)
	g.p.PostFatalCallback(callback, constant.SetGroupInfoRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
}

//todo
func (g *Group) getGroupMemberList(callback open_im_sdk_callback.Base, groupID string, filter, offset, count int32, operationID string) sdk.GetGroupMemberListCallback {
	groupInfoList, err := g.db.GetGroupMemberListSplit(groupID, filter, int(offset), int(count))
	common.CheckDBErrCallback(callback, err, operationID)
	return groupInfoList
}

func (g *Group) getGroupMemberListByJoinTimeFilter(callback open_im_sdk_callback.Base, groupID string, offset, count int32, joinTimeBegin, joinTimeEnd int64, userIDList []string, operationID string) sdk.GetGroupMemberListCallback {
	if joinTimeEnd == 0 {
		joinTimeEnd = utils.GetCurrentTimestampBySecond()
	}
	groupInfoList, err := g.db.GetGroupMemberListSplitByJoinTimeFilter(groupID, int(offset), int(count), joinTimeBegin, joinTimeEnd, userIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return groupInfoList
}

//todo
func (g *Group) getGroupMembersInfo(callback open_im_sdk_callback.Base, groupID string, userIDList sdk.GetGroupMembersInfoParam, operationID string) sdk.GetGroupMembersInfoCallback {
	groupInfoList, err := g.db.GetGroupSomeMemberInfo(groupID, userIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return groupInfoList
}

func (g *Group) kickGroupMember(callback open_im_sdk_callback.Base, groupID string, memberList sdk.KickGroupMemberParam, reason string, operationID string) sdk.KickGroupMemberCallback {
	apiReq := api.KickGroupMemberReq{}
	apiReq.GroupID = groupID
	apiReq.KickedUserIDList = memberList
	apiReq.Reason = reason
	apiReq.OperationID = operationID
	realData := api.KickGroupMemberResp{}
	g.p.PostFatalCallback(callback, constant.KickGroupMemberRouter, apiReq, &realData.UserIDResultList, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, true)
	return realData.UserIDResultList
}

//
////1
func (g *Group) transferGroupOwner(callback open_im_sdk_callback.Base, groupID, newOwnerUserID string, operationID string) {
	apiReq := api.TransferGroupOwnerReq{}
	apiReq.GroupID = groupID
	apiReq.NewOwnerUserID = newOwnerUserID
	apiReq.OperationID = operationID
	apiReq.OldOwnerUserID = g.loginUserID
	g.p.PostFatalCallback(callback, constant.TransferGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupMember(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, true)
}

func (g *Group) inviteUserToGroup(callback open_im_sdk_callback.Base, groupID, reason string, userList sdk.InviteUserToGroupParam, operationID string) sdk.InviteUserToGroupCallback {
	apiReq := api.InviteUserToGroupReq{}
	apiReq.GroupID = groupID
	apiReq.Reason = reason
	apiReq.InvitedUserIDList = userList
	apiReq.OperationID = operationID
	var realData sdk.InviteUserToGroupCallback
	g.p.PostFatalCallback(callback, constant.InviteUserToGroupRouter, apiReq, &realData, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, false)
	return realData
}

//
////1
func (g *Group) getRecvGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetGroupApplicationListCallback {
	applicationList, err := g.db.GetAdminGroupApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	return applicationList
}

func (g *Group) getSendGroupApplicationList(callback open_im_sdk_callback.Base, operationID string) sdk.GetSendGroupApplicationListCallback {
	applicationList, err := g.db.GetSendGroupApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	return applicationList
}

func (g *Group) getRecvGroupApplicationListFromSvr(operationID string) ([]*api.GroupRequest, error) {
	apiReq := api.GetGroupApplicationListReq{}
	apiReq.FromUserID = g.loginUserID
	apiReq.OperationID = operationID
	var realData []*api.GroupRequest
	err := g.p.PostReturn(constant.GetRecvGroupApplicationListRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (g *Group) getSendGroupApplicationListFromSvr(operationID string) ([]*api.GroupRequest, error) {
	apiReq := api.GetUserReqGroupApplicationListReq{}
	apiReq.UserID = g.loginUserID
	apiReq.OperationID = operationID
	var realData []*api.GroupRequest
	err := g.p.PostReturn(constant.GetSendGroupApplicationListRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (g *Group) processGroupApplication(callback open_im_sdk_callback.Base, groupID, fromUserID, handleMsg string, handleResult int32, operationID string) {
	apiReq := api.ApplicationGroupResponseReq{}
	apiReq.GroupID = groupID
	apiReq.OperationID = operationID
	apiReq.FromUserID = fromUserID
	apiReq.HandleResult = handleResult
	apiReq.HandledMsg = handleMsg
	if handleResult == constant.GroupResponseAgree {
		g.p.PostFatalCallback(callback, constant.AcceptGroupApplicationRouter, apiReq, nil, apiReq.OperationID)
		g.syncGroupMemberByGroupID(groupID, operationID, true)
	} else if handleResult == constant.GroupResponseRefuse {
		g.p.PostFatalCallback(callback, constant.RefuseGroupApplicationRouter, apiReq, nil, apiReq.OperationID)
	}
	g.SyncAdminGroupApplication(operationID)
}

func (g *Group) getJoinedGroupListFromSvr(operationID string) ([]*api.GroupInfo, error) {
	apiReq := api.GetJoinedGroupListReq{}
	apiReq.OperationID = operationID
	apiReq.FromUserID = g.loginUserID
	var result []*api.GroupInfo
	log.Debug(operationID, "api args: ", apiReq)
	err := g.p.PostReturn(constant.GetJoinedGroupListRouter, apiReq, &result)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return result, nil
}

func (g *Group) updateMemberCount(groupID string, operationID string) {
	memberCount, err := g.db.GetGroupMemberCount(groupID)
	if err != nil {
		log.Error(operationID, "GetGroupMemberCount failed ", err.Error(), groupID)
		return
	}
	groupInfo, err := g.db.GetGroupInfoByGroupID(groupID)
	if err != nil {
		log.Error(operationID, "GetGroupInfoByGroupID failed ", err.Error(), groupID)
		return
	}
	if groupInfo.MemberCount != int32(memberCount) {
		groupInfo.MemberCount = int32(memberCount)
		log.Info(operationID, "OnGroupInfoChanged , update group info", groupInfo)
		g.db.UpdateGroup(groupInfo)
		g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupInfo))
	}
}

func (g *Group) getGroupMembersInfoFromSvr(groupID string, memberList []string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupMembersInfoReq
	apiReq.OperationID = utils.OperationIDGenerator()
	apiReq.GroupID = groupID
	apiReq.MemberList = memberList
	var realData []*api.GroupMemberFullInfo
	err := g.p.PostReturn(constant.GetGroupMembersInfoRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (g *Group) SyncSelfGroupApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getSendGroupApplicationListFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getSendGroupApplicationListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalSendGroupRequest(svrList)
	onLocal, err := g.db.GetSendGroupApplication()
	if err != nil {
		log.NewError(operationID, "GetSendGroupApplication failed ", err.Error())
		return
	}

	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupRequestDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupRequest failed ", err.Error(), *onServer[index])
			continue
		}
		callbackData := *onServer[index]
		g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == constant.GroupResponseRefuse {
			callbackData := *onServer[index]
			g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))

		} else if onServer[index].HandleResult == constant.GroupResponseAgree {
			callbackData := *onServer[index]
			g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
		} else {
			callbackData := *onServer[index]
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))

		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := *onLocal[index]
		g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnGroupApplicationDeleted", utils.StructToJsonString(callbackData))
	}
}

func (g *Group) SyncAdminGroupApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getRecvGroupApplicationListFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getRecvGroupApplicationListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalAdminGroupRequest(svrList)
	onLocal, err := g.db.GetAdminGroupApplication()
	if err != nil {
		log.NewError(operationID, "GetAdminGroupApplication failed ", err.Error())
		return
	}

	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckAdminGroupRequestDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertAdminGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupRequest failed ", err.Error(), *onServer[index])
			continue
		}
		callbackData := sdk.GroupApplicationAddedCallback(*onServer[index])
		g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnReceiveJoinGroupApplicationAdded", utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateAdminGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == constant.GroupResponseRefuse {
			callbackData := sdk.GroupApplicationRejectCallback(*onServer[index])
			g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))

		} else if onServer[index].HandleResult == constant.GroupResponseAgree {
			callbackData := sdk.GroupApplicationAcceptCallback(*onServer[index])
			g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
		} else {
			callbackData := sdk.GroupApplicationAcceptCallback(*onServer[index])
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveJoinGroupApplicationAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteAdminGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.GroupApplicationDeletedCallback(*onLocal[index])
		g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnReceiveJoinGroupApplicationDeleted", utils.StructToJsonString(callbackData))
	}
}

//func transferGroupInfo(input []*api.GroupInfo) []*api.GroupInfo{
//	var result []*api.GroupInfo
//	for _, v := range input {
//		t := &api.GroupInfo{}
//		copier.Copy(t, &v)
//		if v.NeedVerification != nil {
//			t.NeedVerification = v.NeedVerification.Value
//		}
//		result = append(result, t)
//	}
//	return result
//}
func (g *Group) SyncJoinedGroupList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getJoinedGroupListFromSvr(operationID)
	log.Info(operationID, "getJoinedGroupListFromSvr", svrList, g.loginUserID)
	if err != nil {
		log.NewError(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupInfo(svrList)
	onLocal, err := g.db.GetJoinedGroupList()
	if err != nil {
		log.NewError(operationID, "GetJoinedGroupList failed ", err.Error())
		return
	}

	log.NewInfo(operationID, " onLocal ", onLocal, g.loginUserID)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupInfoDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB, g.loginUserID)
	for _, index := range aInBNot {
		err := g.db.InsertGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroup failed ", err.Error(), *onServer[index])
			continue
		}

		callbackData := sdk.JoinedGroupAddedCallback(*onServer[index])
		g.listener.OnJoinedGroupAdded(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnJoinedGroupAdded", utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroup failed ", err.Error(), onServer[index])
			continue
		}
		callbackData := sdk.GroupInfoChangedCallback(*onServer[index])
		g.listener.OnGroupInfoChanged(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnGroupInfoChanged", utils.StructToJsonString(callbackData))
	}

	for _, index := range bInANot {
		log.Info(operationID, "DeleteGroup: ", onLocal[index].GroupID, g.loginUserID)
		err := g.db.DeleteGroup(onLocal[index].GroupID)
		if err != nil {
			log.NewError(operationID, "DeleteGroup failed ", err.Error(), onLocal[index].GroupID)
			continue
		}
		g.db.DeleteGroupAllMembers(onLocal[index].GroupID)
		callbackData := sdk.JoinedGroupDeletedCallback(*onLocal[index])
		g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnJoinedGroupDeleted", utils.StructToJsonString(callbackData))
	}
}

func (g *Group) syncGroupMemberByGroupID(groupID string, operationID string, onGroupMemberNotification bool) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", groupID)
	svrList, err := g.getGroupAllMemberByGroupIDFromSvr(groupID, operationID)
	if err != nil {
		log.NewError(operationID, "getGroupAllMemberByGroupIDFromSvr failed ", err.Error(), groupID)
		return
	}
	onServer := common.TransferToLocalGroupMember(svrList)
	onLocal, err := g.db.GetGroupMemberListByGroupID(groupID)
	if err != nil {
		log.NewError(operationID, "GetGroupMemberListByGroupID failed ", err.Error(), groupID)
		return
	}
	//log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, _ := common.CheckGroupMemberDiff(onServer, onLocal)
	//log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupMember failed ", err.Error(), *onServer[index])
			continue
		}
		if onGroupMemberNotification == true {
			callbackData := sdk.GroupMemberAddedCallback(*onServer[index])
			g.listener.OnGroupMemberAdded(utils.StructToJsonString(callbackData))
			log.Debug(operationID, "OnGroupMemberAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupMember failed ", err.Error(), onServer[index])
			continue
		}

		callbackData := sdk.GroupMemberInfoChangedCallback(*onServer[index])
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(callbackData))
		log.Info(operationID, "OnGroupMemberInfoChanged", utils.StructToJsonString(callbackData))
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupMember(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupMember failed ", err.Error(), onLocal[index].GroupID, onLocal[index].UserID)
			continue
		}
		if onGroupMemberNotification == true {
			callbackData := sdk.GroupMemberDeletedCallback(*onLocal[index])
			g.listener.OnGroupMemberDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupMemberDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

func (g *Group) SyncJoinedGroupMember(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	groupListOnServer, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		log.Error(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	var wg sync.WaitGroup
	if len(groupListOnServer) == 0 {
		return
	}
	wg.Add(len(groupListOnServer))
	log.Info(operationID, "syncGroupMemberByGroupID begin", len(groupListOnServer))
	for _, v := range groupListOnServer {
		go func(groupID, operationID string) {
			g.syncGroupMemberByGroupID(groupID, operationID, true)
			wg.Done()
		}(v.GroupID, operationID)
	}

	wg.Wait()
	log.Info(operationID, "syncGroupMemberByGroupID end")
}

func (g *Group) SyncJoinedGroupMemberForFirstLogin(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	groupListOnServer, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		log.Error(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	var wg sync.WaitGroup
	if len(groupListOnServer) == 0 {
		return
	}
	wg.Add(len(groupListOnServer))
	log.Info(operationID, "syncGroupMemberByGroupID begin", len(groupListOnServer))
	for _, v := range groupListOnServer {
		go func(groupID, operationID string) {
			g.syncGroupMemberByGroupID(groupID, operationID, false)
			wg.Done()
		}(v.GroupID, operationID)
	}

	wg.Wait()
	log.Info(operationID, "syncGroupMemberByGroupID end")
}

func (g *Group) getGroupAllMemberByGroupIDFromSvr(groupID string, operationID string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupAllMemberReq
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	var realData []*api.GroupMemberFullInfo
	err := g.p.PostReturn(constant.GetGroupAllMemberListRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (g *Group) setGroupMemberNickname(callback open_im_sdk_callback.Base, groupID, userID string, GroupMemberNickname string, operationID string) {
	var apiReq api.SetGroupMemberNicknameReq
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	apiReq.UserID = userID
	apiReq.Nickname = GroupMemberNickname
	g.p.PostFatalCallback(callback, constant.SetGroupMemberNicknameRouter, apiReq, nil, apiReq.OperationID)
	g.syncGroupMemberByGroupID(groupID, operationID, true)
}
