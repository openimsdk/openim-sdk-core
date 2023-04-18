package testv2

type Listener struct{}

func (c *Listener) OnConnecting() {
	// fmt.Println("OnConnecting")
}

func (c *Listener) OnConnectSuccess() {
	// fmt.Println("OnConnectSuccess")
}

func (c *Listener) OnConnectFailed(errCode int32, errMsg string) {
	// fmt.Println("OnConnectFailed")
}

func (c *Listener) OnKickedOffline() {
	// fmt.Println("OnKickedOffline")
}

func (c *Listener) OnUserTokenExpired() {
	// fmt.Println("OnUserTokenExpired")
}
