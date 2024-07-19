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

package third

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/internal/file"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/protocol/third"
	"sync"
)

type Third struct {
	platformID    int32
	loginUserID   string
	version       string
	systemType    string
	LogFilePath   string
	fileUploader  *file.File
	logUploadLock sync.Mutex
}

func NewThird(platformID int32, loginUserID, version, systemType, LogFilePath string, fileUploader *file.File) *Third {
	return &Third{platformID: platformID, loginUserID: loginUserID, version: version, systemType: systemType, LogFilePath: LogFilePath, fileUploader: fileUploader}
}

func (c *Third) UpdateFcmToken(ctx context.Context, fcmToken string, expireTime int64) error {
	req := third.FcmUpdateTokenReq{
		PlatformID: c.platformID,
		FcmToken:   fcmToken,
		Account:    c.loginUserID,
		ExpireTime: expireTime}
	_, err := util.CallApi[third.FcmUpdateTokenResp](ctx, constant.FcmUpdateTokenRouter, &req)
	return err

}

func (c *Third) SetAppBadge(ctx context.Context, appUnreadCount int32) error {
	req := third.SetAppBadgeReq{
		UserID:         c.loginUserID,
		AppUnreadCount: appUnreadCount,
	}
	_, err := util.CallApi[third.SetAppBadgeResp](ctx, constant.SetAppBadgeRouter, &req)
	return err
}
