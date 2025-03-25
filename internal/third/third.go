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
	"sync"

	"github.com/openimsdk/openim-sdk-core/v3/internal/third/file"
)

type Third struct {
	platform      int32
	loginUserID   string
	appFramework  string
	LogFilePath   string
	fileUploader  *file.File
	logUploadLock sync.Mutex
}

func (t *Third) SetPlatform(platform int32) {
	t.platform = platform
}

func (t *Third) SetLoginUserID(loginUserID string) {
	t.loginUserID = loginUserID
}

func (t *Third) SetAppFramework(appFramework string) {
	t.appFramework = appFramework
}

func (t *Third) SetLogFilePath(LogFilePath string) {
	t.LogFilePath = LogFilePath
}

func NewThird(fileUploader *file.File) *Third {
	return &Third{fileUploader: fileUploader}
}
