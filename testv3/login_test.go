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

// @Author BanTanger 2023/7/10 12:30:00
package testv3

import (
	"open_im_sdk/testv3/funcation"
	"testing"
	"time"
)

func Test_LoginOne(t *testing.T) {
	uid := "6506148011"
	res := funcation.LoginOne(uid)
	time.Sleep(1000 * time.Second)
	if res != true {
		t.Errorf("uid [%v] login expected be successful, but fail got", uid)
	}
	t.Logf("uid [%v] login successfully", uid)
}

func Test_LoginBatch(t *testing.T) {
	count := 100
	userIDList := funcation.AllUserID
	funcation.LoginBatch(userIDList[:count])
}
