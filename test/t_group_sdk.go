package test

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/server_api_params"

	//	"encoding/json"
	"fmt"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	//"open_im_sdk/internal/open_im_sdk"
	//"open_im_sdk/pkg/utils"
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

func (testGroupListener) OnReceiveJoinGroupApplicationAdded(callbackInfo string) {
	log.Info(utils.OperationIDGenerator(), utils.GetSelfFuncName(), callbackInfo)

}
func (testGroupListener) OnReceiveJoinGroupApplicationDeleted(callbackInfo string) {
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

//
type testCreateGroup struct {
	OperationID string
}

func (t testCreateGroup) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t testCreateGroup) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
}

var MemberUserID = "openIM101"
var TestgroupID = "19de93b442a1ca3b772aa0f12761939d"

func DoTestCreateGroup() {
	var test testCreateGroup
	test.OperationID = utils.OperationIDGenerator()

	var groupInfo sdk_params_callback.CreateGroupBaseInfoParam
	groupInfo.GroupName = "group name"
	groupInfo.GroupType = 0

	var memberlist []server_api_params.GroupAddMemberInfo
	memberlist = append(memberlist, server_api_params.GroupAddMemberInfo{UserID: MemberUserID, RoleLevel: 1})

	g1 := utils.StructToJsonString(groupInfo)
	g2 := utils.StructToJsonString(memberlist)

	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", g1, g2)
	open_im_sdk.CreateGroup(test, test.OperationID, g1, g2)
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
	input.GroupName = "new group name 111"
	input.Notification = "new notification 222"

	setInfo := utils.StructToJsonString(input)
	open_im_sdk.SetGroupInfo(test, test.OperationID, TestgroupID, setInfo)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", setInfo)

}

//
type testGetGroupsInfo struct {
	OperationID string
}

func (t testGetGroupsInfo) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetGroupsInfo,onSuccess", data)
}

func (t testGetGroupsInfo) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testGetGroupsInfo,onError", errCode, errMsg)
}

func DoTestGetGroupsInfo() {
	var test testGetGroupsInfo
	groupIDList := []string{"8a33030b726bd4792c8410aadfacaa35", "e91805bae94ae3a00eb629f74e83605a"}
	list := utils.StructToJsonString(groupIDList)
	log.Info(test.OperationID, "test getGroupsInfo input", list)
	open_im_sdk.GetGroupsInfo(test, test.OperationID, list)
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
	groupID := "19de93b442a1ca3b772aa0f12761939d"
	reqMsg := "123123"
	log.Info(test.OperationID, "test join group input", groupID, reqMsg)
	open_im_sdk.JoinGroup(test, test.OperationID, groupID, reqMsg)
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

//
func DoTestGetJoinedGroupList() {
	var test testGetJoinedGroupList
	test.OperationID = utils.OperationIDGenerator()
	open_im_sdk.GetJoinedGroupList(test, test.OperationID)
}

type testGetGroupMemberList struct {
	OperationID string
}

func (t testGetGroupMemberList) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t testGetGroupMemberList) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
}

func DotestGetGroupMemberList() {
	var test testGetGroupMemberList
	test.OperationID = utils.OperationIDGenerator()
	var groupId = TestgroupID
	open_im_sdk.GetGroupMemberList(test, test.OperationID, groupId, 1, 30)
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input ", groupId, 1, 30)
}

type testGetGroupMembersInfo struct {
}

func (testGetGroupMembersInfo) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetGroupMembersInfo OnError", errCode, errMsg)
}

func (testGetGroupMembersInfo) OnSuccess(data string) {
	fmt.Println("testGetGroupMembersInfo OnSuccess, output", data)
}

//
//func DotestGetGroupMembersInfo() {
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
}

func (t baseCallback) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t baseCallback) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
}

type testKickGroupMember struct {
	baseCallback
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

func DotestGetGroupApplicationList() string {
	var test testGetGroupApplicationList
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.GetGroupApplicationList(test, test.OperationID)
	return ""
}

//func DoGroupApplicationList() {
//	var test testGroupX
//	fmt.Println("test DoGetGroupApplicationList....")
//	sdk_interface.GetGroupApplicationList(test)
//}
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

func DotestAcceptGroupApplication(uid string) {
	var test testProcessGroupApplication
	test.OperationID = utils.OperationIDGenerator()
	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ")
	open_im_sdk.AcceptGroupApplication(test, TestgroupID, MemberUserID, "ok", test.OperationID)
}

//func DoRefuseGroupApplication(uid string) {
//	str := DoGetGroupApplicationList()
//	var ret open_im_sdk.groupApplicationResult
//	err := json.Unmarshal([]byte(str), &ret)
//	if err != nil {
//		return
//	}
//	var app utils.GroupReqListInfo
//	for i := 0; i < len(ret.GroupApplicationList); i++ {
//		if ret.GroupApplicationList[i].FromUserID == uid {
//			app = ret.GroupApplicationList[i]
//			break
//		}
//	}
//
//	v, err := json.Marshal(app)
//	if err != nil {
//		return
//	}
//
//	fmt.Println(string(v))
//
//	var test testGroupX
//	fmt.Println("refuse", string(v))
//	sdk_interface.RefuseGroupApplication(string(v), "refuse", test)
//}
