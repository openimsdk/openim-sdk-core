package server_api_params

import "open_im_sdk/open_im_sdk"

type GetUserInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUserInfoResp struct {
	CommResp
	UserInfoList []*open_im_sdk.UserInfo
	Data         []map[string]interface{} `json:"data"`
}

type UpdateUserInfoReq struct {
	open_im_sdk.UserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UpdateUserInfoResp struct {
	CommResp
}
