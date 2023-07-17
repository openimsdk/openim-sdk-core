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
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"

	//	"encoding/json"
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
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
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnJoinedGroupDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupMemberAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupMemberDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupApplicationAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupApplicationDeleted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupMemberInfoChanged(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupApplicationAccepted(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnGroupApplicationRejected(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

func (testGroupListener) OnGroupDismissed(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}

type testOrganizationListener struct {
}

func (testOrganizationListener) OnOrganizationUpdated() {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), "on listener callback ")
}

type testWorkMomentsListener struct {
}

func (testWorkMomentsListener) OnRecvNewNotification() {
	log.NewInfo(utils.OperationIDGenerator(), utils.GetSelfFuncName(), "on listener callback ")
}

type testCreateGroup struct {
	OperationID string
}

func (t testCreateGroup) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t testCreateGroup) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
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

	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", g1, g2)
	// open_im_sdk.CreateGroup(test, test.OperationID, g1, g2)
}

type testSetGroupInfo struct {
	OperationID string
}

func (t testSetGroupInfo) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t testSetGroupInfo) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
}

func DoSetGroupInfo() {
	var test testSetGroupInfo
	test.OperationID = utils.OperationIDGenerator()
	var input sdk_params_callback.SetGroupInfoParam
	input.GroupName = "new group name 11111111"
	input.Notification = "new notification 11111"
	var n int32
	n = 1
	input.NeedVerification = &n
	setInfo := utils.StructToJsonString(input)
	// open_im_sdk.SetGroupInfo(test, test.OperationID, TestgroupID, setInfo)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", setInfo)

}

func DoSetGroupVerification() {
	var test testSetGroupInfo
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.SetGroupVerification(test, test.OperationID, TestgroupID, 1)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", TestgroupID, 2)
}
func DoSetGroupLookMemberInfo() {
	var test testSetGroupInfo
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.SetGroupLookMemberInfo(test, test.OperationID, TestgroupID, 0)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", TestgroupID, 1)
}
func DoSetGroupApplyMemberFriend() {
	var test testSetGroupInfo
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.SetGroupApplyMemberFriend(test, test.OperationID, TestgroupID, 1)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", TestgroupID, 1)
}

type testGetGroupsInfo struct {
	OperationID string
}

func (t testGetGroupsInfo) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetGroupsInfo,onSuccess", data)
}

func (t testGetGroupsInfo) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testGetGroupsInfo,onError", errCode, errMsg)
}

type testSearchGroups struct {
	OperationID string
}

func (t testSearchGroups) OnSuccess(data string) {
	log.Info(t.OperationID, "testSearchGroups,onSuccess", data)
}

func (t testSearchGroups) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testSearchGroups,onError", errCode, errMsg)
}
func DoTestGetGroupsInfo() {
	var test testGetGroupsInfo
	groupIDList := []string{TestgroupID}
	list := utils.StructToJsonString(groupIDList)
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, "test getGroupsInfo input", list)
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
	log.Info(t.OperationID, "testJoinGroup,onSuccess", data)
}

func (t testJoinGroup) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testJoinGroup,onError", errCode, errMsg)
}

func DoTestJoinGroup() {
	var test testJoinGroup
	test.OperationID = utils.OperationIDGenerator()
	groupID := "1003105543"
	reqMsg := "121212"
	log.Info(test.OperationID, "test join group input", groupID, reqMsg)
	open_im_sdk.JoinGroup(test, test.OperationID, groupID, reqMsg, constant.JoinBySearch)
}

type testQuitGroup struct {
	OperationID string
}

func (t testQuitGroup) OnSuccess(data string) {
	log.Info(t.OperationID, "testQuitGroup,onSuccess", data)
}

func (t testQuitGroup) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testQuitGroup,onError", errCode, errMsg)
}

func DoTestQuitGroup() {
	var test testQuitGroup
	test.OperationID = utils.OperationIDGenerator()
	groupID := "19de93b442a1ca3b772aa0f12761939d"
	log.Info(test.OperationID, "test quit group input", groupID)
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
	log.Info(t.OperationID, "testGetJoinedGroupList OnError", errCode, errMsg)
}

func (t testGetJoinedGroupList) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetJoinedGroupList OnSuccess, output", data)
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
	log.Info(t.OperationID, utils.GetSelfFuncName(), "testGetGroupMemberList: ", data)

}

func (t testGetGroupMemberList) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), "testGetGroupMemberList", errCode, errMsg)
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
	log.Info(t.OperationID, t.callName, utils.GetSelfFuncName(), data)

}

func (t baseCallback) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, t.callName, utils.GetSelfFuncName(), errCode, errMsg)
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
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ", jlist)
	open_im_sdk.GetGroupMemberListByJoinTimeFilter(test, test.OperationID, "3750066757", 1, 40, 0, 0, jlist)
}
func DotestKickGroupMember() {
	var test testKickGroupMember
	test.OperationID = utils.OperationIDGenerator()
	var memlist []string
	memlist = append(memlist, MemberUserID)
	jlist := utils.StructToJsonString(memlist)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ", jlist)
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
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ", jlist)
	open_im_sdk.InviteUserToGroup(test, test.OperationID, TestgroupID, "come", string(jlist))
}

type testGetGroupApplicationList struct {
	baseCallback
}

func DotestGetRecvGroupApplicationList() string {
	var test testGetGroupApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
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
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.AcceptGroupApplication(test, test.OperationID, TestgroupID, MemberUserID, "ok")
}

func DoTestGetUserReqGroupApplicationList() {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	// open_im_sdk.GetSendGroupApplicationList(test, test.OperationID)
}

// 提示

func DoTestGetRecvGroupApplicationList() {
	var test testProcessGroupApplication
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input:")
	// open_im_sdk.GetRecvGroupApplicationList(test, test.OperationID)
}

func DotestRefuseGroupApplication(uid string) {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.RefuseGroupApplication(test, test.OperationID, TestgroupID, MemberUserID, "no")
}

type testSetGroupMemberNickname struct {
	baseCallback
}

func DotestSetGroupMemberNickname(myUserID string) {
	var test testSetGroupMemberNickname
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
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
