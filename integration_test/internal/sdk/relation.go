package sdk

import (
	"context"

	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/model_struct"
)

func (s *TestSDK) GetAllFriends(ctx context.Context) ([]*model_struct.LocalFriend, error) {
	res, err := s.SDK.Relation().GetFriendList(ctx, false)
	if err != nil {
		return nil, err
	}

	return res, nil

}
