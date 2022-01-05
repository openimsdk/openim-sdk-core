package open_im_sdk

import (
	"open_im_sdk/open_im_sdk/base_info"
)

type GetDesignatedFriendsInfoParams []string
type GetDesignatedFriendsInfoCallback []LocalFriend

type AddFriendParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	ReqMsg   string `json:"reqMsg"`
}
type AddFriendCallback struct{}

type GetRecvFriendApplicationListParams struct{}
type GetRecvFriendApplicationListCallback []LocalFriendRequest

type GetSendFriendApplicationListParams struct{}
type GetSendFriendApplicationListCallback []LocalFriendRequest

type ProcessFriendApplicationParams struct {
	ToUserID     string `json:"toUserID" validate:"required"`
	HandleMsg    string `json:"handleMsg"`
	HandleResult int32  `json:"handleResult" validate:oneof=-1 1`
}
type ProcessFriendApplicationCallback struct{}

type CheckFriendParams []string
type CheckFriendCallback []base_info.UserIDResult

type DeleteFriendParams string
type DeleteFriendCallback struct{}

type GetFriendListParams struct{}
type GetFriendListCallback []LocalFriend

type SetFriendRemarkParams string
type SetFriendRemarkCallback struct{}

type AddBlackParams string
type AddBlackCallback struct{}

type GetBlackListParams struct{}
type GetBlackListCallback []LocalBlack

type DeleteBlackParams string
type DeleteBlackCallback struct{}
