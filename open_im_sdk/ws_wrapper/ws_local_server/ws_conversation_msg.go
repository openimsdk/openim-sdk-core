package ws_local_server

import (
	"encoding/json"
	//"net/http"
	"open_im_sdk/open_im_sdk"
)

func (wsRouter *WsFuncRouter) CreateTextMessage(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	msg := userWorker.CreateTextMessage(input)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

type SendCallback struct {
	BaseSuccFailed
	clientMsgID string
	//uid         string
}

func (s *SendCallback) OnProgress(progress int) {
	mReply := make(map[string]interface{})
	mReply["progress"] = progress
	mReply["clientMsgID"] = s.clientMsgID
	jsonStr, _ := json.Marshal(mReply)

	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", string(jsonStr), "0"}, s.uid)
}

//func SendMessage(callback SendMsgCallBack, message, receiver, groupID string, onlineUserOnly bool) string {
func (wsRouter *WsFuncRouter) SendMessage(input string, operationID string) {

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	sc.operationID = operationID
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "message", "recvID", "groupID", "onlineUserOnly") {
		return
	}
	clientMsgID := userWorker.SendMessage(&sc, m["message"].(string), m["recvID"].(string), m["groupID"].(string), m["onlineUserOnly"].(bool))
	sc.clientMsgID = clientMsgID

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

func (a *AddAdvancedMsgListenerCallback) OnRecvMessageRevoked(msgId string) {
	SendOneUserMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msgId, "0"}, a.uid)
}

func (wsRouter *WsFuncRouter) AddAdvancedMsgListener() {
	var msgCallback AddAdvancedMsgListenerCallback
	msgCallback.uid = wsRouter.uId
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.AddAdvancedMsgListener(&msgCallback)

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
	userWorker.GetAllConversationList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}
func (wsRouter *WsFuncRouter) GetConversationListSplit(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "offset", "count") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetConversationListSplit(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, m["offset"].(int), m["count"].(int))
}
func (wsRouter *WsFuncRouter) SetConversationRecvMessageOpt(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "conversationIDList", "opt") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.SetConversationRecvMessageOpt(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationIDList"].(string), m["opt"].(int))
}
func (wsRouter *WsFuncRouter) GetConversationRecvMessageOpt(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "conversationIDList") {
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetConversationRecvMessageOpt(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, m["conversationIDList"].(string))
}

