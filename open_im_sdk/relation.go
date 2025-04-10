package open_im_sdk

import "github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"

func GetSpecifiedFriendsInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string, filterBlack bool) {
	call(callback, operationID, IMUserContext.Relation().GetSpecifiedFriendsInfo, userIDList, filterBlack)
}

func GetFriendList(callback open_im_sdk_callback.Base, operationID string, filterBlack bool) {
	call(callback, operationID, IMUserContext.Relation().GetFriendList, filterBlack)
}

func GetFriendListPage(callback open_im_sdk_callback.Base, operationID string, offset int32, count int32, filterBlack bool) {
	call(callback, operationID, IMUserContext.Relation().GetFriendListPage, offset, count, filterBlack)
}

func SearchFriends(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, IMUserContext.Relation().SearchFriends, searchParam)
}

func CheckFriend(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, IMUserContext.Relation().CheckFriend, userIDList)
}

func AddFriend(callback open_im_sdk_callback.Base, operationID string, userIDReqMsg string) {
	call(callback, operationID, IMUserContext.Relation().AddFriend, userIDReqMsg)
}

func UpdateFriends(callback open_im_sdk_callback.Base, operationID string, req string) {
	call(callback, operationID, IMUserContext.Relation().UpdateFriends, req)
}

func DeleteFriend(callback open_im_sdk_callback.Base, operationID string, friendUserID string) {
	call(callback, operationID, IMUserContext.Relation().DeleteFriend, friendUserID)
}

func GetFriendApplicationListAsRecipient(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Relation().GetFriendApplicationListAsRecipient)
}

func GetFriendApplicationListAsApplicant(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Relation().GetFriendApplicationListAsApplicant)
}

func AcceptFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, IMUserContext.Relation().AcceptFriendApplication, userIDHandleMsg)
}

func RefuseFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, IMUserContext.Relation().RefuseFriendApplication, userIDHandleMsg)
}

func AddBlack(callback open_im_sdk_callback.Base, operationID string, blackUserID string, ex string) {
	call(callback, operationID, IMUserContext.Relation().AddBlack, blackUserID, ex)
}

func GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.Relation().GetBlackList)
}

func RemoveBlack(callback open_im_sdk_callback.Base, operationID string, removeUserID string) {
	call(callback, operationID, IMUserContext.Relation().RemoveBlack, removeUserID)
}
