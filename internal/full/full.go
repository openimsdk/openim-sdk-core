package full

import (
	"open_im_sdk/internal/friend"
	"open_im_sdk/internal/group"
	"open_im_sdk/internal/user"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	sdk "open_im_sdk/pkg/sdk_params_callback"
)

type Full struct {
	user   *user.User
	friend *friend.Friend
	group  *group.Group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group) *Full {
	return &Full{user: user, friend: friend, group: group}
}

func (u *Full) getUsersInfo(callback open_im_sdk_callback.Base, userIDList sdk.GetUsersInfoParam, operationID string) sdk.GetUsersInfoCallback {
	//from svr
	publicList := u.user.GetUsersInfoFromSvr(callback, userIDList, operationID)
	friendList := u.friend.GetDesignatedFriendListInfo(callback, []string(userIDList), operationID)
	blackList := u.friend.GetDesignatedBlackListInfo(callback, []string(userIDList), operationID)
	return common.MergeUserResult(publicList, friendList, blackList)
}
