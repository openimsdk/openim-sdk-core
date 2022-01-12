package sdk_params_callback

import (
	"open_im_sdk/pkg/constant"
	"open_im_sdk/pkg/db"
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
