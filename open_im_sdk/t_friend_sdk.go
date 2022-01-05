package open_im_sdk

import (
	"encoding/json"
	"fmt"
	X "log"
	"os"
	"runtime"
	"time"
)

var loggerf *X.Logger

func init() {
	loggerf = X.New(os.Stdout, "", X.Llongfile|X.Ltime|X.Ldate)
}

type TestSendImg struct {
}

func (TestSendImg) OnSuccess(data string) {
	fmt.Println("testSendImg, OnSuccess, output: ", data)
}

func (TestSendImg) OnError(code int, msg string) {
	fmt.Println("testSendImg, OnError, ", code, msg)
}

func (TestSendImg) OnProgress(progress int) {
	fmt.Println("progress: ", progress)
}

func TestLog(v ...interface{}) {
	//X.SetFlags(X.Lshortfile | X.LstdFlags)
	loggerf.Println(v)
	a, b, c, d := runtime.Caller(1)
	X.Println(a, b, c, d)
}

var Friend_uid = "18349115126"

///////////////////////////////////////////////////////////

//GetFriendApplicationList

type testGetFriendApplicationList struct {
}

func (testGetFriendApplicationList) OnSuccess(data string) {
	fmt.Println("testGetFriendApplicationList, OnSuccess, output:", data)
}

func (testGetFriendApplicationList) OnError(code int32, msg string) {
	fmt.Println("testGetFriendApplicationList, OnError, ", code, msg)
}

func DoTestGetFriendApplicationList() {
	var test testGetFriendApplicationList
	GetRecvFriendApplicationList(test, "")

}

//////////////////////////////////////////////////////////
type testSetSelfInfo struct {
	ui2UpdateUserInfo
}

func (testSetSelfInfo) OnSuccess(string) {
	fmt.Println("testSetSelfInfo, OnSuccess")
}

func (testSetSelfInfo) OnError(code int32, msg string) {
	fmt.Println("testSetSelfInfo, OnError, ", code, msg)
}

func DoTestSetSelfInfo() {
	var test testSetSelfInfo
	test.Name = "skkkkkkkk"
	test.Email = "3333333@qq.com"
	jsontest, _ := json.Marshal(test)
	fmt.Println("SetSelfInfo, input: ", string(jsontest))
	SetSelfInfo(string(jsontest), test)
}

/////////////////////////////////////////////////////////
type testGetUsersInfo struct {
	ui2ClientCommonReq
}

func (testGetUsersInfo) OnSuccess(data string) {
	fmt.Println("testGetUsersInfo, OnSuccess, output: ", data)
}

func (testGetUsersInfo) OnError(code int32, msg string) {
	fmt.Println("testGetUsersInfo, OnError, ", code, msg)
}

func DoTestGetUsersInfo() {
	var test testGetUsersInfo
	test.UidList = append(test.UidList, Friend_uid)
	jsontest, _ := json.Marshal(test.UidList)
	fmt.Println("testGetUsersInfo, input: ", string(jsontest))
	GetUsersInfo(string(jsontest), test)
}

/////////////////////////////////////////////////////////
type testGetFriendsInfo struct {
	uid []string `json:"uidList"`
}

func (testGetFriendsInfo) OnSuccess(data string) {
	fmt.Println("testGetFriendsInfo, OnSuccess, output: ", data)
}

func (testGetFriendsInfo) OnError(code int32, msg string) {
	fmt.Println("testGetFriendsInfo, OnError, ", code, msg)
}

func DoTestGetFriendsInfo() {
	var test testGetFriendsInfo
	test.uid = append(test.uid, Friend_uid)
	jsontest, _ := json.Marshal(test.uid)
	fmt.Println("testGetFriendsInfo, input: ", string(jsontest))
	GetFriendsInfo(test, string(jsontest), "")
}

///////////////////////////////////////////////////////

type testAddToBlackList struct {
	delUid
}

func (testAddToBlackList) OnSuccess(string) {
	fmt.Println("testAddToBlackList, OnSuccess")
}

