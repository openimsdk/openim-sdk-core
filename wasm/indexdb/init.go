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

// 1. Using the native wasm method, TinyGo is applied to Go's embedded domain,
// but the supported functionality is limited and only supports a subset of Go.
// Even JSON serialization is not supported.
// 2. Function names should follow camelCase convention.
// 3. In the provided SQL generation statements, boolean values need special handling.
// For create statements, boolean values should be explicitly handled because SQLite does not natively support boolean types.
// Instead, integers (1 or 0) are used as substitutes, and GORM converts them automatically.
// 4. In the provided SQL generation statements, field names use snake_case (e.g., recv_id),
// but in the interface-converted data, the JSON tag fields follow camelCase (e.g., recvID).
// All such fields should have a mapped transformation.
// 5. Any GORM operations that involve retrieval (e.g., take and find) should specify whether they need to return an error in the documentation.
// 6. For any update operations, be sure to check GORM's implementation. If there is a select(*) query involved,
// you do not need to use the structures in temp_struct.
// 7. Whenever there's a name conflict with an interface, the DB interface should append the "DB" suffix.
// 8. For any map types, use JSON string conversion, and document this clearly.

type IndexDB struct {
	LocalUsers
	LocalConversations
	*LocalChatLogs
	LocalConversationUnreadMessages
	LocalGroups
	LocalGroupMember
	*Black
	*Friend
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
}
