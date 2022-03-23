package main

//go build -buildmode=c-archive -o openim.lib ./main/open_im_dll.go

/*
#cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-in-object-files -Wl,--allow-shlib-undefined

#include <stdlib.h>

typedef void (*FunWithVoid)();
typedef void (*FunWithInt)(int);
typedef void (*FunWithString)(char *);
typedef void (*FunWithIntString)(int, char *);

static inline void CallFunWithVoid(FunWithVoid cb) {
    return cb();
}

static inline void CallFunWithInt(FunWithInt cb, int i) {
    return cb(i);
}

static inline void CallFunWithString(FunWithString cb, char * s) {
    return cb(s);
}

static inline void CallFunWithIntString(FunWithIntString cb, int i, char * s) {
    return cb(i, s);
}

*/
import "C"
import (
	"open_im_sdk/internal/login"
	openIM "open_im_sdk/open_im_sdk"
	"open_im_sdk/sdk_struct"
	"unsafe"
)

func main() {

}

type GoFunVoid = func()
type GoFunInt = func(int)
type GoFunString = func(string)
type GoFunIntString = func(int, string)

type BaseCallBack struct {
	onError   GoFunIntString
	onSuccess GoFunString
}

func (s *BaseCallBack) OnError(errCode int32, errMsg string) {
	s.onError(int(errCode), errMsg)
}

func (s *BaseCallBack) OnSuccess(data string) {
	s.onSuccess(data)
}

type SendMsgCallBack struct {
	onError    GoFunIntString
	onSuccess  GoFunString
	onProgress GoFunInt
}

func (s *SendMsgCallBack) OnError(errCode int32, errMsg string) {
	s.onError(int(errCode), errMsg)
}

func (s *SendMsgCallBack) OnSuccess(data string) {
	s.onSuccess(data)
}

func (s *SendMsgCallBack) OnProgress(progress int) {
	s.onProgress(progress)
}

type OnConnListener struct {
	onConnecting       GoFunVoid
	onConnectSuccess   GoFunVoid
	onConnectFailed    GoFunIntString
	onKickedOffline    GoFunVoid
	onUserTokenExpired GoFunVoid
}

func (s *OnConnListener) OnConnecting() {
	s.onConnecting()
}

func (s *OnConnListener) OnConnectSuccess() {
	s.onConnectSuccess()
}

func (s *OnConnListener) OnConnectFailed(errCode int32, errMsg string) {
	s.onConnectFailed(int(errCode), errMsg)
}

func (s *OnConnListener) OnKickedOffline() {
	s.onKickedOffline()
}

func (s *OnConnListener) OnUserTokenExpired() {
	s.onUserTokenExpired()
}

type OnGroupListener struct {
	onJoinedGroupAdded         GoFunString
	onJoinedGroupDeleted       GoFunString
	onGroupMemberAdded         GoFunString
	onGroupMemberDeleted       GoFunString
	onGroupApplicationAdded    GoFunString
	onGroupApplicationDeleted  GoFunString
	onGroupInfoChanged         GoFunString
	onGroupMemberInfoChanged   GoFunString
	onGroupApplicationAccepted GoFunString
	onGroupApplicationRejected GoFunString
}

func (s *OnGroupListener) OnJoinedGroupAdded(groupInfo string) {
	s.onJoinedGroupAdded(groupInfo)
}

