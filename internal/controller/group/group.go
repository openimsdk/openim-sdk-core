package group

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
)


type Group struct {
	listener OnGroupListener
	token          string
	loginUserID    string
	db             *db.DataBase
	p              *ws.PostApi
}



func (u *Group) doGroupMsg(msg * api.MsgData) {
	if u.listener == nil {
		return
	}
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == constant.SvrConf.Platform {
		return
	}

	go func() {
		switch msg.ContentType {
		case constant.TransferGroupOwnerTip:
			u.doTransferGroupOwner(msg)
		case constant.CreateGroupTip:
			u.doCreateGroup(msg)
		case constant.JoinGroupTip:
			u.doJoinGroup(msg)
		case constant.QuitGroupTip:
			u.doQuitGroup(msg)
		case constant.SetGroupInfoTip:
			u.doSetGroupInfo(msg)
		case constant.AcceptGroupApplicationTip:
			u.doAcceptGroupApplication(msg)
		case constant.RefuseGroupApplicationTip:
			u.doRefuseGroupApplication(msg)
		case constant.KickGroupMemberTip:
			u.doKickGroupMember(msg)
		case constant.InviteUserToGroupTip:
			u.doInviteUserToGroup(msg)
		default:
			log.Error("0","ContentType tip failed, ", msg.ContentType)
		}
	}()
}

func (u *Group) doCreateGroup(msg *api.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		return
	}
	u.syncJoinedGroupInfo()
	u.syncGroupMemberByGroupID(n.Detail)
	u.onGroupCreated(n.Detail)
}

func (u *Group) doJoinGroup(msg *api.MsgData) {
	//
	//u.syncGroupRequest()
	//
	//var n utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &n)
	//if err != nil {
	//	return
	//}
	//
	//infoSpiltStr := strings.Split(n.Detail, ",")
	//var memberFullInfo open_im_sdk.groupMemberFullInfo
	//memberFullInfo.UserId = msg.SendID
	//memberFullInfo.GroupId = infoSpiltStr[0]
	//u.onReceiveJoinApplication(msg.RecvID, memberFullInfo, infoSpiltStr[1])

}

func (u *Group) doQuitGroup(msg *api.MsgData) {
	//var n utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &n)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
	//	return
	//}
	//
	//utils.sdkLog("syncJoinedGroupInfo start")
	//u.syncJoinedGroupInfo()
	//utils.sdkLog("syncJoinedGroupInfo end")
	//u.syncGroupMemberByGroupId(n.Detail)
	//utils.sdkLog("syncJoinedGroupInfo finish")
	//utils.sdkLog("syncGroupMemberByGroupId finish")
	//
	//var memberFullInfo open_im_sdk.groupMemberFullInfo
	//memberFullInfo.UserId = msg.SendID
	//memberFullInfo.GroupId = n.Detail
	//
	//u.onMemberLeave(n.Detail, memberFullInfo)
}

func (u *Group) doSetGroupInfo(msg *api.MsgData) {
	//var n utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &n)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
	//	return
	//}
	//utils.sdkLog("doSetGroupInfo, ", n)
	//
	//u.syncJoinedGroupInfo()
	//utils.sdkLog("syncJoinedGroupInfo ok")
	//
	//var groupInfo open_im_sdk.setGroupInfoReq
	//err = json.Unmarshal([]byte(n.Detail), &groupInfo)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
	//	return
	//}
	//utils.sdkLog("doSetGroupInfo ok , callback ", groupInfo.GroupId, groupInfo)
	//u.onGroupInfoChanged(groupInfo.GroupId, groupInfo)
}

