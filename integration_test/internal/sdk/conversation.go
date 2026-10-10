package sdk

import (
	"context"
)

func (s *TestSDK) GetTotalConversationCount(ctx context.Context) (int, error) {
	res, err := s.SDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		return 0, err
	}
	return len(res), nil
}
