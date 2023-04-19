package testv2

import (
	"open_im_sdk/open_im_sdk"
	"testing"
)

func Test_GetWorkMomentsUnReadCount(t *testing.T) {
	unreadCount, err := open_im_sdk.UserForSDK.WorkMoments().GetWorkMomentsUnReadCount(ctx)
	if err != nil {
		t.Error(err)
	}
	t.Log(unreadCount)
}

func Test_GetWorkMomentsNotification(t *testing.T) {
	notifications, err := open_im_sdk.UserForSDK.WorkMoments().GetWorkMomentsNotification(ctx, 0, 10)
	if err != nil {
		t.Error(err)
	}
	t.Log(notifications)
}

func Test_ClearWorkMomentsNotification(t *testing.T) {
	err := open_im_sdk.UserForSDK.WorkMoments().ClearWorkMomentsNotification(ctx)
	if err != nil {
		t.Error(err)
	}
}
