package group

import (
	comm "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)

type OnGroupListener interface {
	OnJoinedGroupAdded(groupInfo string)
	OnJoinedGroupDeleted(groupInfo string)

	OnGroupMemberAdded(groupMemberInfo string)
	OnGroupMemberDeleted(groupMemberInfo string)

	OnReceiveJoinGroupApplicationAdded(groupApplication string)
	OnReceiveJoinGroupApplicationDeleted(groupApplication string)

	OnGroupApplicationAdded(groupApplication string)
	OnGroupApplicationDeleted(groupApplication string)

	OnGroupInfoChanged(groupInfo string)
	OnGroupMemberInfoChanged(groupMemberInfo string)

	OnGroupApplicationAccepted(groupApplication string)
	OnGroupApplicationRejected(groupApplication string)
}

type Group struct {
	listener    OnGroupListener
	loginUserID string
	db          *db.DataBase
	p           *ws.PostApi
}

func NewGroup(loginUserID string, db *db.DataBase, p *ws.PostApi) *Group {
	return &Group{loginUserID: loginUserID, db: db, p: p}
}

func (g *Group) DoNotification(msg *api.MsgData) {
	if g.listener == nil {
		return
	}
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	go func() {
		switch msg.ContentType {
		case constant.GroupCreatedNotification:
			g.groupCreatedNotification(msg, operationID)
		case constant.GroupInfoSetNotification:
			g.groupInfoSetNotification(msg, operationID)
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
		default:
			log.Error(operationID, "ContentType tip failed ", msg.ContentType)
		}
	}()
}

func (g *Group) groupCreatedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) groupInfoSetNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupInfoSetTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.SyncJoinedGroupList(operationID) //todo

}

func (g *Group) joinGroupApplicationNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.JoinGroupApplicationTips{Group: &api.GroupInfo{}, Applicant: &api.PublicUserInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}
	if detail.Applicant.UserID == g.loginUserID {
		g.SyncSelfGroupApplication(operationID)
	} else {
		g.SyncGroupApplication(operationID)
	}
}

func (g *Group) memberQuitNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberQuitTips{Group: &api.GroupInfo{}, QuitUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}
	if detail.QuitUser.UserID == g.loginUserID {
		g.SyncJoinedGroupList(operationID)
	} else {
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
	}
}

func (g *Group) groupApplicationAcceptedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationAcceptedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncGroupApplication(operationID)
	} else {
		g.SyncSelfGroupApplication(operationID)
	}
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
}

func (g *Group) groupApplicationRejectedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupApplicationRejectedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}
	if detail.OpUser.UserID == g.loginUserID {
		g.SyncGroupApplication(operationID)
	} else {
		g.SyncSelfGroupApplication(operationID)
	}
}

func (g *Group) groupOwnerTransferredNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.GroupOwnerTransferredTips{Group: &api.GroupInfo{}}
	comm.UnmarshalTips(msg, &detail)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) memberKickedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberKickedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}

	for _, v := range detail.KickedUserList {
		if v.UserID == g.loginUserID {
			g.SyncJoinedGroupList(operationID)
			return
		}
	}
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
}

func (g *Group) memberInvitedNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberInvitedTips{Group: &api.GroupInfo{}, OpUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}

	for _, v := range detail.InvitedUserList {
		if v.UserID == g.loginUserID {
			g.SyncJoinedGroupList(operationID)
			return
		}
	}
	g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
}

func (g *Group) memberEnterNotification(msg *api.MsgData, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", msg.ClientMsgID, msg.ServerMsgID)
	detail := api.MemberEnterTips{Group: &api.GroupInfo{}, EntrantUser: &api.GroupMemberFullInfo{}}
	if err := comm.UnmarshalTips(msg, &detail); err != nil {
		log.Error(operationID, "UnmarshalTips failed ", err.Error(), msg)
	}
	if detail.EntrantUser.UserID == g.loginUserID {
		g.SyncJoinedGroupList(operationID)
	} else {
		g.syncGroupMemberByGroupID(detail.Group.GroupID, operationID)
	}
}

