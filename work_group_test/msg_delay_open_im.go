package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	//	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond
	var groupID *string
	var pressClientNum *int
	//senderNum = flag.Int("sn", 10, "sender num")
	singleSenderMsgNum = flag.Int("mn", 3, "single sender msg num")
	intervalTime = flag.Int("t", 1000, "interval time mill second")
	groupID = flag.String("gID", "3446203278", "groupID")
	pressClientNum = flag.Int("pcn", 10, "press client number ")
	flag.Parse()
	//test.InitMgr(*senderNum)
	log.NewPrivateLog(test.LogName, test.LogLevel)
	//	log.Warn("", "reliability test start, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	test.WorkGroupMsgDelayTest(*singleSenderMsgNum, *intervalTime, 10, *pressClientNum, *pressClientNum+1, *groupID)
}
