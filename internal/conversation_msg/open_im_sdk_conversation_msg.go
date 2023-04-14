package conversation_msg

import (
	"bytes"
	"context"
	"errors"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	pbUser "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"image"
	"open_im_sdk/internal/util"
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
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	imgtype "github.com/shamsher31/goimgtype"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func (c *Conversation) GetAllConversationList(ctx context.Context) ([]*model_struct.LocalConversation, error) {
	return c.db.GetAllConversationListDB(ctx)
}
func (c *Conversation) GetConversationListSplit(ctx context.Context, offset, count int) ([]*model_struct.LocalConversation, error) {
	return c.db.GetConversationListSplitDB(ctx, offset, count)
}
func (c *Conversation) SetConversationRecvMessageOpt(ctx context.Context, conversationIDList []string, opt int) error {
	var conversations []*pbConversation.Conversation
	for _, conversationID := range conversationIDList {
		localConversation, err := c.db.GetConversation(ctx, conversationID)
		if err != nil {
			//log.NewError(operationID, utils.GetSelfFuncName(), "GetConversation failed", err.Error())
			continue
		}
		if localConversation.ConversationType == constant.SuperGroupChatType && opt == constant.NotReceiveMessage {
			return errors.New("super group not support this opt")
		}
		conversations = append(conversations, &pbConversation.Conversation{
			OwnerUserID:      c.loginUserID,
			ConversationID:   conversationID,
			ConversationType: localConversation.ConversationType,
			UserID:           localConversation.UserID,
			GroupID:          localConversation.GroupID,
			RecvMsgOpt:       int32(opt),
			IsPinned:         localConversation.IsPinned,
			IsPrivateChat:    localConversation.IsPrivateChat,
			AttachedInfo:     localConversation.AttachedInfo,
			Ex:               localConversation.Ex,
		})
	}
	req := &pbConversation.BatchSetConversationsReq{
		Conversations: conversations,
		OwnerUserID:   c.loginUserID,
	}
	_, err := util.CallApi[pbConversation.BatchSetConversationsResp](ctx, constant.BatchSetConversationRouter, req)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil
}
func (c *Conversation) SetGlobalRecvMessageOpt(ctx context.Context, opt int) error {
	if err := util.ApiPost(ctx, constant.SetGlobalRecvMessageOptRouter, &pbUser.SetGlobalRecvMessageOptReq{UserID: c.loginUserID, GlobalRecvMsgOpt: int32(opt)}, nil); err != nil {
		return err
	}
	c.user.SyncLoginUserInfo(ctx)
	return nil

}
func (c *Conversation) HideConversation(ctx context.Context, conversationID string) error {
	return c.db.UpdateColumnsConversation(ctx, conversationID, map[string]interface{}{"latest_msg_send_time": 0})
}

// deprecated
func (c *Conversation) GetConversationRecvMessageOpt(ctx context.Context, conversationIDList []string) (resp []*server_api_params.GetConversationRecvMessageOptResp, err error) {
	conversations, err := c.db.GetMultipleConversationDB(ctx, conversationIDList)
	if err != nil {
		return nil, err
	}
	for _, conversation := range conversations {
		resp = append(resp, &server_api_params.GetConversationRecvMessageOptResp{
			ConversationID: conversation.ConversationID,
			Result:         &conversation.RecvMsgOpt,
		})
	}
	return resp, nil

}
func (c *Conversation) GetOneConversation(ctx context.Context, sessionType int32, sourceID string) (*model_struct.LocalConversation, error) {
	conversationID := utils.GetConversationIDBySessionType(sourceID, int(sessionType))
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err == nil {
		return lc, nil
	} else {
		var newConversation model_struct.LocalConversation
		newConversation.ConversationID = conversationID
		newConversation.ConversationType = sessionType
		switch sessionType {
		case constant.SingleChatType:
			newConversation.UserID = sourceID
			faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, sourceID)
			if err != nil {
				return nil, err
			}
			newConversation.ShowName = name
			newConversation.FaceURL = faceUrl
		case constant.GroupChatType, constant.SuperGroupChatType:
			newConversation.GroupID = sourceID
			g, err := c.full.GetGroupInfoFromLocal2Svr(sourceID, sessionType)
			if err != nil {
				return nil, err
			}
			newConversation.ShowName = g.GroupName
			newConversation.FaceURL = g.FaceURL
		}
		lc, errTemp := c.db.GetConversation(ctx, conversationID)
		if errTemp == nil {
			return lc, nil
		}
		err := c.db.InsertConversation(ctx, &newConversation)
		if err != nil {
			return nil, err
		}
		return &newConversation, nil
	}

}
func (c *Conversation) GetMultipleConversation(ctx context.Context, conversationIDList []string) ([]*model_struct.LocalConversation, error) {
	conversations, err := c.db.GetMultipleConversationDB(ctx, conversationIDList)
	if err != nil {
		return nil, err
	}
	return conversations, nil

}
func (c *Conversation) DeleteConversation(ctx context.Context, conversationID string) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	var sourceID string
	switch lc.ConversationType {
	case constant.SingleChatType, constant.NotificationChatType:
		sourceID = lc.UserID
	case constant.GroupChatType, constant.SuperGroupChatType:
		sourceID = lc.GroupID
	}
	if lc.ConversationType == constant.SuperGroupChatType {
		err = c.db.SuperGroupDeleteAllMessage(ctx, lc.GroupID)
		if err != nil {
			return err
		}
	} else {
		//Mark messages related to this conversation for deletion
		err = c.db.UpdateMessageStatusBySourceIDController(ctx, sourceID, constant.MsgStatusHasDeleted, lc.ConversationType)
		if err != nil {
			return err
		}
	}
	//Reset the conversation information, empty conversation
	err = c.db.ResetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	c.doUpdateConversation(common.Cmd2Value{Value: common.UpdateConNode{"", constant.TotalUnreadMessageChanged, ""}})
	return nil

}
func (c *Conversation) DeleteAllConversationFromLocal(ctx context.Context) error {
	err := c.db.ResetAllConversation(ctx)
	if err != nil {
		return err
	}
	return nil

}
func (c *Conversation) SetConversationDraft(ctx context.Context, conversationID, draftText string) error {
	if draftText != "" {
		err := c.db.SetConversationDraftDB(ctx, conversationID, draftText)
		if err != nil {
			return err
		}
	} else {
		err := c.db.RemoveConversationDraft(ctx, conversationID, draftText)
		if err != nil {
			return err
		}
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
	return nil

}
func (c *Conversation) ResetConversationGroupAtType(ctx context.Context, conversationID string) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	if lc.GroupAtType == constant.AtNormal || lc.ConversationType != constant.GroupChatType {
		return errors.New("conversation don't need to reset")
	}
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	apiReq.Conversation.GroupAtType = constant.AtNormal
	apiReq.FieldType = constant.FieldGroupAtType
	err = c.setConversation(ctx, apiReq, lc)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil

}
func (c *Conversation) PinConversation(ctx context.Context, conversationID string, isPinned bool) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	apiReq.Conversation.IsPinned = isPinned
	apiReq.FieldType = constant.FieldIsPinned
	err = c.setConversation(ctx, apiReq, lc)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil
}

