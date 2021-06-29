package open_im_sdk

import (
	"encoding/json"
	"errors"
	"fmt"
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
		case GroupApplicationResponseTip:
			g.doGroupApplicationResponse(&msg)
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
		case AcceptGroupApplicationResultTip:
			g.doAcceptGroupApplicationResult(&msg)
		case RefuseGroupApplicationResultTip:
			g.doRefuseGroupApplicationResult(&msg)
		default:
			sdkLog("tip failed, ", msg.ContentType)
		}
	}()
}

func (g *groupListener) doCreateGroup(msg *MsgData) {
	g.onGroupCreated(msg.RecvID)
}
func (g *groupListener) doJoinGroup(msg *MsgData) {
	var joinGroupReq joinGroupReq
	err := json.Unmarshal([]byte(msg.Content), &joinGroupReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	var memberFullInfo groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = joinGroupReq.GroupID
	g.onReceiveJoinApplication(msg.RecvID, memberFullInfo, joinGroupReq.Message)

}

func (g *groupListener) doQuitGroup(msg *MsgData) {
	var quitGroupReq quitGroupReq
	err := json.Unmarshal([]byte(msg.Content), &quitGroupReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	var memberFullInfo groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = quitGroupReq.GroupID

	g.onMemberLeave(quitGroupReq.GroupID, memberFullInfo)
}
func (g *groupListener) doSetGroupInfo(msg *MsgData) {
	var setGroupInfoReq setGroupInfoReq
	err := json.Unmarshal([]byte(msg.Content), &setGroupInfoReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	sdkLog(setGroupInfoReq)
	g.onGroupInfoChanged(setGroupInfoReq.GroupId, setGroupInfoReq)
}
func (g *groupListener) doTransferGroupOwner(msg *MsgData) {
	var transfer TransferGroupOwnerReq
	err := json.Unmarshal([]byte(msg.Content), &transfer)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	g.onTransferGroupOwner(&transfer)

}
func (g *groupListener) onTransferGroupOwner(transfer *TransferGroupOwnerReq) {
	// need modify sql
	g.listener.OnTransferGroupOwner(transfer.GroupID, transfer.OldOwner, transfer.NewOwner)
}

func (g *groupListener) doGroupApplicationResponse(msg *MsgData) {
	// do nothing
	return
}
func (g *groupListener) doAcceptGroupApplication(msg *MsgData) {
	var groupMember GroupApplicationResponseReq
	err := json.Unmarshal([]byte(msg.Content), &groupMember)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	g.onAcceptGroupApplication(&groupMember)
}
func (g *groupListener) onAcceptGroupApplication(groupMember *GroupApplicationResponseReq) {
	// need modify sql

	member := groupMemberFullInfo{
		GroupId:  groupMember.GroupID,
		Role:     0,
		JoinTime: uint64(groupMember.AddTime),
	}
	if groupMember.ToUserID == "0" {
		member.UserId = groupMember.FromUserID
		member.NickName = groupMember.FromUserNickName
		member.FaceUrl = groupMember.FromUserFaceUrl
	} else {
		member.UserId = groupMember.ToUserID
		member.NickName = groupMember.ToUserNickName
		member.FaceUrl = groupMember.ToUserFaceUrl
	}

	var memberList []groupMemberFullInfo
	memberList = append(memberList, member)
	bMemberListr, err := json.Marshal(memberList)
	if err != nil {
		sdkLog("onAcceptGroupApplication", err.Error())
		return
	}
	g.listener.OnMemberEnter(groupMember.GroupID, string(bMemberListr))
}

func (g *groupListener) doRefuseGroupApplication(msg *MsgData) {
	// do nothing
	return
}

func (g *groupListener) doAcceptGroupApplicationResult(msg *MsgData) {
	var groupMember AgreeOrRejectGroupMember
	err := json.Unmarshal([]byte(msg.Content), &groupMember)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	g.onAcceptGroupApplicationResult(&groupMember)
}
func (g *groupListener) onAcceptGroupApplicationResult(groupMember *AgreeOrRejectGroupMember) {

	op := groupMemberFullInfo{
		GroupId:  groupMember.GroupId,
		UserId:   groupMember.UserId,
		Role:     groupMember.Role,
		JoinTime: groupMember.JoinTime,
		NickName: groupMember.NickName,
		FaceUrl:  groupMember.FaceUrl,
	}

	bOp, err := json.Marshal(op)
	if err != nil {
		sdkLog("Marshal, ", err.Error())
		return
	}
	g.listener.OnApplicationProcessed(groupMember.GroupId, string(bOp), 1, groupMember.Reason)
}

func (g *groupListener) doRefuseGroupApplicationResult(msg *MsgData) {
	var groupMember AgreeOrRejectGroupMember
	err := json.Unmarshal([]byte(msg.Content), &groupMember)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	g.onRefuseGroupApplicationResult(&groupMember)
}
func (g *groupListener) onRefuseGroupApplicationResult(groupMember *AgreeOrRejectGroupMember) {

	op := groupMemberFullInfo{
		GroupId:  groupMember.GroupId,
		UserId:   groupMember.UserId,
		Role:     groupMember.Role,
		JoinTime: groupMember.JoinTime,
		NickName: groupMember.NickName,
		FaceUrl:  groupMember.FaceUrl,
	}

	bOp, err := json.Marshal(op)
	if err != nil {
		sdkLog("Marshal, ", err.Error())
		return
	}
	g.listener.OnApplicationProcessed(groupMember.GroupId, string(bOp), -1, groupMember.Reason)
}

func (g *groupListener) doKickGroupMember(msg *MsgData) {
	var kickReq KickGroupMemberReq
	err := json.Unmarshal([]byte(msg.Content), &kickReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	kickReq.Op = msg.SendID
	memberList, err := g.getGroupMembersInfo(kickReq.GroupID, kickReq.UidList)
	if err != nil {
		return
	}
	tList := make([]string, 1)
	tList = append(tList, kickReq.Op)
	opList, err := g.getGroupMembersInfo(kickReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(memberList) == 0 {
		sdkLog("len: ", len(opList), len(memberList))
		return
	}
	groupManager.OnMemberKicked(kickReq.GroupID, opList[0], memberList)
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
	var inviteReq InviteUserToGroupReq
	err := json.Unmarshal([]byte(msg.Content), &inviteReq)
	if err != nil {
		sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	inviteReq.Op = msg.SendID
	memberList, err := g.getGroupMembersInfo(inviteReq.GroupID, inviteReq.UidList)
	if err != nil {
		return
	}
	tList := make([]string, 1)
	tList = append(tList, inviteReq.Op)
	opList, err := g.getGroupMembersInfo(inviteReq.GroupID, tList)
	if err != nil {
		return
	}
	if len(opList) == 0 || len(memberList) == 0 {
		sdkLog("len: ", len(opList), len(memberList))
		return
	}
	groupManager.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
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

func (g *groupListener) createGroup(group groupInfo, memberList []createGroupMemberInfo) error {
	resp, err := post2Api(createGroupRouter, createGroupReq{memberList, group.GroupName, group.Introduction, group.Notification, group.FaceUrl, operationIDGenerator()}, token)
	if err != nil {
		return err
	}
	var createGroupResp createGroupResp
	if err = json.Unmarshal(resp, &createGroupResp); err != nil {
		return err
	}
	if createGroupResp.ErrCode != 0 {
		return errors.New(createGroupResp.ErrMsg)
	}
	return nil
}

func (g *groupListener) joinGroup(groupId, message string) error {
	resp, err := post2Api(joinGroupRouter, joinGroupReq{groupId, message, operationIDGenerator()}, token)
	if err != nil {
		return err
	}
	var commonResp commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		return err
	}
	if commonResp.ErrCode != 0 {
		return errors.New(commonResp.ErrMsg)
	}
	return nil
}

func (g *groupListener) quitGroup(groupId string) error {
	resp, err := post2Api(quitGroupRouter, quitGroupReq{groupId, operationIDGenerator()}, token)
	if err != nil {
		return err
	}
	var commonResp commonResp
	err = json.Unmarshal(resp, &commonResp)
	if err != nil {
		sdkLog(err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		return errors.New(commonResp.ErrMsg)
	}
	return nil
}

func (g *groupListener) getJoinedGroupList() ([]groupInfo, error) {
	var req getJoinedGroupListReq
	req.OperationID = operationIDGenerator()
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
	resp, err := post2Api(getGroupsInfoRouter, getGroupsInfoReq{groupIdList, operationIDGenerator()}, token)
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

//1
func (g *groupListener) setGroupInfo(newGroupInfo groupInfo) error {
	resp, err := post2Api(setGroupInfoRouter, setGroupInfoReq{newGroupInfo.GroupId, newGroupInfo.GroupName, newGroupInfo.Notification, newGroupInfo.Introduction, newGroupInfo.FaceUrl, operationIDGenerator()}, token)
	if err != nil {
		return err
	}
	var commonResp commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		return err
	}
	if commonResp.ErrCode != 0 {
		return errors.New(commonResp.ErrMsg)
	}
	return nil
}

func (g *groupListener) getGroupMemberList(groupId string, filter int32, next int32) (int32, []groupMemberFullInfo, error) {
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
		return 0, nil, err
	}

	if stcResp.ErrCode != 0 {
		return 0, nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Nextseq, stcResp.Data, nil
}

//1
func (g *groupListener) getGroupMembersInfo(groupId string, memberList []string) ([]groupMemberFullInfo, error) {
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
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		return nil, errors.New(sctResp.ErrMsg)
	}
	return sctResp.Data, nil
}

//1
func (g *groupListener) kickGroupMember(groupId string, memberList []string, reason string) ([]idResult, error) {
	var req kickGroupMemberReq
	req.OperationID = operationIDGenerator()
	req.UidList = memberList
	req.Reason = reason
	req.GroupID = groupId
	resp, err := post2Api(kickGroupMemberRouter, req, token)
	if err != nil {
		return nil, err
	}
	sdkLog("url: ", kickGroupMemberRouter, "req:", req, "resp: ", string(resp))

	var sctResp kickGroupMemberResp
	err = json.Unmarshal(resp, &sctResp)
	if err != nil {
		sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		sdkLog("resp failed, ", sctResp.ErrCode, sctResp.ErrMsg)
		return nil, errors.New(sctResp.ErrMsg)
	}
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
	var stcResp inviteUserToGroupResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		fmt.Println("sssssssssssssssssssss")
		return nil, err
	}
	return stcResp.Data, nil
}

//1
func (g *groupListener) getGroupApplicationList() (*groupApplicationResult, error) {
	resp, err := post2Api(getGroupApplicationListRouter, getGroupApplicationListReq{OperationID: operationIDGenerator()}, token)
	if err != nil {
		return nil, err
	}

	var ret getGroupApplicationListResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return nil, err
	}
	if ret.ErrCode != 0 {
		return nil, errors.New(ret.ErrMsg)
	}

	return &ret.Data, nil
}

//1
func (g *groupListener) acceptGroupApplication(application groupApplication, reason string) error {
	var access accessOrRefuseGroupApplicationReq
	access.OperationID = operationIDGenerator()
	access.GroupId = application.GroupId
	access.FromUser = application.FromUser
	access.FromUserNickName = application.FromUserNickName
	access.FromUserFaceUrl = application.FromUserFaceUrl
	access.ToUser = application.ToUser
	access.AddTime = application.AddTime
	access.RequestMsg = application.RequestMsg
	access.HandledMsg = reason
	access.Type = application.Type
	access.HandleStatus = 2
	access.HandleResult = 1

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
func (g *groupListener) refuseGroupApplication(application groupApplication, reason string) error {
	var access accessOrRefuseGroupApplicationReq
	access.OperationID = operationIDGenerator()
	access.GroupId = application.GroupId
	access.FromUser = application.FromUser
	access.FromUserNickName = application.FromUserNickName
	access.FromUserFaceUrl = application.FromUserFaceUrl
	access.ToUser = application.ToUser
	access.AddTime = application.AddTime
	access.RequestMsg = application.RequestMsg
	access.HandledMsg = reason
	access.Type = application.Type
	access.HandleStatus = 2
	access.HandleResult = 0

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
