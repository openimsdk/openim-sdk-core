package conversation_msg

import (
	//"bytes"
	//"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	imgtype "github.com/shamsher31/goimgtype"
	"image"
	"net/http"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"os"
	"sort"
	"sync"
	"time"
)

func (c *Conversation) GetAllConversationList(callback common.Base) {
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
func (c *Conversation) GetConversationListSplit(callback common.Base, offset, count int) {
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
func (c *Conversation) SetConversationRecvMessageOpt(callback common.Base, conversationIDList string, opt int) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := utils.post2Api(common.setReceiveMessageOptRouter, common.paramsSetReceiveMessageOpt{OperationID: utils.operationIDGenerator(), Option: int32(opt), ConversationIdList: list}, u.token)
		if err != nil {
			utils.sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g common.getReceiveMessageOptResp
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
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, list}})

	}()
}
func (c *Conversation) GetConversationRecvMessageOpt(callback common.Base, conversationIDList string) {
	go func() {
		var list []string
		err := json.Unmarshal([]byte(conversationIDList), &list)
		if err != nil {
			utils.sdkLog("unmarshal failed, ", err.Error())
			callback.OnError(201, err.Error())
			return
		}
		resp, err := utils.post2Api(common.getReceiveMessageOptRouter, common.paramGetReceiveMessageOpt{OperationID: utils.operationIDGenerator(), ConversationIdList: list}, u.token)
		if err != nil {
			utils.sdkLog("post failed, ", err.Error())
			callback.OnError(202, err.Error())
			return
		}
		var g common.getReceiveMessageOptResp
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
func (c *Conversation) GetOneConversation(sourceID string, sessionType int, callback common.Base) {
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
				case constant.SingleChatType:
					c.UserID = sourceID
					faceUrl, name, err := u.getUserNameAndFaceUrlByUid(sourceID)
					if err != nil {
						callback.OnError(301, err.Error())
						utils.sdkLog("getUserNameAndFaceUrlByUid err:", err)
						return
					}
					c.ShowName = name
					c.FaceURL = faceUrl
				case constant.GroupChatType:
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
func (c *Conversation) GetMultipleConversation(conversationIDList string, callback common.Base) {
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
func (c *Conversation) DeleteConversation(conversationID string, callback common.Base) {
	go func() {
		//Transaction operation required
		var sourceID string
		err, c := u.getOneConversationModel(conversationID)
		if err != nil {
			callback.OnError(201, err.Error())
			return
		}
		switch c.ConversationType {
		case constant.SingleChatType:
			sourceID = c.UserID
		case constant.GroupChatType:
			sourceID = c.GroupID
		}
		//Mark messages related to this conversation for deletion
		err = u.setMessageStatusBySourceID(sourceID, constant.MsgStatusHasDeleted, c.ConversationType)
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
			_ = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.TotalUnreadMessageChanged})
		}
	}()
}
func (c *Conversation) SetConversationDraft(conversationID, draftText string, callback common.Base) {
	if draftText != "" {
		err := u.setConversationDraftModel(conversationID, draftText)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	} else {
		err := u.removeConversationDraftModel(conversationID, draftText)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			//_ = u.triggerCmdUpdateConversation(updateConNode{ConId: conversationID, Action: ConAndUnreadChange})
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}
}
func (c *Conversation) PinConversation(conversationID string, isPinned bool, callback common.Base) {
	if isPinned {
		err := u.pinConversationModel(conversationID, constant.Pinned)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	} else {
		err := u.unPinConversationModel(conversationID, constant.NotPinned)
		if err != nil {
			callback.OnError(203, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}

}
func (c *Conversation) GetTotalUnreadMsgCount(callback common.Base) {
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
	OnTotalUnreadMessageCountChanged(totalUnreadCount int64)
}

func (c *Conversation) SetConversationListener(listener OnConversationListener) {
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

func (c *Conversation) AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
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

func (c *Conversation) ForceSyncMsg() bool {
	if u.syncSeq2Msg() == nil {
		return true
	} else {
		return false
	}
}

func (c *Conversation) ForceSyncJoinedGroup() {
	u.syncJoinedGroupInfo()
}

func (c *Conversation) ForceSyncJoinedGroupMember() {

	u.syncJoinedGroupMember()
}

func (c *Conversation) ForceSyncGroupRequest() {
	u.syncGroupRequest()
}

func (c *Conversation) ForceSyncSelfGroupRequest() {
	u.syncSelfGroupRequest()
}

type SendMsgCallBack interface {
	common.Base
	OnProgress(progress int)
}

func (c *Conversation) CreateTextMessage(text string) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Text)
	s.Content = text
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.AtText)
	s.ForceList = users
	s.AtElem.Text = text
	s.AtElem.AtUserList = users
	s.Content = utils.structToJsonString(s.AtElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateLocationMessage(description string, longitude, latitude float64) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Location)
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = utils.structToJsonString(s.LocationElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateCustomMessage(data, extension string, description string) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Custom)
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = utils.structToJsonString(s.CustomElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateQuoteMessage(text string, message string) string {
	s, qs := utils.MsgStruct{}, utils.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	u.initBasicInfo(&s, constant.UserMsgType, constant.Quote)
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.QuoteMessage = &qs
	s.Content = utils.structToJsonString(s.QuoteElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateCardMessage(cardInfo string) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Card)
	s.Content = cardInfo
	return utils.structToJsonString(s)

}
func (c *Conversation) CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
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

	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Video)
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
		imageInfo, err := common.getImageInfo(s.VideoElem.SnapshotPath)
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
func (c *Conversation) CreateFileMessageFromFullPath(fileFullPath string, fileName string) string {
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
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.File)
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
func (c *Conversation) CreateImageMessageFromFullPath(imageFullPath string) string {
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

	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.PictureElem.SourcePath = imageFullPath
	utils.sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := common.getImageInfo(s.PictureElem.SourcePath)
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
func (c *Conversation) CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
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
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
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
func (c *Conversation) CreateImageMessage(imagePath string) string {
	utils.sdkLog("start1: ", time.Now())
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.PictureElem.SourcePath = constant.SvrConf.DbDir + imagePath
	utils.sdkLog("ImageMessage  path:", s.PictureElem.SourcePath)
	utils.sdkLog("end1", time.Now())

	utils.sdkLog("start2 ", time.Now())
	imageInfo, err := common.getImageInfo(s.PictureElem.SourcePath)
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

func (c *Conversation) CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	s := utils.MsgStruct{}
	var p utils.PictureBaseInfo
	_ = json.Unmarshal([]byte(sourcePicture), &p)
	s.PictureElem.SourcePicture = p
	_ = json.Unmarshal([]byte(bigPicture), &p)
	s.PictureElem.BigPicture = p
	_ = json.Unmarshal([]byte(snapshotPicture), &p)
	s.PictureElem.SnapshotPicture = p
	u.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.Content = utils.structToJsonString(s.PictureElem)
	return utils.structToJsonString(s)
}

func (c *Conversation) SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	s := utils.MsgStruct{}
	p := server_api_params.OfflinePushInfo{}
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
			s.SessionType = constant.GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			c.GroupID = groupID
			c.ConversationType = constant.GroupChatType
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
			isExistInGroup := func(target string, groupMemberList []common.groupMemberFullInfo) bool {

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
			s.SessionType = constant.SingleChatType
			s.RecvID = receiver
			conversationID = utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
			c.UserID = receiver
			c.ConversationType = constant.SingleChatType
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
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
				c}})
			//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
			//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		} else {
			options = make(map[string]bool, 2)
			options[constant.IsHistory] = false
			options[constant.IsPersistent] = false
			isRetry = false
		}

		var delFile []string
		//media file handle
		switch s.ContentType {
		case constant.Picture:
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
				s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + constant.ZoomScale + "/h/" + constant.ZoomScale
				s.PictureElem.SnapshotPicture.Width = int32(utils.stringToInt(constant.ZoomScale))
				s.PictureElem.SnapshotPicture.Height = int32(utils.stringToInt(constant.ZoomScale))
				s.Content = utils.structToJsonString(s.PictureElem)
			}
		case constant.Voice:
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
		case constant.Video:
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
		case constant.File:
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
		case constant.Text:
		case constant.AtText:
		case constant.Location:
		case constant.Custom:
		case constant.Merger:
		case constant.Quote:
		case constant.Card:
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
func (c *Conversation) internalSendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string, options map[string]bool) (err error) {
	s := utils.MsgStruct{}
	p := server_api_params.OfflinePushInfo{}
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
		s.SessionType = constant.GroupChatType
		s.RecvID = groupID
		s.GroupID = groupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
		c.GroupID = groupID
		c.ConversationType = constant.GroupChatType
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
		isExistInGroup := func(target string, groupMemberList []common.groupMemberFullInfo) bool {

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
		s.SessionType = constant.SingleChatType
		s.RecvID = receiver
		conversationID = utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
		c.UserID = receiver
		c.ConversationType = constant.SingleChatType
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
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
			c}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
	} else {
		options[constant.IsHistory] = false
		options[constant.IsPersistent] = false
		isRetry = false
	}

	sendMessageToServer(&onlineUserOnly, &s, u, callback, &c, conversationID, []string{}, &p, isRetry, options)
	return nil

}
func (c *Conversation) SendMessageNotOss(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool, offlinePushInfo string) string {
	s := utils.MsgStruct{}
	p := server_api_params.OfflinePushInfo{}
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
			s.SessionType = constant.GroupChatType
			s.RecvID = groupID
			s.GroupID = groupID
			conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			c.GroupID = groupID
			c.ConversationType = constant.GroupChatType
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
			isExistInGroup := func(target string, groupMemberList []common.groupMemberFullInfo) bool {

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
			s.SessionType = constant.SingleChatType
			s.RecvID = receiver
			conversationID = utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
			c.UserID = receiver
			c.ConversationType = constant.SingleChatType
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
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
				c}})
			//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
			//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		} else {
			options = make(map[string]bool, 2)
			options[constant.IsHistory] = false
			options[constant.IsPersistent] = false
			isRetry = false
		}
		sendMessageToServer(&onlineUserOnly, &s, u, callback, &c, conversationID, []string{}, &p, isRetry, options)

	}()
	return s.ClientMsgID
}
func (c *Conversation) autoSendMsg(s *common.MsgStruct, receiver, groupID string, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo bool, offlinePushInfo string) error {
	utils.sdkLog("autoSendMsg input args:", *s, receiver, groupID, onlineUserOnly, isUpdateConversationLatestMsg, isUpdateConversationInfo)
	var conversationID string
	p := server_api_params.OfflinePushInfo{}
	err := json.Unmarshal([]byte(offlinePushInfo), &p)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return err
	}
	r := utils.SendMsgRespFromServer{}
	a := common.paramsUserSendMsg{}
	if receiver == "" {
		s.SessionType = constant.GroupChatType
		s.RecvID = groupID
	} else if groupID == "" {
		s.SessionType = constant.SingleChatType
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
		s.SessionType = constant.GroupChatType
		s.RecvID = groupID
		s.GroupID = groupID
		conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
		c.GroupID = groupID
		faceUrl, name, err := u.getGroupNameAndFaceUrlByUid(groupID)
		if err != nil {
			utils.sdkLog("getGroupNameAndFaceUrlByUid err:", err)
			return err
		}
		c.ShowName = name
		c.FaceURL = faceUrl
	} else {
		s.SessionType = constant.SingleChatType
		s.RecvID = receiver
		conversationID = utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
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
		optionsFlag[constant.IsHistory] = false
		optionsFlag[constant.IsPersistent] = false
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
	bMsg, err := utils.post2Api(common.sendMsgRouter, a, u.token)
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
					_ = u.updateMessageTimeAndMsgIDStatus(r.Data.ClientMsgID, r.Data.SendTime, constant.MsgStatusSendSuccess)
				}
				s.ServerMsgID = r.Data.ServerMsgID
				s.SendTime = r.Data.SendTime
				s.Status = constant.MsgStatusSendSuccess
				c.LatestMsg = utils.structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				if isUpdateConversationLatestMsg {
					u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg, c}})
					u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.IncrUnread, ""}})
				}
				if isUpdateConversationInfo {
					u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.UpdateFaceUrlAndNickName, c}})

				}
				if isUpdateConversationInfo || isUpdateConversationLatestMsg {
					u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
					u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.TotalUnreadMessageChanged, ""}})
				}
			}
		}
	}
	return nil
}
func (c *Conversation) CreateSoundMessageByURL(soundBaseInfo string) string {
	s := utils.MsgStruct{}
	var soundElem utils.SoundBaseInfo
	_ = json.Unmarshal([]byte(soundBaseInfo), &soundElem)
	s.SoundElem = soundElem
	u.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
	s.Content = utils.structToJsonString(s.SoundElem)
	return utils.structToJsonString(s)
}

