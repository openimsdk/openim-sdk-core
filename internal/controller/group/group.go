package group

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	ws "open_im_sdk/internal/controller/interaction"
	"open_im_sdk/internal/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"strings"
)


type Group struct {
	listener OnGroupListener
	token          string
	loginUserID    string
	db             *db.DataBase
	p              *ws.PostApi
}



func (u *Group) doGroupMsg(msg * server_api_params.MsgData) {
	if u.listener == nil {
		utils.sdkLog("group listener is null")
		return
	}
	if msg.SendID == u.loginUserID && msg.SenderPlatformID == u.SvrConf.Platform {
		utils.sdkLog("sync msg ", msg)
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
			utils.sdkLog("ContentType tip failed, ", msg.ContentType)
		}
	}()
}

func (u *Group) doCreateGroup(msg *server_api_params.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		return
	}
	u.syncJoinedGroupInfo()

	u.syncGroupMemberByGroupId(n.Detail)
	u.onGroupCreated(n.Detail)
}

func (u *Group) doJoinGroup(msg *server_api_params.MsgData) {

	u.syncGroupRequest()

	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	infoSpiltStr := strings.Split(n.Detail, ",")
	var memberFullInfo open_im_sdk.groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = infoSpiltStr[0]
	u.onReceiveJoinApplication(msg.RecvID, memberFullInfo, infoSpiltStr[1])

}

func (u *Group) doQuitGroup(msg *server_api_params.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}

	utils.sdkLog("syncJoinedGroupInfo start")
	u.syncJoinedGroupInfo()
	utils.sdkLog("syncJoinedGroupInfo end")
	u.syncGroupMemberByGroupId(n.Detail)
	utils.sdkLog("syncJoinedGroupInfo finish")
	utils.sdkLog("syncGroupMemberByGroupId finish")

	var memberFullInfo open_im_sdk.groupMemberFullInfo
	memberFullInfo.UserId = msg.SendID
	memberFullInfo.GroupId = n.Detail

	u.onMemberLeave(n.Detail, memberFullInfo)
}

func (u *Group) doSetGroupInfo(msg *server_api_params.MsgData) {
	var n utils.NotificationContent
	err := json.Unmarshal([]byte(msg.Content), &n)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	utils.sdkLog("doSetGroupInfo, ", n)

	u.syncJoinedGroupInfo()
	utils.sdkLog("syncJoinedGroupInfo ok")

	var groupInfo open_im_sdk.setGroupInfoReq
	err = json.Unmarshal([]byte(n.Detail), &groupInfo)
	if err != nil {
		utils.sdkLog("unmarshal, ", err.Error(), msg.Content)
		return
	}
	utils.sdkLog("doSetGroupInfo ok , callback ", groupInfo.GroupId, groupInfo)
	u.onGroupInfoChanged(groupInfo.GroupId, groupInfo)
}

func (u *Group) doTransferGroupOwner(msg *server_api_params.MsgData) {
	utils.sdkLog("doTransferGroupOwner start...")
	var transfer server_api_params.TransferGroupOwnerReq
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

func (u *Group) doAcceptGroupApplication(msg *server_api_params.MsgData) {
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
		u.syncJoinedGroupInfo()
		u.listener.OnApplicationProcessed(groupMember.Info.GroupId, string(bOp), 1, groupMember.Info.HandledMsg)
	}
	//g.syncGroupRequest()
	u.syncGroupMemberByGroupId(groupMember.Info.GroupId)
	u.listener.OnMemberEnter(groupMember.Info.GroupId, string(bMemberListr))

	utils.sdkLog("onAcceptGroupApplication success")
}

func (u *Group) doRefuseGroupApplication(msg *server_api_params.MsgData) {
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

func (u *Group) doKickGroupMember(msg *server_api_params.MsgData) {
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
	u.syncJoinedGroupInfo()
	u.syncGroupMemberByGroupId(kickReq.GroupID)
	//u.syncJoinedGroupInfo()
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

func (u *Group) doInviteUserToGroup(msg *open_im_sdk.MsgData) {
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

			u.syncJoinedGroupInfo()
			utils.sdkLog("syncJoinedGroupInfo, ", v)
			break
		}
	}

	u.syncGroupMemberByGroupId(inviteReq.GroupID)
	utils.sdkLog("syncGroupMemberByGroupId, ", inviteReq.GroupID)
	u.OnMemberInvited(inviteReq.GroupID, opList[0], memberList)
}

