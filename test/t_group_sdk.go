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

type testSetGroupInfo struct {
	OperationID string
}

func (t testSetGroupInfo) OnSuccess(data string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), data)

}

func (t testSetGroupInfo) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, utils.GetSelfFuncName(), errCode, errMsg)
}

var TestgroupID = "19de93b442a1ca3b772aa0f12761939d"

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
type testGetGroupMembersInfo struct {
	OperationID string
}

func (t testGetGroupMembersInfo) OnError(errCode int32, errMsg string) {
	log.Info(t.OperationID, "testGetGroupMembersInfo OnError", errCode, errMsg)
}

func (t testGetGroupMembersInfo) OnSuccess(data string) {
	log.Info(t.OperationID, "testGetGroupMembersInfo OnSuccess, output", data)
}

func DoTestGetGroupMembersInfo() {
	var test testGetGroupMembersInfo
	test.OperationID = utils.OperationIDGenerator()
	var memlist []string
	groupID := "e5868dbb42ec8bd559098c92cf72cdf8"
	memlist = append(memlist, "openIM100")
	jlist := utils.StructToJsonStringDefault(memlist)
	log.Info(test.OperationID, "GetGroupMembersInfo input : ", jlist)
	open_im_sdk.GetGroupMembersInfo(test, test.OperationID, groupID, jlist)

}

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
