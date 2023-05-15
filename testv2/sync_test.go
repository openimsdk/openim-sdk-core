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

func Test_SyncSelfGroupApplication(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().SyncSelfGroupApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncAdminGroupApplication(t *testing.T) { // todo failed
	err := open_im_sdk.UserForSDK.Group().SyncAdminGroupApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncSelfFriendApplication(t *testing.T) { // todo
	err := open_im_sdk.UserForSDK.Friend().SyncSelfFriendApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncFriendApplication(t *testing.T) { // todo
	err := open_im_sdk.UserForSDK.Friend().SyncFriendApplication(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncFriend(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncFriendList(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncBlack(t *testing.T) {
	err := open_im_sdk.UserForSDK.Friend().SyncBlackList(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SyncConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SyncConversations(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