func (u *Group) doTransferGroupOwner(msg *api.MsgData) {
	//utils.sdkLog("doTransferGroupOwner start...")
	//var transfer api.TransferGroupOwnerReq
	//var transferContent utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &transferContent)
	//if err != nil {
	//	utils.sdkLog("unmarshal msg.Content, ", err.Error(), msg.Content)
	//	return
	//}
	//if err = json.Unmarshal([]byte(transferContent.Detail), &transfer); err != nil {
	//	utils.sdkLog("unmarshal transferContent", err.Error(), transferContent.Detail)
	//	return
	//}
	//u.onTransferGroupOwner(&transfer)
}
//
//func (u *Group) onTransferGroupOwner(transfer *open_im_sdk.TransferGroupOwnerReq) {
//	//if u.loginUserID == transfer.NewOwner || u.loginUserID == transfer.OldOwner {
//	//	u.syncGroupRequest()
//	//}
//	//u.syncGroupMemberByGroupId(transfer.GroupID)
//	//
//	//gInfo, err := u.getLocalGroupsInfoByGroupID(transfer.GroupID)
//	//if err != nil {
//	//	sdkLog("onTransferGroupOwner, err ", err.Error(), transfer.GroupID, transfer.OldOwner, transfer.NewOwner, transfer.OldOwner)
//	//	return
//	//}
//	//changeInfo := changeGroupInfo{
//	//	data:       *gInfo,
//	//	changeType: 5,
//	//}
//	//bChangeInfo, err := json.Marshal(changeInfo)
//	//if err != nil {
//	//	sdkLog("updateTransferGroupOwner, ", err.Error())
//	//	return
//	//}
//	//u.listener.OnGroupInfoChanged(transfer.GroupID, string(bChangeInfo))
//	//sdkLog("onTransferGroupOwner success")
//}

func (u *Group) doAcceptGroupApplication(msg *api.MsgData) {
	//utils.sdkLog("doAcceptGroupApplication start...")
	//var acceptInfo utils.GroupApplicationInfo
	//var acceptContent utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &acceptContent)
	//if err != nil {
	//	utils.sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
	//	return
	//}
	//err = json.Unmarshal([]byte(acceptContent.Detail), &acceptInfo)
	//if err != nil {
	//	utils.sdkLog("unmarshal acceptContent.Detail", err.Error(), msg.Content)
	//	return
	//}
	//
	//u.onAcceptGroupApplication(&acceptInfo)
}
//func (u *Group) onAcceptGroupApplication(groupMember *open_im_sdk.GroupApplicationInfo) {
//	member := open_im_sdk.groupMemberFullInfo{
//		GroupId:  groupMember.Info.GroupId,
//		Role:     0,
//		JoinTime: uint64(groupMember.Info.AddTime),
//	}
//	if groupMember.Info.ToUser == "0" {
//		member.UserId = groupMember.Info.FromUser
//		member.NickName = groupMember.Info.FromUserNickName
//		member.FaceUrl = groupMember.Info.FromUserFaceUrl
//	} else {
//		member.UserId = groupMember.Info.ToUser
//		member.NickName = groupMember.Info.ToUserNickname
//		member.FaceUrl = groupMember.Info.ToUserFaceUrl
//	}
//
//	bOp, err := json.Marshal(member)
//	if err != nil {
//		utils.sdkLog("Marshal, ", err.Error())
//		return
//	}
//
//	var memberList []open_im_sdk.groupMemberFullInfo
//	memberList = append(memberList, member)
//	bMemberListr, err := json.Marshal(memberList)
//	if err != nil {
//		utils.sdkLog("onAcceptGroupApplication", err.Error())
//		return
//	}
//	if u.loginUserID == member.UserId {
//		u.syncJoinedGroupInfo()
//		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), 1, groupMember.Info.HandledMsg)
//	}
//	//g.syncGroupRequest()
//	u.syncGroupMemberByGroupId(groupMember.Info.GroupId)
//	u.listener.OnMemberEnter(groupMember.Info.GroupId, string(bMemberListr))
//
//	utils.sdkLog("onAcceptGroupApplication success")
//}

