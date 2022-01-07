package conversation_msg

import (
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/http"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/utils"
	"os"
	"sort"
	"sync"
	"time"
)

func (u *open_im_sdk.UserRelated) GetAllConversationList(callback open_im_sdk.Base) {
	go func() {
		err, list := u.getAllConversationListModel()
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(utils.structToJsonString(list))
			} else {
				callback.OnSuccess(utils.structToJsonString([]ConversationStruct{}))
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) GetConversationListSplit(callback open_im_sdk.Base, offset, count int) {
	go func() {
		err, list := u.getConversationListSplitModel(offset, count)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(utils.structToJsonString(list))
			} else {
				callback.OnSuccess(utils.structToJsonString([]ConversationStruct{}))
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) SetConversationRecvMessageOpt(callback open_im_sdk.Base, conversationIDList string, opt int) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := utils.post2Api(open_im_sdk.setReceiveMessageOptRouter, open_im_sdk.paramsSetReceiveMessageOpt{OperationID: utils.operationIDGenerator(), Option: int32(opt), ConversationIdList: list}, u.token)
		if err != nil {
			utils.sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g open_im_sdk.getReceiveMessageOptResp
		err = json.Unmarshal(resp, &g)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		if g.ErrCode != 0 {
			utils.sdkLog("errcode: ", g.ErrCode, g.ErrMsg)
			callback.OnError(g.ErrCode, g.ErrMsg)
			return
		}
		u.receiveMessageOptMutex.Lock()
		for _, v := range list {
			u.receiveMessageOpt[v] = int32(opt)
		}
		u.receiveMessageOptMutex.Unlock()
		_ = u.setMultipleConversationRecvMsgOpt(list, opt)
		callback.OnSuccess("")
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, list}})

	}()
}
func (u *open_im_sdk.UserRelated) GetConversationRecvMessageOpt(callback open_im_sdk.Base, conversationIDList string) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := utils.post2Api(open_im_sdk.getReceiveMessageOptRouter, open_im_sdk.paramGetReceiveMessageOpt{OperationID: utils.operationIDGenerator(), ConversationIdList: list}, u.token)
		if err != nil {
			utils.sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g open_im_sdk.getReceiveMessageOptResp
		err = json.Unmarshal(resp, &g)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		if g.ErrCode != 0 {
			utils.sdkLog("errcode: ", g.ErrCode, g.ErrMsg)
			callback.OnError(g.ErrCode, g.ErrMsg)
			return
		}
		callback.OnSuccess(utils.structToJsonString(g.Data))
	}()
}
func (u *open_im_sdk.UserRelated) GetOneConversation(sourceID string, sessionType int, callback open_im_sdk.Base) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(sourceID, sessionType)
		err, c := u.getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			//
			if c.ConversationID == "" {
				c.ConversationID = conversationID
				c.ConversationType = sessionType
				switch sessionType {
				case open_im_sdk.SingleChatType:
					c.UserID = sourceID
					faceUrl, name, err := u.getUserNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
						return
					}
					c.ShowName = name
					c.FaceURL = faceUrl
				case open_im_sdk.GroupChatType:
					c.GroupID = sourceID
					faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
					}
					c.ShowName = name
					c.FaceURL = faceUrl

				}
				err = u.insertConOrUpdateLatestMsg(&c, conversationID)
				if err != nil {
					callback.OnError(301, err.Error())
					return
				}
				callback.OnSuccess(utils.structToJsonString(c))

			} else {
				callback.OnSuccess(utils.structToJsonString(c))
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) GetMultipleConversation(conversationIDList string, callback open_im_sdk.Base) {
	go func() {
		var c []string
		err := json.Unmarshal([]byte(conversationIDList), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			utils.sdkLog("Unmarshal failed", err.Error())
		}
		err, list := u.getMultipleConversationModel(c)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(utils.structToJsonString(list))
			} else {
				callback.OnSuccess(utils.structToJsonString([]ConversationStruct{}))
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) DeleteConversation(conversationID string, callback open_im_sdk.Base) {
	go func() {
		//Transaction operation required
		var sourceID string
		err, c := u.getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(201, err.Error())
			return
		}
		switch c.ConversationType {
		case open_im_sdk.SingleChatType:
			sourceID = c.UserID
		case open_im_sdk.GroupChatType:
			sourceID = c.GroupID
		}
		//Mark messages related to this conversation for deletion
		err = u.setMessageStatusBySourceID(sourceID, open_im_sdk.MsgStatusHasDeleted, c.ConversationType)
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
			_ = u.triggerCmdUpdateConversation(open_im_sdk.updateConNode{ConId: conversationID, Action: open_im_sdk.TotalUnreadMessageChanged})
		}
	}()
}
func (u *open_im_sdk.UserRelated) SetConversationDraft(conversationID, draftText string, callback open_im_sdk.Base) {
	if draftText != "" {
		err := u.setConversationDraftModel(conversationID, draftText)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	} else {
		err := u.removeConversationDraftModel(conversationID, draftText)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}
}
func (u *open_im_sdk.UserRelated) PinConversation(conversationID string, isPinned bool, callback open_im_sdk.Base) {
	if isPinned {
		err := u.pinConversationModel(conversationID, open_im_sdk.Pinned)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	} else {
		err := u.unPinConversationModel(conversationID, open_im_sdk.NotPinned)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}

}
func (u *open_im_sdk.UserRelated) GetTotalUnreadMsgCount(callback open_im_sdk.Base) {
	count, err := u.getTotalUnreadMsgCountModel()
	if err != nil {
		callback.OnError(203, err.Error())
	} else {
		callback.OnSuccess(utils.int32ToString(count))
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

func (u *open_im_sdk.UserRelated) SetConversationListener(listener OnConversationListener) {
	if u.ConversationListenerx != nil {
		utils.sdkLog("only one ")
		return
	}
	u.ConversationListenerx = listener
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnRecvMessageRevoked(msgId string)
}

func (u *open_im_sdk.UserRelated) AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
	if listener == nil {
		utils.sdkLog("AddAdvancedMsgListener listener is null")
		return
	}
	if len(u.ConversationListener.MsgListenerList) == 1 {
		utils.sdkLog("u.ConversationListener.MsgListenerList == 1")
		return
	}
	u.ConversationListener.MsgListenerList = append(u.ConversationListener.MsgListenerList, listener)
}

func (u *open_im_sdk.UserRelated) ForceSyncMsg() bool {
	if u.syncSeq2Msg() == nil {
		return true
	} else {
		return false
	}
}

func (u *open_im_sdk.UserRelated) ForceSyncJoinedGroup() {
	u.syncJoinedGroupInfo()
}

func (u *open_im_sdk.UserRelated) ForceSyncJoinedGroupMember() {

	u.syncJoinedGroupMember()
}

func (u *open_im_sdk.UserRelated) ForceSyncGroupRequest() {
	u.syncGroupRequest()
}

func (u *open_im_sdk.UserRelated) ForceSyncSelfGroupRequest() {
	u.syncSelfGroupRequest()
}

type SendMsgCallBack interface {
	open_im_sdk.Base
	OnProgress(progress int)
}

func (u *open_im_sdk.UserRelated) CreateTextMessage(text string) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Text)
	s.Content = text
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.AtText)
	s.ForceList = users
	s.AtElem.Text = text
	s.AtElem.AtUserList = users
	s.Content = utils.structToJsonString(s.AtElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateLocationMessage(description string, longitude, latitude float64) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Location)
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = utils.structToJsonString(s.LocationElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateCustomMessage(data, extension string, description string) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Custom)
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = utils.structToJsonString(s.CustomElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateQuoteMessage(text string, message string) string {
	s, qs := open_im_sdk.MsgStruct{}, open_im_sdk.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Quote)
	//Avoid nested references
	if qs.ContentType == open_im_sdk.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = open_im_sdk.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.QuoteMessage = &qs
	s.Content = utils.structToJsonString(s.QuoteElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateCardMessage(cardInfo string) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Card)
	s.Content = cardInfo
	return utils.structToJsonString(s)

}
func (u *open_im_sdk.UserRelated) CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.fileTmpPath(videoFullPath) //a->b
		s, err := utils.copyFile(videoFullPath, dstFile)
		if err != nil {
			utils.sdkLog("open file failed: ", err, videoFullPath)
		}
		utils.sdkLog("videoFullPath dstFile", videoFullPath, dstFile, s)
		dstFile = utils.fileTmpPath(snapshotFullPath) //a->b
		s, err = utils.copyFile(snapshotFullPath, dstFile)
		if err != nil {
			utils.sdkLog("open file failed: ", err, snapshotFullPath)
		}
		utils.sdkLog("snapshotFullPath dstFile", snapshotFullPath, dstFile, s)
		wg.Done()
	}()

	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Video)
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
		utils.sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := open_im_sdk.getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			utils.sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	wg.Wait()
	s.Content = utils.structToJsonString(s.VideoElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateFileMessageFromFullPath(fileFullPath string, fileName string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.fileTmpPath(fileFullPath)
		_, err := utils.copyFile(fileFullPath, dstFile)
		utils.sdkLog("copy file, ", fileFullPath, dstFile)
		if err != nil {
			utils.sdkLog("open file failed: ", err, fileFullPath)

		}
		wg.Done()
	}()
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.File)
	s.FileElem.FilePath = fileFullPath
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		utils.sdkLog("get file info err:", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.FileElem.FileName = fileName
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateImageMessageFromFullPath(imageFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.fileTmpPath(imageFullPath) //a->b
		_, err := utils.copyFile(imageFullPath, dstFile)
		utils.sdkLog("copy file, ", imageFullPath, dstFile)
		if err != nil {
			utils.sdkLog("open file failed: ", err, imageFullPath)
		}
		wg.Done()
	}()

	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Picture)
	s.PictureElem.SourcePath = imageFullPath
	utils.sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := open_im_sdk.getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		utils.sdkLog("getImageInfo err:", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	wg.Wait()
	s.Content = utils.structToJsonString(s.PictureElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.fileTmpPath(soundPath) //a->b
		_, err := utils.copyFile(soundPath, dstFile)
		utils.sdkLog("copy file, ", soundPath, dstFile)
		if err != nil {
			utils.sdkLog("open file failed: ", err, soundPath)
		}
		wg.Done()
	}()
	utils.sdkLog("init base info ")
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Voice)
	s.SoundElem.SoundPath = soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		utils.sdkLog(err.Error(), s.SoundElem.SoundPath)
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	wg.Wait()
	s.Content = utils.structToJsonString(s.SoundElem)
	utils.sdkLog("to string")
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateImageMessage(imagePath string) string {
	utils.sdkLog("start1: ", time.Now())
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Picture)
	s.PictureElem.SourcePath = open_im_sdk.SvrConf.DbDir + imagePath
	utils.sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	utils.sdkLog("end1", time.Now())

	utils.sdkLog("start2 ", time.Now())
	imageInfo, err := open_im_sdk.getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		utils.sdkLog("CreateImageMessage err:", err.Error())
		return ""
	}
	utils.sdkLog("end2", time.Now())

	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	s.Content = utils.structToJsonString(s.PictureElem)
	return utils.structToJsonString(s)
}

