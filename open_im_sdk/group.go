package open_im_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (g *groupListener) doGroupMsg(msg MsgData) {
	if g.listener == nil {
		sdkLog("group listener is null")
		return
	}
	go func() {
		switch msg.ContentType {
		case TransferGroupOwnerTip:
			g.doTransferGroupOwner(&msg)
		case CreateGroupTip:
			g.doCreateGroup(&msg)
		case JoinGroupTip:
			g.doJoinGroup(&msg)
		case QuitGroupTip:
			g.doQuitGroup(&msg)
		case SetGroupInfoTip:
			g.doSetGroupInfo(&msg)
		case AcceptGroupApplicationTip:
			g.doAcceptGroupApplication(&msg)
		case RefuseGroupApplicationTip:
			g.doRefuseGroupApplication(&msg)
		case KickGroupMemberTip:
			g.doKickGroupMember(&msg)
		case InviteUserToGroupTip:
			g.doInviteUserToGroup(&msg)
		default:
			sdkLog("ContentType tip failed, ", msg.ContentType)
		}
	}()
}

func (g *groupListener) doCreateGroup(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("doCreateGroup unmarshal failed", err.Error())
		return
	}
	sdkLog("doCreateGroup, ", msg, n)
	g.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	g.syncGroupMemberByGroupId(n.Detail)
	sdkLog("syncGroupMemberByGroupId ok, ", n.Detail)
	g.onGroupCreated(n.Detail)
	sdkLog("onGroupCreated callback, ", n.Detail)
}

func (g *groupListener) doJoinGroup(msg *MsgData) {

	g.syncGroupRequest()

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
	g.onReceiveJoinApplication(msg.RecvID, memberFullInfo, infoSpiltStr[1])

}

func (g *groupListener) doQuitGroup(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	g.syncJoinedGroupInfo()
	g.syncGroupMemberByGroupId(n.Detail)
	sdkLog("syncJoinedGroupInfo finish")
	sdkLog("syncGroupMemberByGroupId finish")

	var memberFullInfo groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = n.Detail

	g.onMemberLeave(n.Detail, memberFullInfo)
}

func (g *groupListener) doSetGroupInfo(msg *MsgData) {
	var n NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog("doSetGroupInfo, ", n)

	g.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")

	var groupInfo setGroupInfoReq
	err = json.Unmarshal([]byte(n.Detail), &groupInfo)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog("doSetGroupInfo ok , callback ", groupInfo.GroupId, groupInfo)
	g.onGroupInfoChanged(groupInfo.GroupId, groupInfo)
}

func (g *groupListener) doTransferGroupOwner(msg *MsgData) {
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
	g.onTransferGroupOwner(&transfer)
}
func (g *groupListener) onTransferGroupOwner(transfer *TransferGroupOwnerReq) {
	//if err := updateLocalTransferGroupOwner(transfer); err != nil {
	//	sdkLog("updateTransferGroupOwner, ", err.Error(), transfer.GroupID, transfer.OldOwner, transfer.NewOwner, transfer.OldOwner)
	//	return
	//}
	if LoginUid == transfer.NewOwner || LoginUid == transfer.OldOwner {
		g.syncGroupRequest()
	}
	g.syncGroupMemberByGroupId(transfer.GroupID)

	gInfo, err := getLocalGroupsInfoByGroupID(transfer.GroupID)
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
	g.listener.OnGroupInfoChanged(transfer.GroupID, string(bChangeInfo))
	sdkLog("onTransferGroupOwner success")
}

