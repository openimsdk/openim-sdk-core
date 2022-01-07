package open_im_sdk

import (
	"encoding/json"
	"fmt"
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

func (testGroupListener) OnMemberEnter(groupId string, memberList string) {
	fmt.Println("OnMemberEnter", groupId, memberList)
}

func (testGroupListener) OnTransferGroupOwner(groupId string, oldUserID string, newUserID string) {
	fmt.Println("OnTransferGroupOwner", groupId, oldUserID, newUserID)
}

func (testGroupListener) OnMemberLeave(groupId string, member string) {
	fmt.Println("OnMemberLeave", groupId, member)
}
func (testGroupListener) OnMemberInvited(groupId string, opUser string, memberList string) {
	fmt.Println("OnMemberInvited", groupId, opUser, memberList)
}
func (testGroupListener) OnMemberKicked(groupId string, opUser string, memberList string) {
	fmt.Println("OnMemberKicked", groupId, opUser, memberList)
}
func (testGroupListener) OnGroupCreated(groupId string) {
	fmt.Println("OnGroupCreated", groupId)
}

func (testGroupListener) OnGroupInfoChanged(groupId string, groupInfo string) {
	fmt.Println("OnGroupInfoChanged", groupId, groupInfo)
}

func (testGroupListener) OnReceiveJoinApplication(groupId string, member string, opReason string) {
	fmt.Println("OnReceiveJoinApplication", groupId, member, opReason)
}

func (testGroupListener) OnApplicationProcessed(groupId string, opUser string, AgreeOrReject int32, opReason string) {
	fmt.Println("OnApplicationProcessed", groupId, opUser, AgreeOrReject, opReason)
}

type testCreateGroup struct {
	groupInfo
	members []createGroupMemberInfo
}

func (testCreateGroup) OnSuccess(data string) {
	fmt.Println("testCreateGroup,onSuccess", data)
	var uidList []string
	uidList = append(uidList, "307edc814bb0d04a")
	var xb XBase
	j, _ := json.Marshal(uidList)
	GetGroupMembersInfo(data, string(j), xb)
}

func (testCreateGroup) OnError(errCode int32, errMsg string) {
	fmt.Println("testCreateGroup,onError", errCode, errMsg)
}

func DoTestCreateGroup() {
	var test testCreateGroup
	test.groupInfo.GroupName = "chat group1"
	test.groupInfo.FaceUrl = "url address"
	test.groupInfo.Introduction = "xxxxx"
	test.groupInfo.Notification = "zzzzzzzz"
	test.members = append(test.members, createGroupMemberInfo{
		Uid:     "8a97893c744db83e",
		SetRole: 0,
	})

	groupInfo, _ := json.Marshal(test.groupInfo)
	members, _ := json.Marshal(test.members)
	fmt.Println("create groupInfo input:", string(groupInfo), string(members))
	CreateGroup(test, string(groupInfo), string(members), "asfdfasd")
}

type testSetGroupInfo struct {
	groupInfo
}

func (testSetGroupInfo) OnSuccess(data string) {
	fmt.Println("testSetGroupInfo,onSuccess")
}

func (testSetGroupInfo) OnError(errCode int32, errMsg string) {
	fmt.Println("testSetGroupInfo,onError")
}

func DoSetGroupInfo() {
	var test testSetGroupInfo
	test.groupInfo.GroupId = "a411065eedf8bc1830ce544ff51394fe"
	test.GroupName = "test group"
	test.Introduction = "This is an introduction about the test group"
	test.Notification = "this is test bulletins"
	test.FaceUrl = "this is a test face url"
	setInfo, _ := json.Marshal(test.groupInfo)
	fmt.Println("setGroupInfo input", string(setInfo))
	SetGroupInfo(string(setInfo), test)
}

type testGetGroupsInfo struct {
	getGroupsInfoReq
}

func (testGetGroupsInfo) OnSuccess(data string) {
	fmt.Println("testGetGroupsInfo,onSuccess", data)
}

func (testGetGroupsInfo) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetGroupsInfo,onError", errMsg)
}

func DoGetGroupsInfo() {
	var test testGetGroupsInfo
	groupIDList := []string{"a411065eedf8bc1830ce544ff51394fe"}
	test.getGroupsInfoReq.GroupIDList = groupIDList
	groupsIDList, _ := json.Marshal(test.GroupIDList)
	fmt.Println("test getGroupsInfo input", string(groupsIDList))
	GetGroupsInfo(string(groupsIDList), test)
}

type testJoinGroup struct {
	joinGroupReq
}

func (testJoinGroup) OnSuccess(data string) {
	fmt.Println("testJoinGroup,onSuccess", data)
}

func (testJoinGroup) OnError(errCode int32, errMsg string) {
	fmt.Println("testJoinGroup,onError", errMsg)
}

func DoJoinGroup() {
	var test testJoinGroup
	test.joinGroupReq.GroupID = "7149948c2fb143f9ee97e3e9b406b5ec"
	test.joinGroupReq.Message = "jin lai "

	fmt.Println("test join group input", test.GroupID, test.Message)
	JoinGroup(test.GroupID, test.Message, test)
}

type testQuitGroup struct {
	quitGroupReq
}

func (testQuitGroup) OnSuccess(data string) {
	fmt.Println("testQuitGroup,onSuccess", data)
}

func (testQuitGroup) OnError(errCode int32, errMsg string) {
	fmt.Println("testQuitGroup,onError", errMsg)
}

func DoQuitGroup() {
	var test testQuitGroup
	test.quitGroupReq.GroupID = "77215e1b13b75f3ab00cb6570e3d9618"

	fmt.Println("test quit group input", test.GroupID)
	QuitGroup(test.GroupID, test)
}

