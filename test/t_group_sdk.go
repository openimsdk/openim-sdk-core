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
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/sdk_params_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"

	//	"encoding/json"
	"fmt"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
	//"open_im_sdk/internal/open_im_sdk"
	//"open_im_sdk/pkg/utils"
	//	"open_im_sdk/internal/common"
)

type XBase struct {
}

func (XBase) OnError(errCode int32, errMsg string) {
	fmt.Println("get groupmenberinfo OnError", errCode, errMsg)
}
func (XBase) OnSuccess(data string) {
	fmt.Println("get groupmenberinfo OnSuccess, ", data)
}

func (XBase) OnProgress(progress int) {
	fmt.Println("OnProgress, ", progress)
}

type testGroupListener struct {
}

func (testGroupListener) OnJoinedGroupAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnJoinedGroupDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupMemberAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupMemberDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupApplicationAdded(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupApplicationDeleted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupInfoChanged(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupMemberInfoChanged(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupApplicationAccepted(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupApplicationRejected(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

func (testGroupListener) OnGroupDismissed(callbackInfo string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "CallbackInfo", callbackInfo, "operationID", utils.OperationIDGenerator())
}

type testOrganizationListener struct {
}

func (testOrganizationListener) OnOrganizationUpdated() {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "on listener callback", "operationID", utils.OperationIDGenerator())
}

type testWorkMomentsListener struct {
}

func (testWorkMomentsListener) OnRecvNewNotification() {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "on listener callback", "operationID", utils.OperationIDGenerator())
}

type testCreateGroup struct {
	OperationID string
}

func (t testCreateGroup) OnSuccess(data string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "Data", data)
}

func (t testCreateGroup) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "ErrorCode", errCode, "ErrorMsg", errMsg)
}

func SetTestGroupID(groupID, memberID string) {
	MemberUserID = memberID
	TestgroupID = groupID
}

var MemberUserID = "2101502031"
var me = "3984071717"
var TestgroupID = "3109164461"

func DoTestCreateGroup() {
	var test testCreateGroup
	test.OperationID = utils.OperationIDGenerator()

	var groupInfo sdk_params_callback.CreateGroupBaseInfoParam
	groupInfo.GroupName = "聊聊大群测试"
	groupInfo.GroupType = 1

	var memberlist []server_api_params.GroupAddMemberInfo
	memberlist = append(memberlist, server_api_params.GroupAddMemberInfo{UserID: MemberUserID, RoleLevel: 1})
	memberlist = append(memberlist, server_api_params.GroupAddMemberInfo{UserID: me, RoleLevel: 2})

	g1 := utils.StructToJsonString(groupInfo)
	g2 := utils.StructToJsonString(memberlist)

	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID: ", test.OperationID, "g1", g1, "g2", g2)
	// open_im_sdk.CreateGroup(test, test.OperationID, g1, g2)
}

type testSetGroupInfo struct {
	OperationID string
}

func (t testSetGroupInfo) OnSuccess(data string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "data", data)
}

func (t testSetGroupInfo) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "errCode", errCode, "errMsg", errMsg)
}

func DoSetGroupInfo() {
	var test testSetGroupInfo
	operationID := utils.OperationIDGenerator()
	test.OperationID = operationID
	var input sdk_params_callback.SetGroupInfoParam
	input.GroupName = "new group name 11111111"
	input.Notification = "new notification 11111"
	var n int32
	n = 1
	input.NeedVerification = &n
	setInfo := utils.StructToJsonString(input)
	// open_im_sdk.SetGroupInfo(test, operationID, TestgroupID, setInfo)
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", operationID, "input: ", setInfo)
}

func DoSetGroupVerification() {
	var test testSetGroupInfo
	operationID := utils.OperationIDGenerator()
	test.OperationID = operationID
	open_im_sdk.SetGroupVerification(test, operationID, TestgroupID, 1)
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", operationID, "input: ", TestgroupID, 2)
}

func DoSetGroupLookMemberInfo() {
	var test testSetGroupInfo
	operationID := utils.OperationIDGenerator()
	test.OperationID = operationID
	open_im_sdk.SetGroupLookMemberInfo(test, operationID, TestgroupID, 0)
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", operationID, "input: ", TestgroupID, 1)
}

