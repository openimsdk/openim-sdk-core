package sdk

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
)

func (s *TestSDK) GetAllFriends(ctx context.Context) ([]*server_api_params.FullUserInfo, error) {
	res, err := s.SDK.Relation().GetFriendList(ctx, false)
	if err != nil {
		return nil, err
	}

	resp := []*server_api_params.FullUserInfo{}

	for _, v := range res {
		resp = append(resp, &server_api_params.FullUserInfo{FriendInfo: v})
	}

	return resp, nil
}
