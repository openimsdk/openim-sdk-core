// Copyright © 2023 OpenIM SDK.
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

package open_im_sdk_callback

type OnListenerForService interface {
	// 有人申请进群
	OnGroupApplicationAdded(groupApplication string)
	// 进群申请被同意
	OnGroupApplicationAccepted(groupApplication string)
	// 有人申请添加你为好友
	OnFriendApplicationAdded(friendApplication string)
	// 好友申请被同意
	OnFriendApplicationAccepted(groupApplication string)
	// 收到新消息
	OnRecvNewMessage(message string)
}
