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
	senderNum = flag.Int("sn", 2, "sender num")
	singleSenderMsgNum = flag.Int("mn", 10000, "single sender msg num")
	intervalTime = flag.Int("t", 10, "interval time mill second")
	flag.Parse()
	log.NewPrivateLog("", uint32(test.LogLevel))
	log.Warn("", "press test start, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	test.PressTest(*singleSenderMsgNum, *intervalTime, *senderNum)
	log.Warn("", "press test finish, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	select {}
}
