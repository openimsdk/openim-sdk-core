package group

import (
	"errors"
	"math/big"
	comm "open_im_sdk/internal/common"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/db_interface"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	sdk "open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"sort"
	"sync"
	"time"

	"github.com/jinzhu/copier"
)

// //utils.GetCurrentTimestampByMill()
type Group struct {
	listener           open_im_sdk_callback.OnGroupListener
	loginUserID        string
	db                 db_interface.DataBase
	p                  *ws.PostApi
	loginTime          int64
	joinedSuperGroupCh chan common.Cmd2Value
	heartbeatCmdCh     chan common.Cmd2Value

	conversationCh chan common.Cmd2Value
	//	memberSyncMutex sync.RWMutex

	listenerForService open_im_sdk_callback.OnListenerForService
}

func (g *Group) LoginTime() int64 {
	return g.loginTime
}

func (g *Group) SetLoginTime(loginTime int64) {
	g.loginTime = loginTime
}

func (g *Group) SetListenerForService(listener open_im_sdk_callback.OnListenerForService) {
	g.listenerForService = listener
}

func NewGroup(loginUserID string, db db_interface.DataBase, p *ws.PostApi,
	joinedSuperGroupCh chan common.Cmd2Value, heartbeatCmdCh chan common.Cmd2Value,
	conversationCh chan common.Cmd2Value) *Group {
	return &Group{loginUserID: loginUserID, db: db, p: p, joinedSuperGroupCh: joinedSuperGroupCh, heartbeatCmdCh: heartbeatCmdCh, conversationCh: conversationCh}
}

func (g *Group) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
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

func (g *Group) deleteMemberImmediately(groupID string, userID string, operationID string) {
	localMember, err := g.db.GetGroupMemberInfoByGroupIDUserID(groupID, userID)
	if err != nil {
		log.Error(operationID, "GetGroupMemberInfoByGroupIDUserID failed ", err.Error(), groupID, userID)
		return
	}
	err = g.db.DeleteGroupMember(groupID, userID)
	if err != nil {
		log.Error(operationID, "DeleteGroupMember failed ", err.Error(), groupID, userID)
		return
	}
	//err = g.db.SubtractMemberCount(groupID)
	//if err != nil {
	//	log.Error(operationID, "SubtractMemberCount failed ", err.Error(), groupID)
	//}
	//localMember := model_struct.LocalGroupMember{GroupID: groupID, UserID: userID}
	if g.listener != nil {
		g.listener.OnGroupMemberDeleted(utils.StructToJsonString(localMember))
	}

}

func (g *Group) addMemberImmediately(member *api.GroupMemberFullInfo, operationID string) {
	localMember := model_struct.LocalGroupMember{}
	common.GroupMemberCopyToLocal(&localMember, member)
	err := g.db.InsertGroupMember(&localMember)
	if err != nil {
		log.Error(operationID, "InsertGroupMember failed ", err.Error(), *member)
		return
	}
	//err = g.db.AddMemberCount(member.GroupID)
	//if err != nil {
	//	log.Error(operationID, "AddMemberCount failed ", err.Error(), member.GroupID)
	//}
	if g.listener != nil {
		g.listener.OnGroupMemberAdded(utils.StructToJsonString(localMember))
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

func (g *Group) updateMemberImmediately(memberInfo *api.GroupMemberFullInfo, operationID string) {
	localMember := model_struct.LocalGroupMember{}
	common.GroupMemberCopyToLocal(&localMember, memberInfo)
	localMemberGroup, err := g.db.GetGroupMemberInfoByGroupIDUserID(memberInfo.GroupID, memberInfo.UserID)
	if err != nil {
		log.NewError(operationID, "GetGroupMemberInfoByGroupIDUserID failed ", err.Error(), "groupID", memberInfo.GroupID, "userID", memberInfo.UserID)
		return
	}
	err = g.db.UpdateGroupMember(&localMember)
	if err != nil {
		log.Error(operationID, "UpdateGroupMember failed ", err.Error(), localMember)
		return
	}
	if g.listener != nil {
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(localMember))
		log.Info(operationID, "OnGroupMemberInfoChanged", utils.StructToJsonString(localMember))
		if localMemberGroup.Nickname == localMember.Nickname && localMemberGroup.FaceURL == localMember.FaceURL {
			log.NewInfo(operationID, "OnGroupMemberInfoChanged nickname faceURL unchanged", localMember.GroupID, localMember.UserID, localMember.Nickname, localMember.FaceURL)
			return
		}
		_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: localMember.UserID, FaceURL: localMember.FaceURL,
			Nickname: localMember.Nickname, GroupID: localMember.GroupID}}, g.conversationCh)
	}

}

