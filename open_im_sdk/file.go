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
	"open_im_sdk/internal/file"
	"open_im_sdk/open_im_sdk_callback"
)

func UploadFile(callback open_im_sdk_callback.Base, operationID string, req string, progress open_im_sdk_callback.UploadFileCallback) {
	call(callback, operationID, UserForSDK.File().UploadFile, req, file.UploadFileCallback(progress))
}