func (g *Group) onGroupCreated(groupID string) {
	g.listener.OnGroupCreated(groupID)
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

func (u *Group) createGroup(callback common.Base, group sdk_params_callback.CreateGroupBaseInfoParam,
	memberList sdk_params_callback.CreateGroupMemberRoleParam, operationID string) *sdk_params_callback.CreateGroupCallback {
	apiReq := server_api_params.CreateGroupReq{}
	apiReq.OperationID = operationID
	apiReq.OwnerUserID = u.loginUserID
	apiReq.GroupName = group.GroupName
	apiReq.GroupType = group.GroupType
	apiReq.MemberList = memberList
	commData := u.p.PostFatalCallback(callback, constant.CreateGroupRouter, apiReq, u.token)
	realData := server_api_params.CreateGroupResp{}
	err := mapstructure.Decode(commData.Data, &realData.GroupInfo)
	if err != nil{
		callback.OnError(constant.ErrData.ErrCode, constant.ErrData.ErrMsg)
		return nil
	}
	u.syncJoinedGroupInfo()
	u.syncGroupMemberByGroupId(realData.GroupInfo.GroupID)
	return &sdk_params_callback.CreateGroupCallback{realData.GroupInfo}
}

func (u *Group) joinGroup(groupId, message string, callback common.Base, operationID string) error {
	req := open_im_sdk.joinGroupReq{groupId, message, utils.operationIDGenerator()}
	resp, err := utils.post2Api(open_im_sdk.joinGroupRouter, req, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", err.Error(), open_im_sdk.joinGroupRouter, req)
		return err
	}
	var commonResp open_im_sdk.commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		utils.sdkLog("Unmarshal", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		utils.sdkLog("commonResp err", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	utils.sdkLog("psot2api ok", open_im_sdk.joinGroupRouter, req, commonResp)

	u.syncApplyGroupRequest()
	utils.sdkLog("syncApplyGroupRequest ok")

	memberList, err := u.getGroupAllMemberListByGroupIdFromSvr(groupId)
	if err != nil {
		utils.sdkLog("getGroupAllMemberListByGroupIdFromSvr failed", err.Error())
		return err
	}

	var groupAdminUser string
	for _, v := range memberList {
		if v.Role == 1 {
			groupAdminUser = v.UserId
			break
		}
	}
	utils.sdkLog("get admin from svr ok ", groupId, groupAdminUser)
	return nil
}

func (u *Group) quitGroup(groupId string, callback common.Base, operationID string) error {
	req := open_im_sdk.quitGroupReq{groupId, utils.operationIDGenerator()}
	resp, err := utils.post2Api(open_im_sdk.quitGroupRouter, req, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.quitGroupRouter, req)
		return err
	}
	var commonResp open_im_sdk.commonResp
	err = json.Unmarshal(resp, &commonResp)
	if err != nil {
		utils.sdkLog("unmarshal", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		utils.sdkLog("errcode, errmsg", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	utils.sdkLog("post2Api ok ", open_im_sdk.quitGroupRouter, req, commonResp)

	u.syncJoinedGroupInfo()
	utils.sdkLog("syncJoinedGroupInfo ok")
	u.syncGroupMemberByGroupId(groupId) //todo
	utils.sdkLog("syncGroupMemberByGroupId ok ", groupId)
	return nil
}

func (u *Group) getJoinedGroupListFromLocal() ([]open_im_sdk.groupInfo, error) {
	return u.getLocalGroupsInfo()
}

func (u *Group) getJoinedGroupListFromSvr() ([]open_im_sdk.groupInfo, error) {
	var req open_im_sdk.getJoinedGroupListReq
	req.OperationID = utils.operationIDGenerator()
	utils.sdkLog("getJoinedGroupListRouter ", open_im_sdk.getJoinedGroupListRouter, req, u.token)
	resp, err := utils.post2Api(open_im_sdk.getJoinedGroupListRouter, req, u.token)
	if err != nil {
		utils.sdkLog("post api:", err)
		return nil, err
	}

	var stcResp open_im_sdk.getJoinedGroupListResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		utils.sdkLog("unmarshal, ", err)
		return nil, err
	}

	if stcResp.ErrCode != 0 {
		return nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Data, nil
}

func (u *Group) getGroupsInfo(groupIdList []string, callback common.Base, operationID string) ([]open_im_sdk.groupInfo, error) {
	req := open_im_sdk.getGroupsInfoReq{groupIdList, utils.operationIDGenerator()}
	resp, err := utils.post2Api(open_im_sdk.getGroupsInfoRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	var getGroupsInfoResp open_im_sdk.getGroupsInfoResp
	err = json.Unmarshal(resp, &getGroupsInfoResp)
	if err != nil {
		return nil, err
	}
	return getGroupsInfoResp.Data, nil
}

func (u *Group) setGroupInfo(newGroupInfo open_im_sdk.setGroupInfoReq, callback common.Base, operationID string) error {
	g, err := u._getGroupInfoByGroupID(newGroupInfo.GroupId)
	if err != nil {
		utils.sdkLog("findLocalGroupOwnerByGroupId failed, ", newGroupInfo.GroupId, err.Error())
		return err
	}
	if u.loginUserID != g.OwnerUserID {
		utils.sdkLog("no permission, ", u.loginUserID, g.OwnerUserID)
		return errors.New("no permission")
	}
	utils.sdkLog("findLocalGroupOwnerByGroupId ok ", newGroupInfo.GroupId, g.OwnerUserID)

	req := open_im_sdk.setGroupInfoReq{newGroupInfo.GroupId, newGroupInfo.GroupName, newGroupInfo.Notification, newGroupInfo.Introduction, newGroupInfo.FaceUrl, utils.operationIDGenerator()}
	resp, err := utils.post2Api(open_im_sdk.setGroupInfoRouter, req, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.setGroupInfoRouter, req, err.Error())
		return err
	}
	var commonResp open_im_sdk.commonResp
	if err = json.Unmarshal(resp, &commonResp); err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return err
	}
	if commonResp.ErrCode != 0 {
		utils.sdkLog("errcode errmsg: ", commonResp.ErrCode, commonResp.ErrMsg)
		return errors.New(commonResp.ErrMsg)
	}
	utils.sdkLog("post2Api ok, ", open_im_sdk.setGroupInfoRouter, req, commonResp)

	u.syncJoinedGroupInfo()
	utils.sdkLog("syncJoinedGroupInfo ok")
	return nil
}

func (u *Group) getGroupMemberListFromSvr(groupId string, filter int32, next int32) (int32, []open_im_sdk.groupMemberFullInfo, error) {
	var req open_im_sdk.getGroupMemberListReq
	req.OperationID = utils.operationIDGenerator()
	req.GroupID = groupId
	req.NextSeq = next
	req.Filter = filter
	resp, err := utils.post2Api(open_im_sdk.getGroupMemberListRouter, req, u.token)
	if err != nil {
		return 0, nil, err
	}
	var stcResp open_im_sdk.groupMemberInfoResult
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return 0, nil, err
	}

	if stcResp.ErrCode != 0 {
		utils.sdkLog("errcode, errmsg: ", stcResp.ErrCode, stcResp.ErrMsg)
		return 0, nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Nextseq, stcResp.Data, nil
}

func (u *Group) getGroupMemberListFromLocal(groupId string, filter int32, next int32) (int32, []open_im_sdk.groupMemberFullInfo, error) {
	memberList, err := u.getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		return 0, nil, err
	}
	return 0, memberList, nil
}

func (u *open_im_sdk) getGroupMembersInfoFromLocal(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
	var result []open_im_sdk.groupMemberFullInfo
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

func (u *Group) getGroupMembersInfoTry2(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
	result, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil || len(result) == 0 {
		return u.getGroupMembersInfoFromSvr(groupId, memberList)
	} else {
		return result, err
	}
}

func (u *Group) getGroupMembersInfoFromSvr(groupId string, memberList []string) ([]open_im_sdk.groupMemberFullInfo, error) {
	var req open_im_sdk.getGroupMembersInfoReq
	req.GroupID = groupId
	req.OperationID = utils.operationIDGenerator()
	req.MemberList = memberList

	resp, err := utils.post2Api(open_im_sdk.getGroupMembersInfoRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	var sctResp open_im_sdk.getGroupMembersInfoResp
	err = json.Unmarshal(resp, &sctResp)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		utils.sdkLog("errcode, errmsg: ", sctResp.ErrCode, sctResp.ErrMsg)
		return nil, errors.New(sctResp.ErrMsg)
	}
	return sctResp.Data, nil
}

func (u *Group) kickGroupMember(groupId string, memberList []string, reason string, callback common.Base, operationID string) ([]open_im_sdk.idResult, error) {
	var req open_im_sdk.kickGroupMemberApiReq
	req.OperationID = utils.operationIDGenerator()
	memberListInfo, err := u.getGroupMembersInfoFromLocal(groupId, memberList)
	if err != nil {
		utils.sdkLog("getGroupMembersInfoFromLocal, ", err.Error())
		return nil, err
	}
	req.UidListInfo = memberListInfo
	req.Reason = reason
	req.GroupID = groupId

	resp, err := utils.post2Api(open_im_sdk.kickGroupMemberRouter, req, u.token)
	if err != nil {
		utils.sdkLog("post2Api failed, ", open_im_sdk.kickGroupMemberRouter, req, err.Error())
		return nil, err
	}
	utils.sdkLog("url: ", open_im_sdk.kickGroupMemberRouter, "req:", req, "resp: ", string(resp))

	u.syncGroupMemberByGroupId(groupId)
	utils.sdkLog("syncGroupMemberByGroupId: ", groupId)

	var sctResp open_im_sdk.kickGroupMemberApiResp
	err = json.Unmarshal(resp, &sctResp)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error(), resp)
		return nil, err
	}

	if sctResp.ErrCode != 0 {
		utils.sdkLog("resp failed, ", sctResp.ErrCode, sctResp.ErrMsg)
		return nil, errors.New(sctResp.ErrMsg)
	}
	utils.sdkLog("kickGroupMember, ", groupId, memberList, reason, req)
	return sctResp.Data, nil
}

//1
func (u *Group) transferGroupOwner(groupId, userId string, callback common.Base, operationID string) error {
	resp, err := utils.post2Api(open_im_sdk.transferGroupRouter, open_im_sdk.transferGroupReq{GroupID: groupId, Uid: userId, OperationID: utils.operationIDGenerator()}, u.token)
	if err != nil {
		return err
	}
	var ret open_im_sdk.commonResp
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
func (u *Group) inviteUserToGroup(groupId string, reason string, userList []string, callback common.Base, operationID string) ([]open_im_sdk.idResult, error) {
	var req open_im_sdk.inviteUserToGroupReq
	req.GroupID = groupId
	req.OperationID = utils.operationIDGenerator()
	req.Reason = reason
	req.UidList = userList
	resp, err := utils.post2Api(open_im_sdk.inviteUserToGroupRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	u.syncGroupMemberByGroupId(groupId)
	utils.sdkLog("syncGroupMemberByGroupId", groupId)
	var stcResp open_im_sdk.inviteUserToGroupResp
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		utils.sdkLog("unmarshal failed, ", err.Error())
		return nil, err
	}
	if stcResp.ErrCode != 0 {
		utils.sdkLog("errcode, errmsg: ", stcResp.ErrCode, stcResp.ErrMsg)
		return nil, errors.New(stcResp.ErrMsg)
	}

	utils.sdkLog("inviteUserToGroup, autoSendInviteUserToGroupTip", groupId, reason, userList, req, err)
	return stcResp.Data, nil
}

func (u *Group) getLocalGroupApplicationList(groupId string) (*open_im_sdk.groupApplicationResult, error) {
	reply, err := u.getOwnLocalGroupApplicationList(groupId)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (u *Group) delGroupRequestFromGroupRequest(info open_im_sdk.GroupReqListInfo) error {
	return u.delRequestFromGroupRequest(info)
}

//1
func (u *Group) getGroupApplicationList(, callback common.Base, operationID string) (*open_im_sdk.groupApplicationResult, error) {
	resp, err := utils.post2Api(open_im_sdk.getGroupApplicationListRouter, open_im_sdk.getGroupApplicationListReq{OperationID: utils.operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}

	var ret open_im_sdk.getGroupApplicationListResp
	utils.sdkLog("getGroupApplicationListResp", string(resp))
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		utils.sdkLog("unmarshal failed", err.Error())
		return nil, err
	}
	if ret.ErrCode != 0 {
		utils.sdkLog("errcode, errmsg: ", ret.ErrCode, ret.ErrMsg)
		return nil, errors.New(ret.ErrMsg)
	}

	return &ret.Data, nil
}

//1
func (u *Group) acceptGroupApplication(access *open_im_sdk.accessOrRefuseGroupApplicationReq, callback common.Base, operationID string) error {
	resp, err := utils.post2Api(open_im_sdk.acceptGroupApplicationRouter, access, u.token)
	if err != nil {
		return err
	}

	var ret open_im_sdk.commonResp
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
func (u *Group) refuseGroupApplication(access *open_im_sdk.accessOrRefuseGroupApplicationReq, callback common.Base, operationID string) error {
	resp, err := utils.post2Api(open_im_sdk.acceptGroupApplicationRouter, access, u.token)
	if err != nil {
		return err
	}

	var ret open_im_sdk.commonResp
	err = json.Unmarshal(resp, &ret)
	if err != nil {
		return err
	}
	if ret.ErrCode != 0 {
		return errors.New(ret.ErrMsg)
	}
	return nil
}

func (u *Group) getGroupInfoByGroupId(groupId string) (open_im_sdk.groupInfo, error) {
	var gList []string
	gList = append(gList, groupId)
	rList, err := u.getGroupsInfo(gList)
	if err == nil && len(rList) == 1 {
		return rList[0], nil
	} else {
		return open_im_sdk.groupInfo{}, nil
	}

}



func (u *Group) createGroupCallback(node open_im_sdk.updateGroupNode) {
	// member list to json
	jsonMemberList, err := json.Marshal(node.Args.(open_im_sdk.createGroupArgs).initMemberList)
	if err != nil {
		return
	}
	u.listener.OnMemberEnter(node.groupId, string(jsonMemberList))
	u.listener.OnGroupCreated(node.groupId)
}

func (u *Group) joinGroupCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.joinGroupArgs)
	jsonApplyUser, err := json.Marshal(args.applyUser)
	if err != nil {
		return
	}
	u.listener.OnReceiveJoinApplication(node.groupId, string(jsonApplyUser), args.reason)
}

func (u *Group) quitGroupCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.quiteGroupArgs)
	jsonUser, err := json.Marshal(args.quiteUser)
	if err != nil {
		return
	}
	u.listener.OnMemberLeave(node.groupId, string(jsonUser))
}

func (u *Group) setGroupInfoCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.setGroupInfoArgs)
	jsonGroup, err := json.Marshal(args.group)
	if err != nil {
		return
	}
	u.listener.OnGroupInfoChanged(node.groupId, string(jsonGroup))
}

