package open_im_sdk

import (
	"encoding/json"
	"errors"
	"strings"
)

func (u *UserRelated) doGroupMsg(msg *MsgData) {
	if u.listener == nil {
		sdkLog("group listener is null")
		return
	}
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == u.SvrConf.Platform {
		sdkLog("sync msg ", msg)
		return
	}

	go func() {
		switch msg.ContentType {
		case TransferGroupOwnerTip:
			u.doTransferGroupOwner(msg)
		case CreateGroupTip:
			u.doCreateGroup(msg)
		case JoinGroupTip:
			u.doJoinGroup(msg)
		case QuitGroupTip:
			u.doQuitGroup(msg)
		case SetGroupInfoTip:
			u.doSetGroupInfo(msg)
		case AcceptGroupApplicationTip:
			u.doAcceptGroupApplication(msg)
		case RefuseGroupApplicationTip:
			u.doRefuseGroupApplication(msg)
		case KickGroupMemberTip:
			u.doKickGroupMember(msg)
		case InviteUserToGroupTip:
			u.doInviteUserToGroup(msg)
		default:
			sdkLog("ContentType tip failed, ", msg.ContentType)
		}
	}()
}

func (u *UserRelated) doCreateGroup(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("doCreateGroup unmarshal failed", err.Error())
		return
	}
	sdkLog("doCreateGroup, ", msg, n)
	u.syncJoinedGroupInfo()

	sdkLog("syncJoinedGroupInfo ok")
	u.syncGroupMemberByGroupId(n.Detail)
	sdkLog("syncGroupMemberByGroupId ok, ", n.Detail)
	u.onGroupCreated(n.Detail)
	sdkLog("onGroupCreated callback, ", n.Detail)
}

func (u *UserRelated) doJoinGroup(msg *MsgData) {

	u.syncGroupRequest()

	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	infoSpiltStr := strings.Split(n.Detail, ",")
	var memberFullInfo groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = infoSpiltStr[0]
	u.onReceiveJoinApplication(msg.RecvID, memberFullInfo, infoSpiltStr[1])

}

func (u *UserRelated) doQuitGroup(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	sdkLog("syncJoinedGroupInfo start")
	u.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo end")
	u.syncGroupMemberByGroupId(n.Detail)
	sdkLog("syncJoinedGroupInfo finish")
	sdkLog("syncGroupMemberByGroupId finish")

	var memberFullInfo groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = n.Detail

	u.onMemberLeave(n.Detail, memberFullInfo)
}

func (u *UserRelated) doSetGroupInfo(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog("doSetGroupInfo, ", n)

	u.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")

	var groupInfo setGroupInfoReq
	err = json.Unmarshal([]byte(n.Detail), &groupInfo)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog("doSetGroupInfo ok , callback ", groupInfo.GroupId, groupInfo)
	u.onGroupInfoChanged(groupInfo.GroupId, groupInfo)
}

func (u *UserRelated) doTransferGroupOwner(msg *MsgData) {
	sdkLog("doTransferGroupOwner start...")
	var transfer TransferGroupOwnerReq
	var transferContent NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &transferContent)
	if err != nil {
		sdkLog("unmarshal msg.Content, ", err.Error(), msg.Content)
		return
	}
	if err = json.Unmarshal([]byte(transferContent.Detail), &transfer); err != nil {
		sdkLog("unmarshal transferContent", err.Error(), transferContent.Detail)
		return
	}
	u.onTransferGroupOwner(&transfer)
}
func (u *UserRelated) onTransferGroupOwner(transfer *TransferGroupOwnerReq) {
	if u.loginUserID == transfer.NewOwner || u.loginUserID == transfer.OldOwner {
		u.syncGroupRequest()
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

func (u *UserRelated) doAcceptGroupApplication(msg *MsgData) {
	sdkLog("doAcceptGroupApplication start...")
	var acceptInfo GroupApplicationInfo
	var acceptContent NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &acceptContent)
	if err != nil {
		sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
		return
	}
	err = json.Unmarshal([]byte(acceptContent.Detail), &acceptInfo)
	if err != nil {
		sdkLog("unmarshal acceptContent.Detail", err.Error(), msg.Content)
		return
	}

	u.onAcceptGroupApplication(&acceptInfo)
}
func (u *UserRelated) onAcceptGroupApplication(groupMember *GroupApplicationInfo) {
	member := groupMemberFullInfo{
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
		sdkLog("Marshal, ", err.Error())
		return
	}

	var memberList []groupMemberFullInfo
	memberList = append(memberList, member)
	bMemberListr, err := json.Marshal(memberList)
	if err != nil {
		sdkLog("onAcceptGroupApplication", err.Error())
		return
	}
	if u.loginUserID == member.UserId {
		u.syncJoinedGroupInfo()
		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), 1, groupMember.Info.HandledMsg)
	}
	//g.syncGroupRequest()
	u.syncGroupMemberByGroupId(groupMember.Info.GroupId)
	u.listener.OnMemberEnter(groupMember.Info.GroupId, string(bMemberListr))

	sdkLog("onAcceptGroupApplication success")
}