func (testAddToBlackList) OnError(code int, msg string) {
	fmt.Println("testAddToBlackList, OnError, ", code, msg)
}

func DoTestAddToBlackList() {
	var test testAddToBlackList
	test.Uid = Friend_uid
	jsontest, _ := json.Marshal(test.Uid)
	fmt.Println("AddToBlackList, input: ", string(jsontest))
	AddToBlackList(test, string(jsontest), "")
}

///////////////////////////////////////////////////////
type testDeleteFromBlackList struct {
	delUid string
}

func (testDeleteFromBlackList) OnSuccess(string) {
	fmt.Println("testDeleteFromBlackList, OnSuccess")
}

func (testDeleteFromBlackList) OnError(code int, msg string) {
	fmt.Println("testDeleteFromBlackList, OnError, ", code, msg)
}

func DoTestDeleteFromBlackList() {
	var test testDeleteFromBlackList
	test.delUid = Friend_uid
	jsontest, _ := json.Marshal(test.delUid)
	fmt.Println("DeleteFromBlackList, input: ", string(jsontest))
	DeleteBlack(test, string(jsontest), "")
}

//////////////////////////////////////////////////////
type testGetBlackList struct {
}

func (testGetBlackList) OnSuccess(data string) {
	fmt.Println("testGetBlackList, OnSuccess, output: ", data)
}
func (testGetBlackList) OnError(code int, msg string) {
	fmt.Println("testGetBlackList, OnError, ", code, msg)
}
func DoTestGetBlackList() {
	var test testGetBlackList
	GetBlackList(test)
}

//////////////////////////////////////////////////////
type testCheckFriend struct {
	ui2ClientCommonReq
}

func (testCheckFriend) OnSuccess(data string) {
	fmt.Println("testCheckFriend, OnSuccess, output: ", data)
}
func (testCheckFriend) OnError(code int, msg string) {
	fmt.Println("testCheckFriend, OnError, ", code, msg)
}
func DoTestCheckFriend() {
	var test testCheckFriend
	test.UidList = append(test.UidList, Friend_uid)
	jsontest, _ := json.Marshal(test.UidList)
	fmt.Println("CheckFriend, input: ", string(jsontest))
	CheckFriend(test, string(jsontest))
}

/////////////////////////////////////////////////////////
type testSetFriendInfo struct {
	uid2Comment
}

func (testSetFriendInfo) OnSuccess(string) {
	fmt.Println("testSetFriendInfo, OnSucess")
}
func (testSetFriendInfo) OnError(code int, msg string) {
	fmt.Println("testSetFriendInfo, OnError, ", code, msg)
}
func DoTestSetFriendInfo() {
	var test testSetFriendInfo
	test.Uid = Friend_uid
	test.Comment = "MM"
	jsontest, _ := json.Marshal(test)
	fmt.Println("SetFriendInfo, input: ", string(jsontest))
	SetFriendInfo(string(jsontest), test, "")
}

/////////////////////
////////////////////////////////////////////////////////

type TestDeleteFromFriendList struct {
	Uid string `json:"uid"`
}

func (TestDeleteFromFriendList) OnSuccess(string) {
	fmt.Println("testDeleteFromFriendList,  OnSuccess")
}

func (TestDeleteFromFriendList) OnError(code int, msg string) {
	fmt.Println("testDeleteFromFriendList, OnError, ", code, msg)
}

func DoTestDeleteFromFriendList() {
	var test TestDeleteFromFriendList
	test.Uid = Friend_uid
	jsontest, err := json.Marshal(test.Uid)
	fmt.Println("DeleteFromFriendList, input:              sdafasf ", string(jsontest), err)
	DeleteFromFriendList(string(jsontest), test.Uid, "")
}

///////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
type testaddFriend struct {
	UID        string `json:"uid" binding:"required"`
	ReqMessage string `json:"reqMessage"`
}

func (testaddFriend) OnSuccess(data string) {
	fmt.Println("testaddFriend, OnSuccess", data)
}
func (testaddFriend) OnError(code int32, msg string) {
	fmt.Println("testaddFriend, OnError", code, msg)
}

