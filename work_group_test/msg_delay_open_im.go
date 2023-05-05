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
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond
	var groupID *string

	//senderNum = flag.Int("sn", 10, "sender num")
	singleSenderMsgNum = flag.Int("mn", 100, "single sender msg num")
	intervalTime = flag.Int("t", 100, "interval time mill second")
	groupID = flag.String("gid", "3282359177", "groupID")
	//	pressClientNum = flag.Int("pcn", 8, "press client number ")
	flag.Parse()
	//test.InitMgr(*senderNum)
	log.NewPrivateLog(test.LogName, test.LogLevel)
	n := test.GetGroupMemberNum(*groupID)
	var pressClientNum int
	pressClientNum = int(n) - 3

	//	log.Warn("", "reliability test start, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	test.WorkGroupMsgDelayTest(*singleSenderMsgNum, *intervalTime, 10, pressClientNum, pressClientNum+1, *groupID)
}
