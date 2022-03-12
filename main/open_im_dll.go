package main

//go build -buildmode=c-archive -o openim.lib ./main/open_im_dll.go

/*
#cgo LDFLAGS: -Wl,--unresolved-symbols=ignore-in-object-files -Wl,--allow-shlib-undefined

#include <stdlib.h>

typedef void (*FunWithVoid)();
typedef void (*FunWithInt)(int);
typedef void (*FunWithString)(char *);
typedef void (*FunWithIntString)(int, char *);

static inline void CallFunWithVoid(void* cb) {
    return (*(FunWithVoid)cb)();
}

static inline void CallFunWithInt(void* cb, int i) {
    return (*(FunWithVoid)cb)(i);
}

static inline void CallFunWithString(void* cb, char * s) {
    return (*(FunWithVoid)cb)(s);
}

static inline void CallFunWithIntString(void* cb, int i, char * s) {
    return (*(FunWithVoid)cb)(i, s);
}

*/
import "C"
import (
	"errors"
	"open_im_sdk/internal/login"
	openIM "open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/sdk_struct"
	"unsafe"
)

/*
func main() {

}
*/
func FunWithVoid(cb unsafe.Pointer) {
	C.CallFunWithVoid(cb)
}

func FunWithInt(cb unsafe.Pointer, i int) {
	cint := C.int(i)
	C.CallFunWithInt(cb, cint)
}

func FunWithString(cb unsafe.Pointer, str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.CallFunWithString(cb, cstr)
}

func FunWithIntString(cb unsafe.Pointer, i int, str string) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	C.CallFunWithIntString(cb, C.int(i), cstr)
}

//export SdkVersion
func SdkVersion(verCb unsafe.Pointer) {
	ver := openIM.SdkVersion()
	FunWithString(verCb, ver)
}

//export InitSDK
func InitSDK(onConnecting unsafe.Pointer, onConnectSuccess unsafe.Pointer,
	onConnectFailed unsafe.Pointer, onKickedOffline unsafe.Pointer,
	onUserTokenExpired unsafe.Pointer, operationID *C.char, config *C.char) bool {

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

func Login(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userID, token *C.char) {
	openIM.Login(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userID), C.GoString(token))
}

func UploadImage(onReturn unsafe.Pointer, onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, filePath *C.char, token, obj *C.char) {
	ret := openIM.UploadImage(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(filePath), C.GoString(token), C.GoString(obj))

	FunWithString(onReturn, ret)
}

func Logout(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.Logout(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func GetLoginStatus() int32 {
	return openIM.GetLoginStatus()
}

func GetLoginUser(OnReturn unsafe.Pointer) {
	ret := openIM.GetLoginUser()
	FunWithString(OnReturn, ret)
}

///////////////////////user/////////////////////
func GetUsersInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDList *C.char) {
	openIM.GetUsersInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

func SetSelfInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userInfo *C.char) {
	openIM.SetSelfInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userInfo))
}

func GetSelfUserInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetSelfUserInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//////////////////////////group//////////////////////////////////////////
func SetGroupListener(
	onJoinedGroupAdded unsafe.Pointer,
	onJoinedGroupDeleted unsafe.Pointer,
	onGroupMemberAdded unsafe.Pointer,
	onGroupMemberDeleted unsafe.Pointer,
	onGroupApplicationAdded unsafe.Pointer,
	onGroupApplicationDeleted unsafe.Pointer,
	onGroupInfoChanged unsafe.Pointer,
	onGroupMemberInfoChanged unsafe.Pointer,
	onGroupApplicationAccepted unsafe.Pointer,
	onGroupApplicationRejected unsafe.Pointer,
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

func CreateGroup(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupBaseInfo *C.char, memberList *C.char) {
	openIM.CreateGroup(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupBaseInfo), C.GoString(memberList))
}

func JoinGroup(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID, reqMsg *C.char) {
	openIM.JoinGroup(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reqMsg))
}

func QuitGroup(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char) {
	openIM.QuitGroup(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID))
}

func GetJoinedGroupList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetJoinedGroupList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func GetGroupsInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupIDList *C.char) {
	openIM.GetGroupsInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupIDList))
}

func SetGroupInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char, groupInfo *C.char) {
	openIM.SetGroupInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(groupInfo))
}

func GetGroupMemberList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char, filter, offset, count int32) {
	openIM.GetGroupMemberList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), filter, offset, count)
}

func GetGroupMembersInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char, userIDList *C.char) {
	openIM.GetGroupMembersInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(userIDList))
}

func KickGroupMember(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char, reason *C.char, userIDList *C.char) {
	openIM.KickGroupMember(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reason), C.GoString(userIDList))
}

func TransferGroupOwner(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID, newOwnerUserID *C.char) {
	openIM.TransferGroupOwner(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(newOwnerUserID))
}

func InviteUserToGroup(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID, reason *C.char, userIDList *C.char) {
	openIM.InviteUserToGroup(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(reason), C.GoString(userIDList))
}

func GetRecvGroupApplicationList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetRecvGroupApplicationList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func GetSendGroupApplicationList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetSendGroupApplicationList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func AcceptGroupApplication(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID, fromUserID, handleMsg *C.char) {
	openIM.AcceptGroupApplication(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(fromUserID), C.GoString(handleMsg))
}

func RefuseGroupApplication(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID, fromUserID, handleMsg *C.char) {
	openIM.RefuseGroupApplication(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(groupID), C.GoString(fromUserID), C.GoString(handleMsg))
}

////////////////////////////friend/////////////////////////////////////

func GetDesignatedFriendsInfo(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDList *C.char) {
	openIM.GetDesignatedFriendsInfo(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

func GetFriendList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetFriendList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func CheckFriend(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDList *C.char) {
	openIM.CheckFriend(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDList))
}

func AddFriend(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDReqMsg *C.char) {
	openIM.AddFriend(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDReqMsg))
}

func SetFriendRemark(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDRemark *C.char) {
	openIM.SetFriendRemark(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDRemark))
}
func DeleteFriend(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, friendUserID *C.char) {
	openIM.DeleteFriend(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(friendUserID))
}

func GetRecvFriendApplicationList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetRecvFriendApplicationList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func GetSendFriendApplicationList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetSendFriendApplicationList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func AcceptFriendApplication(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDHandleMsg *C.char) {
	openIM.AcceptFriendApplication(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDHandleMsg))
}

func RefuseFriendApplication(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userIDHandleMsg *C.char) {
	openIM.RefuseFriendApplication(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(userIDHandleMsg))
}

func AddBlack(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, blackUserID *C.char) {
	openIM.AddBlack(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(blackUserID))
}

func GetBlackList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetBlackList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

func RemoveBlack(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, removeUserID *C.char) {
	openIM.RemoveBlack(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(removeUserID))
}

func SetFriendListener(
	onFriendApplicationAdded unsafe.Pointer,
	onFriendApplicationDeleted unsafe.Pointer,
	onFriendApplicationAccepted unsafe.Pointer,
	onFriendApplicationRejected unsafe.Pointer,
	onFriendAdded unsafe.Pointer,
	onFriendDeleted unsafe.Pointer,
	onFriendInfoChanged unsafe.Pointer,
	onBlackAdded unsafe.Pointer,
	onBlackDeleted unsafe.Pointer,
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

func GetAllConversationList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetAllConversationList(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}
func GetConversationListSplit(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, offset, count int) {
	openIM.GetConversationListSplit(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), offset, count)
}

func SetConversationRecvMessageOpt(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationIDList *C.char, opt int) {
	openIM.SetConversationRecvMessageOpt(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList), opt)
}

func GetConversationRecvMessageOpt(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationIDList *C.char) {
	openIM.GetConversationRecvMessageOpt(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList))
}
func GetOneConversation(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, sessionType int, sourceID *C.char) {
	openIM.GetOneConversation(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), sessionType, C.GoString(sourceID))
}
func GetMultipleConversation(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationIDList *C.char) {
	openIM.GetMultipleConversation(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationIDList))
}
func DeleteConversation(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationID *C.char) {
	openIM.DeleteConversation(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID))
}
func SetConversationDraft(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationID, draftText *C.char) {
	openIM.SetConversationDraft(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID), C.GoString(draftText))
}
func PinConversation(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, conversationID *C.char, isPinned bool) {
	openIM.PinConversation(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID), C.GoString(conversationID), isPinned)
}
func GetTotalUnreadMsgCount(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char) {
	openIM.GetTotalUnreadMsgCount(&Base{
		onError: func(i int, s string) {
			FunWithIntString(onError, i, s)
		},
		onSuccess: func(s string) {
			FunWithString(onSuccess, s)
		},
	}, C.GoString(operationID))
}

