/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/1 19:16).
 */
package open_im_sdk

import (
	"encoding/json"
	"errors"
	imgtype "github.com/shamsher31/goimgtype"
	"image"
	"os"
)

const TimeOffset = 5

func doNewMsgConversation() {

}

func createTextSystemMessage(n NotificationContent, textType int32) *MsgStruct {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, textType)
	s.Content = structToJsonString(n)
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

func autoSendKickGroupMemberTip(req *kickGroupMemberApiReq) error {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, KickGroupMemberTip)
	jsonReq, err := json.Marshal(req)
	if err != nil {
		sdkLog("marshal failed ", err.Error())
		return err
	}

	var nicknameList string
	for _, v := range req.UidListInfo {
		nicknameList = nicknameList + v.NickName + " "
	}

	var notification NotificationContent
	notification.IsDisplay = 1
	notification.Detail = string(jsonReq)
	notification.DefaultTips = nicknameList + " kicked out of group chat by administrator"
	jsonContentReq, err := json.Marshal(notification)
	if err != nil {
		sdkLog("marshal failed, ", err.Error())
		return err
	}
	s.Content = string(jsonContentReq)
	s.AtElem.AtUserList = []string{}
	autoSendMsg(&s, "", req.GroupID, false, true, false)
	sdkLog("sendmsg, group ", s, req.GroupID)

	/*
		for _, v := range req.UidList {
			notification.DefaultTips = "You are kicked out of group chat by administrator"
			jsonContentReq, err := json.Marshal(notification)
			if err != nil {
				sdkLog("marshal failed, ", err.Error())
				return err
			}
			s.Content = string(jsonContentReq)
			autoSendMsg(&s, v, "", false, false)
			sdkLog("sendmsg, single ", s, v)
		}*/

	return nil
}

func autoSendInviteUserToGroupTip(req inviteUserToGroupReq) error {
	s := MsgStruct{}
	initBasicInfo(&s, UserMsgType, InviteUserToGroupTip)

	jsonReq, err := json.Marshal(req)
	if err != nil {
		sdkLog("marshal failed ", err.Error())
		return err
	}

	var nicknameList string
	for _, v := range req.UidList {
		member, err := getLocalGroupMemberInfoByGroupIdUserId(req.GroupID, v)
		if err != nil || member.GroupId == "" {
			sdkLog("getLocalGroupMemberInfoByGroupIdUserId failed ", err, member.GroupId)
			continue
		}
		nicknameList = nicknameList + member.NickName + " "
	}

	op, err := getLocalGroupMemberInfoByGroupIdUserId(req.GroupID, LoginUid)
	if err != nil {
		sdkLog("nul member, ", req.GroupID, LoginUid)
		return err
	}

	var notification NotificationContent
	notification.IsDisplay = 1
	notification.Detail = string(jsonReq)
	notification.DefaultTips = nicknameList + "  invited into the group chat by " + op.NickName
	jsonContentReq, err := json.Marshal(notification)
	if err != nil {
		sdkLog("marshal failed, ", err.Error())
		return err
	}
	s.Content = string(jsonContentReq)

	autoSendMsg(&s, "", req.GroupID, false, true, false)
	sdkLog("send msg, groupid: ", req.GroupID)
	return nil
}

