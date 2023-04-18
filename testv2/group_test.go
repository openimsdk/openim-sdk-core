package testv2

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/group"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"open_im_sdk/open_im_sdk"
	"testing"
)

func Test_CreateGroupV2(t *testing.T) {
	req := &group.CreateGroupReq{
		InitMembers:  []string{},
		AdminUserIDs: []string{},
		OwnerUserID:  UserID,
		GroupInfo: &sdkws.GroupInfo{
			GroupName: "test-group",
		},
	}
	info, err := open_im_sdk.UserForSDK.Group().CreateGroupV2(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("group info: %s", info.String())
}

func Test_JoinGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().JoinGroup(ctx, "xadxwr24", "1234", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("JoinGroup success")
}

func Test_QuitGroup(t *testing.T) {
	err := open_im_sdk.UserForSDK.Group().QuitGroup(ctx, "xadxwr24")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("QuitGroup success")
}

func Test_GetJoinedGroupList(t *testing.T) {
	info, err := open_im_sdk.UserForSDK.Group().GetJoinedGroupList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("GetJoinedGroupList: %d\n", len(info))
	for _, localGroup := range info {
		t.Logf("%#v", localGroup)
	}
}
