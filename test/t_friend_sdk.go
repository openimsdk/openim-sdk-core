package test

import (
	"encoding/json"
	"fmt"
	//"gorm.io/gorm/callbacks"
	X "log"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
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

var Friend_uid = "openIM101"

///////////////////////////////////////////////////////////

type testGetFriendApplicationList struct {
	baseCallback
}

func DoTestGetFriendApplicationList() {
	var test testGetFriendApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ")
	open_im_sdk.GetRecvFriendApplicationList(test, test.OperationID)

}

//////////////////////////////////////////////////////////
//type testSetSelfInfo struct {
//	open_im_sdk.ui2UpdateUserInfo
//}
//
//func (testSetSelfInfo) OnSuccess(string) {
//	fmt.Println("testSetSelfInfo, OnSuccess")
//}
//
//func (testSetSelfInfo) OnError(code int32, msg string) {
//	fmt.Println("testSetSelfInfo, OnError, ", code, msg)
//}
//
//func DoTestSetSelfInfo() {
//	var test testSetSelfInfo
//	test.Name = "skkkkkkkk"
//	test.Email = "3333333@qq.com"
//	jsontest, _ := json.Marshal(test)
//	fmt.Println("SetSelfInfo, input: ", string(jsontest))
//	open_im_sdk.SetSelfInfo(string(jsontest), test)
//}

/////////////////////////////////////////////////////////
//type testGetUsersInfo struct {
//	open_im_sdk.ui2ClientCommonReq
//}
//
//func (testGetUsersInfo) OnSuccess(data string) {
//	fmt.Println("testGetUsersInfo, OnSuccess, output: ", data)
//}
//
//func (testGetUsersInfo) OnError(code int32, msg string) {
//	fmt.Println("testGetUsersInfo, OnError, ", code, msg)
//}
//
//func DoTestGetUsersInfo() {
//	var test testGetUsersInfo
//	test.UidList = append(test.UidList, Friend_uid)
//	jsontest, _ := json.Marshal(test.UidList)
//	fmt.Println("testGetUsersInfo, input: ", string(jsontest))
//	open_im_sdk.GetUsersInfo(string(jsontest), test)
//}

/////////////////////////////////////////////////////////
type testGetFriendsInfo struct {
	uid []string //`json:"uidList"`
}

func (testGetFriendsInfo) OnSuccess(data string) {
	fmt.Println("DoTestGetDesignatedFriendsInfo, OnSuccess, output: ", data)
}

func (testGetFriendsInfo) OnError(code int32, msg string) {
	fmt.Println("DoTestGetDesignatedFriendsInfo, OnError, ", code, msg)
}

func DoTestGetDesignatedFriendsInfo() {
	var test testGetFriendsInfo
	test.uid = append(test.uid, Friend_uid)

	jsontest, _ := json.Marshal(test.uid)
	fmt.Println("testGetFriendsInfo, input: ", string(jsontest))
	open_im_sdk.GetDesignatedFriendsInfo(test, string(jsontest), "asdffdsfasdfa")
}

///////////////////////////////////////////////////////
//
type testAddToBlackList struct {
	OperationID string
}

func (t testAddToBlackList) OnSuccess(string) {
	log.Info(t.OperationID, "testAddToBlackList, OnSuccess")
}

func (t testAddToBlackList) OnError(code int32, msg string) {
	log.Info(t.OperationID, "testAddToBlackList, OnError, ", code, msg)
}

func DoTestAddToBlackList() {
	var test testAddToBlackList
	test.OperationID = utils.OperationIDGenerator()
	fmt.Println("AddToBlackList, input: ", Friend_uid)
	open_im_sdk.AddBlack(test, test.OperationID, Friend_uid)
}

///////////////////////////////////////////////////////
type testDeleteFromBlackList struct {
	delUid string
}

func (testDeleteFromBlackList) OnSuccess(string) {
	fmt.Println("testDeleteFromBlackList, OnSuccess")
}

func (testDeleteFromBlackList) OnError(code int32, msg string) {
	fmt.Println("testDeleteFromBlackList, OnError, ", code, msg)
}

func DoTestDeleteFromBlackList() {
	var test testDeleteFromBlackList
	test.delUid = Friend_uid
	jsontest, _ := json.Marshal(test.delUid)
	fmt.Println("DeleteFromBlackList, input: ", string(jsontest))
	open_im_sdk.RemoveBlack(test, string(jsontest), "11111111111111asdf11112134dfsa")
}

//////////////////////////////////////////////////////
type testGetBlackList struct {
	OperationID string
}

func (t testGetBlackList) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetBlackList, OnSuccess, output: ", data)
}
func (t testGetBlackList) OnError(code int32, msg string) {
	log.Info(t.OperationID, "testGetBlackList, OnError, ", code, msg)
}
func DoTestGetBlackList() {
	var test testGetBlackList
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.GetBlackList(test, test.OperationID)
}

//////////////////////////////////////////////////////

type testCheckFriend struct {
	OperationID string
}