//
func autoSendMsg(s *MsgStruct, receiver, groupID string, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo bool) error {
	sdkLog("autoSendMsg input args:", *s, receiver, groupID, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo)
	var conversationID string
	r := SendMsgRespFromServer{}
	a := paramsUserSendMsg{}
	op := make(map[string]interface{})
	of := make(map[string]interface{})
	if receiver == "" {
		s.SessionType = GroupChatType
		s.RecvID = groupID
	} else if groupID == "" {
		s.SessionType = SingleChatType
		s.RecvID = receiver
	} else {
		sdkLog("args err: ", receiver, groupID)
		return errors.New("args null")
	}
	c := ConversationStruct{
		ConversationType:  int(s.SessionType),
		RecvMsgOpt:        1,
		LatestMsgSendTime: s.CreateTime,
	}
	if receiver == "" && groupID == "" {
		return errors.New("args error")
	} else if receiver == "" {
		s.SessionType = GroupChatType
		s.RecvID = groupID
		s.GroupID = groupID
		conversationID = GetConversationIDBySessionType(groupID, GroupChatType)
		c.GroupID = groupID
		faceUrl, name, err := getGroupNameAndFaceUrlByUid(groupID)
		if err != nil {
			sdkLog("getGroupNameAndFaceUrlByUid err:", err)
			return err
		}
		c.ShowName = name
		c.FaceURL = faceUrl
	} else {
		s.SessionType = SingleChatType
		s.RecvID = receiver
		conversationID = GetConversationIDBySessionType(receiver, SingleChatType)
		c.UserID = receiver
		faceUrl, name, err := getUserNameAndFaceUrlByUid(receiver)
		if err != nil {
			sdkLog("getUserNameAndFaceUrlByUid err:", err)
			return err
		}
		c.FaceURL = faceUrl
		c.ShowName = name
	}
	userInfo, err := getLoginUserInfoFromLocal()
	if err != nil {
		sdkLog("getLoginUserInfoFromLocal err:", err)
		return err
	}
	s.SenderFaceURL = userInfo.Icon
	s.SenderNickName = userInfo.Name
	c.ConversationID = conversationID
	c.LatestMsg = structToJsonString(s)
	if !onlineUserOnly {
		err = insertMessageToLocalOrUpdateContent(s)
		if err != nil {
			sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			return err
		}
	}

	//Protocol conversion
	a.ReqIdentifier = 1003
	a.PlatformID = s.PlatformID
	a.SendID = s.SendID
	a.OperationID = operationIDGenerator()
	a.Data.SessionType = s.SessionType
	a.SenderNickName = s.SenderNickName
	a.SenderFaceURL = s.SenderFaceURL
	a.Data.MsgFrom = s.MsgFrom
	a.Data.ForceList = []string{}
	a.Data.ContentType = s.ContentType
	a.Data.RecvID = s.RecvID
	a.Data.ForceList = s.ForceList
	a.Data.Content = s.Content
	a.Data.ClientMsgID = s.ClientMsgID
	if onlineUserOnly {
		op["history"] = 0
		op["persistent"] = 0
	}
	a.Data.Options = op
	a.Data.OffLineInfo = of
	bMsg, err := post2Api(sendMsgRouter, a, token)
	if err != nil {
		sdkLog("sendMsgRouter access err:", err.Error())
		updateMessageFailedStatus(s, &c, onlineUserOnly)
		return err
	} else {
		err = json.Unmarshal(bMsg, &r)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			updateMessageFailedStatus(s, &c, onlineUserOnly)
			return err
		} else {
			if r.ErrCode != 0 {
				sdkLog("errcode, errmsg: ", r.ErrCode, r.ErrMsg)
				updateMessageFailedStatus(s, &c, onlineUserOnly)
				return err
			} else {
				if !onlineUserOnly {
					_ = updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.SendTime, MsgStatusSendSuccess)
				}
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				s.Status = MsgStatusSendSuccess
				c.LatestMsg = structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				if isUpdateConversationLatestMsg {
					_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
						c})
					_ = triggerCmdUpdateConversation(updateConNode{c.ConversationID, IncrUnread, ""})
				}
				if isUpdateConversationInfo {
					_ = triggerCmdUpdateConversation(updateConNode{conversationID, UpdateFaceUrlAndNickName, c})

				}
				if isUpdateConversationInfo || isUpdateConversationLatestMsg {
					_ = triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
					_ = triggerCmdUpdateConversation(updateConNode{"", TotalUnreadMessageChanged, ""})
				}
			}
		}
	}
	return nil
}

func updateMessageFailedStatus(s *MsgStruct, c *ConversationStruct, onlineUserOnly bool) {
	if !onlineUserOnly {
		_ = updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
	}
	s.SendTime = s.CreateTime
	s.Status = MsgStatusSendFailed
	c.LatestMsg = structToJsonString(s)
}
func initBasicInfo(message *MsgStruct, msgFrom, contentType int32) {
	message.CreateTime = getCurrentTimestampByNano()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = MsgStatusSending
	message.SendID = LoginUid
	userInfo, _ := getLoginUserInfoFromLocal()
	message.SenderFaceURL = userInfo.Icon
	message.SenderNickName = userInfo.Name
	//Generate client message primary key
	ClientMsgID := getMsgID(message.SendID)
	message.ClientMsgID = ClientMsgID
	message.MsgFrom = msgFrom
	message.ContentType = contentType
	message.PlatformID = SvrConf.Platform
}
func sendMessageFailedHandle(s *MsgStruct, c *ConversationStruct, conversationID string) {
	_ = updateMessageTimeAndMsgIDStatus(s.ClientMsgID, s.CreateTime, MsgStatusSendFailed)
	s.SendTime = s.CreateTime
	s.Status = MsgStatusSendFailed
	c.LatestMsg = structToJsonString(s)
	_ = triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
		*c})
	_ = triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
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
		sdkLog(err.Error())
		return nil, err
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		sdkLog(err.Error())
		return nil, err
	}

	datatype, err := imgtype.Get(filePath)
	if err != nil {
		sdkLog(err.Error())
		return nil, err
	}

	fi, err := os.Stat(filePath)
	if err != nil {
		sdkLog(err.Error())
		return nil, err
	}

	b := img.Bounds()

	return &imageInfo{int32(b.Max.X), int32(b.Max.Y), datatype, fi.Size()}, nil

}
