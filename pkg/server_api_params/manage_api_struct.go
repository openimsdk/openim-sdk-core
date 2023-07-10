package server_api_params

type DeleteUsersReq struct {
	OperationID      string   `json:"operationID" binding:"required"`
	DeleteUserIDList []string `json:"deleteUserIDList" binding:"required"`
}
type DeleteUsersResp struct {
	CommResp
	FailedUserIDList []string `json:"data"`
}
type GetAllUsersUidReq struct {
	OperationID string `json:"operationID" binding:"required"`
}
type GetAllUsersUidResp struct {
	CommResp
	UserIDList []string `json:"data"`
}
type GetUsersOnlineStatusReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type GetUsersOnlineStatusResp struct {
	CommResp
	SuccessResult []GetusersonlinestatusrespSuccessresult `json:"data"`
}
type AccountCheckReq struct {
	OperationID     string   `json:"operationID" binding:"required"`
	CheckUserIDList []string `json:"checkUserIDList" binding:"required,lte=100"`
}
type AccountCheckResp struct {
	CommResp
	ResultList []*AccountCheckResp_SingleUserStatus `json:"data"`
}
type AccountCheckResp_SingleUserStatus struct {
	UserID        string `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	AccountStatus string `protobuf:"bytes,2,opt,name=accountStatus" json:"accountStatus,omitempty"`
}

type GetusersonlinestatusrespSuccessdetail struct {
	Platform             string   `protobuf:"bytes,1,opt,name=platform" json:"platform,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
type GetusersonlinestatusrespSuccessresult struct {
	UserID               string                                   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	Status               string                                   `protobuf:"bytes,2,opt,name=status" json:"status,omitempty"`
	DetailPlatformStatus []*GetusersonlinestatusrespSuccessdetail `protobuf:"bytes,3,rep,name=detailPlatformStatus" json:"detailPlatformStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                 `json:"-"`
	XXX_unrecognized     []byte                                   `json:"-"`
	XXX_sizecache        int32                                    `json:"-"`
}
type AccountcheckrespSingleuserstatus struct {
	UserID               string   `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	AccountStatus        string   `protobuf:"bytes,2,opt,name=accountStatus" json:"accountStatus,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