func (g *groupListener) doAcceptGroupApplication(msg *MsgData) {
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

	g.onAcceptGroupApplication(&acceptInfo)
}
func (g *groupListener) onAcceptGroupApplication(groupMember *GroupApplicationInfo) {
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

	//if err := insertLocalAcceptGroupApplication(&member); err != nil {
	//	sdkLog("insertAcceptGroupApplication, ", err.Error())
	//	return
	//}

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
	if LoginUid == member.UserId {
		g.syncJoinedGroupInfo()
		g.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), 1, groupMember.Info.HandledMsg)
	}
	//g.syncGroupRequest()
	g.syncGroupMemberByGroupId(groupMember.Info.GroupId)
	g.listener.OnMemberEnter(groupMember.Info.GroupId, string(bMemberListr))

	sdkLog("onAcceptGroupApplication success")
}

func (g *groupListener) doRefuseGroupApplication(msg *MsgData) {
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

	g.onRefuseGroupApplication(&refuseInfo)
}

func (g *groupListener) onRefuseGroupApplication(groupMember *GroupApplicationInfo) {
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

	if LoginUid == member.UserId {
		g.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), -1, groupMember.Info.HandledMsg)
	}

	//if err := insertLocalAcceptGroupApplication(&member); err != nil {
	//	sdkLog("insertAcceptGroupApplication, ", err.Error())
	//	return
	//}

	sdkLog("onRefuseGroupApplication success")
}

