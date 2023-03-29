package conversation_msg

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	imgtype "github.com/shamsher31/goimgtype"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func (c *Conversation) GetAllConversationList(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "GetAllConversationList args: ")
		result := c.getAllConversationList(callback, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetAllConversationList callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetConversationListSplit(callback open_im_sdk_callback.Base, offset, count int, operationID string) {
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
func (c *Conversation) SetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList string, opt int, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetConversationRecvMessageOpt args: ", conversationIDList, opt)
		var unmarshalParams sdk_params_callback.SetConversationRecvMessageOptParams
		common.JsonUnmarshalCallback(conversationIDList, &unmarshalParams, callback, operationID)
		c.setConversationRecvMessageOpt(callback, unmarshalParams, opt, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationRecvMessageOptCallback)
		log.NewInfo(operationID, "SetConversationRecvMessageOpt callback: ", sdk_params_callback.SetConversationRecvMessageOptCallback)
	}()
}
func (c *Conversation) SetGlobalRecvMessageOpt(callback open_im_sdk_callback.Base, opt int, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetGlobalRecvMessageOpt args: ", opt)
		c.setGlobalRecvMessageOpt(callback, int32(opt), operationID)
		callback.OnSuccess(sdk_params_callback.SetGlobalRecvMessageOptCallback)
		log.NewInfo(operationID, "SetGlobalRecvMessageOpt callback: ", sdk_params_callback.SetGlobalRecvMessageOptCallback)
	}()
}
func (c *Conversation) HideConversation(callback open_im_sdk_callback.Base, conversationID string, operationID string) {
	c.hideConversation(callback, conversationID, operationID)
	callback.OnSuccess(sdk_params_callback.HideConversationCallback)
	log.NewInfo(operationID, "HideConversation callback: ", sdk_params_callback.HideConversationCallback)

}

// deprecated
func (c *Conversation) GetConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationIDList, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetConversationRecvMessageOpt args: ", conversationIDList)
		var unmarshalParams sdk_params_callback.GetConversationRecvMessageOptParams
		common.JsonUnmarshalCallback(conversationIDList, &unmarshalParams, callback, operationID)
		result := c.getConversationRecvMessageOpt(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetConversationRecvMessageOpt callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetOneConversation(callback open_im_sdk_callback.Base, sessionType int32, sourceID, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetOneConversation args: ", sessionType, sourceID)
		result := c.getOneConversation(callback, sourceID, sessionType, operationID)
		callback.OnSuccess(utils.StructToJsonString(result))
		log.NewInfo(operationID, "GetOneConversation callback: ", utils.StructToJsonString(result))
	}()
}
func (c *Conversation) GetMultipleConversation(callback open_im_sdk_callback.Base, conversationIDList string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetMultipleConversation args: ", conversationIDList)
		var unmarshalParams sdk_params_callback.GetMultipleConversationParams
		common.JsonUnmarshalCallback(conversationIDList, &unmarshalParams, callback, operationID)
		result := c.getMultipleConversation(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetMultipleConversation callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) DeleteConversation(callback open_im_sdk_callback.Base, conversationID string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "DeleteConversation args: ", conversationID)
		c.deleteConversation(callback, conversationID, operationID)
		callback.OnSuccess(sdk_params_callback.DeleteConversationCallback)
		log.NewInfo(operationID, "DeleteConversation callback: ", sdk_params_callback.DeleteConversationCallback)
	}()
}
func (c *Conversation) DeleteAllConversationFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "DeleteAllConversationFromLocal args: ")
		err := c.db.ResetAllConversation()
		common.CheckDBErrCallback(callback, err, operationID)
		callback.OnSuccess(sdk_params_callback.DeleteAllConversationFromLocalCallback)
		log.NewInfo(operationID, "DeleteConversation callback: ", sdk_params_callback.DeleteAllConversationFromLocalCallback)
	}()
}
func (c *Conversation) SetConversationDraft(callback open_im_sdk_callback.Base, conversationID, draftText string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "SetConversationDraft args: ", conversationID)
		c.setConversationDraft(callback, conversationID, draftText, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationDraftCallback)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		log.NewInfo(operationID, "SetConversationDraft callback: ", sdk_params_callback.SetConversationDraftCallback)
	}()
}
func (c *Conversation) ResetConversationGroupAtType(callback open_im_sdk_callback.Base, conversationID, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "ResetConversationGroupAtType args: ", conversationID)
		c.setOneConversationGroupAtType(callback, conversationID, operationID)
		callback.OnSuccess(sdk_params_callback.ResetConversationGroupAtTypeCallback)
		log.NewInfo(operationID, "ResetConversationGroupAtType callback: ", sdk_params_callback.ResetConversationGroupAtTypeCallback)
	}()
}
func (c *Conversation) PinConversation(callback open_im_sdk_callback.Base, conversationID string, isPinned bool, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "PinConversation args: ", conversationID, isPinned)
		c.pinConversation(callback, conversationID, isPinned, operationID)
		callback.OnSuccess(sdk_params_callback.PinConversationDraftCallback)
		log.NewInfo(operationID, "PinConversation callback: ", sdk_params_callback.PinConversationDraftCallback)
	}()
}

func (c *Conversation) SetOneConversationPrivateChat(callback open_im_sdk_callback.Base, conversationID string, isPrivate bool, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", conversationID, isPrivate)
		c.setOneConversationPrivateChat(callback, conversationID, isPrivate, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationMessageOptCallback)
		log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", sdk_params_callback.SetConversationMessageOptCallback)
	}()
}

func (c *Conversation) SetOneConversationBurnDuration(callback open_im_sdk_callback.Base, conversationID string, burnDuration int32, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", conversationID, burnDuration)
		c.setOneConversationBurnDuration(callback, conversationID, burnDuration, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationBurnDurationOptCallback)
		log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", sdk_params_callback.SetConversationBurnDurationOptCallback)
	}()
}

func (c *Conversation) SetOneConversationRecvMessageOpt(callback open_im_sdk_callback.Base, conversationID string, opt int, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", conversationID, opt)
		c.setOneConversationRecvMessageOpt(callback, conversationID, opt, operationID)
		callback.OnSuccess(sdk_params_callback.SetConversationMessageOptCallback)
		log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", sdk_params_callback.SetConversationMessageOptCallback)
	}()
}

func (c *Conversation) GetTotalUnreadMsgCount(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetTotalUnreadMsgCount args: ")
		count, err := c.db.GetTotalUnreadMsgCountDB()
		common.CheckDBErrCallback(callback, err, operationID)
		callback.OnSuccess(utils.Int32ToString(count))
		log.NewInfo(operationID, "GetTotalUnreadMsgCount callback: ", utils.Int32ToString(count))
	}()
}

func (c *Conversation) SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	if c.ConversationListener != nil {
		log.Error("internal", "just only set on listener")
		return
	}
	c.ConversationListener = listener
}

func (c *Conversation) GetConversationsByUserID(callback open_im_sdk_callback.Base, operationID string, UserID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, utils.GetSelfFuncName())
		conversations, err := c.db.GetAllConversationListDB()
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
		}
		var conversationIDs []string
		for _, conversation := range conversations {
			conversationIDs = append(conversationIDs, conversation.ConversationID)
		}
	}()
}

//

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