func (c *Conversation) SetOneConversationPrivateChat(ctx context.Context, conversationID string, isPrivate bool) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	apiReq.Conversation.IsPrivateChat = isPrivate
	apiReq.FieldType = constant.FieldIsPrivateChat
	err = c.setConversation(ctx, apiReq, lc)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil
}

func (c *Conversation) SetOneConversationBurnDuration(ctx context.Context, conversationID string, burnDuration int32) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	apiReq.Conversation.BurnDuration = burnDuration
	apiReq.FieldType = constant.FieldBurnDuration
	err = c.setConversation(ctx, apiReq, lc)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil
}

func (c *Conversation) SetOneConversationRecvMessageOpt(ctx context.Context, conversationID string, opt int) error {
	lc, err := c.db.GetConversation(ctx, conversationID)
	if err != nil {
		return err
	}
	apiReq := &pbConversation.ModifyConversationFieldReq{}
	apiReq.Conversation.BurnDuration = int32(opt)
	apiReq.FieldType = constant.FieldRecvMsgOpt
	err = c.setConversation(ctx, apiReq, lc)
	if err != nil {
		return err
	}
	c.SyncConversations(ctx, 0)
	return nil
}

func (c *Conversation) GetTotalUnreadMsgCount(ctx context.Context) (totalUnreadCount int32, err error) {
	return c.db.GetTotalUnreadMsgCountDB(ctx)
}

