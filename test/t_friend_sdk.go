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

package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/ccontext"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
	X "log"
	"os"
	"runtime"
	"time"

	"github.com/OpenIMSDK/protocol/sdkws"

	"github.com/OpenIMSDK/tools/mcontext"
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

var Friend_uid = "3126758667"

func SetTestFriendID(friendUserID string) {
	Friend_uid = friendUserID
}

///////////////////////////////////////////////////////////

type testGetFriendApplicationList struct {
	baseCallback
}

func DoTestGetFriendApplicationList() {
	var test testGetFriendApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, test.OperationID, utils.GetSelfFuncName(), "input ")
	// open_im_sdk.GetRecvFriendApplicationList(test, test.OperationID)
}

// ////////////////////////////////////////////////////////`
type testSetSelfInfo struct {
	baseCallback
}

func DoTestSetSelfInfo() {
	var test testSetSelfInfo
	test.OperationID = utils.OperationIDGenerator()
	userInfo := sdkws.UserInfo{}
	userInfo.Nickname = "new 4444444444444 Gordon001"
	jsonString := utils.StructToJsonStringDefault(userInfo)
	fmt.Println("SetSelfInfo, input: ")
	open_im_sdk.SetSelfInfo(test, test.OperationID, jsonString)
}

// ///////////////////////////////////////////////////////
type testGetUsersInfo struct {
	baseCallback
}

func DoTestGetUsersInfo() {
	var test testGetUsersInfo
	test.OperationID = utils.OperationIDGenerator()
	userIDList := []string{"4950399653"}
	list := utils.StructToJsonStringDefault(userIDList)
	fmt.Println("testGetUsersInfo, input: ", list)
	open_im_sdk.GetUsersInfo(test, test.OperationID, list)
}

// ///////////////////////////////////////////////////////
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
	// open_im_sdk.GetDesignatedFriendsInfo(test, "xxxxxxxxxxx", string(jsontest))
}

// /////////////////////////////////////////////////////
type testAddToBlackList struct {
	OperationID string
}

func (t testAddToBlackList) OnSuccess(string) {
	log.ZInfo(ctx, t.OperationID, "testAddToBlackList, OnSuccess")
}

func (t testAddToBlackList) OnError(code int32, msg string) {
	log.ZInfo(ctx, t.OperationID, "testAddToBlackList, OnError, ", code, msg)
}

func DoTestAddToBlackList(userID string) {
	var test testAddToBlackList
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.AddBlack(test, test.OperationID, userID, "")
}

// /////////////////////////////////////////////////////
type testDeleteFromBlackList struct {
	baseCallback
}

func DoTestDeleteFromBlackList(userID string) {
	var test testDeleteFromBlackList
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.RemoveBlack(test, test.OperationID, userID)
}

// ////////////////////////////////////////////////////
type testGetBlackList struct {
	OperationID string
}

func (t testGetBlackList) OnSuccess(data string) {
	log.ZInfo(ctx, t.OperationID, "testGetBlackList, OnSuccess, output: ", data)
}
func (t testGetBlackList) OnError(code int32, msg string) {
	log.ZInfo(ctx, t.OperationID, "testGetBlackList, OnError, ", code, msg)
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
	log.ZInfo(ctx, t.OperationID, "testCheckFriend, OnSuccess, output: ", data)
}
func (t testCheckFriend) OnError(code int32, msg string) {
	log.ZInfo(ctx, t.OperationID, "testCheckFriend, OnError, ", code, msg)
}
func DoTestCheckFriend() {
	var test testCheckFriend
	test.OperationID = utils.OperationIDGenerator()
	userIDList := []string{"openIM100"}
	list := utils.StructToJsonString(userIDList)
	fmt.Println("CheckFriend, input: ", list)
	open_im_sdk.CheckFriend(test, test.OperationID, list)
}

// /////////////////////////////////////////////////////////
type testSetFriendRemark struct {
	baseCallback
}

func DotestSetFriendRemark() {
	var test testSetFriendRemark
	test.OperationID = utils.OperationIDGenerator()

	var param sdk_params_callback.SetFriendRemarkParams
	param.ToUserID = Friend_uid
	param.Remark = "4444 "
	jsontest := utils.StructToJsonString(param)
	log.ZInfo(ctx, test.OperationID, utils.GetSelfFuncName(), "input ", jsontest)
	open_im_sdk.SetFriendRemark(test, test.OperationID, jsontest)
}