func (t testCheckFriend) OnSuccess(data string) {
	log.Info(t.OperationID, "testCheckFriend, OnSuccess, output: ", data)
}
func (t testCheckFriend) OnError(code int32, msg string) {
	log.Info(t.OperationID, "testCheckFriend, OnError, ", code, msg)
}
func DoTestCheckFriend() {
	var test testCheckFriend
	test.OperationID = utils.OperationIDGenerator()
	userIDList := []string{"openIM100"}
	list := utils.StructToJsonString(userIDList)
	fmt.Println("CheckFriend, input: ", list)
	open_im_sdk.CheckFriend(test, test.OperationID, list)
}

///////////////////////////////////////////////////////////
type testSetFriendRemark struct {
	baseCallback
}

func DotestSetFriendRemark() {
	var test testSetFriendRemark
	test.OperationID = utils.OperationIDGenerator()

	var param sdk_params_callback.SetFriendRemarkParams
	param.ToUserID = Friend_uid
	param.Remark = "u-gordon"
	jsontest := utils.StructToJsonString(param)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ", jsontest)
	open_im_sdk.SetFriendRemark(test, test.OperationID, jsontest)
}

/////////////////////
////////////////////////////////////////////////////////

type testDeleteFriend struct {
	baseCallback
}

func DotestDeleteFriend() {
	var test testDeleteFriend
	test.OperationID = utils.OperationIDGenerator()

	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ")
	open_im_sdk.DeleteFriend(test, test.OperationID, Friend_uid)
}

///////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
type testaddFriend struct {
	OperationID string
}

func (t testaddFriend) OnSuccess(data string) {
	log.Info(t.OperationID, "testaddFriend, OnSuccess", data)
}
func (t testaddFriend) OnError(code int32, msg string) {
	log.Info(t.OperationID, "testaddFriend, OnError", code, msg)
}

func DoTestAddFriend() {
	var test testaddFriend
	test.OperationID = utils.OperationIDGenerator()
	params := sdk_params_callback.AddFriendParams{
		ToUserID: Friend_uid,
		ReqMsg:   "23232",
	}
	jsontestaddFriend := utils.StructToJsonString(params)
	log.Info(test.OperationID, "addFriend input:", jsontestaddFriend)
	open_im_sdk.AddFriend(test, test.OperationID, jsontestaddFriend)
}

////////////////////////////////////////////////////////////////////

/////////////////////////////////////////////////////////
type testGetSendFriendApplicationList struct {
	OperationID string
}

func (t testGetSendFriendApplicationList) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetSendFriendApplicationList, OnSuccess", data)
}
func (t testGetSendFriendApplicationList) OnError(code int32, msg string) {
	log.Info(t.OperationID, "testGetSendFriendApplicationList, OnError", code, msg)
}

func DoTestGetSendFriendApplicationList() {
	var test testGetSendFriendApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, "GetSendFriendApplicationList input:")
	open_im_sdk.GetSendFriendApplicationList(test, test.OperationID)
}

////////////////////////////////////////////////////////////////////

type testGetFriendList struct {
	baseCallback
}

func DotestGetFriendList() {
	var test testGetFriendList
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ")
	open_im_sdk.GetFriendList(test, test.OperationID)
}

/////////////////////////////////////////////////////////////////////

type testAcceptFriendApplication struct {
	baseCallback
}

func DoTestAcceptFriendApplication() {
	var test testAcceptFriendApplication
	test.OperationID = utils.OperationIDGenerator()
	var param sdk_params_callback.ProcessFriendApplicationParams
	param.HandleMsg = "ok ok "
	param.ToUserID = Friend_uid
	input := utils.StructToJsonString(param)
	open_im_sdk.AcceptFriendApplication(test, test.OperationID, input)
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

//type testRefuseFriendApplication struct {
//	open_im_sdk.ui2AcceptFriend
//}
//
//func (testRefuseFriendApplication) OnSuccess(info string) {
//	fmt.Println("testRefuseFriendApplication OnSuccess", info)
//}
//func (testRefuseFriendApplication) OnError(code int32, msg string) {
//	fmt.Println("testRefuseFriendApplication, OnError, ", code, msg)
//}
//func DoTestRefuseFriendApplication() {
//	var testRefuseFriendApplication testRefuseFriendApplication
//	testRefuseFriendApplication.ui2AcceptFriend = Friend_uid
//
//	jsontestfusetFriendappclicatrion, _ := json.Marshal(testRefuseFriendApplication.UID)
//	open_im_sdk.RefuseFriendApplication(testRefuseFriendApplication, string(jsontestfusetFriendappclicatrion), "")
//	fmt.Println("RefuseFriendApplication, input: ", string(jsontestfusetFriendappclicatrion))
//}

////////////////////////////////////////////////////////////////////

type BaseSuccFailed struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
}

func (b *BaseSuccFailed) OnError(errCode int32, errMsg string) {
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
	callback.funcName = utils.RunFuncName()
	operationID := utils.OperationIDGenerator()
	open_im_sdk.Login(&callback, uid, operationID, tk)
}

func InOutLogou() {
	var callback BaseSuccFailed
	callback.funcName = utils.RunFuncName()
	opretaionID := utils.OperationIDGenerator()
	open_im_sdk.Logout(&callback, opretaionID)
}