func (c *Conversation) SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	if c.ConversationListener != nil {
		log.Error("internal", "just only set on listener")
		return
	}
	c.ConversationListener = listener
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
				faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, recvID, operationID)
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
				soundURL, uuid, err := c.UploadSoundByBuffer(buffer1, s.SoundElem.DataSize, s.SoundElem.SoundType, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.SoundElem.SourceURL = soundURL
				s.SoundElem.UUID = uuid
				s.Content = utils.StructToJsonString(s.SoundElem)

			case constant.Video:

				snapshotURL, snapshotUUID, videoURL, videoUUID, err := c.UploadVideoByBuffer(buffer1, buffer2, s.VideoElem.VideoSize,
					s.VideoElem.SnapshotSize, s.VideoElem.VideoType, s.VideoElem.SnapshotType, callback.OnProgress)
				c.checkErrAndUpdateMessage(callback, 301, err, &s, lc, operationID)
				s.VideoElem.VideoURL = videoURL
				s.VideoElem.SnapshotUUID = snapshotUUID
				s.VideoElem.SnapshotURL = snapshotURL
				s.VideoElem.VideoUUID = videoUUID
				s.Content = utils.StructToJsonString(s.VideoElem)
			case constant.File:
				fileURL, fileUUID, err := c.UploadFileByBuffer(buffer1, s.FileElem.FileSize, s.FileElem.FileType, callback.OnProgress)
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

func (c *Conversation) InternalSendMessage(ctx context.Context, s *sdk_struct.MsgStruct, recvID, groupID string, p *server_api_params.OfflinePushInfo, onlineUserOnly bool, options map[string]bool) (*server_api_params.UserSendMsgResp, error) {
	if recvID == "" && groupID == "" {
		return nil, errors.New("recvID && groupID not both null")
	}
	if recvID == "" {
		g, err := c.full.GetGroupInfoByGroupID(groupID)
		if err != nil {
			return nil, err
		}
		switch g.GroupType {
		case constant.NormalGroup:
			s.SessionType = constant.GroupChatType
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(ctx, groupID)
			if err != nil {
				return nil, err
			}
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				return nil, errors.New("you not exist in this group")
			}

		case constant.SuperGroup:
			s.SessionType = constant.SuperGroupChatType
		case constant.WorkingGroup:
			s.SessionType = constant.SuperGroupChatType
			groupMemberUidList, err := c.db.GetGroupMemberUIDListByGroupID(ctx, groupID)
			if err != nil {
				return nil, err
			}
			if !utils.IsContain(s.SendID, groupMemberUidList) {
				return nil, errors.New("you not exist in this group")
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
	if err != nil {
		return nil, err
	}
	//switch e := err.(type) {
	//case *constant.ErrInfo:
	//	common.CheckAnyErrCallback(callback, e.ErrCode, e, operationID)
	//default:
	//	common.CheckAnyErrCallback(callback, 301, err, operationID)
	//}
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

func (c *Conversation) FindMessageList(ctx context.Context, req []*sdk_params_callback.ConversationArgs) (*sdk_params_callback.FindMessageListCallback, error) {
	var r sdk_params_callback.FindMessageListCallback
	{
	}
	type tempConversationAndMessageList struct {
		conversation *model_struct.LocalConversation
		msgIDList    []string
	}
	var s []*tempConversationAndMessageList
	for _, conversationsArgs := range req {
		localConversation, err := c.db.GetConversation(ctx, conversationsArgs.ConversationID)
		if err == nil {
			t := new(tempConversationAndMessageList)
			t.conversation = localConversation
			t.msgIDList = conversationsArgs.ClientMsgIDList
			s = append(s, t)
		} else {
			//log.Error(operationID, "GetConversation err:", err.Error(), conversationsArgs.ConversationID)
		}
	}
	for _, v := range s {
		messages, err := c.db.GetMultipleMessageController(ctx, v.msgIDList, v.conversation.GroupID, v.conversation.ConversationType)
		if err == nil {
			var tempMessageList []*sdk_struct.MsgStruct
			for _, message := range messages {
				temp := sdk_struct.MsgStruct{}
				temp.ClientMsgID = message.ClientMsgID
				temp.ServerMsgID = message.ServerMsgID
				temp.CreateTime = message.CreateTime
				temp.SendTime = message.SendTime
				temp.SessionType = message.SessionType
				temp.SendID = message.SendID
				temp.RecvID = message.RecvID
				temp.MsgFrom = message.MsgFrom
				temp.ContentType = message.ContentType
				temp.SenderPlatformID = message.SenderPlatformID
				temp.SenderNickname = message.SenderNickname
				temp.SenderFaceURL = message.SenderFaceURL
				temp.Content = message.Content
				temp.Seq = message.Seq
				temp.IsRead = message.IsRead
				temp.Status = message.Status
				temp.AttachedInfo = message.AttachedInfo
				temp.Ex = message.Ex
				err := c.msgHandleByContentType(&temp)
				if err != nil {
					//log.Error(operationID, "Parsing data error:", err.Error(), temp)
					continue
				}
				switch message.SessionType {
				case constant.GroupChatType:
					fallthrough
				case constant.SuperGroupChatType:
					temp.GroupID = temp.RecvID
					temp.RecvID = c.loginUserID
				}
				tempMessageList = append(tempMessageList, &temp)
			}
			findResultItem := sdk_params_callback.SearchByConversationResult{}
			findResultItem.ConversationID = v.conversation.ConversationID
			findResultItem.FaceURL = v.conversation.FaceURL
			findResultItem.ShowName = v.conversation.ShowName
			findResultItem.ConversationType = v.conversation.ConversationType
			findResultItem.MessageList = tempMessageList
			findResultItem.MessageCount = len(findResultItem.MessageList)
			r.FindResultItems = append(r.FindResultItems, &findResultItem)
			r.TotalCount += findResultItem.MessageCount
		} else {
			//log.Error(operationID, "GetMultipleMessageController err:", err.Error(), v)
		}
	}
	return &r, nil

}
func (c *Conversation) GetHistoryMessageList(ctx context.Context, req sdk_params_callback.GetHistoryMessageListParams) ([]*sdk_struct.MsgStruct, error) {
	return c.getHistoryMessageList(ctx, req, false)
}
func (c *Conversation) GetAdvancedHistoryMessageList(ctx context.Context, req sdk_params_callback.GetAdvancedHistoryMessageListParams) (*sdk_params_callback.GetAdvancedHistoryMessageListCallback, error) {
	result, err := c.getAdvancedHistoryMessageList(ctx, req, false)
	if err != nil {
		return nil, err
	}
	if len(result.MessageList) == 0 {
		s := make([]*sdk_struct.MsgStruct, 0)
		result.MessageList = s
	}
	return result, nil
}
func (c *Conversation) GetAdvancedHistoryMessageListReverse(ctx context.Context, req sdk_params_callback.GetAdvancedHistoryMessageListParams) (*sdk_params_callback.GetAdvancedHistoryMessageListCallback, error) {
	result, err := c.getAdvancedHistoryMessageList(ctx, req, true)
	if err != nil {
		return nil, err
	}
	if len(result.MessageList) == 0 {
		s := make([]*sdk_struct.MsgStruct, 0)
		result.MessageList = s
	}
	return result, nil
}

func (c *Conversation) GetHistoryMessageListReverse(ctx context.Context, req sdk_params_callback.GetHistoryMessageListParams) ([]*sdk_struct.MsgStruct, error) {
	return c.getHistoryMessageList(ctx, req, true)
}
func (c *Conversation) RevokeMessage(ctx context.Context, req *sdk_struct.MsgStruct) error {
	return c.revokeOneMessage(ctx, req)
}
func (c *Conversation) NewRevokeMessage(ctx context.Context, req *sdk_struct.MsgStruct) error {
	return c.newRevokeOneMessage(ctx, req)

}
func (c *Conversation) TypingStatusUpdate(ctx context.Context, recvID, msgTip string) error {
	return c.typingStatusUpdate(ctx, recvID, msgTip)
}

func (c *Conversation) MarkC2CMessageAsRead(ctx context.Context, userID string, msgIDList []string) error {
	if len(msgIDList) == 0 {
		conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
		_ = c.setOneConversationUnread(ctx, conversationID, 0)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		return nil
	}
	return c.markC2CMessageAsRead(ctx, msgIDList, userID)
}

func (c *Conversation) MarkMessageAsReadByConID(ctx context.Context, conversationID string, msgIDList []string) error {
	if len(msgIDList) == 0 {
		_ = c.setOneConversationUnread(ctx, conversationID, 0)
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
		_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())
		return nil
	}
	return nil

}

// deprecated
func (c *Conversation) MarkGroupMessageHasRead(ctx context.Context, groupID string) {
	conversationID := utils.GetConversationIDBySessionType(groupID, constant.GroupChatType)
	_ = c.setOneConversationUnread(ctx, conversationID, 0)
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.UnreadCountSetZero}, c.GetCh())
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversationID, Action: constant.ConChange, Args: []string{conversationID}}, c.GetCh())

}
func (c *Conversation) MarkGroupMessageAsRead(ctx context.Context, groupID string, msgIDList []string) error {
	return c.markGroupMessageAsRead(ctx, msgIDList, groupID)

}
func (c *Conversation) DeleteMessageFromLocalStorage(ctx context.Context, message *sdk_struct.MsgStruct) error {
	return c.deleteMessageFromLocalStorage(ctx, message)

}