func (wsRouter *WsFuncRouter) GetOneConversation(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "sourceID", "sessionType") {
		return
	}

	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetOneConversation(m["sourceID"].(string), int(m["sessionType"].(float64)), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetMultipleConversation(conversationIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetMultipleConversation(conversationIDList, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) DeleteConversation(conversationID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.DeleteConversation(conversationID, &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) SetConversationDraft(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "conversationID", "draftText") {
		return
	}
	userWorker.SetConversationDraft(m["conversationID"].(string), m["draftText"].(string), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) PinConversation(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "conversationID", "isPinned") {
		return
	}
	userWorker.PinConversation(m["conversationID"].(string), m["isPinned"].(bool), &BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) GetTotalUnreadMsgCount(input string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetTotalUnreadMsgCount(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId})
}

func (wsRouter *WsFuncRouter) CreateTextAtMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "text", "atUserList") {
		return
	}
	msg := userWorker.CreateTextAtMessage(m["text"].(string), m["atUserList"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateLocationMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "description", "longitude", "latitude") {
		return
	}
	msg := userWorker.CreateLocationMessage(m["description"].(string), m["longitude"].(float64), m["latitude"].(float64))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

/*
func CreateCustomMessage(data, extension []byte, description string) string {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
	}
	msg := open_im_sdk.CreateLocationMessage(m["description"].(string), m["longitude"].(float64), m["latitude"].(float64))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg})
}
CreateCustomMessage(data, extension string, description string) string
*/

func (wsRouter *WsFuncRouter) CreateCustomMessage(input string, operationID string) {
	wrapSdkLog("CreateCustomMessage", input, operationID)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed", operationID)
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "data", "extension", "description") {
		wrapSdkLog("key not in, failed", operationID)
		return
	}
	wrapSdkLog("GlobalSendMessage", operationID)
	msg := userWorker.CreateCustomMessage(m["data"].(string), m["message"].(string), m["description"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateQuoteMessage(input string, operationID string) {
	wrapSdkLog("CreateQuoteMessage", input, operationID)
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed", operationID)
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "text", "message") {
		wrapSdkLog("key not in, failed", operationID)
		return
	}
	wrapSdkLog("GlobalSendMessage", operationID)
	msg := userWorker.CreateQuoteMessage(m["text"].(string), m["message"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateVideoMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "videoFullPath", "videoType", "duration", "snapshotFullPath") {
		return
	}
	msg := userWorker.CreateVideoMessageFromFullPath(m["videoFullPath"].(string), m["videoType"].(string), int64(m["duration"].(float64)), m["snapshotFullPath"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateImageMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "imageFullPath") {
		return
	}
	msg := userWorker.CreateImageMessageFromFullPath(m["imageFullPath"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateSoundMessageFromFullPath(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "soundPath", "duration") {
		return
	}
	msg := userWorker.CreateSoundMessageFromFullPath(m["soundPath"].(string), int64(m["duration"].(float64)))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateMergerMessage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "messageList", "title", "summaryList") {
		return
	}
	msg := userWorker.CreateMergerMessage(m["messageList"].(string), m["title"].(string), m["summaryList"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateForwardMessage(m string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	msg := userWorker.CreateForwardMessage(m)
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) GetHistoryMessageList(getMessageOptions string, operationID string) {
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.GetHistoryMessageList(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, getMessageOptions)
}

func (wsRouter *WsFuncRouter) RevokeMessage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.RevokeMessage(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, message)
}

func (wsRouter *WsFuncRouter) TypingStatusUpdate(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "receiver", "msgTip") {
		return
	}
	userWorker.TypingStatusUpdate(m["receiver"].(string), m["msgTip"].(string))
}

func (wsRouter *WsFuncRouter) MarkC2CMessageAsRead(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "receiver", "msgIDList") {
		return
	}
	userWorker.MarkC2CMessageAsRead(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, m["receiver"].(string), m["msgIDList"].(string))
}

func (wsRouter *WsFuncRouter) MarkSingleMessageHasRead(userID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.MarkSingleMessageHasRead(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, userID)
}

func (wsRouter *WsFuncRouter) MarkGroupMessageHasRead(groupID string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.MarkGroupMessageHasRead(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, groupID)
}

func (wsRouter *WsFuncRouter) DeleteMessageFromLocalStorage(message string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.DeleteMessageFromLocalStorage(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, message)
}

func (wsRouter *WsFuncRouter) InsertSingleMessageToLocalStorage(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "message", "userID", "sender") {
		return
	}
	userWorker.InsertSingleMessageToLocalStorage(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, m["message"].(string), m["userID"].(string), m["sender"].(string))
}

func (wsRouter *WsFuncRouter) FindMessages(messageIDList string, operationID string) {
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	userWorker.FindMessages(&BaseSuccFailed{runFuncName(), operationID, wsRouter.uId}, messageIDList)
}

func (wsRouter *WsFuncRouter) CreateImageMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "sourcePicture", "bigPicture", "snapshotPicture") {
		return
	}
	msg := userWorker.CreateImageMessageByURL(m["sourcePicture"].(string), m["bigPicture"].(string), m["snapshotPicture"].(string))

	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})

}

func (wsRouter *WsFuncRouter) CreateSoundMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "soundBaseInfo") {
		return
	}
	msg := userWorker.CreateSoundMessageByURL(m["soundBaseInfo"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateVideoMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "videoBaseInfo") {
		return
	}
	msg := userWorker.CreateVideoMessageByURL(m["videoBaseInfo"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})
}

func (wsRouter *WsFuncRouter) CreateFileMessageByURL(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "fileBaseInfo") {
		return
	}

	msg := userWorker.CreateFileMessageByURL(m["fileBaseInfo"].(string))
	wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), 0, "", msg, operationID})

}

func (wsRouter *WsFuncRouter) SendMessageNotOss(input string, operationID string) {
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(input), &m); err != nil {
		wrapSdkLog("unmarshal failed")
		wsRouter.GlobalSendMessage(EventData{cleanUpfuncName(runFuncName()), StatusBadParameter, "unmarshal failed", "", operationID})
		return
	}
	var sc SendCallback
	sc.uid = wsRouter.uId
	sc.funcName = runFuncName()
	sc.operationID = operationID
	userWorker := open_im_sdk.GetUserWorker(wsRouter.uId)
	if !wsRouter.checkKeysIn(input, operationID, runFuncName(), m, "message", "receiver", "groupID", "onlineUserOnly") {
		return
	}
	clientMsgID := userWorker.SendMessageNotOss(&sc, m["message"].(string), m["receiver"].(string), m["groupID"].(string), m["onlineUserOnly"].(bool))
	sc.clientMsgID = clientMsgID
}
