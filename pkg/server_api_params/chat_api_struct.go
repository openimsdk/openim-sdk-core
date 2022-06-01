package server_api_params

type DeleteMsgReq struct {
	OpUserID    string   `json:"opUserID"`
	UserID      string   `json:"userID"`
	SeqList     []uint32 `json:"seqList"`
	OperationID string   `json:"operationID"`
}

type DeleteMsgResp struct {
}

type CleanUpMsgReq struct {
	UserID      string `json:"userID"  binding:"required"`
	OperationID string `json:"operationID"  binding:"required"`
}

type CleanUpMsgResp struct {
	CommResp
}
