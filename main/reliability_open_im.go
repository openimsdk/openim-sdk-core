package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond

	senderNum = flag.Int("sn", 10, "sender num")
	singleSenderMsgNum = flag.Int("mn", 1000, "single sender msg num")
	intervalTime = flag.Int("t", 10, "interval time mill second")
	flag.Parse()

	log.NewPrivateLog("", uint32(test.LogLevel))
	log.Warn("", "reliability test start, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	test.ReliabilityTest(*singleSenderMsgNum, *intervalTime, 10, *senderNum)
}
