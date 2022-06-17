package server_api_params

type GetUsersInfoReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required"`
}
type GetUsersInfoResp struct {
	CommResp
	UserInfoList []*PublicUserInfo
	Data         []map[string]interface{} `json:"data"`
}

type UpdateSelfUserInfoReq struct {
	ApiUserInfo
	OperationID string `json:"operationID" binding:"required"`
}

type UpdateUserInfoResp struct {
	CommResp
}
type SetGlobalRecvMessageOptReq struct {
	OperationID      string `json:"operationID" binding:"required"`
	GlobalRecvMsgOpt *int32 `json:"globalRecvMsgOpt" binding:"omitempty,oneof=0 1 2"`
}
type SetGlobalRecvMessageOptResp struct {
	CommResp
}
type GetSelfUserInfoReq struct {
	OperationID string `json:"operationID" binding:"required"`
	UserID      string `json:"userID" binding:"required"`
}
type GetSelfUserInfoResp struct {
	CommResp
	UserInfo *UserInfo              `json:"-"`
	Data     map[string]interface{} `json:"data"`
}
