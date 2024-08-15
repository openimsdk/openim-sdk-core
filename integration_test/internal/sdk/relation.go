package sdk

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/server_api_params"
)

func (s *TestSDK) GetAllFriends(ctx context.Context) ([]*server_api_params.FullUserInfo, error) {
	res, err := s.SDK.Friend().GetFriendList(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}
