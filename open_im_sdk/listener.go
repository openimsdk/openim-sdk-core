package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func SetGroupListener(listener open_im_sdk_callback.OnGroupListener) {
	listenerCall(IMUserContext.SetGroupListener, listener)
}

func SetConversationListener(listener open_im_sdk_callback.OnConversationListener) {
	listenerCall(IMUserContext.SetConversationListener, listener)
}

func SetAdvancedMsgListener(listener open_im_sdk_callback.OnAdvancedMsgListener) {
	listenerCall(IMUserContext.SetAdvancedMsgListener, listener)
}

func SetUserListener(listener open_im_sdk_callback.OnUserListener) {
	listenerCall(IMUserContext.SetUserListener, listener)

}

func SetFriendListener(listener open_im_sdk_callback.OnFriendshipListener) {
	listenerCall(IMUserContext.SetFriendshipListener, listener)
}

func SetCustomBusinessListener(listener open_im_sdk_callback.OnCustomBusinessListener) {
	listenerCall(IMUserContext.SetCustomBusinessListener, listener)
}

func SetMessageKvInfoListener(listener open_im_sdk_callback.OnMessageKvInfoListener) {
	listenerCall(IMUserContext.SetMessageKvInfoListener, listener)
}
