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

package test

import (
	"fmt"
)

type WBase struct {
}

func (WBase) OnError(errCode int32, errMsg string) {
	fmt.Println("get workmoments OnError", errCode, errMsg)
}
func (WBase) OnSuccess(data string) {
	fmt.Println("get workmoments OnSuccess, ", data)
}

func (WBase) OnProgress(progress int) {
	fmt.Println("OnProgress, ", progress)
}

//funcation TestGetWorkMomentsUnReadCount() {
//	operationID := utils.OperationIDGenerator()
//	var cb WBase
//	open_im_sdk.GetWorkMomentsUnReadCount(cb, operationID)
//}
//
//funcation TestGetWorkMomentsNotification() {
//	operationID := utils.OperationIDGenerator()
//	var cb WBase
//	offset := 0
//	count := 10
//	open_im_sdk.GetWorkMomentsNotification(cb, operationID, offset, count)
//}
//
//funcation TestClearWorkMomentsNotification() {
//	operationID := utils.OperationIDGenerator()
//	var cb WBase
//	open_im_sdk.ClearWorkMomentsNotification(cb, operationID)
//}