func DoTestaddFriend() {
	var testaddFriend testaddFriend

	testaddFriend.UID = Friend_uid
	testaddFriend.ReqMessage = "hello"

	jsontestaddFriend, _ := json.Marshal(testaddFriend)
	fmt.Println("addFriend input:", string(jsontestaddFriend))
	AddFriend(testaddFriend, string(jsontestaddFriend), "")
}

/////////////////////////////////////////////////////////////////////

type testGetFriendList struct {
}

func (testGetFriendList) OnSuccess(list string) {
	fmt.Println("testGetFriendList OnSuccess output: ", list)
}
func (testGetFriendList) OnError(code int, msg string) {
	fmt.Println("testGetFriendList, OnError, ", code, msg)
}
func DoTestGetFriendList() {
	var testGetFriendList testGetFriendList
	GetFriendList(testGetFriendList, "")
}

/////////////////////////////////////////////////////////////////////

type testAcceptFriendApplication struct {
	ui2AcceptFriend
}

func (testAcceptFriendApplication) OnSuccess(info string) {
	fmt.Println("testAcceptFriendApplication OnSuccess", info)
}
func (testAcceptFriendApplication) OnError(code int, msg string) {
	fmt.Println("testAcceptFriendApplication, OnError, ", code, msg)
}

func DoTestAcceptFriendApplication() {
	var testAcceptFriendApplication testAcceptFriendApplication
	testAcceptFriendApplication.UID = Friend_uid

	jsontestAcceptFriendappclicatrion, _ := json.Marshal(testAcceptFriendApplication.UID)
	AcceptFriendApplication(testAcceptFriendApplication, string(jsontestAcceptFriendappclicatrion), "")
	fmt.Println("AcceptFriendApplication, input: ", string(jsontestAcceptFriendappclicatrion))
}

/*
type testRefuseFriendApplication struct {
	ui2AcceptFriend
}

func (testRefuseFriendApplication) OnSuccess(info string) {
	fmt.Println("RefuseFriendApplication OnSuccess", info)
}
func (testRefuseFriendApplication) OnError(code int, msg string) {
	fmt.Println("RefuseFriendApplication, OnError, ", code, msg)
}
*/
/*
func DoTestRefuseFriendApplication() {

	var test testRefuseFriendApplication
	test.UID = Friend_uid

	js, _ := json.Marshal(test)
	RefuseFriendApplication(test, string(js))
	fmt.Println("RefuseFriendApplication, input: ", string(js))
}


*/

/////////////////////////////////////////////////////////////////////

type testRefuseFriendApplication struct {
	ui2AcceptFriend
}

func (testRefuseFriendApplication) OnSuccess(info string) {
	fmt.Println("testRefuseFriendApplication OnSuccess", info)
}
func (testRefuseFriendApplication) OnError(code int, msg string) {
	fmt.Println("testRefuseFriendApplication, OnError, ", code, msg)
}
func DoTestRefuseFriendApplication() {
	var testRefuseFriendApplication testRefuseFriendApplication
	testRefuseFriendApplication.UID = Friend_uid

	jsontestfusetFriendappclicatrion, _ := json.Marshal(testRefuseFriendApplication.UID)
	RefuseFriendApplication(testRefuseFriendApplication, string(jsontestfusetFriendappclicatrion), "")
	fmt.Println("RefuseFriendApplication, input: ", string(jsontestfusetFriendappclicatrion))
}

////////////////////////////////////////////////////////////////////

type BaseSuccFailed struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
}

func (b *BaseSuccFailed) OnError(errCode int, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	fmt.Println("onError ", b.funcName)
	fmt.Println("test_openim: ", "login failed ", errCode, errMsg)

}

func (b *BaseSuccFailed) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	fmt.Println("test_openim: ", "login success")
}

func InOutlllogin(uid, tk string) {
	var callback BaseSuccFailed
	callback.funcName = RunFuncName()
	Login(uid, tk, &callback)
}

func InOutLogou() {
	var callback BaseSuccFailed
	callback.funcName = RunFuncName()
	Logout(&callback)
}

