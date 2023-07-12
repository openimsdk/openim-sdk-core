// @Author BanTanger 2023/7/11 10:27
package funcation

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/sdk_params_callback"
	"open_im_sdk/pkg/utils"
)

func GetConversation(uid, conversationID, startMsgID string, count int) {
	var params sdk_params_callback.GetAdvancedHistoryMessageListParams
	params.UserID = AllLoginMgr[uid].UserID
	params.ConversationID = conversationID
	params.StartClientMsgID = startMsgID
	params.Count = count
	open_im_sdk.GetAdvancedHistoryMessageList(&testConversation, utils.OperationIDGenerator(), utils.StructToJsonString(params))
}
