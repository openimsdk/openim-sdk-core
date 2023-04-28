// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

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