func (c *Conversation) CreateTextMessage(text, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Text, operationID)
	s.Content = text
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateAdvancedTextMessage(text, messageEntityList, operationID string) string {
	var messageEntitys []*sdk_struct.MessageEntity
	s := sdk_struct.MsgStruct{}
	err := json.Unmarshal([]byte(messageEntityList), &messageEntitys)
	if err != nil {
		log.Error("internal", "messages unmarshal err", err.Error())
		return ""
	}
	c.initBasicInfo(&s, constant.UserMsgType, constant.AdvancedText, operationID)
	s.MessageEntityElem.Text = text
	s.MessageEntityElem.MessageEntityList = messageEntitys
	s.Content = utils.StructToJsonString(s.MessageEntityElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateTextAtMessage(text, atUserList, atUsersInfo, message, operationID string) string {
	var usersInfo []*sdk_struct.AtInfo
	var userIDList []string
	if text == "" {
		return ""
	}
	_ = json.Unmarshal([]byte(atUsersInfo), &usersInfo)
	_ = json.Unmarshal([]byte(atUserList), &userIDList)
	s, qs := sdk_struct.MsgStruct{}, sdk_struct.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	c.initBasicInfo(&s, constant.UserMsgType, constant.AtText, operationID)
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.AtElem.Text = text
	s.AtElem.AtUserList = userIDList
	s.AtElem.AtUsersInfo = usersInfo
	s.AtElem.QuoteMessage = &qs
	if message == "" {
		s.AtElem.QuoteMessage = nil
	}
	s.Content = utils.StructToJsonString(s.AtElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateLocationMessage(description string, longitude, latitude float64, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Location, operationID)
	s.LocationElem.Description = description
	s.LocationElem.Longitude = longitude
	s.LocationElem.Latitude = latitude
	s.Content = utils.StructToJsonString(s.LocationElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateCustomMessage(data, extension string, description, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Custom, operationID)
	s.CustomElem.Data = data
	s.CustomElem.Extension = extension
	s.CustomElem.Description = description
	s.Content = utils.StructToJsonString(s.CustomElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateQuoteMessage(text string, message, operationID string) string {
	s, qs := sdk_struct.MsgStruct{}, sdk_struct.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	c.initBasicInfo(&s, constant.UserMsgType, constant.Quote, operationID)
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
func (c *Conversation) CreateAdvancedQuoteMessage(text string, message, messageEntityList, operationID string) string {
	var messageEntities []*sdk_struct.MessageEntity
	s, qs := sdk_struct.MsgStruct{}, sdk_struct.MsgStruct{}
	_ = json.Unmarshal([]byte(message), &qs)
	_ = json.Unmarshal([]byte(messageEntityList), &messageEntities)
	c.initBasicInfo(&s, constant.UserMsgType, constant.Quote, operationID)
	//Avoid nested references
	if qs.ContentType == constant.Quote {
		qs.Content = qs.QuoteElem.Text
		qs.ContentType = constant.Text
	}
	s.QuoteElem.Text = text
	s.QuoteElem.MessageEntityList = messageEntities
	s.QuoteElem.QuoteMessage = &qs
	s.Content = utils.StructToJsonString(s.QuoteElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateCardMessage(cardInfo, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Card, operationID)
	s.Content = cardInfo
	return utils.StructToJsonString(s)

}
func (c *Conversation) CreateVideoMessageFromFullPath(videoFullPath string, videoType string, duration int64, snapshotFullPath, operationID string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(videoFullPath, c.DataDir) //a->b
		s, err := utils.CopyFile(videoFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, videoFullPath)
		}
		log.Info("internal", "videoFullPath dstFile", videoFullPath, dstFile, s)
		dstFile = utils.FileTmpPath(snapshotFullPath, c.DataDir) //a->b
		s, err = utils.CopyFile(snapshotFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, snapshotFullPath)
		}
		log.Info("internal", "snapshotFullPath dstFile", snapshotFullPath, dstFile, s)
		wg.Done()
	}()

	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video, operationID)
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
func (c *Conversation) CreateFileMessageFromFullPath(fileFullPath string, fileName, operationID string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(fileFullPath, c.DataDir)
		_, err := utils.CopyFile(fileFullPath, dstFile)
		log.Info(operationID, "copy file, ", fileFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err.Error(), fileFullPath)
		}
		wg.Done()
	}()
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.File, operationID)
	s.FileElem.FilePath = fileFullPath
	fi, err := os.Stat(fileFullPath)
	if err != nil {
		log.Error("internal", "get file Attributes error", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.FileElem.FileName = fileName
	s.Content = utils.StructToJsonString(s.FileElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateImageMessageFromFullPath(imageFullPath, operationID string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(imageFullPath, c.DataDir) //a->b
		_, err := utils.CopyFile(imageFullPath, dstFile)
		log.Info("internal", "copy file, ", imageFullPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, imageFullPath)
		}
		wg.Done()
	}()

	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture, operationID)
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
func (c *Conversation) CreateSoundMessageFromFullPath(soundPath string, duration int64, operationID string) string {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dstFile := utils.FileTmpPath(soundPath, c.DataDir) //a->b
		_, err := utils.CopyFile(soundPath, dstFile)
		log.Info("internal", "copy file, ", soundPath, dstFile)
		if err != nil {
			log.Error("internal", "open file failed: ", err, soundPath)
		}
		wg.Done()
	}()
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice, operationID)
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
func (c *Conversation) CreateImageMessage(imagePath, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture, operationID)
	s.PictureElem.SourcePath = c.DataDir + imagePath
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
func (c *Conversation) CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture, operationID string) string {
	s := sdk_struct.MsgStruct{}
	var p sdk_struct.PictureBaseInfo
	_ = json.Unmarshal([]byte(sourcePicture), &p)
	s.PictureElem.SourcePicture = p
	_ = json.Unmarshal([]byte(bigPicture), &p)
	s.PictureElem.BigPicture = p
	_ = json.Unmarshal([]byte(snapshotPicture), &p)
	s.PictureElem.SnapshotPicture = p
	c.initBasicInfo(&s, constant.UserMsgType, constant.Picture, operationID)
	s.Content = utils.StructToJsonString(s.PictureElem)
	return utils.StructToJsonString(s)
}

func msgStructToLocalChatLog(dst *model_struct.LocalChatLog, src *sdk_struct.MsgStruct) {
	copier.Copy(dst, src)
	if src.SessionType == constant.GroupChatType || src.SessionType == constant.SuperGroupChatType {
		dst.RecvID = src.GroupID
	}
}
func localChatLogToMsgStruct(dst *sdk_struct.NewMsgList, src []*model_struct.LocalChatLog) {
	copier.Copy(dst, &src)

}
func (c *Conversation) checkErrAndUpdateMessage(callback open_im_sdk_callback.SendMsgCallBack, errCode int32, err error, s *sdk_struct.MsgStruct, lc *model_struct.LocalConversation, operationID string) {
	if err != nil {
		if callback != nil {
			c.updateMsgStatusAndTriggerConversation(s.ClientMsgID, "", s.CreateTime, constant.MsgStatusSendFailed, s, lc, operationID)
			errInfo := "operationID[" + operationID + "], " + "info[" + err.Error() + "]" + s.ClientMsgID + " " + s.ServerMsgID
			log.NewError(operationID, "checkErr ", errInfo)
			callback.OnError(errCode, errInfo)
			runtime.Goexit()
		}
	}
}
func (c *Conversation) updateMsgStatusAndTriggerConversation(clientMsgID, serverMsgID string, sendTime int64, status int32, s *sdk_struct.MsgStruct, lc *model_struct.LocalConversation, operationID string) {
	//log.NewDebug(operationID, "this is test send message ", sendTime, status, clientMsgID, serverMsgID)
	s.SendTime = sendTime
	s.Status = status
	s.ServerMsgID = serverMsgID
	err := c.db.UpdateMessageTimeAndStatusController(s)
	if err != nil {
		log.Error(operationID, "send message update message status error", sendTime, status, clientMsgID, serverMsgID, err.Error())
	}
	lc.LatestMsg = utils.StructToJsonString(s)
	lc.LatestMsgSendTime = sendTime
	log.Info(operationID, "2 send message come here", *lc)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: lc.ConversationID, Action: constant.AddConOrUpLatMsg, Args: *lc}, c.GetCh())

}
func (c *Conversation) SendMessage(callback open_im_sdk_callback.SendMsgCallBack, message, recvID, groupID string, offlinePushInfo string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.Debug(operationID, "SendMessage start ")
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		s.SendID = c.loginUserID
		s.SenderPlatformID = c.platformID
		p := &server_api_params.OfflinePushInfo{}
		if offlinePushInfo == "" {
			p = nil
		} else {
			common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
		}
		if recvID == "" && groupID == "" {
			common.CheckAnyErrCallback(callback, 201, errors.New("recvID && groupID not both null"), operationID)
		}
		var localMessage model_struct.LocalChatLog
		var conversationID string
		options := make(map[string]bool, 2)
		lc := &model_struct.LocalConversation{LatestMsgSendTime: s.CreateTime}
		//根据单聊群聊类型组装消息和会话
		if recvID == "" {
			g, err := c.full.GetGroupInfoByGroupID(groupID)
			common.CheckAnyErrCallback(callback, 202, err, operationID)
			lc.ShowName = g.GroupName
			lc.FaceURL = g.FaceURL
			switch g.GroupType {
			case constant.NormalGroup:
				s.SessionType = constant.GroupChatType
				lc.ConversationType = constant.GroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			case constant.SuperGroup, constant.WorkingGroup:
				s.SessionType = constant.SuperGroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
				lc.ConversationType = constant.SuperGroupChatType
			}
			s.GroupID = groupID
			lc.GroupID = groupID
			gm, err := c.db.GetGroupMemberInfoByGroupIDUserID(groupID, c.loginUserID)
			if err == nil && gm != nil {
				log.Debug(operationID, "group chat test", *gm)
				if gm.Nickname != "" {
					s.SenderNickname = gm.Nickname
				}
			}
			//groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			//common.CheckAnyErrCallback(callback, 202, err, operationID)
			//if !utils.IsContain(s.SendID, groupMemberUidList) {
			//	common.CheckAnyErrCallback(callback, 208, errors.New("you not exist in this group"), operationID)
			//}
			s.AttachedInfoElem.GroupHasReadInfo.GroupMemberCount = g.MemberCount
			s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
		} else {
			log.Debug(operationID, "send msg single chat come here")
			s.SessionType = constant.SingleChatType
			s.RecvID = recvID
			conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			lc.UserID = recvID
			lc.ConversationType = constant.SingleChatType
			//faceUrl, name, err := c.friend.GetUserNameAndFaceUrlByUid(recvID, operationID)
			oldLc, err := c.db.GetConversation(conversationID)
			if err == nil && oldLc.IsPrivateChat {
				options[constant.IsNotPrivate] = false
				s.AttachedInfoElem.IsPrivateChat = true
				s.AttachedInfoElem.BurnDuration = oldLc.BurnDuration
				s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
			}
			if err != nil {
				t := time.Now()
				faceUrl, name, err := c.cache.GetUserNameAndFaceURL(recvID, operationID)
				log.Debug(operationID, "GetUserNameAndFaceURL cost time:", time.Since(t))
				common.CheckAnyErrCallback(callback, 301, err, operationID)
				lc.FaceURL = faceUrl
				lc.ShowName = name
			}

		}
		t := time.Now()
		log.Debug(operationID, "before insert  message is ", s)
		oldMessage, err := c.db.GetMessageController(&s)
		log.Debug(operationID, "GetMessageController cost time:", time.Since(t), err)
		if err != nil {
			msgStructToLocalChatLog(&localMessage, &s)
			err := c.db.InsertMessageController(&localMessage)
			common.CheckAnyErrCallback(callback, 201, err, operationID)
		} else {
			if oldMessage.Status != constant.MsgStatusSendFailed {
				common.CheckAnyErrCallback(callback, 202, errors.New("only failed message can be repeatedly send"), operationID)
			} else {
				s.Status = constant.MsgStatusSending
			}
		}
		lc.ConversationID = conversationID
		lc.LatestMsg = utils.StructToJsonString(s)
		log.Info(operationID, "send message come here", *lc)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.AddConOrUpLatMsg, Args: *lc}, c.GetCh())
		var delFile []string
		//media file handle
		if s.Status != constant.MsgStatusSendSuccess { //filter forward message
			switch s.ContentType {
			case constant.Picture:
				var sourcePath string
				if utils.FileExist(s.PictureElem.SourcePath) {
					sourcePath = s.PictureElem.SourcePath
					delFile = append(delFile, utils.FileTmpPath(s.PictureElem.SourcePath, c.DataDir))
				} else {
					sourcePath = utils.FileTmpPath(s.PictureElem.SourcePath, c.DataDir)
					delFile = append(delFile, sourcePath)
				}
				log.Info(operationID, "file", sourcePath, delFile)
				sourceUrl, uuid, err := c.UploadImage(sourcePath, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
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
					delFile = append(delFile, utils.FileTmpPath(s.SoundElem.SoundPath, c.DataDir))
				} else {
					sourcePath = utils.FileTmpPath(s.SoundElem.SoundPath, c.DataDir)
					delFile = append(delFile, sourcePath)
				}
				log.Info(operationID, "file", sourcePath, delFile)
				soundURL, uuid, err := c.UploadSound(sourcePath, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = utils.StructToJsonString(s.SoundElem)

			case constant.Video:
				var videoPath string
				var snapPath string
				if utils.FileExist(s.VideoElem.VideoPath) {
					videoPath = s.VideoElem.VideoPath
					snapPath = s.VideoElem.SnapshotPath
					delFile = append(delFile, utils.FileTmpPath(s.VideoElem.VideoPath, c.DataDir))
					delFile = append(delFile, utils.FileTmpPath(s.VideoElem.SnapshotPath, c.DataDir))
				} else {
					videoPath = utils.FileTmpPath(s.VideoElem.VideoPath, c.DataDir)
					snapPath = utils.FileTmpPath(s.VideoElem.SnapshotPath, c.DataDir)
					delFile = append(delFile, videoPath)
					delFile = append(delFile, snapPath)
				}
				log.Info(operationID, "file: ", videoPath, snapPath, delFile)
				snapshotURL, snapshotUUID, videoURL, videoUUID, err := c.UploadVideo(videoPath, snapPath, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.VideoElem.VideoURL = videoURL
				s.VideoElem.SnapshotUUID = snapshotUUID
				s.VideoElem.SnapshotURL = snapshotURL
				s.VideoElem.VideoUUID = videoUUID
				s.Content = utils.StructToJsonString(s.VideoElem)
			case constant.File:
				fileURL, fileUUID, err := c.UploadFile(s.FileElem.FilePath, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
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
			case constant.Face:
			case constant.AdvancedText:
			default:
				common.CheckAnyErrCallback(callback, 202, errors.New("contentType not currently supported"+utils.Int32ToString(s.ContentType)), operationID)
			}
			oldMessage, err := c.db.GetMessageController(&s)
			if err != nil {
				log.Warn(operationID, "get message err")
			} else {
				log.Debug(operationID, "before update database message is ", *oldMessage)
			}
			if utils.IsContainInt(int(s.ContentType), []int{constant.Picture, constant.Voice, constant.Video, constant.File}) {
				msgStructToLocalChatLog(&localMessage, &s)
				log.Warn(operationID, "update message is ", s, localMessage)
				err = c.db.UpdateMessageController(&localMessage)
				common.CheckAnyErrCallback(callback, 201, err, operationID)
			}
		}
		c.sendMessageToServer(&s, lc, callback, delFile, p, options, operationID)
	}()

}
func (c *Conversation) SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, message, recvID, groupID string, offlinePushInfo string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		s.SendID = c.loginUserID
		s.SenderPlatformID = c.platformID
		p := &server_api_params.OfflinePushInfo{}
		if offlinePushInfo == "" {
			p = nil
		} else {
			common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
		}
		if recvID == "" && groupID == "" {
			common.CheckAnyErrCallback(callback, 201, errors.New("recvID && groupID not both null"), operationID)
		}
		var localMessage model_struct.LocalChatLog
		var conversationID string
		options := make(map[string]bool, 2)
		lc := &model_struct.LocalConversation{LatestMsgSendTime: s.CreateTime}
		//根据单聊群聊类型组装消息和会话
		if recvID == "" {
			g, err := c.full.GetGroupInfoByGroupID(groupID)
			common.CheckAnyErrCallback(callback, 202, err, operationID)
			lc.ShowName = g.GroupName
			lc.FaceURL = g.FaceURL
			switch g.GroupType {
			case constant.NormalGroup:
				s.SessionType = constant.GroupChatType
				lc.ConversationType = constant.GroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			case constant.SuperGroup, constant.WorkingGroup:
				s.SessionType = constant.SuperGroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
				lc.ConversationType = constant.SuperGroupChatType
			}
			s.GroupID = groupID
			lc.GroupID = groupID
			gm, err := c.db.GetGroupMemberInfoByGroupIDUserID(groupID, c.loginUserID)
			if err == nil && gm != nil {
				log.Debug(operationID, "group chat test", *gm)
				if gm.Nickname != "" {
					s.SenderNickname = gm.Nickname
				}
			}
			//groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			//common.CheckAnyErrCallback(callback, 202, err, operationID)
			//if !utils.IsContain(s.SendID, groupMemberUidList) {
			//	common.CheckAnyErrCallback(callback, 208, errors.New("you not exist in this group"), operationID)
			//}
			s.AttachedInfoElem.GroupHasReadInfo.GroupMemberCount = g.MemberCount
			s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
		} else {
			s.SessionType = constant.SingleChatType
			s.RecvID = recvID
			conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			lc.UserID = recvID
			lc.ConversationType = constant.SingleChatType
			//faceUrl, name, err := c.friend.GetUserNameAndFaceUrlByUid(recvID, operationID)
			oldLc, err := c.db.GetConversation(conversationID)
			if err == nil && oldLc.IsPrivateChat {
				options[constant.IsNotPrivate] = false
				s.AttachedInfoElem.IsPrivateChat = true
				s.AttachedInfoElem.BurnDuration = oldLc.BurnDuration
				s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
			}
			if err != nil {
				faceUrl, name, err := c.cache.GetUserNameAndFaceURL(recvID, operationID)
				common.CheckAnyErrCallback(callback, 301, err, operationID)
				lc.FaceURL = faceUrl
				lc.ShowName = name
			}

		}
		oldMessage, err := c.db.GetMessageController(&s)
		if err != nil {
			msgStructToLocalChatLog(&localMessage, &s)
			err := c.db.InsertMessageController(&localMessage)
			common.CheckAnyErrCallback(callback, 201, err, operationID)
		} else {
			if oldMessage.Status != constant.MsgStatusSendFailed {
				common.CheckAnyErrCallback(callback, 202, errors.New("only failed message can be repeatedly send"), operationID)
			} else {
				s.Status = constant.MsgStatusSending
			}
		}
		lc.ConversationID = conversationID
		lc.LatestMsg = utils.StructToJsonString(s)
		//u.doUpdateConversation(common.cmd2Value{Value: common.updateConNode{conversationID, constant.AddConOrUpLatMsg,
		//c}})
		//u.doUpdateConversation(cmd2Value{Value: updateConNode{"", ConChange, []string{conversationID}}})
		//_ = u.triggerCmdUpdateConversation(updateConNode{conversationID, ConChange, ""})
		var delFile []string
		if utils.IsContainInt(int(s.ContentType), []int{constant.Picture, constant.Voice, constant.Video, constant.File}) {
			msgStructToLocalChatLog(&localMessage, &s)
			err = c.db.UpdateMessageController(&localMessage)
			common.CheckAnyErrCallback(callback, 201, err, operationID)
		}
		c.sendMessageToServer(&s, lc, callback, delFile, p, options, operationID)
	}()
}
func (c *Conversation) SendMessageByBuffer(callback open_im_sdk_callback.SendMsgCallBack, message, recvID, groupID string, offlinePushInfo string, operationID string, buffer1, buffer2 *bytes.Buffer) {
	if callback == nil {
		return
	}
	go func() {
		log.Debug(operationID, "SendMessageByBuffer start ")
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		s.SendID = c.loginUserID
		s.SenderPlatformID = c.platformID
		p := &server_api_params.OfflinePushInfo{}
		if offlinePushInfo == "" {
			p = nil
		} else {
			common.JsonUnmarshalAndArgsValidate(offlinePushInfo, &p, callback, operationID)
		}
		if recvID == "" && groupID == "" {
			common.CheckAnyErrCallback(callback, 201, errors.New("recvID && groupID not both null"), operationID)
		}
		var localMessage model_struct.LocalChatLog
		var conversationID string
		options := make(map[string]bool, 2)
		lc := &model_struct.LocalConversation{LatestMsgSendTime: s.CreateTime}
		//根据单聊群聊类型组装消息和会话
		if recvID == "" {
			g, err := c.full.GetGroupInfoByGroupID(groupID)
			common.CheckAnyErrCallback(callback, 202, err, operationID)
			lc.ShowName = g.GroupName
			lc.FaceURL = g.FaceURL
			switch g.GroupType {
			case constant.NormalGroup:
				s.SessionType = constant.GroupChatType
				lc.ConversationType = constant.GroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
			case constant.SuperGroup, constant.WorkingGroup:
				s.SessionType = constant.SuperGroupChatType
				conversationID = utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType)
				lc.ConversationType = constant.SuperGroupChatType
			}
			s.GroupID = groupID
			lc.GroupID = groupID
			gm, err := c.db.GetGroupMemberInfoByGroupIDUserID(groupID, c.loginUserID)
			if err == nil && gm != nil {
				log.Debug(operationID, "group chat test", *gm)
				if gm.Nickname != "" {
					s.SenderNickname = gm.Nickname
				}
			}
			//groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			//common.CheckAnyErrCallback(callback, 202, err, operationID)
			//if !utils.IsContain(s.SendID, groupMemberUidList) {
			//	common.CheckAnyErrCallback(callback, 208, errors.New("you not exist in this group"), operationID)
			//}
			s.AttachedInfoElem.GroupHasReadInfo.GroupMemberCount = g.MemberCount
			s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
		} else {
			log.Debug(operationID, "send msg single chat come here")
			s.SessionType = constant.SingleChatType
			s.RecvID = recvID
			conversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			lc.UserID = recvID
			lc.ConversationType = constant.SingleChatType
			//faceUrl, name, err := c.friend.GetUserNameAndFaceUrlByUid(recvID, operationID)
			oldLc, err := c.db.GetConversation(conversationID)
			if err == nil && oldLc.IsPrivateChat {
				options[constant.IsNotPrivate] = false
				s.AttachedInfoElem.IsPrivateChat = true
				s.AttachedInfo = utils.StructToJsonString(s.AttachedInfoElem)
			}
			if err != nil {
				t := time.Now()
				faceUrl, name, err := c.cache.GetUserNameAndFaceURL(recvID, operationID)
				log.Debug(operationID, "GetUserNameAndFaceURL cost time:", time.Since(t))
				common.CheckAnyErrCallback(callback, 301, err, operationID)
				lc.FaceURL = faceUrl
				lc.ShowName = name
			}

		}
		t := time.Now()
		log.Debug(operationID, "before insert  message is ", s)
		oldMessage, err := c.db.GetMessageController(&s)
		log.Debug(operationID, "GetMessageController cost time:", time.Since(t), err)
		if err != nil {
			msgStructToLocalChatLog(&localMessage, &s)
			err := c.db.InsertMessageController(&localMessage)
			common.CheckAnyErrCallback(callback, 201, err, operationID)
		} else {
			if oldMessage.Status != constant.MsgStatusSendFailed {
				common.CheckAnyErrCallback(callback, 202, errors.New("only failed message can be repeatedly send"), operationID)
			} else {
				s.Status = constant.MsgStatusSending
			}
		}
		lc.ConversationID = conversationID
		lc.LatestMsg = utils.StructToJsonString(s)
		log.Info(operationID, "send message come here", *lc)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.AddConOrUpLatMsg, Args: *lc}, c.GetCh())
		var delFile []string
		//media file handle
		if s.Status != constant.MsgStatusSendSuccess { //filter forward message
			switch s.ContentType {
			case constant.Picture:
				sourceUrl, uuid, err := c.UploadImageByBuffer(buffer1, s.PictureElem.SourcePicture.Size, s.PictureElem.SourcePicture.Type, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.PictureElem.SourcePicture.Url = sourceUrl
				s.PictureElem.SourcePicture.UUID = uuid
				s.PictureElem.SnapshotPicture.Url = sourceUrl + "?imageView2/2/w/" + constant.ZoomScale + "/h/" + constant.ZoomScale
				s.PictureElem.SnapshotPicture.Width = int32(utils.StringToInt(constant.ZoomScale))
				s.PictureElem.SnapshotPicture.Height = int32(utils.StringToInt(constant.ZoomScale))
				s.Content = utils.StructToJsonString(s.PictureElem)

			case constant.Voice:
				soundURL, uuid, err := c.UploadSoundByBuffer(buffer1, s.SoundElem.DataSize, "sound", callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = utils.StructToJsonString(s.SoundElem)

			case constant.Video:

				snapshotURL, snapshotUUID, videoURL, videoUUID, err := c.UploadVideoByBuffer(buffer1, buffer2, s.VideoElem.VideoSize,
					s.VideoElem.SnapshotSize, s.VideoElem.VideoType, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.VideoElem.VideoURL = videoURL
				s.VideoElem.SnapshotUUID = snapshotUUID
				s.VideoElem.SnapshotURL = snapshotURL
				s.VideoElem.VideoUUID = videoUUID
				s.Content = utils.StructToJsonString(s.VideoElem)
			case constant.File:
				fileURL, fileUUID, err := c.UploadFileByBuffer(buffer1, s.FileElem.FileSize, "file", callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
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
			case constant.Face:
			case constant.AdvancedText:
			default:
				common.CheckAnyErrCallback(callback, 202, errors.New("contentType not currently supported"+utils.Int32ToString(s.ContentType)), operationID)
			}
			oldMessage, err := c.db.GetMessageController(&s)
			if err != nil {
				log.Warn(operationID, "get message err")
			} else {
				log.Debug(operationID, "before update database message is ", *oldMessage)
			}
			if utils.IsContainInt(int(s.ContentType), []int{constant.Picture, constant.Voice, constant.Video, constant.File}) {
				msgStructToLocalChatLog(&localMessage, &s)
				log.Warn(operationID, "update message is ", s, localMessage)
				err = c.db.UpdateMessageController(&localMessage)
				common.CheckAnyErrCallback(callback, 201, err, operationID)
			}
		}
		c.sendMessageToServer(&s, lc, callback, delFile, p, options, operationID)
	}()
}

func (c *Conversation) InternalSendMessage(callback open_im_sdk_callback.Base, s *sdk_struct.MsgStruct, recvID, groupID, operationID string, p *server_api_params.OfflinePushInfo, onlineUserOnly bool, options map[string]bool) (*server_api_params.UserSendMsgResp, error) {
	if recvID == "" && groupID == "" {
		common.CheckAnyErrCallback(callback, 201, errors.New("recvID && groupID not both null"), operationID)
	}
	if recvID == "" {
		g, err := c.full.GetGroupInfoByGroupID(groupID)
		common.CheckAnyErrCallback(callback, 202, err, operationID)
		switch g.GroupType {
		case constant.NormalGroup:
			s.SessionType = constant.GroupChatType
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			common.CheckAnyErrCallback(callback, 202, err, operationID)
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				common.CheckAnyErrCallback(callback, 208, errors.New("you not exist in this group"), operationID)
			}

		case constant.SuperGroup:
			s.SessionType = constant.SuperGroupChatType
		case constant.WorkingGroup:
			s.SessionType = constant.SuperGroupChatType
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(groupID)
			common.CheckAnyErrCallback(callback, 202, err, operationID)
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				common.CheckAnyErrCallback(callback, 208, errors.New("you not exist in this group"), operationID)
			}
		}
		s.GroupID = groupID

	} else {
		s.SessionType = constant.SingleChatType
		s.RecvID = recvID
	}

	if onlineUserOnly {
		options[constant.IsHistory] = false
		options[constant.IsPersistent] = false
		options[constant.IsOfflinePush] = false
		options[constant.IsSenderSync] = false
	}

	var wsMsgData server_api_params.MsgData
	copier.Copy(&wsMsgData, s)
	wsMsgData.Content = []byte(s.Content)
	wsMsgData.CreateTime = s.CreateTime
	wsMsgData.Options = options
	wsMsgData.OfflinePushInfo = p
	timeout := 10
	retryTimes := 0
	g, err := c.SendReqWaitResp(&wsMsgData, constant.WSSendMsg, timeout, retryTimes, c.loginUserID, operationID)
	switch e := err.(type) {
	case *constant.ErrInfo:
		common.CheckAnyErrCallback(callback, e.ErrCode, e, operationID)
	default:
		common.CheckAnyErrCallback(callback, 301, err, operationID)
	}
	var sendMsgResp server_api_params.UserSendMsgResp
	_ = proto.Unmarshal(g.Data, &sendMsgResp)
	return &sendMsgResp, nil

}

