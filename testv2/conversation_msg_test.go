package testv2

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/sdk_struct"
	"testing"
)

func Test_CreateTextMessage(t *testing.T) {
	resp, err := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "test message")
	if err != nil {
		t.Error(err)
	}
	log.ZInfo(ctx, "CreateTextMessage success", "resp", resp)
}

func Test_CreateAdvancedTextMessage(t *testing.T) {
	m := []*sdk_struct.MessageEntity{
		{Type: "text",
			Offset: 1,
			Length: 2,
			Url:    "http://test.com",
			Info:   ""},
		{Type: "text2",
			Offset: 1,
			Length: 2,
			Url:    "http://test.com",
			Info:   ""},
	}
	resp, err := open_im_sdk.UserForSDK.Conversation().CreateAdvancedTextMessage(ctx, "test message", m)
	if err != nil {
		t.Error(err)
	}
	log.ZInfo(ctx, "CreateTextMessage success", "resp", resp)
}

func Test_CreateTextAtMessage(t *testing.T) {
	resp, err := open_im_sdk.UserForSDK.Conversation().CreateTextAtMessage(ctx, "test message", []string{"test"}, []*sdk_struct.AtInfo{}, &sdk_struct.MsgStruct{})
	if err != nil {
		t.Error(err)
	}
	log.ZInfo(ctx, "CreateTextAtMessage success", "resp", resp)
}
