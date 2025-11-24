package converter

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

func ServerConversationToLocal(info *pbConversation.Conversation) *model_struct.LocalConversation {
	if info == nil {
		return nil
	}
	return &model_struct.LocalConversation{
		ConversationID:   info.ConversationID,
		ConversationType: info.ConversationType,
		UserID:           info.UserID,
		GroupID:          info.GroupID,
		RecvMsgOpt:       info.RecvMsgOpt,
		GroupAtType:      info.GroupAtType,
		IsPinned:         info.IsPinned,
		BurnDuration:     info.BurnDuration,
		IsPrivateChat:    info.IsPrivateChat,
		AttachedInfo:     info.AttachedInfo,
		Ex:               info.Ex,
		MsgDestructTime:  info.MsgDestructTime,
		IsMsgDestruct:    info.IsMsgDestruct,
	}
}

func LocalConversationToServer(info *model_struct.LocalConversation) *pbConversation.Conversation {
	if info == nil {
		return nil
	}
	return &pbConversation.Conversation{
		ConversationID:   info.ConversationID,
		ConversationType: info.ConversationType,
		UserID:           info.UserID,
		GroupID:          info.GroupID,
		RecvMsgOpt:       info.RecvMsgOpt,
		GroupAtType:      info.GroupAtType,
		IsPinned:         info.IsPinned,
		BurnDuration:     info.BurnDuration,
		IsPrivateChat:    info.IsPrivateChat,
		AttachedInfo:     info.AttachedInfo,
		MsgDestructTime:  info.MsgDestructTime,
		Ex:               info.Ex,
		IsMsgDestruct:    info.IsMsgDestruct,
	}
}

func MsgDataToLocalChatLog(info *sdkws.MsgData) *model_struct.LocalChatLog {
	if info == nil {
		return nil
	}
	local := &model_struct.LocalChatLog{
		ClientMsgID:      info.ClientMsgID,
		ServerMsgID:      info.ServerMsgID,
		SendID:           info.SendID,
		RecvID:           info.RecvID,
		SenderPlatformID: info.SenderPlatformID,
		SenderNickname:   info.SenderNickname,
		SenderFaceURL:    info.SenderFaceURL,
		SessionType:      info.SessionType,
		MsgFrom:          info.MsgFrom,
		ContentType:      info.ContentType,
		Content:          string(info.Content),
		IsRead:           info.IsRead,
		Seq:              info.Seq,
		SendTime:         info.SendTime,
		CreateTime:       info.CreateTime,
		AttachedInfo:     info.AttachedInfo,
		Ex:               info.Ex,
	}

	if info.Status >= constant.MsgStatusHasDeleted {
		local.Status = info.Status
	} else {
		local.Status = constant.MsgStatusSendSuccess
	}

	if info.SessionType == constant.WriteGroupChatType || info.SessionType == constant.ReadGroupChatType {
		local.RecvID = info.GroupID
	}
	return local
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
	}

	var attachedInfo sdk_struct.AttachedInfoElem
	if err := utils.JsonStringToStruct(local.AttachedInfo, &attachedInfo); err != nil {
		log.ZWarn(context.Background(), "JsonStringToStruct error", err, "localMessage.AttachedInfo", local.AttachedInfo)
	}
	msg.AttachedInfoElem = &attachedInfo

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
	default:
		elem := sdk_struct.NotificationElem{}
		err = utils.JsonStringToStruct(msg.Content, &elem)
		msg.NotificationElem = &elem
	}
	msg.Content = ""
	return errs.Wrap(err)
}