func (c *Conversation) sendMessageToServer(s *sdk_struct.MsgStruct, lc *model_struct.LocalConversation, callback open_im_sdk_callback.SendMsgCallBack,
	delFile []string, offlinePushInfo *server_api_params.OfflinePushInfo, options map[string]bool, operationID string) {
	log.Debug(operationID, "sendMessageToServer ", s.ServerMsgID, " ", s.ClientMsgID)
	//Protocol conversion
	var wsMsgData server_api_params.MsgData
	copier.Copy(&wsMsgData, s)
	if wsMsgData.ContentType == constant.Text && c.encryptionKey != "" {
		ciphertext, err := utils.AesEncrypt([]byte(s.Content), []byte(c.encryptionKey))
		c.checkErrAndUpdateMessage(callback, 302, err, s, lc, operationID)
		attachInfo := sdk_struct.AttachedInfoElem{}
		_ = utils.JsonStringToStruct(s.AttachedInfo, &attachInfo)
		attachInfo.IsEncryption = true
		attachInfo.InEncryptStatus = true
		wsMsgData.Content = ciphertext
		wsMsgData.AttachedInfo = utils.StructToJsonString(attachInfo)
	} else {
		wsMsgData.Content = []byte(s.Content)
	}
	wsMsgData.CreateTime = s.CreateTime
	wsMsgData.Options = options
	wsMsgData.AtUserIDList = s.AtElem.AtUserList
	wsMsgData.OfflinePushInfo = offlinePushInfo
	timeout := 300
	retryTimes := 60
	resp, err := c.SendReqWaitResp(&wsMsgData, constant.WSSendMsg, timeout, retryTimes, c.loginUserID, operationID)
	switch e := err.(type) {
	case *constant.ErrInfo:
		c.checkErrAndUpdateMessage(callback, e.ErrCode, e, s, lc, operationID)
	default:
		c.checkErrAndUpdateMessage(callback, 302, err, s, lc, operationID)
	}
	var sendMsgResp server_api_params.UserSendMsgResp
	_ = proto.Unmarshal(resp.Data, &sendMsgResp)
	s.SendTime = sendMsgResp.SendTime
	s.Status = constant.MsgStatusSendSuccess
	s.ServerMsgID = sendMsgResp.ServerMsgID
	callback.OnProgress(100)
	callback.OnSuccess(utils.StructToJsonString(s))
	log.Debug(operationID, "callback OnSuccess", s.ClientMsgID, s.ServerMsgID)
	//remove media cache file
	for _, v := range delFile {
		err := os.Remove(v)
		if err != nil {
			log.Error(operationID, "remove failed,", err.Error(), v)
		}
		log.Debug(operationID, "remove file: ", v)
	}
	c.updateMsgStatusAndTriggerConversation(sendMsgResp.ClientMsgID, sendMsgResp.ServerMsgID, sendMsgResp.SendTime, constant.MsgStatusSendSuccess, s, lc, operationID)

}

