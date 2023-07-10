package main

import (
	"flag"
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond
	var groupID *string
	senderNum = flag.Int("sn", 10, "sender num")
	singleSenderMsgNum = flag.Int("mn", 10, "single sender msg num")
	intervalTime = flag.Int("t", 0, "interval time mill second")
	groupID = flag.String("gid", "1607724699", "groupID")
	flag.Parse()
	constant.OnlyForTest = 1
	log.NewPrivateLog("", test.LogLevel)
	n := test.GetGroupMemberNum(*groupID)
	if n-3 < uint32(*senderNum) {
		log.Error("", "args failed n-3 < *senderNum ", n-3, *senderNum)
	}
	log.Warn("", "WorkGroupPressTest begin, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum, "groupID: ", *groupID)
	test.WorkGroupPressTest(*singleSenderMsgNum, *intervalTime, *senderNum, *groupID)
	log.Warn("", "WorkGroupPressTest end")
	select {}
}
