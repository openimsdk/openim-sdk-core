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

// @Author BanTanger 2023/7/23 15:02
package testv3new

import (
	"fmt"
	"testing"
)

func Test_userRegister(t *testing.T) {
	userID := "bantanger123"
	manager := NewRegisterManager()
	_ = manager.RegisterOne(userID)
	token, _ := manager.GetToken(userID)
	t.Log(token)
}

func Test_userRegisterBatch(t *testing.T) {
	count := 1001
	var users []string
	for i := 0; i <= count; i++ {
		users = append(users, fmt.Sprintf("register_test_%v", i))
	}
	manager := NewRegisterManager()
	_ = manager.RegisterBatch(users)
	token, _ := manager.GetToken(users[0])
	t.Log(token)
}
