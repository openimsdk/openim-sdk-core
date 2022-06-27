package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/sdk_struct"
)

//type GetAllConversationListParam null
type GetAllConversationListCallback []*model_struct.LocalConversation

//type GetAllConversationListParam offset count
type GetConversationListSplitCallback []*model_struct.LocalConversation

type SetConversationRecvMessageOptParams []string

const SetConversationRecvMessageOptCallback = constant.SuccessCallbackDefault
const SetGlobalRecvMessageOptCallback = constant.SuccessCallbackDefault

type GetConversationRecvMessageOptParams []string

type GetMultipleConversationParams []string
type GetMultipleConversationCallback []*model_struct.LocalConversation

const DeleteConversationCallback = constant.SuccessCallbackDefault
const DeleteAllConversationFromLocalCallback = constant.SuccessCallbackDefault

const SetConversationDraftCallback = constant.SuccessCallbackDefault
const ResetConversationGroupAtTypeCallback = constant.SuccessCallbackDefault

const PinConversationDraftCallback = constant.SuccessCallbackDefault

const SetConversationMessageOptCallback = constant.SuccessCallbackDefault

const SetConversationPrivateChatOptCallback = constant.SuccessCallbackDefault

type GetHistoryMessageListParams struct {
	UserID           string `json:"userID"`
	GroupID          string `json:"groupID"`
	ConversationID   string `json:"conversationID"`
	StartClientMsgID string `json:"startClientMsgID"`
	Count            int    `json:"count"`
}
type GetHistoryMessageListCallback []*sdk_struct.MsgStruct

type RevokeMessageParams sdk_struct.MsgStruct

const RevokeMessageCallback = constant.SuccessCallbackDefault

const TypingStatusUpdateCallback = constant.SuccessCallbackDefault

type MarkC2CMessageAsReadParams []string

const MarkC2CMessageAsReadCallback = constant.SuccessCallbackDefault

const MarkGroupMessageHasReadCallback = constant.SuccessCallbackDefault

type MarkGroupMessageAsReadParams []string

const MarkGroupMessageAsReadCallback = constant.SuccessCallbackDefault

type MarkMessageAsReadByConIDParams []string

const MarkMessageAsReadByConIDCallback = constant.SuccessCallbackDefault

type SetConversationStatusParams struct {
	UserId string `json:"userID" validate:"required"`
	Status int    `json:"status" validate:"required"`
}
type SearchLocalMessagesParams struct {
	ConversationID       string   `json:"conversationID"`
	KeywordList          []string `json:"keywordList"`
	KeywordListMatchType int      `json:"keywordListMatchType"`
	SenderUserIDList     []string `json:"senderUserIDList"`
	MessageTypeList      []int    `json:"messageTypeList"`
	SearchTimePosition   int64    `json:"searchTimePosition"`
	SearchTimePeriod     int64    `json:"searchTimePeriod"`
	PageIndex            int      `json:"pageIndex"`
	Count                int      `json:"count"`
}
type SearchLocalMessagesCallback struct {
	TotalCount        int                           `json:"totalCount"`
	SearchResultItems []*SearchByConversationResult `json:"searchResultItems"`
}
type SearchByConversationResult struct {
	ConversationID   string                  `json:"conversationID"`
	ConversationType int32                   `json:"conversationType"`
	ShowName         string                  `json:"showName"`
	FaceURL          string                  `json:"faceURL"`
	MessageCount     int                     `json:"messageCount"`
	MessageList      []*sdk_struct.MsgStruct `json:"messageList"`
}
