// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.

package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/server_api_params"
)

//other user
type GetUsersInfoParam []string
type GetUsersInfoCallback []server_api_params.FullUserInfo

//type GetSelfUserInfoParam string
type GetSelfUserInfoCallback *model_struct.LocalUser

type SetSelfUserInfoParam server_api_params.ApiUserInfo

const SetSelfUserInfoCallback = constant.SuccessCallbackDefault
