package module

import (
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"

	"github.com/openimsdk/protocol/friend"
)

type TestFriendManager struct {
	*MetaManager
}

func (t *TestFriendManager) ImportFriends(ownerUserID string, friendUserIDs []string) error {
	req := &friend.ImportFriendReq{
		OwnerUserID:   ownerUserID,
		FriendUserIDs: friendUserIDs,
	}
	return t.postWithCtx(constant.ImportFriendListRouter, &req, nil)
}