func InOutDoTest(uid, tk, ws, api string) {
	var cf IMConfig
	cf.IpApiAddr = api

	cf.IpWsAddr = ws
	cf.Platform = 1
	cf.DbDir = "./"

	var s string
	b, _ := json.Marshal(cf)
	s = string(b)
	fmt.Println(s)
	var testinit testInitLister
	InitSDK(s, testinit)

	var testConversation conversationCallBack
	SetConversationListener(testConversation)

	var msgCallBack MsgListenerCallBak
	AddAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	SetFriendListener(friendListener)

	var groupListener testGroupListener
	SetGroupListener(groupListener)

	InOutlllogin(uid, tk)
}

func lllogin(uid, tk string) bool {
	var callback BaseSuccFailed
	callback.funcName = RunFuncName()
	Login(uid, tk, &callback)

	for true {
		if callback.errCode == 1 {
			fmt.Println("success code 1")
			return true
		} else if callback.errCode == -1 {
			fmt.Println("failed code -1")
			return false
		} else {
			fmt.Println("code sleep")
			time.Sleep(1 * time.Second)
			continue
		}
	}
	return true
}

func DoTest(uid, tk, ws, api string) {
	var cf IMConfig
	cf.IpApiAddr = api // "http://120.24.45.199:10000"
	//	cf.IpWsAddr = "wss://open-im.rentsoft.cn/wss"
	cf.IpWsAddr = ws //"ws://120.24.45.199:17778"
	cf.Platform = 2
	cf.DbDir = "./"

	var s string
	b, _ := json.Marshal(cf)
	s = string(b)
	fmt.Println(s)
	var testinit testInitLister
	InitSDK(s, testinit)

	var testConversation conversationCallBack
	SetConversationListener(testConversation)

	var msgCallBack MsgListenerCallBak
	AddAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	SetFriendListener(friendListener)

	var groupListener testGroupListener
	SetGroupListener(groupListener)

	time.Sleep(1 * time.Second)

	for !lllogin(uid, tk) {
		fmt.Println("lllogin, failed, login...")
		time.Sleep(1 * time.Second)
	}

}

////////////////////////////////////////////////////////////////////
type TestSendMsgCallBack struct {
	msg string
}

func (t *TestSendMsgCallBack) OnError(errCode int, errMsg string) {
	fmt.Println("test_openim: send msg failed: ", errCode, errMsg, "|", t.msg, "|")
}

func (t *TestSendMsgCallBack) OnSuccess(data string) {
	fmt.Println("test_openim: send msg success: |", t.msg, "|")
}

func (t *TestSendMsgCallBack) OnProgress(progress int) {
	//	fmt.Printf("msg_send , onProgress %d\n", progress)
}

type BaseSuccFailedTest struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
}

func (b *BaseSuccFailedTest) OnError(errCode int, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	fmt.Println("22onError ", b.funcName, errCode, errMsg)
}

func (b *BaseSuccFailedTest) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	fmt.Println("22OnSuccess: ", b.funcName, data)
}

func DotestSetConversationRecvMessageOpt() {
	var callback BaseSuccFailedTest
	callback.funcName = RunFuncName()
	var idList []string
	idList = append(idList, "18567155635")
	jsontest, _ := json.Marshal(idList)
	SetConversationRecvMessageOpt(&callback, string(jsontest), 2)
	fmt.Println("SetConversationRecvMessageOpt", string(jsontest))
}

func DoTestGetConversationRecvMessageOpt() {
	var callback BaseSuccFailedTest
	callback.funcName = RunFuncName()
	var idList []string
	idList = append(idList, "18567155635")
	jsontest, _ := json.Marshal(idList)
	GetConversationRecvMessageOpt(&callback, string(jsontest))
	fmt.Println("GetConversationRecvMessageOpt", string(jsontest))
}

