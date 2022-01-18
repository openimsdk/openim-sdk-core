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

func (c *Conversation) initBasicInfo(message *sdk_struct.MsgStruct, msgFrom, contentType int32) {
	message.CreateTime = utils.GetCurrentTimestampByMill()
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

type MsgFormats []*db.LocalChatLog

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