func (u *UserRelated) doRefuseGroupApplication(msg *MsgData) {
	// do nothing
	sdkLog("doRefuseGroupApplication start...")
	var refuseInfo GroupApplicationInfo
	var refuseContent NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &refuseContent)
	if err != nil {
		sdkLog("unmarshal msg.Content ", err.Error(), msg.Content)
		return
	}
	err = json.Unmarshal([]byte(refuseContent.Detail), &refuseInfo)
	if err != nil {
		sdkLog("unmarshal RefuseContent.Detail", err.Error(), msg.Content)
		return
	}

	u.onRefuseGroupApplication(&refuseInfo)
}

func (u *UserRelated) onRefuseGroupApplication(groupMember *GroupApplicationInfo) {
	member := groupMemberFullInfo{
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
		sdkLog("Marshal, ", err.Error())
		return
	}

	if u.loginUserID == member.UserId {
		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), -1, groupMember.Info.HandledMsg)
	}

	sdkLog("onRefuseGroupApplication success")
}

func (u *UserRelated) doKickGroupMember(msg *MsgData) {
	var notification NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &notification)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog("doKickGroupMember ", *msg, msg.Content)
	var kickReq kickGroupMemberApiReq
	err = json.Unmarshal([]byte(notification.Detail), &kickReq)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := u.getGroupMembersInfoFromLocal(kickReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(kickReq.UidListInfo) == 0 {
		sdkLog("len: ", len(opList), len(kickReq.UidListInfo))
	}
	//	g.syncGroupMember()
	u.syncJoinedGroupInfo()
	u.syncGroupMemberByGroupId(kickReq.GroupID)
	//u.syncJoinedGroupInfo()
	//u.syncGroupMemberByGroupId(kickReq.GroupID)
	if len(opList) > 0 {
		u.OnMemberKicked(kickReq.GroupID, opList[0], kickReq.UidListInfo)
	} else {
		var op groupMemberFullInfo
		op.NickName = "manager"
		u.OnMemberKicked(kickReq.GroupID, op, kickReq.UidListInfo)
	}

}

func (g *groupListener) OnMemberKicked(groupId string, op groupMemberFullInfo, memberList []groupMemberFullInfo) {
	jsonOp, err := json.Marshal(op)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), op)
		return
	}

	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		sdkLog("marshal faile, ", err.Error(), memberList)
		return
	}
	g.listener.OnMemberKicked(groupId, string(jsonOp), string(jsonMemberList))
}

