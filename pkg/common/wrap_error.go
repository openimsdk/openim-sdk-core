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

package common

//import (
//	"github.com/mitchellh/mapstructure"
//	"open_im_sdk/open_im_sdk_callback"
//	"open_im_sdk/pkg/db"
//	"open_im_sdk/pkg/db/model_struct"
//)
//
//funcation GetGroupMemberListByGroupID(callback open_im_sdk_callback.Base, operationID string, db *db.DataBase, groupID string) []*model_struct.LocalGroupMember {
//	memberList, err := db.GetGroupMemberListByGroupID(groupID)
//	CheckDBErrCallback(callback, err, operationID)
//	return memberList
//}
//
//funcation MapstructureDecode(input interface{}, output interface{}, callback open_im_sdk_callback.Base, oprationID string) {
//	err := mapstructure.Decode(input, output)
//	CheckDataErrCallback(callback, err, oprationID)
//}