func MsgStructToLocalChatLog(message *sdk_struct.MsgStruct) *model_struct.LocalChatLog {
	if message == nil {
		return nil
	}
	local := &model_struct.LocalChatLog{
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
		local.Content = utils.StructToJsonString(message.TextElem)
	case constant.Picture:
		local.Content = utils.StructToJsonString(message.PictureElem)
	case constant.Sound:
		local.Content = utils.StructToJsonString(message.SoundElem)
	case constant.Video:
		local.Content = utils.StructToJsonString(message.VideoElem)
	case constant.File:
		local.Content = utils.StructToJsonString(message.FileElem)
	case constant.AtText:
		local.Content = utils.StructToJsonString(message.AtTextElem)
	case constant.Merger:
		local.Content = utils.StructToJsonString(message.MergeElem)
	case constant.Card:
		local.Content = utils.StructToJsonString(message.CardElem)
	case constant.Location:
		local.Content = utils.StructToJsonString(message.LocationElem)
	case constant.Custom, constant.CustomMsgNotTriggerConversation, constant.CustomMsgOnlineOnly:
		local.Content = utils.StructToJsonString(message.CustomElem)
	case constant.Quote:
		local.Content = utils.StructToJsonString(message.QuoteElem)
	case constant.Face:
		local.Content = utils.StructToJsonString(message.FaceElem)
	case constant.AdvancedText:
		local.Content = utils.StructToJsonString(message.AdvancedTextElem)
	default:
		local.Content = utils.StructToJsonString(message.NotificationElem)
	}

	if message.SessionType == constant.WriteGroupChatType || message.SessionType == constant.ReadGroupChatType {
		local.RecvID = message.GroupID
	}

	local.AttachedInfo = utils.StructToJsonString(message.AttachedInfoElem)
	return local
}

func MsgStructToMsgData(message *sdk_struct.MsgStruct, options map[string]bool) *sdkws.MsgData {
	if message == nil {
		return nil
	}
	data := &sdkws.MsgData{
		SendID:           message.SendID,
		RecvID:           message.RecvID,
		GroupID:          message.GroupID,
		ClientMsgID:      message.ClientMsgID,
		ServerMsgID:      message.ServerMsgID,
		SenderPlatformID: message.SenderPlatformID,
		SenderNickname:   message.SenderNickname,
		SenderFaceURL:    message.SenderFaceURL,
		SessionType:      message.SessionType,
		MsgFrom:          message.MsgFrom,
		ContentType:      message.ContentType,
		Content:          []byte(message.Content),
		Seq:              message.Seq,
		SendTime:         message.SendTime,
		CreateTime:       message.CreateTime,
		Status:           message.Status,
		IsRead:           message.IsRead,
		OfflinePushInfo:  message.OfflinePush,
		AttachedInfo:     message.AttachedInfo,
		Ex:               message.Ex,
	}
	if atElem := message.AtTextElem; atElem != nil && len(atElem.AtUserList) > 0 {
		data.AtUserIDList = append([]string(nil), atElem.AtUserList...)
	}
	if message.AttachedInfoElem != nil {
		data.AttachedInfo = utils.StructToJsonString(message.AttachedInfoElem)
	}
	if len(options) > 0 {
		data.Options = make(map[string]bool, len(options))
		for k, v := range options {
			data.Options[k] = v
		}
	}
	return data
}

func MsgDataToMsgStruct(serverMessage *sdkws.MsgData) *sdk_struct.MsgStruct {
	if serverMessage == nil {
		return nil
	}
	return &sdk_struct.MsgStruct{
		ClientMsgID:      serverMessage.ClientMsgID,
		ServerMsgID:      serverMessage.ServerMsgID,
		CreateTime:       serverMessage.CreateTime,
		SendTime:         serverMessage.SendTime,
		SessionType:      serverMessage.SessionType,
		SendID:           serverMessage.SendID,
		RecvID:           serverMessage.RecvID,
		GroupID:          serverMessage.GroupID,
		MsgFrom:          serverMessage.MsgFrom,
		ContentType:      serverMessage.ContentType,
		SenderPlatformID: serverMessage.SenderPlatformID,
		SenderNickname:   serverMessage.SenderNickname,
		SenderFaceURL:    serverMessage.SenderFaceURL,
		Content:          string(serverMessage.Content),
		Seq:              serverMessage.Seq,
		IsRead:           serverMessage.IsRead,
		Status:           serverMessage.Status,
		OfflinePush:      serverMessage.OfflinePushInfo,
		AttachedInfo:     serverMessage.AttachedInfo,
		Ex:               serverMessage.Ex,
	}
}

func MsgHandleByContentType(msg *sdk_struct.MsgStruct) (err error) {
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
	default:
		t := sdk_struct.NotificationElem{}
		err = utils.JsonStringToStruct(msg.Content, &t)
		msg.NotificationElem = &t
	}
	msg.Content = ""

	return errs.Wrap(err)
}
