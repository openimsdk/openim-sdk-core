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
	log.Warn("", "CreateWorkGroup  start, group member number: ", *groupMemberNumber)
	*groupMemberNumber = *groupMemberNumber + 2
	groupID := test.CreateWorkGroup(*groupMemberNumber)
	log.Warn("", "CreateWorkGroup finish, group member number: ", *groupMemberNumber, groupID)
}