func DoSetGroupApplyMemberFriend() {
	var test testSetGroupInfo
	operationID := utils.OperationIDGenerator()
	test.OperationID = operationID
	open_im_sdk.SetGroupApplyMemberFriend(test, operationID, TestgroupID, 1)
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", operationID, "input: ", TestgroupID, 1)
}

type testGetGroupsInfo struct {
	OperationID string
}

func (t testGetGroupsInfo) OnSuccess(data string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "testGetGroupsInfo,onSuccess", data)
}

func (t testGetGroupsInfo) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "testGetGroupsInfo,onError", errCode, errMsg)
}

type testSearchGroups struct {
	OperationID string
}

func (t testSearchGroups) OnSuccess(data string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "data", data)
}

func (t testSearchGroups) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, utils.GetSelfFuncName(), "operationID", t.OperationID, "errCode", errCode, "errMsg", errMsg)
}

func DoTestGetGroupsInfo() {
	var test testGetGroupsInfo
	groupIDList := []string{TestgroupID}
	list := utils.StructToJsonString(groupIDList)
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DoTestGetGroupsInfo", "operationID", test.OperationID, "input", list)
	// open_im_sdk.GetGroupsInfo(test, test.OperationID, list)
}

func DoTestSearchGroups() {
	var test testGetGroupsInfo
	var params sdk_params_callback.SearchGroupsParam
	params.KeywordList = []string{"17"}
	//params.IsSearchGroupID =true
	params.IsSearchGroupName = true
	open_im_sdk.SearchGroups(test, test.OperationID, utils.StructToJsonString(params))
}

type testJoinGroup struct {
	OperationID string
}

func (t testJoinGroup) OnSuccess(data string) {
	log.ZInfo(ctx, "testJoinGroup", "operationID", t.OperationID, "onSuccess", data)
}

func (t testJoinGroup) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, "testJoinGroup", "operationID", t.OperationID, "onError", errCode, errMsg)
}

func DoTestJoinGroup() {
	var test testJoinGroup
	test.OperationID = utils.OperationIDGenerator()
	groupID := "1003105543"
	reqMsg := "121212"
	ex := "ex"
	log.ZInfo(ctx, "testJoinGroup", "operationID", test.OperationID, "input", groupID, reqMsg, ex)
	open_im_sdk.JoinGroup(test, test.OperationID, groupID, reqMsg, constant.JoinBySearch, ex)
}

type testQuitGroup struct {
	OperationID string
}

func (t testQuitGroup) OnSuccess(data string) {
	log.ZInfo(ctx, "testQuitGroup", "operationID", t.OperationID, "onSuccess", data)
}

func (t testQuitGroup) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, "testQuitGroup", "operationID", t.OperationID, "onError", errCode, errMsg)
}
func DoTestQuitGroup() {
	var test testQuitGroup
	test.OperationID = utils.OperationIDGenerator()
	groupID := "19de93b442a1ca3b772aa0f12761939d"
	log.ZInfo(ctx, "testQuitGroup", "operationID", test.OperationID, "input", groupID)
	open_im_sdk.QuitGroup(test, test.OperationID, groupID)
}

type testGetJoinedGroupList struct {
	OperationID string
}

/*
OnError(errCode int, errMsg string)
OnSuccess(data string)
*/
func (t testGetJoinedGroupList) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, "testGetJoinedGroupList", "operationID", t.OperationID, "OnError", errCode, errMsg)
}

func (t testGetJoinedGroupList) OnSuccess(data string) {
	log.ZInfo(ctx, "testGetJoinedGroupList", "operationID", t.OperationID, "OnSuccess", "output", data)
}

func DoTestGetJoinedGroupList() {
	var test testGetJoinedGroupList
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.GetJoinedGroupList(test, test.OperationID)
}

type testGetGroupMemberList struct {
	OperationID string
}

func (t testGetGroupMemberList) OnSuccess(data string) {
	log.ZInfo(ctx, "testGetGroupMemberList", "operationID", t.OperationID, "function", utils.GetSelfFuncName(), "data", data)
}

func (t testGetGroupMemberList) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, "testGetGroupMemberList", "operationID", t.OperationID, "function", utils.GetSelfFuncName(), "errCode", errCode, "errMsg", errMsg)
}

