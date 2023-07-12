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

package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
)

// import (
//
//	"encoding/json"
//	"open_im_sdk/open_im_sdk"
//
// )
func (wsRouter *WsFuncRouter) CreateTextMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	log.Info(operationID, "CreateTextMessage start ", input)
	msg := userWorker.Conversation().CreateTextMessage(input, operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateAdvancedTextMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "messageEntityList") {
		return
	}
	msg := userWorker.Conversation().CreateAdvancedTextMessage(m["text"].(string), m["messageEntityList"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

type SendCallback struct {
	BaseSuccessFailed
	clientMsgID string
}

func (s *SendCallback) OnProgress(progress int) {
	mReply := make(map[string]interface{})
	mReply["progress"] = progress
	mReply["clientMsgID"] = s.clientMsgID
	jsonStr, _ := json.Marshal(mReply)

	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(jsonStr), ""}, s.uid)
}

func (wsRouter *WsFuncRouter) SendMessage(input string, operationID string) {

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.operationID = operationID
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "message", "recvID", "groupID", "offlinePushInfo") {
		return
	}
	s := sdk_struct.MsgStruct{}
	common.JsonUnmarshalAndArgsValidate(m["message"].(string), &s, &sc, operationID)
	sc.clientMsgID = s.ClientMsgID
	userWorker.Conversation().SendMessage(&sc, m["message"].(string), m["recvID"].(string), m["groupID"].(string), m["offlinePushInfo"].(string), operationID)

}

type AddAdvancedMsgListenerCallback struct {
	uid string
}

type BatchMsgListenerCallback struct {
	uid string
}

func (b *BatchMsgListenerCallback) OnRecvNewMessages(messageList string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", messageList, "0"}, b.uid)
}

func (a *AddAdvancedMsgListenerCallback) OnRecvNewMessage(message string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", message, "0"}, a.uid)
}
func (a *BatchMsgListenerCallback) OnRecvOfflineNewMessages(messageList string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", messageList, "0"}, a.uid)
}

func (a *AddAdvancedMsgListenerCallback) OnRecvC2CReadReceipt(msgReceiptList string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msgReceiptList, "0"}, a.uid)
}
func (a *AddAdvancedMsgListenerCallback) OnRecvGroupReadReceipt(groupMsgReceiptList string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", groupMsgReceiptList, "0"}, a.uid)
}
func (a *AddAdvancedMsgListenerCallback) OnRecvMessageRevoked(msgId string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msgId, "0"}, a.uid)
}
func (a *AddAdvancedMsgListenerCallback) OnNewRecvMessageRevoked(messageRevoked string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", messageRevoked, "0"}, a.uid)
}
func (a *AddAdvancedMsgListenerCallback) OnRecvMessageExtensionsChanged(msgID string, reactionExtensionList string) {
	m := make(map[string]interface{})
	m["msgID"] = msgID
	m["reactionExtensionList"] = reactionExtensionList
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", utils.StructToJsonString(m), "0"}, a.uid)

}

func (a *AddAdvancedMsgListenerCallback) OnRecvMessageExtensionsDeleted(msgID string, reactionExtensionKeyList string) {
	m := make(map[string]interface{})
	m["msgID"] = msgID
	m["reactionExtensionKeyList"] = reactionExtensionKeyList
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", utils.StructToJsonString(m), "0"}, a.uid)
}
func (a *AddAdvancedMsgListenerCallback) OnRecvMessageExtensionsAdded(msgID string, reactionExtensionList string) {
	m := make(map[string]interface{})
	m["msgID"] = msgID
	m["reactionExtensionKeyList"] = reactionExtensionList
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", utils.StructToJsonString(m), "0"}, a.uid)
}
func (wsRouter *WsFuncRouter) SetAdvancedMsgListener() {
	var msgCallback AddAdvancedMsgListenerCallback
	msgCallback.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetAdvancedMsgListener(&msgCallback)
}

func (wsRouter *WsFuncRouter) SetBatchMsgListener() {
	var callback BatchMsgListenerCallback
	callback.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetBatchMsgListener(&callback)
}

type ConversationCallback struct {
	uid string
}

//funcation (c *ConversationCallback) OnSyncServerProgress(progress int) {
//	var ed EventData
//	ed.Event = cleanUpfuncName(runFuncName())
//	ed.ErrCode = 0
//	ed.Data = utils.IntToString(progress)
//	SendOneUserMessage(ed, c.uid)
//}

