package ws_local_server

import (
	"encoding/json"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"os"
	"runtime/pprof"
)

//
//import (
//	"encoding/json"
//	"open_im_sdk/open_im_sdk"
//)
//
func (wsRouter *WsFuncRouter) CreateTextMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	msg := userWorker.Conversation().CreateTextMessage(input, operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

type SendCallback struct {
	BaseSuccessFailed
	clientMsgID string
	//uid         string
}

func (s *SendCallback) OnProgress(progress int) {
	mReply := make(map[string]interface{})
	mReply["progress"] = progress
	jsonStr, _ := json.Marshal(mReply)

	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(jsonStr), "0"}, s.uid)
}

func (wsRouter *WsFuncRouter) SendMessage(input string, operationID string) {
	f, _ := os.OpenFile("operationID", os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
	userWorker.Conversation().SendMessage(&sc, m["message"].(string), m["recvID"].(string), m["groupID"].(string), m["offlinePushInfo"].(string), operationID)
	pprof.Lookup("heap").WriteTo(f, 0)
}

type AddAdvancedMsgListenerCallback struct {
	uid string
}

func (a *AddAdvancedMsgListenerCallback) OnRecvNewMessage(message string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", message, "0"}, a.uid)
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

func (wsRouter *WsFuncRouter) SetAdvancedMsgListener() {
	var msgCallback AddAdvancedMsgListenerCallback
	msgCallback.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetAdvancedMsgListener(&msgCallback)

}

type ConversationCallback struct {
	uid string
}

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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationIDList", "opt") {
		return
	}
	userWorker.Conversation().SetConversationRecvMessageOpt(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationIDList"].(string), int(m["opt"].(float64)), operationID)
}

func (wsRouter *WsFuncRouter) GetConversationRecvMessageOpt(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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

func (wsRouter *WsFuncRouter) SetConversationDraft(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "draftText") {
		return
	}
	userWorker.Conversation().SetConversationDraft(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationID"].(string), m["draftText"].(string), operationID)
}

func (wsRouter *WsFuncRouter) PinConversation(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "conversationID", "privateChat") {
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "atUserIDList") {
		return
	}
	msg := userWorker.Conversation().CreateTextAtMessage(m["text"].(string), m["atUserIDList"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateLocationMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
	wrapSdkLog(operationID, utils.GetSelfFuncName(), "CreateCustomMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "data", "extension", "description") {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "key not in, failed", input, m)
		return
	}
	wrapSdkLog(operationID, utils.GetSelfFuncName(), "GlobalSendMessage", input)
	msg := userWorker.Conversation().CreateCustomMessage(m["data"].(string), m["extension"].(string), m["description"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateQuoteMessage(input string, operationID string) {
	wrapSdkLog(operationID, utils.GetSelfFuncName(), "CreateQuoteMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "text", "message") {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "key not in, failed", input)
		return
	}
	wrapSdkLog(operationID, utils.GetSelfFuncName(), "GlobalSendMessage")
	msg := userWorker.Conversation().CreateQuoteMessage(m["text"].(string), m["message"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}
func (wsRouter *WsFuncRouter) CreateCardMessage(input string, operationID string) {
	wrapSdkLog(operationID, "CreateCardMessage", input)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "imageFullPath") {
		return
	}
	msg := userWorker.Conversation().CreateImageMessageFromFullPath(m["imageFullPath"].(string), operationID)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateSoundMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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

func (wsRouter *WsFuncRouter) RevokeMessage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().RevokeMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}

func (wsRouter *WsFuncRouter) TypingStatusUpdate(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "userID", "msgIDList") {
		return
	}
	userWorker.Conversation().MarkC2CMessageAsRead(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["userID"].(string), m["msgIDList"].(string), operationID)
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if userWorker.AdvancedFunction() == nil {
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusNotSupportFunction, "unsupported function", "", operationID})
	}
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "groupID", "msgIDList") {
		return
	}
	userWorker.AdvancedFunction().MarkGroupMessageAsRead(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["groupID"].(string), m["msgIDList"].(string), operationID)
}

func (wsRouter *WsFuncRouter) DeleteMessageFromLocalStorage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteMessageFromLocalStorage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}

func (wsRouter *WsFuncRouter) DeleteMessage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, message, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().DeleteMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, message, operationID)
}

func (wsRouter *WsFuncRouter) InsertSingleMessageToLocalStorage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), m, "message", "groupID", "sendID") {
		return
	}
	userWorker.Conversation().InsertGroupMessageToLocalStorage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, m["message"].(string), m["groupID"].(string), m["sendID"].(string), operationID)
}

//func (wsRouter *WsFuncRouter) FindMessages(messageIDList string, operationID string) {
//	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
//	userWorker.Conversation().FindMessages(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, messageIDList)
//}

func (wsRouter *WsFuncRouter) SearchLocalMessages(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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

func (wsRouter *WsFuncRouter) SendMessageNotOss(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog(operationID, utils.GetSelfFuncName(), "unmarshal failed", input, err.Error())
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

func (wsRouter *WsFuncRouter) ClearGroupHistoryMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkResourceLoadingAndKeysIn(userWorker, input, operationID, runFuncName(), nil) {
		return
	}
	userWorker.Conversation().ClearGroupHistoryMessage(&BaseSuccessFailed{runFuncName(), operationID, wsRouter.uId}, input, operationID)
}

//func (wsRouter *WsFuncRouter) SetSdkLog(input string, operationID string) {
//	m := make(map[string]interface{})
//	if err := json.Unmarshal([]byte(input), &m); err != nil {
//		wrapSdkLog("unmarshal failed")
//		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
//		return
//	}
//	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "flag") {
//		return
//	}
//	userWorker := init.GetUserWorker(wsRouter.uId)
//	userWorker.SetSdkLog(m["flag"].(int32))
//}