func InOutDoTest(uid, tk, ws, api string) {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api

	cf.WsAddr = ws
	cf.Platform = 1
	cf.DataDir = "./"
	cf.LogLevel = 6

	var s string
	b, _ := json.Marshal(cf)
	s = string(b)
	fmt.Println(s)
	var testinit testInitLister
	operationID := utils.OperationIDGenerator()
	open_im_sdk.InitSDK(s, operationID, testinit)

	var testConversation conversationCallBack
	open_im_sdk.SetConversationListener(testConversation)

	var testUser userCallback
	open_im_sdk.SetUserListener(testUser)

	//var msgCallBack MsgListenerCallBak
	//open_im_sdk.AddAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	open_im_sdk.SetFriendListener(friendListener)

	var groupListener testGroupListener
	open_im_sdk.SetGroupListener(groupListener)

	InOutlllogin(uid, tk)
	time.Sleep(5 * time.Second)
}

func lllogin(uid, tk string) bool {
	var callback BaseSuccFailed
	callback.funcName = utils.RunFuncName()
	operationID := utils.OperationIDGenerator()
	open_im_sdk.Login(&callback, uid, operationID, tk)

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
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api // "http://120.24.45.199:10000"
	//	cf.IpWsAddr = "wss://open-im.rentsoft.cn/wss"
	cf.WsAddr = ws //"ws://120.24.45.199:17778"
	cf.Platform = 2
	cf.DataDir = "./"

	var s string
	b, _ := json.Marshal(cf)
	s = string(b)
	fmt.Println(s)
	var testinit testInitLister
	operationID := utils.OperationIDGenerator()
	open_im_sdk.InitSDK(s, operationID, testinit)

	var testConversation conversationCallBack
	open_im_sdk.SetConversationListener(testConversation)

	var testUser userCallback
	open_im_sdk.SetUserListener(testUser)

	//var msgCallBack MsgListenerCallBak
	//open_im_sdk.AddAdvancedMsgListener(msgCallBack)

	var friendListener testFriendListener
	open_im_sdk.SetFriendListener(friendListener)

	var groupListener testGroupListener
	open_im_sdk.SetGroupListener(groupListener)

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

func (t *TestSendMsgCallBack) OnError(errCode int32, errMsg string) {
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

func (b *BaseSuccFailedTest) OnError(errCode int32, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	fmt.Println("22onError ", b.funcName, errCode, errMsg)
}

func (b *BaseSuccFailedTest) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	fmt.Println("22OnSuccess: ", b.funcName, data)
}

//func DotestSetConversationRecvMessageOpt() {
//	var callback BaseSuccFailedTest
//	callback.funcName = utils.RunFuncName()
//	var idList []string
//	idList = append(idList, "18567155635")
//	jsontest, _ := json.Marshal(idList)
//	open_im_sdk.SetConversationRecvMessageOpt(&callback, string(jsontest), 2)
//	fmt.Println("SetConversationRecvMessageOpt", string(jsontest))
//}
//
//func DoTestGetMultipleConversation() {
//	var callback BaseSuccFailedTest
//	callback.funcName = utils.RunFuncName()
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
//	callback.funcName = utils.RunFuncName()
//	var idList []string
//	idList = append(idList, "18567155635")
//	jsontest, _ := json.Marshal(idList)
//	open_im_sdk.GetConversationRecvMessageOpt(&callback, string(jsontest))
//	fmt.Println("GetConversationRecvMessageOpt", string(jsontest))
//}

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

//func DoTestGetAllConversationList() {
//	var test TestGetAllConversationListCallBack
//	open_im_sdk.GetAllConversationList(test)
//}

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

type userCallback struct {
}

func (userCallback) OnSelfInfoUpdated(callbackData string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackData)
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

func (testInitLister) OnConnectFailed(ErrCode int32, ErrMsg string) {
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

func (testInitLister) OnError(code int32, msg string) {
	fmt.Println("testInitLister, OnError", code, msg)
}

type testLogin struct {
}

func (testLogin) OnSuccess(string) {
	fmt.Println("testLogin OnSuccess")
}

func (testLogin) OnError(code int32, msg string) {
	fmt.Println("testLogin, OnError", code, msg)
}

type testFriendListener struct {
	x int
}

//OnFriendApplicationAdded(friendApplication string)
//	OnFriendApplicationDeleted(friendApplication string)
//	OnFriendApplicationAccepted(groupApplication string)
//	OnFriendApplicationRejected(friendApplication string)
//	OnFriendAdded(friendInfo string)
//	OnFriendDeleted(friendInfo string)
//	OnFriendInfoChanged(friendInfo string)
//	OnBlackAdded(blackInfo string)
//	OnBlackDeleted(blackInfo string)

func (testFriendListener) OnFriendApplicationAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}
func (testFriendListener) OnFriendApplicationDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendApplicationAccepted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendApplicationRejected(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnBlackAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}
func (testFriendListener) OnBlackDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnFriendInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)
}

func (testFriendListener) OnSuccess() {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName())
}

func (testFriendListener) OnError(code int32, msg string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), code, msg)
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
