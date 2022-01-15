package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/server_api_params"
)

type GetDesignatedFriendsInfoParams []string
type GetDesignatedFriendsInfoCallback []db.LocalFriend

type AddFriendParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	ReqMsg   string `json:"reqMsg"`
}

const AddFriendCallback = constant.SuccessCallbackDefault

//type GetRecvFriendApplicationListParams struct{}
type GetRecvFriendApplicationListCallback []*db.LocalFriendRequest

//type GetSendFriendApplicationListParams struct{}
type GetSendFriendApplicationListCallback []db.LocalFriendRequest

type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

//type ProcessFriendApplicationCallback struct{}
const ProcessFriendApplicationCallback = constant.SuccessCallbackDefault

type CheckFriendParams []string
type CheckFriendCallback []server_api_params.UserIDResult

type DeleteFriendParams string

//type DeleteFriendCallback struct{}
const DeleteFriendCallback = constant.SuccessCallbackDefault

//type GetFriendListParams struct{}
type GetFriendListCallback []db.LocalFriend

type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}

//type SetFriendRemarkCallback struct{}
const SetFriendRemarkCallback = constant.SuccessCallbackDefault

type AddBlackParams string

//type AddBlackCallback struct{}
const AddBlackCallback = constant.SuccessCallbackDefault

//type GetBlackListParams struct{}

//type GetBlackListCallback []LocalBlack
const GetBlackListCallback = constant.SuccessCallbackDefault

type DeleteBlackParams string

//type DeleteBlackCallback struct{}
const DeleteBlackCallback = constant.SuccessCallbackDefault