func (g *Group) createGroup(callback common.Base, group sdk.CreateGroupBaseInfoParam,
	memberList sdk.CreateGroupMemberRoleParam, operationID string) *sdk.CreateGroupCallback {
	apiReq := api.CreateGroupReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = g.loginUserID
	apiReq.GroupName = group.GroupName
	apiReq.GroupType = group.GroupType
	apiReq.MemberList = memberList
	realData := api.CreateGroupResp{}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "api req args: ", apiReq)
	g.p.PostFatalCallback(callback, constant.CreateGroupRouter, apiReq, &realData.GroupInfo, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(realData.GroupInfo.GroupID, operationID)
	var temp sdk.CreateGroupCallback
	temp = sdk.CreateGroupCallback(realData.GroupInfo)
	return &temp
}

func (g *Group) joinGroup(groupID, reqMsg string, callback common.Base, operationID string) {
	apiReq := api.JoinGroupReq{}
	apiReq.OperationID = operationID
	apiReq.ReqMessage = reqMsg
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.JoinGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) quitGroup(groupID string, callback common.Base, operationID string) {
	apiReq := api.QuitGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.QuitGroupRouter, apiReq, nil, apiReq.OperationID)
	g.syncGroupMemberByGroupID(groupID, operationID) //todo
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) getJoinedGroupList(callback common.Base, operationID string) sdk.GetJoinedGroupListCallback {
	groupList, err := g.db.GetJoinedGroupList()
	common.CheckDBErrCallback(callback, err, operationID)
	return groupList
}

func (g *Group) getGroupsInfo(groupIdList sdk.GetGroupsInfoParam, callback common.Base, operationID string) sdk.GetGroupsInfoCallback {
	groupList, err := g.db.GetJoinedGroupList()
	common.CheckDBErrCallback(callback, err, operationID)
	var result sdk.GetGroupsInfoCallback
	for _, v := range groupList {
		in := false
		for _, k := range groupIdList {
			if v.GroupID == k {
				in = true
				break
			}
		}
		if in {
			result = append(result, v)
		}
	}
	return result
}

func (g *Group) setGroupInfo(callback common.Base, groupInfo sdk.SetGroupInfoParam, groupID, operationID string) {
	apiReq := api.SetGroupInfoReq{}
	apiReq.GroupName = groupInfo.GroupName
	apiReq.FaceURL = groupInfo.FaceUrl
	apiReq.Notification = groupInfo.Notification
	apiReq.Introduction = groupInfo.Introduction
	apiReq.Ex = groupInfo.Ex
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.SetGroupInfoRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
}

//todo
func (g *Group) getGroupMemberList(callback common.Base, groupID string, filter int32, next int32, operationID string) sdk.GetGroupMemberListCallback {
	groupInfoList, err := g.db.GetGroupMemberListByGroupID(groupID)
	common.CheckDBErrCallback(callback, err, operationID)
	return sdk.GetGroupMemberListCallback{MemberList: groupInfoList, NextSeq: 0}
}

//todo
func (g *Group) getGroupMembersInfo(callback common.Base, groupID string, userIDList sdk.GetGroupMembersInfoParam, operationID string) sdk.GetGroupMembersInfoCallback {
	groupInfoList, err := g.db.GetGroupSomeMemberInfo(groupID, userIDList)
	common.CheckDBErrCallback(callback, err, operationID)
	return groupInfoList
}

