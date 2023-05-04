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

package open_im_sdk

import (
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/log"
)

func SetSignalingListener(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetSignalingListener(callback)
}

func SetSignalingListenerForService(callback open_im_sdk_callback.OnSignalingListener) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetSignalingListenerForService(callback)
}

func SetListenerForService(callback open_im_sdk_callback.OnListenerForService) {
	if callback == nil || UserForSDK == nil {
		log.Error("callback or UserForSDK is nil")
		return
	}
	UserForSDK.SetListenerForService(callback)
}