func DotestGetGroupMemberList() {
	var test testGetGroupMemberList
	test.OperationID = utils.OperationIDGenerator()
	var groupId = TestgroupID
	open_im_sdk.GetGroupMemberList(test, test.OperationID, groupId, 4, 0, 100)
}

func DotestCos() {
	//var callback baseCallback
	//p := ws.NewPostApi(token, UserForSDK.ImConfig().ApiAddr)
	//var storage common.ObjectStorage = common.NewCOS(p)
	//test(storage, callback)
}

//funcation DotestMinio() {
//	var callback baseCallback
//	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVSUQiOiIxMzkwMDAwMDAwMCIsIlBsYXRmb3JtIjoiSU9TIiwiZXhwIjoxNjQ1NzgyNDY0LCJuYmYiOjE2NDUxNzc2NjQsImlhdCI6MTY0NTE3NzY2NH0.T-SDoLxdlwRGOMZPIKriPtAlOGWCLodsGi1dWxN8kto"
//	p := ws.NewPostApi(token, "https://storage.rentsoft.cn")
//	minio := common.NewMinio(p)
//	var storage common.ObjectStorage = minio
//	log.NewInfo("", *minio)
//	test(storage, callback)
//}
//
//funcation test(storage common.ObjectStorage, callback baseCallback) {
//	dir, newName, err := storage.UploadFile("./cmd/main.go", funcation(progress int) {
//		if progress == 100 {
//			callback.OnSuccess("")
//		}
//	})
//	log.NewInfo("0", dir, newName, err)
//	dir, newName, err = storage.UploadImage("C:\\Users\\Administrator\\Desktop\\1.jpg", funcation(progress int) {
//		if progress == 100 {
//			callback.OnSuccess("")
//		}
//	})
//	log.NewInfo("0", dir, newName, err, err)
//	dir, newName, err = storage.UploadSound("./cmd/main.go", funcation(progress int) {
//		if progress == 100 {
//			callback.OnSuccess("")
//		}
//	})
//	log.NewInfo("0", dir, newName, err, err)
//	snapshotURL, snapshotUUID, videoURL, videoUUID, err := storage.UploadVideo("./cmd/main.go", "C:\\Users\\Administrator\\Desktop\\1.jpg", funcation(progress int) {
//		if progress == 100 {
//			callback.OnSuccess("")
//		}
//	})
//	log.NewInfo(snapshotURL, snapshotUUID, videoURL, videoUUID, err)
//}

type testGetGroupMembersInfo struct {
}

func (testGetGroupMembersInfo) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetGroupMembersInfo OnError", errCode, errMsg)
}

func (testGetGroupMembersInfo) OnSuccess(data string) {
	fmt.Println("testGetGroupMembersInfo OnSuccess, output", data)
}

//
//funcation DotestGetGroupMembersInfo() {
//	var test testGetGroupMembersInfo
//	var memlist []string
//	memlist = append(memlist, "307edc814bb0d04a")
//	//memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
//	jlist, _ := json.Marshal(memlist)
//	fmt.Println("GetGroupMembersInfo input : ", string(jlist))
//	sdk_interface.GetGroupMembersInfo("7ff61d8f9d4a8a0d6d70a14e2683aad5", string(jlist), test)
//	//GetGroupMemberList("05dc84b52829e82242a710ecf999c72c", 0, 0, test)
//}
//

type baseCallback struct {
	OperationID string
	callName    string
}

func (t baseCallback) OnSuccess(data string) {
	log.ZInfo(ctx, t.callName, "operationID", t.OperationID, "function", utils.GetSelfFuncName(), "data", data)
}

func (t baseCallback) OnError(errCode int32, errMsg string) {
	log.ZInfo(ctx, t.callName, "operationID", t.OperationID, "function", utils.GetSelfFuncName(), "errCode", errCode, "errMsg", errMsg)
}

type testKickGroupMember struct {
	baseCallback
}
type testGetGroupMemberListByJoinTimeFilter struct {
	baseCallback
}