func (u *open_im_sdk.UserRelated) CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	s := open_im_sdk.MsgStruct{}
	var p open_im_sdk.PictureBaseInfo
	_ = json.Unmarshal([]byte(sourcePicture), &p)
	s.PictureElem.SourcePicture = p
	_ = json.Unmarshal([]byte(bigPicture), &p)
	s.PictureElem.BigPicture = p
	_ = json.Unmarshal([]byte(snapshotPicture), &p)
	s.PictureElem.SnapshotPicture = p
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Picture)
	s.Content = utils.structToJsonString(s.PictureElem)
	return utils.structToJsonString(s)
}

func (u *open_im_sdk.UserRelated) SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	s := open_im_sdk.MsgStruct{}
	p := open_im_sdk.OfflinePushInfo{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(2038, err.Error())
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	err = json.Unmarshal([]byte(offlinePushInfo), &p)
	if err != nil {
		callback.OnError(2038, err.Error())
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	go func() {
		var conversationID string
		var options map[string]bool
		isRetry := true
		c := ConversationStruct{
			LatestMsgSendTime: s.CreateTime,
		}
		if receiver == "" && groupID == "" {
			callback.OnError(201, "args err")
			return
		} else if receiver == "" {
			s.SessionType = open_im_sdk.GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
			c.GroupID = groupID
			c.ConversationType = open_im_sdk.GroupChatType
			faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
			if err != nil {
				utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
				callback.OnError(202, err.Error())
				return
			}
			c.ShowName = name
			c.FaceURL = faceUrl
			groupMemberList, err := u.getLocalGroupMemberListByGroupID(groupID)
			if err != nil {
				utils.sdkLog("getLocalGroupMemberListByGroupID err:", err)
				callback.OnError(202, err.Error())
				return
			}
			isExistInGroup := func(target string, groupMemberList []open_im_sdk.groupMemberFullInfo) bool {

				for _, element := range groupMemberList {

					if target == element.UserId {
						return true
					}
				}
				return false

			}(s.SendID, groupMemberList)
			if !isExistInGroup {
				utils.sdkLog("SendGroupMessage err:", "not exist in this group")
				callback.OnError(208, "not exist in this group")
				return
			}

		} else {
			s.SessionType = open_im_sdk.SingleChatType
			s.RecvID = receiver
			conversationID = utils.GetConversationIDBySessionType(receiver, open_im_sdk.SingleChatType)
			c.UserID = receiver
			c.ConversationType = open_im_sdk.SingleChatType
			faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
			if err != nil {
				utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
				callback.OnError(301, err.Error())
				return
			}
			c.FaceURL = faceUrl
			c.ShowName = name

		}
		c.ConversationID = conversationID
		c.LatestMsg = utils.structToJsonString(s)
		if !onlineUserOnly {
			err = u.insertMessageToLocalOrUpdateContent(&s)
			if err != nil {
				utils.sdkLog("insertMessageToLocalOrUpdateContent err:", err)
				callback.OnError(202, err.Error())
				return
			}
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.AddConOrUpLatMsg,
				c}})
			//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
			//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		} else {
			options = make(map[string]bool, 2)
			options[open_im_sdk.IsHistory] = false
			options[open_im_sdk.IsPersistent] = false
			isRetry = false
		}

		var delFile []string
		//media file handle
		switch s.ContentType {
		case open_im_sdk.Picture:
			var sourcePath string
			if utils.fileExist(s.PictureElem.SourcePath) {
				sourcePath = s.PictureElem.SourcePath
				delFile = append(delFile, utils.fileTmpPath(s.PictureElem.SourcePath))
			} else {
				sourcePath = utils.fileTmpPath(s.PictureElem.SourcePath)
				delFile = append(delFile, sourcePath)
			}
			utils.sdkLog("file: ", sourcePath, delFile)
			sourceUrl, uuid, err := u.uploadImage(sourcePath, callback)
			if err != nil {
				utils.sdkLog("oss Picture upload err", err.Error())
				callback.OnError(301, err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.PictureElem.SourcePicture.Url = sourceUrl
				s.PictureElem.SourcePicture.UUID = uuid
				s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + open_im_sdk.ZoomScale + "/h/" + open_im_sdk.ZoomScale
				s.PictureElem.SnapshotPicture.Width = int32(utils.stringToInt(open_im_sdk.ZoomScale))
				s.PictureElem.SnapshotPicture.Height = int32(utils.stringToInt(open_im_sdk.ZoomScale))
				s.Content = utils.structToJsonString(s.PictureElem)
			}
		case open_im_sdk.Voice:
			var sourcePath string
			if utils.fileExist(s.SoundElem.SoundPath) {
				sourcePath = s.SoundElem.SoundPath
				delFile = append(delFile, utils.fileTmpPath(s.SoundElem.SoundPath))
			} else {
				sourcePath = utils.fileTmpPath(s.SoundElem.SoundPath)
				delFile = append(delFile, sourcePath)
			}
			utils.sdkLog("file: ", sourcePath, delFile)
			soundURL, uuid, err := u.uploadSound(sourcePath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				utils.sdkLog("uploadSound err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = utils.structToJsonString(s.SoundElem)
			}
		case open_im_sdk.Video:
			var videoPath string
			var snapPath string
			if utils.fileExist(s.VideoElem.VideoPath) {
				videoPath = s.VideoElem.VideoPath
				snapPath = s.VideoElem.SnapshotPath
				delFile = append(delFile, utils.fileTmpPath(s.VideoElem.VideoPath))
				delFile = append(delFile, utils.fileTmpPath(s.VideoElem.SnapshotPath))
			} else {
				videoPath = utils.fileTmpPath(s.VideoElem.VideoPath)
				snapPath = utils.fileTmpPath(s.VideoElem.SnapshotPath)
				delFile = append(delFile, videoPath)
				delFile = append(delFile, snapPath)
			}
			utils.sdkLog("file: ", videoPath, snapPath, delFile)
			snapshotURL, snapshotUUID, videoURL, videoUUID, err := u.uploadVideo(videoPath, snapPath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				utils.sdkLog("oss  Video upload err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return
			} else {
				s.VideoElem.VideoURL = videoURL
				s.VideoElem.SnapshotUUID = snapshotUUID
				s.VideoElem.SnapshotURL = snapshotURL
				s.VideoElem.VideoUUID = videoUUID
				s.Content = utils.structToJsonString(s.VideoElem)
			}
		case open_im_sdk.File:
			fileURL, fileUUID, err := u.uploadFile(s.FileElem.FilePath, callback)
			if err != nil {
				callback.OnError(301, err.Error())
				utils.sdkLog("oss  File upload err:", err.Error())
				u.sendMessageFailedHandle(&s, &c, conversationID)
				return

			} else {
				s.FileElem.SourceURL = fileURL
				s.FileElem.UUID = fileUUID
				s.Content = utils.structToJsonString(s.FileElem)
			}
		case open_im_sdk.Text:
		case open_im_sdk.AtText:
		case open_im_sdk.Location:
		case open_im_sdk.Custom:
		case open_im_sdk.Merger:
		case open_im_sdk.Quote:
		case open_im_sdk.Card:
		default:
			callback.OnError(2038, "Not currently supported ")
			utils.sdkLog("Not currently supported ", s.ContentType)
			return
		}
		if !onlineUserOnly {
			//Store messages to local database
			err = u.insertMessageToLocalOrUpdateContent(&s)
			if err != nil {
				callback.OnError(202, err.Error())
				return
			}
		}
		sendMessageToServer(&onlineUserOnly, &s, u, callback, &c, conversationID, delFile, &p, isRetry, options)
	}()
	return s.ClientMsgID
}
func (u *open_im_sdk.UserRelated) internalSendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string, options map[string]bool) (err error) {
	s := open_im_sdk.MsgStruct{}
	p := open_im_sdk.OfflinePushInfo{}
	err = json.Unmarshal([]byte(message), &s)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return err
	}
	err = json.Unmarshal([]byte(offlinePushInfo), &p)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return err
	}

	var conversationID string
	isRetry := true
	c := ConversationStruct{
		LatestMsgSendTime: s.CreateTime,
	}
	if receiver == "" && groupID == "" {
		return errors.New("args err")
	} else if receiver == "" {
		s.SessionType = open_im_sdk.GroupChatType
		s.RecvID = groupID
		s.GroupID = groupID
		conversationID = utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
		c.GroupID = groupID
		c.ConversationType = open_im_sdk.GroupChatType
		faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
		if err != nil {
			utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
			return errors.New("getGroupNameAndFaceUrlByUid err")
		}
		c.ShowName = name
		c.FaceURL = faceUrl
		groupMemberList, err := u.getLocalGroupMemberListByGroupID(groupID)
		if err != nil {
			utils.sdkLog("getLocalGroupMemberListByGroupID err:", err)
			return errors.New("getLocalGroupMemberListByGroupID err")
		}
		isExistInGroup := func(target string, groupMemberList []open_im_sdk.groupMemberFullInfo) bool {

			for _, element := range groupMemberList {

				if target == element.UserId {
					return true
				}
			}
			return false

		}(s.SendID, groupMemberList)
		if !isExistInGroup {
			utils.sdkLog("SendGroupMessage err:", "not exist in this group")
			return errors.New("not exist in this group")
		}

	} else {
		s.SessionType = open_im_sdk.SingleChatType
		s.RecvID = receiver
		conversationID = utils.GetConversationIDBySessionType(receiver, open_im_sdk.SingleChatType)
		c.UserID = receiver
		c.ConversationType = open_im_sdk.SingleChatType
		faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
		if err != nil {
			utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
			return errors.New("getUserNameAndFaceUrlByUid err")
		}
		c.FaceURL = faceUrl
		c.ShowName = name

	}
	c.ConversationID = conversationID
	c.LatestMsg = utils.structToJsonString(s)
	if !onlineUserOnly {
		err = u.insertMessageToLocalOrUpdateContent(&s)
		if err != nil {
			utils.sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			return errors.New("insertMessageToLocalOrUpdateContent err:")
		}
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.AddConOrUpLatMsg,
			c}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
	} else {
		options[open_im_sdk.IsHistory] = false
		options[open_im_sdk.IsPersistent] = false
		isRetry = false
	}

	sendMessageToServer(&onlineUserOnly, &s, u, callback, &c, conversationID, []string{}, &p, isRetry, options)
	return nil

}
func (u *open_im_sdk.UserRelated) SendMessageNotOss(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	s := open_im_sdk.MsgStruct{}
	p := open_im_sdk.OfflinePushInfo{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(2038, err.Error())
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	err = json.Unmarshal([]byte(offlinePushInfo), &p)
	if err != nil {
		callback.OnError(2038, err.Error())
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}

	go func() {
		var conversationID string
		var options map[string]bool
		isRetry := true
		c := ConversationStruct{
			LatestMsgSendTime: s.CreateTime,
		}
		if receiver == "" && groupID == "" {
			callback.OnError(201, "args err")
			return
		} else if receiver == "" {
			s.SessionType = open_im_sdk.GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
			c.GroupID = groupID
			c.ConversationType = open_im_sdk.GroupChatType
			faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
			if err != nil {
				utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
				callback.OnError(202, err.Error())
				return
			}
			c.ShowName = name
			c.FaceURL = faceUrl
			groupMemberList, err := u.getLocalGroupMemberListByGroupID(groupID)
			if err != nil {
				utils.sdkLog("getLocalGroupMemberListByGroupID err:", err)
				callback.OnError(202, err.Error())
				return
			}
			isExistInGroup := func(target string, groupMemberList []open_im_sdk.groupMemberFullInfo) bool {

				for _, element := range groupMemberList {

					if target == element.UserId {
						return true
					}
				}
				return false

			}(s.SendID, groupMemberList)
			if !isExistInGroup {
				utils.sdkLog("SendGroupMessage err:", "not exist in this group")
				callback.OnError(208, "not exist in this group")
				return
			}

		} else {
			s.SessionType = open_im_sdk.SingleChatType
			s.RecvID = receiver
			conversationID = utils.GetConversationIDBySessionType(receiver, open_im_sdk.SingleChatType)
			c.UserID = receiver
			c.ConversationType = open_im_sdk.SingleChatType
			faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
			if err != nil {
				utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
				callback.OnError(301, err.Error())
				return
			}
			c.FaceURL = faceUrl
			c.ShowName = name
		}
		c.ConversationID = conversationID
		c.LatestMsg = utils.structToJsonString(s)

		if !onlineUserOnly {
			err = u.insertMessageToLocalOrUpdateContent(&s)
			if err != nil {
				utils.sdkLog("insertMessageToLocalOrUpdateContent err:", err)
				callback.OnError(202, err.Error())
				return
			}
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.AddConOrUpLatMsg,
				c}})
			//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
			//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		} else {
			options = make(map[string]bool, 2)
			options[open_im_sdk.IsHistory] = false
			options[open_im_sdk.IsPersistent] = false
			isRetry = false
		}
		sendMessageToServer(&onlineUserOnly, &s, u, callback, &c, conversationID, []string{}, &p, isRetry, options)

	}()
	return s.ClientMsgID
}
func (u *open_im_sdk.UserRelated) autoSendMsg(s *open_im_sdk.MsgStruct, receiver, groupID string, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo bool, offlinePushInfo string) error {
	utils.sdkLog("autoSendMsg input args:", *s, receiver, groupID, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo)
	var conversationID string
	p := open_im_sdk.OfflinePushInfo{}
	err := json.Unmarshal([]byte(offlinePushInfo), &p)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return err
	}
	r := open_im_sdk.SendMsgRespFromServer{}
	a := open_im_sdk.paramsUserSendMsg{}
	if receiver == "" {
		s.SessionType = open_im_sdk.GroupChatType
		s.RecvID = groupID
	} else if groupID == "" {
		s.SessionType = open_im_sdk.SingleChatType
		s.RecvID = receiver
	} else {
		utils.sdkLog("args err: ", receiver, groupID)
		return errors.New("args null")
	}
	c := ConversationStruct{
		ConversationType:  int(s.SessionType),
		LatestMsgSendTime: s.CreateTime,
	}
	if receiver == "" && groupID == "" {
		return errors.New("args error")
	} else if receiver == "" {
		s.SessionType = open_im_sdk.GroupChatType
		s.RecvID = groupID
		s.GroupID = groupID
		conversationID = utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
		c.GroupID = groupID
		faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
		if err != nil {
			utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
			return err
		}
		c.ShowName = name
		c.FaceURL = faceUrl
	} else {
		s.SessionType = open_im_sdk.SingleChatType
		s.RecvID = receiver
		conversationID = utils.GetConversationIDBySessionType(receiver, open_im_sdk.SingleChatType)
		c.UserID = receiver
		faceUrl, name, err := u.getUserNameAndFaceUrlByUid(receiver)
		if err != nil {
			utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
			return err
		}
		c.FaceURL = faceUrl
		c.ShowName = name
	}
	userInfo, err := u.getLoginUserInfoFromLocal()
	if err != nil {
		utils.sdkLog("getLoginUserInfoFromLocal err:", err)
		return err
	}
	s.SenderFaceURL = userInfo.Icon
	s.SenderNickname = userInfo.Name
	c.ConversationID = conversationID
	c.LatestMsg = utils.structToJsonString(s)
	if !onlineUserOnly {
		err = u.insertMessageToLocalOrUpdateContent(s)
		if err != nil {
			utils.sdkLog("insertMessageToLocalOrUpdateContent err:", err)
			return err
		}
	}
	optionsFlag := make(map[string]bool, 2)
	if onlineUserOnly {
		optionsFlag[open_im_sdk.IsHistory] = false
		optionsFlag[open_im_sdk.IsPersistent] = false
	}

	//Protocol conversion
	a.SenderPlatformID = s.SenderPlatformID
	a.SendID = s.SendID
	a.SenderNickName = s.SenderNickname
	a.SenderFaceURL = s.SenderFaceURL
	a.OperationID = utils.operationIDGenerator()
	a.Data.SessionType = s.SessionType
	a.Data.MsgFrom = s.MsgFrom
	a.Data.ContentType = s.ContentType
	a.Data.RecvID = s.RecvID
	a.Data.GroupID = s.GroupID
	a.Data.ForceList = s.ForceList
	a.Data.Content = []byte(s.Content)
	a.Data.Options = optionsFlag
	a.Data.ClientMsgID = s.ClientMsgID
	a.Data.CreateTime = s.CreateTime
	a.Data.OffLineInfo = p
	bMsg, err := utils.post2Api(open_im_sdk.sendMsgRouter, a, u.token)
	if err != nil {
		utils.sdkLog("sendMsgRouter access err:", err.Error())
		u.updateMessageFailedStatus(s, &c, onlineUserOnly)
		return err
	} else {
		err = json.Unmarshal(bMsg, &r)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			u.updateMessageFailedStatus(s, &c, onlineUserOnly)
			return err
		} else {
			if r.ErrCode != 0 {
				utils.sdkLog("errcode, errmsg: ", r.ErrCode, r.ErrMsg)
				u.updateMessageFailedStatus(s, &c, onlineUserOnly)
				return err
			} else {
				if !onlineUserOnly {
					_ = u.updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.SendTime, open_im_sdk.MsgStatusSendSuccess)
				}
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				s.Status = open_im_sdk.MsgStatusSendSuccess
				c.LatestMsg = utils.structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				if isUpdateConversationLatestMsg {
					u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.AddConOrUpLatMsg, c}})
					u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.IncrUnread, ""}})
				}
				if isUpdateConversationInfo {
					u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.UpdateFaceUrlAndNickName, c}})

				}
				if isUpdateConversationInfo || isUpdateConversationLatestMsg {
					u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
					u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.TotalUnreadMessageChanged, ""}})
				}
			}
		}
	}
	return nil
}
func (u *open_im_sdk.UserRelated) CreateSoundMessageByURL(soundBaseInfo string) string {
	s := open_im_sdk.MsgStruct{}
	var soundElem open_im_sdk.SoundBaseInfo
	_ = json.Unmarshal([]byte(soundBaseInfo), &soundElem)
	s.SoundElem = soundElem
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Voice)
	s.Content = utils.structToJsonString(s.SoundElem)
	return utils.structToJsonString(s)
}

