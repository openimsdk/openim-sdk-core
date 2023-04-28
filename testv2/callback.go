// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package testv2

import "fmt"

type TestSendMsg struct {
}

func (TestSendMsg) OnSuccess(data string) {
	fmt.Println("testSendImg, OnSuccess, output: ", data)
}

func (TestSendMsg) OnError(code int32, msg string) {
	fmt.Println("testSendImg, OnError, ", code, msg)
}

func (TestSendMsg) OnProgress(progress int) {
	fmt.Println("progress: ", progress)
}
