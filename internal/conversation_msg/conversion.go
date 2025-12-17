// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conversation_msg

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"

	pbConversation "github.com/openimsdk/protocol/conversation"
)

func ServerConversationToLocal(conversation *pbConversation.Conversation) *model_struct.LocalConversation {
	return &model_struct.LocalConversation{
		ConversationID:   conversation.ConversationID,
		ConversationType: conversation.ConversationType,
		UserID:           conversation.UserID,
		GroupID:          conversation.GroupID,
		RecvMsgOpt:       conversation.RecvMsgOpt,
		GroupAtType:      conversation.GroupAtType,
		IsPinned:         conversation.IsPinned,
		BurnDuration:     conversation.BurnDuration,
		IsPrivateChat:    conversation.IsPrivateChat,
		AttachedInfo:     conversation.AttachedInfo,
		Ex:               conversation.Ex,
		MsgDestructTime:  conversation.MsgDestructTime,
		IsMsgDestruct:    conversation.IsMsgDestruct,
	}
}

func LocalConversationToServer(conversation *model_struct.LocalConversation) *pbConversation.Conversation {
	return &pbConversation.Conversation{
		ConversationID:   conversation.ConversationID,
		ConversationType: conversation.ConversationType,
		UserID:           conversation.UserID,
		GroupID:          conversation.GroupID,
		RecvMsgOpt:       conversation.RecvMsgOpt,
		GroupAtType:      conversation.GroupAtType,
		IsPinned:         conversation.IsPinned,
		BurnDuration:     conversation.BurnDuration,
		IsPrivateChat:    conversation.IsPrivateChat,
		AttachedInfo:     conversation.AttachedInfo,
		MsgDestructTime:  conversation.MsgDestructTime,
		Ex:               conversation.Ex,
		IsMsgDestruct:    conversation.IsMsgDestruct,
	}
}

func MsgDataToLocalChatLog(serverMessage *sdkws.MsgData) *model_struct.LocalChatLog {
	localMessage := &model_struct.LocalChatLog{
		ClientMsgID:      serverMessage.ClientMsgID,
		ServerMsgID:      serverMessage.ServerMsgID,
		SendID:           serverMessage.SendID,
		RecvID:           serverMessage.RecvID,
		SenderPlatformID: serverMessage.SenderPlatformID,
		SenderNickname:   serverMessage.SenderNickname,
		SenderFaceURL:    serverMessage.SenderFaceURL,
		SessionType:      serverMessage.SessionType,
		MsgFrom:          serverMessage.MsgFrom,
		ContentType:      serverMessage.ContentType,
		Content:          string(serverMessage.Content),
		IsRead:           serverMessage.IsRead,
		Seq:              serverMessage.Seq,
		SendTime:         serverMessage.SendTime,
		CreateTime:       serverMessage.CreateTime,
		AttachedInfo:     serverMessage.AttachedInfo,
		Ex:               serverMessage.Ex,
	}

	if serverMessage.Status >= constant.MsgStatusHasDeleted {
		localMessage.Status = serverMessage.Status
	} else {
		localMessage.Status = constant.MsgStatusSendSuccess
	}

	if serverMessage.SessionType == constant.WriteGroupChatType || serverMessage.SessionType == constant.ReadGroupChatType {
		localMessage.RecvID = serverMessage.GroupID
	}
	return localMessage
}

func LocalChatLogToMsgStruct(local *model_struct.LocalChatLog) *sdk_struct.MsgStruct {
	if local == nil {
		return nil
	}
	msg := &sdk_struct.MsgStruct{
		ClientMsgID:      local.ClientMsgID,
		ServerMsgID:      local.ServerMsgID,
		CreateTime:       local.CreateTime,
		SendTime:         local.SendTime,
		SessionType:      local.SessionType,
		SendID:           local.SendID,
		RecvID:           local.RecvID,
		MsgFrom:          local.MsgFrom,
		ContentType:      local.ContentType,
		SenderPlatformID: local.SenderPlatformID,
		SenderNickname:   local.SenderNickname,
		SenderFaceURL:    local.SenderFaceURL,
		Content:          local.Content,
		Seq:              local.Seq,
		IsRead:           local.IsRead,
		Status:           local.Status,
		Ex:               local.Ex,
		LocalEx:          local.LocalEx,
		AttachedInfo:     local.AttachedInfo,
	}

	if err := PopulateMsgStructByContentType(msg); err != nil {
		log.ZWarn(context.Background(), "Parsing data error", err, "messageContent", msg.Content)
	}

	switch local.SessionType {
	case constant.WriteGroupChatType, constant.ReadGroupChatType:
		msg.GroupID = local.RecvID
	}
	return msg
}

