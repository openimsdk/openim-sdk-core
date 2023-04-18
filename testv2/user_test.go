package testv2

import (
	"open_im_sdk/open_im_sdk"
	"testing"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func Test_GetSelfUserInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(userInfo)
}

func Test_GetUsersInfo(t *testing.T) {
	userInfo, err := open_im_sdk.UserForSDK.User().GetUsersInfo(ctx, []string{UserID})
	if err != nil {
		t.Error(err)
	}
	t.Log(userInfo)
}

func Test_SetSelfInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.User().SetSelfInfo(ctx, &sdkws.UserInfo{
		Nickname: "test",
		FaceURL:  "http://test.com",
	})
	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	if userInfo.UserID != UserID {
		t.Error("user id not match")
	}
	t.Log(userInfo)
}

func Test_UpdateMsgSenderInfo(t *testing.T) {
	err := open_im_sdk.UserForSDK.User().UpdateMsgSenderInfo(ctx, "test", "http://test.com")
	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.UserForSDK.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(userInfo)
}