//
func SetConversationListener(
	onSyncServerStart unsafe.Pointer,
	onSyncServerFinish unsafe.Pointer,
	onSyncServerFailed unsafe.Pointer,
	onNewConversation unsafe.Pointer,
	onConversationChanged unsafe.Pointer,
	onTotalUnreadMessageCountChanged unsafe.Pointer,
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
func SetAdvancedMsgListener(
	onRecvNewMessage unsafe.Pointer,
	onRecvC2CReadReceipt unsafe.Pointer,
	onRecvMessageRevoked unsafe.Pointer,
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

func SetUserListener(onSelfInfoUpdated unsafe.Pointer) {
	openIM.SetUserListener(&OnUserListener{
		onSelfInfoUpdated: func(s string) {
			FunWithString(onSelfInfoUpdated, s)
		},
	})
}

func CreateTextAtMessage(onReturn unsafe.Pointer, operationID *C.char, text, atUserList *C.char) {
	ret := openIM.CreateTextAtMessage(C.GoString(operationID), C.GoString(text), C.GoString(atUserList))
	FunWithString(onReturn, ret)
}

//
func CreateTextMessage(operationID *C.char, text *C.char) *C.char {
	return userForSDK.Conversation().CreateTextMessage(text, operationID)
}

func CreateLocationMessage(operationID *C.char, description *C.char, longitude, latitude float64) *C.char {
	return userForSDK.Conversation().CreateLocationMessage(description, longitude, latitude, operationID)
}
func CreateCustomMessage(operationID *C.char, data, extension *C.char, description *C.char) *C.char {
	return userForSDK.Conversation().CreateCustomMessage(data, extension, description, operationID)
}
func CreateQuoteMessage(operationID *C.char, text *C.char, message *C.char) *C.char {
	return userForSDK.Conversation().CreateQuoteMessage(text, message, operationID)
}
func CreateCardMessage(operationID *C.char, cardInfo *C.char) *C.char {
	return userForSDK.Conversation().CreateCardMessage(cardInfo, operationID)

}
func CreateVideoMessageFromFullPath(operationID *C.char, videoFullPath *C.char, videoType *C.char, duration int64, snapshotFullPath *C.char) *C.char {
	return userForSDK.Conversation().CreateVideoMessageFromFullPath(videoFullPath, videoType, duration, snapshotFullPath, operationID)
}
func CreateImageMessageFromFullPath(operationID *C.char, imageFullPath *C.char) *C.char {
	return userForSDK.Conversation().CreateImageMessageFromFullPath(imageFullPath, operationID)
}
func CreateSoundMessageFromFullPath(operationID *C.char, soundPath *C.char, duration int64) *C.char {
	return userForSDK.Conversation().CreateSoundMessageFromFullPath(soundPath, duration, operationID)
}
func CreateFileMessageFromFullPath(operationID *C.char, fileFullPath, fileName *C.char) *C.char {
	return userForSDK.Conversation().CreateFileMessageFromFullPath(fileFullPath, fileName, operationID)
}
func CreateImageMessage(operationID *C.char, imagePath *C.char) *C.char {
	return userForSDK.Conversation().CreateImageMessage(imagePath, operationID)
}
func CreateImageMessageByURL(operationID *C.char, sourcePicture, bigPicture, snapshotPicture *C.char) *C.char {
	return userForSDK.Conversation().CreateImageMessageByURL(sourcePicture, bigPicture, snapshotPicture, operationID)
}

func CreateSoundMessageByURL(operationID *C.char, soundBaseInfo *C.char) *C.char {
	return userForSDK.Conversation().CreateSoundMessageByURL(soundBaseInfo, operationID)
}
func CreateSoundMessage(operationID *C.char, soundPath *C.char, duration int64) *C.char {
	return userForSDK.Conversation().CreateSoundMessage(soundPath, duration, operationID)
}
func CreateVideoMessageByURL(operationID *C.char, videoBaseInfo *C.char) *C.char {
	return userForSDK.Conversation().CreateVideoMessageByURL(videoBaseInfo, operationID)
}
func CreateVideoMessage(operationID *C.char, videoPath *C.char, videoType *C.char, duration int64, snapshotPath *C.char) *C.char {
	return userForSDK.Conversation().CreateVideoMessage(videoPath, videoType, duration, snapshotPath, operationID)
}
func CreateFileMessageByURL(operationID *C.char, fileBaseInfo *C.char) *C.char {
	return userForSDK.Conversation().CreateFileMessageByURL(fileBaseInfo, operationID)
}
func CreateFileMessage(operationID *C.char, filePath *C.char, fileName *C.char) *C.char {
	return userForSDK.Conversation().CreateFileMessage(filePath, fileName, operationID)
}
func CreateMergerMessage(operationID *C.char, messageList, title, summaryList *C.char) *C.char {
	return userForSDK.Conversation().CreateMergerMessage(messageList, title, summaryList, operationID)
}
func CreateForwardMessage(operationID *C.char, m *C.char) *C.char {
	return userForSDK.Conversation().CreateForwardMessage(m, operationID)
}

func SendMessage(callback open_im_sdk_callback.SendMsgCallBack, operationID, message, recvID, groupID, offlinePushInfo *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SendMessage(callback, message, recvID, groupID, offlinePushInfo, operationID)
}
func SendMessageNotOss(callback open_im_sdk_callback.SendMsgCallBack, operationID *C.char, message, recvID, groupID *C.char, offlinePushInfo *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SendMessageNotOss(callback, message, recvID, groupID, offlinePushInfo, operationID)
}

func GetHistoryMessageList(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, getMessageOptions *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().GetHistoryMessageList(callback, getMessageOptions, operationID)
}

func RevokeMessage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, message *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().RevokeMessage(callback, message, operationID)
}
func TypingStatusUpdate(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, recvID, msgTip *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().TypingStatusUpdate(callback, recvID, msgTip, operationID)
}
func MarkC2CMessageAsRead(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userID *C.char, msgIDList *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkC2CMessageAsRead(callback, userID, msgIDList, operationID)
}

func MarkGroupMessageHasRead(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().MarkGroupMessageHasRead(callback, groupID, operationID)
}
func DeleteMessageFromLocalStorage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, message *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().DeleteMessageFromLocalStorage(callback, message, operationID)
}
func ClearC2CHistoryMessage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, userID *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearC2CHistoryMessage(callback, userID, operationID)
}
func ClearGroupHistoryMessage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, groupID *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().ClearGroupHistoryMessage(callback, groupID, operationID)
}
func InsertSingleMessageToLocalStorage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, message, recvID, sendID *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().InsertSingleMessageToLocalStorage(callback, message, recvID, sendID, operationID)
}
func InsertGroupMessageToLocalStorage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, message, groupID, sendID *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().InsertGroupMessageToLocalStorage(callback, message, groupID, sendID, operationID)
}
func SearchLocalMessages(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, searchParam *C.char) {
	if err := CheckResourceLoad(userForSDK); err != nil {
		log.Error(operationID, "resource loading is not completed ", err.Error())
		callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return
	}
	userForSDK.Conversation().SearchLocalMessages(callback, searchParam, operationID)
}

