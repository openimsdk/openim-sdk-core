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

//go:build js && wasm
// +build js,wasm

package indexdb

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/wasm/exec"
)

//1.Use native WASM methods, tinyGo is applied in the Go embedded domain with limited supported features, offering only
//  a subset of the Go language. It doesn't even support JSON serialization.
//2.Function names follow camelCase naming convention.
//3.In the provided SQL generation statements, boolean values require special handling. For CREATE statements, any field
//  designed to hold a boolean value should be explicitly noted. This is because SQLite natively does not support
//  boolean types and uses integers 1 or 0 as substitutes. GORM handles this by converting booleans into integers.
//4.In the provided SQL generation statements, the field names use snake_case, for example: recv_id. However, in the
//  interface, the JSON tag field style follows camelCase, like recvID. You need to create a mapping between all similar
//  fields to properly handle the conversion between the two formats.
//5.Any operation involving GORM, such as Take and Find, must clearly indicate in the documentation whether it returns
//  an error. This includes explaining when and how these functions might return errors, ensuring proper error handling
//  is accounted for in the implementation.
//6.For any operations involving Update, it's essential to review the GORM prototype implementation. If select(*) is
//  present, you don't need to use the structure from temp_struct. Make sure to handle updates accordingly based on the
//  GORM behavior for optimal performance and accuracy.
//7.Whenever a name conflicts with an interface, append the suffix DB to the database-related interface names for
//  clarity and to avoid conflicts. This ensures consistency and avoids naming collisions in the codebase.
//8.For any map types, always use JSON string conversion. This should be clearly documented to ensure consistency in
//  handling map data types across the project.

type IndexDB struct {
	LocalUsers
	LocalConversations
	*LocalChatLogs
	LocalConversationUnreadMessages
	LocalGroups
	LocalGroupMember
	LocalGroupRequest
	LocalCacheMessage
	LocalUserCommand
	*FriendRequest
	*Black
	*Friend
	LocalChatLogReactionExtensions
	loginUserID string
}

func (i IndexDB) Close(ctx context.Context) error {
	_, err := exec.Exec()
	return err
}

func (i IndexDB) InitDB(ctx context.Context, userID string, dataDir string) error {
	_, err := exec.Exec(userID, dataDir)
	return err
}

func (i IndexDB) SetChatLogFailedStatus(ctx context.Context) {
	//msgList, err := i.GetSendingMessageList()
	//if err != nil {
	//	log.Error("", "GetSendingMessageList failed ", err.Error())
	//	return
	//}
	//for _, v := range msgList {
	//	v.Status = constant.MsgStatusSendFailed
	//	err := i.UpdateMessage(v)
	//	if err != nil {
	//		log.Error("", "UpdateMessage failed ", err.Error(), v)
	//		continue
	//	}
	//}
	//groupIDList, err := i.GetReadDiffusionGroupIDList()
	//if err != nil {
	//	log.Error("", "GetReadDiffusionGroupIDList failed ", err.Error())
	//	return
	//}
	//for _, v := range groupIDList {
	//	msgList, err := i.SuperGroupGetSendingMessageList(v)
	//	if err != nil {
	//		log.Error("", "GetSendingMessageList failed ", err.Error())
	//		return
	//	}
	//	if len(msgList) > 0 {
	//		for _, v := range msgList {
	//			v.Status = constant.MsgStatusSendFailed
	//			err := i.SuperGroupUpdateMessage(v)
	//			if err != nil {
	//				log.Error("", "UpdateMessage failed ", err.Error(), v)
	//				continue
	//			}
	//		}
	//	}
	//
	//}
}