func (u *Group) kickGroupMemberCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.kickGroupAgrs)
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

func (u *open_im_sdk) transferGroupOwnerCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.transferGroupArgs)

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

func (u *Group) inviteUserToGroupCallback(node open_im_sdk.updateGroupNode) {
	args := node.Args.(open_im_sdk.inviteUserToGroupArgs)
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

func (u *Group) GroupApplicationProcessedCallback(node open_im_sdk.updateGroupNode, process int32) {
	args := node.Args.(open_im_sdk.applyGroupProcessedArgs)
	list := make([]open_im_sdk.groupMemberFullInfo, 0)
	for _, v := range args.applyList {
		list = append(list, v.member)
	}
	jsonApplyList, err := json.Marshal(list)
	if err != nil {
		return
	}

	processed := node.Args.(open_im_sdk.applyGroupProcessedArgs) //receiver : all group member
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

func (u *Group) acceptGroupApplicationCallback(node open_im_sdk.updateGroupNode) {
	u.GroupApplicationProcessedCallback(node, 1)
}

func (u *Group) refuseGroupApplicationCallback(node open_im_sdk.updateGroupNode) {
	u.GroupApplicationProcessedCallback(node, -1)
}

func (u *Group) syncSelfGroupRequest() {

}

func (u *Group) syncGroupRequest() {
	groupRequestOnServerResp, err := u.getGroupApplicationList()
	if err != nil {
		utils.sdkLog("groupRequestOnServerResp failed", err.Error())
		return
	}
	groupRequestOnServer := groupRequestOnServerResp.GroupApplicationList
	groupRequestOnServerInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupRequestOnServer {
		groupRequestOnServerInterface = append(groupRequestOnServerInterface, v)
	}

	groupRequestOnLocalResp, err := u.getLocalGroupApplicationList("")
	if err != nil {
		utils.sdkLog("groupRequestOnLocalResp failed", err.Error())
		return
	}
	groupRequestOnLocal := groupRequestOnLocalResp.GroupApplicationList
	groupRequestOnLocalInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupRequestOnLocal {
		groupRequestOnLocalInterface = append(groupRequestOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := utils.checkDiff(groupRequestOnServerInterface, groupRequestOnLocalInterface)

	utils.sdkLog("len ", len(aInBNot), len(bInANot), len(sameA))
	for _, index := range aInBNot {
		err = u.insertIntoRequestToGroupRequest(groupRequestOnServer[index])
		if err != nil {
			utils.sdkLog("insertIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		utils.sdkLog("insertIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

	for _, index := range bInANot {
		err = u.delGroupRequestFromGroupRequest(groupRequestOnLocal[index])
		if err != nil {
			utils.sdkLog("delGroupRequestFromGroupRequest failed", err.Error())
			continue
		}
		utils.sdkLog("delGroupRequestFromGroupRequest ", groupRequestOnLocal[index])
	}
	for _, index := range sameA {
		if err = u.replaceIntoRequestToGroupRequest(groupRequestOnServer[index]); err != nil {
			utils.sdkLog("replaceIntoRequestToGroupRequest failed", err.Error())
			continue
		}
		utils.sdkLog("replaceIntoRequestToGroupRequest ", groupRequestOnServer[index])
	}

}

func (g *groupListener) syncApplyGroupRequest() {

}

func (u *Group) syncJoinedGroupInfo() {
	groupListOnServer, err := u.getJoinedGroupListFromSvr()
	if err != nil {
		utils.sdkLog("groupListOnServer failed", err.Error())
		return
	}
	groupListOnServerInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupListOnServer {
		groupListOnServerInterface = append(groupListOnServerInterface, v)
	}

	groupListOnLocal, err := u.getJoinedGroupListFromLocal()
	if err != nil {
		utils.sdkLog("groupListOnLocal failed", err.Error())
		return
	}
	groupListOnLocalInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupListOnLocal {
		groupListOnLocalInterface = append(groupListOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := utils.checkDiff(groupListOnServerInterface, groupListOnLocalInterface)

	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupInfo(groupListOnServer[index])
		if err != nil {
			utils.sdkLog("insertIntoLocalGroupInfo failed", err.Error(), groupListOnServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupInfo(groupListOnLocal[index].GroupId)
		if err != nil {
			utils.sdkLog("delLocalGroupInfo failed", err.Error())
			continue
		}
	}
	for _, index := range sameA {
		if err = u.replaceLocalGroupInfo(groupListOnServer[index]); err != nil {
			utils.sdkLog("replaceLocalGroupInfo failed", err.Error())
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

func (u *Group) getLocalGroupInfoByGroupId1(groupId string) (*Group.groupInfo, error) {
	return u.getLocalGroupsInfoByGroupID(groupId)
}

func (u *Group) syncGroupMemberByGroupId(groupId string) {
	groupMemberOnServer, err := u.getGroupAllMemberListByGroupIdFromSvr(groupId)
	if err != nil {
		utils.sdkLog("syncGroupMemberByGroupId failed", err.Error())
		return
	}
	utils.sdkLog("getGroupAllMemberListByGroupIdFromSvr, ", groupId, len(groupMemberOnServer))

	groupMemberOnServerInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := u.getLocalGroupMemberListByGroupID(groupId)
	if err != nil {
		utils.sdkLog("getLocalGroupMemberListByGroupID failed", err.Error())
		return
	}
	utils.sdkLog("getLocalGroupMemberListByGroupID, ", groupId, len(groupMemberOnLocal))

	groupMemberOnLocalInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupMemberOnLocal {
		groupMemberOnLocalInterface = append(groupMemberOnLocalInterface, v)
	}
	aInBNot, bInANot, sameA, _ := utils.checkDiff(groupMemberOnServerInterface, groupMemberOnLocalInterface)
	//0 0 2 2 3
	utils.sdkLog("diff len: ", len(aInBNot), len(bInANot), len(sameA), len(groupMemberOnServerInterface), len(groupMemberOnLocalInterface))
	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			utils.sdkLog("insertIntoLocalGroupMember failed", err.Error(), "index", index, groupMemberOnServer[index])
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			utils.sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = u.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			utils.sdkLog("replaceLocalGroupMemberInfo failed", err.Error())
			continue
		}
	}

}

func (u *Group) syncJoinedGroupMember() {
	groupMemberOnServer, err := u.getJoinGroupAllMemberList()
	if err != nil {
		utils.sdkLog("getJoinGroupAllMemberList failed", err.Error())
		return
	}
	groupMemberOnServerInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupMemberOnServer {
		groupMemberOnServerInterface = append(groupMemberOnServerInterface, v)
	}
	groupMemberOnLocal, err := u.getLocalGroupMemberList()
	if err != nil {
		utils.sdkLog("getLocalGroupMemberList failed", err.Error())
		return
	}
	groupMemberOnLocalInterface := make([]open_im_sdk.diff, 0)
	for _, v := range groupMemberOnLocal {
		groupMemberOnLocalInterface = append(groupMemberOnLocalInterface, v)
	}

	aInBNot, bInANot, sameA, _ := utils.checkDiff(groupMemberOnServerInterface, groupMemberOnLocalInterface)

	for _, index := range aInBNot {
		err = u.insertIntoLocalGroupMember(groupMemberOnServer[index])
		if err != nil {
			utils.sdkLog("insertIntoLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range bInANot {
		err = u.delLocalGroupMember(groupMemberOnLocal[index])
		if err != nil {
			utils.sdkLog("delLocalGroupMember failed", err.Error())
			continue
		}
	}

	for _, index := range sameA {
		err = u.replaceLocalGroupMemberInfo(groupMemberOnServer[index])
		if err != nil {
			utils.sdkLog(err.Error())
			continue
		}
	}

}

func (u *Group) getJoinGroupAllMemberList(callback common.Base, operationID string) ([]open_im_sdk.groupMemberFullInfo, error) {
	groupInfoList, err := u.getJoinedGroupListFromLocal()
	if err != nil {
		return nil, err
	}
	joinGroupMemberList := make([]open_im_sdk.groupMemberFullInfo, 0)
	for _, v := range groupInfoList {
		theGroupMemberList, err := u.getGroupAllMemberListByGroupIdFromSvr(v.GroupId)
		if err != nil {
			utils.sdkLog(err.Error())
			continue
		}
		for _, v := range theGroupMemberList {
			joinGroupMemberList = append(joinGroupMemberList, v)
		}
	}
	return joinGroupMemberList, nil
}

func (u *Group) getGroupAllMemberListByGroupIdFromSvr(groupId string) ([]open_im_sdk.groupMemberFullInfo, error) {
	var req open_im_sdk.getGroupAllMemberReq
	req.OperationID = utils.operationIDGenerator()
	req.GroupID = groupId

	resp, err := utils.post2Api(open_im_sdk.getGroupAllMemberListRouter, req, u.token)
	if err != nil {
		return nil, err
	}
	utils.sdkLog("getGroupAllMemberListRouter", open_im_sdk.getGroupAllMemberListRouter, req, string(resp))
	var stcResp open_im_sdk.groupMemberInfoResult
	err = json.Unmarshal(resp, &stcResp)
	if err != nil {
		utils.sdkLog("Unmarshal failed, ", err.Error())
		return nil, err
	}

	if stcResp.ErrCode != 0 {
		utils.sdkLog("errcode errmsg ", stcResp.ErrCode, stcResp.ErrMsg)
		return nil, errors.New(stcResp.ErrMsg)
	}
	return stcResp.Data, nil
}

func (u *Group) getLocalGroupMemberListNew() ([]open_im_sdk.groupMemberFullInfo, error) {
	return u.getLocalGroupMemberList()
}

func (u *Group) getLocalGroupMemberListByGroupIDNew(groupId string) ([]open_im_sdk.groupMemberFullInfo, error) {
	return u.getLocalGroupMemberListByGroupID(groupId)
}
func (u *Group) insertIntoLocalGroupMemberNew(info open_im_sdk.groupMemberFullInfo) error {
	return u.insertIntoLocalGroupMember(info)
}
func (u *Group) delLocalGroupMemberNew(info open_im_sdk.groupMemberFullInfo) error {
	return u.delLocalGroupMember(info)
}
func (u *Group) replaceLocalGroupMemberInfoNew(info open_im_sdk.groupMemberFullInfo) error {
	return u.replaceLocalGroupMemberInfo(info)
}

func (u *Group) insertIntoSelfApplyToGroupRequestNew(groupId, message string) error {
	return u.insertIntoSelfApplyToGroupRequest(groupId, message)
}