func (g *Group) updateLocalMemberImmediately(groupID, userID string, args map[string]interface{}, operationID string) {
	err := g.db.UpdateGroupMemberField(groupID, userID, args)
	if err != nil {
		log.Error(operationID, "UpdateGroupMemberField failed ", err.Error(), groupID, userID, args)
		return
	}
	member, err := g.db.GetGroupMemberInfoByGroupIDUserID(groupID, userID)
	if err != nil {
		log.Error(operationID, "GetGroupMemberInfoByGroupIDUserID failed ", err.Error(), groupID, userID)
		return
	}
	if g.listener != nil {
		g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(member))
	}
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

func (g *Group) createGroup(callback open_im_sdk_callback.Base, group sdk.CreateGroupBaseInfoParam, memberList sdk.CreateGroupMemberRoleParam, operationID string) *sdk.CreateGroupCallback {
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
	g.p.PostFatalCallbackPenetrate(callback, constant.CreateGroupRouter, apiReq, &realData.GroupInfo, apiReq.OperationID)
	m := utils.JsonDataOne(&realData.GroupInfo)
	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(realData.GroupInfo.GroupID, operationID, false)
	return (*sdk.CreateGroupCallback)(&m)
}

func (g *Group) joinGroup(groupID, reqMsg string, joinSource int32, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.JoinGroupReq{}
	apiReq.OperationID = operationID
	apiReq.ReqMessage = reqMsg
	apiReq.GroupID = groupID
	apiReq.JoinSource = joinSource
	g.p.PostFatalCallback(callback, constant.JoinGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncSelfGroupApplication(operationID)
}

func (g *Group) GetGroupOwnerIDAndAdminIDList(groupID, operationID string) (ownerID string, adminIDList []string, err error) {
	localGroup, err := g.db.GetGroupInfoByGroupID(groupID)
	if err != nil {
		return "", nil, err
	}
	adminIDList, err = g.db.GetGroupAdminID(groupID)
	if err != nil {
		return "", nil, err
	}
	return localGroup.OwnerUserID, adminIDList, nil
}

func (g *Group) quitGroup(groupID string, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.QuitGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.QuitGroupRouter, apiReq, nil, apiReq.OperationID)
	g.db.DeleteGroupAllMembers(groupID)
	g.SyncJoinedGroupList(operationID)
}