func (u *open_im_sdk.UserRelated) CreateSoundMessage(soundPath string, duration int64) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Voice)
	s.SoundElem.SoundPath = open_im_sdk.SvrConf.DbDir + soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		utils.sdkLog(err.Error())
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	s.Content = utils.structToJsonString(s.SoundElem)
	return utils.structToJsonString(s)
}

func (u *open_im_sdk.UserRelated) CreateVideoMessageByURL(videoBaseInfo string) string {
	s := open_im_sdk.MsgStruct{}
	var videoElem open_im_sdk.VideoBaseInfo
	_ = json.Unmarshal([]byte(videoBaseInfo), &videoElem)
	s.VideoElem = videoElem
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Video)
	s.Content = utils.structToJsonString(s.VideoElem)
	return utils.structToJsonString(s)
}

func (u *open_im_sdk.UserRelated) CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Video)
	s.VideoElem.VideoPath = open_im_sdk.SvrConf.DbDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = open_im_sdk.SvrConf.DbDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		utils.sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := open_im_sdk.getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			utils.sdkLog("CreateVideoMessage err:", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	s.Content = utils.structToJsonString(s.VideoElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateFileMessageByURL(fileBaseInfo string) string {
	s := open_im_sdk.MsgStruct{}
	var fileElem open_im_sdk.FileBaseInfo
	_ = json.Unmarshal([]byte(fileBaseInfo), &fileElem)
	s.FileElem = fileElem
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.File)
	s.Content = utils.structToJsonString(s.FileElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateFileMessage(filePath string, fileName string) string {
	s := open_im_sdk.MsgStruct{}
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.File)
	s.FileElem.FilePath = open_im_sdk.SvrConf.DbDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		utils.sdkLog(err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	utils.sdkLog("CreateForwardMessage new input: ", utils.structToJsonString(s))
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateMergerMessage(messageList, title, summaryList string) string {
	var messages []*open_im_sdk.MsgStruct
	var summaries []string
	s := open_im_sdk.MsgStruct{}
	err := json.Unmarshal([]byte(messageList), &messages)
	if err != nil {
		utils.sdkLog("CreateMergerMessage err:", err.Error())
		return ""
	}
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Merger)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.structToJsonString(s.MergeElem)
	return utils.structToJsonString(s)
}
func (u *open_im_sdk.UserRelated) CreateForwardMessage(m string) string {
	utils.sdkLog("CreateForwardMessage input: ", m)
	s := open_im_sdk.MsgStruct{}
	err := json.Unmarshal([]byte(m), &s)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	if s.Status != open_im_sdk.MsgStatusSendSuccess {
		utils.sdkLog("only send success message can be revoked")
		return ""
	}

	u.initBasicInfo(&s, open_im_sdk.UserMsgType, s.ContentType)
	//Forward message seq is set to 0
	s.Seq = 0
	return utils.structToJsonString(s)
}

func sendMessageToServer(onlineUserOnly *bool, s *open_im_sdk.MsgStruct, u *open_im_sdk.UserRelated, callback SendMsgCallBack,
	c *ConversationStruct, conversationID string, delFile []string, offlinePushInfo *open_im_sdk.OfflinePushInfo, isRetry bool, options map[string]bool) {
	//Protocol conversion
	wsMsgData := open_im_sdk.MsgData{
		SendID:           s.SendID,
		RecvID:           s.RecvID,
		GroupID:          s.GroupID,
		ClientMsgID:      s.ClientMsgID,
		ServerMsgID:      s.ServerMsgID,
		SenderPlatformID: s.SenderPlatformID,
		SenderNickname:   s.SenderNickname,
		SenderFaceURL:    s.SenderFaceURL,
		SessionType:      s.SessionType,
		MsgFrom:          s.MsgFrom,
		ContentType:      s.ContentType,
		Content:          []byte(s.Content),
		ForceList:        s.ForceList,
		CreateTime:       s.CreateTime,
		Options:          options,
		OfflinePushInfo:  offlinePushInfo,
	}
	msgIncr, ch := u.AddCh()
	var wsReq open_im_sdk.GeneralWsReq
	var err error
	wsReq.ReqIdentifier = open_im_sdk.WSSendMsg
	wsReq.OperationID = utils.operationIDGenerator()
	wsReq.SendID = s.SendID
	//wsReq.Token = u.token
	wsReq.MsgIncr = msgIncr
	wsReq.Data, err = proto.Marshal(&wsMsgData)
	if err != nil {
		utils.sdkLog("Marshal failed ", err.Error())
		utils.LogFReturn(nil)
		callback.OnError(http.StatusInternalServerError, err.Error())
		u.sendMessageFailedHandle(s, c, conversationID)
		return
	}

	SendFlag := false
	var connSend *websocket.Conn
	for tr := 0; tr < 30; tr++ {
		utils.LogBegin("WriteMsg", wsReq.OperationID)
		err, connSend = u.WriteMsg(wsReq)
		utils.LogEnd("WriteMsg ", wsReq.OperationID, connSend)
		if err != nil {
			if !isRetry {
				break
			}
			utils.sdkLog("ws writeMsg  err:,", wsReq.OperationID, err.Error(), tr)
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			utils.sdkLog("writeMsg  retry ok", wsReq.OperationID, tr)
			SendFlag = true
			break
		}
	}
	//onlineUserOnly end after send message to ws
	if *onlineUserOnly {
		return
	}
	if SendFlag == false {
		u.DelCh(msgIncr)
		callback.OnError(http.StatusInternalServerError, err.Error())
		u.sendMessageFailedHandle(s, c, conversationID)
		return
	}

	timeout := 300
	breakFlag := 0

	for {
		if breakFlag == 1 {
			utils.sdkLog("break ", wsReq.OperationID)
			break
		}
		select {
		case r := <-ch:
			utils.sdkLog("ws  ch recvMsg success:,", wsReq.OperationID)
			if r.ErrCode != 0 {
				callback.OnError(int32(r.ErrCode), r.ErrMsg)
				u.sendMessageFailedHandle(s, c, conversationID)
			} else {
				callback.OnProgress(100)
				callback.OnSuccess("")
				//remove media cache file
				for _, v := range delFile {
					err := os.Remove(v)
					if err != nil {
						utils.sdkLog("remove failed,", err.Error(), v)
					}
					utils.sdkLog("remove file: ", v)
				}
				var sendMsgResp open_im_sdk.UserSendMsgResp
				err = proto.Unmarshal(r.Data, &sendMsgResp)
				if err != nil {
					utils.sdkLog("Unmarshal failed ", err.Error())
					//	callback.OnError(http.StatusInternalServerError, err.Error())
					//	u.sendMessageFailedHandle(&s, &c, conversationID)
					//	u.DelCh(msgIncr)
				}
				_ = u.updateMessageTimeAndMsgIDStatus(sendMsgResp.ClientMsgID, sendMsgResp.SendTime, open_im_sdk.MsgStatusSendSuccess)

				s.ServerMsgID = sendMsgResp.ServerMsgID
				s.SendTime = sendMsgResp.SendTime
				s.Status = open_im_sdk.MsgStatusSendSuccess
				c.LatestMsg = utils.structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				_ = u.triggerCmdUpdateConversation(open_im_sdk.updateConNode{conversationID, open_im_sdk.AddConOrUpLatMsg,
					c})
				u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
			}
			breakFlag = 1
		case <-time.After(time.Second * time.Duration(timeout)):
			var flag bool
			utils.sdkLog("ws ch recvMsg err: ", wsReq.OperationID)
			if connSend != u.conn {
				utils.sdkLog("old conn != current conn  ", connSend, u.conn)
				flag = false // error
			} else {
				flag = false //error
				for tr := 0; tr < 3; tr++ {
					err = u.sendPingMsg()
					if err != nil {
						utils.sdkLog("sendPingMsg failed ", wsReq.OperationID, err.Error(), tr)
						time.Sleep(time.Duration(30) * time.Second)
					} else {
						utils.sdkLog("sendPingMsg ok ", wsReq.OperationID)
						flag = true //wait continue
						break
					}
				}
			}
			if flag == false {
				callback.OnError(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))
				u.sendMessageFailedHandle(s, c, conversationID)
				utils.sdkLog("onError callback ", wsReq.OperationID)
				breakFlag = 1
				break
			} else {
				utils.sdkLog("wait resp continue", wsReq.OperationID)
				breakFlag = 0
				continue
			}
		}
	}

	u.DelCh(msgIncr)
}

func (u *open_im_sdk.UserRelated) GetHistoryMessageList(callback open_im_sdk.Base, getMessageOptions string) {
	go func() {
		utils.sdkLog("GetHistoryMessageList", getMessageOptions)
		var sourceID string
		var conversationID string
		var startTime int64
		var latestMsg open_im_sdk.MsgStruct
		var sessionType int
		p := open_im_sdk.PullMsgReq{}
		err := json.Unmarshal([]byte(getMessageOptions), &p)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		if p.UserID == "" {
			sourceID = p.GroupID
			conversationID = utils.GetConversationIDBySessionType(sourceID, open_im_sdk.GroupChatType)
			sessionType = open_im_sdk.GroupChatType
		} else {
			sourceID = p.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, open_im_sdk.SingleChatType)
			sessionType = open_im_sdk.SingleChatType
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
					utils.sdkLog("get history err :", err)
					callback.OnError(200, err.Error())
					return
				}
				startTime = latestMsg.SendTime + open_im_sdk.TimeOffset
			}

		} else {
			startTime = p.StartMsg.SendTime
		}
		utils.sdkLog("sourceID:", sourceID, "startTime:", startTime, "count:", p.Count)
		err, list := u.getHistoryMessage(sourceID, startTime, p.Count, sessionType)
		sort.Sort(list)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(utils.structToJsonString(list))
			} else {
				callback.OnSuccess(utils.structToJsonString([]open_im_sdk.MsgStruct{}))
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) RevokeMessage(callback open_im_sdk.Base, message string) {
	go func() {
		//var receiver, groupID string
		c := open_im_sdk.MsgStruct{}
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
		if s.Status != open_im_sdk.MsgStatusSendSuccess {
			callback.OnError(201, "only send success message can be revoked")
			return
		}
		utils.sdkLog("test data", s)
		//Send message internally
		switch s.SessionType {
		case open_im_sdk.SingleChatType:
			//receiver = s.RecvID
		case open_im_sdk.GroupChatType:
			//groupID = s.GroupID
		default:
			callback.OnError(200, "args err")
		}
		s.Content = s.ClientMsgID
		s.ClientMsgID = utils.getMsgID(s.SendID)
		s.ContentType = open_im_sdk.Revoke
		//err = u.autoSendMsg(s, receiver, groupID, false, true, false)
		if err != nil {
			utils.sdkLog("autoSendMsg revokeMessage err:", err.Error())
			callback.OnError(300, err.Error())

		} else {
			err = u.setMessageStatus(s.Content, open_im_sdk.MsgStatusRevoked)
			if err != nil {
				utils.sdkLog("setLocalMessageStatus revokeMessage err:", err.Error())
				callback.OnError(300, err.Error())
			} else {
				callback.OnSuccess("")
			}
		}
	}()
}
func (u *open_im_sdk.UserRelated) TypingStatusUpdate(receiver, msgTip string) {
	go func() {
		s := open_im_sdk.MsgStruct{}
		u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.Typing)
		s.Content = msgTip
		//err := u.autoSendMsg(&s, receiver, "", true, false, false)
		//if err != nil {
		//	sdkLog("TypingStatusUpdate err:", err)
		//} else {
		//	sdkLog("TypingStatusUpdate success!!!")
		//}
	}()
}

