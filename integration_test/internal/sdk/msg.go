package sdk

import (
	"context"
	"github.com/openimsdk/openim-sdk-core/v3/integration_test/internal/vars"
	"github.com/openimsdk/openim-sdk-core/v3/sdk_struct"
)

func (s *TestSDK) SendSingleMsg(ctx context.Context, msg *sdk_struct.MsgStruct, receiveID string) (*sdk_struct.MsgStruct, error) {
	vars.SendMsgCount.Add(1)
	res, err := s.SDK.Conversation().SendMessage(ctx, msg, receiveID, "", nil, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TestSDK) SendGroupMsg(ctx context.Context, msg *sdk_struct.MsgStruct, groupID string) (*sdk_struct.MsgStruct, error) {
	vars.SendMsgCount.Add(1)
	res, err := s.SDK.Conversation().SendMessage(ctx, msg, "", groupID, nil, false)
	if err != nil {
		return nil, err
	}
	return res, nil
}
