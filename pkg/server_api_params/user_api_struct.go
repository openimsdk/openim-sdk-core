package server_api_params

type GetUserInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUserInfoResp struct {
	CommResp
	UserInfoList []*UserInfo
	Data         []map[string]interface{} `json:"data"`
}

type UpdateUserInfoReq struct {
	UserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UpdateUserInfoResp struct {
	CommResp
}