func (u *UserRelated) doInviteUserToGroup(msg *MsgData) {
	var notification NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &notification)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	var inviteReq inviteUserToGroupReq
	err = json.Unmarshal([]byte(notification.Detail), &inviteReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), notification.Detail)
		return
	}

	memberList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, inviteReq.UidList)
	if err != nil {
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := u.getGroupMembersInfoTry2(inviteReq.GroupID, tList)
	sdkLog("getGroupMembersInfoFromSvr, ", inviteReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(memberList) == 0 {
		sdkLog("len: ", len(opList), len(memberList))
		return
	}
	for _, v := range inviteReq.UidList {
		if u.loginUserID == v {

			u.syncJoinedGroupInfo()
			sdkLog("syncJoinedGroupInfo, ", v)
			break
		}
	}

	u.syncGroupMemberByGroupId(inviteReq.GroupID)
	sdkLog("syncGroupMemberByGroupId, ", inviteReq.GroupID)
	u.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
}

func (g *groupListener) onGroupCreated(groupID string) {
	g.listener.OnGroupCreated(groupID)
}
func (g *groupListener) onMemberEnter(groupId string, memberList []groupMemberFullInfo) {
	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), jsonMemberList)
		return
	}
	g.listener.OnMemberEnter(groupId, string(jsonMemberList))
}
func (g *groupListener) onReceiveJoinApplication(groupAdminId string, member groupMemberFullInfo, opReason string) {
	jsonMember, err := json.Marshal(member)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), jsonMember)
		return
	}
	g.listener.OnReceiveJoinApplication(groupAdminId, string(jsonMember), opReason)
}
func (g *groupListener) onMemberLeave(groupId string, member groupMemberFullInfo) {
	jsonMember, err := json.Marshal(member)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), jsonMember)
		return
	}
	g.listener.OnMemberLeave(groupId, string(jsonMember))
}

func (g *groupListener) onGroupInfoChanged(groupId string, changeInfos setGroupInfoReq) {
	jsonGroupInfo, err := json.Marshal(changeInfos)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), jsonGroupInfo)
		return
	}
	sdkLog(string(jsonGroupInfo))
	g.listener.OnGroupInfoChanged(groupId, string(jsonGroupInfo))
}
func (g *groupListener) OnMemberInvited(groupId string, op groupMemberFullInfo, memberList []groupMemberFullInfo) {
	jsonOp, err := json.Marshal(op)
	if err != nil {
		sdkLog("marshal failed, ", err.Error(), op)
		return
	}

	jsonMemberList, err := json.Marshal(memberList)
	if err != nil {
		sdkLog("marshal faile, ", err.Error(), memberList)
		return
	}
	g.listener.OnMemberInvited(groupId, string(jsonOp), string(jsonMemberList))
}

func (u *UserRelated) createGroup(group groupInfo, memberList []createGroupMemberInfo) (*createGroupResp, error) {
	req := createGroupReq{memberList, group.GroupName, group.Introduction, group.Notification, group.FaceUrl, operationIDGenerator(), group.Ex}
	resp, err := post2Api(createGroupRouter, req, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", createGroupRouter, req)
		return nil, err
	}
	var createGroupResp createGroupResp
	if err = json.Unmarshal(resp, &createGroupResp); err != nil {
		sdkLog("Unmarshal failed, ", err.Error())
		return nil, err
	}
	sdkLog("post2Api ok ", createGroupRouter, req, createGroupResp)

	if createGroupResp.ErrCode != 0 {
		sdkLog("errcode errmsg: ", createGroupResp.ErrCode, createGroupResp.ErrMsg)
		return nil, errors.New(createGroupResp.ErrMsg)
	}

	u.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	u.syncGroupMemberByGroupId(createGroupResp.Data.GroupId)
	sdkLog("syncGroupMemberByGroupId ok")
	return &createGroupResp, nil
}

