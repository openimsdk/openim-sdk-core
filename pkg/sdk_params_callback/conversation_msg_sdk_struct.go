package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
	"open_im_sdk/sdk_struct"
)

//type GetAllConversationListParam null
type GetAllConversationListCallback []*db.LocalConversation

//type GetAllConversationListParam offset count
type GetConversationListSplitCallback []*db.LocalConversation

type SetConversationRecvMessageOptParams []string

const SetConversationRecvMessageOptCallback = constant.SuccessCallbackDefault

type GetConversationRecvMessageOptParams []string

type GetMultipleConversationParams []string
type GetMultipleConversationCallback []*db.LocalConversation

const DeleteConversationCallback = constant.SuccessCallbackDefault

const SetConversationDraftCallback = constant.SuccessCallbackDefault

const PinConversationDraftCallback = constant.SuccessCallbackDefault

type GetHistoryMessageListParams struct {
	UserID           string `json:"userID"`
	GroupID          string `json:"groupID"`
	StartClientMsgID string `json:"startClientMsgID"`
	Count            int    `json:"count"`
}
type GetHistoryMessageListCallback []*db.LocalChatLog

type RevokeMessageParams sdk_struct.MsgStruct

const RevokeMessageCallback = constant.SuccessCallbackDefault

const TypingStatusUpdateCallback = constant.SuccessCallbackDefault

type MarkC2CMessageAsReadParams []string

const MarkC2CMessageAsReadCallback = constant.SuccessCallbackDefault
