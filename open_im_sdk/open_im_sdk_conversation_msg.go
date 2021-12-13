package open_im_sdk

import (
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

func (u *UserRelated) GetAllConversationList(callback Base) {
	go func() {
		err, list := u.getAllConversationListModel()
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(structToJsonString(list))
			} else {
				callback.OnSuccess(structToJsonString([]ConversationStruct{}))
			}
		}
	}()
}
func (u *UserRelated) SetConversationRecvMessageOpt(callback Base, conversationIDList string, opt int) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := post2Api(setReceiveMessageOptRouter, paramsSetReceiveMessageOpt{OperationID: operationIDGenerator(), Option: int32(opt), ConversationIdList: list}, u.token)
		if err != nil {
			sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g getReceiveMessageOptResp
		err = json.Unmarshal(resp, &g)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		if g.ErrCode != 0 {
			sdkLog("errcode: ", g.ErrCode, g.ErrMsg)
			callback.OnError(int(g.ErrCode), g.ErrMsg)
			return
		}
		u.receiveMessageOptMutex.Lock()
		for _, v := range list {
			u.receiveMessageOpt[v] = int32(opt)
		}
		u.receiveMessageOptMutex.Unlock()
		_ = u.setMultipleConversationRecvMsgOpt(list, opt)
		callback.OnSuccess("")
		//_ = u.triggerCmdUpdateConversation(updateConNode{Action: ConChange})
		u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, list}})
	}()
}
func (u *UserRelated) GetConversationRecvMessageOpt(callback Base, conversationIDList string) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := post2Api(getReceiveMessageOptRouter, paramGetReceiveMessageOpt{OperationID: operationIDGenerator(), ConversationIdList: list}, u.token)
		if err != nil {
			sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g getReceiveMessageOptResp
		err = json.Unmarshal(resp, &g)
		if err != nil {
			sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		if g.ErrCode != 0 {
			sdkLog("errcode: ", g.ErrCode, g.ErrMsg)
			callback.OnError(int(g.ErrCode), g.ErrMsg)
			return
		}
		callback.OnSuccess(structToJsonString(g.Data))
	}()
}
func (u *UserRelated) GetOneConversation(sourceID string, sessionType int, callback Base) {
	go func() {
		conversationID := GetConversationIDBySessionType(sourceID, sessionType)
		err, c := u.getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			//
			if c.ConversationID == "" {
				c.ConversationID = conversationID
				c.ConversationType = sessionType
				switch sessionType {
				case SingleChatType:
					c.UserID = sourceID
					faceUrl, name, err := u.getUserNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						sdkLog("getUserNameAndFaceUrlByUid err:", err)
						return
					}
					c.ShowName = name
					c.FaceURL = faceUrl
				case GroupChatType:
					c.GroupID = sourceID
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					}
					c.ShowName = name
					c.FaceURL = faceUrl

				}
				err = u.addConversationOrUpdateLatestMsg(&c, conversationID)
				if err != nil {
					callback.OnError(301, err.Error())
					return
				}
				callback.OnSuccess(structToJsonString(c))

			} else {
				callback.OnSuccess(structToJsonString(c))
			}
		}
	}()
}
func (u *UserRelated) GetMultipleConversation(conversationIDList string, callback Base) {
	go func() {
		var c []string
		err := json.Unmarshal([]byte(conversationIDList), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			sdkLog("Unmarshal failed", err.Error())
		}
		err, list := u.getMultipleConversationModel(c)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(structToJsonString(list))
			} else {
				callback.OnSuccess(structToJsonString([]ConversationStruct{}))
			}
		}
	}()
}
func (u *UserRelated) DeleteConversation(conversationID string, callback Base) {
	go func() {
		//Transaction operation required
		var sourceID string
		err, c := u.getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(201, err.Error())
			return
		}
		switch c.ConversationType {
		case SingleChatType:
			sourceID = c.UserID
		case GroupChatType:
			sourceID = c.GroupID
		}
		//Mark messages related to this conversation for deletion
		err = u.setMessageStatusBySourceID(sourceID, MsgStatusHasDeleted, c.ConversationType)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		//Reset the session information, empty session
		err = u.ResetConversation(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
			return
		} else {
			callback.OnSuccess("")
			_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: TotalUnreadMessageChanged})
		}
	}()
}
func (u *UserRelated) SetConversationDraft(conversationID, draftText string, callback Base) {
	var time int64
	if draftText == "" {
		time = 0
	} else {
		time = getCurrentTimestampByNano()
	}
	err := u.setConversationDraftModel(conversationID, draftText, time)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess("")
		//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
	}
}
func (u *UserRelated) PinConversation(conversationID string, isPinned bool, callback Base) {
	var i int
	if isPinned {
		i = 1
	} else {
		i = 0
	}
	err := u.pinConversationModel(conversationID, i)
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess("")
		//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConChange})
	}

}
func (u *UserRelated) GetTotalUnreadMsgCount(callback Base) {
	count, err := u.getTotalUnreadMsgCountModel()
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess(int32ToString(count))
	}

}