func (s *OnGroupListener) OnJoinedGroupDeleted(groupInfo string) {
	s.onJoinedGroupDeleted(groupInfo)
}
func (s *OnGroupListener) OnGroupMemberAdded(groupMemberInfo string) {
	s.onGroupMemberAdded(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupMemberDeleted(groupMemberInfo string) {
	s.onGroupMemberDeleted(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupApplicationAdded(groupApplication string) {
	s.onGroupApplicationAdded(groupApplication)
}
func (s *OnGroupListener) OnGroupApplicationDeleted(groupApplication string) {
	s.onGroupApplicationDeleted(groupApplication)
}
func (s *OnGroupListener) OnGroupInfoChanged(groupInfo string) {
	s.onGroupInfoChanged(groupInfo)
}
func (s *OnGroupListener) OnGroupMemberInfoChanged(groupMemberInfo string) {
	s.onGroupMemberInfoChanged(groupMemberInfo)
}
func (s *OnGroupListener) OnGroupApplicationAccepted(groupApplication string) {
	s.onGroupApplicationAccepted(groupApplication)
}
func (s *OnGroupListener) OnGroupApplicationRejected(groupApplication string) {
	s.onGroupApplicationRejected(groupApplication)
}

type OnFriendshipListener struct {
	onFriendApplicationAdded    GoFunString
	onFriendApplicationDeleted  GoFunString
	onFriendApplicationAccepted GoFunString
	onFriendApplicationRejected GoFunString
	onFriendAdded               GoFunString
	onFriendDeleted             GoFunString
	onFriendInfoChanged         GoFunString
	onBlackAdded                GoFunString
	onBlackDeleted              GoFunString
}

func (s *OnFriendshipListener) OnFriendApplicationAdded(friendApplication string) {
	s.onFriendApplicationAdded(friendApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationDeleted(friendApplication string) {
	s.onFriendApplicationDeleted(friendApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationAccepted(groupApplication string) {
	s.onFriendApplicationAccepted(groupApplication)
}
func (s *OnFriendshipListener) OnFriendApplicationRejected(friendApplication string) {
	s.onFriendApplicationRejected(friendApplication)
}
func (s *OnFriendshipListener) OnFriendAdded(friendInfo string) {
	s.onFriendAdded(friendInfo)
}
func (s *OnFriendshipListener) OnFriendDeleted(friendInfo string) {
	s.onFriendDeleted(friendInfo)
}
func (s *OnFriendshipListener) OnFriendInfoChanged(friendInfo string) {
	s.onFriendInfoChanged(friendInfo)
}
func (s *OnFriendshipListener) OnBlackAdded(blackInfo string) {
	s.onBlackAdded(blackInfo)
}
func (s *OnFriendshipListener) OnBlackDeleted(blackInfo string) {
	s.onBlackDeleted(blackInfo)
}

type OnConversationListener struct {
	onSyncServerStart                GoFunVoid
	onSyncServerFinish               GoFunVoid
	onSyncServerFailed               GoFunVoid
	onNewConversation                GoFunString
	onConversationChanged            GoFunString
	onTotalUnreadMessageCountChanged GoFunInt
}

func (s *OnConversationListener) OnSyncServerStart() {
	s.onSyncServerStart()
}
func (s *OnConversationListener) OnSyncServerFinish() {
	s.onSyncServerFinish()
}
func (s *OnConversationListener) OnSyncServerFailed() {
	s.onSyncServerFailed()
}
func (s *OnConversationListener) OnNewConversation(conversationList string) {
	s.onNewConversation(conversationList)
}
func (s *OnConversationListener) OnConversationChanged(conversationList string) {
	s.onConversationChanged(conversationList)
}
func (s *OnConversationListener) OnTotalUnreadMessageCountChanged(totalUnreadCount int32) {
	s.onTotalUnreadMessageCountChanged(int(totalUnreadCount))
}

type OnAdvancedMsgListener struct {
	onRecvNewMessage     GoFunString
	onRecvC2CReadReceipt GoFunString
	onRecvMessageRevoked GoFunString
}

func (s *OnAdvancedMsgListener) OnRecvNewMessage(message string) {
	s.onRecvNewMessage(message)
}
func (s *OnAdvancedMsgListener) OnRecvC2CReadReceipt(msgReceiptList string) {
	s.onRecvC2CReadReceipt(msgReceiptList)
}
func (s *OnAdvancedMsgListener) OnRecvMessageRevoked(msgId string) {
	s.onRecvMessageRevoked(msgId)
}

type OnUserListener struct {
	onSelfInfoUpdated GoFunString
}

func (s *OnUserListener) OnSelfInfoUpdated(userInfo string) {
	s.onSelfInfoUpdated(userInfo)
}

func FunWithVoid(cb C.FunWithVoid) {
	C.CallFunWithVoid(cb)
}

func FunWithInt(cb C.FunWithInt, i int) {
	cint := C.int(i)
	C.CallFunWithInt(cb, cint)
}

func FunWithString(cb C.FunWithString, str string) {
	//go会保证在c函数返回之前不移动内存
	var bstr = []byte(str)
	cstr := (*C.char)(unsafe.Pointer(&bstr[0]))
	//cstr := C.CString(str)
	//defer C.free(unsafe.Pointer(cstr))

	C.CallFunWithString(cb, cstr)
}

func FunWithIntString(cb C.FunWithIntString, i int, str string) {
	var bstr = []byte(str)
	cstr := (*C.char)(unsafe.Pointer(&bstr[0]))
	//cstr := C.CString(str)
	//defer C.free(unsafe.Pointer(cstr))

	C.CallFunWithIntString(cb, C.int(i), cstr)
}

//export SdkVersion
func SdkVersion(verCb C.FunWithString) {
	ver := openIM.SdkVersion()
	FunWithString(verCb, ver)
}

//export InitSDK
func InitSDK(onConnecting C.FunWithVoid, onConnectSuccess C.FunWithVoid,
	onConnectFailed C.FunWithIntString, onKickedOffline C.FunWithVoid,
	onUserTokenExpired C.FunWithVoid, operationID *C.char, config *C.char) bool {

	return openIM.InitSDK(&OnConnListener{
		onConnecting: func() {
			FunWithVoid(onConnecting)
		},
		onConnectSuccess: func() {
			FunWithVoid(onConnectSuccess)
		},
		onConnectFailed: func(errCode int, errMsg string) {
			FunWithIntString(onConnectFailed, int(errCode), errMsg)
		},
		onKickedOffline: func() {
			FunWithVoid(onKickedOffline)
		},
		onUserTokenExpired: func() {
			FunWithVoid(onUserTokenExpired)
		},
	},
		C.GoString(operationID),
		C.GoString(config),
	)
}

//export Login
func Login(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userID, token *C.char) {
	openIM.Login(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userID), C.GoString(token))
}

//export UploadImage
func UploadImage(onReturn C.FunWithString, onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, filePath *C.char, token, obj *C.char) {
	ret := openIM.UploadImage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(filePath), C.GoString(token), C.GoString(obj))

	FunWithString(onReturn, ret)
}

//export Logout
func Logout(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.Logout(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export GetLoginStatus
func GetLoginStatus() int32 {
	return openIM.GetLoginStatus()
}

//export GetLoginUser
func GetLoginUser(OnReturn C.FunWithString) {
	ret := openIM.GetLoginUser()
	FunWithString(OnReturn, ret)
}

///////////////////////user/////////////////////
//export GetUsersInfo
func GetUsersInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDList *C.char) {
	openIM.GetUsersInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

//export SetSelfInfo
func SetSelfInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userInfo *C.char) {
	openIM.SetSelfInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userInfo))
}

//export GetSelfUserInfo
func GetSelfUserInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetSelfUserInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//////////////////////////group//////////////////////////////////////////
//export SetGroupListener
func SetGroupListener(
	onJoinedGroupAdded C.FunWithString,
	onJoinedGroupDeleted C.FunWithString,
	onGroupMemberAdded C.FunWithString,
	onGroupMemberDeleted C.FunWithString,
	onGroupApplicationAdded C.FunWithString,
	onGroupApplicationDeleted C.FunWithString,
	onGroupInfoChanged C.FunWithString,
	onGroupMemberInfoChanged C.FunWithString,
	onGroupApplicationAccepted C.FunWithString,
	onGroupApplicationRejected C.FunWithString,
) {
	openIM.SetGroupListener(&OnGroupListener{
		onJoinedGroupAdded: func(s string) {
			FunWithString(onJoinedGroupAdded, s)
		},
		onJoinedGroupDeleted: func(s string) {
			FunWithString(onJoinedGroupDeleted, s)
		},
		onGroupMemberAdded: func(s string) {
			FunWithString(onGroupMemberAdded, s)
		},
		onGroupMemberDeleted: func(s string) {
			FunWithString(onGroupMemberDeleted, s)
		},
		onGroupApplicationAdded: func(s string) {
			FunWithString(onGroupApplicationAdded, s)
		},
		onGroupApplicationDeleted: func(s string) {
			FunWithString(onGroupApplicationDeleted, s)
		},
		onGroupInfoChanged: func(s string) {
			FunWithString(onGroupInfoChanged, s)
		},
		onGroupMemberInfoChanged: func(s string) {
			FunWithString(onGroupMemberInfoChanged, s)
		},
		onGroupApplicationAccepted: func(s string) {
			FunWithString(onGroupApplicationAccepted, s)
		},
		onGroupApplicationRejected: func(s string) {
			FunWithString(onGroupApplicationRejected, s)
		},
	})
}

//export CreateGroup
func CreateGroup(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupBaseInfo *C.char, memberList *C.char) {
	openIM.CreateGroup(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupBaseInfo), C.GoString(memberList))
}

//export JoinGroup
func JoinGroup(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID, reqMsg *C.char) {
	openIM.JoinGroup(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reqMsg))
}

//export QuitGroup
func QuitGroup(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char) {
	openIM.QuitGroup(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID))
}

//export GetJoinedGroupList
func GetJoinedGroupList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetJoinedGroupList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export GetGroupsInfo
func GetGroupsInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupIDList *C.char) {
	openIM.GetGroupsInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupIDList))
}

