package sdk_params_callback

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/open_im_sdk/server_api_params"
)

type GetDesignatedFriendsInfoParams []string
type GetDesignatedFriendsInfoCallback []open_im_sdk.LocalFriend

type AddFriendParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	ReqMsg   string `json:"reqMsg"`
}

const AddFriendCallback = open_im_sdk.SuccessCallbackDefault

type GetRecvFriendApplicationListParams struct{}
type GetRecvFriendApplicationListCallback []*open_im_sdk.LocalFriendRequest

type GetSendFriendApplicationListParams struct{}
type GetSendFriendApplicationListCallback []open_im_sdk.LocalFriendRequest

type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

//type ProcessFriendApplicationCallback struct{}
const ProcessFriendApplicationCallback = open_im_sdk.SuccessCallbackDefault

type CheckFriendParams []string
type CheckFriendCallback []server_api_params.UserIDResult

type DeleteFriendParams string

//type DeleteFriendCallback struct{}
const DeleteFriendCallback = open_im_sdk.SuccessCallbackDefault

type GetFriendListParams struct{}
type GetFriendListCallback []open_im_sdk.LocalFriend

type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}

//type SetFriendRemarkCallback struct{}
const SetFriendRemarkCallback = open_im_sdk.SuccessCallbackDefault

type AddBlackParams string

//type AddBlackCallback struct{}
const AddBlackCallback = open_im_sdk.SuccessCallbackDefault

type GetBlackListParams struct{}

//type GetBlackListCallback []LocalBlack
const GetBlackListCallback = open_im_sdk.SuccessCallbackDefault

type DeleteBlackParams string

//type DeleteBlackCallback struct{}
const DeleteBlackCallback = open_im_sdk.SuccessCallbackDefault
