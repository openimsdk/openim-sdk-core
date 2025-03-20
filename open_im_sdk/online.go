package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

// SubscribeUsersStatus Presence status of subscribed users.
func SubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, IMUserContext.LongConnMgr().SubscribeUsersStatus, userIDs)
}

// UnsubscribeUsersStatus Unsubscribe a user's presence.
func UnsubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, IMUserContext.LongConnMgr().UnsubscribeUsersStatus, userIDs)
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func GetSubscribeUsersStatus(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.LongConnMgr().GetSubscribeUsersStatus)
}

// GetUserStatus Get the online status of users.
func GetUserStatus(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, IMUserContext.LongConnMgr().SubscribeUsersStatus, userIDs)
}
