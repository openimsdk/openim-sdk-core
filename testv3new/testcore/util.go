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

// @Author BanTanger 2023/7/21 15:04
package testcore

import (
	"context"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/utils"
)

func InitContext(uid string) (context.Context, ccontext.GlobalConfig) {
	config := ccontext.GlobalConfig{
		UserID:   uid,
		Token:    AdminToken,
		IMConfig: Config,
	}
	ctx := ccontext.WithInfo(context.Background(), &config)
	ctx = ccontext.WithOperationID(ctx, utils.OperationIDGenerator())
	return ctx, config
}

func AddUserID(uid string) {
	userLock.Lock()
	AllUserID = append(AllUserID, uid)
	userLock.Unlock()
}