/////////////////////
////////////////////////////////////////////////////////

type testDeleteFriend struct {
	baseCallback
}

func DotestDeleteFriend(userID string) {
	var test testDeleteFriend
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.DeleteFriend(test, test.OperationID, userID)
}

// /////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
type testaddFriend struct {
	OperationID string
}

func (t testaddFriend) OnSuccess(data string) {
	log.ZInfo(ctx, t.OperationID, "testaddFriend, OnSuccess", data)
}
func (t testaddFriend) OnError(code int32, msg string) {
	log.ZInfo(ctx, t.OperationID, "testaddFriend, OnError", code, msg)
}

func DoTestAddFriend() {
	var test testaddFriend
	test.OperationID = utils.OperationIDGenerator()
	params := sdk_params_callback.AddFriendParams{
		ToUserID: Friend_uid,
		ReqMsg:   "777777777777777777777777",
	}
	jsontestaddFriend := utils.StructToJsonString(params)
	log.ZInfo(ctx, test.OperationID, "addFriend input:", jsontestaddFriend)
	open_im_sdk.AddFriend(test, test.OperationID, jsontestaddFriend)
}

////////////////////////////////////////////////////////////////////

// ///////////////////////////////////////////////////////
type testGetSendFriendApplicationList struct {
	OperationID string
}

func (t testGetSendFriendApplicationList) OnSuccess(data string) {
	log.ZInfo(ctx, t.OperationID, "testGetSendFriendApplicationList, OnSuccess", data)
}
func (t testGetSendFriendApplicationList) OnError(code int32, msg string) {
	log.ZInfo(ctx, t.OperationID, "testGetSendFriendApplicationList, OnError", code, msg)
}

func DoTestGetSendFriendApplicationList() {
	var test testGetSendFriendApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, test.OperationID, "GetSendFriendApplicationList input:")
	// open_im_sdk.GetSendFriendApplicationList(test, test.OperationID)
}

////////////////////////////////////////////////////////////////////

type testGetFriendList struct {
	baseCallback
}

func DotestGetFriendList() {
	var test testGetFriendList
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, test.OperationID, utils.GetSelfFuncName(), "input ")
	open_im_sdk.GetFriendList(test, test.OperationID)
}

type testSearchFriends struct {
	baseCallback
}

