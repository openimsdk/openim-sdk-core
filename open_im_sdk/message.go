/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/6/1 19:16).
 */
package open_im_sdk

import (
	"encoding/json"

	imgtype "github.com/shamsher31/goimgtype"
	"image"
	"os"
)

const TimeOffset = 5

func doNewMsgConversation() {

}

func internalSendMsg(callback Base, s MsgStruct, receiver, groupID string, onlineUserOnly bool) {
	r := SendMsgRespFromServer{}
	a := paramsUserSendMsg{}
	m := make(map[string]interface{})

	switch s.ContentType {
	case Text:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	case Picture:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	case Sound:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	case Video:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	case File:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	case C2CMessageAsRead:
	case RevokeMessageTip:
		ConListener.MsgListenerList[0].OnRecvMessageRevoked(s.RevokeMessage.ServerMsgID)
	default:
		sdkLog("err type: ", s.ContentType)
		callback.OnError(ErrCodeConversation, "err type")
		return
	}

	if receiver == "" {
		s.SessionType = GroupChatType
		s.RecvID = groupID
	} else if groupID == "" {
		s.SessionType = SingleChatType
		s.RecvID = receiver
	} else {
		sdkLog("args err: ", receiver, groupID)
		callback.OnError(ErrCodeConversation, "args, err")
		return
	}
	go func() {
		//Protocol conversion
		a.ReqIdentifier = 1003
		a.PlatformID = s.PlatformID
		a.SendID = s.SendID
		a.OperationID = operationIDGenerator()
		a.Data.SessionType = s.SessionType
		a.Data.MsgFrom = s.MsgFrom
		a.MsgIncr = 1
		a.Data.ForceList = []string{}
		a.Data.ContentType = s.ContentType
		a.Data.RecvID = s.RecvID
		a.Data.ForceList = s.ForceList
		a.Data.Content = s.Content
		a.Data.ClientMsgID = s.ClientMsgID
		if onlineUserOnly {
			a.Data.Options["history"] = 0
			a.Data.Options["persistent"] = 0
		} else {
			a.Data.Options = m
		}
		a.Data.OffLineInfo = m
		bMsg, err := post2Api(sendMsgRouter, a, token)
		if err != nil {
			callback.OnError(ErrCodeConversation, err.Error())
			return
		} else {
			err = json.Unmarshal(bMsg, &r)
			if err != nil {
				callback.OnError(ErrCodeConversation, err.Error())
				return
			} else {
				if r.ErrCode != 0 {
					callback.OnError(ErrCodeConversation, r.ErrMsg)
					return
				}
			}
		}
		callback.OnSuccess("")
		return
	}()
}

func initBasicInfo(message *MsgStruct, msgFrom, contentType int32) {
	message.CreateTime = getCurrentTimestampBySecond()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = MsgStatusSending
	message.SendID = LoginUid
	userInfo, _ := getUserInfoFromLocal()
	message.SenderFaceURL = userInfo.Icon
	message.SenderNickName = userInfo.Name
	//Generate client message primary key
	ClientMsgID := getMsgID(message.SendID)
	message.ClientMsgID = ClientMsgID
	message.MsgFrom = msgFrom
	if contentType != 0 {
		message.ContentType = contentType
	}
	message.PlatformID = SvrConf.Platform
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
