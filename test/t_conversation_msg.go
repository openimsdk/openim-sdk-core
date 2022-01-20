package test

import (
	"encoding/json"
	"fmt"
	"open_im_sdk/sdk_struct"
)

//func DotestSetConversationRecvMessageOpt() {
//	var callback BaseSuccFailedTest
//	callback.funcName = utils.GetSelfFuncName()
//	var idList []string
//	idList = append(idList, "18567155635")
//	jsontest, _ := json.Marshal(idList)
//	open_im_sdk.SetConversationRecvMessageOpt(&callback, string(jsontest), 2)
//	fmt.Println("SetConversationRecvMessageOpt", string(jsontest))
//}
//
//func DoTestGetMultipleConversation() {
//	var callback BaseSuccFailedTest
//	callback.funcName = utils.GetSelfFuncName()
//	var idList []string
//	fmt.Println("DoTestGetMultipleConversation come here")
//	idList = append(idList, "single_13977954313", "group_77215e1b13b75f3ab00cb6570e3d9618")
//	jsontest, _ := json.Marshal(idList)
//	open_im_sdk.GetMultipleConversation(string(jsontest), &callback)
//	fmt.Println("GetMultipleConversation", string(jsontest))
//}
//
//func DoTestGetConversationRecvMessageOpt() {
//	var callback BaseSuccFailedTest
//	callback.funcName = utils.GetSelfFuncName()
//	var idList []string
//	idList = append(idList, "18567155635")
//	jsontest, _ := json.Marshal(idList)
//	open_im_sdk.GetConversationRecvMessageOpt(&callback, string(jsontest))
//	fmt.Println("GetConversationRecvMessageOpt", string(jsontest))
//}

//func DoTestGetHistoryMessage(userID string) {
//	var testGetHistoryCallBack GetHistoryCallBack
//	open_im_sdk.GetHistoryMessageList(testGetHistoryCallBack, utils.structToJsonString(&utils.PullMsgReq{
//		UserID: userID,
//		Count:  50,
//	}))
//}
//func DoTestDeleteConversation(conversationID string) {
//	var testDeleteConversation DeleteConversationCallBack
//	open_im_sdk.DeleteConversation(conversationID, testDeleteConversation)
//
//}

type DeleteConversationCallBack struct {
}

func (d DeleteConversationCallBack) OnError(errCode int32, errMsg string) {
	fmt.Printf("DeleteConversationCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (d DeleteConversationCallBack) OnSuccess(data string) {
	fmt.Printf("DeleteConversationCallBack , success,data:%v\n", data)
}

type DeleteMessageFromLocalStorageCallBack struct {
}

func (d DeleteMessageFromLocalStorageCallBack) OnError(errCode int32, errMsg string) {
	fmt.Printf("DeleteMessageFromLocalStorageCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (d DeleteMessageFromLocalStorageCallBack) OnSuccess(data string) {
	fmt.Printf("DeleteMessageFromLocalStorageCallBack , success,data:%v\n", data)
}

type TestGetAllConversationListCallBack struct {
}

func (t TestGetAllConversationListCallBack) OnError(errCode int32, errMsg string) {
	fmt.Printf("TestGetAllConversationListCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (t TestGetAllConversationListCallBack) OnSuccess(data string) {
	fmt.Printf("TestGetAllConversationListCallBack , success,data:%v\n", data)
}

type TestGetOneConversationCallBack struct {
}

func (t TestGetOneConversationCallBack) OnError(errCode int32, errMsg string) {
	fmt.Printf("TestGetOneConversationCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (t TestGetOneConversationCallBack) OnSuccess(data string) {
	fmt.Printf("TestGetOneConversationCallBack , success,data:%v\n", data)
}

//func DoTestGetOneConversation(sourceID string, sessionType int) {
//	var test TestGetOneConversationCallBack
//	//GetOneConversation(Friend_uid, SingleChatType, test)
//	open_im_sdk.GetOneConversation(sourceID, sessionType, test)
//
//}
//func DoTestCreateImageMessage(path string) string {
//	return open_im_sdk.CreateImageMessage(path)
//}
//func DoTestSetConversationDraft() {
//	var test TestSetConversationDraft
//	open_im_sdk.SetConversationDraft("single_c93bc8b171cce7b9d1befb389abfe52f", "hah", test)
//
//}

type TestSetConversationDraft struct {
}

func (t TestSetConversationDraft) OnError(errCode int32, errMsg string) {
	fmt.Printf("SetConversationDraft , OnError %v\n", errMsg)
}

func (t TestSetConversationDraft) OnSuccess(data string) {
	fmt.Printf("SetConversationDraft , OnSuccess %v\n", data)
}

type GetHistoryCallBack struct {
}

func (g GetHistoryCallBack) OnError(errCode int32, errMsg string) {
	fmt.Printf("GetHistoryCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (g GetHistoryCallBack) OnSuccess(data string) {
	fmt.Printf("get History , OnSuccessData: %v\n", data)
}

type MsgListenerCallBak struct {
}

func (m MsgListenerCallBak) OnRecvNewMessage(msg string) {
	var mm sdk_struct.MsgStruct
	err := json.Unmarshal([]byte(msg), &mm)
	if err != nil {
		fmt.Println("Unmarshal failed")
	} else {
		fmt.Println("test_openim: ", "recv time: ", time.Now().UnixNano(), "send time: ", mm.SendTime, " msgid: ", mm.ClientMsgID)
	}

}
func (m MsgListenerCallBak) OnRecvC2CReadReceipt(data string) {
	fmt.Println("OnRecvC2CReadReceipt , ", data)
}

func (m MsgListenerCallBak) OnRecvMessageRevoked(msgId string) {
	fmt.Println("OnRecvMessageRevoked ", msgId)
}

type conversationCallBack struct {
}

func (c conversationCallBack) OnSyncServerStart() {
	panic("implement me")
}

func (c conversationCallBack) OnSyncServerFinish() {
	panic("implement me")
}

func (c conversationCallBack) OnSyncServerFailed() {
	panic("implement me")
}

func (c conversationCallBack) OnNewConversation(conversationList string) {
	fmt.Printf("OnNewConversation returnList is %s\n", conversationList)
}

func (c conversationCallBack) OnConversationChanged(conversationList string) {
	fmt.Printf("OnConversationChanged returnList is %s\n", conversationList)
}

func (c conversationCallBack) OnTotalUnreadMessageCountChanged(totalUnreadCount int64) {
	fmt.Printf("OnTotalUnreadMessageCountChanged returnTotalUnreadCount is %d\n", totalUnreadCount)
}

type testMarkC2CMessageAsRead struct {
}

func (testMarkC2CMessageAsRead) OnSuccess(data string) {
	fmt.Println(" testMarkC2CMessageAsRead  OnSuccess", data)
}

func (testMarkC2CMessageAsRead) OnError(code int32, msg string) {
	fmt.Println("testMarkC2CMessageAsRead, OnError", code, msg)
}

//func DoTestMarkC2CMessageAsRead() {
//	var test testMarkC2CMessageAsRead
//	readid := "2021-06-23 12:25:36-7eefe8fc74afd7c6adae6d0bc76929e90074d5bc-8522589345510912161"
//	var xlist []string
//	xlist = append(xlist, readid)
//	jsonid, _ := json.Marshal(xlist)
//	open_im_sdk.MarkC2CMessageAsRead(test, Friend_uid, string(jsonid))
//}
