package conversation_msg

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"image"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"os"
	"runtime"
	"sort"
	"sync"
)

//
//import "C"
import (
	//	//"bytes"
	//	//"encoding/gob"
	//	"encoding/json"
	//	"errors"
	//	"github.com/golang/protobuf/proto"
	//	"github.com/gorilla/websocket"
	imgtype "github.com/shamsher31/goimgtype"
	//	"image"
	//	"net/http"
	//	"open_im_sdk/pkg/db"
	//
	//	"open_im_sdk/pkg/common"
	//	"open_im_sdk/pkg/constant"
	//	"open_im_sdk/pkg/log"
	//	"open_im_sdk/pkg/sdk_params_callback"
	//	"open_im_sdk/pkg/server_api_params"
	//	"open_im_sdk/pkg/utils"
	//	"os"
	//	"sort"
	//	"sync"
	//	"time"
)

func (c *Conversation) GetAllConversationList(callback common.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetAllConversationList args: ")
		result := c.getAllConversationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetAllConversationList callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetConversationListSplit(callback common.Base, offset, count int, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetConversationListSplit args: ", offset, count)
		result := c.getConversationListSplit(callback, offset, count, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetConversationListSplit callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) SetConversationRecvMessageOpt(callback common.Base, conversationIDList string, opt int, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetConversationRecvMessageOpt args: ", conversationIDList, opt)
		var unmarshalParams sdk_params_callback.SetConversationRecvMessageOptParams
		common.JsonUnmarshal(conversationIDList, &unmarshalParams, callback, operationID)
		c.setConversationRecvMessageOpt(callback, unmarshalParams, opt, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationRecvMessageOptCallback)
		log.NewInfo(operationID, "SetConversationRecvMessageOpt callback: ", sdk_params_callback.AddFriendCallback)
	}()

}
func (c *Conversation) GetConversationRecvMessageOpt(callback common.Base, conversationIDList, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetConversationRecvMessageOpt args: ", conversationIDList)
		var unmarshalParams sdk_params_callback.GetConversationRecvMessageOptParams
		common.JsonUnmarshal(conversationIDList, &unmarshalParams, callback, operationID)
		result := c.getConversationRecvMessageOpt(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetConversationRecvMessageOpt callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetOneConversation(callback common.Base, sessionType int32, sourceID, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetOneConversation args: ", sessionType, sourceID)
		result := c.getOneConversation(callback, sourceID, sessionType, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "GetRecvFriendApplicationList callback: ", utils.StructToJsonString(result))
	}()
}
func (c *Conversation) GetMultipleConversation(callback common.Base, conversationIDList string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetMultipleConversation args: ", conversationIDList)
		var unmarshalParams sdk_params_callback.GetMultipleConversationParams
		common.JsonUnmarshal(conversationIDList, &unmarshalParams, callback, operationID)
		result := c.getMultipleConversation(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetMultipleConversation callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) DeleteConversation(callback common.Base, conversationID string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "DeleteConversation args: ", conversationID)
		c.deleteConversation(callback, conversationID, operationID)
		callback.OnSuccess(sdk_params_callback.DeleteConversationCallback)
		//_ = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.TotalUnreadMessageChanged})
		log.NewInfo(operationID, "DeleteConversation callback: ", sdk_params_callback.DeleteConversationCallback)
	}()
}
func (c *Conversation) SetConversationDraft(callback common.Base, conversationID, draftText string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetConversationDraft args: ", conversationID)
		c.setConversationDraft(callback, conversationID, draftText, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationDraftCallback)
		//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		log.NewInfo(operationID, "SetConversationDraft callback: ", sdk_params_callback.SetConversationDraftCallback)
	}()
}
func (c *Conversation) PinConversation(callback common.Base, conversationID string, isPinned bool, operationID string) {

	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "PinConversation args: ", conversationID)
		c.pinConversation(callback, conversationID, isPinned, operationID)
		callback.OnSuccess(sdk_params_callback.PinConversationDraftCallback)
		//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
		log.NewInfo(operationID, "PinConversation callback: ", sdk_params_callback.PinConversationDraftCallback)
	}()
}
func (c *Conversation) GetTotalUnreadMsgCount(callback common.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetTotalUnreadMsgCount args: ")
		count, err := c.db.GetTotalUnreadMsgCount()
		common.CheckErr(callback, err, operationID)
		callback.OnSuccess(utils.Int64ToString(count))
		log.NewInfo(operationID, "GetTotalUnreadMsgCount callback: ", utils.Int64ToString(count))
	}()
}

type OnConversationListener interface {
	OnSyncServerStart()
	OnSyncServerFinish()
	OnSyncServerFailed()
	OnNewConversation(conversationList string)
	OnConversationChanged(conversationList string)
	OnTotalUnreadMessageCountChanged(totalUnreadCount int64)
}