func (u *UserRelated) joinGroup(groupId, message string) error {
	req := joinGroupReq{groupId, message, operationIDGenerator()}
	resp, err := post2Api(joinGroupRouter, req, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", err.Error(), joinGroupRouter, req)
		return err
	}
	var commonResp commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		sdkLog("Unmarshal", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		sdkLog("commonResp err", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	sdkLog("psot2api ok", joinGroupRouter, req, commonResp)

	u.syncApplyGroupRequest()
	sdkLog("syncApplyGroupRequest ok")

	memberList, err := u.getGroupAllMemberListByGroupIdFromSvr(groupId)
	if err != nil {
		sdkLog("getGroupAllMemberListByGroupIdFromSvr failed", err.Error())
		return err
	}

	var groupAdminUser string
	for _, v := range memberList {
		if v.Role == 1 {
			groupAdminUser = v.UserId
			break
		}
	}
	sdkLog("get admin from svr ok ", groupId, groupAdminUser)
	return nil
}

func (u *UserRelated) quitGroup(groupId string) error {
	req := quitGroupReq{groupId, operationIDGenerator()}
	resp, err := post2Api(quitGroupRouter, req, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", quitGroupRouter, req)
		return err
	}
	var commonResp commonResp
	err = json.Unmarshal(resp, &commonResp)
	if err != nil {
		sdkLog("unmarshal", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		sdkLog("errcode, errmsg", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	sdkLog("post2Api ok ", quitGroupRouter, req, commonResp)

	u.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	u.syncGroupMemberByGroupId(groupId) //todo
	sdkLog("syncGroupMemberByGroupId ok ", groupId)
	return nil
}

func (u *UserRelated) getJoinedGroupListFromLocal() ([]groupInfo, error) {
	return u.getLocalGroupsInfo()
}

func (u *UserRelated) getJoinedGroupListFromSvr() ([]groupInfo, error) {
	var req getJoinedGroupListReq
	req.OperationID = operationIDGenerator()
	sdkLog("getJoinedGroupListRouter ", getJoinedGroupListRouter, req, u.token)
	resp, err := post2Api(getJoinedGroupListRouter, req, u.token)
	if err != nil {
		sdkLog("post api:", err)
		return nil, err
	}

	var stcResp getJoinedGroupListResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		sdkLog("unmarshal, ", err)
		return nil, err
	}

	if stcResp.ErrCode != 0 {
		return nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Data, nil
}

func (u *UserRelated) getGroupsInfo(groupIdList []string) ([]groupInfo, error) {
	req := getGroupsInfoReq{groupIdList, operationIDGenerator()}
	resp, err := post2Api(getGroupsInfoRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	var getGroupsInfoResp getGroupsInfoResp
	err = json.Unmarshal(resp, &getGroupsInfoResp)
	if err != nil {
		return nil, err
	}
	return getGroupsInfoResp.Data, nil
}

func (u *UserRelated) setGroupInfo(newGroupInfo setGroupInfoReq) error {
	g, err := u._getGroupInfoByGroupID(newGroupInfo.GroupId)
	if err != nil {
		sdkLog("findLocalGroupOwnerByGroupId failed, ", newGroupInfo.GroupId, err.Error())
		return err
	}
	if u.loginUserID != g.OwnerUserID {
		sdkLog("no permission, ", u.loginUserID, g.OwnerUserID)
		return errors.New("no permission")
	}
	sdkLog("findLocalGroupOwnerByGroupId ok ", newGroupInfo.GroupId, g.OwnerUserID)

	req := setGroupInfoReq{newGroupInfo.GroupId, newGroupInfo.GroupName, newGroupInfo.Notification, newGroupInfo.Introduction, newGroupInfo.FaceUrl, operationIDGenerator()}
	resp, err := post2Api(setGroupInfoRouter, req, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", setGroupInfoRouter, req, err.Error())
		return err
	}
	var commonResp commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		sdkLog("errcode errmsg: ", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	sdkLog("post2Api ok, ", setGroupInfoRouter, req, commonResp)

	u.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	return nil
}

func (u *UserRelated) getGroupMemberListFromSvr(groupId string, filter int32, next int32) (int32, []groupMemberFullInfo, error) {
	var req getGroupMemberListReq
	req.OperationID = operationIDGenerator()
	req.GroupID = groupId
	req.NextSeq = next
	req.Filter = filter
	resp, err := post2Api(getGroupMemberListRouter, req, u.token)
	if err != nil {
		return 0, nil, err
	}
	var stcResp groupMemberInfoResult
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return 0, nil, err
	}

	if stcResp.ErrCode != 0 {
		sdkLog("errcode, errmsg: ", stcResp.ErrCode, stcResp.ErrMsg)
		return 0, nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Nextseq, stcResp.Data, nil
}

func (u *UserRelated) getGroupMemberListFromLocal(groupId string, filter int32, next int32) (int32, []groupMemberFullInfo, error) {
	memberList, err := u.getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		return 0, nil, err
	}
	return 0, memberList, nil
}

func (u *UserRelated) getGroupMembersInfoFromLocal(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	var result []groupMemberFullInfo
	localMemberList, err := u.getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		return nil, err
	}
	for _, i := range localMemberList {
		for _, j := range memberList {
			if i.UserId == j {
				result = append(result, i)
			}
		}
	}
	return result, nil
}

func (u *UserRelated) getGroupMembersInfoTry2(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	result, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil || len(result) == 0 {
		return u.getGroupMembersInfoFromSvr(groupId, memberList)
	} else {
		return result, err
	}
}

func (u *UserRelated) getGroupMembersInfoFromSvr(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	var req getGroupMembersInfoReq
	req.GroupID = groupId
	req.OperationID = operationIDGenerator()
	req.MemberList = memberList

	resp, err := post2Api(getGroupMembersInfoRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	var sctResp getGroupMembersInfoResp
	err = json.Unmarshal(resp, &sctResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		sdkLog("errcode, errmsg: ", sctResp.ErrCode, sctResp.ErrMsg)
		return nil, errors.New(sctResp.ErrMsg)
	}
	return sctResp.Data, nil
}

func (u *UserRelated) kickGroupMember(groupId string, memberList []string, reason string) ([]idResult, error) {
	var req kickGroupMemberApiReq
	req.OperationID = operationIDGenerator()
	memberListInfo, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil {
		sdkLog("getGroupMembersInfoFromLocal, ", err.Error())
		return nil, err
	}
	req.UidListInfo = memberListInfo
	req.Reason = reason
	req.GroupID = groupId

	resp, err := post2Api(kickGroupMemberRouter, req, u.token)
	if err != nil {
		sdkLog("post2Api failed, ", kickGroupMemberRouter, req, err.Error())
		return nil, err
	}
	sdkLog("url: ", kickGroupMemberRouter, "req:", req, "resp: ", string(resp))

	u.syncGroupMemberByGroupId(groupId)
	sdkLog("syncGroupMemberByGroupId: ", groupId)

	var sctResp kickGroupMemberApiResp
	err = json.Unmarshal(resp, &sctResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error(), resp)
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		sdkLog("resp failed, ", sctResp.ErrCode, sctResp.ErrMsg)
		return nil, errors.New(sctResp.ErrMsg)
	}
	sdkLog("kickGroupMember, ", groupId, memberList, reason, req)
	return sctResp.Data, nil
}

//1
func (u *UserRelated) transferGroupOwner(groupId, userId string) error {
	resp, err := post2Api(transferGroupRouter, transferGroupReq{GroupID: groupId, Uid: userId, OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return err
	}
	var ret commonResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}
	if ret.ErrCode != 0 {
		return errors.New(ret.ErrMsg)
	}

	return nil
}

//1
func (u *UserRelated) inviteUserToGroup(groupId string, reason string, userList []string) ([]idResult, error) {
	var req inviteUserToGroupReq
	req.GroupID = groupId
	req.OperationID = operationIDGenerator()
	req.Reason = reason
	req.UidList = userList
	resp, err := post2Api(inviteUserToGroupRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	u.syncGroupMemberByGroupId(groupId)
	sdkLog("syncGroupMemberByGroupId", groupId)
	var stcResp inviteUserToGroupResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if stcResp.ErrCode != 0 {
		sdkLog("errcode, errmsg: ", stcResp.ErrCode, stcResp.ErrMsg)
		return nil, errors.New(stcResp.ErrMsg)
	}

	sdkLog("inviteUserToGroup, autoSendInviteUserToGroupTip", groupId, reason, userList, req, err)
	return stcResp.Data, nil
}

func (u *UserRelated) getLocalGroupApplicationList(groupId string) (*groupApplicationResult, error) {
	reply, err := u.getOwnLocalGroupApplicationList(groupId)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (u *UserRelated) delGroupRequestFromGroupRequest(info GroupReqListInfo) error {
	return u.delRequestFromGroupRequest(info)
}

//1
func (u *UserRelated) getGroupApplicationList() (*groupApplicationResult, error) {
	resp, err := post2Api(getGroupApplicationListRouter, getGroupApplicationListReq{OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}

	var ret getGroupApplicationListResp
	sdkLog("getGroupApplicationListResp", string(resp))
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		sdkLog("unmarshal failed", err.Error())
		return nil, err
	}
	if ret.ErrCode != 0 {
		sdkLog("errcode, errmsg: ", ret.ErrCode, ret.ErrMsg)
		return nil, errors.New(ret.ErrMsg)
	}

	return &ret.Data, nil
}

//1
func (u *UserRelated) acceptGroupApplication(access *accessOrRefuseGroupApplicationReq) error {
	resp, err := post2Api(acceptGroupApplicationRouter, access, u.token)
	if err != nil {
		return err
	}

	var ret commonResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}
	if ret.ErrCode != 0 {
		return errors.New(ret.ErrMsg)
	}
	return nil
}

//1
func (u *UserRelated) refuseGroupApplication(access *accessOrRefuseGroupApplicationReq) error {
	resp, err := post2Api(acceptGroupApplicationRouter, access, u.token)
	if err != nil {
		return err
	}

	var ret commonResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}
	if ret.ErrCode != 0 {
		return errors.New(ret.ErrMsg)
	}
	return nil
}

func (u *UserRelated) getGroupInfoByGroupId(groupId string) (groupInfo, error) {
	var gList []string
	gList = append(gList, groupId)
	rList, err := u.getGroupsInfo(gList)
	if err == nil && len(rList) == 1 {
		return rList[0], nil
	} else {
		return groupInfo{}, nil
	}

}

type groupListener struct {
	listener OnGroupListener
}

func (u *UserRelated) createGroupCallback(node updateGroupNode) {
	// member list to json
	jsonMemberList, err := json.Marshal(node.Args.(createGroupArgs).initMemberList)
	if err != nil {
		return
	}
	u.listener.OnMemberEnter(node.groupId, string(jsonMemberList))
	u.listener.OnGroupCreated(node.groupId)
}

func (u *UserRelated) joinGroupCallback(node updateGroupNode) {
	args := node.Args.(joinGroupArgs)
	jsonApplyUser, err := json.Marshal(args.applyUser)
	if err != nil {
		return
	}
	u.listener.OnReceiveJoinApplication(node.groupId, string(jsonApplyUser), args.reason)
}

func (u *UserRelated) quitGroupCallback(node updateGroupNode) {
	args := node.Args.(quiteGroupArgs)
	jsonUser, err := json.Marshal(args.quiteUser)
	if err != nil {
		return
	}
	u.listener.OnMemberLeave(node.groupId, string(jsonUser))
}

func (u *UserRelated) setGroupInfoCallback(node updateGroupNode) {
	args := node.Args.(setGroupInfoArgs)
	jsonGroup, err := json.Marshal(args.group)
	if err != nil {
		return
	}
	u.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
}

func (u *UserRelated) kickGroupMemberCallback(node updateGroupNode) {
	args := node.Args.(kickGroupAgrs)
	jsonop, err := json.Marshal(args.op)
	if err != nil {
		return
	}

	jsonKickedList, err := json.Marshal(args.kickedList)
	if err != nil {
		return
	}

	u.listener.OnMemberKicked(node.groupId, string(jsonop), string(jsonKickedList))
}

func (u *UserRelated) transferGroupOwnerCallback(node updateGroupNode) {
	args := node.Args.(transferGroupArgs)

	group, err := u.getGroupInfoByGroupId(node.groupId)
	if err != nil {
		return
	}
	group.OwnerId = args.newOwner.UserId

	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return
	}
	u.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
}

func (u *UserRelated) inviteUserToGroupCallback(node updateGroupNode) {
	args := node.Args.(inviteUserToGroupArgs)
	jsonInvitedList, err := json.Marshal(args.invited)
	if err != nil {
		return
	}
	jsonOp, err := json.Marshal(args.op)
	if err != nil {
		return
	}
	u.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonInvitedList))
}

func (u *UserRelated) GroupApplicationProcessedCallback(node updateGroupNode, process int32) {
	args := node.Args.(applyGroupProcessedArgs)
	list := make([]groupMemberFullInfo, 0)
	for _, v := range args.applyList {
		list = append(list, v.member)
	}
	jsonApplyList, err := json.Marshal(list)
	if err != nil {
		return
	}

	processed := node.Args.(applyGroupProcessedArgs) //receiver : all group member
	var flag = 0
	var idx = 0
	for i, v := range processed.applyList {
		if v.member.UserId == u.loginUserID {
			flag = 1
			idx = i
			break
		}
	}

	if flag == 1 {
		jsonOp, err := json.Marshal(processed.op)
		if err != nil {
			return
		}
		u.listener.OnApplicationProcessed(node.groupId, string(jsonOp), process, processed.applyList[idx].reason)
	}

	if process == 1 {
		jsonOp, err := json.Marshal(processed.op)
		if err != nil {
			return
		}
		u.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonApplyList))
	}
}