func (c *Conversation) ClearC2CHistoryMessage(ctx context.Context, userID string, operationID string) error {
	return c.clearC2CHistoryMessage(ctx, userID)
}
func (c *Conversation) ClearGroupHistoryMessage(ctx context.Context, groupID string) error {
	return c.clearGroupHistoryMessage(ctx, groupID)

}
func (c *Conversation) ClearC2CHistoryMessageFromLocalAndSvr(ctx context.Context, userID string) error {

	conversationID := utils.GetConversationIDBySessionType(userID, constant.SingleChatType)
	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearC2CHistoryMessage(ctx, userID)

}

// fixme
func (c *Conversation) ClearGroupHistoryMessageFromLocalAndSvr(ctx context.Context, groupID string) error {
	conversationID, _, err := c.getConversationTypeByGroupID(groupID)
	if err != nil {
		return err
	}
	err = c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.clearGroupHistoryMessage(ctx, groupID)
}

func (c *Conversation) InsertSingleMessageToLocalStorage(ctx context.Context, s *sdk_struct.MsgStruct, recvID, sendID string) (*sdk_struct.MsgStruct, error) {
	if recvID == "" || sendID == "" {
		return nil, errors.New("recvID or sendID is null")
	}
	var conversation model_struct.LocalConversation
	if sendID != c.loginUserID {
		faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, sendID)
		if err != nil {
			//log.Error(operationID, "GetUserNameAndFaceURL err", err.Error(), sendID)
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
		_, err := c.db.GetConversation(ctx, conversation.ConversationID)
		if err != nil {
			faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, recvID)
			if err != nil {
				return nil, err
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
	msgStructToLocalChatLog(&localMessage, s)
	conversation.LatestMsg = utils.StructToJsonString(s)
	conversation.ConversationType = constant.SingleChatType
	conversation.LatestMsgSendTime = s.SendTime
	err := c.insertMessageToLocalStorage(ctx, &localMessage)
	if err != nil {
		return nil, err
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: conversation}, c.GetCh())
	return s, nil

}

func (c *Conversation) InsertGroupMessageToLocalStorage(ctx context.Context, s *sdk_struct.MsgStruct, groupID, sendID string) (*sdk_struct.MsgStruct, error) {
	if groupID == "" || sendID == "" {
		return nil, errors.New("groupID or sendID is null")
	}
	var conversation model_struct.LocalConversation
	var err error
	_, conversation.ConversationType, err = c.getConversationTypeByGroupID(groupID)
	if err != nil {
		return nil, err
	}

	conversation.ConversationID = utils.GetConversationIDBySessionType(groupID, int(conversation.ConversationType))
	if sendID != c.loginUserID {
		faceUrl, name, err := c.cache.GetUserNameAndFaceURL(ctx, sendID)
		if err != nil {
			log.Error("", "getUserNameAndFaceUrlByUid err", err.Error(), sendID)
		}
		s.SenderFaceURL = faceUrl
		s.SenderNickname = name
	}
	localMessage := model_struct.LocalChatLog{}
	s.SendID = sendID
	s.RecvID = groupID
	s.GroupID = groupID
	s.ClientMsgID = utils.GetMsgID(s.SendID)
	s.SendTime = utils.GetCurrentTimestampByMill()
	s.SessionType = conversation.ConversationType
	s.Status = constant.MsgStatusSendSuccess
	msgStructToLocalChatLog(&localMessage, s)
	conversation.LatestMsg = utils.StructToJsonString(s)
	conversation.LatestMsgSendTime = s.SendTime
	conversation.FaceURL = s.SenderFaceURL
	conversation.ShowName = s.SenderNickname
	err = c.insertMessageToLocalStorage(ctx, &localMessage)
	if err != nil {
		return nil, err
	}
	_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{ConID: conversation.ConversationID, Action: constant.AddConOrUpLatMsg, Args: conversation}, c.GetCh())
	return s, nil

}