func (g *groupListener) doKickGroupMember(msg *MsgData) {
	var notification NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &notification)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	var kickReq kickGroupMemberApiReq
	err = json.Unmarshal([]byte(notification.Detail), &kickReq)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := g.getGroupMembersInfoFromLocal(kickReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(kickReq.UidListInfo) == 0 {
		sdkLog("len: ", len(opList), len(kickReq.UidListInfo))
		return
	}
	//	g.syncGroupMember()
	g.syncJoinedGroupInfo()
	g.syncGroupMemberByGroupId(kickReq.GroupID)
	g.OnMemberKicked(kickReq.GroupID, opList[0], kickReq.UidListInfo)
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

func (g *groupListener) doInviteUserToGroup(msg *MsgData) {
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

	memberList, err := g.getGroupMembersInfoTry2(inviteReq.GroupID, inviteReq.UidList)
	if err != nil {
		return
	}

	tList := make([]string, 1)
	tList = append(tList, msg.SendID)
	opList, err := g.getGroupMembersInfoTry2(inviteReq.GroupID, tList)
	sdkLog("getGroupMembersInfoFromSvr, ", inviteReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(memberList) == 0 {
		sdkLog("len: ", len(opList), len(memberList))
		return
	}
	for _, v := range inviteReq.UidList {
		if LoginUid == v {

			g.syncJoinedGroupInfo()
			sdkLog("syncJoinedGroupInfo, ", v)
			break
		}
	}

	g.syncGroupMemberByGroupId(inviteReq.GroupID)
	sdkLog("syncGroupMemberByGroupId, ", inviteReq.GroupID)
	g.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
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

func (g *groupListener) createGroup(group groupInfo, memberList []createGroupMemberInfo) (*createGroupResp, error) {
	req := createGroupReq{memberList, group.GroupName, group.Introduction, group.Notification, group.FaceUrl, operationIDGenerator()}
	resp, err := post2Api(createGroupRouter, req, token)
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

	g.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	g.syncGroupMemberByGroupId(createGroupResp.Data.GroupId)
	sdkLog("syncGroupMemberByGroupId ok")

	n := NotificationContent{
		IsDisplay:   1,
		DefaultTips: "You have joined the group chat:" + createGroupResp.Data.GroupName,
		Detail:      createGroupResp.Data.GroupId,
	}
	msg := createTextSystemMessage(n, CreateGroupTip)
	autoSendMsg(msg, "", createGroupResp.Data.GroupId, false, true, true)
	sdkLog("sendMsg, groupId: ", createGroupResp.Data.GroupId)
	return &createGroupResp, nil
}

func (g *groupListener) joinGroup(groupId, message string) error {
	req := joinGroupReq{groupId, message, operationIDGenerator()}
	resp, err := post2Api(joinGroupRouter, req, token)
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

	g.syncApplyGroupRequest()
	sdkLog("syncApplyGroupRequest ok")

	n := NotificationContent{
		IsDisplay:   1,
		DefaultTips: "User：" + LoginUid + " application to join your group",
		Detail:      groupId + "," + message,
	}
	memberList, err := g.getGroupAllMemberListByGroupIdFromSvr(groupId)
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

	msg := createTextSystemMessage(n, JoinGroupTip)
	err = autoSendMsg(msg, groupAdminUser, "", false, false, false)
	sdkLog("sendMsg ", n, groupAdminUser, err)
	return nil
}

func (g *groupListener) quitGroup(groupId string) error {
	req := quitGroupReq{groupId, operationIDGenerator()}
	resp, err := post2Api(quitGroupRouter, req, token)
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

	g.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")
	g.syncGroupMemberByGroupId(groupId) //todo
	sdkLog("syncGroupMemberByGroupId ok ", groupId)

	userInfo, err := getLoginUserInfoFromLocal()
	if err != nil {
		sdkLog("getLoginUserInfoFromLocal failed", err.Error())
		return err
	}
	n2Group := NotificationContent{
		IsDisplay:   1,
		DefaultTips: "User: " + userInfo.Name + " have quit group chat",
		Detail:      "",
	}
	msg2Group := createTextSystemMessage(n2Group, QuitGroupTip)
	err = autoSendMsg(msg2Group, "", groupId, false, true, false)
	sdkLog("sendMsg, ", n2Group, groupId, err)
	return nil
}

func (g *groupListener) getJoinedGroupListFromLocal() ([]groupInfo, error) {
	return getLocalGroupsInfo()
}

func (g *groupListener) getJoinedGroupListFromSvr() ([]groupInfo, error) {
	var req getJoinedGroupListReq
	req.OperationID = operationIDGenerator()
	sdkLog("getJoinedGroupListRouter ", getJoinedGroupListRouter, req, token)
	resp, err := post2Api(getJoinedGroupListRouter, req, token)
	if err != nil {
		fmt.Println("post api:", err)
		return nil, err
	}

	var stcResp getJoinedGroupListResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		fmt.Println("unmarshal, ", err)
		return nil, err
	}

	if stcResp.ErrCode != 0 {
		return nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Data, nil
}

func (g *groupListener) getGroupsInfo(groupIdList []string) ([]groupInfo, error) {
	req := getGroupsInfoReq{groupIdList, operationIDGenerator()}
	resp, err := post2Api(getGroupsInfoRouter, req, token)
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

func (g *groupListener) setGroupInfo(newGroupInfo setGroupInfoReq) error {
	uid, err := findLocalGroupOwnerByGroupId(newGroupInfo.GroupId)
	if err != nil {
		sdkLog("findLocalGroupOwnerByGroupId failed, ", newGroupInfo.GroupId, err.Error())
		return err
	}
	if LoginUid != uid {
		sdkLog("no permission, ", LoginUid, uid)
		return errors.New("no permission")
	}
	sdkLog("findLocalGroupOwnerByGroupId ok ", newGroupInfo.GroupId, uid)

	req := setGroupInfoReq{newGroupInfo.GroupId, newGroupInfo.GroupName, newGroupInfo.Notification, newGroupInfo.Introduction, newGroupInfo.FaceUrl, operationIDGenerator()}
	resp, err := post2Api(setGroupInfoRouter, req, token)
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

	g.syncJoinedGroupInfo()
	sdkLog("syncJoinedGroupInfo ok")

	groupInfo, err := g.getLocalGroupInfoByGroupId(newGroupInfo.GroupId)
	if err != nil {
		sdkLog("getLocalGroupInfoByGroupId", err.Error())
		return err
	}
	jsonInfo, err := json.Marshal(groupInfo)
	if err != nil {
		sdkLog("marshal failed", err.Error())
		return err
	}

	n := NotificationContent{
		IsDisplay:   1,
		DefaultTips: "Group Info has been changed",
		Detail:      string(jsonInfo),
	}
	msg := createTextSystemMessage(n, SetGroupInfoTip)
	err = autoSendMsg(msg, "", newGroupInfo.GroupId, false, false, true)
	sdkLog("sendMsg: ", n, string(jsonInfo), err)
	return nil
}

func (g *groupListener) getGroupMemberListFromSvr(groupId string, filter int32, next int32) (int32, []groupMemberFullInfo, error) {
	var req getGroupMemberListReq
	req.OperationID = operationIDGenerator()
	req.GroupID = groupId
	req.NextSeq = next
	req.Filter = filter
	resp, err := post2Api(getGroupMemberListRouter, req, token)
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

func (g *groupListener) getGroupMemberListFromLocal(groupId string, filter int32, next int32) (int32, []groupMemberFullInfo, error) {
	memberList, err := getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		return 0, nil, err
	}
	return 0, memberList, nil
}

func (g *groupListener) getGroupMembersInfoFromLocal(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	var result []groupMemberFullInfo
	localMemberList, err := g.getLocalGroupMemberListByGroupID(groupId)
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

func (g *groupListener) getGroupMembersInfoTry2(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	result, err := g.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil || len(result) == 0 {
		return g.getGroupMembersInfoFromSvr(groupId, memberList)
	} else {
		return result, err
	}
}

func (g *groupListener) getGroupMembersInfoFromSvr(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
	var req getGroupMembersInfoReq
	req.GroupID = groupId
	req.OperationID = operationIDGenerator()
	req.MemberList = memberList

	resp, err := post2Api(getGroupMembersInfoRouter, req, token)
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

func (g *groupListener) kickGroupMember(groupId string, memberList []string, reason string) ([]idResult, error) {
	var req kickGroupMemberApiReq
	req.OperationID = operationIDGenerator()
	memberListInfo, err := g.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil {
		sdkLog("getGroupMembersInfoFromLocal, ", err.Error())
		return nil, err
	}
	req.UidListInfo = memberListInfo
	req.Reason = reason
	req.GroupID = groupId
	//type KickGroupMemberReq struct {
	//	GroupID     string   `json:"groupID"`
	//	UidList     []string `json:"uidList" binding:"required"`
	//	Reason      string   `json:"reason"`
	//	OperationID string   `json:"operationID" binding:"required"`
	//}

	resp, err := post2Api(kickGroupMemberRouter, req, token)
	if err != nil {
		sdkLog("post2Api failed, ", kickGroupMemberRouter, req, err.Error())
		return nil, err
	}
	sdkLog("url: ", kickGroupMemberRouter, "req:", req, "resp: ", string(resp))

	g.syncGroupMemberByGroupId(groupId)
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
	err = autoSendKickGroupMemberTip(&req)
	sdkLog("kickGroupMember, ", groupId, memberList, reason, req)
	return sctResp.Data, nil
}

//1
func (g *groupListener) transferGroupOwner(groupId, userId string) error {
	resp, err := post2Api(transferGroupRouter, transferGroupReq{GroupID: groupId, Uid: userId, OperationID: operationIDGenerator()}, token)
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
func (g *groupListener) inviteUserToGroup(groupId string, reason string, userList []string) ([]idResult, error) {
	var req inviteUserToGroupReq
	req.GroupID = groupId
	req.OperationID = operationIDGenerator()
	req.Reason = reason
	req.UidList = userList
	resp, err := post2Api(inviteUserToGroupRouter, req, token)
	if err != nil {
		return nil, err
	}
	g.syncGroupMemberByGroupId(groupId)
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

	err = autoSendInviteUserToGroupTip(req)
	sdkLog("inviteUserToGroup, autoSendInviteUserToGroupTip", groupId, reason, userList, req, err)
	return stcResp.Data, nil
}

func (g *groupListener) getLocalGroupApplicationList(groupId string) (*groupApplicationResult, error) {
	reply, err := getOwnLocalGroupApplicationList(groupId)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
func (g *groupListener) insertIntoRequestToGroupRequest(info GroupReqListInfo) error {
	return insertIntoRequestToGroupRequest(info)
}
func (g *groupListener) delGroupRequestFromGroupRequest(info GroupReqListInfo) error {
	return delRequestFromGroupRequest(info)
}
func (g *groupListener) replaceIntoRequestToGroupRequest(info GroupReqListInfo) error {
	return replaceIntoRequestToGroupRequest(info)
}

//1
func (g *groupListener) getGroupApplicationList() (*groupApplicationResult, error) {
	resp, err := post2Api(getGroupApplicationListRouter, getGroupApplicationListReq{OperationID: operationIDGenerator()}, token)
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
func (g *groupListener) acceptGroupApplication(access *accessOrRefuseGroupApplicationReq) error {
	resp, err := post2Api(acceptGroupApplicationRouter, access, token)
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
func (g *groupListener) refuseGroupApplication(access *accessOrRefuseGroupApplicationReq) error {
	resp, err := post2Api(acceptGroupApplicationRouter, access, token)
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

func (g *groupListener) getGroupInfoByGroupId(groupId string) (groupInfo, error) {
	var gList []string
	gList = append(gList, groupId)
	rList, err := g.getGroupsInfo(gList)
	if err == nil && len(rList) == 1 {
		return rList[0], nil
	} else {
		return groupInfo{}, nil
	}

}

type groupListener struct {
	listener OnGroupListener
	//	ch       chan cmd2Value
}

/*
func (g *groupListener) getCh() chan cmd2Value {
	return g.ch
}

func (g *groupListener) work(c2v cmd2Value) {
	//triggered by conversation
	switch c2v.Cmd {
	case CmdGroup:
		if g.listener == nil {
			return
		}
		g.doCmdGroup(c2v)
	}
}
*/

func (g *groupListener) createGroupCallback(node updateGroupNode) {
	// member list to json
	jsonMemberList, err := json.Marshal(node.Args.(createGroupArgs).initMemberList)
	if err != nil {
		return
	}
	groupManager.listener.OnMemberEnter(node.groupId, string(jsonMemberList))
	groupManager.listener.OnGroupCreated(node.groupId)
}

func (g *groupListener) joinGroupCallback(node updateGroupNode) {
	args := node.Args.(joinGroupArgs)
	jsonApplyUser, err := json.Marshal(args.applyUser)
	if err != nil {
		return
	}
	groupManager.listener.OnReceiveJoinApplication(node.groupId, string(jsonApplyUser), args.reason)
}

func (g *groupListener) quitGroupCallback(node updateGroupNode) {
	args := node.Args.(quiteGroupArgs)
	jsonUser, err := json.Marshal(args.quiteUser)
	if err != nil {
		return
	}
	groupManager.listener.OnMemberLeave(node.groupId, string(jsonUser))
}

func (g *groupListener) setGroupInfoCallback(node updateGroupNode) {
	args := node.Args.(setGroupInfoArgs)
	jsonGroup, err := json.Marshal(args.group)
	if err != nil {
		return
	}
	groupManager.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
}

func (g *groupListener) kickGroupMemberCallback(node updateGroupNode) {
	args := node.Args.(kickGroupAgrs)
	jsonop, err := json.Marshal(args.op)
	if err != nil {
		return
	}

	jsonKickedList, err := json.Marshal(args.kickedList)
	if err != nil {
		return
	}

	groupManager.listener.OnMemberKicked(node.groupId, string(jsonop), string(jsonKickedList))
}

func (g *groupListener) transferGroupOwnerCallback(node updateGroupNode) {
	args := node.Args.(transferGroupArgs)

	group, err := g.getGroupInfoByGroupId(node.groupId)
	if err != nil {
		return
	}
	group.OwnerId = args.newOwner.UserId

	jsonGroup, err := json.Marshal(group)
	if err != nil {
		return
	}
	groupManager.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
}

func (g *groupListener) inviteUserToGroupCallback(node updateGroupNode) {
	args := node.Args.(inviteUserToGroupArgs)
	jsonInvitedList, err := json.Marshal(args.invited)
	if err != nil {
		return
	}
	jsonOp, err := json.Marshal(args.op)
	if err != nil {
		return
	}
	groupManager.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonInvitedList))
}

func (g *groupListener) GroupApplicationProcessedCallback(node updateGroupNode, process int32) {
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
		if v.member.UserId == LoginUid {
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
		groupManager.listener.OnApplicationProcessed(node.groupId, string(jsonOp), process, processed.applyList[idx].reason)
	}

	if process == 1 {
		jsonOp, err := json.Marshal(processed.op)
		if err != nil {
			return
		}
		g.listener.OnMemberInvited(node.groupId, string(jsonOp), string(jsonApplyList))
	}
}

func (g *groupListener) acceptGroupApplicationCallback(node updateGroupNode) {
	g.GroupApplicationProcessedCallback(node, 1)
}

func (g *groupListener) refuseGroupApplicationCallback(node updateGroupNode) {
	g.GroupApplicationProcessedCallback(node, -1)
}

func (g *groupListener) syncGroupRequest() {
	groupRequestOnServerResp, err := g.getGroupApplicationList()
	if err != nil {
		sdkLog("groupRequestOnServerResp failed", err.Error())
		return
	}
	groupRequestOnServer := groupRequestOnServerResp.GroupApplicationList
	groupRequestOnServerInterface := make([]diff, 0)
	for _, v := range groupRequestOnServer {
		groupRequestOnServerInterface = append(groupRequestOnServerInterface, v)
	}

	groupRequestOnLocalResp, err := g.getLocalGroupApplicationList("")
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
		err = g.insertIntoRequestToGroupRequest(groupRequestOnServer[index])
		if err != nil {
			sdkLog("insertIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		sdkLog("insertIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

	for _, index := range bInANot {
		err = g.delGroupRequestFromGroupRequest(groupRequestOnLocal[index])
		if err != nil {
			sdkLog("delGroupRequestFromGroupRequest failed", err.Error())
			continue
		}
		sdkLog("delGroupRequestFromGroupRequest ", groupRequestOnLocal[index])
	}
	for _, index := range sameA {
		if err = g.replaceIntoRequestToGroupRequest(groupRequestOnServer[index]); err != nil {
			sdkLog("replaceIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		sdkLog("replaceIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

}

func (g *groupListener) syncApplyGroupRequest() {

}

func (g *groupListener) syncJoinedGroupInfo() {
	groupListOnServer, err := g.getJoinedGroupListFromSvr()
	if err != nil {
		sdkLog("groupListOnServer failed", err.Error())
		return
	}
	groupListOnServerInterface := make([]diff, 0)
	for _, v := range groupListOnServer {
		groupListOnServerInterface = append(groupListOnServerInterface, v)
	}

	groupListOnLocal, err := g.getJoinedGroupListFromLocal()
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
		err = g.insertIntoLocalGroupInfo(groupListOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupInfo failed", err.Error(), groupListOnServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		err = g.delLocalGroupInfo(groupListOnLocal[index].GroupId)
		if err != nil {
			sdkLog("delLocalGroupInfo failed", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		if err = g.replaceLocalGroupInfo(groupListOnServer[index]); err != nil {
			sdkLog("replaceLocalGroupInfo failed", err.Error())
			continue
		}
	}
}

func (g *groupListener) getLocalGroupsInfo() ([]groupInfo, error) {
	localGroupsInfo, err := getLocalGroupsInfo()
	if err != nil {
		return nil, err
	}
	groupId2Owner := make(map[string]string)
	groupId2MemberNum := make(map[string]uint32)
	for index, v := range localGroupsInfo {
		if _, ok := groupId2Owner[v.GroupId]; !ok {
			ownerId, err := findLocalGroupOwnerByGroupId(v.GroupId)
			if err != nil {
				sdkLog(err.Error())
			}
			groupId2Owner[v.GroupId] = ownerId
		}
		localGroupsInfo[index].OwnerId = groupId2Owner[v.GroupId]
		if _, ok := groupId2MemberNum[v.GroupId]; !ok {
			num, err := getLocalGroupMemberNumByGroupId(v.GroupId)
			if err != nil {
				sdkLog(err.Error())
			}
			groupId2MemberNum[v.GroupId] = uint32(num)
		}
		localGroupsInfo[index].MemberCount = groupId2MemberNum[v.GroupId]
	}
	return localGroupsInfo, nil
}
func (g *groupListener) getLocalGroupInfoByGroupId(groupId string) (*groupInfo, error) {
	return getLocalGroupsInfoByGroupID(groupId)
}
func (g *groupListener) insertIntoLocalGroupInfo(info groupInfo) error {
	return insertIntoLocalGroupInfo(info)
}
func (g *groupListener) delLocalGroupInfo(groupId string) error {
	return delLocalGroupInfo(groupId)
}
func (g *groupListener) replaceLocalGroupInfo(info groupInfo) error {
	return replaceLocalGroupInfo(info)
}

func (g *groupListener) syncGroupMemberByGroupId(groupId string) {
	groupMemberOnServer, err := g.getGroupAllMemberListByGroupIdFromSvr(groupId)
	if err != nil {
		sdkLog("syncGroupMemberByGroupId failed", err.Error())
		return
	}
	sdkLog("getGroupAllMemberListByGroupIdFromSvr, ", groupId, len(groupMemberOnServer))

	groupMemberOnServerInterface := make([]diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := g.getLocalGroupMemberListByGroupID(groupId)
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
		err = g.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range bInANot {
		err = g.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = g.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			sdkLog("replaceLocalGroupMemberInfo failed", err.Error())
			continue
		}
	}

}

func (g *groupListener) syncJoinedGroupMember() {
	groupMemberOnServer, err := g.getJoinGroupAllMemberList()
	if err != nil {
		sdkLog("getJoinGroupAllMemberList failed", err.Error())
		return
	}
	groupMemberOnServerInterface := make([]diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := g.getLocalGroupMemberList()
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
		err = g.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			sdkLog("insertIntoLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range bInANot {
		err = g.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = g.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			sdkLog(err.Error())
			continue
		}
	}

}

func (g *groupListener) getJoinGroupAllMemberList() ([]groupMemberFullInfo, error) {
	groupInfoList, err := g.getJoinedGroupListFromLocal()
	if err != nil {
		return nil, err
	}
	joinGroupMemberList := make([]groupMemberFullInfo, 0)
	for _, v := range groupInfoList {
		theGroupMemberList, err := g.getGroupAllMemberListByGroupIdFromSvr(v.GroupId)
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

func (g *groupListener) getGroupAllMemberListByGroupIdFromSvr(groupId string) ([]groupMemberFullInfo, error) {
	var req getGroupAllMemberReq
	req.OperationID = operationIDGenerator()
	req.GroupID = groupId

	resp, err := post2Api(getGroupAllMemberListRouter, req, token)
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

func (g *groupListener) getLocalGroupMemberList() ([]groupMemberFullInfo, error) {
	return getLocalGroupMemberList()
}

func (g *groupListener) getLocalGroupMemberListByGroupID(groupId string) ([]groupMemberFullInfo, error) {
	return getLocalGroupMemberListByGroupID(groupId)
}
func (g *groupListener) insertIntoLocalGroupMember(info groupMemberFullInfo) error {
	return insertIntoLocalGroupMember(info)
}
func (g *groupListener) delLocalGroupMember(info groupMemberFullInfo) error {
	return delLocalGroupMember(info)
}
func (g *groupListener) replaceLocalGroupMemberInfo(info groupMemberFullInfo) error {
	return replaceLocalGroupMemberInfo(info)
}

func (g *groupListener) insertIntoSelfApplyToGroupRequest(groupId, message string) error {
	return insertIntoSelfApplyToGroupRequest(groupId, message)
}

/*
func (g *groupListener) doCmdGroup(c2v cmd2Value) {
	node := c2v.Value.(updateGroupNode)
	//	 2: apply to join group; 3:quit group; 4:set group info; 5:kick group member;6:transfer group owner;7:invite user to group 8:accept group application; 9:refuse group application;
	switch node.Action {
	case GroupActionCreateGroup: //receivers : creator and init member
		g.doGroupJoinedSync(node)
		g.createGroupCallback(node) //callback  onMemberEnter and onGroupCreated
	case GroupActionApplyJoinGroup: //receiver : group creator
		g.doGroupApplicationSync(node)
		g.joinGroupCallback(node) //callback OnReceiveJoinApplication
	case GroupActionQuitGroup: //receiver: all group member include operator
		quit := node.Args.(quiteGroupArgs)
		if quit.quiteUser.UserId == LoginUid {
			g.doGroupJoinedSync(node)
		} else {
			g.doGroupMemberSync(node)
		}
		g.quitGroupCallback(node) //callback OnMemberLeave
	case GroupActionSetGroupInfo: //receiver: all group member
		g.doGroupInfoSync(node)
		g.setGroupInfoCallback(node) //callback onGroupInfoChanged
	case GroupActionKickGroupMember: //receiver: all group member include kicked
		kick := node.Args.(kickGroupAgrs)
		if kick.op.UserId == LoginUid {
			g.doGroupJoinedSync(node)
		} else {
			g.doGroupMemberSync(node)
		}
		g.kickGroupMemberCallback(node) //callback OnMemberKicked
	case GroupActionTransferGroupOwner: // receiver : all group member
		g.doGroupInfoSync(node)
		g.transferGroupOwnerCallback(node) //callback OnGroupInfoChanged
	case GroupActionInviteUserToGroup: //receiver: group owner
		invite := node.Args.(inviteUserToGroupArgs)
		var flag = 0
		for _, v := range invite.invited {
			if v.UserId == LoginUid {
				flag = 1
				break
			}
		}
		if flag == 1 {
			g.doGroupJoinedSync(node)
		} else {
			g.doGroupMemberSync(node)
		}
		g.inviteUserToGroupCallback(node) //callback  OnReceiveJoinApplication
	case GroupActionAcceptGroupApplication:
		accept := node.Args.(applyGroupProcessedArgs) //receiver : all group member
		var flag = 0
		for _, v := range accept.applyList {
			if v.member.UserId == LoginUid {
				flag = 1
				break
			}
		}
		if flag == 1 {
			g.doGroupApplicationSync(node)
			g.doGroupJoinedSync(node)
		} else if accept.op.UserId == LoginUid {
			g.doGroupMemberSync(node)
			g.doGroupApplicationSync(node)
		} else {
			g.doGroupMemberSync(node)
		}
		g.acceptGroupApplicationCallback(node)
	case GroupActionRefuseGroupApplication:
		g.doGroupApplicationSync(node)
		g.refuseGroupApplicationCallback(node)
	}
}
*/

/*
//svr->local
func (g *groupListener) doGroupApplicationSync(node updateGroupNode) {

}

//svr->local
func (g *groupListener) doGroupMemberSync(node updateGroupNode) {

}

//svr->local
func (g *groupListener) doGroupInfoSync(node updateGroupNode) {

}

//svr->local
func (g *groupListener) doGroupJoinedSync(node updateGroupNode) {

}
*/
