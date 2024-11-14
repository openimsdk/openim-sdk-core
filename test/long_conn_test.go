package test

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"
	"testing"
	"time"
)

func Test_SubscribeUsersStatus(t *testing.T) {
	time.Sleep(time.Second)
	message, err := open_im_sdk.UserForSDK.LongConnMgr().SubscribeUsersStatus(ctx, []string{"5975996883"})
	if err != nil {
		t.Error(err)
	}
	t.Log(message)
	ch := make(chan struct{})
	select {
	case <-ch:
	}
}