func (u *Group) doRefuseGroupApplication(msg *api.MsgData) {
	//// do nothing
	//utils.sdkLog("doRefuseGroupApplication start...")
	//var refuseInfo utils.GroupApplicationInfo
	//var refuseContent utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &refuseContent)
	//if err != nil {
	//	utils.sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
	//	return
	//}
	//err = json.Unmarshal([]byte(refuseContent.Detail), &refuseInfo)
	//if err != nil {
	//	utils.sdkLog("unmarshal RefuseContent.Detail", err.Error(), msg.Content)
	//	return
	//}
	//
	//u.onRefuseGroupApplication(&refuseInfo)
}
//
//func (u *Group) onRefuseGroupApplication(groupMember *open_im_sdk.GroupApplicationInfo) {
//	//member := open_im_sdk.groupMemberFullInfo{
//	//	GroupId:  groupMember.Info.GroupId,
//	//	Role:     0,
//	//	JoinTime: uint64(groupMember.Info.AddTime),
//	//}
//	//if groupMember.Info.ToUser == "0" {
//	//	member.UserId = groupMember.Info.FromUser
//	//	member.NickName = groupMember.Info.FromUserNickName
//	//	member.FaceUrl = groupMember.Info.FromUserFaceUrl
//	//} else {
//	//	member.UserId = groupMember.Info.ToUser
//	//	member.NickName = groupMember.Info.ToUserNickname
//	//	member.FaceUrl = groupMember.Info.ToUserFaceUrl
//	//}
//	//
//	//bOp, err := json.Marshal(member)
//	//if err != nil {
//	//	utils.sdkLog("Marshal, ", err.Error())
//	//	return
//	//}
//	//
//	//if u.loginUserID == member.UserId {
//	//	u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), -1, groupMember.Info.HandledMsg)
//	//}
//	//
//	//utils.sdkLog("onRefuseGroupApplication success")
//}

func (u *Group) doKickGroupMember(msg *api.MsgData) {
	//var notification utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &notification)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
	//	return
	//}
	//utils.sdkLog("doKickGroupMember ", *msg, msg.Content)
	//var kickReq open_im_sdk.kickGroupMemberApiReq
	//err = json.Unmarshal([]byte(notification.Detail), &kickReq)
	//if err != nil {
	//	utils.sdkLog("unmarshal failed, ", err.Error())
	//	return
	//}
	//
	//tList := make([]string, 1)
	//tList = append(tList, msg.SendID)
	//opList, err := u.getGroupMembersInfoFromLocal(kickReq.GroupID, tList)
	//if err != nil {
	//	return
	//}
	//if len(opList) == 0 || len(kickReq.UidListInfo) == 0 {
	//	utils.sdkLog("len: ", len(opList), len(kickReq.UidListInfo))
	//}
	////	g.syncGroupMember()
	//u.syncJoinedGroupInfo()
	//u.syncGroupMemberByGroupId(kickReq.GroupID)
	////u.syncJoinedGroupInfo()
	////u.syncGroupMemberByGroupId(kickReq.GroupID)
	//if len(opList) > 0 {
	//	u.OnMemberKicked(kickReq.GroupID, opList[0], kickReq.UidListInfo)
	//} else {
	//	var op open_im_sdk.groupMemberFullInfo
	//	op.NickName = "manager"
	//	u.OnMemberKicked(kickReq.GroupID, op, kickReq.UidListInfo)
	//}

}
//
//func (g *Group) OnMemberKicked(groupId string, op open_im_sdk.groupMemberFullInfo, memberList []open_im_sdk.groupMemberFullInfo) {
//	//jsonOp, err := json.Marshal(op)
//	//if err != nil {
//	//	utils.sdkLog("marshal failed, ", err.Error(), op)
//	//	return
//	//}
//	//
//	//jsonMemberList, err := json.Marshal(memberList)
//	//if err != nil {
//	//	utils.sdkLog("marshal faile, ", err.Error(), memberList)
//	//	return
//	//}
//	//g.listener.OnMemberKicked(groupId, string(jsonOp), string(jsonMemberList))
//}

