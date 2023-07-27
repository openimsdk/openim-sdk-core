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

// @Author BanTanger 2023/7/10 15:30:00
package funcation

import (
	"context"

	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"

	"open_im_sdk/sdk_struct"
)

func SendMsg(ctx context.Context, fromID, recvID, groupID, msg string) (*sdk_struct.MsgStruct, error) {
	message, err := AllLoginMgr[fromID].Mgr.Conversation().CreateTextMessage(ctx, msg)
	if err != nil {
		log.ZError(ctx, "CreateTextMessage ", err)
		return nil, err
	}
	o := sdkws.OfflinePushInfo{Title: "title", Desc: "desc"}
	log.ZInfo(ctx, "SendMessage ", "fromID", fromID, "recvID", recvID, "groupID", groupID, "message", message)
	AllLoginMgr[fromID].Mgr.Conversation().SendMessage(ctx, message, recvID, groupID, &o)
	return message, nil
}
