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

var memberUserID = "openIM101"

func DoTestCreateGroup() {
	var test testCreateGroup
	test.OperationID = utils.OperationIDGenerator()

	var groupInfo sdk_params_callback.CreateGroupBaseInfoParam
	groupInfo.GroupName = "group name"
	groupInfo.GroupType = 0

	var memberlist []server_api_params.GroupAddMemberInfo
	memberlist = append(memberlist, server_api_params.GroupAddMemberInfo{UserID: memberUserID, RoleLevel: 1})

	g1 := utils.StructToJsonString(groupInfo)
	g2 := utils.StructToJsonString(memberlist)

	log.Info(test.OperationID, utils.GetSelfFuncName(), "input: ", g1, g2)
	open_im_sdk.CreateGroup(test, test.OperationID, g1, g2)
}

//
//type testSetGroupInfo struct {
//	open_im_sdk.groupInfo
//}
//
//func (testSetGroupInfo) OnSuccess(data string) {
//	fmt.Println("testSetGroupInfo,onSuccess")
//}
//
//func (testSetGroupInfo) OnError(errCode int32, errMsg string) {
//	fmt.Println("testSetGroupInfo,onError")
//}
//
//func DoSetGroupInfo() {
//	var test testSetGroupInfo
//	test.groupInfo.GroupId = "a411065eedf8bc1830ce544ff51394fe"
//	test.GroupName = "test group"
//	test.Introduction = "This is an introduction about the test group"
//	test.Notification = "this is test bulletins"
//	test.FaceUrl = "this is a test face url"
//	setInfo, _ := json.Marshal(test.groupInfo)
//	fmt.Println("setGroupInfo input", string(setInfo))
//	sdk_interface.SetGroupInfo(string(setInfo), test)
//}
//
//type testGetGroupsInfo struct {
//	open_im_sdk.getGroupsInfoReq
//}
//
//func (testGetGroupsInfo) OnSuccess(data string) {
//	fmt.Println("testGetGroupsInfo,onSuccess", data)
//}
//
//func (testGetGroupsInfo) OnError(errCode int32, errMsg string) {
//	fmt.Println("testGetGroupsInfo,onError", errMsg)
//}
//
//func DoGetGroupsInfo() {
//	var test testGetGroupsInfo
//	groupIDList := []string{"a411065eedf8bc1830ce544ff51394fe"}
//	test.getGroupsInfoReq.GroupIDList = groupIDList
//	groupsIDList, _ := json.Marshal(test.GroupIDList)
//	fmt.Println("test getGroupsInfo input", string(groupsIDList))
//	sdk_interface.GetGroupsInfo(string(groupsIDList), test)
//}
//
//type testJoinGroup struct {
//	open_im_sdk.joinGroupReq
//}
//
//func (testJoinGroup) OnSuccess(data string) {
//	fmt.Println("testJoinGroup,onSuccess", data)
//}
//
//func (testJoinGroup) OnError(errCode int32, errMsg string) {
//	fmt.Println("testJoinGroup,onError", errMsg)
//}
//
//func DoJoinGroup() {
//	var test testJoinGroup
//	test.joinGroupReq.GroupID = "7149948c2fb143f9ee97e3e9b406b5ec"
//	test.joinGroupReq.Message = "jin lai "
//
//	fmt.Println("test join group input", test.GroupID, test.Message)
//	sdk_interface.JoinGroup(test.GroupID, test.Message, test)
//}
//
//type testQuitGroup struct {
//	open_im_sdk.quitGroupReq
//}
//
//func (testQuitGroup) OnSuccess(data string) {
//	fmt.Println("testQuitGroup,onSuccess", data)
//}
//
//func (testQuitGroup) OnError(errCode int32, errMsg string) {
//	fmt.Println("testQuitGroup,onError", errMsg)
//}
//
//func DoQuitGroup() {
//	var test testQuitGroup
//	test.quitGroupReq.GroupID = "77215e1b13b75f3ab00cb6570e3d9618"
//
//	fmt.Println("test quit group input", test.GroupID)
//	sdk_interface.QuitGroup(test.GroupID, test)
//}
//
//type testGetJoinedGroupList struct {
//}
//
///*
//	OnError(errCode int, errMsg string)
//	OnSuccess(data string)
//*/
//func (testGetJoinedGroupList) OnError(errCode int32, errMsg string) {
//	fmt.Println("testGetJoinedGroupList OnError", errCode, errMsg)
//}
//
//func (testGetJoinedGroupList) OnSuccess(data string) {
//	fmt.Println("testGetJoinedGroupList OnSuccess, output", data)
//}
//
//func DotestGetJoinedGroupList() {
//	var test testGetJoinedGroupList
//	sdk_interface.GetJoinedGroupList(test)
//}
//
//type testGetGroupMemberList struct {
//}
//
//func (testGetGroupMemberList) OnError(errCode int32, errMsg string) {
//	fmt.Println("testGetGroupMemberList OnError", errCode, errMsg)
//}
//
//func (testGetGroupMemberList) OnSuccess(data string) {
//	fmt.Println("testGetGroupMemberList OnSuccess, output", data)
//}
//
//func DotestGetGroupMemberList() {
//	var test testGetGroupMemberList
//	var groupId string = ""
//	sdk_interface.GetGroupMemberList(groupId, 0, 0, test)
//}
//
//type testGetGroupMembersInfo struct {
//}
//
//func (testGetGroupMembersInfo) OnError(errCode int32, errMsg string) {
//	fmt.Println("testGetGroupMembersInfo OnError", errCode, errMsg)
//}
//
//func (testGetGroupMembersInfo) OnSuccess(data string) {
//	fmt.Println("testGetGroupMembersInfo OnSuccess, output", data)
//}
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
//type testKickGroupMember struct {
//}
//
//func (testKickGroupMember) OnError(errCode int32, errMsg string) {
//	fmt.Println("testKickGroupMember OnError", errCode, errMsg)
//}
//
//func (testKickGroupMember) OnSuccess(data string) {
//	fmt.Println("testKickGroupMember OnSuccess, output", data)
//}
//
//func DotestKickGroupMember() {
//	var test testKickGroupMember
//	var memlist []string
//	//memlist = append(memlist, "e7b437c8b05a1fb8875e7196c636f327")
//	memlist = append(memlist, "307edc814bb0d04a")
//	jlist, _ := json.Marshal(memlist)
//
//	fmt.Println("KickGroupMember input", string(jlist))
//	sdk_interface.KickGroupMember("f4cc5c9b556221b92992538f7e6ac26e", "kkkkkkk", string(jlist), test)
//}
//
//type testInviteUserToGroup struct {
//}
//
//func (testInviteUserToGroup) OnError(errCode int32, errMsg string) {
//	fmt.Println("testInviteUserToGroup OnError", errCode, errMsg)
//}
//
//func (testInviteUserToGroup) OnSuccess(data string) {
//	fmt.Println("testInviteUserToGroup OnSuccess, output", data)
//}
//
//func DotesttestInviteUserToGroup() {
//	var test testInviteUserToGroup
//	var memlist []string
//	memlist = append(memlist, "307edc814bb0d04a")
//	//memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
//	jlist, _ := json.Marshal(memlist)
//	fmt.Println("DotesttestInviteUserToGroup, input: ", string(jlist))
//	sdk_interface.InviteUserToGroup("f4cc5c9b556221b92992538f7e6ac26e", "friend", string(jlist), test)
//}
//
//type testGroupX struct {
//}
//
//func (testGroupX) OnSuccess(data string) {
//	fmt.Println("testGroupX,onSuccess", data)
//}
//
//func (testGroupX) OnError(errCode int32, errMsg string) {
//	fmt.Println("testGroupX,onError", errMsg)
//}
//func (testGroupX) OnProgress(progress int) {
//	fmt.Println("testGroupX  ", progress)
//}
//
//func DoGetGroupApplicationList() string {
//	//	var test testGroupX
//	fmt.Println("test DoGetGroupApplicationList....")
//
//	return ""
//}
//func DoGroupApplicationList() {
//	var test testGroupX
//	fmt.Println("test DoGetGroupApplicationList....")
//	sdk_interface.GetGroupApplicationList(test)
//}
//func DoTransferGroupOwner(groupid, userid string) {
//	var test testGroupX
//	fmt.Println("test DoTransferGroupOwner....")
//	sdk_interface.TransferGroupOwner(groupid, userid, test)
//}
//func DoAcceptGroupApplication(uid string) {
//
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
//	var test testGroupX
//	fmt.Println("accept", string(v))
//	sdk_interface.AcceptGroupApplication(string(v), "accept", test)
//}
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