//export SetGroupInfo
func SetGroupInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char, groupInfo *C.char) {
	openIM.SetGroupInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(groupInfo))
}

//export GetGroupMemberList
func GetGroupMemberList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char, filter, offset, count int32) {
	openIM.GetGroupMemberList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), filter, offset, count)
}

//export GetGroupMembersInfo
func GetGroupMembersInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char, userIDList *C.char) {
	openIM.GetGroupMembersInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(userIDList))
}

//export KickGroupMember
func KickGroupMember(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char, reason *C.char, userIDList *C.char) {
	openIM.KickGroupMember(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reason), C.GoString(userIDList))
}

//export TransferGroupOwner
func TransferGroupOwner(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID, newOwnerUserID *C.char) {
	openIM.TransferGroupOwner(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(newOwnerUserID))
}

//export InviteUserToGroup
func InviteUserToGroup(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID, reason *C.char, userIDList *C.char) {
	openIM.InviteUserToGroup(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reason), C.GoString(userIDList))
}

//export GetRecvGroupApplicationList
func GetRecvGroupApplicationList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetRecvGroupApplicationList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export GetSendGroupApplicationList
func GetSendGroupApplicationList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetSendGroupApplicationList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export AcceptGroupApplication
func AcceptGroupApplication(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID, fromUserID, handleMsg *C.char) {
	openIM.AcceptGroupApplication(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(fromUserID), C.GoString(handleMsg))
}