func (c *Conversation) CreateSoundMessageByURL(soundBaseInfo, operationID string) string {
	s := sdk_struct.MsgStruct{}
	var soundElem sdk_struct.SoundBaseInfo
	_ = json.Unmarshal([]byte(soundBaseInfo), &soundElem)
	s.SoundElem = soundElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice, operationID)
	s.Content = utils.StructToJsonString(s.SoundElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateSoundMessage(soundPath string, duration int64, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Voice, operationID)
	s.SoundElem.SoundPath = c.DataDir + soundPath
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
func (c *Conversation) CreateVideoMessageByURL(videoBaseInfo, operationID string) string {
	s := sdk_struct.MsgStruct{}
	var videoElem sdk_struct.VideoBaseInfo
	_ = json.Unmarshal([]byte(videoBaseInfo), &videoElem)
	s.VideoElem = videoElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video, operationID)
	s.Content = utils.StructToJsonString(s.VideoElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateVideoMessage(videoPath string, videoType string, duration int64, snapshotPath, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Video, operationID)
	s.VideoElem.VideoPath = c.DataDir + videoPath
	s.VideoElem.VideoType = videoType
	s.VideoElem.Duration = duration
	if snapshotPath == "" {
		s.VideoElem.SnapshotPath = ""
	} else {
		s.VideoElem.SnapshotPath = c.DataDir + snapshotPath
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
func (c *Conversation) CreateFileMessageByURL(fileBaseInfo, operationID string) string {
	s := sdk_struct.MsgStruct{}
	var fileElem sdk_struct.FileBaseInfo
	_ = json.Unmarshal([]byte(fileBaseInfo), &fileElem)
	s.FileElem = fileElem
	c.initBasicInfo(&s, constant.UserMsgType, constant.File, operationID)
	s.Content = utils.StructToJsonString(s.FileElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateFileMessage(filePath string, fileName, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.File, operationID)
	s.FileElem.FilePath = c.DataDir + filePath
	s.FileElem.FileName = fileName
	fi, err := os.Stat(s.FileElem.FilePath)
	if err != nil {
		log.Error("internal", "get file message err", err.Error())
		return ""
	}
	s.FileElem.FileSize = fi.Size()
	s.Content = utils.StructToJsonString(s.FileElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateMergerMessage(messageList, title, summaryList, operationID string) string {
	var messages []*sdk_struct.MsgStruct
	var summaries []string
	s := sdk_struct.MsgStruct{}
	err := json.Unmarshal([]byte(messageList), &messages)
	if err != nil {
		log.Error("internal", "messages Unmarshal err", err.Error())
		return ""
	}
	_ = json.Unmarshal([]byte(summaryList), &summaries)
	c.initBasicInfo(&s, constant.UserMsgType, constant.Merger, operationID)
	s.MergeElem.AbstractList = summaries
	s.MergeElem.Title = title
	s.MergeElem.MultiMessage = messages
	s.Content = utils.StructToJsonString(s.MergeElem)
	return utils.StructToJsonString(s)
}
func (c *Conversation) CreateFaceMessage(index int, data, operationID string) string {
	s := sdk_struct.MsgStruct{}
	c.initBasicInfo(&s, constant.UserMsgType, constant.Face, operationID)
	s.FaceElem.Data = data
	s.FaceElem.Index = index
	s.Content = utils.StructToJsonString(s.FaceElem)
	return utils.StructToJsonString(s)

}
func (c *Conversation) CreateForwardMessage(m, operationID string) string {
	s := sdk_struct.MsgStruct{}
	err := json.Unmarshal([]byte(m), &s)
	if err != nil {
		log.Error("internal", "messages Unmarshal err", err.Error())
		return ""
	}
	if s.Status != constant.MsgStatusSendSuccess {
		log.Error("internal", "only send success message can be Forward")
		return ""
	}
	c.initBasicInfo(&s, constant.UserMsgType, s.ContentType, operationID)
	//Forward message seq is set to 0
	s.Seq = 0
	s.Status = constant.MsgStatusSendSuccess
	return utils.StructToJsonString(s)
}
func (c *Conversation) FindMessageList(callback open_im_sdk_callback.Base, findMessageOptions, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		t := time.Now()
		log.NewInfo(operationID, "FindMessageList args: ", findMessageOptions)
		var unmarshalParams sdk_params_callback.FindMessageListParams
		common.JsonUnmarshalCallback(findMessageOptions, &unmarshalParams, callback, operationID)
		result := c.findMessageList(unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "FindMessageList callback: ", utils.StructToJsonStringDefault(result), "cost time", time.Since(t))
	}()
}
func (c *Conversation) GetHistoryMessageList(callback open_im_sdk_callback.Base, getMessageOptions, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		t := time.Now()
		log.NewInfo(operationID, "GetHistoryMessageList args: ", getMessageOptions)
		var unmarshalParams sdk_params_callback.GetHistoryMessageListParams
		common.JsonUnmarshalCallback(getMessageOptions, &unmarshalParams, callback, operationID)
		result := c.getHistoryMessageList(callback, unmarshalParams, operationID, false)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "cost time", time.Since(t), "GetHistoryMessageList callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetAdvancedHistoryMessageList(callback open_im_sdk_callback.Base, getMessageOptions, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		t := time.Now()
		log.NewInfo(operationID, "GetAdvancedHistoryMessageList args: ", getMessageOptions)
		var unmarshalParams sdk_params_callback.GetAdvancedHistoryMessageListParams
		common.JsonUnmarshalCallback(getMessageOptions, &unmarshalParams, callback, operationID)
		result := c.getAdvancedHistoryMessageList(callback, unmarshalParams, operationID, false)
		if len(result.MessageList) == 0 {
			s := make([]*sdk_struct.MsgStruct, 0)
			result.MessageList = s
		}
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.Error(operationID, "length:", len(result.MessageList), "cost time", time.Since(t), "GetAdvancedHistoryMessageList callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) GetAdvancedHistoryMessageListReverse(callback open_im_sdk_callback.Base, getMessageOptions, operationID string) {
	var unmarshalParams sdk_params_callback.GetAdvancedHistoryMessageListParams
	common.JsonUnmarshalCallback(getMessageOptions, &unmarshalParams, callback, operationID)
	result := c.getAdvancedHistoryMessageList(callback, unmarshalParams, operationID, true)
	callback.OnSuccess(utils.StructToJsonStringDefault(result))
	log.NewInfo(operationID, "GetAdvancedHistoryMessageListReverse callback: ", utils.StructToJsonStringDefault(result))
}

func (c *Conversation) GetHistoryMessageListReverse(callback open_im_sdk_callback.Base, getMessageOptions, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "GetHistoryMessageListReverse args: ", getMessageOptions)
		var unmarshalParams sdk_params_callback.GetHistoryMessageListParams
		common.JsonUnmarshalCallback(getMessageOptions, &unmarshalParams, callback, operationID)
		result := c.getHistoryMessageList(callback, unmarshalParams, operationID, true)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "GetHistoryMessageListReverse callback: ", utils.StructToJsonStringDefault(result))
	}()
}
func (c *Conversation) RevokeMessage(callback open_im_sdk_callback.Base, message string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "RevokeMessage args: ", message)
		var unmarshalParams sdk_params_callback.RevokeMessageParams
		common.JsonUnmarshalCallback(message, &unmarshalParams, callback, operationID)
		c.revokeOneMessage(callback, unmarshalParams, operationID)
		callback.OnSuccess(sdk_params_callback.RevokeMessageCallback)
		log.NewInfo(operationID, "RevokeMessage callback: ", sdk_params_callback.RevokeMessageCallback)
	}()
}
func (c *Conversation) NewRevokeMessage(callback open_im_sdk_callback.Base, message string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "RevokeMessage args: ", message)
		var unmarshalParams sdk_params_callback.RevokeMessageParams
		common.JsonUnmarshalCallback(message, &unmarshalParams, callback, operationID)
		c.newRevokeOneMessage(callback, unmarshalParams, operationID)
		callback.OnSuccess(sdk_params_callback.RevokeMessageCallback)
		log.NewInfo(operationID, "RevokeMessage callback: ", sdk_params_callback.RevokeMessageCallback)
	}()
}
func (c *Conversation) TypingStatusUpdate(callback open_im_sdk_callback.Base, recvID, msgTip, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "TypingStatusUpdate args: ", recvID, msgTip)
		c.typingStatusUpdate(callback, recvID, msgTip, operationID)
		callback.OnSuccess(sdk_params_callback.TypingStatusUpdateCallback)
		log.NewInfo(operationID, "TypingStatusUpdate callback: ", sdk_params_callback.TypingStatusUpdateCallback)
	}()
}