func (g *Group) dismissGroup(groupID string, callback open_im_sdk_callback.Base, operationID string) {
	apiReq := api.DismissGroupReq{}
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	g.p.PostFatalCallback(callback, constant.DismissGroupRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
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
	g.SyncJoinedGroupList(operationID)
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
	g.updateLocalMemberImmediately(groupID, userID, map[string]interface{}{"mute_end_time": uint32(int64(time.Now().Second()) + int64(mutedSeconds))}, operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, true)
}

// func (g *Group) updateGroup(groupID, userID string, updateField map[string]interface{}, operationID string) {
//
//		g.updateLocalMemberImmediately(member, operationID)
//	}
func (g *Group) setGroupMemberRoleLevel(callback open_im_sdk_callback.Base, groupID, userID string, roleLevel int, operationID string) {
	apiReq := api.SetGroupMemberRoleLevelReq{
		SetGroupMemberBaseInfoReq: api.SetGroupMemberBaseInfoReq{
			OperationID: operationID,
			UserID:      userID,
			GroupID:     groupID,
		},
		RoleLevel: roleLevel,
	}
	g.p.PostFatalCallback(callback, constant.SetGroupMemberInfoRouter, apiReq, nil, apiReq.OperationID)
	g.updateLocalMemberImmediately(groupID, userID, map[string]interface{}{"role_level": int32(roleLevel)}, operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, true)
}

func (g *Group) setGroupMemberInfo(callback open_im_sdk_callback.Base, param sdk.SetGroupMemberInfoParam, operationID string) {
	apiReq := api.SetGroupMemberInfoReq{OperationID: operationID, Ex: param.Ex, UserID: param.UserID, GroupID: param.GroupID}
	g.p.PostFatalCallback(callback, constant.SetGroupMemberInfoRouter, apiReq, nil, apiReq.OperationID)
	if param.Ex != nil {
		g.updateLocalMemberImmediately(param.GroupID, param.UserID, map[string]interface{}{"ex": *param.Ex}, operationID)
	} else {
		g.updateLocalMemberImmediately(param.GroupID, param.UserID, map[string]interface{}{"ex": ""}, operationID)
	}
	g.syncGroupMemberByGroupID(param.GroupID, operationID, true)
}

func (g *Group) getJoinedGroupList(callback open_im_sdk_callback.Base, operationID string) sdk.GetJoinedGroupListCallback {
	groupList, err := g.db.GetJoinedGroupListDB()
	log.Info(operationID, utils.GetSelfFuncName(), " args ", groupList)
	common.CheckDBErrCallback(callback, err, operationID)
	superGroupList, _ := g.db.GetJoinedSuperGroupList()
	groupList = append(groupList, superGroupList...)
	return groupList
}

func (g *Group) GetGroupInfoFromLocal2Svr(groupID string) (*model_struct.LocalGroup, error) {
	localGroup, err1 := g.db.GetGroupInfoByGroupID(groupID)
	if err1 == nil {
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
		return nil, utils.Wrap(err, "get groupInfo from server err")
	} else {
		return nil, utils.Wrap(errors.New("server not this group"), "")
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
	groupList, err := g.db.GetJoinedGroupListDB()
	common.CheckDBErrCallback(callback, err, operationID)
	superGroupList, err := g.db.GetJoinedSuperGroupList()
	common.CheckDBErrCallback(callback, err, operationID)
	if len(superGroupList) > 0 {
		groupList = append(groupList, superGroupList...)
	}
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

func (g *Group) getGroupAbstractInfoFromSvr(groupID string, operationID string) (*api.GetGroupAbstractInfoResp, error) {
	apiReq := api.GetGroupAbstractInfoReq{}
	apiReq.GroupID = groupID
	apiReq.OperationID = operationID
	var groupAbstractInfoResp api.GetGroupAbstractInfoResp
	err := g.p.Post2UnmarshalRespReturn(constant.GetGroupAbstractInfoRouter, apiReq, &groupAbstractInfoResp)
	if err != nil {
		return nil, utils.Wrap(err, apiReq.OperationID+" "+groupID)
	}
	return &groupAbstractInfoResp, nil
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

func (g *Group) modifyGroupInfo(callback open_im_sdk_callback.Base, apiReq api.SetGroupInfoReq, operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", apiReq)
	g.p.PostFatalCallback(callback, constant.SetGroupInfoRouter, apiReq, nil, apiReq.OperationID)
	g.SyncJoinedGroupList(operationID)
}

// todo
func (g *Group) getGroupMemberList(callback open_im_sdk_callback.Base, groupID string, filter, offset, count int32, operationID string) sdk.GetGroupMemberListCallback {
	groupInfoList, err := g.db.GetGroupMemberListSplit(groupID, filter, int(offset), int(count))
	common.CheckDBErrCallback(callback, err, operationID)
	return groupInfoList
}

func (g *Group) getGroupMemberOwnerAndAdmin(callback open_im_sdk_callback.Base, groupID string, operationID string) sdk.GetGroupMemberListCallback {
	groupInfoList, err := g.db.GetGroupMemberOwnerAndAdmin(groupID)
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

// todo
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
	//	g.SyncJoinedGroupList(operationID)
	for _, v := range memberList {
		g.deleteMemberImmediately(groupID, v, operationID)
	}
	g.syncGroupMemberByGroupID(groupID, operationID, true)
	return realData.UserIDResultList
}

func (g *Group) transferGroupOwner(callback open_im_sdk_callback.Base, groupID, newOwnerUserID string, operationID string) {
	apiReq := api.TransferGroupOwnerReq{}
	apiReq.GroupID = groupID
	apiReq.NewOwnerUserID = newOwnerUserID
	apiReq.OperationID = operationID
	apiReq.OldOwnerUserID = g.loginUserID
	g.p.PostFatalCallback(callback, constant.TransferGroupRouter, apiReq, nil, apiReq.OperationID)

	g.updateLocalMemberImmediately(groupID, g.loginUserID, map[string]interface{}{"role_level": constant.GroupOrdinaryUsers}, operationID)
	g.updateLocalMemberImmediately(groupID, newOwnerUserID, map[string]interface{}{"role_level": constant.GroupOwner}, operationID)
	g.SyncJoinedGroupList(operationID)
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
	//	g.SyncJoinedGroupList(operationID)
	g.syncGroupMemberByGroupID(groupID, operationID, false)
	return realData
}

// //1
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

func (g *Group) GetJoinedDiffusionGroupIDListFromSvr(operationID string) ([]string, error) {
	result, err := g.getJoinedGroupListFromSvr(operationID)
	if err != nil {
		return nil, utils.Wrap(err, "working group get err")
	}
	var groupIDList []string
	for _, v := range result {
		if v.GroupType == constant.WorkingGroup {
			groupIDList = append(groupIDList, v.GroupID)
		}
	}
	return groupIDList, nil
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
		log.Info(operationID, "OnGroupInfoChanged, update group info ", groupInfo)
		g.db.UpdateGroup(groupInfo)
		if g.listener != nil {
			g.listener.OnGroupInfoChanged(utils.StructToJsonString(groupInfo))
		}
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

	log.NewInfo(operationID, "svrList onServer onLocal ", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupRequestDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB)
	for _, index := range aInBNot {
		err := g.db.InsertGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupRequest failed ", err.Error(), *onServer[index])
			continue
		}
		callbackData := *onServer[index]
		if g.listener != nil {
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAdded ", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := g.db.UpdateGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == constant.GroupResponseRefuse {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))
			}

		} else if onServer[index].HandleResult == constant.GroupResponseAgree {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
			if g.listenerForService != nil {
				g.listenerForService.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
		} else {
			callbackData := *onServer[index]
			if g.listener != nil {
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := *onLocal[index]
		if g.listener != nil {
			g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
		}
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
		if g.listener != nil {
			g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
		}
		if g.listenerForService != nil {
			g.listenerForService.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupApplicationAdded ", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		err := g.db.UpdateAdminGroupRequest(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupRequest failed ", err.Error())
			continue
		}
		if onServer[index].HandleResult == constant.GroupResponseRefuse {
			callbackData := sdk.GroupApplicationRejectCallback(*onServer[index])
			if g.listener != nil {
				g.listener.OnGroupApplicationRejected(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationRejected", utils.StructToJsonString(callbackData))
			}

		} else if onServer[index].HandleResult == constant.GroupResponseAgree {
			callbackData := sdk.GroupApplicationAcceptCallback(*onServer[index])
			if g.listener != nil {
				g.listener.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
			if g.listenerForService != nil {
				g.listenerForService.OnGroupApplicationAccepted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAccepted", utils.StructToJsonString(callbackData))
			}
		} else {
			callbackData := sdk.GroupApplicationAcceptCallback(*onServer[index])
			if g.listener != nil {
				g.listener.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAdded", utils.StructToJsonString(callbackData))
			}
			if g.listenerForService != nil {
				g.listenerForService.OnGroupApplicationAdded(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupApplicationAdded ", utils.StructToJsonString(callbackData))
			}
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteAdminGroupRequest(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupRequest failed ", err.Error())
			continue
		}
		callbackData := sdk.GroupApplicationDeletedCallback(*onLocal[index])
		if g.listener != nil {
			g.listener.OnGroupApplicationDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnReceiveJoinGroupApplicationDeleted", utils.StructToJsonString(callbackData))
		}
	}
}

//	func transferGroupInfo(input []*api.GroupInfo) []*api.GroupInfo{
//		var result []*api.GroupInfo
//		for _, v := range input {
//			t := &api.GroupInfo{}
//			copier.Copy(t, &v)
//			if v.NeedVerification != nil {
//				t.NeedVerification = v.NeedVerification.Value
//			}
//			result = append(result, t)
//		}
//		return result
//	}
func (g *Group) SyncJoinedGroupList(operationID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
	svrList, err := g.getJoinedGroupListFromSvr(operationID)
	log.Info(operationID, "getJoinedGroupListFromSvr", svrList, g.loginUserID)
	if err != nil {
		log.NewError(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
		return
	}
	onServer := common.TransferToLocalGroupInfo(svrList)
	onLocal, err := g.db.GetJoinedGroupListDB()
	if err != nil {
		log.NewError(operationID, "GetJoinedGroupList failed ", err.Error())
		return
	}

	log.NewInfo(operationID, " onLocal ", onLocal, g.loginUserID)
	aInBNot, bInANot, sameA, sameB := common.CheckGroupInfoDiff(onServer, onLocal)
	log.Info(operationID, "diff ", aInBNot, bInANot, sameA, sameB, g.loginUserID)
	var isReadDiffusion bool
	for _, index := range aInBNot {
		err := g.db.InsertGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroup failed ", err.Error(), *onServer[index])
			continue
		}

		callbackData := sdk.JoinedGroupAddedCallback(*onServer[index])
		if (*onServer[index]).GroupType == int32(constant.WorkingGroup) {
			isReadDiffusion = true
		}
		if g.listener != nil {
			g.listener.OnJoinedGroupAdded(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnJoinedGroupAdded", utils.StructToJsonString(callbackData))
		}
	}
	for _, index := range sameA {
		callbackData := sdk.GroupInfoChangedCallback(*onServer[index])
		localGroup, err := g.db.GetGroupInfoByGroupID(callbackData.GroupID)
		if err != nil {
			log.NewError(operationID, "GetGroupInfoByGroupID failed ", err.Error(), callbackData.GroupID)
			continue
		}
		err = g.db.UpdateGroup(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroup failed ", err.Error(), onServer[index])
			continue
		}
		if g.listener != nil {
			g.listener.OnGroupInfoChanged(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupInfoChanged", utils.StructToJsonString(callbackData))

			if localGroup.GroupName == callbackData.GroupName && localGroup.FaceURL == callbackData.FaceURL {
				log.NewInfo(operationID, "OnGroupInfoChanged nickname faceURL unchanged", callbackData.GroupID, callbackData.GroupName, callbackData.FaceURL)
				continue
			}
			common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.SourceIDAndSessionType{SourceID: callbackData.GroupID, SessionType: constant.GroupChatType}}, g.conversationCh)
		}
	}

	for _, index := range bInANot {
		log.Info(operationID, "DeleteGroup: ", onLocal[index].GroupID, g.loginUserID)
		err := g.db.DeleteGroup(onLocal[index].GroupID)
		if err != nil {
			log.NewError(operationID, "DeleteGroup failed ", err.Error(), onLocal[index].GroupID)
			continue
		}
		if (*onLocal[index]).GroupType == int32(constant.WorkingGroup) {
			isReadDiffusion = true
		}
		g.db.DeleteGroupAllMembers(onLocal[index].GroupID)
		callbackData := sdk.JoinedGroupDeletedCallback(*onLocal[index])
		if g.listener != nil {
			g.listener.OnJoinedGroupDeleted(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnJoinedGroupDeleted", utils.StructToJsonString(callbackData))
		}
	}
	if isReadDiffusion {
		cmd := sdk_struct.CmdJoinedSuperGroup{OperationID: operationID}
		err := common.TriggerCmdJoinedSuperGroup(cmd, g.joinedSuperGroupCh)
		if err != nil {
			log.Error(operationID, "TriggerCmdJoinedSuperGroup failed ", err.Error())
		}
		err = common.TriggerCmdWakeUp(g.heartbeatCmdCh)
		if err != nil {
			log.Error(operationID, "TriggerCmdWakeUp failed ", err.Error())
		}
	}
}

func (g *Group) calculateGroupMemberHash(groupID string, operationID string) (uint64, error) {
	userIDList, err := g.db.GetGroupMemberUIDListByGroupID(groupID)
	if err != nil {
		return 0, utils.Wrap(err, "GetGroupMemberUIDListByGroupID")
	}
	log.NewInfo(operationID, "calculateGroupMemberHash userIDList len: ", len(userIDList), " groupID: ", groupID)
	if len(userIDList) == 0 {
		return 0, nil
	}
	sort.Strings(userIDList)
	all := ""
	for _, v := range userIDList {
		all += v
	}
	bi := big.NewInt(0)
	bi.SetString(utils.Md5(all)[0:8], 16)
	//	log.Info(operationID, "md5 all: ", all, "bi64: ", bi.Uint64())
	return bi.Uint64(), nil
}

func (g *Group) isContinueSyncGroupMember(groupID string, operationID string) bool {
	groupAbstractInfo, err := g.getGroupAbstractInfoFromSvr(groupID, operationID)
	if err != nil {
		log.Error(operationID, "getGroupAbstractInfoFromSvr failed ", groupID, err.Error())
		return true
	}
	log.Debug(operationID, "getGroupAbstractInfoFromSvr success", groupID, "groupAbstractInfo ", *groupAbstractInfo)
	if groupAbstractInfo.GroupMemberNumber < constant.UseHashGroupMemberNum {
		log.Info(operationID, "groupAbstractInfo  group member number: ", groupAbstractInfo.GroupMemberNumber)
		return true
	}

	localGroupHash, err := g.calculateGroupMemberHash(groupID, operationID)
	if err != nil {
		log.Error(operationID, "calculateGroupMemberHash failed ", err.Error(), groupID)
		return true
	}
	log.Info(operationID, "calculateGroupMemberHash ", localGroupHash, groupID)
	if localGroupHash == groupAbstractInfo.GroupMemberListHash {
		log.Info(operationID, " localGroupHash == groupAbstractInfo.GroupMemberListHash ", localGroupHash)
		return false
	} else {
		log.Info(operationID, "localGroupHash != groupAbstractInfo.GroupMemberListHash", localGroupHash, groupAbstractInfo.GroupMemberListHash)
		return true
	}
}

func (g *Group) syncGroupMemberByGroupID(groupID string, operationID string, onGroupMemberNotification bool) {
	//	g.memberSyncMutex.Lock()
	//	defer g.memberSyncMutex.Unlock()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", groupID)
	conSync := g.isContinueSyncGroupMember(groupID, operationID)
	if conSync == false {
		log.Info(operationID, "isContinueSyncGroupMember false, don't sync")
		return
	}
	svrList, err := g.getGroupAllMemberSplitByGroupIDFromSvr(groupID, operationID)
	if err != nil {
		log.NewError(operationID, "getGroupAllMemberSplitByGroupIDFromSvr failed ", err.Error(), groupID)
		return
	}
	log.Info(operationID, "getGroupAllMemberByGroupIDFromSvr ", len(svrList), groupID)
	onServer := common.TransferToLocalGroupMember(svrList)
	onLocal, err := g.db.GetGroupMemberListByGroupID(groupID)
	if err != nil {
		log.NewError(operationID, "GetGroupMemberListByGroupID failed ", err.Error(), groupID)
		return
	}
	//log.NewInfo(operationID, "svrList onServer onLocal", svrList, onServer, onLocal)
	aInBNot, bInANot, sameA, _ := common.CheckGroupMemberDiff(onServer, onLocal)
	log.Info(operationID, "getGroupAllMemberByGroupIDFromSvr  diff ", aInBNot, bInANot, sameA, len(onLocal), len(onServer))
	var insertGroupMemberList []*model_struct.LocalGroupMember
	for _, index := range aInBNot {
		if onGroupMemberNotification == false {
			insertGroupMemberList = append(insertGroupMemberList, onServer[index])
			continue
		}
		err := g.db.InsertGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "InsertGroupMember failed ", err.Error(), *onServer[index])
			continue
		}
		if onGroupMemberNotification == true {
			callbackData := sdk.GroupMemberAddedCallback(*onServer[index])
			if g.listener != nil {
				g.listener.OnGroupMemberAdded(utils.StructToJsonString(callbackData))
				log.Debug(operationID, "OnGroupMemberAdded", utils.StructToJsonString(callbackData))
			}
		}
	}
	if len(insertGroupMemberList) > 0 {
		split := 1000
		idx := 0
		remain := len(insertGroupMemberList) % split
		log.Info(operationID, "BatchInsertGroupMember all: ", len(insertGroupMemberList))
		for idx = 0; idx < len(insertGroupMemberList)/split; idx++ {
			sub := insertGroupMemberList[idx*split : (idx+1)*split]
			err = g.db.BatchInsertGroupMember(sub)
			log.Info(operationID, "BatchInsertGroupMember len: ", len(sub))
			if err != nil {
				log.Error(operationID, "BatchInsertGroupMember failed ", err.Error(), len(sub))
				for again := 0; again < len(sub); again++ {
					if err = g.db.InsertGroupMember(sub[again]); err != nil {
						log.Error(operationID, "InsertGroupMember failed ", err.Error(), sub[again])
					}
				}
			}
		}
		if remain > 0 {
			sub := insertGroupMemberList[idx*split:]
			log.Info(operationID, "BatchInsertGroupMember len: ", len(sub), groupID)
			err = g.db.BatchInsertGroupMember(sub)
			if err != nil {
				log.Error(operationID, "BatchInsertGroupMember failed ", err.Error(), len(sub))
				for again := 0; again < len(sub); again++ {
					if err = g.db.InsertGroupMember(sub[again]); err != nil {
						log.Error(operationID, "InsertGroupMember failed ", err.Error(), sub[again])
					}
				}
			}
		}
	}

	for _, index := range sameA {
		callbackData := sdk.GroupMemberInfoChangedCallback(*onServer[index])
		localMemberGroup, err := g.db.GetGroupMemberInfoByGroupIDUserID(callbackData.GroupID, callbackData.UserID)
		if err != nil {
			log.NewError(operationID, "GetGroupMemberInfoByGroupIDUserID failed ", err.Error(), "groupID", callbackData.GroupID, "userID", callbackData.UserID)
			continue
		}
		err = g.db.UpdateGroupMember(onServer[index])
		if err != nil {
			log.NewError(operationID, "UpdateGroupMember failed ", err.Error(), *onServer[index])
			continue
		}
		callbackData = sdk.GroupMemberInfoChangedCallback(*onServer[index])
		if g.listener != nil {
			g.listener.OnGroupMemberInfoChanged(utils.StructToJsonString(callbackData))
			log.Info(operationID, "OnGroupMemberInfoChanged", utils.StructToJsonString(callbackData))
			if localMemberGroup.Nickname == callbackData.Nickname && localMemberGroup.FaceURL == callbackData.FaceURL {
				log.NewInfo(operationID, "OnGroupMemberInfoChanged nickname faceURL unchanged", callbackData.GroupID, callbackData.UserID, callbackData.Nickname, callbackData.FaceURL)
				continue
			}
			_ = common.TriggerCmdUpdateMessage(common.UpdateMessageNode{Action: constant.UpdateMsgFaceUrlAndNickName, Args: common.UpdateMessageInfo{UserID: callbackData.UserID, FaceURL: callbackData.FaceURL,
				Nickname: callbackData.Nickname, GroupID: callbackData.GroupID}}, g.conversationCh)
		}
	}
	for _, index := range bInANot {
		err := g.db.DeleteGroupMember(onLocal[index].GroupID, onLocal[index].UserID)
		if err != nil {
			log.NewError(operationID, "DeleteGroupMember failed ", err.Error(), onLocal[index].GroupID, onLocal[index].UserID)
			continue
		}
		if onGroupMemberNotification == true {
			callbackData := sdk.GroupMemberDeletedCallback(*onLocal[index])
			if g.listener != nil {
				g.listener.OnGroupMemberDeleted(utils.StructToJsonString(callbackData))
				log.Info(operationID, "OnGroupMemberDeleted", utils.StructToJsonString(callbackData))
			}
		}
	}
}

//func (g *Group) SyncJoinedGroupMember(operationID string) {
//	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ")
//	groupListOnServer, err := g.getJoinedGroupListFromSvr(operationID)
//	if err != nil {
//		log.Error(operationID, "getJoinedGroupListFromSvr failed ", err.Error())
//		return
//	}
//	var wg sync.WaitGroup
//	if len(groupListOnServer) == 0 {
//		return
//	}
//	wg.Add(len(groupListOnServer))
//	log.Info(operationID, "syncGroupMemberByGroupID begin", len(groupListOnServer))
//	for _, v := range groupListOnServer {
//		go func(groupID, operationID string) {
//			g.syncGroupMemberByGroupID(groupID, operationID, true)
//			wg.Done()
//		}(v.GroupID, operationID)
//	}
//
//	wg.Wait()
//	log.Info(operationID, "syncGroupMemberByGroupID end")
//}

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

func (g *Group) getGroupAllMemberSplitByGroupIDFromSvr(groupID string, operationID string) ([]*api.GroupMemberFullInfo, error) {
	var apiReq api.GetGroupAllMemberReq
	apiReq.OperationID = operationID
	apiReq.GroupID = groupID
	var result []*api.GroupMemberFullInfo
	var page int32
	for {
		apiReq.Offset = page * constant.SplitGetGroupMemberNum
		apiReq.Count = constant.SplitGetGroupMemberNum
		var realData []*api.GroupMemberFullInfo
		err := g.p.PostReturn(constant.GetGroupAllMemberListRouter, apiReq, &realData)
		if err != nil {
			log.Error(operationID, "GetGroupAllMemberListRouter failed ", constant.GetGroupAllMemberListRouter, apiReq)
			return result, utils.Wrap(err, apiReq.OperationID)
		}
		log.Info(operationID, "GetGroupAllMemberListRouter result len: ", len(realData), groupID)
		result = append(result, realData...)
		if apiReq.Count > int32(len(realData)) {
			break
		}
		page++
	}
	return result, nil
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

func (g *Group) searchGroupMembers(callback open_im_sdk_callback.Base, searchParam sdk.SearchGroupMembersParam, operationID string) sdk.SearchGroupMembersCallback {
	if len(searchParam.KeywordList) == 0 {
		log.Error(operationID, "len keywordList == 0")
		common.CheckArgsErrCallback(callback, errors.New("no keyword"), operationID)
	}
	members, err := g.db.SearchGroupMembersDB(searchParam.KeywordList[0], searchParam.GroupID, searchParam.IsSearchMemberNickname, searchParam.IsSearchUserID, searchParam.Offset, searchParam.Count)
	common.CheckDBErrCallback(callback, err, operationID)
	return members
}
