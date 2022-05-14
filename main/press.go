package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	var onlineNum *int          //Number of online users
	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond
	onlineNum = flag.Int("on", 10, "online num")
	senderNum = flag.Int("sn", 100, "sender num")
	singleSenderMsgNum = flag.Int("mn", 100, "single sender msg num")
	intervalTime = flag.Int("t", 10, "interval time mill second")
	flag.Parse()

	log.NewPrivateLog("open_im_test", 4)
	log.Warn("", "online test start, number of online users: ", *onlineNum)
	test.OnlineTest(*onlineNum)
	log.Warn("", "online test finish, number of online users: ", *onlineNum)
	log.Warn("", "reliability test start, user: ", *senderNum, "message number: ", *singleSenderMsgNum)
	test.ReliabilityTest(*singleSenderMsgNum, *intervalTime, 10, *senderNum)
}