//export RefuseGroupApplication
func RefuseGroupApplication(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID, fromUserID, handleMsg *C.char) {
	openIM.RefuseGroupApplication(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(fromUserID), C.GoString(handleMsg))
}

////////////////////////////friend/////////////////////////////////////
//export GetDesignatedFriendsInfo
func GetDesignatedFriendsInfo(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDList *C.char) {
	openIM.GetDesignatedFriendsInfo(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

//export GetFriendList
func GetFriendList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetFriendList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export CheckFriend
func CheckFriend(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDList *C.char) {
	openIM.CheckFriend(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

//export AddFriend
func AddFriend(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDReqMsg *C.char) {
	openIM.AddFriend(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDReqMsg))
}

//export SetFriendRemark
func SetFriendRemark(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDRemark *C.char) {
	openIM.SetFriendRemark(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDRemark))
}

//export DeleteFriend
func DeleteFriend(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, friendUserID *C.char) {
	openIM.DeleteFriend(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(friendUserID))
}

//export GetRecvFriendApplicationList
func GetRecvFriendApplicationList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetRecvFriendApplicationList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export GetSendFriendApplicationList
func GetSendFriendApplicationList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetSendFriendApplicationList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export AcceptFriendApplication
func AcceptFriendApplication(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDHandleMsg *C.char) {
	openIM.AcceptFriendApplication(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDHandleMsg))
}

//export RefuseFriendApplication
func RefuseFriendApplication(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userIDHandleMsg *C.char) {
	openIM.RefuseFriendApplication(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDHandleMsg))
}

