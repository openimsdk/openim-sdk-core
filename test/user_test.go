package test

import (
	"testing"
	"time"

	"github.com/openimsdk/protocol/wrapperspb"

	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk"

	"github.com/openimsdk/protocol/sdkws"
)

func Test_GetSelfUserInfo(t *testing.T) {
	userInfo, err := open_im_sdk.IMUserContext.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}

	t.Log(userInfo)
}

func Test_SetSelfInfoEx(t *testing.T) {
	newNickName := "test"
	//newFaceURL := "http://test.com"
	err := open_im_sdk.IMUserContext.User().SetSelfInfo(ctx, &sdkws.UserInfoWithEx{
		Nickname: &wrapperspb.StringValue{
			Value: newNickName,
		},
		//FaceURL:  newFaceURL,
		Ex: &wrapperspb.StringValue{
			Value: "ASD",
		},
	})
	newFaceURL := "http://test.com"

	if err != nil {
		t.Error(err)
	}
	userInfo, err := open_im_sdk.IMUserContext.User().GetSelfUserInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	if userInfo.UserID != UserID && userInfo.Nickname != newNickName && userInfo.FaceURL != newFaceURL {
		t.Error("user id not match")
	}
	t.Log(userInfo)
	time.Sleep(time.Second * 10)
}