//
func (c *Conversation) SetConversationListener(listener OnConversationListener) {
	if c.ConversationListener != nil {
		log.Error("internal", "just only set on listener")
		return
	}
	c.ConversationListener = listener
}

type OnAdvancedMsgListener interface {
	OnRecvNewMessage(message string)
	OnRecvC2CReadReceipt(msgReceiptList string)
	OnRecvMessageRevoked(msgId string)
}

//
func (c *Conversation) AddAdvancedMsgListener(listener OnAdvancedMsgListener) {
	if listener == nil {
		log.Error("internal", "AddAdvancedMsgListener listener is null")
		return
	}
	if len(c.MsgListenerList) == 1 {
		log.Error("internal", "u.ConversationListener.MsgListenerList == 1")
		return
	}
	c.MsgListenerList = append(c.MsgListenerList, listener)
}

////
////func (c *Conversation) ForceSyncMsg() bool {
////	if c.syncSeq2Msg() == nil {
////		return true
////	} else {
////		return false
////	}
////}
////
////func (c *Conversation) ForceSyncJoinedGroup() {
////	u.syncJoinedGroupInfo()
////}
////
////func (c *Conversation) ForceSyncJoinedGroupMember() {
////
////	u.syncJoinedGroupMember()
////}
////
////func (c *Conversation) ForceSyncGroupRequest() {
////	u.syncGroupRequest()
////}
////
////func (c *Conversation) ForceSyncSelfGroupRequest() {
////	u.syncSelfGroupRequest()
////}
//