func PopulateMsgStructByContentType(msg *sdk_struct.MsgStruct) (err error) {
	switch msg.ContentType {
	case constant.Text:
		elem := sdk_struct.TextElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.TextElem = &elem
	case constant.Picture:
		elem := sdk_struct.PictureElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.PictureElem = &elem
	case constant.Sound:
		elem := sdk_struct.SoundElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.SoundElem = &elem
	case constant.Video:
		elem := sdk_struct.VideoElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.VideoElem = &elem
	case constant.File:
		elem := sdk_struct.FileElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.FileElem = &elem
	case constant.AdvancedText:
		elem := sdk_struct.AdvancedTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.AdvancedTextElem = &elem
	case constant.AtText:
		elem := sdk_struct.AtTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.AtTextElem = &elem
	case constant.Location:
		elem := sdk_struct.LocationElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.LocationElem = &elem
	case constant.Custom, constant.CustomMsgNotTriggerConversation, constant.CustomMsgOnlineOnly:
		elem := sdk_struct.CustomElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.CustomElem = &elem
	case constant.Typing:
		elem := sdk_struct.TypingElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.TypingElem = &elem
	case constant.Quote:
		elem := sdk_struct.QuoteElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.QuoteElem = &elem
	case constant.Merger:
		elem := sdk_struct.MergeElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.MergeElem = &elem
	case constant.Face:
		elem := sdk_struct.FaceElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.FaceElem = &elem
	case constant.Card:
		elem := sdk_struct.CardElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.CardElem = &elem
	case constant.MarkdownText:
		elem := sdk_struct.MarkdownTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.MarkdownTextElem = &elem
	default:
		elem := sdk_struct.NotificationElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.NotificationElem = &elem
	}
	var attachedInfo sdk_struct.AttachedInfoElem
	if msg.AttachedInfo != "" {
		if err := utils.JsonStringToStruct(msg.AttachedInfo, &attachedInfo); err != nil {
			log.ZWarn(context.Background(), "JsonStringToStruct error", err, "localMessage.AttachedInfo", msg.AttachedInfo)
		}
	}
	msg.AttachedInfoElem = &attachedInfo
	msg.Content = ""
	return errs.Wrap(err)
}

func msgHandleByContentType(msg *sdk_struct.MsgStruct) (err error) {
	switch msg.ContentType {
	case constant.Text:
		t := sdk_struct.TextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.TextElem = &t
	case constant.Picture:
		t := sdk_struct.PictureElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.PictureElem = &t
	case constant.Sound:
		t := sdk_struct.SoundElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.SoundElem = &t
	case constant.Video:
		t := sdk_struct.VideoElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.VideoElem = &t
	case constant.File:
		t := sdk_struct.FileElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.FileElem = &t
	case constant.AdvancedText:
		t := sdk_struct.AdvancedTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.AdvancedTextElem = &t
	case constant.AtText:
		t := sdk_struct.AtTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.AtTextElem = &t
	case constant.Location:
		t := sdk_struct.LocationElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.LocationElem = &t
	case constant.Custom:
		fallthrough
	case constant.CustomMsgNotTriggerConversation:
		fallthrough
	case constant.CustomMsgOnlineOnly:
		t := sdk_struct.CustomElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.CustomElem = &t
	case constant.Typing:
		t := sdk_struct.TypingElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.TypingElem = &t
	case constant.Quote:
		t := sdk_struct.QuoteElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.QuoteElem = &t
	case constant.Merger:
		t := sdk_struct.MergeElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.MergeElem = &t
	case constant.Face:
		t := sdk_struct.FaceElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.FaceElem = &t
	case constant.Card:
		t := sdk_struct.CardElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.CardElem = &t
	case constant.MarkdownText:
		t := sdk_struct.MarkdownTextElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.MarkdownTextElem = &t
	default:
		t := sdk_struct.NotificationElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.NotificationElem = &t
	}
	msg.Content = ""

	return errs.Wrap(err)
}

func MsgStructToLocalChatLog(message *sdk_struct.MsgStruct) *model_struct.LocalChatLog {
	localMessage := &model_struct.LocalChatLog{
		ClientMsgID:      message.ClientMsgID,
		ServerMsgID:      message.ServerMsgID,
		SendID:           message.SendID,
		RecvID:           message.RecvID,
		SenderPlatformID: message.SenderPlatformID,
		SenderNickname:   message.SenderNickname,
		SenderFaceURL:    message.SenderFaceURL,
		SessionType:      message.SessionType,
		MsgFrom:          message.MsgFrom,
		ContentType:      message.ContentType,
		IsRead:           message.IsRead,
		Status:           message.Status,
		Seq:              message.Seq,
		SendTime:         message.SendTime,
		CreateTime:       message.CreateTime,
		AttachedInfo:     message.AttachedInfo,
		Ex:               message.Ex,
		LocalEx:          message.LocalEx,
	}
	switch message.ContentType {
	case constant.Text:
		localMessage.Content = utils.StructToJsonString(message.TextElem)
	case constant.Picture:
		localMessage.Content = utils.StructToJsonString(message.PictureElem)
	case constant.Sound:
		localMessage.Content = utils.StructToJsonString(message.SoundElem)
	case constant.Video:
		localMessage.Content = utils.StructToJsonString(message.VideoElem)
	case constant.File:
		localMessage.Content = utils.StructToJsonString(message.FileElem)
	case constant.AtText:
		localMessage.Content = utils.StructToJsonString(message.AtTextElem)
	case constant.Merger:
		localMessage.Content = utils.StructToJsonString(message.MergeElem)
	case constant.Card:
		localMessage.Content = utils.StructToJsonString(message.CardElem)
	case constant.Location:
		localMessage.Content = utils.StructToJsonString(message.LocationElem)
	case constant.Custom:
		localMessage.Content = utils.StructToJsonString(message.CustomElem)
	case constant.Quote:
		localMessage.Content = utils.StructToJsonString(message.QuoteElem)
	case constant.Face:
		localMessage.Content = utils.StructToJsonString(message.FaceElem)
	case constant.AdvancedText:
		localMessage.Content = utils.StructToJsonString(message.AdvancedTextElem)
	case constant.MarkdownText:
		localMessage.Content = utils.StructToJsonString(message.MarkdownTextElem)
	default:
		localMessage.Content = utils.StructToJsonString(message.NotificationElem)
	}
	if message.SessionType == constant.WriteGroupChatType || message.SessionType == constant.ReadGroupChatType {
		localMessage.RecvID = message.GroupID
	}
	localMessage.AttachedInfo = utils.StructToJsonString(message.AttachedInfoElem)
	return localMessage
}
