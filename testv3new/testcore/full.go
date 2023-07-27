package testcore

import "open_im_sdk/internal/login"

type FullSDKCore struct {
	login.LoginMgr
}

func NewFullSDKCore(userID string) *FullSDKCore {
	return &FullSDKCore{}
}

func (b *FullSDKCore) SetCallback() {

}

func (b *FullSDKCore) InitConn() error {
	return nil
}

func (b *FullSDKCore) SendMsg() error {
	return nil
}
