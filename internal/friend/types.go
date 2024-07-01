package friend

import "github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"

type GetFriendInfoListV2 struct {
	FullUserInfoList []*server_api_params.FullUserInfo
	IsEnd            bool `json:"isEnd"`
}