func InOutDoTestSendMsg(sendId, receiverID string) {
	m := "test:" + sendId + ":" + receiverID + ":"
	//s := CreateTextMessage(m)
	var testSendMsg TestSendMsgCallBack
	//	testSendMsg.msg = SendMessage(&testSendMsg, s, receiverID, "", false)
	fmt.Println("func send ", m, testSendMsg.msg)
	fmt.Println("test to recv : ", receiverID)
}

func DoTestSendMsg(sendId, receiverID string, idx string) {
	m := "test:" + sendId + ":" + receiverID + ":" + idx
	//s := CreateTextMessage(m)
	var testSendMsg TestSendMsgCallBack
	//	testSendMsg.msg = SendMessage(&testSendMsg, s, receiverID, "", false)
	fmt.Println("func send ", m, testSendMsg.msg)
	fmt.Println("test to recv : ", receiverID)
}

func DoTestGetHistoryMessage(userID string) {
	var testGetHistoryCallBack GetHistoryCallBack
	GetHistoryMessageList(testGetHistoryCallBack, structToJsonString(&PullMsgReq{
		UserID: userID,
		Count:  50,
	}))
}
func DoTestDeleteConversation(conversationID string) {
	var testDeleteConversation DeleteConversationCallBack
	DeleteConversation(conversationID, testDeleteConversation)

}

type DeleteConversationCallBack struct {
}