func (c *Conversation) MarkC2CMessageAsRead(callback open_im_sdk_callback.Base, userID string, msgIDList, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "MarkC2CMessageAsRead args: ", userID, msgIDList)
		var unmarshalParams sdk_params_callback.MarkC2CMessageAsReadParams
		common.JsonUnmarshalCallback(msgIDList, &unmarshalParams, callback, operationID)
		if len(unmarshalParams) == 0 {
			conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
			c.setOneConversationUnread(callback, conversationID, 0, operationID)
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
			callback.OnSuccess(sdk_params_callback.MarkC2CMessageAsReadCallback)
			return
		}
		c.markC2CMessageAsRead(callback, unmarshalParams, userID, operationID)
		callback.OnSuccess(sdk_params_callback.MarkC2CMessageAsReadCallback)
		log.NewInfo(operationID, "MarkC2CMessageAsRead callback: ", sdk_params_callback.MarkC2CMessageAsReadCallback)
	}()

}

func (c *Conversation) MarkMessageAsReadByConID(callback open_im_sdk_callback.Base, conversationID, msgIDList, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "MarkMessageAsReadByConID args: ", conversationID, msgIDList)
		var unmarshalParams sdk_params_callback.MarkMessageAsReadByConIDParams
		common.JsonUnmarshalCallback(msgIDList, &unmarshalParams, callback, operationID)
		if len(unmarshalParams) == 0 {
			c.setOneConversationUnread(callback, conversationID, 0, operationID)
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
			_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
			callback.OnSuccess(sdk_params_callback.MarkMessageAsReadByConIDCallback)
			return
		}
		//c.markMessageAsReadByConID(callback, unmarshalParams, conversationID, operationID)
		callback.OnSuccess(sdk_params_callback.MarkMessageAsReadByConIDCallback)
		log.NewInfo(operationID, "MarkMessageAsReadByConID callback: ", sdk_params_callback.MarkMessageAsReadByConIDCallback)
	}()

}

