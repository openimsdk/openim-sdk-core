package open_im_sdk

import "open_im_sdk/open_im_sdk_callback"

func GetDesignatedFriendsInfo(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, userForSDK.Friend().GetDesignatedFriendsInfo, userIDList)
}

func GetFriendList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Friend().GetFriendList)
}

func SearchFriends(callback open_im_sdk_callback.Base, operationID string, searchParam string) {
	call(callback, operationID, userForSDK.Friend().SearchFriends, searchParam)
}

func CheckFriend(callback open_im_sdk_callback.Base, operationID string, userIDList string) {
	call(callback, operationID, userForSDK.Friend().CheckFriend, userIDList)
}

func AddFriend(callback open_im_sdk_callback.Base, operationID string, userIDReqMsg string) {
	call(callback, operationID, userForSDK.Friend().AddFriend, userIDReqMsg)
}

func SetFriendRemark(callback open_im_sdk_callback.Base, operationID string, userIDRemark string) {
	call(callback, operationID, userForSDK.Friend().SetFriendRemark, userIDRemark)
}

func DeleteFriend(callback open_im_sdk_callback.Base, operationID string, friendUserID string) {
	call(callback, operationID, userForSDK.Friend().DeleteFriend, friendUserID)
}

func GetRecvFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Friend().GetRecvFriendApplicationList)
}

func GetSendFriendApplicationList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Friend().GetSendFriendApplicationList)
}

func AcceptFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, userForSDK.Friend().AcceptFriendApplication, userIDHandleMsg)
}

func RefuseFriendApplication(callback open_im_sdk_callback.Base, operationID string, userIDHandleMsg string) {
	call(callback, operationID, userForSDK.Friend().RefuseFriendApplication, userIDHandleMsg)
}

func AddBlack(callback open_im_sdk_callback.Base, operationID string, blackUserID string) {
	call(callback, operationID, userForSDK.Friend().AddBlack, blackUserID)
}

func GetBlackList(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, userForSDK.Friend().GetBlackList)
}

func RemoveBlack(callback open_im_sdk_callback.Base, operationID string, removeUserID string) {
	call(callback, operationID, userForSDK.Friend().RemoveBlack, removeUserID)
}