func (c *Conversation) CreateSoundMessage(soundPath string, duration int64) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
	s.SoundElem.SoundPath = constant.SvrConf.DbDir + soundPath
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

func (c *Conversation) CreateVideoMessageByURL(videoBaseInfo string) string {
	s := utils.MsgStruct{}
	var videoElem utils.VideoBaseInfo
	_ = json.Unmarshal([]byte(videoBaseInfo), &videoElem)
	s.VideoElem = videoElem
	u.initBasicInfo(&s, constant.UserMsgType, constant.Video)
	s.Content = utils.structToJsonString(s.VideoElem)
	return utils.structToJsonString(s)
}

func (c *Conversation) CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.Video)
	s.VideoElem.VideoPath = constant.SvrConf.DbDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = constant.SvrConf.DbDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		utils.sdkLog(err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := common.getImageInfo(s.VideoElem.SnapshotPath)
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
func (c *Conversation) CreateFileMessageByURL(fileBaseInfo string) string {
	s := utils.MsgStruct{}
	var fileElem utils.FileBaseInfo
	_ = json.Unmarshal([]byte(fileBaseInfo), &fileElem)
	s.FileElem = fileElem
	u.initBasicInfo(&s, constant.UserMsgType, constant.File)
	s.Content = utils.structToJsonString(s.FileElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateFileMessage(filePath string, fileName string) string {
	s := utils.MsgStruct{}
	u.initBasicInfo(&s, constant.UserMsgType, constant.File)
	s.FileElem.FilePath = constant.SvrConf.DbDir + filePath
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
func (c *Conversation) CreateMergerMessage(messageList, title, summaryList string) string {
	var messages []*utils.MsgStruct
	var summaries []string
	s := utils.MsgStruct{}
	err := json.Unmarshal([]byte(messageList), &messages)
	if err != nil {
		utils.sdkLog("CreateMergerMessage err:", err.Error())
		return ""
	}
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	u.initBasicInfo(&s, constant.UserMsgType, constant.Merger)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.structToJsonString(s.MergeElem)
	return utils.structToJsonString(s)
}
func (c *Conversation) CreateForwardMessage(m string) string {
	utils.sdkLog("CreateForwardMessage input: ", m)
	s := utils.MsgStruct{}
	err := json.Unmarshal([]byte(m), &s)
	if err != nil {
		utils.sdkLog("json unmarshal err:", err.Error())
		return ""
	}
	if s.Status != constant.MsgStatusSendSuccess {
		utils.sdkLog("only send success message can be revoked")
		return ""
	}

	u.initBasicInfo(&s, constant.UserMsgType, s.ContentType)
	//Forward message seq is set to 0
	s.Seq = 0
	return utils.structToJsonString(s)
}

func sendMessageToServer(onlineUserOnly *bool, s *common.MsgStruct, c *Conversation, callback SendMsgCallBack,
	c *ConversationStruct, conversationID string, delFile []string, offlinePushInfo *server_api_params.OfflinePushInfo, isRetry bool, options map[string]bool) {
	//Protocol conversion
	wsMsgData := server_api_params.MsgData{
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
	var wsReq utils.GeneralWsReq
	var err error
	wsReq.ReqIdentifier = constant.WSSendMsg
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
				var sendMsgResp server_api_params.UserSendMsgResp
				err = proto.Unmarshal(r.Data, &sendMsgResp)
				if err != nil {
					utils.sdkLog("Unmarshal failed ", err.Error())
					//	callback.OnError(http.StatusInternalServerError, err.Error())
					//	u.sendMessageFailedHandle(&s, &c, conversationID)
					//	u.DelCh(msgIncr)
				}
				_ = u.updateMessageTimeAndMsgIDStatus(sendMsgResp.ClientMsgID, sendMsgResp.SendTime, constant.MsgStatusSendSuccess)

				s.ServerMsgID = sendMsgResp.ServerMsgID
				s.SendTime = sendMsgResp.SendTime
				s.Status = constant.MsgStatusSendSuccess
				c.LatestMsg = utils.structToJsonString(s)
				c.LatestMsgSendTime = s.SendTime
				_ = u.triggerCmdUpdateConversation(common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
					c})
				u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
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

func (c *Conversation) GetHistoryMessageList(callback common.Base, getMessageOptions string) {
	go func() {
		utils.sdkLog("GetHistoryMessageList", getMessageOptions)
		var sourceID string
		var conversationID string
		var startTime int64
		var latestMsg utils.MsgStruct
		var sessionType int
		p := utils.PullMsgReq{}
		err := json.Unmarshal([]byte(getMessageOptions), &p)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		if p.UserID == "" {
			sourceID = p.GroupID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.GroupChatType)
			sessionType = constant.GroupChatType
		} else {
			sourceID = p.UserID
			conversationID = utils.GetConversationIDBySessionType(sourceID, constant.SingleChatType)
			sessionType = constant.SingleChatType
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
				startTime = latestMsg.SendTime + TimeOffset
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
				callback.OnSuccess(utils.structToJsonString([]utils.MsgStruct{}))
			}
		}
	}()
}
func (c *Conversation) RevokeMessage(callback common.Base, message string) {
	go func() {
		//var receiver, groupID string
		c := utils.MsgStruct{}
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
		if s.Status != constant.MsgStatusSendSuccess {
			callback.OnError(201, "only send success message can be revoked")
			return
		}
		utils.sdkLog("test data", s)
		//Send message internally
		switch s.SessionType {
		case constant.SingleChatType:
			//receiver = s.RecvID
		case constant.GroupChatType:
			//groupID = s.GroupID
		default:
			callback.OnError(200, "args err")
		}
		s.Content = s.ClientMsgID
		s.ClientMsgID = utils.getMsgID(s.SendID)
		s.ContentType = constant.Revoke
		//err = u.autoSendMsg(s, receiver, groupID, false, true, false)
		if err != nil {
			utils.sdkLog("autoSendMsg revokeMessage err:", err.Error())
			callback.OnError(300, err.Error())

		} else {
			err = u.setMessageStatus(s.Content, constant.MsgStatusRevoked)
			if err != nil {
				utils.sdkLog("setLocalMessageStatus revokeMessage err:", err.Error())
				callback.OnError(300, err.Error())
			} else {
				callback.OnSuccess("")
			}
		}
	}()
}
func (c *Conversation) TypingStatusUpdate(receiver, msgTip string) {
	go func() {
		s := utils.MsgStruct{}
		u.initBasicInfo(&s, constant.UserMsgType, constant.Typing)
		s.Content = msgTip
		//err := u.autoSendMsg(&s, receiver, "", true, false, false)
		//if err != nil {
		//	sdkLog("TypingStatusUpdate err:", err)
		//} else {
		//	sdkLog("TypingStatusUpdate success!!!")
		//}
	}()
}

func (c *Conversation) MarkC2CMessageAsRead(callback common.Base, receiver string, msgIDList string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
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
		s := utils.MsgStruct{}
		u.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt)
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
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.UpdateLatestMessageChange, ""}})
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}()
}