//export AddBlack
func AddBlack(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, blackUserID *C.char) {
	openIM.AddBlack(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(blackUserID))
}

//export GetBlackList
func GetBlackList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetBlackList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export RemoveBlack
func RemoveBlack(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, removeUserID *C.char) {
	openIM.RemoveBlack(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(removeUserID))
}

//export SetFriendListener
func SetFriendListener(
	onFriendApplicationAdded C.FunWithString,
	onFriendApplicationDeleted C.FunWithString,
	onFriendApplicationAccepted C.FunWithString,
	onFriendApplicationRejected C.FunWithString,
	onFriendAdded C.FunWithString,
	onFriendDeleted C.FunWithString,
	onFriendInfoChanged C.FunWithString,
	onBlackAdded C.FunWithString,
	onBlackDeleted C.FunWithString,
) {
	openIM.SetFriendListener(&OnFriendshipListener{
		onFriendApplicationAdded: func(s string) {
			FunWithString(onFriendApplicationAdded, s)
		},
		onFriendApplicationDeleted: func(s string) {
			FunWithString(onFriendApplicationDeleted, s)
		},
		onFriendApplicationAccepted: func(s string) {
			FunWithString(onFriendApplicationAccepted, s)
		},
		onFriendApplicationRejected: func(s string) {
			FunWithString(onFriendApplicationRejected, s)
		},
		onFriendAdded: func(s string) {
			FunWithString(onFriendAdded, s)
		},
		onFriendDeleted: func(s string) {
			FunWithString(onFriendDeleted, s)
		},
		onFriendInfoChanged: func(s string) {
			FunWithString(onFriendInfoChanged, s)
		},
		onBlackAdded: func(s string) {
			FunWithString(onBlackAdded, s)
		},
		onBlackDeleted: func(s string) {
			FunWithString(onBlackDeleted, s)
		},
	})
}

///////////////////////conversation////////////////////////////////////
//export GetAllConversationList
func GetAllConversationList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetAllConversationList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export GetConversationListSplit
func GetConversationListSplit(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, offset, count int) {
	openIM.GetConversationListSplit(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), offset, count)
}

//export SetConversationRecvMessageOpt
func SetConversationRecvMessageOpt(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationIDList *C.char, opt int) {
	openIM.SetConversationRecvMessageOpt(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList), opt)
}

//export GetConversationRecvMessageOpt
func GetConversationRecvMessageOpt(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationIDList *C.char) {
	openIM.GetConversationRecvMessageOpt(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList))
}

//export GetOneConversation
func GetOneConversation(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, sessionType int, sourceID *C.char) {
	openIM.GetOneConversation(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), sessionType, C.GoString(sourceID))
}

//export GetMultipleConversation
func GetMultipleConversation(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationIDList *C.char) {
	openIM.GetMultipleConversation(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList))
}

//export DeleteConversation
func DeleteConversation(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationID *C.char) {
	openIM.DeleteConversation(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID))
}

//export SetConversationDraft
func SetConversationDraft(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationID, draftText *C.char) {
	openIM.SetConversationDraft(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID), C.GoString(draftText))
}

//export PinConversation
func PinConversation(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, conversationID *C.char, isPinned bool) {
	openIM.PinConversation(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID), isPinned)
}

