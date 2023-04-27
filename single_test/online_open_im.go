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
	var onlineNum *int //Number of online users
	onlineNum = flag.Int("on", 10, "online num")
	flag.Parse()
	log.Warn("", "online test start, online num: ", *onlineNum)
	test.OnlineTest(*onlineNum)
	log.Warn("", "online test finish, online num: ", *onlineNum)
	select {}
}