// fixme
func (c *Conversation) MarkAllConversationHasRead(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		var lc model_struct.LocalConversation
		lc.UnreadCount = 0
		err := c.db.UpdateAllConversation(&lc)
		common.CheckDBErrCallback(callback, err, operationID)
		callback.OnSuccess("")

	}()
}

// deprecated
func (c *Conversation) MarkGroupMessageHasRead(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "MarkGroupMessageHasRead args: ", groupID)
		conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
		c.setOneConversationUnread(callback, conversationID, 0, operationID)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		callback.OnSuccess(sdk_params_callback.MarkGroupMessageHasReadCallback)
	}()
}
func (c *Conversation) MarkGroupMessageAsRead(callback open_im_sdk_callback.Base, groupID string, msgIDList, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		log.NewInfo(operationID, "MarkGroupMessageAsRead args: ", groupID, msgIDList)
		var unmarshalParams sdk_params_callback.MarkGroupMessageAsReadParams
		common.JsonUnmarshalCallback(msgIDList, &unmarshalParams, callback, operationID)
		c.markGroupMessageAsRead(callback, unmarshalParams, groupID, operationID)
		callback.OnSuccess(sdk_params_callback.MarkGroupMessageAsReadCallback)
		log.NewInfo(operationID, "MarkGroupMessageAsRead callback: ", sdk_params_callback.MarkGroupMessageAsReadCallback)
	}()
}
func (c *Conversation) DeleteMessageFromLocalStorage(callback open_im_sdk_callback.Base, message string, operationID string) {
	go func() {
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		c.deleteMessageFromLocalStorage(callback, &s, operationID)
		callback.OnSuccess("")
	}()
}