func (g *Group) kickGroupMember(callback common.Base, groupID string, memberList sdk.KickGroupMemberParam, reason string, operationID string) sdk.KickGroupMemberCallback {
	apiReq := api.KickGroupMemberReq{}
	apiReq.GroupID = groupID
	apiReq.KickedUserIDList = memberList
	apiReq.Reason = reason
	apiReq.OperationID = operationID
	realData := api.KickGroupMemberResp{}
	g.p.PostFatalCallback(callback, constant.KickGroupMemberRouter, apiReq, &realData.UserIDResultList, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
	return realData.UserIDResultList
}

//
////1
func (g *Group) transferGroupOwner(callback common.Base, groupID, newOwnerUserID string, operationID string) {
	apiReq := api.TransferGroupOwnerReq{}
	apiReq.GroupID = groupID
	apiReq.NewOwnerUserID = newOwnerUserID
	apiReq.OperationID = operationID
	apiReq.OldOwnerUserID = g.loginUserID
	g.p.PostFatalCallback(callback, constant.TransferGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupMember(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID)
}

func (g *Group) inviteUserToGroup(callback common.Base, groupID, reason string, userList sdk.InviteUserToGroupParam, operationID string) sdk.InviteUserToGroupCallback {
	apiReq := api.InviteUserToGroupReq{}
	apiReq.GroupID = groupID
	apiReq.Reason = reason
	apiReq.InvitedUserIDList = userList
	apiReq.OperationID = operationID
	var realData sdk.InviteUserToGroupCallback
	g.p.PostFatalCallback(callback, constant.InviteUserToGroupRouter, apiReq, &realData, apiReq.OperationID)
	g.SyncJoinedGroupMember(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID)
	return realData
}

//
////1
func (g *Group) getGroupApplicationList(callback common.Base, operationID string) sdk.GetGroupApplicationListCallback {
	applicationList, err := g.db.GetRecvGroupApplication()
	common.CheckDBErrCallback(callback, err, operationID)
	return applicationList
}

func (g *Group) getGroupApplicationListFromSvr(operationID string) ([]*api.GroupRequest, error) {
	apiReq := api.GetGroupApplicationListReq{}
	apiReq.FromUserID = g.loginUserID
	apiReq.OperationID = operationID
	var realData []*api.GroupRequest
	err := g.p.PostReturn(constant.GetGroupApplicationListRouter, apiReq, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}

func (g *Group) processGroupApplication(callback common.Base, groupID, fromUserID, handleMsg string, handleResult int32, operationID string) {
	apiReq := api.ApplicationGroupResponseReq{}
	apiReq.GroupID = groupID
	apiReq.OperationID = operationID
	apiReq.FromUserID = fromUserID
	apiReq.HandleResult = handleResult
	apiReq.HandledMsg = handleMsg
	if handleResult == 1 {
		g.p.PostFatalCallback(callback, constant.AcceptGroupApplicationRouter, apiReq, nil, apiReq.OperationID)
		g.syncGroupMemberByGroupID(groupID, operationID)
	} else if handleResult == -1 {
		g.p.PostFatalCallback(callback, constant.RefuseGroupApplicationRouter, apiReq, nil, apiReq.OperationID)
	}
	g.SyncGroupApplication(operationID)
}

//
//func (u *Group) GroupApplicationProcessedCallback(node open_im_sdk.updateGroupNode, process int32) {
//	args := node.Args.(open_im_sdk.applyGroupProcessedArgs)
//	list := make([]open_im_sdk.groupMemberFullInfo, 0)
//	for _, v := range args.applyList {
//		list = append(list, v.member)
//	}
//	jsonApplyList, err := json.Marshal(list)
//	if err != nil {
//		return
//	}
//
//	processed := node.Args.(open_im_sdk.applyGroupProcessedArgs) //receiver : all group member
//	var flag = 0
//	var idx = 0
//	for i, v := range processed.applyList {
//		if v.member.UserId == u.loginUserID {
//			flag = 1
//			idx = i
//			break
//		}
//	}
//
//	if flag == 1 {
//		jsonOp, err := json.Marshal(processed.op)
//		if err != nil {
//			return
//		}
//		u.listener.OnApplicationProcessed(node.groupId, string(jsonOp), process, processed.applyList[idx].reason)
//	}
//
//	if process == 1 {
//		jsonOp, err := json.Marshal(processed.op)
//		if err != nil {
//			return
//		}
//		u.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonApplyList))
//	}
//}

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

//
//func (u *Group) getGroupMembersInfoFromLocal(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
//	var result []open_im_sdk.groupMemberFullInfo
//	localMemberList, err := u.getLocalGroupMemberListByGroupID(groupId)
//	if err != nil {
//		return nil, err
//	}
//	for _, i := range localMemberList {
//		for _, j := range memberList {
//			if i.UserId == j {
//				result = append(result, i)
//			}
//		}
//	}
//	return result, nil
//}

//func (u *Group) getGroupMembersInfoTry2(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
//	result, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
//	if err != nil || len(result) == 0 {
//		return u.getGroupMembersInfoFromSvr(groupId, memberList)
//	} else {
//		return result, err
//	}
//}

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

//todo
func (g *Group) SyncSelfGroupApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")

}

func (g *Group) SyncGroupApplication(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getGroupApplicationListFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getGroupApplicationListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupRequest(svrList)
	onLocal, err := g.db.GetRecvGroupApplication()
	if err != nil {
		log.NewError(operationID, "GetJoinedGroupList failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupRequestDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.GroupApplicationAddedCallback(*onServer[index])
		g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == -1 {
			callbackData := sdk.GroupApplicationRejectCallback(*onServer[index])
			g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))

		} else if onServer[index].HandleResult == 1 {
			callbackData := sdk.GroupApplicationAcceptCallback(*onServer[index])
			g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupRequest(onServer[index].GroupID, onServer[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.GroupApplicationDeletedCallback(*onLocal[index])
		g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
	}
}

func (g *Group) SyncJoinedGroupList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		log.NewError(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupInfo(svrList)
	onLocal, err := g.db.GetJoinedGroupList()
	if err != nil {
		log.NewError(operationID, "GetRecvFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupInfoDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroup failed ", err.Error(), onServer[index])
			continue
		}
		callbackData := sdk.JoinedGroupAddedCallback(*onServer[index])
		g.listener.OnJoinedGroupAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroup failed ", err.Error(), onServer[index])
			continue
		}
		callbackData := sdk.GroupInfoChangedCallback(*onServer[index])
		g.listener.OnGroupInfoChanged(utils.StructToJsonString(callbackData))
	}

	for _, index := range bInANot {
		err := g.db.DeleteGroup(onLocal[index].GroupID)
		if err != nil {
			log.NewError(operationID, "DeleteGroup failed ", err.Error(), onLocal[index].GroupID)
			continue
		}
		callbackData := sdk.JoinedGroupDeletedCallback(*onLocal[index])
		g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(callbackData))
	}
}

func (g *Group) syncGroupMemberByGroupID(groupID string, operationID string) {
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
	log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupMemberDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupMember failed ", err.Error(), onServer[index])
			continue
		}
		callbackData := sdk.GroupMemberAddedCallback(*onServer[index])
		g.listener.OnGroupMemberAdded(utils.StructToJsonString(callbackData))
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupMember failed ", err.Error(), onServer[index])
			continue
		}
		callbackData := sdk.GroupMemberInfoChangedCallback(*onServer[index])
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(callbackData))
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupMember(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupMember failed ", err.Error(), onLocal[index].GroupID, onLocal[index].UserID)
			continue
		}
		callbackData := sdk.GroupMemberDeletedCallback(*onLocal[index])
		g.listener.OnGroupMemberDeleted(utils.StructToJsonString(callbackData))
	}
}

