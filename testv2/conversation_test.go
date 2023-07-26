// Copyright © 2023 OpenIM SDK. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testv2

import (
	"context"
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/sdk_struct"
	"testing"

	"github.com/OpenIMSDK/protocol/sdkws"
)

func Test_GetAllConversationList(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetAllConversationList(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

func Test_GetConversationListSplit(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetConversationListSplit(ctx, 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	for _, conversation := range conversations {
		t.Log(conversation)
	}
}

//func Test_SetConversationRecvMessageOpt(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Conversation().SetConversationRecvMessageOpt(ctx, []string{"asdasd"}, 1)
//	if err != nil {
//		t.Fatal(err)
//	}
//}

func Test_SetSetGlobalRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetGlobalRecvMessageOpt(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_HideConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().HideConversation(ctx, "asdasd")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetConversationRecvMessageOpt(t *testing.T) {
	opts, err := open_im_sdk.UserForSDK.Conversation().GetConversationRecvMessageOpt(ctx, []string{"asdasd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range opts {
		t.Log(v.ConversationID, *v.Result)
	}
}

func Test_GetGlobalRecvMessageOpt(t *testing.T) {
	opt, err := open_im_sdk.UserForSDK.Conversation().GetOneConversation(ctx, 2, "1772958501")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*opt)
}

func Test_GetGetMultipleConversation(t *testing.T) {
	conversations, err := open_im_sdk.UserForSDK.Conversation().GetMultipleConversation(ctx, []string{"asdasd"})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range conversations {
		t.Log(v)
	}
}

// funcation Test_DeleteConversation(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Conversation().DeleteConversation(ctx, "group_17729585012")
//	if err != nil {
//		if !strings.Contains(err.Error(), "no update") {
//			t.Fatal(err)
//		}
//	}
// }

func Test_DeleteAllConversationFromLocal(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllConversationFromLocal(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetConversationDraft(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetConversationDraft(ctx, "group_17729585012", "draft")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ResetConversationGroupAtType(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ResetConversationGroupAtType(ctx, "group_17729585012")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_PinConversation(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().PinConversation(ctx, "group_17729585012", true)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationPrivateChat(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationPrivateChat(ctx, "single_3411008330", true)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationBurnDuration(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationBurnDuration(ctx, "single_3411008330", 10)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SetOneConversationRecvMessageOpt(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().SetOneConversationRecvMessageOpt(ctx, "single_3411008330", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_GetTotalUnreadMsgCount(t *testing.T) {
	count, err := open_im_sdk.UserForSDK.Conversation().GetTotalUnreadMsgCount(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func Test_SendMessage(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessage(ctx, msg, "3411008330", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendMessageNotOss(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessageNotOss(ctx, msg, "3411008330", "", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendMessageByBuffer(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, _ := open_im_sdk.UserForSDK.Conversation().CreateTextMessage(ctx, "textMsg")
	_, err := open_im_sdk.UserForSDK.Conversation().SendMessageByBuffer(ctx, msg, "3411008330", "", &sdkws.OfflinePushInfo{}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FindMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().FindMessageList(ctx, []*sdk_params_callback.ConversationArgs{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(msgs.TotalCount)
	for _, v := range msgs.FindResultItems {
		t.Log(v)
	}
}

func Test_GetHistoryMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetHistoryMessageList(ctx, sdk_params_callback.GetHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageList(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageList(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_GetAdvancedHistoryMessageListReverse(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().GetAdvancedHistoryMessageListReverse(ctx, sdk_params_callback.GetAdvancedHistoryMessageListParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.MessageList {
		t.Log(v)
	}
}

func Test_InsertSingleMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertSingleMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "3411008330", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_InsertGroupMessageToLocalStorage(t *testing.T) {
	_, err := open_im_sdk.UserForSDK.Conversation().InsertGroupMessageToLocalStorage(ctx, &sdk_struct.MsgStruct{}, "group_17729585012", "")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SearchLocalMessages(t *testing.T) {
	msgs, err := open_im_sdk.UserForSDK.Conversation().SearchLocalMessages(ctx, &sdk_params_callback.SearchLocalMessagesParams{})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range msgs.SearchResultItems {
		t.Log(v)
	}
}

// // delete
// funcation Test_DeleteMessageFromLocalStorage(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Conversation().DeleteMessageFromLocalStorage(ctx, &sdk_struct.MsgStruct{SessionType: 1, ContentType: 1203,
//		ClientMsgID: "ef02943b05b02d02f92b0e92516099a3", Seq: 31, SendID: "kernaltestuid8", RecvID: "kernaltestuid9"})
//	if err != nil {
//		t.Fatal(err)
//	}
// }
//
// funcation Test_DeleteMessage(t *testing.T) {
//	err := open_im_sdk.UserForSDK.Conversation().DeleteMessage(ctx, &sdk_struct.MsgStruct{SessionType: 1, ContentType: 1203,
//		ClientMsgID: "ef02943b05b02d02f92b0e92516099a3", Seq: 31, SendID: "kernaltestuid8", RecvID: "kernaltestuid9"})
//	if err != nil {
//		t.Fatal(err)
//	}
// }

func Test_DeleteAllMessage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMessage(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_DeleteAllMessageFromLocalStorage(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().DeleteAllMessageFromLocalStorage(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_ClearConversationAndDeleteAllMsg(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().ClearConversationAndDeleteAllMsg(ctx, "si_3271407977_7152307910")
	if err != nil {
		t.Fatal(err)
	}
}

// func Test_RevokeMessage(t *testing.T) {
// 	err := open_im_sdk.UserForSDK.Conversation().RevokeMessage(ctx, &sdk_struct.MsgStruct{SessionType: 1, ContentType: 101,
// 		ClientMsgID: "380e2eb1709875340d769880982ebb21", Seq: 57, SendID: "9169012630", RecvID: "2456093263"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	time.Sleep(time.Second * 10)
// }

func Test_MarkConversationMessageAsRead(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().MarkConversationMessageAsRead(ctx, "si_2688118337_7249315132")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MarkMsgsAsRead(t *testing.T) {
	err := open_im_sdk.UserForSDK.Conversation().MarkMessagesAsReadByMsgID(ctx, "si_2688118337_7249315132",
		[]string{"fb56ed151b675e0837ed3af79dbf66b1",
			"635715c539be2e7812a0fc802f0cdc54", "1aba3fae3dc3f61c17e8eb09519cf8e1"})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendImgMsg(t *testing.T) {
	ctx = context.WithValue(ctx, "callback", TestSendMsg{})
	msg, err := open_im_sdk.UserForSDK.Conversation().CreateImageMessage(ctx, "C:\\Users\\Admin\\Desktop\\test.png")
	if err != nil {
		t.Fatal(err)
	}
	res, err := open_im_sdk.UserForSDK.Conversation().SendMessage(ctx, msg, "1919501984", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("send smg => %+v\n", res)
}
