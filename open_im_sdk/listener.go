// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