func (g *Group) SyncJoinedGroupMember(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	groupListOnServer, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		log.Error(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	for _, v := range groupListOnServer {
		g.syncGroupMemberByGroupID(v.GroupID, operationID)
	}
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

//func (u *Group) getLocalGroupMemberListNew() ([]open_im_sdk.groupMemberFullInfo, error) {
//	return u.getLocalGroupMemberList()
//}

//func (u *Group) getLocalGroupMemberListByGroupIDNew(groupId string) ([]open_im_sdk.groupMemberFullInfo, error) {
//	return u.getLocalGroupMemberListByGroupID(groupId)
//}
//func (u *Group) insertIntoLocalGroupMemberNew(info open_im_sdk.groupMemberFullInfo) error {
//	return u.insertIntoLocalGroupMember(info)
//}
//func (u *Group) delLocalGroupMemberNew(info open_im_sdk.groupMemberFullInfo) error {
//	return u.delLocalGroupMember(info)
//}
//func (u *Group) replaceLocalGroupMemberInfoNew(info open_im_sdk.groupMemberFullInfo) error {
//	return u.replaceLocalGroupMemberInfo(info)
//}
//
//func (u *Group) insertIntoSelfApplyToGroupRequestNew(groupId, message string) error {
//	return u.insertIntoSelfApplyToGroupRequest(groupId, message)
//}

/*
func (u *UserRelated) getLocalGroupsInfo1() ([]groupInfo, error) {
	localGroupsInfo, err := u.getLocalGroupsInfo()
	if err != nil {
		return nil, err
	}
	groupId2Owner := make(map[string]string)
	groupId2MemberNum := make(map[string]uint32)
	for index, v := range localGroupsInfo {
		if _, ok := groupId2Owner[v.GroupId]; !ok {
			ownerId, err := u.findLocalGroupOwnerByGroupId(v.GroupId)
			if err != nil {
				sdkLog(err.Error())
			}
			groupId2Owner[v.GroupId] = ownerId
		}
		localGroupsInfo[index].OwnerId = groupId2Owner[v.GroupId]
		if _, ok := groupId2MemberNum[v.GroupId]; !ok {
			num, err := u.getLocalGroupMemberNumByGroupId(v.GroupId)
			if err != nil {
				sdkLog(err.Error())
			}
			groupId2MemberNum[v.GroupId] = uint32(num)
		}
		localGroupsInfo[index].MemberCount = groupId2MemberNum[v.GroupId]
	}
	return localGroupsInfo, nil
}
*/

//
//func (u *Group) createGroupCallback(node open_im_sdk.updateGroupNode) {
//	// member list to json
//	jsonMemberList, err := json.Marshal(node.Args.(open_im_sdk.createGroupArgs).initMemberList)
//	if err != nil {
//		return
//	}
//	u.listener.OnMemberEnter(node.groupId, string(jsonMemberList))
//	u.listener.OnGroupCreated(node.groupId)
//}
//
//func (u *Group) joinGroupCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.joinGroupArgs)
//	jsonApplyUser, err := json.Marshal(args.applyUser)
//	if err != nil {
//		return
//	}
//	u.listener.OnReceiveJoinApplication(node.groupId, string(jsonApplyUser), args.reason)
//}

//func (u *Group) quitGroupCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.quiteGroupArgs)
//	jsonUser, err := json.Marshal(args.quiteUser)
//	if err != nil {
//		return
//	}
//	u.listener.OnMemberLeave(node.groupId, string(jsonUser))
//}

//func (u *Group) setGroupInfoCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.setGroupInfoArgs)
//	jsonGroup, err := json.Marshal(args.group)
//	if err != nil {
//		return
//	}
//	u.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
//}
////
//func (u *Group) kickGroupMemberCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.kickGroupAgrs)
//	jsonop, err := json.Marshal(args.op)
//	if err != nil {
//		return
//	}
//
//	jsonKickedList, err := json.Marshal(args.kickedList)
//	if err != nil {
//		return
//	}
//
//	u.listener.OnMemberKicked(node.groupId, string(jsonop), string(jsonKickedList))
//}
//
//func (u *open_im_sdk) transferGroupOwnerCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.transferGroupArgs)
//
//	group, err := u.getGroupInfoByGroupId(node.groupId)
//	if err != nil {
//		return
//	}
//	group.OwnerId = args.newOwner.UserId
//
//	jsonGroup, err := json.Marshal(group)
//	if err != nil {
//		return
//	}
//	u.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
//}
//
//func (u *Group) inviteUserToGroupCallback(node open_im_sdk.updateGroupNode) {
//	args := node.Args.(open_im_sdk.inviteUserToGroupArgs)
//	jsonInvitedList, err := json.Marshal(args.invited)
//	if err != nil {
//		return
//	}
//	jsonOp, err := json.Marshal(args.op)
//	if err != nil {
//		return
//	}
//	u.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonInvitedList))
//}

//
//func (u *Group) acceptGroupApplicationCallback(node open_im_sdk.updateGroupNode) {
//	u.GroupApplicationProcessedCallback(node, 1)
//}
//
//func (u *Group) refuseGroupApplicationCallback(node open_im_sdk.updateGroupNode) {
//	u.GroupApplicationProcessedCallback(node, -1)
//}

/*
func (g *Group) doCreateGroup(msg *api.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		return
	}
	u.SyncJoinedGroupList()
	u.syncGroupMemberByGroupID(n.Detail)
	u.onGroupCreated(n.Detail)
}

func (g *Group) doJoinGroup(msg *api.MsgData) {

	u.SyncGroupApplication()

	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		return
	}

	infoSpiltStr := strings.Split(n.Detail, ",")
	var memberFullInfo open_im_sdk.groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = infoSpiltStr[0]
	u.onReceiveJoinApplication(msg.RecvID, memberFullInfo, infoSpiltStr[1])

}

func (g *Group) doQuitGroup(msg *api.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	utils.sdkLog("SyncJoinedGroupList start")
	u.SyncJoinedGroupList()
	utils.sdkLog("SyncJoinedGroupList end")
	u.syncGroupMemberByGroupId(n.Detail)
	utils.sdkLog("SyncJoinedGroupList finish")
	utils.sdkLog("syncGroupMemberByGroupId finish")

	var memberFullInfo open_im_sdk.groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = n.Detail

	u.onMemberLeave(n.Detail, memberFullInfo)
}

func (g *Group) doSetGroupInfo(msg *api.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	utils.sdkLog("doSetGroupInfo, ", n)

	u.SyncJoinedGroupList()
	utils.sdkLog("SyncJoinedGroupList ok")

	var groupInfo open_im_sdk.setGroupInfoReq
	err = json.Unmarshal([]byte(n.Detail), &groupInfo)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	utils.sdkLog("doSetGroupInfo ok , callback ", groupInfo.GroupId, groupInfo)
	u.onGroupInfoChanged(groupInfo.GroupId, groupInfo)
}

func (g *Group) doTransferGroupOwner(msg *api.MsgData) {
	utils.sdkLog("doTransferGroupOwner start...")
	var transfer api.TransferGroupOwnerReq
	var transferContent utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &transferContent)
	if err != nil {
		utils.sdkLog("unmarshal msg.Content, ", err.Error(), msg.Content)
		return
	}
	if err = json.Unmarshal([]byte(transferContent.Detail), &transfer); err != nil {
		utils.sdkLog("unmarshal transferContent", err.Error(), transferContent.Detail)
		return
	}
	u.onTransferGroupOwner(&transfer)
}


func (u *Group) onTransferGroupOwner(transfer *open_im_sdk.TransferGroupOwnerReq) {
	if u.loginUserID == transfer.NewOwner || u.loginUserID == transfer.OldOwner {
		u.SyncGroupApplication()
	}
	u.syncGroupMemberByGroupId(transfer.GroupID)

	gInfo, err := u.getLocalGroupsInfoByGroupID(transfer.GroupID)
	if err != nil {
		sdkLog("onTransferGroupOwner, err ", err.Error(), transfer.GroupID, transfer.OldOwner, transfer.NewOwner, transfer.OldOwner)
		return
	}
	changeInfo := changeGroupInfo{
		data:       *gInfo,
		changeType: 5,
	}
	bChangeInfo, err := json.Marshal(changeInfo)
	if err != nil {
		sdkLog("updateTransferGroupOwner, ", err.Error())
		return
	}
	u.listener.OnGroupInfoChanged(transfer.GroupID, string(bChangeInfo))
	sdkLog("onTransferGroupOwner success")
}

func (g *Group) doAcceptGroupApplication(msg *api.MsgData) {
	utils.sdkLog("doAcceptGroupApplication start...")
	var acceptInfo utils.GroupApplicationInfo
	var acceptContent utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &acceptContent)
	if err != nil {
		utils.sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
		return
	}
	err = json.Unmarshal([]byte(acceptContent.Detail), &acceptInfo)
	if err != nil {
		utils.sdkLog("unmarshal acceptContent.Detail", err.Error(), msg.Content)
		return
	}

	u.onAcceptGroupApplication(&acceptInfo)
}

func (u *Group) onAcceptGroupApplication(groupMember *open_im_sdk.GroupApplicationInfo) {
	member := open_im_sdk.groupMemberFullInfo{
		GroupId:  groupMember.Info.GroupId,
		Role:     0,
		JoinTime: uint64(groupMember.Info.AddTime),
	}
	if groupMember.Info.ToUser == "0" {
		member.UserId = groupMember.Info.FromUser
		member.NickName = groupMember.Info.FromUserNickName
		member.FaceUrl = groupMember.Info.FromUserFaceUrl
	} else {
		member.UserId = groupMember.Info.ToUser
		member.NickName = groupMember.Info.ToUserNickname
		member.FaceUrl = groupMember.Info.ToUserFaceUrl
	}

	bOp, err := json.Marshal(member)
	if err != nil {
		utils.sdkLog("Marshal, ", err.Error())
		return
	}

	var memberList []open_im_sdk.groupMemberFullInfo
	memberList = append(memberList, member)
	bMemberListr, err := json.Marshal(memberList)
	if err != nil {
		utils.sdkLog("onAcceptGroupApplication", err.Error())
		return
	}
	if u.loginUserID == member.UserId {
		u.SyncJoinedGroupList()
		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), 1, groupMember.Info.HandledMsg)
	}
	//g.SyncGroupApplication()
	u.syncGroupMemberByGroupId(groupMember.Info.GroupId)
	u.listener.OnMemberEnter(groupMember.Info.GroupId, string(bMemberListr))

	utils.sdkLog("onAcceptGroupApplication success")
}

func (g *Group) doRefuseGroupApplication(msg *api.MsgData) {
	// do nothing
	utils.sdkLog("doRefuseGroupApplication start...")
	var refuseInfo utils.GroupApplicationInfo
	var refuseContent utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &refuseContent)
	if err != nil {
		utils.sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
		return
	}
	err = json.Unmarshal([]byte(refuseContent.Detail), &refuseInfo)
	if err != nil {
		utils.sdkLog("unmarshal RefuseContent.Detail", err.Error(), msg.Content)
		return
	}

	u.onRefuseGroupApplication(&refuseInfo)
}


func (u *Group) onRefuseGroupApplication(groupMember *open_im_sdk.GroupApplicationInfo) {
	member := open_im_sdk.groupMemberFullInfo{
		GroupId:  groupMember.Info.GroupId,
		Role:     0,
		JoinTime: uint64(groupMember.Info.AddTime),
	}
	if groupMember.Info.ToUser == "0" {
		member.UserId = groupMember.Info.FromUser
		member.NickName = groupMember.Info.FromUserNickName
		member.FaceUrl = groupMember.Info.FromUserFaceUrl
	} else {
		member.UserId = groupMember.Info.ToUser
		member.NickName = groupMember.Info.ToUserNickname
		member.FaceUrl = groupMember.Info.ToUserFaceUrl
	}

	bOp, err := json.Marshal(member)
	if err != nil {
		utils.sdkLog("Marshal, ", err.Error())
		return
	}

	if u.loginUserID == member.UserId {
		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), -1, groupMember.Info.HandledMsg)
	}

	utils.sdkLog("onRefuseGroupApplication success")
}

func (g *Group) doKickGroupMember(msg *api.MsgData) {
	var notification utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &notification)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	utils.sdkLog("doKickGroupMember ", *msg, msg.Content)
	var kickReq open_im_sdk.kickGroupMemberApiReq
	err = json.Unmarshal([]byte(notification.Detail), &kickReq)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := u.getGroupMembersInfoFromLocal(kickReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(kickReq.UidListInfo) == 0 {
		utils.sdkLog("len: ", len(opList), len(kickReq.UidListInfo))
	}
	//	g.syncGroupMember()
	u.SyncJoinedGroupList()
	u.syncGroupMemberByGroupId(kickReq.GroupID)
	//u.SyncJoinedGroupList()
	//u.syncGroupMemberByGroupId(kickReq.GroupID)
	if len(opList) > 0 {
		u.OnMemberKicked(kickReq.GroupID, opList[0], kickReq.UidListInfo)
	} else {
		var op open_im_sdk.groupMemberFullInfo
		op.NickName = "manager"
		u.OnMemberKicked(kickReq.GroupID, op, kickReq.UidListInfo)
	}

}


func (g *Group) OnMemberKicked(groupId string, op open_im_sdk.groupMemberFullInfo, memberList []open_im_sdk.groupMemberFullInfo) {
	jsonOp, err := json.Marshal(op)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), op)
		return
	}

	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		utils.sdkLog("marshal faile, ", err.Error(), memberList)
		return
	}
	g.listener.OnMemberKicked(groupId, string(jsonOp), string(jsonMemberList))
}

func (g *Group) doInviteUserToGroup(msg *api.MsgData) {
	var notification utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &notification)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	var inviteReq open_im_sdk.inviteUserToGroupReq
	err = json.Unmarshal([]byte(notification.Detail), &inviteReq)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), notification.Detail)
		return
	}

	memberList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, inviteReq.UidList)
	if err != nil {
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, tList)
	utils.sdkLog("getGroupMembersInfoFromSvr, ", inviteReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(memberList) == 0 {
		utils.sdkLog("len: ", len(opList), len(memberList))
		return
	}
	for _, v := range inviteReq.UidList {
		if u.loginUserID == v {

			u.SyncJoinedGroupList()
			utils.sdkLog("SyncJoinedGroupList, ", v)
			break
		}
	}

	u.syncGroupMemberByGroupId(inviteReq.GroupID)
	utils.sdkLog("syncGroupMemberByGroupId, ", inviteReq.GroupID)
	u.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
}



func (g *Group) onMemberEnter(groupId string, memberList []open_im_sdk.groupMemberFullInfo) {
	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), jsonMemberList)
		return
	}
	g.listener.OnMemberEnter(groupId, string(jsonMemberList))
}
func (g *Group) onReceiveJoinApplication(groupAdminId string, member open_im_sdk.groupMemberFullInfo, opReason string) {
	jsonMember, err := json.Marshal(member)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), jsonMember)
		return
	}
	g.listener.OnReceiveJoinApplication(groupAdminId, string(jsonMember), opReason)
}
func (g *Group) onMemberLeave(groupId string, member open_im_sdk.groupMemberFullInfo) {
	jsonMember, err := json.Marshal(member)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), jsonMember)
		return
	}
	g.listener.OnMemberLeave(groupId, string(jsonMember))
}

func (g *Group) onGroupInfoChanged(groupId string, changeInfos open_im_sdk.setGroupInfoReq) {
	jsonGroupInfo, err := json.Marshal(changeInfos)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), jsonGroupInfo)
		return
	}
	utils.sdkLog(string(jsonGroupInfo))
	g.listener.OnGroupInfoChanged(groupId, string(jsonGroupInfo))
}
func (g *Group) OnMemberInvited(groupId string, op open_im_sdk.groupMemberFullInfo, memberList []open_im_sdk.groupMemberFullInfo) {
	jsonOp, err := json.Marshal(op)
	if err != nil {
		utils.sdkLog("marshal failed, ", err.Error(), op)
		return
	}

	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		utils.sdkLog("marshal faile, ", err.Error(), memberList)
		return
	}
	g.listener.OnMemberInvited(groupId, string(jsonOp), string(jsonMemberList))
}


*/

//func (u *Group) delGroupRequestFromGroupRequest(info open_im_sdk.GroupReqListInfo) error {
//	return u.delRequestFromGroupRequest(info)
//}

//1
//func (u *Group) refuseGroupApplication(access *open_im_sdk.accessOrRefuseGroupApplicationReq, callback common.Base, operationID string) error {
//	resp, err := utils.post2Api(open_im_sdk.acceptGroupApplicationRouter, access, u.token)
//	if err != nil {
//		return err
//	}
//
//	var ret open_im_sdk.commonResp
//	err = json.Unmarshal(resp, &ret)
//	if err != nil {
//		return err
//	}
//	if ret.ErrCode != 0 {
//		return errors.New(ret.ErrMsg)
//	}
//	return nil
//}