func (u *UserRelated) acceptGroupApplicationCallback(node updateGroupNode) {
	u.GroupApplicationProcessedCallback(node, 1)
}

func (u *UserRelated) refuseGroupApplicationCallback(node updateGroupNode) {
	u.GroupApplicationProcessedCallback(node, -1)
}

func (u *UserRelated) syncSelfGroupRequest() {

}

func (u *UserRelated) syncGroupRequest() {
	groupRequestOnServerResp, err := u.getGroupApplicationList()
	if err != nil {
		sdkLog("groupRequestOnServerResp failed", err.Error())
		return
	}
	groupRequestOnServer := groupRequestOnServerResp.GroupApplicationList
	groupRequestOnServerInterface := make([]diff, 0)
	for _, v := range groupRequestOnServer {
		groupRequestOnServerInterface = append(groupRequestOnServerInterface, v)
	}

	groupRequestOnLocalResp, err := u.getLocalGroupApplicationList("")
	if err != nil {
		sdkLog("groupRequestOnLocalResp failed", err.Error())
		return
	}
	groupRequestOnLocal := groupRequestOnLocalResp.GroupApplicationList
	groupRequestOnLocalInterface := make([]diff, 0)
	for _, v := range groupRequestOnLocal {
		groupRequestOnLocalInterface = append(groupRequestOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := checkDiff(groupRequestOnServerInterface, groupRequestOnLocalInterface)

	sdkLog("len ", len(aInBNot), len(bInANot), len(sameA))
	for _, index := range aInBNot {
		err = u.insertIntoRequestToGroupRequest(groupRequestOnServer[index])
		if err != nil {
			sdkLog("insertIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		sdkLog("insertIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

	for _, index := range bInANot {
		err = u.delGroupRequestFromGroupRequest(groupRequestOnLocal[index])
		if err != nil {
			sdkLog("delGroupRequestFromGroupRequest failed", err.Error())
			continue
		}
		sdkLog("delGroupRequestFromGroupRequest ", groupRequestOnLocal[index])
	}
	for _, index := range sameA {
		if err = u.replaceIntoRequestToGroupRequest(groupRequestOnServer[index]); err != nil {
			sdkLog("replaceIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		sdkLog("replaceIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

}

func (g *groupListener) syncApplyGroupRequest() {

}

func (u *UserRelated) syncJoinedGroupInfo() {
	groupListOnServer, err := u.getJoinedGroupListFromSvr()
	if err != nil {
		sdkLog("groupListOnServer failed", err.Error())
		return
	}
	groupListOnServerInterface := make([]diff, 0)
	for _, v := range groupListOnServer {
		groupListOnServerInterface = append(groupListOnServerInterface, v)
	}

	groupListOnLocal, err := u.getJoinedGroupListFromLocal()
	if err != nil {
		sdkLog("groupListOnLocal failed", err.Error())
		return
	}
	groupListOnLocalInterface := make([]diff, 0)
	for _, v := range groupListOnLocal {
		groupListOnLocalInterface = append(groupListOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := checkDiff(groupListOnServerInterface, groupListOnLocalInterface)

	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupInfo(groupListOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupInfo failed", err.Error(), groupListOnServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupInfo(groupListOnLocal[index].GroupId)
		if err != nil {
			sdkLog("delLocalGroupInfo failed", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		if err = u.replaceLocalGroupInfo(groupListOnServer[index]); err != nil {
			sdkLog("replaceLocalGroupInfo failed", err.Error())
			continue
		}
	}
}

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

func (u *UserRelated) getLocalGroupInfoByGroupId1(groupId string) (*groupInfo, error) {
	return u.getLocalGroupsInfoByGroupID(groupId)
}

func (u *UserRelated) syncGroupMemberByGroupId(groupId string) {
	groupMemberOnServer, err := u.getGroupAllMemberListByGroupIdFromSvr(groupId)
	if err != nil {
		sdkLog("syncGroupMemberByGroupId failed", err.Error())
		return
	}
	sdkLog("getGroupAllMemberListByGroupIdFromSvr, ", groupId, len(groupMemberOnServer))

	groupMemberOnServerInterface := make([]diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := u.getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		sdkLog("getLocalGroupMemberListByGroupID failed", err.Error())
		return
	}
	sdkLog("getLocalGroupMemberListByGroupID, ", groupId, len(groupMemberOnLocal))

	groupMemberOnLocalInterface := make([]diff, 0)
	for _, v := range groupMemberOnLocal {
		groupMemberOnLocalInterface = append(groupMemberOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := checkDiff(groupMemberOnServerInterface, groupMemberOnLocalInterface)
	//0 0 2 2 3
	sdkLog("diff len: ", len(aInBNot), len(bInANot), len(sameA), len(groupMemberOnServerInterface), len(groupMemberOnLocalInterface))
	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupMember failed", err.Error(), "index", index, groupMemberOnServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = u.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			sdkLog("replaceLocalGroupMemberInfo failed", err.Error())
			continue
		}
	}

}

func (u *UserRelated) syncJoinedGroupMember() {
	groupMemberOnServer, err := u.getJoinGroupAllMemberList()
	if err != nil {
		sdkLog("getJoinGroupAllMemberList failed", err.Error())
		return
	}
	groupMemberOnServerInterface := make([]diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := u.getLocalGroupMemberList()
	if err != nil {
		sdkLog("getLocalGroupMemberList failed", err.Error())
		return
	}
	groupMemberOnLocalInterface := make([]diff, 0)
	for _, v := range groupMemberOnLocal {
		groupMemberOnLocalInterface = append(groupMemberOnLocalInterface, v)
	}

	aInBNot, bInANot, sameA, _ := checkDiff(groupMemberOnServerInterface, groupMemberOnLocalInterface)

	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = u.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			sdkLog(err.Error())
			continue
		}
	}

}

func (u *UserRelated) getJoinGroupAllMemberList() ([]groupMemberFullInfo, error) {
	groupInfoList, err := u.getJoinedGroupListFromLocal()
	if err != nil {
		return nil, err
	}
	joinGroupMemberList := make([]groupMemberFullInfo, 0)
	for _, v := range groupInfoList {
		theGroupMemberList, err := u.getGroupAllMemberListByGroupIdFromSvr(v.GroupId)
		if err != nil {
			sdkLog(err.Error())
			continue
		}
		for _, v := range theGroupMemberList {
			joinGroupMemberList = append(joinGroupMemberList, v)
		}
	}
	return joinGroupMemberList, nil
}

func (u *UserRelated) getGroupAllMemberListByGroupIdFromSvr(groupId string) ([]groupMemberFullInfo, error) {
	var req getGroupAllMemberReq
	req.OperationID = operationIDGenerator()
	req.GroupID = groupId

	resp, err := post2Api(getGroupAllMemberListRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	sdkLog("getGroupAllMemberListRouter", getGroupAllMemberListRouter, req, string(resp))
	var stcResp groupMemberInfoResult
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		sdkLog("Unmarshal failed, ", err.Error())
		return nil, err
	}

	if stcResp.ErrCode != 0 {
		sdkLog("errcode errmsg ", stcResp.ErrCode, stcResp.ErrMsg)
		return nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Data, nil
}

func (u *UserRelated) getLocalGroupMemberListNew() ([]groupMemberFullInfo, error) {
	return u.getLocalGroupMemberList()
}

func (u *UserRelated) getLocalGroupMemberListByGroupIDNew(groupId string) ([]groupMemberFullInfo, error) {
	return u.getLocalGroupMemberListByGroupID(groupId)
}
func (u *UserRelated) insertIntoLocalGroupMemberNew(info groupMemberFullInfo) error {
	return u.insertIntoLocalGroupMember(info)
}
func (u *UserRelated) delLocalGroupMemberNew(info groupMemberFullInfo) error {
	return u.delLocalGroupMember(info)
}
func (u *UserRelated) replaceLocalGroupMemberInfoNew(info groupMemberFullInfo) error {
	return u.replaceLocalGroupMemberInfo(info)
}

func (u *UserRelated) insertIntoSelfApplyToGroupRequestNew(groupId, message string) error {
	return u.insertIntoSelfApplyToGroupRequest(groupId, message)
}
