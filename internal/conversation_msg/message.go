/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/1 19:16).
 */
package conversation_msg

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

const TimeOffset = 5

func doNewMsgConversation() {

}

//func (u *open_im_sdk.UserRelated) createTextSystemMessage(n open_im_sdk.NotificationContent, textType int32) *open_im_sdk.MsgStruct {
//	s := utils.MsgStruct{}
//	u.initBasicInfo(&s, constant.UserMsgType, textType)
//	s.Content = utils.structToJsonString(n)
//	s.AtElem.AtUserList = []string{}
//	return &s
//}

/*
func autoSendTransferGroupOwnerTip(groupId string, oldOwner, newOwner string) error{
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, TransferGroupOwnerTip)
	var tReq TransferGroupReq
	jsonReq, err := json.Marshal(tReq)
	if err != nil {
		sdkLog("marshal failed ", err.Error())
		return err
	}
	s.Content = string(jsonReq)

	return autoSendMsg(s, "", groupId, false)
}*/

//func (u *UserRelated) autoSendKickGroupMemberTip(req *kickGroupMemberApiReq) error {
//	s := MsgStruct{}
//	u.initBasicInfo(&s, UserMsgType, KickGroupMemberTip)
//	jsonReq, err := json.Marshal(req)
//	if err != nil {
//		sdkLog("marshal failed ", err.Error())
//		return err
//	}
//
//	var nicknameList string
//	for _, v := range req.UidListInfo {
//		nicknameList = nicknameList + v.NickName + " "
//	}
//
//	var notification NotificationContent
//	notification.IsDisplay = 1
//	notification.Detail = string(jsonReq)
//	notification.DefaultTips = nicknameList + " kicked out of group chat by administrator"
//	jsonContentReq, err := json.Marshal(notification)
//	if err != nil {
//		sdkLog("marshal failed, ", err.Error())
//		return err
//	}
//	s.Content = string(jsonContentReq)
//	s.AtElem.AtUserList = []string{}
//	u.autoSendMsg(&s, "", req.GroupID, false, true, false)
//	sdkLog("sendmsg, group ", s, req.GroupID)
//
//	/*
//		for _, v := range req.UidList {
//			notification.DefaultTips = "You are kicked out of group chat by administrator"
//			jsonContentReq, err := json.Marshal(notification)
//			if err != nil {
//				sdkLog("marshal failed, ", err.Error())
//				return err
//			}
//			s.Content = string(jsonContentReq)
//			autoSendMsg(&s, v, "", false, false)
//			sdkLog("sendmsg, single ", s, v)
//		}*/
//
//	return nil
//}

//func (u *UserRelated) autoSendInviteUserToGroupTip(req inviteUserToGroupReq) error {
//	s := MsgStruct{}
//	u.initBasicInfo(&s, UserMsgType, InviteUserToGroupTip)
//
//	jsonReq, err := json.Marshal(req)
//	if err != nil {
//		sdkLog("marshal failed ", err.Error())
//		return err
//	}
//
//	var nicknameList string
//	for _, v := range req.UidList {
//		member, err := u.getLocalGroupMemberInfoByGroupIdUserId(req.GroupID, v)
//		if err != nil || member.GroupId == "" {
//			sdkLog("getLocalGroupMemberInfoByGroupIdUserId failed ", err, member.GroupId)
//			continue
//		}
//		nicknameList = nicknameList + member.NickName + " "
//	}
//
//	op, err := u.getLocalGroupMemberInfoByGroupIdUserId(req.GroupID, u.LoginUid)
//	if err != nil {
//		sdkLog("nul member, ", req.GroupID, u.LoginUid)
//		return err
//	}
//
//	var notification NotificationContent
//	notification.IsDisplay = 1
//	notification.Detail = string(jsonReq)
//	notification.DefaultTips = nicknameList + "  invited into the group chat by " + op.NickName
//	jsonContentReq, err := json.Marshal(notification)
//	if err != nil {
//		sdkLog("marshal failed, ", err.Error())
//		return err
//	}
//	s.Content = string(jsonContentReq)
//
//	u.autoSendMsg(&s, "", req.GroupID, false, true, false)
//	sdkLog("send msg, groupid: ", req.GroupID)
//	return nil
//}

//func (u *open_im_sdk.UserRelated) updateMessageFailedStatus(s *open_im_sdk.MsgStruct, c *ConversationStruct, onlineUserOnly bool) {
//	if !onlineUserOnly {
//		_ = u.updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, constant.MsgStatusSendFailed)
//	}
//	s.SendTime = s.CreateTime
//	s.Status = constant.MsgStatusSendFailed
//	c.LatestMsg = utils.structToJsonString(s)
//}
func (c *Conversation) initBasicInfo(message *sdk_struct.MsgStruct, msgFrom, contentType int32) {
	message.CreateTime = utils.GetCurrentTimestampByNano()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = constant.MsgStatusSending
	message.SendID = c.loginUserID
	userInfo, _ := c.db.GetLoginUser()
	message.SenderFaceURL = userInfo.FaceURL
	message.SenderNickname = userInfo.Nickname
	ClientMsgID := utils.GetMsgID(message.SendID)
	message.ClientMsgID = ClientMsgID
	message.MsgFrom = msgFrom
	message.ContentType = contentType
	message.SenderPlatformID = c.platformID

}

func (c *Conversation) sendMessageFailedHandle(s *sdk_struct.MsgStruct, lc *db.LocalConversation, conversationID string) {
	_ = c.updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, constant.MsgStatusSendFailed)
	s.SendTime = s.CreateTime
	s.Status = constant.MsgStatusSendFailed
	lc.LatestMsg = utils.StructToJsonString(s)
	_ = u.triggerCmdUpdateConversation(open_im_sdk.updateConNode{conversationID, constant.AddConOrUpLatMsg,
		*c})
	u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", constant.NewConChange, []string{conversationID}}})
}

//
//type MsgFormats []*open_im_sdk.MsgStruct
//
//// Implement the sort.Interface interface to get the number of elements method
//func (s MsgFormats) Len() int {
//	return len(s)
//}
//
////Implement the sort.Interface interface comparison element method
//func (s MsgFormats) Less(i, j int) bool {
//	return s[i].SendTime < s[j].SendTime
//}
//
////Implement the sort.Interface interface exchange element method
//func (s MsgFormats) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}

func getIsRead(b bool) int {
	if b {
		return constant.HasRead
	} else {
		return constant.NotRead
	}
}
func getIsFilter(b bool) int {
	if b {
		return constant.IsFilter
	} else {
		return constant.NotFilter
	}
}
func getIsReadB(i int) bool {
	if i == constant.HasRead {
		return true
	} else {
		return false
	}

}

func (c *Conversation) GetConversationIDBySessionType(sourceID string, sessionType int32) string {
	switch sessionType {
	case constant.SingleChatType:
		return "single_" + sourceID
	case constant.GroupChatType:
		return "group_" + sourceID
	}
	return ""
}