type OnConversationListener interface {
	OnSyncServerStart()
	OnSyncServerFinish()
	OnSyncServerFailed()
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int32)
}

func (u *UserRelated) SetConversationListener(listener OnConversationListener) {
	if u.ConversationListenerx != nil {
		sdkLog("only one ")
		return
	}
	u.ConversationListenerx = listener
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnRecvMessageRevoked(msgId string)
}

func (u *UserRelated) AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
	if listener == nil {
		sdkLog("AddAdvancedMsgListener listener is null")
		return
	}
	if len(u.ConversationListener.MsgListenerList) == 1 {
		sdkLog("u.ConversationListener.MsgListenerList == 1")
		return
	}
	u.ConversationListener.MsgListenerList = append(u.ConversationListener.MsgListenerList, listener)
}

func (u *UserRelated) ForceSyncMsg() bool {
	if u.syncSeq2Msg() == nil {
		return true
	} else {
		return false
	}
}

func (u *UserRelated) ForceSyncJoinedGroup() {
	u.syncJoinedGroupInfo()
}

func (u *UserRelated) ForceSyncJoinedGroupMember() {

	u.syncJoinedGroupMember()
}

func (u *UserRelated) ForceSyncGroupRequest() {
	u.syncGroupRequest()
}

func (u *UserRelated) ForceSyncApplyGroupRequest() {
	u.syncApplyGroupRequest()
}

type SendMsgCallBack interface {
	Base
	OnProgress(progress int)
}

func (u *UserRelated) CreateTextMessage(text string) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Text)
	s.Content = text
	s.AtElem.AtUserList = []string{}
	return structToJsonString(s)
}
func (u *UserRelated) CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, AtText)
	s.ForceList = users
	s.AtElem.Text = text
	s.AtElem.AtUserList = users
	s.Content = structToJsonString(s.AtElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateLocationMessage(description string, longitude, latitude float64) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Location)
	s.AtElem.AtUserList = []string{}
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = structToJsonString(s.LocationElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateCustomMessage(data, extension string, description string) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Custom)
	s.AtElem.AtUserList = []string{}
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = structToJsonString(s.CustomElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateQuoteMessage(text string, message string) string {
	s, qs := MsgStruct{}, MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	u.initBasicInfo(&s, UserMsgType, Quote)
	//Avoid nested references
	if qs.ContentType == Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = Text
	}
	s.AtElem.AtUserList = []string{}
	s.QuoteElem.Text = text
	s.QuoteElem.QuoteMessage = &qs
	s.Content = structToJsonString(s.QuoteElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateCardMessage(cardInfo string) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Card)
	s.Content = cardInfo
	s.AtElem.AtUserList = []string{}
	return structToJsonString(s)

}
func (u *UserRelated) CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(videoFullPath) //a->b
		s, err := copyFile(videoFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, videoFullPath)
		}
		sdkLog("videoFullPath dstFile", videoFullPath, dstFile, s)
		dstFile = fileTmpPath(snapshotFullPath) //a->b
		s, err = copyFile(snapshotFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, snapshotFullPath)
		}
		sdkLog("snapshotFullPath dstFile", snapshotFullPath, dstFile, s)
		wg.Done()
	}()

	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Video)
	s.VideoElem.VideoPath = videoFullPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotFullPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = snapshotFullPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	wg.Wait()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.VideoElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateFileMessageFromFullPath(fileFullPath string, fileName string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(fileFullPath)
		_, err := copyFile(fileFullPath, dstFile)
		sdkLog("copy file, ", fileFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, fileFullPath)

		}
		wg.Done()
	}()
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, File)
	s.FileElem.FilePath = fileFullPath
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		sdkLog("get file info err:", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.FileElem.FileName = fileName
	s.AtElem.AtUserList = []string{}
	return structToJsonString(s)
}
func (u *UserRelated) CreateImageMessageFromFullPath(imageFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(imageFullPath) //a->b
		_, err := copyFile(imageFullPath, dstFile)
		sdkLog("copy file, ", imageFullPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, imageFullPath)
		}
		wg.Done()
	}()

	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Picture)
	s.PictureElem.SourcePath = imageFullPath
	sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		sdkLog("getImageInfo err:", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	wg.Wait()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.PictureElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := fileTmpPath(soundPath) //a->b
		_, err := copyFile(soundPath, dstFile)
		sdkLog("copy file, ", soundPath, dstFile)
		if err != nil {
			sdkLog("open file failed: ", err, soundPath)
		}
		wg.Done()
	}()
	sdkLog("init base info ")
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Voice)
	s.SoundElem.SoundPath = soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		sdkLog(err.Error(), s.SoundElem.SoundPath)
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	wg.Wait()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.SoundElem)
	sdkLog("to string")
	return structToJsonString(s)
}
func (u *UserRelated) CreateImageMessage(imagePath string) string {
	sdkLog("start1: ", time.Now())
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Picture)
	s.PictureElem.SourcePath = SvrConf.DbDir + imagePath
	sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	sdkLog("end1", time.Now())

	sdkLog("start2 ", time.Now())
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		sdkLog("CreateImageMessage err:", err.Error())
		return ""
	}
	sdkLog("end2", time.Now())

	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.PictureElem)
	return structToJsonString(s)
}

