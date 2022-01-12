package server_api_params

type PullUserMsgResp struct {
	ErrCode       int                       `json:"errCode"`
	ErrMsg        string                    `json:"errMsg"`
	ReqIdentifier int                       `json:"reqIdentifier"`
	MsgIncr       int                       `json:"msgIncr"`
	Data          paramsPullUserMsgDataResp `json:"data"`
}

type paramsPullUserMsgDataResp struct {
	Group  []*GatherFormat `json:"group"`
	MaxSeq int64           `json:"maxSeq"`
	MinSeq int64           `json:"minSeq"`
	Single []*GatherFormat `json:"single"`
}

type ArrMsg struct {
	SingleData []MsgData
	GroupData  []MsgData
}