func (c *Conversation) CreateTextMessage(text string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Text)
	s.Content = text
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateTextAtMessage(text, atUserList string) string {
	var users []string
	_ = json.Unmarshal([]byte(atUserList), &users)
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.AtText)
	s.AtElem.Text = text
	s.AtElem.AtUserList = users
	s.Content = utils.StructToJsonString(s.AtElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateLocationMessage(description string, longitude, latitude float64) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Location)
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = utils.StructToJsonString(s.LocationElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateCustomMessage(data, extension string, description string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Custom)
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = utils.StructToJsonString(s.CustomElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateQuoteMessage(text string, message string) string {
	s, qs := sdk_struct.MsgStruct{}, sdk_struct.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	c.initBasicInfo(&s, constant.UserMsgType, constant.Quote)
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.QuoteMessage = &qs
	s.Content = utils.StructToJsonString(s.QuoteElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateCardMessage(cardInfo string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Card)
	s.Content = cardInfo
	return utils.StructToJsonString(s)

}
func (c *Conversation) CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(videoFullPath, c.DbDir) //a->b
		s, err := utils.CopyFile(videoFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, videoFullPath)
		}
		log.Info("internal", "videoFullPath dstFile", videoFullPath, dstFile, s)
		dstFile = utils.FileTmpPath(snapshotFullPath, c.DbDir) //a->b
		s, err = utils.CopyFile(snapshotFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, snapshotFullPath)
		}
		log.Info("internal", "snapshotFullPath dstFile", snapshotFullPath, dstFile, s)
		wg.Done()
	}()

	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video)
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
		log.Error("internal", "get file Attributes error", err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotFullPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			log.Error("internal", "get Image Attributes error", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	wg.Wait()
	s.Content = utils.StructToJsonString(s.VideoElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateFileMessageFromFullPath(fileFullPath string, fileName string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(fileFullPath, c.DbDir)
		_, err := utils.CopyFile(fileFullPath, dstFile)
		log.Info("internal", "copy file, ", fileFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err.Error(), fileFullPath)
		}
		wg.Done()
	}()
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.File)
	s.FileElem.FilePath = fileFullPath
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		log.Error("internal", "get file Attributes error", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.FileElem.FileName = fileName
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateImageMessageFromFullPath(imageFullPath string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(imageFullPath, c.DbDir) //a->b
		_, err := utils.CopyFile(imageFullPath, dstFile)
		log.Info("internal", "copy file, ", imageFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, imageFullPath)
		}
		wg.Done()
	}()

	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.PictureElem.SourcePath = imageFullPath
	log.Info("internal", "ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		log.Error("internal", "getImageInfo err:", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	wg.Wait()
	s.Content = utils.StructToJsonString(s.PictureElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateSoundMessageFromFullPath(soundPath string, duration int64) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(soundPath, c.DbDir) //a->b
		_, err := utils.CopyFile(soundPath, dstFile)
		log.Info("internal", "copy file, ", soundPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, soundPath)
		}
		wg.Done()
	}()
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
	s.SoundElem.SoundPath = soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		log.Error("internal", "getSoundInfo err:", err.Error(), s.SoundElem.SoundPath)
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	wg.Wait()
	s.Content = utils.StructToJsonString(s.SoundElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateImageMessage(imagePath string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.PictureElem.SourcePath = c.DbDir + imagePath
	log.Debug("internal", "ImageMessage  path:", s.PictureElem.SourcePath)
	imageInfo, err := getImageInfo(s.PictureElem.SourcePath)
	if err != nil {
		log.Error("internal", "get imageInfo err", err.Error())
		return ""
	}
	s.PictureElem.SourcePicture.Width = imageInfo.Width
	s.PictureElem.SourcePicture.Height = imageInfo.Height
	s.PictureElem.SourcePicture.Type = imageInfo.Type
	s.PictureElem.SourcePicture.Size = imageInfo.Size
	s.Content = utils.StructToJsonString(s.PictureElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture string) string {
	s := sdk_struct.MsgStruct{}
	var p sdk_struct.PictureBaseInfo
	_ = json.Unmarshal([]byte(sourcePicture), &p)
	s.PictureElem.SourcePicture = p
	_ = json.Unmarshal([]byte(bigPicture), &p)
	s.PictureElem.BigPicture = p
	_ = json.Unmarshal([]byte(snapshotPicture), &p)
	s.PictureElem.SnapshotPicture = p
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture)
	s.Content = utils.StructToJsonString(s.PictureElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) SendMessage(callback SendMsgCallBack, message, recvID, groupID string, offlinePushInfo string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		//参数校验
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		p := server_api_params.OfflinePushInfo{}
		common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
		if recvID == "" && groupID == "" {
			common.CheckAnyErr(callback, 201, errors.New("recvID && groupID not be allowed"), operationID)
		}
		var localMessage db.LocalChatLog
		var conversationID string
		var options map[string]bool
		lc := db.LocalConversation{
			LatestMsgSendTime: s.CreateTime,
		}
		//根据单聊群聊类型组装消息和会话
		if recvID == "" {
			s.SessionType = constant.GroupChatType
			s.GroupID = groupID
			conversationID = c.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			lc.GroupID = groupID
			lc.ConversationType = constant.GroupChatType
			g, err := c.db.GetGroupInfoByGroupID(groupID)
			common.CheckAnyErr(callback, 202, err, operationID)
			lc.ShowName = g.GroupName
			lc.FaceURL = g.FaceURL
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			common.CheckAnyErr(callback, 202, err, operationID)
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				common.CheckAnyErr(callback, 208, errors.New("you not exist in this group"), operationID)
			}
		} else {
			s.SessionType = constant.SingleChatType
			s.RecvID = recvID
			conversationID = c.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			lc.UserID = recvID
			lc.ConversationType = constant.SingleChatType
			faceUrl, name, err := c.getUserNameAndFaceUrlByUid(callback, recvID, operationID)
			common.CheckAnyErr(callback, 301, err, operationID)
			lc.FaceURL = faceUrl
			lc.ShowName = name
		}
		lc.ConversationID = conversationID
		lc.LatestMsg = utils.StructToJsonString(s)
		msgStructToLocalChatLog(&s, &localMessage)
		err := c.db.InsertMessage(&localMessage)
		common.CheckAnyErr(callback, 201, err, operationID)
		//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
		//c}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		options = make(map[string]bool, 2)

		var delFile []string
		//media file handle
		switch s.ContentType {
		case constant.Picture:
			var sourcePath string
			if utils.FileExist(s.PictureElem.SourcePath) {
				sourcePath = s.PictureElem.SourcePath
				delFile = append(delFile, utils.FileTmpPath(s.PictureElem.SourcePath, c.DbDir))
			} else {
				sourcePath = utils.FileTmpPath(s.PictureElem.SourcePath, c.DbDir)
				delFile = append(delFile, sourcePath)
			}
			log.Info(operationID, "file", sourcePath, delFile)
			sourceUrl, uuid, err := c.uploadImage(sourcePath, callback)
			c.checkErrAndUpdateMessage(callback, 301, err, &s, &lc, operationID)
			s.PictureElem.SourcePicture.Url = sourceUrl
			s.PictureElem.SourcePicture.UUID = uuid
			s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + constant.ZoomScale + "/h/" + constant.ZoomScale
			s.PictureElem.SnapshotPicture.Width = int32(utils.StringToInt(constant.ZoomScale))
			s.PictureElem.SnapshotPicture.Height = int32(utils.StringToInt(constant.ZoomScale))
			s.Content = utils.StructToJsonString(s.PictureElem)

		case constant.Voice:
			var sourcePath string
			if utils.FileExist(s.SoundElem.SoundPath) {
				sourcePath = s.SoundElem.SoundPath
				delFile = append(delFile, utils.FileTmpPath(s.SoundElem.SoundPath, c.DbDir))
			} else {
				sourcePath = utils.FileTmpPath(s.SoundElem.SoundPath, c.DbDir)
				delFile = append(delFile, sourcePath)
			}
			log.Info(operationID, "file", sourcePath, delFile)
			soundURL, uuid, err := c.uploadSound(sourcePath, callback)
			c.checkErrAndUpdateMessage(callback, 301, err, &s, &lc, operationID)
			s.SoundElem.SourceURL = soundURL
			s.SoundElem.UUID = uuid
			s.Content = utils.StructToJsonString(s.SoundElem)

		case constant.Video:
			var videoPath string
			var snapPath string
			if utils.FileExist(s.VideoElem.VideoPath) {
				videoPath = s.VideoElem.VideoPath
				snapPath = s.VideoElem.SnapshotPath
				delFile = append(delFile, utils.FileTmpPath(s.VideoElem.VideoPath, c.DbDir))
				delFile = append(delFile, utils.FileTmpPath(s.VideoElem.SnapshotPath, c.DbDir))
			} else {
				videoPath = utils.FileTmpPath(s.VideoElem.VideoPath, c.DbDir)
				snapPath = utils.FileTmpPath(s.VideoElem.SnapshotPath, c.DbDir)
				delFile = append(delFile, videoPath)
				delFile = append(delFile, snapPath)
			}
			log.Info(operationID, "file: ", videoPath, snapPath, delFile)
			snapshotURL, snapshotUUID, videoURL, videoUUID, err := c.uploadVideo(videoPath, snapPath, callback)
			c.checkErrAndUpdateMessage(callback, 301, err, &s, &lc, operationID)
			s.VideoElem.VideoURL = videoURL
			s.VideoElem.SnapshotUUID = snapshotUUID
			s.VideoElem.SnapshotURL = snapshotURL
			s.VideoElem.VideoUUID = videoUUID
			s.Content = utils.StructToJsonString(s.VideoElem)
		case constant.File:
			fileURL, fileUUID, err := c.uploadFile(s.FileElem.FilePath, callback)
			c.checkErrAndUpdateMessage(callback, 301, err, &s, &lc, operationID)
			s.FileElem.SourceURL = fileURL
			s.FileElem.UUID = fileUUID
			s.Content = utils.StructToJsonString(s.FileElem)
		case constant.Text:
		case constant.AtText:
		case constant.Location:
		case constant.Custom:
		case constant.Merger:
		case constant.Quote:
		case constant.Card:
		default:
			common.CheckAnyErr(callback, 202, errors.New("contentType not currently supported"+utils.Int32ToString(s.ContentType)), operationID)
		}
		msgStructToLocalChatLog(&s, &localMessage)
		err = c.db.UpdateMessage(&localMessage)
		common.CheckAnyErr(callback, 201, err, operationID)
		c.sendMessageToServer(&s, &lc, callback, delFile, &p, options, operationID)
	}()
}
func msgStructToLocalChatLog(s *sdk_struct.MsgStruct, localMessage *db.LocalChatLog) {
	copier.Copy(localMessage, s)
}
func (c *Conversation) checkErrAndUpdateMessage(callback SendMsgCallBack, errCode int32, err error, s *sdk_struct.MsgStruct, lc *db.LocalConversation, operationID string) {
	if err != nil {
		if callback != nil {
			c.updateMsgStatusAndTriggerConversation(s.ClientMsgID, "", s.CreateTime, constant.MsgStatusSendFailed, s, lc, operationID)
			errInfo := "operationID[" + operationID + "], " + "info[" + err.Error() + "]"
			log.NewError(operationID, "checkErr ", errInfo)
			callback.OnError(errCode, errInfo)
			runtime.Goexit()
		}
	}
}
func (c *Conversation) updateMsgStatusAndTriggerConversation(clientMsgID, serverMsgID string, sendTime uint32, status int32, s *sdk_struct.MsgStruct, lc *db.LocalConversation, operationID string) {
	_ = c.db.UpdateMessageTimeAndStatus(clientMsgID, sendTime, status)
	s.SendTime = sendTime
	s.Status = status
	s.ServerMsgID = serverMsgID
	lc.LatestMsg = utils.StructToJsonString(s)
	lc.LatestMsgSendTime = sendTime
	//会话数据库操作，触发UI会话回调
}
func (c *Conversation) getUserNameAndFaceUrlByUid(callback SendMsgCallBack, friendUserID, operationID string) (faceUrl, name string, err error) {
	friendInfo, _ := c.db.GetFriendInfoByFriendUserID(friendUserID)
	if friendInfo != nil {
		if friendInfo.Remark != "" {
			return friendInfo.FaceURL, friendInfo.Remark, nil
		} else {
			return friendInfo.FaceURL, friendInfo.Nickname, nil
		}
	} else {
		userInfos := c.user.GetUsersInfoFromSvr(callback, []string{friendUserID}, operationID)
		for _, v := range userInfos {
			return v.FaceURL, v.Nickname, nil
		}
	}
	return "", "", errors.New("getUserNameAndFaceUrlByUid err")

}

func (c *Conversation) internalSendMessage(callback common.Base, message, recvID, groupID, offlinePushInfo, operationID string, onlineUserOnly bool, options map[string]bool) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	p := server_api_params.OfflinePushInfo{}
	common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
	if recvID == "" && groupID == "" {
		common.CheckAnyErr(callback, 201, errors.New("recvID && groupID not be allowed"), operationID)
	}
	if recvID == "" {
		s.SessionType = constant.GroupChatType
		s.GroupID = groupID
		groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
		common.CheckAnyErr(callback, 202, err, operationID)
		if !utils.IsContain(s.SendID, groupMemberUidList) {
			common.CheckAnyErr(callback, 208, errors.New("you not exist in this group"), operationID)
		}

	} else {
		s.SessionType = constant.SingleChatType
		s.RecvID = recvID
	}

	if onlineUserOnly {
		options[constant.IsHistory] = false
		options[constant.IsPersistent] = false
	}

	var wsMsgData server_api_params.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = int64(s.CreateTime)
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = &p
	timeout := 300
	retryTimes := 0
	_, err := c.SendReqWaitResp(&wsMsgData, constant.WSSendMsg, timeout, retryTimes, c.loginUserID, operationID)
	common.CheckAnyErr(callback, 301, err, operationID)

}

