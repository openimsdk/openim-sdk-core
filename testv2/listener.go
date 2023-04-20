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