func (u *UserRelated) CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	s := MsgStruct{}
	var p PictureBaseInfo
	_ = json.Unmarshal([]byte(sourcePicture), &p)
	s.PictureElem.SourcePicture = p
	_ = json.Unmarshal([]byte(bigPicture), &p)
	s.PictureElem.BigPicture = p
	_ = json.Unmarshal([]byte(snapshotPicture), &p)
	s.PictureElem.SnapshotPicture = p
	u.initBasicInfo(&s, UserMsgType, Picture)
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.PictureElem)
	return structToJsonString(s)
}

func (u *UserRelated) SendMessageNotOss(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool) string {
	var conversationID string
	r := SendMsgRespFromServer{}
	a := paramsUserSendMsg{}
	s := MsgStruct{}
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(2038, err.Error())
		sdkLog("json unmarshal err:", err.Error())
		return ""
	}

	go func() {
		c := ConversationStruct{
			LatestMsgSendTime: s.CreateTime,
		}
		if receiver == "" && groupID == "" {
			callback.OnError(201, "args err")
			return
		} else if receiver == "" {
			s.SessionType = GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = GetConversationIDBySessionType(groupID, GroupChatType)
			c.GroupID = groupID
			c.ConversationType = GroupChatType
			faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
			if err != nil {
				sdkLog("getGroupNameAndFaceUrlByUid err:", err)
				callback.OnError(202, err.Error())
				return
			}
			c.ShowName = name
			c.FaceURL = faceUrl
			groupMemberList, err := u.getLocalGroupMemberListByGroupID(groupID)
			if err != nil {
				sdkLog("getLocalGroupMemberListByGroupID err:", err)
				callback.OnError(202, err.Error())
				return
			}
			isExistInGroup := func(target string, groupMemberList []groupMemberFullInfo) bool {

				for _, element := range groupMemberList {

					if target == element.UserId {
						return true
					}
				}
				return false

			}(s.SendID, groupMemberList)
			if !isExistInGroup {
				sdkLog("SendGroupMessage err:", "not exist in this group")
				callback.OnError(208, "not exist in this group")
				return
			}

		} else {
			s.SessionType = SingleChatType
			s.RecvID = receiver
			conversationID = GetConversationIDBySessionType(receiver, SingleChatType)
			c.UserID = receiver
			c.ConversationType = SingleChatType
			faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
			if err != nil {
				sdkLog("getUserNameAndFaceUrlByUid err:", err)
				callback.OnError(301, err.Error())
				return
			}
			c.FaceURL = faceUrl
			c.ShowName = name
		}
		userInfo, err := u.getLoginUserInfoFromLocal()
		if err != nil {
			sdkLog("getLoginUserInfoFromLocal err:", err)
			return
		}
		s.SenderFaceURL = userInfo.Icon
		s.SenderNickName = userInfo.Name
		c.ConversationID = conversationID
		c.LatestMsg = structToJsonString(s)
		err = u.insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			callback.OnError(202, err.Error())
			return
		}
		_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
			c})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})

		//Protocol conversion
		a.ReqIdentifier = 1003
		a.PlatformID = s.PlatformID
		a.SendID = s.SendID
		a.SenderFaceURL = s.SenderFaceURL
		a.SenderNickName = s.SenderNickName
		a.OperationID = operationIDGenerator()
		a.Data.SessionType = s.SessionType
		a.Data.MsgFrom = s.MsgFrom
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
		bMsg, err := post2Api(sendMsgRouter, a, u.token)
		if err != nil {
			callback.OnError(http.StatusInternalServerError, err.Error())
			u.sendMessageFailedHandle(&s, &c, conversationID)
		} else if err = json.Unmarshal(bMsg, &r); err != nil {
			callback.OnError(200, err.Error()+"  "+string(bMsg))
			u.sendMessageFailedHandle(&s, &c, conversationID)
		} else {
			if r.ErrCode != 0 {
				callback.OnError(r.ErrCode, r.ErrMsg)
				u.sendMessageFailedHandle(&s, &c, conversationID)
			} else {
				callback.OnSuccess("")
				callback.OnProgress(100)
				_ = u.updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.SendTime, MsgStatusSendSuccess)
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				s.Status = MsgStatusSendSuccess
				c.LatestMsg = structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
					c})
				u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
			}
		}
	}()
	return s.ClientMsgID
}
func (u *UserRelated) CreateSoundMessageByURL(soundBaseInfo string) string {
	s := MsgStruct{}
	var soundElem SoundBaseInfo
	_ = json.Unmarshal([]byte(soundBaseInfo), &soundElem)
	s.SoundElem = soundElem
	u.initBasicInfo(&s, UserMsgType, Voice)
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.SoundElem)
	return structToJsonString(s)
}

