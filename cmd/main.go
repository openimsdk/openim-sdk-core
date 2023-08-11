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
	"fmt"
	"open_im_sdk/test"
	"time"
)

func main() {
	APIADDR := "http://59.36.173.89:10002"
	WSADDR := "ws://59.36.173.89:10001"
	REGISTERADDR := APIADDR + "/user_register"
	ACCOUNTCHECK := APIADDR + "/manager/account_check"
	TOKENADDR := APIADDR + "/auth/user_token"
	SECRET := "openIM123"
	SENDINTERVAL := 20
	test.REGISTERADDR = REGISTERADDR
	test.TOKENADDR = TOKENADDR
	test.SECRET = SECRET
	test.SENDINTERVAL = SENDINTERVAL
	test.WSADDR = WSADDR
	test.ACCOUNTCHECK = ACCOUNTCHECK
	strMyUidx := "9226250128"

	tokenx := test.RunGetToken(strMyUidx)
	fmt.Println(tokenx)
	test.InOutDoTest(strMyUidx, tokenx, WSADDR, APIADDR)
	time.Sleep(time.Second * 10)
	// test.DoTestGetUsersInfo()
	// test.DoTestSetMsgDestructTime("sg_1012596513")
	// test.DoTestRevoke()
	// test.DotestDeleteFriend("8303492153")
	// test.TestMarkGroupMessageAsRead()
	// test.DoTestRevoke()
	// time.Sleep(time.Second * 5)
	// test.DoTestAddToBlackList("9169012630")
	// test.DoTestDeleteFromBlackList("9169012630")
	// test.DotestDeleteFriend("9169012630")
	// test.DoTestSetConversationPinned("si_2456093263_9169012630", true)
	// test.DoTestSetOneConversationRecvMessageOpt("si_2456093263_9169012630", 2)
	// test.DoTestGetConversationRecvMessageOpt("si_2456093263_9169012630")
	// test.DoTestDeleteConversationMsgFromLocalAndSvr("sg_537415520")
	for {
		time.Sleep(10000 * time.Millisecond)
	}

}
