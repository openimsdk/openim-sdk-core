package common

import "open_im_sdk/pkg/db"

type Base interface {
	OnError(errCode int32, errMsg string)
	OnSuccess(data string)
}

type FriendApplicationListAddedCallback db.LocalFriendRequest
type FriendApplicationListAcceptCallback db.LocalFriendRequest
type FriendApplicationListRejectCallback db.LocalFriendRequest
type FriendListAddedCallback db.LocalFriend
type FriendListDeletedCallback db.LocalFriend
type BlackListAddCallback db.LocalBlack
type BlackListDeletedCallback db.LocalBlack
type FriendInfoChangedCallback db.LocalFriend