func (d DeleteConversationCallBack) OnError(errCode int, errMsg string) {
	fmt.Printf("DeleteConversationCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (d DeleteConversationCallBack) OnSuccess(data string) {
	fmt.Printf("DeleteConversationCallBack , success,data:%v\n", data)
}

type DeleteMessageFromLocalStorageCallBack struct {
}

func (d DeleteMessageFromLocalStorageCallBack) OnError(errCode int, errMsg string) {
	fmt.Printf("DeleteMessageFromLocalStorageCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (d DeleteMessageFromLocalStorageCallBack) OnSuccess(data string) {
	fmt.Printf("DeleteMessageFromLocalStorageCallBack , success,data:%v\n", data)
}

type TestGetAllConversationListCallBack struct {
}

func (t TestGetAllConversationListCallBack) OnError(errCode int, errMsg string) {
	fmt.Printf("TestGetAllConversationListCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (t TestGetAllConversationListCallBack) OnSuccess(data string) {
	fmt.Printf("TestGetAllConversationListCallBack , success,data:%v\n", data)
}

func DoTestGetAllConversationList() {
	var test TestGetAllConversationListCallBack
	GetAllConversationList(test)
}

type TestGetOneConversationCallBack struct {
}

func (t TestGetOneConversationCallBack) OnError(errCode int, errMsg string) {
	fmt.Printf("TestGetOneConversationCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (t TestGetOneConversationCallBack) OnSuccess(data string) {
	fmt.Printf("TestGetOneConversationCallBack , success,data:%v\n", data)
}

func DoTestGetOneConversation(sourceID string, sessionType int) {
	var test TestGetOneConversationCallBack
	//GetOneConversation(Friend_uid, SingleChatType, test)
	GetOneConversation(sourceID, sessionType, test)

}
func DoTestCreateImageMessage(path string) string {
	return CreateImageMessage(path)
}
func DoTestSetConversationDraft() {
	var test TestSetConversationDraft
	SetConversationDraft("single_c93bc8b171cce7b9d1befb389abfe52f", "hah", test)

}

type TestSetConversationDraft struct {
}

func (t TestSetConversationDraft) OnError(errCode int, errMsg string) {
	fmt.Printf("SetConversationDraft , OnError %v\n", errMsg)
}

func (t TestSetConversationDraft) OnSuccess(data string) {
	fmt.Printf("SetConversationDraft , OnSuccess %v\n", data)
}

type GetHistoryCallBack struct {
}

func (g GetHistoryCallBack) OnError(errCode int, errMsg string) {
	fmt.Printf("GetHistoryCallBack , errCode:%v,errMsg:%v\n", errCode, errMsg)
}

func (g GetHistoryCallBack) OnSuccess(data string) {
	fmt.Printf("get History , OnSuccessData: %v\n", data)
}

type MsgListenerCallBak struct {
}

func (m MsgListenerCallBak) OnRecvNewMessage(msg string) {
	var mm MsgStruct
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

func (c conversationCallBack) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	fmt.Printf("OnTotalUnreadMessageCountChanged returnTotalUnreadCount is %d\n", totalUnreadCount)
}

////////////////////////////////////////////////////////////////////
type testInitLister struct {
}

func (testInitLister) OnUserTokenExpired() {
	fmt.Println("testInitLister, OnUserTokenExpired")
}
func (testInitLister) OnConnecting() {
	fmt.Println("testInitLister, OnConnecting")
}

func (testInitLister) OnConnectSuccess() {
	fmt.Println("testInitLister, OnConnectSuccess")
}

func (testInitLister) OnConnectFailed(ErrCode int, ErrMsg string) {
	fmt.Println("testInitLister, OnConnectFailed", ErrCode, ErrMsg)
}

func (testInitLister) OnKickedOffline() {
	fmt.Println("testInitLister, OnKickedOffline")
}

func (testInitLister) OnSelfInfoUpdated(info string) {
	fmt.Println("testInitLister, OnSelfInfoUpdated, ", info)
}

func (testInitLister) OnSucess() {
	fmt.Println("testInitLister, OnSucess")
}

func (testInitLister) OnError(code int, msg string) {
	fmt.Println("testInitLister, OnError", code, msg)
}

type testLogin struct {
}

func (testLogin) OnSuccess(string) {
	fmt.Println("testLogin OnSuccess")
}

func (testLogin) OnError(code int, msg string) {
	fmt.Println("testLogin, OnError", code, msg)
}

type testFriendListener struct {
	x int
}

func (testFriendListener) OnFriendApplicationListAdded(friendAdded string) {
	fmt.Println("testFriendListener,OnFriendApplicationListAdded", friendAdded)
}
func (testFriendListener) OnFriendApplicationListDeleted(friendDeleted string) {
	fmt.Println("testFriendListener,OnFriendApplicationListDeleted", friendDeleted)
}

func (testFriendListener) OnFriendApplicationListAccept(friendAccept string) {
	fmt.Println("testFriendListener,OnFriendApplicationListAccept", friendAccept)
}

func (testFriendListener) OnFriendApplicationListReject(info string) {
	fmt.Println("testFriendListener,OnFriendApplicationListReject", info)
}

func (testFriendListener) OnFriendListAdded(info string) {
	fmt.Println("testFriendListener,OnFriendListAdded", info)
}

func (testFriendListener) OnFriendListDeleted(info string) {
	fmt.Println("testFriendListener,OnFriendListDeleted", info)
}

func (testFriendListener) OnBlackListAdd(info string) {
	fmt.Println("testFriendListener, OnBlackListAdd", info)
}
func (testFriendListener) OnBlackListDeleted(info string) {
	fmt.Println("testFriendListener, OnBlackListDeleted", info)
}

func (testFriendListener) OnFriendInfoChanged(InfoList string) {
	fmt.Println("testFriendListener, OnFriendInfoChanged")
}

func (testFriendListener) OnSuccess() {
	fmt.Println("testLogin OnSuccess")
}

func (testFriendListener) OnError(code int, msg string) {
	fmt.Println("testLogin, OnError", code, msg)
}

type testMarkC2CMessageAsRead struct {
}

func (testMarkC2CMessageAsRead) OnSuccess(data string) {
	fmt.Println(" testMarkC2CMessageAsRead  OnSuccess", data)
}

func (testMarkC2CMessageAsRead) OnError(code int, msg string) {
	fmt.Println("testMarkC2CMessageAsRead, OnError", code, msg)
}

func DoTestMarkC2CMessageAsRead() {
	var test testMarkC2CMessageAsRead
	readid := "2021-06-23 12:25:36-7eefe8fc74afd7c6adae6d0bc76929e90074d5bc-8522589345510912161"
	var xlist []string
	xlist = append(xlist, readid)
	jsonid, _ := json.Marshal(xlist)
	MarkC2CMessageAsRead(test, Friend_uid, string(jsonid))
}