func (u *open_im_sdk.UserRelated) MarkC2CMessageAsRead(callback open_im_sdk.Base, receiver string, msgIDList string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(receiver, open_im_sdk.SingleChatType)
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
		s := open_im_sdk.MsgStruct{}
		u.initBasicInfo(&s, open_im_sdk.UserMsgType, open_im_sdk.HasReadReceipt)
		s.Content = msgIDList
		utils.sdkLog("MarkC2CMessageAsRead: send Message")
		//err = u.autoSendMsg(&s, receiver, "", false, false, false)
		if err != nil {
			utils.sdkLog("MarkC2CMessageAsRead  err:", err.Error())
			callback.OnError(300, err.Error())
		} else {
			callback.OnSuccess("")
			err = u.setSingleMessageHasReadByMsgIDList(receiver, list)
			if err != nil {
				utils.sdkLog("setSingleMessageHasReadByMsgIDList  err:", err.Error())
			}
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{conversationID, open_im_sdk.UpdateLatestMessageChange, ""}})
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}()
}

//Deprecated
func (u *open_im_sdk.UserRelated) MarkSingleMessageHasRead(callback open_im_sdk.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, open_im_sdk.SingleChatType)
		//if err := u.setSingleMessageHasRead(userID); err != nil {
		//	callback.OnError(201, err.Error())
		//} else {
		callback.OnSuccess("")
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{ConId: conversationID, Action: open_im_sdk.UnreadCountSetZero}})
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		//}
	}()
}
func (u *open_im_sdk.UserRelated) MarkAllConversationHasRead(callback open_im_sdk.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, open_im_sdk.SingleChatType)
		//if err := u.setSingleMessageHasRead(userID); err != nil {
		//	callback.OnError(201, err.Error())
		//} else {
		callback.OnSuccess("")
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{ConId: conversationID, Action: open_im_sdk.UnreadCountSetZero}})
		u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		//}
	}()
}
func (u *open_im_sdk.UserRelated) MarkGroupMessageHasRead(callback open_im_sdk.Base, groupID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
		if err := u.setGroupMessageHasRead(groupID); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{ConId: conversationID, Action: open_im_sdk.UnreadCountSetZero}})
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}()
}
func (u *open_im_sdk.UserRelated) DeleteMessageFromLocalStorage(callback open_im_sdk.Base, message string) {
	go func() {
		var conversation ConversationStruct
		var latestMsg open_im_sdk.MsgStruct
		var conversationID string
		var sourceID string
		s := open_im_sdk.MsgStruct{}
		err := json.Unmarshal([]byte(message), &s)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		err = u.setMessageStatus(s.ClientMsgID, open_im_sdk.MsgStatusHasDeleted)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		callback.OnSuccess("")
		if s.SessionType == open_im_sdk.GroupChatType {
			conversationID = utils.GetConversationIDBySessionType(s.RecvID, open_im_sdk.GroupChatType)
			sourceID = s.RecvID

		} else if s.SessionType == open_im_sdk.SingleChatType {
			if s.SendID != u.loginUserID {
				conversationID = utils.GetConversationIDBySessionType(s.SendID, open_im_sdk.SingleChatType)
				sourceID = s.SendID
			} else {
				conversationID = utils.GetConversationIDBySessionType(s.RecvID, open_im_sdk.SingleChatType)
				sourceID = s.RecvID
			}
		}
		_, m := u.getConversationLatestMsgModel(conversationID)
		if m != "" {
			err := json.Unmarshal([]byte(m), &latestMsg)
			if err != nil {
				utils.sdkLog("DeleteMessage err :", err)
				callback.OnError(200, err.Error())
				return
			}
		} else {
			utils.sdkLog("err ,conversation has been deleted")
		}

		if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
			err, list := u.getHistoryMessage(sourceID, s.SendTime+open_im_sdk.TimeOffset, 1, int(s.SessionType))
			if err != nil {
				utils.sdkLog("DeleteMessageFromLocalStorage database err:", err.Error())
			}
			conversation.ConversationID = conversationID
			if list == nil {
				conversation.LatestMsg = ""
				conversation.LatestMsgSendTime = utils.getCurrentTimestampByNano()
			} else {
				conversation.LatestMsg = utils.structToJsonString(list[0])
				conversation.LatestMsgSendTime = list[0].SendTime
			}
			err = u.triggerCmdUpdateConversation(open_im_sdk.updateConNode{ConId: conversationID, Action: open_im_sdk.AddConOrUpLatMsg, Args: conversation})
			if err != nil {
				utils.sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
			}
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})

		}
	}()
}
func (u *open_im_sdk.UserRelated) ClearC2CHistoryMessage(callback open_im_sdk.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, open_im_sdk.SingleChatType)
		err := u.setMessageStatusBySourceID(userID, open_im_sdk.MsgStatusHasDeleted, open_im_sdk.SingleChatType)
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
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}()
}
func (u *open_im_sdk.UserRelated) ClearGroupHistoryMessage(callback open_im_sdk.Base, groupID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(groupID, open_im_sdk.GroupChatType)
		err := u.setMessageStatusBySourceID(groupID, open_im_sdk.MsgStatusHasDeleted, open_im_sdk.GroupChatType)
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
			u.doUpdateConversation(open_im_sdk.cmd2Value{Value: open_im_sdk.updateConNode{"", open_im_sdk.NewConChange, []string{conversationID}}})
		}
	}()
}

