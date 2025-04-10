package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func GetUsersInfo(callback open_im_sdk_callback.Base, operationID string, userIDs string) {
	call(callback, operationID, IMUserContext.User().GetUsersInfo, userIDs)
}

// SetSelfInfo sets the user's own information.
func SetSelfInfo(callback open_im_sdk_callback.Base, operationID string, userInfo string) {
	call(callback, operationID, IMUserContext.User().SetSelfInfo, userInfo)
}

// GetSelfUserInfo obtains the user's own information.
func GetSelfUserInfo(callback open_im_sdk_callback.Base, operationID string) {
	call(callback, operationID, IMUserContext.User().GetSelfUserInfo)
}
