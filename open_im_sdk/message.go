/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/1 19:16).
 */
package open_im_sdk

import (
	imgtype "github.com/shamsher31/goimgtype"
	"image"
	"open_im_sdk/open_im_sdk/conversation_msg"
	"open_im_sdk/open_im_sdk/utils"
	"os"
)

const TimeOffset = 5

func doNewMsgConversation() {

}

func (u *UserRelated) createTextSystemMessage(n NotificationContent, textType int32) *MsgStruct {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, textType)
	s.Content = utils.structToJsonString(n)
	s.AtElem.AtUserList = []string{}
	return &s
}

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

func (u *UserRelated) updateMessageFailedStatus(s *MsgStruct, c *conversation_msg.ConversationStruct, onlineUserOnly bool) {
	if !onlineUserOnly {
		_ = u.updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
	}
	s.SendTime = s.CreateTime
	s.Status = MsgStatusSendFailed
	c.LatestMsg = utils.structToJsonString(s)
}
func (u *UserRelated) initBasicInfo(message *MsgStruct, msgFrom, contentType int32) {
	message.CreateTime = utils.getCurrentTimestampByNano()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = MsgStatusSending
	message.SendID = u.loginUserID
	userInfo, _ := u.getLoginUserInfoFromLocal()
	message.SenderFaceURL = userInfo.Icon
	message.SenderNickname = userInfo.Name
	ClientMsgID := utils.getMsgID(message.SendID)
	message.ClientMsgID = ClientMsgID
	message.MsgFrom = msgFrom
	message.ContentType = contentType
	message.SenderPlatformID = SvrConf.Platform

}
func (u *UserRelated) sendMessageFailedHandle(s *MsgStruct, c *conversation_msg.ConversationStruct, conversationID string) {
	_ = u.updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
	s.SendTime = s.CreateTime
	s.Status = MsgStatusSendFailed
	c.LatestMsg = utils.structToJsonString(s)
	_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
		*c})
	u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
}

type MsgFormats []*MsgStruct

// Implement the sort.Interface interface to get the number of elements method
func (s MsgFormats) Len() int {
	return len(s)
}

//Implement the sort.Interface interface comparison element method
func (s MsgFormats) Less(i, j int) bool {
	return s[i].SendTime < s[j].SendTime
}

//Implement the sort.Interface interface exchange element method
func (s MsgFormats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getImageInfo(filePath string) (*imageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		utils.sdkLog(err.Error())
		return nil, err
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		utils.sdkLog(err.Error())
		return nil, err
	}

	datatype, err := imgtype.Get(filePath)
	if err != nil {
		utils.sdkLog(err.Error())
		return nil, err
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		utils.sdkLog(err.Error())
		return nil, err
	}

	b := img.Bounds()

	return &imageInfo{int32(b.Max.X), int32(b.Max.Y), datatype, fi.Size()}, nil

}
