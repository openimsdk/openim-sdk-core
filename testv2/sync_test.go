package testv2

import (
	"open_im_sdk/open_im_sdk"
	"testing"
)

func Test_SyncJoinGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SyncJoinedGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncGroupMember(t *testing.T) {
	groups, err := open_im_sdk.UserForSDK.Group().GetServerJoinGroup(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range groups {
		err := open_im_sdk.UserForSDK.Group().SyncGroupMember(ctx, group.GroupID)
		if err != nil {
			t.Fatal(err)
		}
	}
}