func (u *Group) doInviteUserToGroup(msg *api.MsgData) {
	//var notification utils.NotificationContent
	//err := json.Unmarshal([]byte(msg.Content), &notification)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
	//	return
	//}
	//
	//var inviteReq open_im_sdk.inviteUserToGroupReq
	//err = json.Unmarshal([]byte(notification.Detail), &inviteReq)
	//if err != nil {
	//	utils.sdkLog("unmarshal, ", err.Error(), notification.Detail)
	//	return
	//}
	//
	//memberList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, inviteReq.UidList)
	//if err != nil {
	//	return
	//}
	//
	//tList := make([]string, 1)
	//tList = append(tList, msg.SendID)
	//opList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, tList)
	//utils.sdkLog("getGroupMembersInfoFromSvr, ", inviteReq.GroupID, tList)
	//if err != nil {
	//	return
	//}
	//if len(opList) == 0 || len(memberList) == 0 {
	//	utils.sdkLog("len: ", len(opList), len(memberList))
	//	return
	//}
	//for _, v := range inviteReq.UidList {
	//	if u.loginUserID == v {
	//
	//		u.syncJoinedGroupInfo()
	//		utils.sdkLog("syncJoinedGroupInfo, ", v)
	//		break
	//	}
	//}
	//
	//u.syncGroupMemberByGroupId(inviteReq.GroupID)
	//utils.sdkLog("syncGroupMemberByGroupId, ", inviteReq.GroupID)
	//u.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
}

func (g *Group) onGroupCreated(groupID string) {
	g.listener.OnGroupCreated(groupID)
}
//
//func (g *Group) onMemberEnter(groupId string, memberList []open_im_sdk.groupMemberFullInfo) {
//	jsonMemberList, err := json.Marshal(memberList)
//	if err != nil {
//		utils.sdkLog("marshal failed, ", err.Error(), jsonMemberList)
//		return
//	}
//	g.listener.OnMemberEnter(groupId, string(jsonMemberList))
//}
//func (g *Group) onReceiveJoinApplication(groupAdminId string, member open_im_sdk.groupMemberFullInfo, opReason string) {
//	jsonMember, err := json.Marshal(member)
//	if err != nil {
//		utils.sdkLog("marshal failed, ", err.Error(), jsonMember)
//		return
//	}
//	g.listener.OnReceiveJoinApplication(groupAdminId, string(jsonMember), opReason)
//}
//func (g *Group) onMemberLeave(groupId string, member open_im_sdk.groupMemberFullInfo) {
//	jsonMember, err := json.Marshal(member)
//	if err != nil {
//		utils.sdkLog("marshal failed, ", err.Error(), jsonMember)
//		return
//	}
//	g.listener.OnMemberLeave(groupId, string(jsonMember))
//}
//
//func (g *Group) onGroupInfoChanged(groupId string, changeInfos open_im_sdk.setGroupInfoReq) {
//	jsonGroupInfo, err := json.Marshal(changeInfos)
//	if err != nil {
//		utils.sdkLog("marshal failed, ", err.Error(), jsonGroupInfo)
//		return
//	}
//	utils.sdkLog(string(jsonGroupInfo))
//	g.listener.OnGroupInfoChanged(groupId, string(jsonGroupInfo))
//}
//func (g *Group) OnMemberInvited(groupId string, op open_im_sdk.groupMemberFullInfo, memberList []open_im_sdk.groupMemberFullInfo) {
//	jsonOp, err := json.Marshal(op)
//	if err != nil {
//		utils.sdkLog("marshal failed, ", err.Error(), op)
//		return
//	}
//
//	jsonMemberList, err := json.Marshal(memberList)
//	if err != nil {
//		utils.sdkLog("marshal faile, ", err.Error(), memberList)
//		return
//	}
//	g.listener.OnMemberInvited(groupId, string(jsonOp), string(jsonMemberList))
//}

func (u *Group) createGroup(callback common.Base, group sdk.CreateGroupBaseInfoParam,
	memberList sdk.CreateGroupMemberRoleParam, operationID string) *sdk.CreateGroupCallback {
	apiReq := api.CreateGroupReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = u.loginUserID
	apiReq.GroupName = group.GroupName
	apiReq.GroupType = group.GroupType
	apiReq.MemberList = memberList
	commData := u.p.PostFatalCallback(callback, constant.CreateGroupRouter, apiReq, u.token)
	realData := api.CreateGroupResp{}
	err := mapstructure.Decode(commData.Data, &realData.GroupInfo)
	if err != nil{
		callback.OnError(constant.ErrData.ErrCode, constant.ErrData.ErrMsg)
		return nil
	}
	u.syncJoinedGroupInfo()
	u.syncGroupMemberByGroupID(realData.GroupInfo.GroupID)
	return &sdk.CreateGroupCallback{GroupInfo:realData.GroupInfo}
}

