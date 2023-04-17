package testv2

import (
	"fmt"
	"open_im_sdk/open_im_sdk"
	"testing"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
)

func Test_GetSelfUserInfo(t *testing.T) {
	fmt.Println("xxxx")
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(userInfo, "ss")
	log.ZDebug(ctx, "Test_GetSelfUserInfo", userInfo)
}
