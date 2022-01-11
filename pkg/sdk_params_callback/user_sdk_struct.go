package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/pkg/server_api_params"
)

//other user
type GetUsersInfoParam []string
type GetUsersInfoCallback []*server_api_params.PublicUserInfo


//type GetSelfUserInfoParam string
type GetSelfUserInfoCallback *db.LocalUser


type SetSelfUserInfoParam server_api_params.UserInfo
const SetSelfUserInfoCallback = constant.SuccessCallbackDefault