func DotestGetGroupMemberListByJoinTimeFilter() {
	var test testGetGroupMemberListByJoinTimeFilter
	test.OperationID = utils.OperationIDGenerator()
	var memlist []string
	jlist := utils.StructToJsonString(memlist)
	log.ZInfo(ctx, "DotestGetGroupMemberListByJoinTimeFilter", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", jlist)
	open_im_sdk.GetGroupMemberListByJoinTimeFilter(test, test.OperationID, "3750066757", 1, 40, 0, 0, jlist)
}

func DotestKickGroupMember() {
	var test testKickGroupMember
	test.OperationID = utils.OperationIDGenerator()
	var memlist []string
	memlist = append(memlist, MemberUserID)
	jlist := utils.StructToJsonString(memlist)
	log.ZInfo(ctx, "DotestKickGroupMember", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", jlist)
	open_im_sdk.KickGroupMember(test, test.OperationID, TestgroupID, "kkk", jlist)
}

type testInviteUserToGroup struct {
	baseCallback
}

func DotestInviteUserToGroup() {
	var test testInviteUserToGroup
	test.OperationID = utils.OperationIDGenerator()
	var memlist []string
	memlist = append(memlist, MemberUserID)
	jlist := utils.StructToJsonString(memlist)
	log.ZInfo(ctx, "DotestInviteUserToGroup", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", jlist)
	open_im_sdk.InviteUserToGroup(test, test.OperationID, TestgroupID, "come", string(jlist))
}

type testGetGroupApplicationList struct {
	baseCallback
}

func DotestGetRecvGroupApplicationList() string {
	var test testGetGroupApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DotestGetRecvGroupApplicationList", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	// open_im_sdk.GetRecvGroupApplicationList(test, test.OperationID)
	return ""
}

//	funcation DoGroupApplicationList() {
//		var test testGroupX
//		fmt.Println("test DoGetGroupApplicationList....")
//		sdk_interface.GetGroupApplicationList(test)
//	}
type testTransferGroupOwner struct {
	baseCallback
}

func DotestTransferGroupOwner() {
	var test testTransferGroupOwner
	test.OperationID = utils.OperationIDGenerator()

	open_im_sdk.TransferGroupOwner(test, test.OperationID, TestgroupID, MemberUserID)

}

type testProcessGroupApplication struct {
	baseCallback
}

func DoTestAcceptGroupApplication(uid string) {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DoTestAcceptGroupApplication", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	open_im_sdk.AcceptGroupApplication(test, test.OperationID, TestgroupID, MemberUserID, "ok")
}

func DoTestGetUserReqGroupApplicationList() {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DoTestGetUserReqGroupApplicationList", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	// open_im_sdk.GetSendGroupApplicationList(test, test.OperationID)
}

func DoTestGetRecvGroupApplicationList() {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DoTestGetRecvGroupApplicationList", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	// open_im_sdk.GetRecvGroupApplicationList(test, test.OperationID)
}

func DotestRefuseGroupApplication(uid string) {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DotestRefuseGroupApplication", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	open_im_sdk.RefuseGroupApplication(test, test.OperationID, TestgroupID, MemberUserID, "no")
}

type testSetGroupMemberNickname struct {
	baseCallback
}

func DotestSetGroupMemberNickname(myUserID string) {
	var test testSetGroupMemberNickname
	test.OperationID = utils.OperationIDGenerator()
	log.ZInfo(ctx, "DotestSetGroupMemberNickname", "operationID", test.OperationID, "function", utils.GetSelfFuncName(), "input", "")
	open_im_sdk.SetGroupMemberNickname(test, test.OperationID, TestgroupID, myUserID, "")
}

func DoTestSetGroupMemberRoleLevel(groupID, userID string, roleLevel int) {
	var test testSetGroupMemberNickname
	test.OperationID = utils.OperationIDGenerator()
	fmt.Println(test.OperationID, utils.GetSelfFuncName(), "inputx: ")
	open_im_sdk.SetGroupMemberRoleLevel(test, test.OperationID, groupID, userID, roleLevel)
}

func DoTestSetGroupMemberInfo(groupID, userID string, ex string) {
	var test testSetGroupMemberNickname
	test.OperationID = utils.OperationIDGenerator()
	param := sdk_params_callback.SetGroupMemberInfoParam{GroupID: groupID, UserID: userID}
	if ex != "" {
		param.Ex = &ex
	}
	g1 := utils.StructToJsonString(param)
	fmt.Println(test.OperationID, utils.GetSelfFuncName(), "inputx: ", g1)

	open_im_sdk.SetGroupMemberInfo(test, test.OperationID, g1)
}