func (c *Conversation) ClearC2CHistoryMessage(callback open_im_sdk_callback.Base, userID string, operationID string) {
	go func() {
		c.clearC2CHistoryMessage(callback, userID, operationID)
		callback.OnSuccess("")

	}()
}
func (c *Conversation) ClearGroupHistoryMessage(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	go func() {
		c.clearGroupHistoryMessage(callback, groupID, operationID)
		callback.OnSuccess("")

	}()
}
func (c *Conversation) ClearC2CHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, userID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userID)
		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
		c.deleteConversationAndMsgFromSvr(callback, conversationID, operationID)
		c.clearC2CHistoryMessage(callback, userID, operationID)
		callback.OnSuccess("")
	}()
}

// fixme
func (c *Conversation) ClearGroupHistoryMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, groupID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", groupID)
		conversationID, _, err := c.getConversationTypeByGroupID(groupID)
		common.CheckDBErrCallback(callback, err, operationID)
		c.deleteConversationAndMsgFromSvr(callback, conversationID, operationID)
		c.clearGroupHistoryMessage(callback, groupID, operationID)
		callback.OnSuccess("")
	}()
}

func (c *Conversation) InsertSingleMessageToLocalStorage(callback open_im_sdk_callback.Base, message, recvID, sendID, operationID string) {
	go func() {
		log.NewInfo(operationID, "InsertSingleMessageToLocalStorage args: ", message, recvID, sendID)
		if recvID == "" || sendID == "" {
			common.CheckAnyErrCallback(callback, 208, errors.New("recvID or sendID is null"), operationID)
		}
		var conversation model_struct.LocalConversation

		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		if sendID != c.loginUserID {
			faceUrl, name, err := c.cache.GetUserNameAndFaceURL(sendID, operationID)
			if err != nil {
				log.Error(operationID, "GetUserNameAndFaceURL err", err.Error(), sendID)
			}
			s.SenderFaceURL = faceUrl
			s.SenderNickname = name
			conversation.FaceURL = faceUrl
			conversation.ShowName = name
			conversation.UserID = sendID
			conversation.ConversationID = utils.GetConversationIDBySessionType(sendID, constant.SingleChatType)

		} else {
			conversation.UserID = recvID
			conversation.ConversationID = utils.GetConversationIDBySessionType(recvID, constant.SingleChatType)
			_, err := c.db.GetConversation(conversation.ConversationID)
			if err != nil {
				faceUrl, name, err := c.cache.GetUserNameAndFaceURL(recvID, operationID)
				if err != nil {
					common.CheckAnyErrCallback(callback, 208, err, operationID)
				}
				conversation.FaceURL = faceUrl
				conversation.ShowName = name
			}
		}

		localMessage := model_struct.LocalChatLog{}
		s.SendID = sendID
		s.RecvID = recvID
		s.ClientMsgID = utils.GetMsgID(s.SendID)
		s.SendTime = utils.GetCurrentTimestampByMill()
		s.SessionType = constant.SingleChatType
		s.Status = constant.MsgStatusSendSuccess
		msgStructToLocalChatLog(&localMessage, &s)
		conversation.LatestMsg = utils.StructToJsonString(s)
		conversation.ConversationType = constant.SingleChatType
		conversation.LatestMsgSendTime = s.SendTime
		_ = c.insertMessageToLocalStorage(callback, &localMessage, operationID)
		callback.OnSuccess(utils.StructToJsonString(&s))
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: conversation}, c.GetCh())
	}()
}

