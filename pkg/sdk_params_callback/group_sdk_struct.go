package sdk_params_callback

type SearchGroupsParam struct {
	KeywordList       []string `json:"keywordList"`
	IsSearchGroupID   bool     `json:"isSearchGroupID"`
	IsSearchGroupName bool     `json:"isSearchGroupName"`
}

type SearchGroupMembersParam struct {
	GroupID                string   `json:"groupID"`
	KeywordList            []string `json:"keywordList"`
	IsSearchUserID         bool     `json:"isSearchUserID"`
	IsSearchMemberNickname bool     `json:"isSearchMemberNickname"`
	Offset                 int      `json:"offset"`
	Count                  int      `json:"count"`
	PageNumber             int      `json:"pageNumber"`
}

type GetGroupApplicationListAsRecipientReq struct {
	GroupIDs      []string `json:"groupIDs"`
	HandleResults []int32  `json:"handleResults"`
	Offset        int32    `json:"offset"`
	Count         int32    `json:"count"`
}

type GetGroupApplicationListAsApplicantReq struct {
	GroupIDs      []string `json:"groupIDs"`
	HandleResults []int32  `json:"handleResults"`
	Offset        int32    `json:"offset"`
	Count         int32    `json:"count"`
}

type GetGroupApplicationUnhandledCountReq struct {
	Time int64 `json:"time"`
}
