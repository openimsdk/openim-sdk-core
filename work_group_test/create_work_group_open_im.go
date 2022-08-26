package main

import (
	"flag"
	"open_im_sdk/pkg/log"
	"open_im_sdk/test"
)

func main() {
	var onlineNum *int
	onlineNum = flag.Int("gmn", 10, "group member number ")
	flag.Parse()
	log.Warn("", "CreateWorkGroup  start, group member number: ", *onlineNum)
	*onlineNum = *onlineNum + 2
	groupID := test.CreateWorkGroup(*onlineNum)
	log.Warn("", "CreateWorkGroup finish, group member number: ", *onlineNum, groupID)
}
