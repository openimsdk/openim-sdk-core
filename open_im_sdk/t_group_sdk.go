package open_im_sdk

import (
	"encoding/json"
	"fmt"
)

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
	fmt.Println("testCreateGroup,onSuccess")
}

func (testCreateGroup) OnError(errCode int, errMsg string) {
	fmt.Println("testCreateGroup,onError")
}

func DoCreateGroup() {
	var test testCreateGroup
	test.groupInfo.GroupName = "chat group1"
	test.groupInfo.FaceUrl = "url address"
	test.groupInfo.Introduction = "简介"
	test.groupInfo.Notification = "公告"
	test.members = append(test.members, createGroupMemberInfo{
		Uid:     "fea5bc2bf77813c4",
		SetRole: 0,
	})

	test.members = append(test.members, createGroupMemberInfo{
		Uid:     "ce8e840e8a28df31",
		SetRole: 2,
	})
	groupInfo, _ := json.Marshal(test.groupInfo)
	members, _ := json.Marshal(test.members)
	fmt.Println("create groupInfo input:", string(groupInfo), string(members))
	CreateGroup(string(groupInfo), string(members), test)
}

type testSetGroupInfo struct {
	groupInfo
}

func (testSetGroupInfo) OnSuccess(data string) {
	fmt.Println("testSetGroupInfo,onSuccess")
}

func (testSetGroupInfo) OnError(errCode int, errMsg string) {
	fmt.Println("testSetGroupInfo,onError")
}

func DoSetGroupInfo() {
	var test testSetGroupInfo
	test.groupInfo.GroupId = "d81a54e757fa824be04abd1451ac8c64"
	test.GroupName = "测试群"
	test.Introduction = "这是测试群简介"
	test.Notification = "这是测试群公告"
	test.FaceUrl = "这是测试群头像"
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

func (testGetGroupsInfo) OnError(errCode int, errMsg string) {
	fmt.Println("testGetGroupsInfo,onError", errMsg)
}

func DoGetGroupsInfo() {
	var test testGetGroupsInfo
	groupIDList := []string{"15db7d7cfa97a90113265e77a8fce984", "e13148505c158d0d4b3642942300dbf4"}
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

func (testJoinGroup) OnError(errCode int, errMsg string) {
	fmt.Println("testJoinGroup,onError", errMsg)
}

func DoJoinGroup() {
	var test testJoinGroup
	test.joinGroupReq.GroupID = "d81a54e757fa824be04abd1451ac8c64"
	test.joinGroupReq.Message = "小果子"
	groupID, _ := json.Marshal(test.GroupID)
	message, _ := json.Marshal(test.Message)
	fmt.Println("test join group input", string(groupID), string(message))
	JoinGroup(string(groupID), string(message), test)
}

type testQuitGroup struct {
	quitGroupReq
}

func (testQuitGroup) OnSuccess(data string) {
	fmt.Println("testQuitGroup,onSuccess", data)
}

func (testQuitGroup) OnError(errCode int, errMsg string) {
	fmt.Println("testQuitGroup,onError", errMsg)
}

func DoQuitGroup() {
	var test testQuitGroup
	test.quitGroupReq.GroupID = "bade7d67f2375e36f902a57aa797e588"
	groupID, _ := json.Marshal(test.GroupID)

	fmt.Println("test quit group input", string(groupID))
	QuitGroup(string(groupID), test)
}

type testGetJoinedGroupList struct {
}

/*
	OnError(errCode int, errMsg string)
	OnSuccess(data string)
*/
func (testGetJoinedGroupList) OnError(errCode int, errMsg string) {
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

func (testGetGroupMemberList) OnError(errCode int, errMsg string) {
	fmt.Println("testGetGroupMemberList OnError", errCode, errMsg)
}

func (testGetGroupMemberList) OnSuccess(data string) {
	fmt.Println("testGetGroupMemberList OnSuccess, output", data)
}

func DotestGetGroupMemberList() {
	var test testGetGroupMemberList

	GetGroupMemberList("05dc84b52829e82242a710ecf999c72c", 0, 0, test)
}

type testGetGroupMembersInfo struct {
}

func (testGetGroupMembersInfo) OnError(errCode int, errMsg string) {
	fmt.Println("testGetGroupMembersInfo OnError", errCode, errMsg)
}

func (testGetGroupMembersInfo) OnSuccess(data string) {
	fmt.Println("testGetGroupMembersInfo OnSuccess, output", data)
}

func DotestGetGroupMembersInfo() {
	var test testGetGroupMembersInfo
	var memlist []string
	memlist = append(memlist, "e7b437c8b05a1fb8875e7196c636f327")
	memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
	jlist, _ := json.Marshal(memlist)
	fmt.Println("GetGroupMembersInfo input : ", string(jlist))
	GetGroupMembersInfo("7ff61d8f9d4a8a0d6d70a14e2683aad5", string(jlist), test)
	//GetGroupMemberList("05dc84b52829e82242a710ecf999c72c", 0, 0, test)
}

type testKickGroupMember struct {
}

func (testKickGroupMember) OnError(errCode int, errMsg string) {
	fmt.Println("testKickGroupMember OnError", errCode, errMsg)
}

func (testKickGroupMember) OnSuccess(data string) {
	fmt.Println("testKickGroupMember OnSuccess, output", data)
}

func DotestKickGroupMember() {
	var test testKickGroupMember
	var memlist []string
	memlist = append(memlist, "e7b437c8b05a1fb8875e7196c636f327")
	memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
	jlist, _ := json.Marshal(memlist)

	fmt.Println("KickGroupMember input", string(jlist))
	KickGroupMember("7ff61d8f9d4a8a0d6d70a14e2683aad5", string(jlist), "kkkkkkk", test)
}

type testInviteUserToGroup struct {
}

func (testInviteUserToGroup) OnError(errCode int, errMsg string) {
	fmt.Println("testInviteUserToGroup OnError", errCode, errMsg)
}

func (testInviteUserToGroup) OnSuccess(data string) {
	fmt.Println("testInviteUserToGroup OnSuccess, output", data)
}

func DotesttestInviteUserToGroup() {
	var test testInviteUserToGroup
	var memlist []string
	memlist = append(memlist, "fea5bc2bf77813c4")
	//memlist = append(memlist, "ded01dfe543700402608e19d4e2f839e")
	jlist, _ := json.Marshal(memlist)
	fmt.Println("DotesttestInviteUserToGroup, input: ", string(jlist))
	InviteUserToGroup("d81a54e757fa824be04abd1451ac8c64", "friend", string(jlist), test)
}