//export GetTotalUnreadMsgCount
func GetTotalUnreadMsgCount(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char) {
	openIM.GetTotalUnreadMsgCount(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//export SetConversationListener
func SetConversationListener(
	onSyncServerStart C.FunWithVoid,
	onSyncServerFinish C.FunWithVoid,
	onSyncServerFailed C.FunWithVoid,
	onNewConversation C.FunWithString,
	onConversationChanged C.FunWithString,
	onTotalUnreadMessageCountChanged C.FunWithInt,
) {
	openIM.SetConversationListener(&OnConversationListener{
		onSyncServerStart:  func() { FunWithVoid(onSyncServerStart) },
		onSyncServerFinish: func() { FunWithVoid(onSyncServerFinish) },
		onSyncServerFailed: func() {
			FunWithVoid(onSyncServerFailed)
		},
		onNewConversation: func(s string) {
			FunWithString(onNewConversation, s)
		},
		onConversationChanged: func(s string) {
			FunWithString(onConversationChanged, s)
		},
		onTotalUnreadMessageCountChanged: func(i int) {
			FunWithInt(onTotalUnreadMessageCountChanged, i)
		},
	})
}

//export SetAdvancedMsgListener
func SetAdvancedMsgListener(
	onRecvNewMessage C.FunWithString,
	onRecvC2CReadReceipt C.FunWithString,
	onRecvMessageRevoked C.FunWithString,
) {
	openIM.SetAdvancedMsgListener(
		&OnAdvancedMsgListener{
			onRecvNewMessage: func(s string) {
				FunWithString(onRecvNewMessage, s)
			},
			onRecvC2CReadReceipt: func(s string) {
				FunWithString(onRecvC2CReadReceipt, s)
			},
			onRecvMessageRevoked: func(s string) {
				FunWithString(onRecvMessageRevoked, s)
			},
		},
	)
}

//export SetUserListener
func SetUserListener(onSelfInfoUpdated C.FunWithString) {
	openIM.SetUserListener(&OnUserListener{
		onSelfInfoUpdated: func(s string) {
			FunWithString(onSelfInfoUpdated, s)
		},
	})
}

//export CreateTextAtMessage
func CreateTextAtMessage(onReturn C.FunWithString, operationID *C.char, text, atUserList *C.char) {
	ret := openIM.CreateTextAtMessage(C.GoString(operationID), C.GoString(text), C.GoString(atUserList))
	FunWithString(onReturn, ret)
}

//export CreateTextMessage
func CreateTextMessage(onReturn C.FunWithString, operationID *C.char, text *C.char) {
	ret := openIM.CreateTextMessage(C.GoString(operationID), C.GoString(text))
	FunWithString(onReturn, ret)
}

//export CreateLocationMessage
func CreateLocationMessage(onReturn C.FunWithString, operationID *C.char, description *C.char, longitude, latitude float64) {
	ret := openIM.CreateLocationMessage(C.GoString(operationID), C.GoString(description), longitude, latitude)
	FunWithString(onReturn, ret)
}

//export CreateCustomMessage
func CreateCustomMessage(onReturn C.FunWithString, operationID *C.char, data, extension *C.char, description *C.char) {
	ret := openIM.CreateCustomMessage(C.GoString(operationID), C.GoString(data), C.GoString(extension), C.GoString(description))
	FunWithString(onReturn, ret)
}

//export CreateQuoteMessage
func CreateQuoteMessage(onReturn C.FunWithString, operationID *C.char, text *C.char, message *C.char) {
	ret := openIM.CreateQuoteMessage(C.GoString(operationID), C.GoString(text), C.GoString(message))
	FunWithString(onReturn, ret)
}

//export CreateCardMessage
func CreateCardMessage(onReturn C.FunWithString, operationID *C.char, cardInfo *C.char) {
	ret := openIM.CreateCardMessage(C.GoString(operationID), C.GoString(cardInfo))
	FunWithString(onReturn, ret)
}

//export CreateVideoMessageFromFullPath
func CreateVideoMessageFromFullPath(onReturn C.FunWithString, operationID *C.char, videoFullPath *C.char, videoType *C.char, duration int64, snapshotFullPath *C.char) {
	ret := openIM.CreateVideoMessageFromFullPath(C.GoString(operationID), C.GoString(videoFullPath), C.GoString(videoType), duration, C.GoString(snapshotFullPath))
	FunWithString(onReturn, ret)
}

//export CreateImageMessageFromFullPath
func CreateImageMessageFromFullPath(onReturn C.FunWithString, operationID *C.char, imageFullPath *C.char) {
	ret := openIM.CreateImageMessageFromFullPath(C.GoString(operationID), C.GoString(imageFullPath))
	FunWithString(onReturn, ret)
}

//export CreateSoundMessageFromFullPath
func CreateSoundMessageFromFullPath(onReturn C.FunWithString, operationID *C.char, soundPath *C.char, duration int64) {
	ret := openIM.CreateSoundMessageFromFullPath(C.GoString(operationID), C.GoString(soundPath), duration)
	FunWithString(onReturn, ret)
}

//export CreateFileMessageFromFullPath
func CreateFileMessageFromFullPath(onReturn C.FunWithString, operationID *C.char, fileFullPath, fileName *C.char) {
	ret := openIM.CreateFileMessageFromFullPath(C.GoString(operationID), C.GoString(fileFullPath), C.GoString(fileName))
	FunWithString(onReturn, ret)
}

//export CreateImageMessage
func CreateImageMessage(onReturn C.FunWithString, operationID *C.char, imagePath *C.char) {
	ret := openIM.CreateImageMessage(C.GoString(operationID), C.GoString(imagePath))
	FunWithString(onReturn, ret)
}

//export CreateImageMessageByURL
func CreateImageMessageByURL(onReturn C.FunWithString, operationID *C.char, sourcePicture, bigPicture, snapshotPicture *C.char) {
	ret := openIM.CreateImageMessageByURL(C.GoString(operationID), C.GoString(sourcePicture), C.GoString(bigPicture), C.GoString(snapshotPicture))
	FunWithString(onReturn, ret)
}

//export CreateSoundMessageByURL
func CreateSoundMessageByURL(onReturn C.FunWithString, operationID *C.char, soundBaseInfo *C.char) {
	ret := openIM.CreateSoundMessageByURL(C.GoString(operationID), C.GoString(soundBaseInfo))
	FunWithString(onReturn, ret)
}

//export CreateSoundMessage
func CreateSoundMessage(onReturn C.FunWithString, operationID *C.char, soundPath *C.char, duration int64) {
	ret := openIM.CreateSoundMessage(C.GoString(operationID), C.GoString(soundPath), duration)
	FunWithString(onReturn, ret)
}

//export CreateVideoMessageByURL
func CreateVideoMessageByURL(onReturn C.FunWithString, operationID *C.char, videoBaseInfo *C.char) {
	ret := openIM.CreateVideoMessageByURL(C.GoString(operationID), C.GoString(videoBaseInfo))
	FunWithString(onReturn, ret)
}

//export CreateVideoMessage
func CreateVideoMessage(onReturn C.FunWithString, operationID *C.char, videoPath *C.char, videoType *C.char, duration int64, snapshotPath *C.char) {
	ret := openIM.CreateVideoMessage(C.GoString(operationID), C.GoString(videoPath), C.GoString(videoType), duration, C.GoString(snapshotPath))
	FunWithString(onReturn, ret)
}

//export CreateFileMessageByURL
func CreateFileMessageByURL(onReturn C.FunWithString, operationID *C.char, fileBaseInfo *C.char) {
	ret := openIM.CreateFileMessageByURL(C.GoString(operationID), C.GoString(fileBaseInfo))
	FunWithString(onReturn, ret)
}

//export CreateFileMessage
func CreateFileMessage(onReturn C.FunWithString, operationID *C.char, filePath *C.char, fileName *C.char) {
	ret := openIM.CreateFileMessage(C.GoString(operationID), C.GoString(filePath), C.GoString(fileName))
	FunWithString(onReturn, ret)
}

//export CreateMergerMessage
func CreateMergerMessage(onReturn C.FunWithString, operationID *C.char, messageList, title, summaryList *C.char) {
	ret := openIM.CreateMergerMessage(C.GoString(operationID), C.GoString(messageList), C.GoString(title), C.GoString(summaryList))
	FunWithString(onReturn, ret)
}

//export CreateForwardMessage
func CreateForwardMessage(onReturn C.FunWithString, operationID *C.char, m *C.char) {
	ret := openIM.CreateForwardMessage(C.GoString(operationID), C.GoString(m))
	FunWithString(onReturn, ret)
}

//export SendMessage
func SendMessage(onError C.FunWithIntString, onSuccess C.FunWithString, onProgress C.FunWithInt, operationID, message, recvID, groupID, offlinePushInfo *C.char) {
	openIM.SendMessage(&SendMsgCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
		onProgress: func(i int) {
			FunWithInt(onProgress, i)
		},
	}, C.GoString(operationID), C.GoString(message), C.GoString(recvID), C.GoString(groupID), C.GoString(offlinePushInfo))
}

//export SendMessageNotOss
func SendMessageNotOss(onError C.FunWithIntString, onSuccess C.FunWithString, onProgress C.FunWithInt, operationID *C.char, message, recvID, groupID *C.char, offlinePushInfo *C.char) {
	openIM.SendMessageNotOss(&SendMsgCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
		onProgress: func(i int) {
			FunWithInt(onProgress, i)
		},
	}, C.GoString(operationID), C.GoString(message), C.GoString(recvID), C.GoString(groupID), C.GoString(offlinePushInfo))
}