type testGetJoinedGroupList struct {
}

/*
	OnError(errCode int, errMsg string)
	OnSuccess(data string)
*/
func (testGetJoinedGroupList) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetJoinedGroupList OnError", errCode, errMsg)
}

func (testGetJoinedGroupList) OnSuccess(data string) {
	fmt.Println("testGetJoinedGroupList OnSuccess, output", data)
}

func DotestGetJoinedGroupList() {
	var test testGetJoinedGroupList
	GetJoinedGroupList(test)
}

type testGetGroupMemberList struct {
}

func (testGetGroupMemberList) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetGroupMemberList OnError", errCode, errMsg)
}

func (testGetGroupMemberList) OnSuccess(data string) {
	fmt.Println("testGetGroupMemberList OnSuccess, output", data)
}

func DotestGetGroupMemberList() {
	var test testGetGroupMemberList
	var groupId string = ""
	GetGroupMemberList(groupId, 0, 0, test)
}

type testGetGroupMembersInfo struct {
}

func (testGetGroupMembersInfo) OnError(errCode int32, errMsg string) {
	fmt.Println("testGetGroupMembersInfo OnError", errCode, errMsg)
}

func (testGetGroupMembersInfo) OnSuccess(data string) {
	fmt.Println("testGetGroupMembersInfo OnSuccess, output", data)
}

func DotestGetGroupMembersInfo() {
	var test testGetGroupMembersInfo
	var memlist []string
	memlist = append(memlist, "307edc814bb0d04a")
	//memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
	jlist, _ := json.Marshal(memlist)
	fmt.Println("GetGroupMembersInfo input : ", string(jlist))
	GetGroupMembersInfo("7ff61d8f9d4a8a0d6d70a14e2683aad5", string(jlist), test)
	//GetGroupMemberList("05dc84b52829e82242a710ecf999c72c", 0, 0, test)
}

type testKickGroupMember struct {
}

func (testKickGroupMember) OnError(errCode int32, errMsg string) {
	fmt.Println("testKickGroupMember OnError", errCode, errMsg)
}

func (testKickGroupMember) OnSuccess(data string) {
	fmt.Println("testKickGroupMember OnSuccess, output", data)
}

func DotestKickGroupMember() {
	var test testKickGroupMember
	var memlist []string
	//memlist = append(memlist, "e7b437c8b05a1fb8875e7196c636f327")
	memlist = append(memlist, "307edc814bb0d04a")
	jlist, _ := json.Marshal(memlist)

	fmt.Println("KickGroupMember input", string(jlist))
	KickGroupMember("f4cc5c9b556221b92992538f7e6ac26e", "kkkkkkk", string(jlist), test)
}

type testInviteUserToGroup struct {
}

func (testInviteUserToGroup) OnError(errCode int32, errMsg string) {
	fmt.Println("testInviteUserToGroup OnError", errCode, errMsg)
}

func (testInviteUserToGroup) OnSuccess(data string) {
	fmt.Println("testInviteUserToGroup OnSuccess, output", data)
}

func DotesttestInviteUserToGroup() {
	var test testInviteUserToGroup
	var memlist []string
	memlist = append(memlist, "307edc814bb0d04a")
	//memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
	jlist, _ := json.Marshal(memlist)
	fmt.Println("DotesttestInviteUserToGroup, input: ", string(jlist))
	InviteUserToGroup("f4cc5c9b556221b92992538f7e6ac26e", "friend", string(jlist), test)
}

type testGroupX struct {
}

func (testGroupX) OnSuccess(data string) {
	fmt.Println("testGroupX,onSuccess", data)
}

func (testGroupX) OnError(errCode int32, errMsg string) {
	fmt.Println("testGroupX,onError", errMsg)
}
func (testGroupX) OnProgress(progress int) {
	fmt.Println("testGroupX  ", progress)
}

func DoGetGroupApplicationList() string {
	//	var test testGroupX
	fmt.Println("test DoGetGroupApplicationList....")

	return ""
}
func DoGroupApplicationList() {
	var test testGroupX
	fmt.Println("test DoGetGroupApplicationList....")
	GetGroupApplicationList(test)
}
func DoTransferGroupOwner(groupid, userid string) {
	var test testGroupX
	fmt.Println("test DoTransferGroupOwner....")
	TransferGroupOwner(groupid, userid, test)
}
func DoAcceptGroupApplication(uid string) {

	str := DoGetGroupApplicationList()
	var ret groupApplicationResult
	err := json.Unmarshal([]byte(str), &ret)
	if err != nil {
		return
	}
	var app GroupReqListInfo
	for i := 0; i < len(ret.GroupApplicationList); i++ {
		if ret.GroupApplicationList[i].FromUserID == uid {
			app = ret.GroupApplicationList[i]
			break
		}
	}

	v, err := json.Marshal(app)
	if err != nil {
		return
	}

	var test testGroupX
	fmt.Println("accept", string(v))
	AcceptGroupApplication(string(v), "accept", test)
}
func DoRefuseGroupApplication(uid string) {
	str := DoGetGroupApplicationList()
	var ret groupApplicationResult
	err := json.Unmarshal([]byte(str), &ret)
	if err != nil {
		return
	}
	var app GroupReqListInfo
	for i := 0; i < len(ret.GroupApplicationList); i++ {
		if ret.GroupApplicationList[i].FromUserID == uid {
			app = ret.GroupApplicationList[i]
			break
		}
	}

	v, err := json.Marshal(app)
	if err != nil {
		return
	}

	fmt.Println(string(v))

	var test testGroupX
	fmt.Println("refuse", string(v))
	RefuseGroupApplication(string(v), "refuse", test)
}