func (c *Conversation) SendMessageNotOss(callback SendMsgCallBack, message, recvID, groupID string, offlinePushInfo string, operationID string) {
	go func() {
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		p := server_api_params.OfflinePushInfo{}
		common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
		if recvID == "" && groupID == "" {
			common.CheckAnyErr(callback, 201, errors.New("recvID && groupID not be allowed"), operationID)
		}
		var localMessage db.LocalChatLog
		var conversationID string
		var options map[string]bool
		lc := db.LocalConversation{
			LatestMsgSendTime: s.CreateTime,
		}
		//根据单聊群聊类型组装消息和会话
		if recvID == "" {
			s.SessionType = constant.GroupChatType
			s.GroupID = groupID
			conversationID = c.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			lc.GroupID = groupID
			lc.ConversationType = constant.GroupChatType
			g, err := c.db.GetGroupInfoByGroupID(groupID)
			common.CheckAnyErr(callback, 202, err, operationID)
			lc.ShowName = g.GroupName
			lc.FaceURL = g.FaceURL
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			common.CheckAnyErr(callback, 202, err, operationID)
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				common.CheckAnyErr(callback, 208, errors.New("you not exist in this group"), operationID)
			}
		} else {
			s.SessionType = constant.SingleChatType
			s.RecvID = recvID
			conversationID = c.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			lc.UserID = recvID
			lc.ConversationType = constant.SingleChatType
			faceUrl, name, err := c.getUserNameAndFaceUrlByUid(callback, recvID, operationID)
			common.CheckAnyErr(callback, 301, err, operationID)
			lc.FaceURL = faceUrl
			lc.ShowName = name
		}
		lc.ConversationID = conversationID
		lc.LatestMsg = utils.StructToJsonString(s)
		msgStructToLocalChatLog(&s, &localMessage)
		err := c.db.InsertMessage(&localMessage)
		common.CheckAnyErr(callback, 201, err, operationID)
		//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
		//c}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", NewConChange, []string{conversationID}}})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		options = make(map[string]bool, 2)
		var delFile []string

		msgStructToLocalChatLog(&s, &localMessage)
		err = c.db.UpdateMessage(&localMessage)
		common.CheckAnyErr(callback, 201, err, operationID)
		c.sendMessageToServer(&s, &lc, callback, delFile, &p, options, operationID)

	}()
}