func (u *UserRelated) CreateSoundMessage(soundPath string, duration int64) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Voice)
	s.SoundElem.SoundPath = SvrConf.DbDir + soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.SoundElem)
	return structToJsonString(s)
}

func (u *UserRelated) CreateVideoMessageByURL(videoBaseInfo string) string {
	s := MsgStruct{}
	var videoElem VideoBaseInfo
	_ = json.Unmarshal([]byte(videoBaseInfo), &videoElem)
	s.VideoElem = videoElem
	u.initBasicInfo(&s, UserMsgType, Video)
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.VideoElem)
	return structToJsonString(s)
}

func (u *UserRelated) CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, Video)
	s.VideoElem.VideoPath = SvrConf.DbDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = SvrConf.DbDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.VideoElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateFileMessageByURL(fileBaseInfo string) string {
	s := MsgStruct{}
	var fileElem FileBaseInfo
	_ = json.Unmarshal([]byte(fileBaseInfo), &fileElem)
	s.FileElem = fileElem
	u.initBasicInfo(&s, UserMsgType, File)
	s.AtElem.AtUserList = []string{}
	s.Content = structToJsonString(s.FileElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateFileMessage(filePath string, fileName string) string {
	s := MsgStruct{}
	u.initBasicInfo(&s, UserMsgType, File)
	s.FileElem.FilePath = SvrConf.DbDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		sdkLog(err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.AtElem.AtUserList = []string{}

	sdkLog("CreateForwardMessage new input: ", structToJsonString(s))
	return structToJsonString(s)
}
func (u *UserRelated) CreateMergerMessage(messageList, title, summaryList string) string {
	var messages []*MsgStruct
	var summaries []string
	s := MsgStruct{}
	err := json.Unmarshal([]byte(messageList), &messages)
	if err != nil {
		sdkLog("CreateMergerMessage err:", err.Error())
		return ""
	}
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	u.initBasicInfo(&s, UserMsgType, Merger)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = structToJsonString(s.MergeElem)
	return structToJsonString(s)
}
func (u *UserRelated) CreateForwardMessage(m string) string {
	sdkLog("CreateForwardMessage input: ", m)
	s := MsgStruct{}
	err := json.Unmarshal([]byte(m), &s)
	if err != nil {
		sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	if s.Status != MsgStatusSendSuccess {
		sdkLog("only send success message can be revoked")
		return ""
	}

	u.initBasicInfo(&s, UserMsgType, s.ContentType)
	//Forward message seq is set to 0
	s.Seq = 0
	return structToJsonString(s)
}

func (u *UserRelated) SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool) string {
	var conversationID string
	//r := SendMsgRespFromServer{}
	//a := paramsUserSendMsg{}
	s := MsgStruct{}
	//m := make(map[string]interface{})
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(2038, err.Error())
		sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	go func() {
		c := ConversationStruct{
			LatestMsgSendTime: s.CreateTime,
		}
		if receiver == "" && groupID == "" {
			callback.OnError(201, "args err")
			return
		} else if receiver == "" {
			s.SessionType = GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = GetConversationIDBySessionType(groupID, GroupChatType)
			c.GroupID = groupID
			c.ConversationType = GroupChatType
			faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
			if err != nil {
				sdkLog("getGroupNameAndFaceUrlByUid err:", err)
				callback.OnError(202, err.Error())
				return
			}
			c.ShowName = name
			c.FaceURL = faceUrl
			groupMemberList, err := u.getLocalGroupMemberListByGroupID(groupID)
			if err != nil {
				sdkLog("getLocalGroupMemberListByGroupID err:", err)
				callback.OnError(202, err.Error())
				return
			}
			isExistInGroup := func(target string, groupMemberList []groupMemberFullInfo) bool {

				for _, element := range groupMemberList {

					if target == element.UserId {
						return true
					}
				}
				return false

			}(s.SendID, groupMemberList)
			if !isExistInGroup {
				sdkLog("SendGroupMessage err:", "not exist in this group")
				callback.OnError(208, "not exist in this group")
				return
			}

		} else {
			s.SessionType = SingleChatType
			s.RecvID = receiver
			conversationID = GetConversationIDBySessionType(receiver, SingleChatType)
			c.UserID = receiver
			c.ConversationType = SingleChatType
			faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
			if err != nil {
				sdkLog("getUserNameAndFaceUrlByUid err:", err)
				callback.OnError(301, err.Error())
				return
			}
			c.FaceURL = faceUrl
			c.ShowName = name

		}
		userInfo, err := u.getLoginUserInfoFromLocal()
		if err != nil {
			sdkLog("getLoginUserInfoFromLocal err:", err)
			return
		}
		s.SenderFaceURL = userInfo.Icon
		s.SenderNickName = userInfo.Name
		c.ConversationID = conversationID
		c.LatestMsg = structToJsonString(s)
		err = u.insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			callback.OnError(202, err.Error())
			return
		}
		_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
			c})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		var delFile []string
		switch s.ContentType {
		case Text:
		case AtText:
		case Location:
		case Custom:
		case Merger:
		case Quote:
		case Card:
		case Picture:
			var sourcePath string
			if fileExist(s.PictureElem.SourcePath) {
				sourcePath = s.PictureElem.SourcePath
				delFile = append(delFile, fileTmpPath(s.PictureElem.SourcePath))
			} else {
				sourcePath = fileTmpPath(s.PictureElem.SourcePath)
				delFile = append(delFile, sourcePath)
			}
			sdkLog("file: ", sourcePath, delFile)
			sourceUrl, uuid, err := u.uploadImage(sourcePath, callback)
			if err != nil {
				sdkLog("oss Picture upload err", err.Error())
				callback.OnError(301, err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.PictureElem.SourcePicture.Url = sourceUrl
				s.PictureElem.SourcePicture.UUID = uuid
				s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + ZoomScale + "/h/" + ZoomScale
				s.PictureElem.SnapshotPicture.Width = int32(stringToInt(ZoomScale))
				s.PictureElem.SnapshotPicture.Height = int32(stringToInt(ZoomScale))
				s.Content = structToJsonString(s.PictureElem)
			}
		case Voice:
			var sourcePath string
			if fileExist(s.SoundElem.SoundPath) {
				sourcePath = s.SoundElem.SoundPath
				delFile = append(delFile, fileTmpPath(s.SoundElem.SoundPath))
			} else {
				sourcePath = fileTmpPath(s.SoundElem.SoundPath)
				delFile = append(delFile, sourcePath)
			}
			sdkLog("file: ", sourcePath, delFile)
			soundURL, uuid, err := u.uploadSound(sourcePath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				sdkLog("uploadSound err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = structToJsonString(s.SoundElem)
			}
		case Video:
			var videoPath string
			var snapPath string
			if fileExist(s.VideoElem.VideoPath) {
				videoPath = s.VideoElem.VideoPath
				snapPath = s.VideoElem.SnapshotPath
				delFile = append(delFile, fileTmpPath(s.VideoElem.VideoPath))
				delFile = append(delFile, fileTmpPath(s.VideoElem.SnapshotPath))
			} else {
				videoPath = fileTmpPath(s.VideoElem.VideoPath)
				snapPath = fileTmpPath(s.VideoElem.SnapshotPath)
				delFile = append(delFile, videoPath)
				delFile = append(delFile, snapPath)
			}
			sdkLog("file: ", videoPath, snapPath, delFile)
			snapshotURL, snapshotUUID, videoURL, videoUUID, err := u.uploadVideo(videoPath, snapPath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				sdkLog("oss  Video upload err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.VideoElem.VideoURL = videoURL
				s.VideoElem.SnapshotUUID = snapshotUUID
				s.VideoElem.SnapshotURL = snapshotURL
				s.VideoElem.VideoUUID = videoUUID
				s.Content = structToJsonString(s.VideoElem)
			}
		case File:
			fileURL, fileUUID, err := u.uploadFile(s.FileElem.FilePath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				sdkLog("oss  File upload err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return

			} else {
				s.FileElem.SourceURL = fileURL
				s.FileElem.UUID = fileUUID
				s.Content = structToJsonString(s.FileElem)
			}
		default:
			callback.OnError(2038, "Not currently supported ")
			sdkLog("Not currently supported ", s.ContentType)
			return
		}
		//Store messages to local database
		err = u.insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		//Protocol conversion
		//a.ReqIdentifier = 1003
		//a.PlatformID = s.PlatformID
		//a.SendID = s.SendID
		//a.SenderFaceURL = s.SenderFaceURL
		//a.SenderNickName = s.SenderNickName
		//a.OperationID = operationIDGenerator()
		//a.Data.SessionType = s.SessionType
		//a.Data.MsgFrom = s.MsgFrom
		//a.Data.ContentType = s.ContentType
		//a.Data.RecvID = s.RecvID
		//a.Data.ForceList = s.ForceList
		//a.Data.Content = s.Content
		//a.Data.ClientMsgID = s.ClientMsgID
		//if onlineUserOnly {
		//	a.Data.Options["history"] = 0
		//	a.Data.Options["persistent"] = 0
		//} else {
		//	a.Data.Options = m
		//}
		//a.Data.OffLineInfo = m

		optionsFlag := make(map[string]int32, 2)
		if onlineUserOnly {
			optionsFlag["history"] = 0
			optionsFlag["persistent"] = 0
		}
		wsMsgData := UserSendMsgReq{
			Options:        optionsFlag,
			SenderNickName: s.SenderNickName,
			SenderFaceURL:  s.SenderFaceURL,
			PlatformID:     s.PlatformID,
			SessionType:    s.SessionType,
			MsgFrom:        s.MsgFrom,
			ContentType:    s.ContentType,
			RecvID:         s.RecvID,
			ForceList:      s.ForceList,
			Content:        s.Content,
			ClientMsgID:    s.ClientMsgID,
		}
		msgIncr, ch := u.AddCh()
		var wsReq GeneralWsReq
		wsReq.ReqIdentifier = WSSendMsg
		wsReq.OperationID = operationIDGenerator()
		wsReq.SendID = s.SendID
		//wsReq.Token = u.token
		wsReq.MsgIncr = msgIncr
		wsReq.Data, err = proto.Marshal(&wsMsgData)
		if err != nil {
			sdkLog("Marshal failed ", err.Error())
			LogFReturn(nil)
			callback.OnError(http.StatusInternalServerError, err.Error())
			u.sendMessageFailedHandle(&s, &c, conversationID)
			return
		}

		SendFlag := false
		var connSend *websocket.Conn
		for tr := 0; tr < 30; tr++ {
			LogBegin("WriteMsg", wsReq.OperationID)
			err, connSend = u.WriteMsg(wsReq)
			LogEnd("WriteMsg ", wsReq.OperationID, connSend)
			if err != nil {
				sdkLog("ws writeMsg  err:,", wsReq.OperationID, err.Error(), tr)
				time.Sleep(time.Duration(5) * time.Second)
			} else {
				sdkLog("writeMsg  retry ok", wsReq.OperationID, tr)
				SendFlag = true
				break
			}
		}

		if SendFlag == false {
			u.DelCh(msgIncr)
			callback.OnError(http.StatusInternalServerError, err.Error())
			u.sendMessageFailedHandle(&s, &c, conversationID)
			return
		}

		timeout := 300
		breakFlag := 0
		for {
			if breakFlag == 1 {
				sdkLog("break ", wsReq.OperationID)
				break
			}
			select {
			case r := <-ch:
				sdkLog("ws  ch recvMsg success:,", wsReq.OperationID)
				if r.ErrCode != 0 {
					callback.OnError(r.ErrCode, r.ErrMsg)
					u.sendMessageFailedHandle(&s, &c, conversationID)
				} else {
					callback.OnSuccess("")
					callback.OnProgress(100)
					for _, v := range delFile {
						err := os.Remove(v)
						if err != nil {
							sdkLog("remove failed,", err.Error(), v)
						}
						sdkLog("remove file: ", v)
					}
					var sendMsgResp UserSendMsgResp
					err = proto.Unmarshal(r.Data, &sendMsgResp)
					if err != nil {
						sdkLog("Unmarshal failed ", err.Error())
						//	callback.OnError(http.StatusInternalServerError, err.Error())
						//	u.sendMessageFailedHandle(&s, &c, conversationID)
						//	u.DelCh(msgIncr)
					}
					_ = u.updateMessageTimeAndMsgIDStatus(sendMsgResp.ClientMsgID, sendMsgResp.SendTime, MsgStatusSendSuccess)

					s.ServerMsgID = sendMsgResp.ServerMsgID
					s.SendTime = sendMsgResp.SendTime
					s.Status = MsgStatusSendSuccess
					c.LatestMsg = structToJsonString(s)
					c.LatestMsgSendTime = s.SendTime
					_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, AddConOrUpLatMsg,
						c})
					u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
				}
				breakFlag = 1
			case <-time.After(time.Second * time.Duration(timeout)):
				var flag bool
				sdkLog("ws ch recvMsg err: ", wsReq.OperationID)
				if connSend != u.conn {
					sdkLog("old conn != current conn  ", connSend, u.conn)
					flag = false // error
				} else {
					flag = false //error
					for tr := 0; tr < 3; tr++ {
						err = u.sendPingMsg()
						if err != nil {
							sdkLog("sendPingMsg failed ", wsReq.OperationID, err.Error(), tr)
							time.Sleep(time.Duration(30) * time.Second)
						} else {
							sdkLog("sendPingMsg ok ", wsReq.OperationID)
							flag = true //wait continue
							break
						}
					}
				}
				if flag == false {
					callback.OnError(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
					u.sendMessageFailedHandle(&s, &c, conversationID)
					sdkLog("onError callback ", wsReq.OperationID)
					breakFlag = 1
					break
				} else {
					sdkLog("wait resp continue", wsReq.OperationID)
					breakFlag = 0
					continue
				}
			}
		}

		u.DelCh(msgIncr)
	}()
	return s.ClientMsgID
}

func (u *UserRelated) GetHistoryMessageList(callback Base, getMessageOptions string) {
	go func() {
		sdkLog("GetHistoryMessageList", getMessageOptions)
		var sourceID string
		var conversationID string
		var startTime int64
		var latestMsg MsgStruct
		var sessionType int
		p := PullMsgReq{}
		err := json.Unmarshal([]byte(getMessageOptions), &p)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		if p.UserID == "" {
			sourceID = p.GroupID
			conversationID = GetConversationIDBySessionType(sourceID, GroupChatType)
			sessionType = GroupChatType
		} else {
			sourceID = p.UserID
			conversationID = GetConversationIDBySessionType(sourceID, SingleChatType)
			sessionType = SingleChatType
		}
		if p.StartMsg == nil {
			err, m := u.getConversationLatestMsgModel(conversationID)
			if err != nil {
				callback.OnError(200, err.Error())
				return
			}
			if m == "" {
				startTime = 0
			} else {
				err := json.Unmarshal([]byte(m), &latestMsg)
				if err != nil {
					sdkLog("get history err :", err)
					callback.OnError(200, err.Error())
					return
				}
				startTime = latestMsg.SendTime + TimeOffset
			}

		} else {
			startTime = p.StartMsg.SendTime
		}
		sdkLog("sourceID:", sourceID, "startTime:", startTime, "count:", p.Count)
		err, list := u.getHistoryMessage(sourceID, startTime, p.Count, sessionType)
		sort.Sort(list)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(structToJsonString(list))
			} else {
				callback.OnSuccess(structToJsonString([]MsgStruct{}))
			}
		}
	}()
}
func (u *UserRelated) RevokeMessage(callback Base, message string) {
	go func() {
		var receiver, groupID string
		c := MsgStruct{}
		err := json.Unmarshal([]byte(message), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		s, err := u.getOneMessage(c.ClientMsgID)
		if err != nil || s == nil {
			callback.OnError(201, "getOneMessage err")
			return
		}
		if s.Status != MsgStatusSendSuccess {
			callback.OnError(201, "only send success message can be revoked")
			return
		}
		sdkLog("test data", s)
		//Send message internally
		switch s.SessionType {
		case SingleChatType:
			receiver = s.RecvID
		case GroupChatType:
			groupID = s.RecvID
		default:
			callback.OnError(200, "args err")
		}
		s.Content = s.ClientMsgID
		s.ClientMsgID = getMsgID(s.SendID)
		s.ContentType = Revoke
		err = u.autoSendMsg(s, receiver, groupID, false, true, false)
		if err != nil {
			sdkLog("autoSendMsg revokeMessage err:", err.Error())
			callback.OnError(300, err.Error())

		} else {
			err = u.setMessageStatus(s.Content, MsgStatusRevoked)
			if err != nil {
				sdkLog("setLocalMessageStatus revokeMessage err:", err.Error())
				callback.OnError(300, err.Error())
			} else {
				callback.OnSuccess("")
			}
		}
	}()
}
func (u *UserRelated) TypingStatusUpdate(receiver, msgTip string) {
	go func() {
		s := MsgStruct{}
		u.initBasicInfo(&s, UserMsgType, Typing)
		s.Content = msgTip
		err := u.autoSendMsg(&s, receiver, "", true, false, false)
		if err != nil {
			sdkLog("TypingStatusUpdate err:", err)
		} else {
			sdkLog("TypingStatusUpdate success!!!")
		}
	}()
}

func (u *UserRelated) MarkC2CMessageAsRead(callback Base, receiver string, msgIDList string) {
	go func() {
		conversationID := GetConversationIDBySessionType(receiver, SingleChatType)
		var list []string
		err := json.Unmarshal([]byte(msgIDList), &list)
		if err != nil {
			callback.OnError(201, "json unmarshal err")
			return
		}
		if len(list) == 0 {
			callback.OnError(200, "msg list is null")
			return
		}
		s := MsgStruct{}
		u.initBasicInfo(&s, UserMsgType, HasReadReceipt)
		s.Content = msgIDList
		sdkLog("MarkC2CMessageAsRead: send Message")
		err = u.autoSendMsg(&s, receiver, "", false, false, false)
		if err != nil {
			sdkLog("MarkC2CMessageAsRead  err:", err.Error())
			callback.OnError(300, err.Error())
		} else {
			callback.OnSuccess("")
			var msgIDs []string
			_ = json.Unmarshal([]byte(msgIDList), &msgIDs)
			_ = u.setSingleMessageHasReadByMsgIDList(receiver, msgIDs)
			u.doUpdateConversation(cmd2Value{Value: updateConNode{conversationID, UpdateLatestMessageChange, ""}})
			u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		}
	}()
}

//Deprecated
func (u *UserRelated) MarkSingleMessageHasRead(callback Base, userID string) {
	go func() {
		conversationID := GetConversationIDBySessionType(userID, SingleChatType)
		//if err := u.setSingleMessageHasRead(userID); err != nil {
		//	callback.OnError(201, err.Error())
		//} else {
		callback.OnSuccess("")
		u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
		u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		//}
	}()
}
func (u *UserRelated) MarkGroupMessageHasRead(callback Base, groupID string) {
	go func() {
		conversationID := GetConversationIDBySessionType(groupID, GroupChatType)
		if err := u.setGroupMessageHasRead(groupID); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
			u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: UnreadCountSetZero})
			u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		}
	}()
}
func (u *UserRelated) DeleteMessageFromLocalStorage(callback Base, message string) {
	go func() {
		var conversation ConversationStruct
		var latestMsg MsgStruct
		var conversationID string
		var sourceID string
		s := MsgStruct{}
		err := json.Unmarshal([]byte(message), &s)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		err = u.setMessageStatus(s.ClientMsgID, MsgStatusHasDeleted)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		callback.OnSuccess("")
		if s.SessionType == GroupChatType {
			conversationID = GetConversationIDBySessionType(s.RecvID, GroupChatType)
			sourceID = s.RecvID

		} else if s.SessionType == SingleChatType {
			if s.SendID != u.LoginUid {
				conversationID = GetConversationIDBySessionType(s.SendID, SingleChatType)
				sourceID = s.SendID
			} else {
				conversationID = GetConversationIDBySessionType(s.RecvID, SingleChatType)
				sourceID = s.RecvID
			}
		}
		_, m := u.getConversationLatestMsgModel(conversationID)
		if m != "" {
			err := json.Unmarshal([]byte(m), &latestMsg)
			if err != nil {
				sdkLog("DeleteMessage err :", err)
				callback.OnError(200, err.Error())
				return
			}
		} else {
			sdkLog("err ,conversation has been deleted")
		}

		if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
			err, list := u.getHistoryMessage(sourceID, s.SendTime+TimeOffset, 1, int(s.SessionType))
			if err != nil {
				sdkLog("DeleteMessageFromLocalStorage database err:", err.Error())
			}
			conversation.ConversationID = conversationID
			if list == nil {
				conversation.LatestMsg = ""
				conversation.LatestMsgSendTime = getCurrentTimestampByNano()
			} else {
				conversation.LatestMsg = structToJsonString(list[0])
				conversation.LatestMsgSendTime = list[0].SendTime
			}
			err = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: AddConOrUpLatMsg, Args: conversation})
			if err != nil {
				sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
			}
			u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})

		}
	}()
}
func (u *UserRelated) ClearC2CHistoryMessage(callback Base, userID string) {
	go func() {
		conversationID := GetConversationIDBySessionType(userID, SingleChatType)
		err := u.setMessageStatusBySourceID(userID, MsgStatusHasDeleted, SingleChatType)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		err = u.clearConversation(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
			return
		} else {
			callback.OnSuccess("")
			_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
		}
	}()
}
func (u *UserRelated) ClearGroupHistoryMessage(callback Base, groupID string) {
	go func() {
		conversationID := GetConversationIDBySessionType(groupID, GroupChatType)
		err := u.setMessageStatusBySourceID(groupID, MsgStatusHasDeleted, GroupChatType)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		err = u.clearConversation(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
			return
		} else {
			callback.OnSuccess("")
			_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
		}
	}()
}

func (u *UserRelated) InsertSingleMessageToLocalStorage(callback Base, message, userID, sender string) string {
	s := MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return ""
	}
	s.SendID = sender
	s.RecvID = userID
	//Generate client message primary key
	s.ClientMsgID = getMsgID(s.SendID)
	s.SendTime = getCurrentTimestampByNano()
	go func() {
		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
		}
	}()
	return s.ClientMsgID
}

func (u *UserRelated) InsertGroupMessageToLocalStorage(callback Base, message, groupID, sender string) string {
	s := MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return ""
	}
	s.SendID = sender
	s.RecvID = groupID
	//Generate client message primary key
	s.ClientMsgID = getMsgID(s.SendID)
	s.SendTime = getCurrentTimestampByNano()
	go func() {
		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
		}
	}()
	return s.ClientMsgID
}

func (u *UserRelated) FindMessages(callback Base, messageIDList string) {
	go func() {
		var c []string
		err := json.Unmarshal([]byte(messageIDList), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			sdkLog("Unmarshal failed, ", err.Error())

		}
		err, list := u.getMultipleMessageModel(c)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(structToJsonString(list))
			} else {
				callback.OnSuccess(structToJsonString([]MsgStruct{}))
			}
		}
	}()
}