//func FindMessages(callback common.Base, operationID *C.char, messageIDList *C.char) {
//	userForSDK.Conversation().FindMessages(callback, messageIDList)
//}

func InitOnce(config *sdk_struct.IMConfig) bool {
	sdk_struct.SvrConf = *config
	return true
}

func CheckToken(userID, token *C.char) error {
	return login.CheckToken(userID, token, "")
}

func CheckResourceLoad(uSDK *login.LoginMgr) error {
	if uSDK == nil {
		//	callback.OnError(constant.ErrResourceLoadNotComplete.ErrCode, constant.ErrResourceLoadNotComplete.ErrMsg)
		return utils.Wrap(errors.New("CheckResourceLoad failed uSDK == nil "), "")
	}
	if uSDK.Friend() == nil || uSDK.User() == nil || uSDK.Group() == nil || uSDK.Conversation() == nil ||
		uSDK.Full() == nil {
		return utils.Wrap(errors.New("CheckResourceLoad failed, resource nil "), "")
	}
	return nil
}

func uploadImage(onError unsafe.Pointer, onSuccess unsafe.Pointer, operationID *C.char, filePath *C.char, token, obj *C.char) *C.char {
	if obj == "cos" {
		p := ws.NewPostApi(token, userForSDK.ImConfig().ApiAddr)
		o := common2.NewCOS(p)
		url, _, err := o.UploadFile(filePath, func(progress int) {
			if progress == 100 {
				callback.OnSuccess("")
			}
		})

		if err != nil {
			callback.OnError(100, err.Error())
			return ""
		}
		return url

	} else {
		return ""
	}
}
func GetConversationIDBySessionType(sourceID *C.char, sessionType int) *C.char {
	return utils.GetConversationIDBySessionType(sourceID, sessionType)
}
