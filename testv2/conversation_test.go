package testv2

import (
	"open_im_sdk/open_im_sdk"
	"testing"
)

func Test_GetAllConversationList(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_GetConversationListSplit(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetConversationListSplit(ctx, 0, 20)
	if err != nil {
		t.Error(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_SetConversationRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversationRecvMessageOpt(ctx, []string{"asdasd"}, 1)
	if err != nil {
		t.Error(err)
	}
}

func Test_SetSetGlobalRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetGlobalRecvMessageOpt(ctx, 1)
	if err != nil {
		t.Error(err)
	}
}

func Test_HideConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().HideConversation(ctx, "asdasd")
	if err != nil {
		t.Error(err)
	}
}
