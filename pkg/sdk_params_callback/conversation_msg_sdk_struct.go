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

const SetConversationStatusCallback = constant.SuccessCallbackDefault

type GetHistoryMessageListParams struct {
	UserID           string `json:"userID"`
	GroupID          string `json:"groupID"`
	StartClientMsgID string `json:"startClientMsgID"`
	Count            int    `json:"count"`
}
type GetHistoryMessageListCallback []*sdk_struct.MsgStruct

type RevokeMessageParams sdk_struct.MsgStruct

const RevokeMessageCallback = constant.SuccessCallbackDefault

const TypingStatusUpdateCallback = constant.SuccessCallbackDefault

type MarkC2CMessageAsReadParams []string

const MarkC2CMessageAsReadCallback = constant.SuccessCallbackDefault

const MarkGroupMessageHasRead = constant.SuccessCallbackDefault

type SetConversationStatusParams struct {
	UserId string `json:"userID" validate:"required"`
	Status int    `json:"status" validate:"required"`
}
type SearchLocalMessagesParams struct {
	SourceID             string   `json:"sourceID"`
	SessionType          int      `json:"sessionType"`
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
	ConversationID string                  `json:"conversationID"`
	MessageCount   int                     `json:"messageCount"`
	MessageList    []*sdk_struct.MsgStruct `json:"messageList"`
}
