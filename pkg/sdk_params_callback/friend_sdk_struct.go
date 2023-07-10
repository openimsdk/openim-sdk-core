package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
)

//1
type GetDesignatedFriendsInfoParams []string
type GetDesignatedFriendsInfoCallback []server_api_params.FullUserInfo

//1
type AddFriendParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	ReqMsg   string `json:"reqMsg"`
}

const AddFriendCallback = constant.SuccessCallbackDefault

//1
//type GetRecvFriendApplicationListParams struct{}
type GetRecvFriendApplicationListCallback []*model_struct.LocalFriendRequest

//1
//type GetSendFriendApplicationListParams struct{}
type GetSendFriendApplicationListCallback []*model_struct.LocalFriendRequest

//1
type ProcessFriendApplicationParams struct {
	ToUserID  string `json:"toUserID" validate:"required"`
	HandleMsg string `json:"handleMsg"`
}

const ProcessFriendApplicationCallback = constant.SuccessCallbackDefault

//1
type CheckFriendParams []string
type CheckFriendCallback []server_api_params.UserIDResult

//1
type DeleteFriendParams string

//type DeleteFriendCallback struct{}
const DeleteFriendCallback = constant.SuccessCallbackDefault

//1
//type GetFriendListParams struct{}
type GetFriendListCallback []server_api_params.FullUserInfo

type SearchFriendsParam struct {
	KeywordList      []string `json:"keywordList"`
	IsSearchUserID   bool     `json:"isSearchUserID"`
	IsSearchNickname bool     `json:"isSearchNickname"`
	IsSearchRemark   bool     `json:"isSearchRemark"`
}
type SearchFriendsCallback []*SearchFriendItem
type SearchFriendItem struct {
	model_struct.LocalFriend
	Relationship int `json:"relationship"`
}

//1
type SetFriendRemarkParams struct {
	ToUserID string `json:"toUserID" validate:"required"`
	Remark   string `json:"remark" validate:"required"`
}

//type SetFriendRemarkCallback struct{}
const SetFriendRemarkCallback = constant.SuccessCallbackDefault

//1
type AddBlackParams string

//type AddBlackCallback struct{}
const AddBlackCallback = constant.SuccessCallbackDefault

//1
//type GetBlackListParams struct{}
type GetBlackListCallback []server_api_params.FullUserInfo

//1
type RemoveBlackParams string

//type DeleteBlackCallback struct{}
const RemoveBlackCallback = constant.SuccessCallbackDefault