func DotestSearchFriends() {
	var test testSearchFriends
	test.OperationID = utils.OperationIDGenerator()
	test.callName = "SearchFriends"
	var params sdk_params_callback.SearchFriendsParam
	params.KeywordList = []string{"G"}
	params.IsSearchUserID = true
	params.IsSearchNickname = true
	params.IsSearchRemark = true
	log.ZInfo(ctx, test.OperationID, utils.GetSelfFuncName(), "input ", params)
	open_im_sdk.SearchFriends(test, test.OperationID, utils.StructToJsonString(params))
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

type testRefuseFriendApplication struct {
	baseCallback
}

func DoTestRefuseFriendApplication() {
	var test testRefuseFriendApplication
	test.OperationID = utils.OperationIDGenerator()
	var param sdk_params_callback.ProcessFriendApplicationParams
	param.HandleMsg = "nonono"
	param.ToUserID = Friend_uid
	input := utils.StructToJsonString(param)
	open_im_sdk.RefuseFriendApplication(test, test.OperationID, input)
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

func SetListenerAndLogin(uid, tk string) {
	//
	//var testConversation conversationCallBack
	//open_im_sdk.SetConversationListener(&testConversation)
	//
	//var testUser userCallback
	//open_im_sdk.SetUserListener(testUser)
	//
	//var msgCallBack MsgListenerCallBak
	//open_im_sdk.SetAdvancedMsgListener(&msgCallBack)
	//
	//var batchMsg BatchMsg
	//open_im_sdk.SetBatchMsgListener(&batchMsg)
	//
	//var friendListener testFriendListener
	//open_im_sdk.SetFriendListener(friendListener)
	//
	//var groupListener testGroupListener
	//open_im_sdk.SetGroupListener(groupListener)
	//var signalingListener testSignalingListener
	//open_im_sdk.SetSignalingListener(&signalingListener)
	//
	//var organizationListener testOrganizationListener
	//open_im_sdk.SetOrganizationListener(organizationListener)
	//
	//var workMomentsListener testWorkMomentsListener
	//open_im_sdk.SetWorkMomentsListener(workMomentsListener)

	//InOutlllogin(uid, tk)

	log.ZWarn(ctx, "SetListenerAndLogin fin", errors.New(""))
}

func lllogin(uid, tk string) bool {
	var callback BaseSuccessFailed
	callback.funcName = utils.GetSelfFuncName()
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
func ReliabilityInitAndLogin(index int, uid, tk, ws, api string) {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api
	cf.WsAddr = ws
	cf.PlatformID = 1
	cf.DataDir = "./"
	cf.IsLogStandardOutput = true
	cf.LogLevel = uint32(LogLevel)

	operationID := utils.OperationIDGenerator()

	ctx := mcontext.NewCtx(operationID)
	var testinit testInitLister
	lg := new(open_im_sdk.LoginMgr)
	log.ZInfo(ctx, "DoReliabilityTest", "UID", uid, "Token", tk, "WS", ws, "API", api)
	log.ZInfo(ctx, "New login manager ", "OperationID", operationID)

	allLoginMgr[index].mgr = lg
	lg.InitSDK(cf, &testinit)

	ctx = ccontext.WithOperationID(lg.Context(), operationID)

	log.ZInfo(ctx, "Initialized SDK with config", "Config", cf)

	var testConversation conversationCallBack
	lg.SetConversationListener(&testConversation)

	var testUser userCallback
	lg.SetUserListener(testUser)

	var msgCallBack MsgListenerCallBak
	lg.SetAdvancedMsgListener(&msgCallBack)

	var friendListener testFriendListener
	lg.SetFriendListener(friendListener)

	var groupListener testGroupListener
	lg.SetGroupListener(groupListener)

	var callback BaseSuccessFailed
	callback.funcName = utils.GetSelfFuncName()

	for {
		if callback.errCode == 1 && testConversation.SyncFlag == 1 {
			lg.User().GetSelfUserInfo(ctx)
			return
		}
	}
}

func PressInitAndLogin(index int, uid, tk, ws, api string) {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api
	cf.WsAddr = ws
	cf.PlatformID = 1
	cf.DataDir = "./"
	cf.LogLevel = uint32(LogLevel)

	operationID := utils.OperationIDGenerator()
	ctx := mcontext.NewCtx(operationID)
	var testinit testInitLister
	lg := new(open_im_sdk.LoginMgr)
	log.ZInfo(ctx, "DoReliabilityTest", "UID", uid, "Token", tk, "WS", ws, "API", api)

	log.ZInfo(ctx, "New login manager ", "OperationID", operationID)

	allLoginMgr[index].mgr = lg
	lg.InitSDK(cf, &testinit)

	log.ZInfo(ctx, "Initialized SDK with config", "Config", cf)

	var testConversation conversationCallBack
	lg.SetConversationListener(&testConversation)

	var testUser userCallback
	lg.SetUserListener(testUser)

	var msgCallBack MsgListenerCallBak
	lg.SetAdvancedMsgListener(&msgCallBack)

	var friendListener testFriendListener
	lg.SetFriendListener(friendListener)

	var groupListener testGroupListener
	lg.SetGroupListener(groupListener)

	err := lg.Login(ctx, uid, tk)
	if err != nil {
		log.ZError(ctx, "Login failed", err, "OperationID", operationID)
	}
}

func DoTest(uid, tk, ws, api string) {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = api // "http://120.24.45.199:10000"
	cf.WsAddr = ws   //"ws://120.24.45.199:17778"
	cf.PlatformID = 1
	cf.DataDir = "./"

	var s string
	b, _ := json.Marshal(cf)
	s = string(b)
	fmt.Println(s)
	var testinit testInitLister
	operationID := utils.OperationIDGenerator()
	if !open_im_sdk.InitSDK(&testinit, operationID, s) {
		log.ZError(ctx, "InitSDK failed", errors.New("InitSDK failed"))
		return
	}

	var testConversation conversationCallBack
	open_im_sdk.SetConversationListener(&testConversation)

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
	msg         string
	OperationID string
	sendID      string
	recvID      string
	msgID       string
	sendTime    int64
	recvTime    int64
	groupID     string
}

func (t *TestSendMsgCallBack) OnError(errCode int32, errMsg string) {
	log.ZError(ctx, "test_openim: send msg failed: ", errors.New("test_openim: send msg failed: "),
		"operationID", t.OperationID, "errCode", errCode, "errMsg", errMsg, "msgID", t.msgID, "msg", t.msg)
	SendMsgMapLock.Lock()
	defer SendMsgMapLock.Unlock()
	SendFailedAllMsg[t.msgID] = t.sendID + t.recvID

}

func (t *TestSendMsgCallBack) OnSuccess(data string) {
	log.ZInfo(ctx, "test_openim: send msg success: |", "operationID", t.OperationID,
		"msgID", t.msgID, "msg", t.msg, "data", data)
	SendMsgMapLock.Lock()
	defer SendMsgMapLock.Unlock()
	//k, _ := SendSuccAllMsg[t.msgID]
	//k.SendSeccCallbackTime = utils.GetCurrentTimestampByMill()
	//k.SendIDRecvID = t.sendID + t.recvID
}

func (t *TestSendMsgCallBack) OnProgress(progress int) {
	//	fmt.Printf("msg_send , onProgress %d\n", progress)
}

type TestSendMsgCallBackPress struct {
	msg         string
	OperationID string
	sendID      string
	recvID      string
	msgID       string
}

func (t *TestSendMsgCallBackPress) OnError(errCode int32, errMsg string) {
	log.ZWarn(ctx, "TestSendMsgCallBackPress: send msg failed: |", errors.New(""), "operationID", t.OperationID, "errCode",
		errCode, "errMsg", errMsg, "msgID", t.msgID, "msg", t.msg)
}

func (t *TestSendMsgCallBackPress) OnSuccess(data string) {
	log.ZInfo(ctx, "TestSendMsgCallBackPress: send msg success: |", "operationID", t.OperationID, "msgID", t.msgID, "msg", t.msg)
}

func (t *TestSendMsgCallBackPress) OnProgress(progress int) {
	//	fmt.Printf("msg_send , onProgress %d\n", progress)
}

type BaseSuccessFailedTest struct {
	successData string
	errCode     int
	errMsg      string
	funcName    string
}

func (b *BaseSuccessFailedTest) OnError(errCode int32, errMsg string) {
	b.errCode = -1
	b.errMsg = errMsg
	fmt.Println("22onError ", b.funcName, errCode, errMsg)
}

func (b *BaseSuccessFailedTest) OnSuccess(data string) {
	b.errCode = 1
	b.successData = data
	fmt.Println("22OnSuccess: ", b.funcName, data)
}

func InOutDoTestSendMsg(sendId, receiverID string) {
	m := "test:" + sendId + ":" + receiverID + ":"
	//s := CreateTextMessage(m)
	var testSendMsg TestSendMsgCallBack
	//	testSendMsg.msg = SendMessage(&testSendMsg, s, receiverID, "", false)
	fmt.Println("func send ", m, testSendMsg.msg)
	fmt.Println("test to recv : ", receiverID)
}

//func DoTestGetAllConversationList() {
//	var test TestGetAllConversationListCallBack
//	open_im_sdk.GetAllConversationList(test)
//}

type userCallback struct {
}

func (c userCallback) OnUserStatusChanged(statusMap string) {
	log.ZInfo(ctx, "User Status Changed", "statusMap", statusMap)
}

func (userCallback) OnSelfInfoUpdated(callbackData string) {
	log.ZInfo(ctx, "Self Info Updated", "callbackData", callbackData)
}

// //////////////////////////////////////////////////////////////////
type testInitLister struct {
}

func (t *testInitLister) OnUserTokenInvalid(errMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "errMsg", errMsg)
}

func (t *testInitLister) OnUserTokenExpired() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}
func (t *testInitLister) OnConnecting() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (t *testInitLister) OnConnectSuccess() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (t *testInitLister) OnConnectFailed(ErrCode int32, ErrMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "errCode", ErrCode, "errMsg", ErrMsg)
}

func (t *testInitLister) OnKickedOffline() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (t *testInitLister) OnSelfInfoUpdated(info string) {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (t *testInitLister) OnSuccess() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (t *testInitLister) OnError(code int32, msg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "code", code, "msg", msg)
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

func (testFriendListener) OnFriendApplicationAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendApplicationDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendApplicationAccepted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendApplicationRejected(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnBlackAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnBlackDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnFriendInfoChanged(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo)
}

func (testFriendListener) OnSuccess() {
	log.ZInfo(ctx, utils.GetSelfFuncName())
}

func (testFriendListener) OnError(code int32, msg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "Code", code, "Message", msg)
}