func (c *Conversation) CreateSoundMessageByURL(soundBaseInfo string) string {
	s := sdk_struct.MsgStruct{}
	var soundElem sdk_struct.SoundBaseInfo
	_ = json.Unmarshal([]byte(soundBaseInfo), &soundElem)
	s.SoundElem = soundElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
	s.Content = utils.StructToJsonString(s.SoundElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateSoundMessage(soundPath string, duration int64) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice)
	s.SoundElem.SoundPath = c.DbDir + soundPath
	s.SoundElem.Duration = duration
	fi, err := os.Stat(s.SoundElem.SoundPath)
	if err != nil {
		log.Error("internal", "get sound info err", err.Error())
		return ""
	}
	s.SoundElem.DataSize = fi.Size()
	s.Content = utils.StructToJsonString(s.SoundElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateVideoMessageByURL(videoBaseInfo string) string {
	s := sdk_struct.MsgStruct{}
	var videoElem sdk_struct.VideoBaseInfo
	_ = json.Unmarshal([]byte(videoBaseInfo), &videoElem)
	s.VideoElem = videoElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video)
	s.Content = utils.StructToJsonString(s.VideoElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video)
	s.VideoElem.VideoPath = c.DbDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = c.DbDir + snapshotPath
	}
	fi, err := os.Stat(s.VideoElem.VideoPath)
	if err != nil {
		log.Error("internal", "get video file error", err.Error())
		return ""
	}
	s.VideoElem.VideoSize = fi.Size()
	if snapshotPath != "" {
		imageInfo, err := getImageInfo(s.VideoElem.SnapshotPath)
		if err != nil {
			log.Error("internal", "get snapshot info ", err.Error())
			return ""
		}
		s.VideoElem.SnapshotHeight = imageInfo.Height
		s.VideoElem.SnapshotWidth = imageInfo.Width
		s.VideoElem.SnapshotSize = imageInfo.Size
	}
	s.Content = utils.StructToJsonString(s.VideoElem)
	return utils.StructToJsonString(s)
}