func (u *Group) joinGroup(groupID, reqMsg string, callback common.Base, operationID string) *api.CommDataResp {
	apiReq := api.JoinGroupReq{}
	apiReq.OperationID = operationID
	apiReq.ReqMessage = reqMsg
	apiReq.GroupID = groupID
	commData := u.p.PostFatalCallback(callback, constant.JoinGroupRouter, apiReq, u.token)
	u.syncApplyGroupRequest()
	return commData
}

func (u *Group) quitGroup(groupID string, callback common.Base, operationID string) *api.CommDataResp {
	apiReq := api.QuitGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	commData := u.p.PostFatalCallback(callback, constant.QuitGroupRouter, apiReq, u.token)
	u.syncGroupMemberByGroupID(groupID) //todo
	u.syncApplyGroupRequest()
	return commData
}


func (u *Group) getJoinedGroupList(callback common.Base, operationID string) sdk.GetJoinedGroupListCallback {
	groupList, err := u.db.GetJoinedGroupList()
	common.CheckErr(callback, err, operationID)
	return groupList
}


func (u *Group) getGroupsInfo(groupIdList sdk.GetGroupsInfoParam, callback common.Base, operationID string) sdk.GetGroupsInfoCallback {
	groupList, err := u.db.GetJoinedGroupList()
	common.CheckErr(callback, err, operationID)
	var result sdk.GetGroupsInfoCallback
	for _, v := range groupList{
		in := false
		for _, k := range groupIdList{
			if v.GroupID == k{
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


func (u *Group) setGroupInfo(callback common.Base, groupInfo sdk.SetGroupInfoParam, groupID, operationID string)  *api.CommDataResp{
	apiReq := api.SetGroupInfoReq{}
	apiReq.GroupName = groupInfo.GroupName
	apiReq.FaceUrl = groupInfo.FaceUrl
	apiReq.Notification = groupInfo.Notification
	apiReq.Introduction = groupInfo.Introduction
	apiReq.Ex = groupInfo.Ex
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	commData := u.p.PostFatalCallback(callback, constant.SetGroupInfoRouter, apiReq, u.token)
	u.syncJoinedGroupInfo()
	return commData
}

//todo
func (u *Group) getGroupMemberList(callback common.Base, groupID string, filter int32, next int32,  operationID string) sdk.GetGroupMemberListCallback{
	groupInfoList, err := u.db.GetGroupMemberListByGroupID(groupID)
	common.CheckErr(callback, err, operationID)
	return sdk.GetGroupMemberListCallback{MemberList: groupInfoList, NextSeq:0}
}

//todo
func (u *Group) getGroupMembersInfo(callback common.Base, groupID string, userList sdk.GetGroupMembersInfoParam, operationID string) sdk.GetGroupMembersInfoCallback {
	groupInfoList, err := u.db.GetGroupMemberListByGroupID(groupID)
	common.CheckErr(callback, err, operationID)
	return groupInfoList
}

func (u *Group) kickGroupMember(callback common.Base, groupID string, memberList sdk.KickGroupMemberParam, reason string,  operationID string) sdk.KickGroupMemberCallback {
	apiReq := api.KickGroupMemberReq{}
	apiReq.GroupID = groupID
	apiReq.KickedUserIDList =  memberList
	apiReq.Reason = reason
	apiReq.OperationID = operationID
	commData := u.p.PostFatalCallback(callback, constant.KickGroupMemberRouter, apiReq, u.token)
	u.syncJoinedGroupInfo()
	realData := api.KickGroupMemberResp{}
	err := mapstructure.Decode(commData.Data, &realData.UserIDResultList)
	common.CheckDataErr(callback, err, operationID)
	return realData.UserIDResultList
}

//1
func (u *Group) transferGroupOwner(callback common.Base, groupID, newOwnerUserID string,  operationID string) *api.CommDataResp {
	apiReq := api.TransferGroupOwnerReq{}
	apiReq.GroupID = groupID
	apiReq.NewOwnerUserID = newOwnerUserID
	apiReq.OperationID = operationID
	apiReq.OldOwnerUserID = u.loginUserID
	commData := u.p.PostFatalCallback(callback, constant.TransferGroupRouter, apiReq, u.token)
	u.syncJoinedGroupMember()
	u.syncGroupMemberByGroupID(groupID)
	return commData
}


func (u *Group) inviteUserToGroup(callback common.Base, groupID, reason string, userList sdk.InviteUserToGroupParam,  operationID string) sdk.InviteUserToGroupCallback {
	apiReq := api.InviteUserToGroupReq{}
	apiReq.GroupID = groupID
	apiReq.Reason = reason
	apiReq.InvitedUserIDList = userList
	apiReq.OperationID = operationID
	commData := u.p.PostFatalCallback(callback, constant.InviteUserToGroupRouter, apiReq, u.token)
	u.syncJoinedGroupMember()
	u.syncGroupMemberByGroupID(groupID)
	var realData sdk.InviteUserToGroupCallback
	err := mapstructure.Decode(commData.Data, &realData)
	common.CheckDataErr(callback, err, operationID)
	return realData
}


//1
func (u *Group) getGroupApplicationList(callback common.Base, operationID string) sdk.GetGroupApplicationListCallback {
	applicationList, err:= u.db.GetRecvGroupApplication()
	common.CheckErr(callback, err, operationID)
	return applicationList
}

func (u *Group) getGroupApplicationListFromSvr() ([]*api.GroupRequest, error) {
	apiReq := api.GetGroupApplicationListReq{}
	apiReq.FromUserID = u.loginUserID
	apiReq.OperationID = utils.OperationIDGenerator()
	commData, err := u.p.PostReturn(constant.GetGroupApplicationListRouter, apiReq, u.token)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	var  realData []*api.GroupRequest
	err = mapstructure.Decode(commData.Data, &realData)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}





func (u *Group) processGroupApplication(callback common.Base, groupID, fromUserID,  handleMsg string,  handleResult int32, operationID string) *api.CommDataResp{
	apiReq := api.ApplicationGroupResponseReq{}
	apiReq.GroupID = groupID
	apiReq.OperationID= operationID
	apiReq.FromUserID = fromUserID
	apiReq.HandleResult = handleResult
	apiReq.HandledMsg = handleMsg
	var commData *api.CommDataResp
	if handleResult == 1{
		commData = u.p.PostFatalCallback(callback, constant.AcceptGroupApplicationRouter, apiReq, u.token)
		u.syncGroupMemberByGroupID(groupID)
	} else if handleResult == -1 {
		commData = u.p.PostFatalCallback(callback, constant.RefuseGroupApplicationRouter, apiReq, u.token)
	}
	u.syncGroupRequest()
	return commData
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





func (u *Group) getJoinedGroupListFromSvr() ([]*api.GroupInfo, error) {
	apiReq := api.GetJoinedGroupListReq{}
	apiReq.OperationID = utils.OperationIDGenerator()
	apiReq.FromUserID = u.loginUserID
	commData, err := u.p.PostReturn(constant.GetJoinedGroupListRouter, apiReq, u.token)
	if err != nil{
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	var result []*api.GroupInfo
	err = mapstructure.Decode(commData.Data, &result)
	if err != nil{
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
//
//func (u *Group) getGroupMembersInfoTry2(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
//	result, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
//	if err != nil || len(result) == 0 {
//		return u.getGroupMembersInfoFromSvr(groupId, memberList)
//	} else {
//		return result, err
//	}
//}

func (u *Group) getGroupMembersInfoFromSvr(groupID string, memberList []string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupMembersInfoReq
	apiReq.OperationID = utils.OperationIDGenerator()
	apiReq.GroupID = groupID
	apiReq.MemberList = memberList
	commData, err := u.p.PostReturn(constant.GetGroupMembersInfoRouter, apiReq, apiReq.OperationID)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	var realData []*api.GroupMemberFullInfo
	err = mapstructure.Decode(commData.Data, &realData){
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	return realData, nil
}







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


func (u *Group) SyncSelfGroupRequest() {

}

func (u *Group) SyncGroupRequest() {
	svrList, err := u.getGroupApplicationListFromSvr()
	if err != nil {
		log.NewError("0", "getGroupApplicationListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupRequest(svrList)
	onLocal, err := u.db.GetRecvGroupApplication()
	if err != nil {
		log.NewError("0", "GetJoinedGroupList failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, _ := common.CheckGroupRequestDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := u.db.InsertGroupRequest(onServer[index])
		if err != nil {
			log.NewError("0", "InsertGroupRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u.db.UpdateGroupRequest(onServer[index])
		if err != nil {
			log.NewError("0", "UpdateGroupRequest failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u.db.DeleteGroupRequest(onServer[index].GroupID, onServer[index].UserID)
		if err != nil {
			log.NewError("0", "DeleteGroupRequest failed ", err.Error())
			continue
		}
	}
}

func (g *Group) SyncApplyGroupRequest() {

}

func (u *Group) SyncJoinedGroupInfo() {
	svrList, err := u.getJoinedGroupListFromSvr()
	if err != nil {
		log.NewError("0", "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupInfo(svrList)
	onLocal, err := u.db.GetJoinedGroupList()
	if err != nil {
		log.NewError("0", "GetRecvFriendApplication failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, _ := common.CheckGroupInfoDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := u.db.InsertGroup(onServer[index])
		if err != nil {
			log.NewError("0", "InsertGroup failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u.db.UpdateGroup(onServer[index])
		if err != nil {
			log.NewError("0", "UpdateGroup failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u.db.DeleteGroup(onServer[index].GroupID)
		if err != nil {
			log.NewError("0", "DeleteGroup failed ", err.Error())
			continue
		}
	}
}


//func (u *Group) getLocalGroupInfoByGroupId1(groupId string) (*Group.groupInfo, error) {
//	return u.getLocalGroupsInfoByGroupID(groupId)
//}

func (u *Group) syncGroupMemberByGroupID(groupID string) {
	svrList, err := u.getGroupAllMemberByGroupIDFromSvr(groupID)
	if err != nil {
		log.NewError("0", "getGroupAllMemberByGroupIDFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupMember(svrList)
	onLocal, err := u.db.GetGroupMemberListByGroupID(groupID)
	if err != nil {
		log.NewError("0", "GetGroupMemberListByGroupID failed ", err.Error())
		return
	}
	log.NewInfo("0", "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, _ := common.CheckGroupMemberDiff(onServer, onLocal)
	for _, index := range aInBNot {
		err := u.db.InsertGroupMember(onServer[index])
		if err != nil {
			log.NewError("0", "InsertGroupMember failed ", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		err := u.db.UpdateGroupMember(onServer[index])
		if err != nil {
			log.NewError("0", "UpdateGroupMember failed ", err.Error())
			continue
		}
	}
	for _, index := range bInANot {
		err := u.db.DeleteGroupMember(onServer[index].GroupID, onServer[index].UserID)
		if err != nil {
			log.NewError("0", "DeleteGroupMember failed ", err.Error())
			continue
		}
	}
}

func (u *Group) syncJoinedGroupMember() {
	groupListOnServer, err := u.getJoinedGroupListFromSvr()
	if err != nil {
		log.Error("0", "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	for _, v := range groupListOnServer{
		u.syncGroupMemberByGroupID(v.GroupID)
	}
}


func (u *Group) getGroupAllMemberByGroupIDFromSvr(groupID string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupAllMemberReq
	apiReq.OperationID = utils.OperationIDGenerator()
	apiReq.GroupID = groupID
	commData, err := u.p.PostReturn(constant.GetGroupAllMemberListRouter, apiReq, apiReq.OperationID)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID)
	}
	var realData []*api.GroupMemberFullInfo
	err = mapstructure.Decode(commData.Data, &realData){
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