//export GetHistoryMessageList
func GetHistoryMessageList(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, getMessageOptions *C.char) {
	openIM.GetHistoryMessageList(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(getMessageOptions))
}

//export RevokeMessage
func RevokeMessage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, message *C.char) {
	openIM.RevokeMessage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(message))
}

//export UpdateMessage
func UpdateMessage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, message *C.char) {
	openIM.UpdateMessage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(message))
}

//export TypingStatusUpdate
func TypingStatusUpdate(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, recvID, msgTip *C.char) {
	openIM.TypingStatusUpdate(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(recvID), C.GoString(msgTip))
}

//export MarkC2CMessageAsRead
func MarkC2CMessageAsRead(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userID *C.char, msgIDList *C.char) {
	openIM.MarkC2CMessageAsRead(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userID), C.GoString(msgIDList))
}

//export MarkGroupMessageHasRead
func MarkGroupMessageHasRead(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char) {
	openIM.MarkGroupMessageHasRead(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID))
}

//export DeleteMessageFromLocalStorage
func DeleteMessageFromLocalStorage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, message *C.char) {
	openIM.DeleteMessageFromLocalStorage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(message))
}

//export func ClearC2CHistoryMessage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userID *C.char) {

func ClearC2CHistoryMessage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, userID *C.char) {
	openIM.ClearC2CHistoryMessage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userID))
}