func (c *Conversation) InsertGroupMessageToLocalStorage(callback open_im_sdk_callback.Base, message, groupID, sendID, operationID string) {
	go func() {
		log.NewInfo(operationID, "InsertSingleMessageToLocalStorage args: ", message, groupID, sendID)
		if groupID == "" || sendID == "" {
			common.CheckAnyErrCallback(callback, 208, errors.New("groupID or sendID is null"), operationID)
		}
		var conversation model_struct.LocalConversation
		var err error
		_, conversation.ConversationType, err = c.getConversationTypeByGroupID(groupID)
		common.CheckAnyErrCallback(callback, 202, err, operationID)
		conversation.ConversationID = utils.GetConversationIDBySessionType(groupID, int(conversation.ConversationType))
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		if sendID != c.loginUserID {
			faceUrl, name, err, isFromSvr := c.friend.GetUserNameAndFaceUrlByUid(sendID, operationID)
			if err != nil {
				log.Error(operationID, "getUserNameAndFaceUrlByUid err", err.Error(), sendID)
			}
			s.SenderFaceURL = faceUrl
			s.SenderNickname = name
			if isFromSvr {
				c.cache.Update(sendID, faceUrl, name)
			}
		}
		localMessage := model_struct.LocalChatLog{}
		s.SendID = sendID
		s.RecvID = groupID
		s.GroupID = groupID
		s.ClientMsgID = utils.GetMsgID(s.SendID)
		s.SendTime = utils.GetCurrentTimestampByMill()
		s.SessionType = conversation.ConversationType
		s.Status = constant.MsgStatusSendSuccess
		msgStructToLocalChatLog(&localMessage, &s)
		conversation.LatestMsg = utils.StructToJsonString(s)
		conversation.LatestMsgSendTime = s.SendTime
		conversation.FaceURL = s.SenderFaceURL
		conversation.ShowName = s.SenderNickname
		_ = c.insertMessageToLocalStorage(callback, &localMessage, operationID)
		callback.OnSuccess(utils.StructToJsonString(&s))
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: conversation}, c.GetCh())
	}()

}

//modifyLocalMessages(callback open_im_sdk_callback.Base, message, groupID, sendID, operationID string)

func (c *Conversation) SetConversationStatus(callback open_im_sdk_callback.Base, operationID string, userID string, status int) {
	if callback == nil {
		log.Error(operationID, "callback is nil")
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userID, status)
		//var unmarshalParams sdk.SetConversationStatusParams
		//common.JsonUnmarshalAndArgsValidate(userIDRemark, &unmarshalParams, callback, operationID)
		//f.setConversationStatus(unmarshalParams, callback, operationID)
		//callback.OnSuccess(utils.StructToJsonString(sdk.SetFriendRemarkCallback))
		//log.NewInfo(operationID, fName, " callback: ", utils.StructToJsonString(sdk.SetFriendRemarkCallback))
	}()
}

func (c *Conversation) SearchLocalMessages(callback open_im_sdk_callback.Base, searchParam, operationID string) {
	if callback == nil {
		return
	}
	go func() {
		s := time.Now()
		log.NewInfo(operationID, "SearchLocalMessages args: ", searchParam)
		var unmarshalParams sdk_params_callback.SearchLocalMessagesParams
		common.JsonUnmarshalCallback(searchParam, &unmarshalParams, callback, operationID)
		unmarshalParams.KeywordList = utils.TrimStringList(unmarshalParams.KeywordList)
		result := c.searchLocalMessages(callback, unmarshalParams, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, "cost time", time.Since(s))
		log.NewInfo(operationID, "SearchLocalMessages callback: ", result.TotalCount, len(result.SearchResultItems))
	}()
}
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

const TimeOffset = 5

func (c *Conversation) initBasicInfo(message *sdk_struct.MsgStruct, msgFrom, contentType int32, operationID string) {
	message.CreateTime = utils.GetCurrentTimestampByMill()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = constant.MsgStatusSending
	message.SendID = c.loginUserID
	userInfo, err := c.db.GetLoginUser(c.loginUserID)
	if err != nil {
		log.Error(operationID, "GetLoginUser ", err.Error(), c.loginUserID)
	} else {
		message.SenderFaceURL = userInfo.FaceURL
		message.SenderNickname = userInfo.Nickname
	}
	ClientMsgID := utils.GetMsgID(message.SendID)
	message.ClientMsgID = ClientMsgID
	message.MsgFrom = msgFrom
	message.ContentType = contentType
	message.SenderPlatformID = c.platformID
	message.IsExternalExtensions = c.IsExternalExtensions
}

func (c *Conversation) DeleteConversationFromLocalAndSvr(callback open_im_sdk_callback.Base, conversationID string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", conversationID)
		c.deleteConversationAndMsgFromSvr(callback, conversationID, operationID)
		c.deleteConversation(callback, conversationID, operationID)
		callback.OnSuccess(sdk_params_callback.DeleteConversationCallback)
		log.NewInfo(operationID, fName, "callback: ", sdk_params_callback.DeleteConversationCallback)
	}()

}

func (c *Conversation) DeleteMessageFromLocalAndSvr(callback open_im_sdk_callback.Base, message string, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", message)
		s := sdk_struct.MsgStruct{}
		common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
		c.deleteMessageFromSvr(callback, &s, operationID)
		c.deleteMessageFromLocalStorage(callback, &s, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, fName, "callback: ", "")
	}()
}

func (c *Conversation) DeleteAllMsgFromLocalAndSvr(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName)
		//c.deleteAllMsgFromSvr(callback, operationID)
		c.clearMessageFromSvr(callback, operationID)
		c.deleteAllMsgFromLocal(callback, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, fName, "callback: ", "")
	}()
}

func (c *Conversation) DeleteAllMsgFromLocal(callback open_im_sdk_callback.Base, operationID string) {
	if callback == nil {
		return
	}
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName)
		c.deleteAllMsgFromLocal(callback, operationID)
		callback.OnSuccess("")
		log.NewInfo(operationID, fName, "callback: ", "12")
	}()
}
func (c *Conversation) getConversationTypeByGroupID(groupID string) (conversationID string, conversationType int32, err error) {
	g, err := c.full.GetGroupInfoByGroupID(groupID)
	if err != nil {
		return "", 0, utils.Wrap(err, "get group info error")
	}
	switch g.GroupType {
	case constant.NormalGroup:
		return utils.GetConversationIDBySessionType(groupID, constant.GroupChatType), constant.GroupChatType, nil
	case constant.SuperGroup, constant.WorkingGroup:
		return utils.GetConversationIDBySessionType(groupID, constant.SuperGroupChatType), constant.SuperGroupChatType, nil
	default:
		return "", 0, utils.Wrap(errors.New("err groupType"), "group type err")
	}
}

func (c *Conversation) SetMessageReactionExtensions(callback open_im_sdk_callback.Base, message, reactionExtensionList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.SetMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionList, &unmarshalParams, callback, operationID)
	result := c.setMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
	//c.modifyGroupMessageReaction(callback, counter, reactionType, operationType, groupID, msgID, operationID)
}

func (c *Conversation) AddMessageReactionExtensions(callback open_im_sdk_callback.Base, message, reactionExtensionList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.AddMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionList, &unmarshalParams, callback, operationID)
	result := c.addMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

func (c *Conversation) DeleteMessageReactionExtensions(callback open_im_sdk_callback.Base, message, reactionExtensionKeyList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.DeleteMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionKeyList, &unmarshalParams, callback, operationID)
	result := c.deleteMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

func (c *Conversation) GetMessageListReactionExtensions(callback open_im_sdk_callback.Base, messageList, operationID string) {
	var list []*sdk_struct.MsgStruct
	common.JsonUnmarshalAndArgsValidate(messageList, &list, callback, operationID)
	result := c.getMessageListReactionExtensions(callback, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

/**
**Get some reaction extensions in reactionExtensionKeyList of message list
 */
func (c *Conversation) GetMessageListSomeReactionExtensions(callback open_im_sdk_callback.Base, messageList, reactionExtensionKeyList, operationID string) {
	var messagelist []*sdk_struct.MsgStruct
	common.JsonUnmarshalAndArgsValidate(messageList, &messagelist, callback, operationID)
	var list []string
	common.JsonUnmarshalAndArgsValidate(reactionExtensionKeyList, &list, callback, operationID)
	result := c.getMessageListSomeReactionExtensions(callback, messagelist, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) SetTypeKeyInfo(callback open_im_sdk_callback.Base, message, typeKey, ex string, isCanRepeat bool, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	result := c.setTypeKeyInfo(callback, &s, typeKey, ex, isCanRepeat, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) GetTypeKeyListInfo(callback open_im_sdk_callback.Base, message, typeKeyList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var list []string
	common.JsonUnmarshalAndArgsValidate(typeKeyList, &list, callback, operationID)
	result := c.getTypeKeyListInfo(callback, &s, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) GetAllTypeKeyInfo(callback open_im_sdk_callback.Base, message, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	result := c.getAllTypeKeyInfo(callback, &s, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