func (u *open_im_sdk.UserRelated) InsertSingleMessageToLocalStorage(callback open_im_sdk.Base, message, userID, sender string) string {
	s := open_im_sdk.MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return ""
	}
	s.SendID = sender
	s.RecvID = userID
	//Generate client message primary key
	s.ClientMsgID = utils.getMsgID(s.SendID)
	s.SendTime = utils.getCurrentTimestampByNano()
	go func() {
		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
		}
	}()
	return s.ClientMsgID
}

func (u *open_im_sdk.UserRelated) InsertGroupMessageToLocalStorage(callback open_im_sdk.Base, message, groupID, sender string) string {
	s := open_im_sdk.MsgStruct{}
	err := json.Unmarshal([]byte(message), &s)
	if err != nil {
		callback.OnError(200, err.Error())
		return ""
	}
	s.SendID = sender
	s.RecvID = groupID
	//Generate client message primary key
	s.ClientMsgID = utils.getMsgID(s.SendID)
	s.SendTime = utils.getCurrentTimestampByNano()
	go func() {
		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
		}
	}()
	return s.ClientMsgID
}

func (u *open_im_sdk.UserRelated) FindMessages(callback open_im_sdk.Base, messageIDList string) {
	go func() {
		var c []string
		err := json.Unmarshal([]byte(messageIDList), &c)
		if err != nil {
			callback.OnError(200, err.Error())
			utils.sdkLog("Unmarshal failed, ", err.Error())

		}
		err, list := u.getMultipleMessageModel(c)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			if list != nil {
				callback.OnSuccess(utils.structToJsonString(list))
			} else {
				callback.OnSuccess(utils.structToJsonString([]open_im_sdk.MsgStruct{}))
			}
		}
	}()
}