func (c *ConversationCallback) OnSyncServerStart() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, c.uid)
}
func (c *ConversationCallback) OnSyncServerFinish() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, c.uid)
}
func (c *ConversationCallback) OnSyncServerFailed() {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	SendOneUserMessage(ed, c.uid)
}
func (c *ConversationCallback) OnNewConversation(conversationList string) {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	ed.Data = conversationList
	SendOneUserMessage(ed, c.uid)
}

func (c *ConversationCallback) OnConversationChanged(conversationList string) {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	ed.Data = conversationList
	SendOneUserMessage(ed, c.uid)
}
func (c *ConversationCallback) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	var ed EventData
	ed.Event = cleanUpfuncName(runFuncName())
	ed.ErrCode = 0
	ed.Data = int32ToString(totalUnreadCount)
	SendOneUserMessage(ed, c.uid)
}

func (wsRouter *WsFuncRouter) SetConversationListener() {
	var ccb ConversationCallback
	ccb.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetConversationListener(&ccb)
}

func (wsRouter *WsFuncRouter) GetAllConversationList(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetAllConversationList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}
func (wsRouter *WsFuncRouter) GetConversationListSplit(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "offset", "count") {
		return
	}
	userWorker.Conversation().GetConversationListSplit(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, int(m["offset"].(float64)), int(m["count"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) SetOneConversationRecvMessageOpt(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationIDList", "opt") {
		return
	}
	userWorker.Conversation().SetOneConversationRecvMessageOpt(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), int(m["opt"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) SetConversationRecvMessageOpt(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationIDList", "opt") {
		return
	}
	userWorker.Conversation().SetConversationRecvMessageOpt(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationIDList"].(string), int(m["opt"].(float64)), operationID)
}
func (wsRouter *WsFuncRouter) SetGlobalRecvMessageOpt(input, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "opt") {
		return
	}
	userWorker.Conversation().SetGlobalRecvMessageOpt(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, int(m["opt"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) GetConversationRecvMessageOpt(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationIDList") {
		return
	}
	userWorker.Conversation().GetConversationRecvMessageOpt(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetOneConversation(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "sourceID", "sessionType") {
		return
	}
	userWorker.Conversation().GetOneConversation(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, int32(m["sessionType"].(float64)), m["sourceID"].(string), operationID)
}

func (wsRouter *WsFuncRouter) GetMultipleConversation(conversationIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, conversationIDList, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetMultipleConversation(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, conversationIDList, operationID)
}

func (wsRouter *WsFuncRouter) DeleteConversation(conversationID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, conversationID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteConversation(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, conversationID, operationID)
}
func (wsRouter *WsFuncRouter) DeleteAllConversationFromLocal(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteAllConversationFromLocal(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) SetConversationDraft(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "draftText") {
		return
	}
	userWorker.Conversation().SetConversationDraft(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), m["draftText"].(string), operationID)
}
func (wsRouter *WsFuncRouter) ResetConversationGroupAtType(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Conversation().ResetConversationGroupAtType(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) PinConversation(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "isPinned") {
		return
	}
	userWorker.Conversation().PinConversation(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), m["isPinned"].(bool), operationID)
}

func (wsRouter *WsFuncRouter) SetOneConversationPrivateChat(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "isPrivate") {
		return
	}
	userWorker.Conversation().SetOneConversationPrivateChat(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), m["isPrivate"].(bool), operationID)
}

func (wsRouter *WsFuncRouter) GetTotalUnreadMsgCount(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetTotalUnreadMsgCount(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) CreateTextAtMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "atUserIDList", "atUsersInfo", "message") {
		return
	}
	msg := userWorker.Conversation().CreateTextAtMessage(m["text"].(string), m["atUserIDList"].(string), m["atUsersInfo"].(string), m["message"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateLocationMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "description", "longitude", "latitude") {
		return
	}
	msg := userWorker.Conversation().CreateLocationMessage(m["description"].(string), m["longitude"].(float64), m["latitude"].(float64), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateCustomMessage(input string, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "CreateCustomMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "data", "extension", "description") {
		log.Info(operationID, utils.GetSelfFuncName(), "key not in, failed", input, m)
		return
	}
	log.Info(operationID, utils.GetSelfFuncName(), "GlobalSendMessage", input)
	msg := userWorker.Conversation().CreateCustomMessage(m["data"].(string), m["extension"].(string), m["description"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateQuoteMessage(input string, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "CreateQuoteMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "message") {
		log.Info(operationID, utils.GetSelfFuncName(), "key not in, failed", input)
		return
	}
	log.Info(operationID, utils.GetSelfFuncName(), "GlobalSendMessage")
	msg := userWorker.Conversation().CreateQuoteMessage(m["text"].(string), m["message"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateAdvancedQuoteMessage(input string, operationID string) {
	log.Info(operationID, utils.GetSelfFuncName(), "CreateQuoteMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "message", "messageEntityList") {
		log.Info(operationID, utils.GetSelfFuncName(), "key not in, failed", input)
		return
	}
	log.Info(operationID, utils.GetSelfFuncName(), "GlobalSendMessage")
	msg := userWorker.Conversation().CreateAdvancedQuoteMessage(m["text"].(string), m["message"].(string), m["messageEntityList"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateCardMessage(input string, operationID string) {
	log.Info(operationID, "CreateCardMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	msg := userWorker.Conversation().CreateCardMessage(input, operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateVideoMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "videoFullPath", "videoType", "duration", "snapshotFullPath") {
		return
	}
	msg := userWorker.Conversation().CreateVideoMessageFromFullPath(m["videoFullPath"].(string), m["videoType"].(string), int64(m["duration"].(float64)), m["snapshotFullPath"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateImageMessageFromFullPath(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	msg := userWorker.Conversation().CreateImageMessageFromFullPath(input, operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateSoundMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "soundPath", "duration") {
		return
	}
	msg := userWorker.Conversation().CreateSoundMessageFromFullPath(m["soundPath"].(string), int64(m["duration"].(float64)), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateMergerMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "messageList", "title", "summaryList") {
		return
	}
	msg := userWorker.Conversation().CreateMergerMessage(m["messageList"].(string), m["title"].(string), m["summaryList"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateFaceMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "index", "data") {
		return
	}
	msg := userWorker.Conversation().CreateFaceMessage(int(m["index"].(float64)), m["data"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateForwardMessage(m string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, m, operationID, runFuncName(), nil) {
		return
	}
	msg := userWorker.Conversation().CreateForwardMessage(m, operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) FindMessageList(findMessageOptions string, operationID string) {
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, findMessageOptions, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().FindMessageList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, findMessageOptions, operationID)
}
func (wsRouter *WsFuncRouter) GetHistoryMessageList(getMessageOptions string, operationID string) {
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, getMessageOptions, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetHistoryMessageList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, getMessageOptions, operationID)
}
func (wsRouter *WsFuncRouter) GetAdvancedHistoryMessageList(getMessageOptions string, operationID string) {
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, getMessageOptions, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetAdvancedHistoryMessageList(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, getMessageOptions, operationID)
}
func (wsRouter *WsFuncRouter) GetHistoryMessageListReverse(getMessageOptions string, operationID string) {
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, getMessageOptions, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().GetHistoryMessageListReverse(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, getMessageOptions, operationID)
}

// deprecated
func (wsRouter *WsFuncRouter) RevokeMessage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().RevokeMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}
func (wsRouter *WsFuncRouter) NewRevokeMessage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().NewRevokeMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}
func (wsRouter *WsFuncRouter) TypingStatusUpdate(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "recvID", "msgTip") {
		return
	}
	userWorker.Conversation().TypingStatusUpdate(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["recvID"].(string), m["msgTip"].(string), operationID)
}

func (wsRouter *WsFuncRouter) MarkC2CMessageAsRead(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "userID", "msgIDList") {
		return
	}
	userWorker.Conversation().MarkC2CMessageAsRead(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["userID"].(string), m["msgIDList"].(string), operationID)
}
func (wsRouter *WsFuncRouter) MarkMessageAsReadByConID(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "msgIDList") {
		return
	}
	userWorker.Conversation().MarkMessageAsReadByConID(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), m["msgIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) MarkGroupMessageHasRead(groupID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, groupID, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().MarkGroupMessageHasRead(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, groupID, operationID)
}
func (wsRouter *WsFuncRouter) MarkGroupMessageAsRead(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "msgIDList") {
		return
	}
	userWorker.Conversation().MarkGroupMessageAsRead(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["msgIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) DeleteMessageFromLocalStorage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteMessageFromLocalStorage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}

func (wsRouter *WsFuncRouter) DeleteMessageFromLocalAndSvr(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteMessageFromLocalAndSvr(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}

func (wsRouter *WsFuncRouter) DeleteAllMsgFromLocalAndSvr(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Conversation().DeleteAllMsgFromLocalAndSvr(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) DeleteAllMsgFromLocal(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Conversation().DeleteAllMsgFromLocal(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, operationID)
}

func (wsRouter *WsFuncRouter) DeleteConversationFromLocalAndSvr(conversationID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.Conversation().DeleteConversationFromLocalAndSvr(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, conversationID, operationID)
}

func (wsRouter *WsFuncRouter) InsertSingleMessageToLocalStorage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "message", "recvID", "sendID") {
		return
	}
	userWorker.Conversation().InsertSingleMessageToLocalStorage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["message"].(string), m["recvID"].(string), m["sendID"].(string), operationID)
}
func (wsRouter *WsFuncRouter) InsertGroupMessageToLocalStorage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "message", "groupID", "sendID") {
		return
	}
	userWorker.Conversation().InsertGroupMessageToLocalStorage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["message"].(string), m["groupID"].(string), m["sendID"].(string), operationID)
}

//funcation (wsRouter *WsFuncRouter) FindMessages(messageIDList string, operationID string) {
//	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
//	userWorker.Conversation().FindMessages(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, messageIDList)
//}

func (wsRouter *WsFuncRouter) SearchLocalMessages(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	//if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "searchParam") {
	//	return
	//}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().SearchLocalMessages(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId},
		input, operationID)
}

func (wsRouter *WsFuncRouter) CreateImageMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "sourcePicture", "bigPicture", "snapshotPicture") {
		return
	}
	msg := userWorker.Conversation().CreateImageMessageByURL(m["sourcePicture"].(string), m["bigPicture"].(string), m["snapshotPicture"].(string), operationID)

	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})

}

func (wsRouter *WsFuncRouter) CreateSoundMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "soundBaseInfo") {
		return
	}
	msg := userWorker.Conversation().CreateSoundMessageByURL(m["soundBaseInfo"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateVideoMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "videoBaseInfo") {
		return
	}
	msg := userWorker.Conversation().CreateVideoMessageByURL(m["videoBaseInfo"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateFileMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "fileBaseInfo") {
		return
	}

	msg := userWorker.Conversation().CreateFileMessageByURL(m["fileBaseInfo"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})

}

func (wsRouter *WsFuncRouter) CreateFileMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "fileFullPath", "fileName") {
		return
	}

	msg := userWorker.Conversation().CreateFileMessageFromFullPath(m["fileFullPath"].(string), m["fileName"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})

}

func (wsRouter *WsFuncRouter) SendMessageNotOss(input string, operationID string) {

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		log.Info(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	sc.operationID = operationID
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "message", "recvID", "groupID", "offlinePushInfo") {
		return
	}
	userWorker.Conversation().SendMessageNotOss(&sc, m["message"].(string), m["recvID"].(string), m["groupID"].(string), m["offlinePushInfo"].(string), operationID)

}

func (wsRouter *WsFuncRouter) ClearC2CHistoryMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().ClearC2CHistoryMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}
func (wsRouter *WsFuncRouter) ClearC2CHistoryMessageFromLocalAndSvr(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().ClearC2CHistoryMessageFromLocalAndSvr(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

func (wsRouter *WsFuncRouter) ClearGroupHistoryMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().ClearGroupHistoryMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}
func (wsRouter *WsFuncRouter) ClearGroupHistoryMessageFromLocalAndSvr(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().ClearGroupHistoryMessageFromLocalAndSvr(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

//	funcation (wsRouter *WsFuncRouter) SetSdkLog(input string, operationID string) {
//		m := make(map[string]interface{})
//		if err := json.Unmarshal([]byte(input), &m); err != nil {
//			log.Info("unmarshal failed")
//			wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
//			return
//		}
//		if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "flag") {
//			return
//		}
//		userWorker := init.GetUserWorker(wsRouter.uId)
//		userWorker.SetSdkLog(m["flag"].(int32))
//	}
func (wsRouter *WsFuncRouter) GetAtAllTag(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	msg := constant.AtAllString
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