//export ClearGroupHistoryMessage
func ClearGroupHistoryMessage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, groupID *C.char) {
	openIM.ClearGroupHistoryMessage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID))
}

//export InsertSingleMessageToLocalStorage
func InsertSingleMessageToLocalStorage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, message, recvID, sendID *C.char) {
	openIM.InsertSingleMessageToLocalStorage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(message), C.GoString(recvID), C.GoString(sendID))
}

//export InsertGroupMessageToLocalStorage
func InsertGroupMessageToLocalStorage(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, message, groupID, sendID *C.char) {
	openIM.InsertGroupMessageToLocalStorage(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(message), C.GoString(groupID), C.GoString(sendID))
}

//export SearchLocalMessages
func SearchLocalMessages(onError C.FunWithIntString, onSuccess C.FunWithString, operationID *C.char, searchParam *C.char) {
	openIM.SearchLocalMessages(&BaseCallBack{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(searchParam))
}

//func FindMessages(callback common.Base, operationID *C.char, messageIDList *C.char) {
//	userForSDK.Conversation().FindMessages(callback, messageIDList)
//}

func InitOnce(config *sdk_struct.IMConfig) bool {
	return openIM.InitOnce(config)
}

func CheckToken(userID, token *C.char) error {
	return openIM.CheckToken(C.GoString(userID), C.GoString(token))
}

func CheckResourceLoad(uSDK *login.LoginMgr) error {
	return openIM.CheckResourceLoad(uSDK)
}

//export GetConversationIDBySessionType
func GetConversationIDBySessionType(onReturn C.FunWithString, sourceID *C.char, sessionType int) {
	ret := openIM.GetConversationIDBySessionType(C.GoString(sourceID), sessionType)
	FunWithString(onReturn, ret)
}