//Deprecated
func (c *Conversation) MarkSingleMessageHasRead(callback common.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
		//if err := u.setSingleMessageHasRead(userID); err != nil {
		//	callback.OnError(201, err.Error())
		//} else {
		callback.OnSuccess("")
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		//}
	}()
}
func (c *Conversation) MarkAllConversationHasRead(callback common.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
		//if err := u.setSingleMessageHasRead(userID); err != nil {
		//	callback.OnError(201, err.Error())
		//} else {
		callback.OnSuccess("")
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		//}
	}()
}
func (c *Conversation) MarkGroupMessageHasRead(callback common.Base, groupID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
		if err := u.setGroupMessageHasRead(groupID); err != nil {
			callback.OnError(201, err.Error())
		} else {
			callback.OnSuccess("")
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}()
}
func (c *Conversation) DeleteMessageFromLocalStorage(callback common.Base, message string) {
	go func() {
		var conversation ConversationStruct
		var latestMsg utils.MsgStruct
		var conversationID string
		var sourceID string
		s := utils.MsgStruct{}
		err := json.Unmarshal([]byte(message), &s)
		if err != nil {
			callback.OnError(200, err.Error())
			return
		}
		err = u.setMessageStatus(s.ClientMsgID, constant.MsgStatusHasDeleted)
		if err != nil {
			callback.OnError(202, err.Error())
			return
		}
		callback.OnSuccess("")
		if s.SessionType == constant.GroupChatType {
			conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.GroupChatType)
			sourceID = s.RecvID

		} else if s.SessionType == constant.SingleChatType {
			if s.SendID != u.loginUserID {
				conversationID = utils.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
				sourceID = s.SendID
			} else {
				conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
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
			err, list := u.getHistoryMessage(sourceID, s.SendTime+TimeOffset, 1, int(s.SessionType))
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
			err = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.AddConOrUpLatMsg, Args: conversation})
			if err != nil {
				utils.sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
			}
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})

		}
	}()
}
func (c *Conversation) ClearC2CHistoryMessage(callback common.Base, userID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
		err := u.setMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
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
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}()
}
func (c *Conversation) ClearGroupHistoryMessage(callback common.Base, groupID string) {
	go func() {
		conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
		err := u.setMessageStatusBySourceID(groupID, constant.MsgStatusHasDeleted, constant.GroupChatType)
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
			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		}
	}()
}

func (c *Conversation) InsertSingleMessageToLocalStorage(callback common.Base, message, userID, sender string) string {
	s := utils.MsgStruct{}
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

func (c *Conversation) InsertGroupMessageToLocalStorage(callback common.Base, message, groupID, sender string) string {
	s := utils.MsgStruct{}
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

func (c *Conversation) FindMessages(callback common.Base, messageIDList string) {
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
				callback.OnSuccess(utils.structToJsonString([]utils.MsgStruct{}))
			}
		}
	}()
}
func getImageInfo(filePath string) (*utils.ImageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, utils.Wrap(err, "open file err")
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, utils.Wrap(err, "image file  Decode err")
	}

	datatype, err := imgtype.Get(filePath)
	if err != nil {
		return nil, utils.Wrap(err, "image file  get type err")
	}
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, utils.Wrap(err, "image file  Stat err")
	}

	b := img.Bounds()

	return &utils.ImageInfo{int32(b.Max.X), int32(b.Max.Y), datatype, fi.Size()}, nil

}