func (c *Conversation) CreateFileMessageByURL(fileBaseInfo string) string {
	s := sdk_struct.MsgStruct{}
	var fileElem sdk_struct.FileBaseInfo
	_ = json.Unmarshal([]byte(fileBaseInfo), &fileElem)
	s.FileElem = fileElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.File)
	s.Content = utils.StructToJsonString(s.FileElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateFileMessage(filePath string, fileName string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.File)
	s.FileElem.FilePath = c.DbDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		log.Error("internal", "get file message err", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateMergerMessage(messageList, title, summaryList string) string {
	var messages []*sdk_struct.MsgStruct
	var summaries []string
	s := sdk_struct.MsgStruct{}
	err := json.Unmarshal([]byte(messageList), &messages)
	if err != nil {
		log.Error("internal", "messages Unmarshal err", err.Error())
		return ""
	}
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	c.initBasicInfo(&s, constant.UserMsgType, constant.Merger)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.StructToJsonString(s.MergeElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateForwardMessage(m string) string {
	s := sdk_struct.MsgStruct{}
	err := json.Unmarshal([]byte(m), &s)
	if err != nil {
		log.Error("internal", "messages Unmarshal err", err.Error())
		return ""
	}
	if s.Status != constant.MsgStatusSendSuccess {
		log.Error("internal", "only send success message can be revoked")
		return ""
	}
	c.initBasicInfo(&s, constant.UserMsgType, s.ContentType)
	//Forward message seq is set to 0
	s.Seq = 0
	return utils.StructToJsonString(s)
}

func (c *Conversation) sendMessageToServer(s *sdk_struct.MsgStruct, lc *db.LocalConversation, callback SendMsgCallBack,
	delFile []string, offlinePushInfo *server_api_params.OfflinePushInfo, options map[string]bool, operationID string) {
	//Protocol conversion
	var wsMsgData server_api_params.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = int64(s.CreateTime)
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = offlinePushInfo
	timeout := 300
	retryTimes := 6
	resp, err := c.SendReqWaitResp(&wsMsgData, constant.WSSendMsg, timeout, retryTimes, c.loginUserID, operationID)
	c.checkErrAndUpdateMessage(callback, 302, err, s, lc, operationID)
	callback.OnProgress(100)
	callback.OnSuccess("")
	//remove media cache file
	for _, v := range delFile {
		err := os.Remove(v)
		if err != nil {
			log.Error(operationID, "remove failed,", err.Error(), v)
		}
		log.Debug(operationID, "remove file: ", v)
	}
	var sendMsgResp server_api_params.UserSendMsgResp
	_ = proto.Unmarshal(resp.Data, &sendMsgResp)
	c.updateMsgStatusAndTriggerConversation(sendMsgResp.ClientMsgID, sendMsgResp.ServerMsgID, uint32(sendMsgResp.SendTime), constant.MsgStatusSendSuccess, s, lc, operationID)

}

func (c *Conversation) GetHistoryMessageList(callback common.Base, getMessageOptions, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetHistoryMessageList args: ", getMessageOptions)
		var unmarshalParams sdk_params_callback.GetHistoryMessageListParams
		common.JsonUnmarshal(getMessageOptions, &unmarshalParams, callback, operationID)
		result := c.getHistoryMessageList(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetMultipleConversation callback: ", utils.StructToJsonStringDefault(result))
	}()
	go func() {

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

//func (c *Conversation) RevokeMessage(callback common.Base, message string) {
//	go func() {
//		//var receiver, groupID string
//		c := utils.MsgStruct{}
//		err := json.Unmarshal([]byte(message), &c)
//		if err != nil {
//			callback.OnError(200, err.Error())
//			return
//		}
//		s, err := u.getOneMessage(c.ClientMsgID)
//		if err != nil || s == nil {
//			callback.OnError(201, "getOneMessage err")
//			return
//		}
//		if s.Status != constant.MsgStatusSendSuccess {
//			callback.OnError(201, "only send success message can be revoked")
//			return
//		}
//		utils.sdkLog("test data", s)
//		//Send message internally
//		switch s.SessionType {
//		case constant.SingleChatType:
//			//receiver = s.RecvID
//		case constant.GroupChatType:
//			//groupID = s.GroupID
//		default:
//			callback.OnError(200, "args err")
//		}
//		s.Content = s.ClientMsgID
//		s.ClientMsgID = utils.getMsgID(s.SendID)
//		s.ContentType = constant.Revoke
//		//err = u.autoSendMsg(s, receiver, groupID, false, true, false)
//		if err != nil {
//			utils.sdkLog("autoSendMsg revokeMessage err:", err.Error())
//			callback.OnError(300, err.Error())
//
//		} else {
//			err = u.setMessageStatus(s.Content, constant.MsgStatusRevoked)
//			if err != nil {
//				utils.sdkLog("setLocalMessageStatus revokeMessage err:", err.Error())
//				callback.OnError(300, err.Error())
//			} else {
//				callback.OnSuccess("")
//			}
//		}
//	}()
//}
//func (c *Conversation) TypingStatusUpdate(receiver, msgTip string) {
//	go func() {
//		s := utils.MsgStruct{}
//		u.initBasicInfo(&s, constant.UserMsgType, constant.Typing)
//		s.Content = msgTip
//		//err := u.autoSendMsg(&s, receiver, "", true, false, false)
//		//if err != nil {
//		//	sdkLog("TypingStatusUpdate err:", err)
//		//} else {
//		//	sdkLog("TypingStatusUpdate success!!!")
//		//}
//	}()
//}
//
//func (c *Conversation) MarkC2CMessageAsRead(callback common.Base, receiver string, msgIDList string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(receiver, constant.SingleChatType)
//		var list []string
//		err := json.Unmarshal([]byte(msgIDList), &list)
//		if err != nil {
//			callback.OnError(201, "json unmarshal err")
//			return
//		}
//		if len(list) == 0 {
//			callback.OnError(200, "msg list is null")
//			return
//		}
//		s := utils.MsgStruct{}
//		u.initBasicInfo(&s, constant.UserMsgType, constant.HasReadReceipt)
//		s.Content = msgIDList
//		utils.sdkLog("MarkC2CMessageAsRead: send Message")
//		//err = u.autoSendMsg(&s, receiver, "", false, false, false)
//		if err != nil {
//			utils.sdkLog("MarkC2CMessageAsRead  err:", err.Error())
//			callback.OnError(300, err.Error())
//		} else {
//			callback.OnSuccess("")
//			err = u.setSingleMessageHasReadByMsgIDList(receiver, list)
//			if err != nil {
//				utils.sdkLog("setSingleMessageHasReadByMsgIDList  err:", err.Error())
//			}
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.UpdateLatestMessageChange, ""}})
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		}
//	}()
//}
//
////Deprecated
//func (c *Conversation) MarkSingleMessageHasRead(callback common.Base, userID string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
//		//if err := u.setSingleMessageHasRead(userID); err != nil {
//		//	callback.OnError(201, err.Error())
//		//} else {
//		callback.OnSuccess("")
//		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
//		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		//}
//	}()
//}
//func (c *Conversation) MarkAllConversationHasRead(callback common.Base, userID string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
//		//if err := u.setSingleMessageHasRead(userID); err != nil {
//		//	callback.OnError(201, err.Error())
//		//} else {
//		callback.OnSuccess("")
//		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
//		u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		//}
//	}()
//}
//func (c *Conversation) MarkGroupMessageHasRead(callback common.Base, groupID string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
//		if err := u.setGroupMessageHasRead(groupID); err != nil {
//			callback.OnError(201, err.Error())
//		} else {
//			callback.OnSuccess("")
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{ConId: conversationID, Action: constant.UnreadCountSetZero}})
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		}
//	}()
//}
//func (c *Conversation) DeleteMessageFromLocalStorage(callback common.Base, message string) {
//	go func() {
//		var conversation ConversationStruct
//		var latestMsg utils.MsgStruct
//		var conversationID string
//		var sourceID string
//		s := utils.MsgStruct{}
//		err := json.Unmarshal([]byte(message), &s)
//		if err != nil {
//			callback.OnError(200, err.Error())
//			return
//		}
//		err = u.setMessageStatus(s.ClientMsgID, constant.MsgStatusHasDeleted)
//		if err != nil {
//			callback.OnError(202, err.Error())
//			return
//		}
//		callback.OnSuccess("")
//		if s.SessionType == constant.GroupChatType {
//			conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.GroupChatType)
//			sourceID = s.RecvID
//
//		} else if s.SessionType == constant.SingleChatType {
//			if s.SendID != u.loginUserID {
//				conversationID = utils.GetConversationIDBySessionType(s.SendID, constant.SingleChatType)
//				sourceID = s.SendID
//			} else {
//				conversationID = utils.GetConversationIDBySessionType(s.RecvID, constant.SingleChatType)
//				sourceID = s.RecvID
//			}
//		}
//		_, m := u.getConversationLatestMsgModel(conversationID)
//		if m != "" {
//			err := json.Unmarshal([]byte(m), &latestMsg)
//			if err != nil {
//				utils.sdkLog("DeleteMessage err :", err)
//				callback.OnError(200, err.Error())
//				return
//			}
//		} else {
//			utils.sdkLog("err ,conversation has been deleted")
//		}
//
//		if s.ClientMsgID == latestMsg.ClientMsgID { //If the deleted message is the latest message of the conversation, update the latest message of the conversation
//			err, list := u.getHistoryMessage(sourceID, s.SendTime+TimeOffset, 1, int(s.SessionType))
//			if err != nil {
//				utils.sdkLog("DeleteMessageFromLocalStorage database err:", err.Error())
//			}
//			conversation.ConversationID = conversationID
//			if list == nil {
//				conversation.LatestMsg = ""
//				conversation.LatestMsgSendTime = utils.getCurrentTimestampByNano()
//			} else {
//				conversation.LatestMsg = utils.structToJsonString(list[0])
//				conversation.LatestMsgSendTime = list[0].SendTime
//			}
//			err = u.triggerCmdUpdateConversation(common.updateConNode{ConId: conversationID, Action: constant.AddConOrUpLatMsg, Args: conversation})
//			if err != nil {
//				utils.sdkLog("DeleteMessageFromLocalStorage triggerCmdUpdateConversation err:", err.Error())
//			}
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//
//		}
//	}()
//}
//func (c *Conversation) ClearC2CHistoryMessage(callback common.Base, userID string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
//		err := u.setMessageStatusBySourceID(userID, constant.MsgStatusHasDeleted, constant.SingleChatType)
//		if err != nil {
//			callback.OnError(202, err.Error())
//			return
//		}
//		err = u.clearConversation(conversationID)
//		if err != nil {
//			callback.OnError(203, err.Error())
//			return
//		} else {
//			callback.OnSuccess("")
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		}
//	}()
//}
//func (c *Conversation) ClearGroupHistoryMessage(callback common.Base, groupID string) {
//	go func() {
//		conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
//		err := u.setMessageStatusBySourceID(groupID, constant.MsgStatusHasDeleted, constant.GroupChatType)
//		if err != nil {
//			callback.OnError(202, err.Error())
//			return
//		}
//		err = u.clearConversation(conversationID)
//		if err != nil {
//			callback.OnError(203, err.Error())
//			return
//		} else {
//			callback.OnSuccess("")
//			u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{"", constant.NewConChange, []string{conversationID}}})
//		}
//	}()
//}
//
//func (c *Conversation) InsertSingleMessageToLocalStorage(callback common.Base, message, userID, sender string) string {
//	s := utils.MsgStruct{}
//	err := json.Unmarshal([]byte(message), &s)
//	if err != nil {
//		callback.OnError(200, err.Error())
//		return ""
//	}
//	s.SendID = sender
//	s.RecvID = userID
//	//Generate client message primary key
//	s.ClientMsgID = utils.getMsgID(s.SendID)
//	s.SendTime = utils.getCurrentTimestampByNano()
//	go func() {
//		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
//			callback.OnError(201, err.Error())
//		} else {
//			callback.OnSuccess("")
//		}
//	}()
//	return s.ClientMsgID
//}
//
//func (c *Conversation) InsertGroupMessageToLocalStorage(callback common.Base, message, groupID, sender string) string {
//	s := utils.MsgStruct{}
//	err := json.Unmarshal([]byte(message), &s)
//	if err != nil {
//		callback.OnError(200, err.Error())
//		return ""
//	}
//	s.SendID = sender
//	s.RecvID = groupID
//	//Generate client message primary key
//	s.ClientMsgID = utils.getMsgID(s.SendID)
//	s.SendTime = utils.getCurrentTimestampByNano()
//	go func() {
//		if err = u.insertMessageToLocalOrUpdateContent(&s); err != nil {
//			callback.OnError(201, err.Error())
//		} else {
//			callback.OnSuccess("")
//		}
//	}()
//	return s.ClientMsgID
//}
//
//func (c *Conversation) FindMessages(callback common.Base, messageIDList string) {
//	go func() {
//		var c []string
//		err := json.Unmarshal([]byte(messageIDList), &c)
//		if err != nil {
//			callback.OnError(200, err.Error())
//			utils.sdkLog("Unmarshal failed, ", err.Error())
//
//		}
//		err, list := u.getMultipleMessageModel(c)
//		if err != nil {
//			callback.OnError(203, err.Error())
//		} else {
//			if list != nil {
//				callback.OnSuccess(utils.structToJsonString(list))
//			} else {
//				callback.OnSuccess(utils.structToJsonString([]utils.MsgStruct{}))
//			}
//		}
//	}()
//}
func getImageInfo(filePath string) (*sdk_struct.ImageInfo, error) {
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

	return &sdk_struct.ImageInfo{int32(b.Max.X), int32(b.Max.Y), datatype, fi.Size()}, nil

}