func (c *Conversation) SearchLocalMessages(ctx context.Context, searchParam *sdk_params_callback.SearchLocalMessagesParams) (*sdk_params_callback.SearchLocalMessagesCallback, error) {

	searchParam.KeywordList = utils.TrimStringList(searchParam.KeywordList)
	return c.searchLocalMessages(ctx, searchParam)

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

func (c *Conversation) initBasicInfo(ctx context.Context, message *sdk_struct.MsgStruct, msgFrom, contentType int32) error {
	message.CreateTime = utils.GetCurrentTimestampByMill()
	message.SendTime = message.CreateTime
	message.IsRead = false
	message.Status = constant.MsgStatusSending
	message.SendID = c.loginUserID
	userInfo, err := c.db.GetLoginUser(ctx, c.loginUserID)
	if err != nil {
		return err
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
	return nil
}

func (c *Conversation) DeleteConversationFromLocalAndSvr(ctx context.Context, conversationID string) error {

	err := c.deleteConversationAndMsgFromSvr(ctx, conversationID)
	if err != nil {
		return err
	}
	return c.deleteConversation(ctx, conversationID)

}

func (c *Conversation) DeleteMessageFromLocalAndSvr(ctx context.Context, s *sdk_struct.MsgStruct) error {
	err := c.deleteMessageFromSvr(ctx, s)
	if err != nil {
		return err
	}
	return c.deleteMessageFromLocalStorage(ctx, s)
}

func (c *Conversation) DeleteAllMsgFromLocalAndSvr(ctx context.Context, operationID string) {
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

func (c *Conversation) DeleteAllMsgFromLocal(ctx context.Context, operationID string) {
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

func (c *Conversation) SetMessageReactionExtensions(ctx context.Context, message, reactionExtensionList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.SetMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionList, &unmarshalParams, callback, operationID)
	result := c.setMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
	//c.modifyGroupMessageReaction(callback, counter, reactionType, operationType, groupID, msgID, operationID)
}

func (c *Conversation) AddMessageReactionExtensions(ctx context.Context, message, reactionExtensionList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.AddMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionList, &unmarshalParams, callback, operationID)
	result := c.addMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

func (c *Conversation) DeleteMessageReactionExtensions(ctx context.Context, message, reactionExtensionKeyList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var unmarshalParams sdk_params_callback.DeleteMessageReactionExtensionsParams
	common.JsonUnmarshalCallback(reactionExtensionKeyList, &unmarshalParams, callback, operationID)
	result := c.deleteMessageReactionExtensions(callback, &s, unmarshalParams, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

func (c *Conversation) GetMessageListReactionExtensions(ctx context.Context, messageList, operationID string) {
	var list []*sdk_struct.MsgStruct
	common.JsonUnmarshalAndArgsValidate(messageList, &list, callback, operationID)
	result := c.getMessageListReactionExtensions(callback, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}

/**
**Get some reaction extensions in reactionExtensionKeyList of message list
 */
func (c *Conversation) GetMessageListSomeReactionExtensions(ctx context.Context, messageList, reactionExtensionKeyList, operationID string) {
	var messagelist []*sdk_struct.MsgStruct
	common.JsonUnmarshalAndArgsValidate(messageList, &messagelist, callback, operationID)
	var list []string
	common.JsonUnmarshalAndArgsValidate(reactionExtensionKeyList, &list, callback, operationID)
	result := c.getMessageListSomeReactionExtensions(callback, messagelist, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) SetTypeKeyInfo(ctx context.Context, message, typeKey, ex string, isCanRepeat bool, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	result := c.setTypeKeyInfo(callback, &s, typeKey, ex, isCanRepeat, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) GetTypeKeyListInfo(ctx context.Context, message, typeKeyList, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	var list []string
	common.JsonUnmarshalAndArgsValidate(typeKeyList, &list, callback, operationID)
	result := c.getTypeKeyListInfo(callback, &s, list, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
func (c *Conversation) GetAllTypeKeyInfo(ctx context.Context, message, operationID string) {
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(message, &s, callback, operationID)
	result := c.getAllTypeKeyInfo(callback, &s, operationID)
	callback.OnSuccess(utils.StructToJsonString(result))
	log.NewInfo(operationID, utils.GetSelfFuncName(), "callback: ", utils.StructToJsonString(result))
}
