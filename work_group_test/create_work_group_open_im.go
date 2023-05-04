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

package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	var groupMemberNumber *int
	groupMemberNumber = flag.Int("gmn", 1000, "group member number ")
	flag.Parse()
	log.NewPrivateLog("", test.LogLevel)
	log.Warn("", "CreateWorkGroup  start, group member number: ", *groupMemberNumber)
	*groupMemberNumber = *groupMemberNumber + 2

	test.CreateWorkGroup(*groupMemberNumber)
	log.Warn("", "CreateWorkGroup finish, group member number: ", *groupMemberNumber+1)

}
