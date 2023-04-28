// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package testv2

type OnConnListener struct{}

func (c *OnConnListener) OnConnecting() {
	// fmt.Println("OnConnecting")
}

func (c *OnConnListener) OnConnectSuccess() {
	// fmt.Println("OnConnectSuccess")
}

func (c *OnConnListener) OnConnectFailed(errCode int32, errMsg string) {
	// fmt.Println("OnConnectFailed")
}

func (c *OnConnListener) OnKickedOffline() {
	// fmt.Println("OnKickedOffline")
}

func (c *OnConnListener) OnUserTokenExpired() {
	// fmt.Println("OnUserTokenExpired")
}
