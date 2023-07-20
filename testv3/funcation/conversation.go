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

// @Author BanTanger 2023/7/11 10:27
package funcation

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func GetConversation(uid, conversationID, startMsgID string, count int) {
	var params sdk_params_callback.GetAdvancedHistoryMessageListParams
	params.UserID = AllLoginMgr[uid].UserID
	params.ConversationID = conversationID
	params.StartClientMsgID = startMsgID
	params.Count = count
	open_im_sdk.GetAdvancedHistoryMessageList(&testConversation, utils.OperationIDGenerator(), utils.StructToJsonString(params))
}
