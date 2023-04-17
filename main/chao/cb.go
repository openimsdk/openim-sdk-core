package chao

import (
	"fmt"
	"open_im_sdk/open_im_sdk_callback"
)

type Listener struct{}

func (c *Listener) OnConnecting() {
	fmt.Println("OnConnecting")
}

func (c *Listener) OnConnectSuccess() {
	fmt.Println("OnConnectSuccess")
}

func (c *Listener) OnConnectFailed(errCode int32, errMsg string) {
	fmt.Println("OnConnectFailed")
}

func (c *Listener) OnKickedOffline() {
	fmt.Println("OnKickedOffline")
}

func (c *Listener) OnUserTokenExpired() {
	fmt.Println("OnUserTokenExpired")
}

func NewCallback(msg string) open_im_sdk_callback.Base {
	return &Callback{msg: msg}
}

type Callback struct {
	msg string
}

func (c *Callback) OnError(errCode int32, errMsg string) {
	fmt.Printf("[%s] OnError errCode: %d, errMsg: %s\n", c.msg, errCode, errMsg)
}
func (c *Callback) OnSuccess(data string) {
	fmt.Printf("[%s] OnSuccess onSuccess: %s\n", c.msg, data)
}

func NewCbFn(onError func(errCode int32, errMsg string), onSuccess func(data string)) open_im_sdk_callback.Base {
	return &CbFn{onError: onError, onSuccess: onSuccess}
}

type CbFn struct {
	onError   func(errCode int32, errMsg string)
	onSuccess func(data string)
}

func (c *CbFn) OnError(errCode int32, errMsg string) {
	c.onError(errCode, errMsg)
}

func (c *CbFn) OnSuccess(data string) {
	c.onSuccess(data)
}
