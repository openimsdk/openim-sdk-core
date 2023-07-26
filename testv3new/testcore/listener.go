package testcore

type ConnListner struct {
}

func (c *ConnListner) OnConnecting()     {}
func (c *ConnListner) OnConnectSuccess() {}
func (c *ConnListner) OnConnectFailed(errCode int32, errMsg string) {
	// log.ZError(context.Background(), "connect failed", nil, "errCode", errCode, "errMsg", errMsg)
}
func (c *ConnListner) OnKickedOffline()    {}
func (c *ConnListner) OnUserTokenExpired() {}
